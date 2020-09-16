// copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package axis implements the Axis protocol.
package axis

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/axis-cash/go-axis/zero/wallet/stakeservice"

	"github.com/axis-cash/go-axis-import/c_type"

	"github.com/axis-cash/go-axis-import/superzk"

	"github.com/axis-cash/go-axis/common/address"

	"github.com/axis-cash/go-axis/voter"
	"github.com/axis-cash/go-axis/zero/txtool"
	"github.com/axis-cash/go-axis/zero/zconfig"

	"github.com/axis-cash/go-axis/internal/ethapi"
	"github.com/axis-cash/go-axis/zero/wallet/exchange"

	"github.com/axis-cash/go-axis/accounts"
	"github.com/axis-cash/go-axis/common"
	"github.com/axis-cash/go-axis/common/hexutil"
	"github.com/axis-cash/go-axis/consensus"
	"github.com/axis-cash/go-axis/consensus/ethash"
	"github.com/axis-cash/go-axis/core"
	"github.com/axis-cash/go-axis/core/bloombits"
	"github.com/axis-cash/go-axis/core/rawdb"
	"github.com/axis-cash/go-axis/core/types"
	"github.com/axis-cash/go-axis/core/vm"
	"github.com/axis-cash/go-axis/event"
	"github.com/axis-cash/go-axis/log"
	"github.com/axis-cash/go-axis/miner"
	"github.com/axis-cash/go-axis/node"
	"github.com/axis-cash/go-axis/p2p"
	"github.com/axis-cash/go-axis/params"
	"github.com/axis-cash/go-axis/rlp"
	"github.com/axis-cash/go-axis/rpc"
	"github.com/axis-cash/go-axis/axis/downloader"
	"github.com/axis-cash/go-axis/axis/filters"
	"github.com/axis-cash/go-axis/axis/gasprice"
	"github.com/axis-cash/go-axis/axisdb"
	"github.com/axis-cash/go-axis/zero/wallet/light"
)

type LesServer interface {
	Start(srvr *p2p.Server)
	Stop()
	Protocols() []p2p.Protocol
	SetBloomBitsIndexer(bbIndexer *core.ChainIndexer)
}

// Axis implements the Axis full node service.
type Axis struct {
	config      *Config
	chainConfig *params.ChainConfig

	// Channel for shutting down the service
	shutdownChan chan bool // Channel for shutting down the Axis

	// Handlers
	txPool          *core.TxPool
	voter           *voter.Voter
	blockchain      *core.BlockChain
	exchange        *exchange.Exchange
	lightNode       *light.LightNode
	protocolManager *ProtocolManager
	lesServer       LesServer

	// DB interfaces
	chainDb axisdb.Database // Block chain database

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *core.ChainIndexer             // Bloom indexer operating during block imports

	APIBackend *AxisAPIBackend

	miner    *miner.Miner
	gasPrice *big.Int
	axisbase accounts.Account

	networkID     uint64
	netRPCService *ethapi.PublicNetAPI

	lock sync.RWMutex // Protects the variadic fields (s.g. gas price and axisbase)
}

func (s *Axis) AddLesServer(ls LesServer) {
	s.lesServer = ls
	ls.SetBloomBitsIndexer(s.bloomIndexer)
}

var AxisInstance *Axis

// New creates a new Axis object (including the
// initialisation of the common Axis object)
func New(ctx *node.ServiceContext, config *Config) (*Axis, error) {
	if config.SyncMode == downloader.LightSync {
		return nil, errors.New("can't run axis.Axis in light sync mode, use les.LightEthereum")
	}
	if !config.SyncMode.IsValid() {
		return nil, fmt.Errorf("invalid sync mode %d", config.SyncMode)
	}
	chainDb, err := CreateDB(ctx, config, "chaindata")
	if err != nil {
		return nil, err
	}
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlock(chainDb, config.Genesis)
	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {
		return nil, genesisErr
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	axis := &Axis{
		config:         config,
		chainDb:        chainDb,
		chainConfig:    chainConfig,
		eventMux:       ctx.EventMux,
		accountManager: ctx.AccountManager,
		engine:         CreateConsensusEngine(ctx, &config.Ethash, chainConfig, chainDb),
		shutdownChan:   make(chan bool),
		networkID:      config.NetworkId,
		gasPrice:       config.GasPrice,
		bloomRequests:  make(chan chan *bloombits.Retrieval),
		bloomIndexer:   NewBloomIndexer(chainDb, params.BloomBitsBlocks),
	}

	log.Info("Initialising Axis protocol", "versions", ProtocolVersions, "network", config.NetworkId)

	if !config.SkipBcVersionCheck {
		bcVersion := rawdb.ReadDatabaseVersion(chainDb)
		if bcVersion != core.BlockChainVersion && bcVersion != 0 {
			return nil, fmt.Errorf("Blockchain DB version mismatch (%d / %d). Run gaxis upgradedb.\n", bcVersion, core.BlockChainVersion)
		}
		rawdb.WriteDatabaseVersion(chainDb, core.BlockChainVersion)
	}
	var (
		vmConfig    = vm.Config{EnablePreimageRecording: config.EnablePreimageRecording}
		cacheConfig = &core.CacheConfig{Disabled: config.NoPruning, TrieNodeLimit: config.TrieCache, TrieTimeLimit: config.TrieTimeout}
	)
	axis.blockchain, err = core.NewBlockChain(chainDb, cacheConfig, axis.chainConfig, axis.engine, vmConfig, axis.accountManager)

	txtool.Ref_inst.SetBC(&core.State1BlockChain{axis.blockchain})

	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		axis.blockchain.SetHead(compat.RewindTo, core.DelFn)
		rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}
	axis.bloomIndexer.Start(axis.blockchain)

	// if config.TxPool.Journal != "" {
	//	config.TxPool.Journal = ctx.ResolvePath(config.TxPool.Journal)
	// }

	config.TxPool.StartLight = config.StartLight

	axis.txPool = core.NewTxPool(config.TxPool, axis.chainConfig, axis.blockchain)

	axis.voter = voter.NewVoter(axis.chainConfig, axis.blockchain, axis)

	if axis.protocolManager, err = NewProtocolManager(axis.chainConfig, config.SyncMode, config.NetworkId, axis.eventMux, axis.voter, axis.txPool, axis.engine, axis.blockchain, chainDb); err != nil {
		return nil, err
	}
	axis.miner = miner.New(axis, axis.chainConfig, axis.EventMux(), axis.voter, axis.engine)
	axis.miner.SetExtra(makeExtraData(config.ExtraData))

	axis.APIBackend = &AxisAPIBackend{axis, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.GasPrice
	}
	axis.APIBackend.gpo = gasprice.NewOracle(axis.APIBackend, gpoParams)

	ethapi.Backend_Instance = axis.APIBackend

	// init exchange
	if config.StartExchange {
		axis.exchange = exchange.NewExchange(zconfig.Exchange_dir(), axis.txPool, axis.accountManager, config.AutoMerge)
	}

	if config.StartStake {
		stakeservice.NewStakeService(zconfig.Stake_dir(), axis.blockchain, axis.accountManager)
	}

	// init light
	if config.StartLight {
		axis.lightNode = light.NewLightNode(zconfig.Light_dir(), axis.txPool, axis.blockchain.GetDB())
	}

	// if config.Proof != nil {
	// 	if config.Proof.PKr == (c_type.PKr{}) {
	// 		wallets := axis.accountManager.Wallets()
	// 		if len(wallets) == 0 {
	// 			// panic("init proofService error")
	// 		}
	//
	// 		account := wallets[0].Accounts()
	// 		config.Proof.PKr = superzk.Pk2PKr(account[0].Address.ToUint512(), &c_type.Uint256{1})
	// 	}
	// 	proofservice.NewProofService("", axis.APIBackend, config.Proof);
	// }

	AxisInstance = axis
	return axis, nil
}

func makeExtraData(extra []byte) []byte {
	if len(extra) == 0 {
		// create default extradata
		extra, _ = rlp.EncodeToBytes([]interface{}{
			uint(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch),
			"gaxis",
			runtime.Version(),
			runtime.GOOS,
		})
	}
	if uint64(len(extra)) > params.MaximumExtraDataSize {
		log.Warn("Miner extra data exceed limit", "extra", hexutil.Bytes(extra), "limit", params.MaximumExtraDataSize)
		extra = nil
	}
	return extra
}

// CreateDB creates the chain database.
func CreateDB(ctx *node.ServiceContext, config *Config, name string) (axisdb.Database, error) {
	db, err := ctx.OpenDatabase(name, config.DatabaseCache, config.DatabaseHandles)
	if err != nil {
		return nil, err
	}
	if db, ok := db.(*axisdb.LDBDatabase); ok {
		db.Meter("axis/db/chaindata/")
	}
	return db, nil
}

// CreateConsensusEngine creates the required type of consensus engine instance for an Axis service
func CreateConsensusEngine(ctx *node.ServiceContext, config *ethash.Config, chainConfig *params.ChainConfig, db axisdb.Database) consensus.Engine { // If proof-of-authority is requested, set it up
	// Otherwise assume proof-of-work
	switch config.PowMode {
	case ethash.ModeFake:
		log.Warn("Ethash used in fake mode")
		return ethash.NewFaker()
	case ethash.ModeTest:
		log.Warn("Ethash used in test mode")
		return ethash.NewTester()
	case ethash.ModeShared:
		log.Warn("Ethash used in shared mode")
		return ethash.NewShared()
	default:
		engine := ethash.New(ethash.Config{
			CacheDir:       ctx.ResolvePath(config.CacheDir),
			CachesInMem:    config.CachesInMem,
			CachesOnDisk:   config.CachesOnDisk,
			DatasetDir:     config.DatasetDir,
			DatasetsInMem:  config.DatasetsInMem,
			DatasetsOnDisk: config.DatasetsOnDisk,
		})
		engine.SetThreads(-1) // Disable CPU mining
		return engine
	}
}

// APIs return the collection of RPC services the ethereum package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *Axis) APIs() []rpc.API {
	apis := ethapi.GetAPIs(s.APIBackend)

	// Append any APIs exposed explicitly by the consensus engine
	apis = append(apis, s.engine.APIs(s.BlockChain())...)

	// Append all the local APIs and return
	return append(apis, []rpc.API{
		{
			Namespace: "axis",
			Version:   "1.0",
			Service:   NewPublicAxisAPI(s),
			Public:    true,
		}, {
			Namespace: "axis",
			Version:   "1.0",
			Service:   NewPublicMinerAPI(s),
			Public:    true,
		}, {
			Namespace: "axis",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "miner",
			Version:   "1.0",
			Service:   NewPrivateMinerAPI(s),
			Public:    false,
		}, {
			Namespace: "axis",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.APIBackend, false),
			Public:    true,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateAdminAPI(s),
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(s),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(s.chainConfig, s),
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
}

func (s *Axis) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *Axis) Axisbase() (eb accounts.Account, err error) {
	s.lock.RLock()
	axisbase := s.axisbase
	s.lock.RUnlock()

	if axisbase != (accounts.Account{}) {
		return axisbase, nil
	}
	if wallets := s.AccountManager().Wallets(); len(wallets) > 0 {
		if accounts := wallets[0].Accounts(); len(accounts) > 0 {
			axisbase := accounts[0]

			s.lock.Lock()
			s.axisbase = axisbase
			s.lock.Unlock()

			log.Info("Axisbase automatically configured", "address", axisbase)
			return axisbase, nil
		}
	}
	return accounts.Account{}, fmt.Errorf("Axisbase must be explicitly specified")
}

// SetAxisbase sets the mining reward address.
func (s *Axis) SetAxisbase(axisbase address.MixBase58Adrress) {
	s.lock.Lock()
	account, _ := s.accountManager.FindAccountByPkr(axisbase.ToPkr())
	s.axisbase = account
	s.lock.Unlock()

	s.miner.SetAxisbase(account)
}

func (s *Axis) StartMining(local bool) error {
	eb, err := s.Axisbase()
	if err != nil {
		log.Error("Cannot start mining without axisbase", "err", err)
		return fmt.Errorf("axisbase missing: %v", err)
	}
	current_height := s.blockchain.CurrentHeader().Number.Uint64()
	pkr, lic, ret := superzk.Pk2PKrAndLICr(eb.Address.ToUint512().NewRef(), current_height)
	ret = superzk.CheckLICr(&pkr, &lic, current_height)
	if !ret {
		lic_t := c_type.LICr{}
		if bytes.Equal(lic.Proof[:32], lic_t.Proof[:32]) {
			log.Error("Cannot start mining , miner license does not exists, or coinBase is error", "", "")
			return fmt.Errorf(" miner license does not exists")
		} else {
			log.Error("Cannot start mining ,invalid miner license, or coinBase is error", "", "")
			return fmt.Errorf("invalid miner license: %v", common.Bytes2Hex(lic.Proof[:]))
		}
	}

	if local {
		// If local (CPU) mining is started, we can disable the transaction rejection
		// mechanism introduced to speed sync times. CPU mining on mainnet is ludicrous
		// so none will ever hit this path, whereas marking sync done on CPU mining
		// will ensure that private networks work in single miner mode too.
		atomic.StoreUint32(&s.protocolManager.acceptTxs, 1)
	}
	go s.miner.Start(eb)
	return nil
}

func (s *Axis) StopMining()         { s.miner.Stop() }
func (s *Axis) IsMining() bool      { return s.miner.Mining() }
func (s *Axis) Miner() *miner.Miner { return s.miner }

func (s *Axis) AccountManager() *accounts.Manager  { return s.accountManager }
func (s *Axis) BlockChain() *core.BlockChain       { return s.blockchain }
func (s *Axis) TxPool() *core.TxPool               { return s.txPool }
func (s *Axis) Voter() *voter.Voter                { return s.voter }
func (s *Axis) EventMux() *event.TypeMux           { return s.eventMux }
func (s *Axis) Engine() consensus.Engine           { return s.engine }
func (s *Axis) ChainDb() axisdb.Database           { return s.chainDb }
func (s *Axis) IsListening() bool                  { return true } // Always listening
func (s *Axis) EthVersion() int                    { return int(s.protocolManager.SubProtocols[0].Version) }
func (s *Axis) NetVersion() uint64                 { return s.networkID }
func (s *Axis) Downloader() *downloader.Downloader { return s.protocolManager.downloader }

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *Axis) Protocols() []p2p.Protocol {
	if s.lesServer == nil {
		return s.protocolManager.SubProtocols
	}
	return append(s.protocolManager.SubProtocols, s.lesServer.Protocols()...)
}

// Start implements node.Service, starting all internal goroutines needed by the
// Axis protocol implementation.
func (s *Axis) Start(srvr *p2p.Server) error {
	// Start the bloom bits servicing goroutines
	s.startBloomHandlers()

	// Start the RPC service
	s.netRPCService = ethapi.NewPublicNetAPI(srvr, s.NetVersion())

	// Figure out a max peers count based on the server limits
	maxPeers := srvr.MaxPeers
	if s.config.LightServ > 0 {
		if s.config.LightPeers >= srvr.MaxPeers {
			return fmt.Errorf("invalid peer config: light peer count (%d) >= total peer count (%d)", s.config.LightPeers, srvr.MaxPeers)
		}
		maxPeers -= s.config.LightPeers
	}
	// Start the networking layer and the light server if requested
	s.protocolManager.Start(maxPeers)
	if s.lesServer != nil {
		s.lesServer.Start(srvr)
	}
	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Axis protocol.
func (s *Axis) Stop() error {
	s.bloomIndexer.Close()
	s.blockchain.Stop()
	s.protocolManager.Stop()
	if s.lesServer != nil {
		s.lesServer.Stop()
	}
	s.txPool.Stop()
	s.miner.Stop()
	s.eventMux.Stop()

	s.chainDb.Close()
	close(s.shutdownChan)

	return nil
}
