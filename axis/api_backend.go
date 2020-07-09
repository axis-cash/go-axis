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

package axis

import (
	"context"
	"errors"
	"math/big"

	"github.com/axis-cash/go-axis/zero/txtool/flight"

	"github.com/axis-cash/go-axis/zero/txtool"
	"github.com/axis-cash/go-axis/zero/txtool/prepare"

	"github.com/axis-cash/go-axis/zero/wallet/exchange"

	"github.com/axis-cash/go-axis/log"

	"github.com/axis-cash/go-axis-import/c_type"

	"github.com/axis-cash/go-axis/consensus"
	"github.com/axis-cash/go-axis/miner"

	"github.com/axis-cash/go-axis/accounts"
	"github.com/axis-cash/go-axis/common"
	"github.com/axis-cash/go-axis/core"
	"github.com/axis-cash/go-axis/core/bloombits"
	"github.com/axis-cash/go-axis/core/rawdb"
	"github.com/axis-cash/go-axis/core/state"
	"github.com/axis-cash/go-axis/core/types"
	"github.com/axis-cash/go-axis/core/vm"
	"github.com/axis-cash/go-axis/event"
	"github.com/axis-cash/go-axis/params"
	"github.com/axis-cash/go-axis/rpc"
	"github.com/axis-cash/go-axis/axis/downloader"
	"github.com/axis-cash/go-axis/axis/gasprice"
	"github.com/axis-cash/go-axis/axisdb"
	"github.com/axis-cash/go-axis/zero/wallet/light"
)

// AxisAPIBackend implements ethapi.Backend for full nodes
type AxisAPIBackend struct {
	axis *Axis
	gpo  *gasprice.Oracle
}

// ChainConfig returns the active chain configuration.
func (b *AxisAPIBackend) ChainConfig() *params.ChainConfig {
	return b.axis.chainConfig
}

func (b *AxisAPIBackend) CurrentBlock() *types.Block {
	return b.axis.blockchain.CurrentBlock()
}

func (b *AxisAPIBackend) GetEngin() consensus.Engine {
	return b.axis.engine
}

func (b *AxisAPIBackend) GetMiner() *miner.Miner {
	return b.axis.miner
}

func (b *AxisAPIBackend) SetHead(number uint64) {
	b.axis.protocolManager.downloader.Cancel()
	b.axis.blockchain.SetHead(number, core.DelFn)
}

func (b *AxisAPIBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.axis.miner.PendingBlock()
		return block.Header(), nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.axis.blockchain.CurrentBlock().Header(), nil
	}
	return b.axis.blockchain.GetHeaderByNumber(uint64(blockNr)), nil
}

func (b *AxisAPIBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	return b.axis.blockchain.GetHeaderByHash(hash), nil
}

func (b *AxisAPIBackend) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.axis.miner.PendingBlock()
		return block, nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.axis.blockchain.CurrentBlock(), nil
	}
	return b.axis.blockchain.GetBlockByNumber(uint64(blockNr)), nil
}

func (b *AxisAPIBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	// Pending state is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block, state := b.axis.miner.Pending()
		return state, block.Header(), nil
	}
	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, nil, err
	}
	stateDb, err := b.axis.BlockChain().StateAt(header)
	return stateDb, header, err
}

func (b *AxisAPIBackend) GetBlock(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return b.axis.blockchain.GetBlockByHash(hash), nil
}

func (b *AxisAPIBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	if number := rawdb.ReadHeaderNumber(b.axis.chainDb, hash); number != nil {
		return rawdb.ReadReceipts(b.axis.chainDb, hash, *number), nil
	}
	return nil, nil
}

func (b *AxisAPIBackend) GetLogs(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	number := rawdb.ReadHeaderNumber(b.axis.chainDb, hash)
	if number == nil {
		return nil, nil
	}
	receipts := rawdb.ReadReceipts(b.axis.chainDb, hash, *number)
	if receipts == nil {
		return nil, nil
	}
	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

func (b *AxisAPIBackend) GetTd(blockHash common.Hash) *big.Int {
	return b.axis.blockchain.GetTdByHash(blockHash)
}

func (b *AxisAPIBackend) GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error) {
	vmError := func() error { return nil }

	context := core.NewEVMContext(msg, header, b.axis.BlockChain(), nil)
	return vm.NewEVM(context, state, b.axis.chainConfig, vmCfg), vmError, nil
}

func (b *AxisAPIBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.axis.BlockChain().SubscribeRemovedLogsEvent(ch)
}

func (b *AxisAPIBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.axis.BlockChain().SubscribeChainEvent(ch)
}

func (b *AxisAPIBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return b.axis.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *AxisAPIBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return b.axis.BlockChain().SubscribeChainSideEvent(ch)
}

func (b *AxisAPIBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.axis.BlockChain().SubscribeLogsEvent(ch)
}

func (b *AxisAPIBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return b.axis.txPool.AddLocal(signedTx)
}

func (b *AxisAPIBackend) GetPoolTransactions() (types.Transactions, error) {
	pending, err := b.axis.txPool.Pending()
	if err != nil {
		return nil, err
	}

	return pending, nil
}

func (b *AxisAPIBackend) GetPoolTransaction(hash common.Hash) *types.Transaction {
	return b.axis.txPool.Get(hash)
}

//func (b *AxisAPIBackend) GetPoolNonce(ctx context.Context, addr common.Data) (uint64, error) {
//	return b.axis.txPool.State().GetNonce(addr), nil
//}

func (b *AxisAPIBackend) Stats() (pending int, queued int) {
	return b.axis.txPool.Stats()
}

func (b *AxisAPIBackend) TxPoolContent() (types.Transactions, types.Transactions) {
	return b.axis.TxPool().Content()
}

func (b *AxisAPIBackend) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return b.axis.TxPool().SubscribeNewTxsEvent(ch)
}

func (b *AxisAPIBackend) Downloader() *downloader.Downloader {
	return b.axis.Downloader()
}

func (b *AxisAPIBackend) ProtocolVersion() int {
	return b.axis.EthVersion()
}

func (b *AxisAPIBackend) PeerCount() uint {
	return uint(b.axis.netRPCService.PeerCount())
}

func (b *AxisAPIBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

func (b *AxisAPIBackend) ChainDb() axisdb.Database {
	return b.axis.ChainDb()
}

func (b *AxisAPIBackend) EventMux() *event.TypeMux {
	return b.axis.EventMux()
}

func (b *AxisAPIBackend) AccountManager() *accounts.Manager {
	return b.axis.AccountManager()
}

func (b *AxisAPIBackend) BloomStatus() (uint64, uint64) {
	sections, _, _ := b.axis.bloomIndexer.Sections()
	return params.BloomBitsBlocks, sections
}

func (b *AxisAPIBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.axis.bloomRequests)
	}
}

func (b *AxisAPIBackend) GetBlocksInfo(start uint64, count uint64) ([]txtool.Block, error) {
	return flight.SRI_Inst.GetBlocksInfo(start, count)

}
func (b *AxisAPIBackend) GetAnchor(roots []c_type.Uint256) ([]txtool.Witness, error) {
	return flight.SRI_Inst.GetAnchor(roots)

}
func (b *AxisAPIBackend) CommitTx(tx *txtool.GTx) error {
	gasPrice := big.Int(tx.GasPrice)
	gas := uint64(tx.Gas)
	signedTx := types.NewTxWithGTx(gas, &gasPrice, &tx.Tx)
	log.Info("commitTx", "txhash", signedTx.Hash().String())
	return b.axis.txPool.AddLocal(signedTx)
}

func (b *AxisAPIBackend) GetPkNumber(pk c_type.Uint512) (number uint64, e error) {
	if b.axis.exchange == nil {
		e = errors.New("not start exchange")
		return
	}
	return b.axis.exchange.GetCurrencyNumber(pk), nil
}

func (b *AxisAPIBackend) GetPkr(pk *c_type.Uint512, index *c_type.Uint256) (pkr c_type.PKr, e error) {
	if b.axis.exchange == nil {
		e = errors.New("not start exchange")
		return
	}
	return b.axis.exchange.GetPkr(pk, index)
}

func (b *AxisAPIBackend) GetPkrEx(pk *c_type.Uint512, index *c_type.Uint256) (pkr c_type.PKrEx, e error) {
	if b.axis.exchange == nil {
		e = errors.New("not start exchange")
		return
	}
	return b.axis.exchange.GetPkrEx(pk, index)
}

func (b *AxisAPIBackend) GetLockedBalances(pk c_type.Uint512) (balances map[string]*big.Int) {
	if b.axis.exchange == nil {
		return
	}
	return b.axis.exchange.GetLockedBalances(pk)
}

func (b *AxisAPIBackend) GetMaxAvailable(pk c_type.Uint512, currency string) (amount *big.Int) {
	if b.axis.exchange == nil {
		return
	}
	return b.axis.exchange.GetMaxAvailable(pk, currency)
}

func (b *AxisAPIBackend) GetBalances(pk c_type.Uint512) (balances map[string]*big.Int, tickets map[string][]*common.Hash) {
	if b.axis.exchange == nil {
		return
	}
	return b.axis.exchange.GetBalances(pk)
}

func (b *AxisAPIBackend) GenTx(param prepare.PreTxParam) (txParam *txtool.GTxParam, e error) {
	if b.axis.exchange == nil {
		e = errors.New("not start exchange")
		return
	}
	return b.axis.exchange.GenTx(param)
}

func (b *AxisAPIBackend) GetRecordsByPkr(pkr c_type.PKr, begin, end uint64) (records []exchange.Utxo, err error) {
	if b.axis.exchange == nil {
		err = errors.New("not start exchange")
		return
	}
	return b.axis.exchange.GetRecordsByPkr(pkr, begin, end)
}

func (b *AxisAPIBackend) GetRecordsByPk(pk *c_type.Uint512, begin, end uint64) (records []exchange.Utxo, err error) {
	if b.axis.exchange == nil {
		err = errors.New("not start exchange")
		return
	}
	return b.axis.exchange.GetRecordsByPk(pk, begin, end)
}

func (b *AxisAPIBackend) GetRecordsByTxHash(txHash c_type.Uint256) (records []exchange.Utxo, err error) {
	if b.axis.exchange == nil {
		err = errors.New("not start exchange")
		return
	}
	return b.axis.exchange.GetRecordsByTxHash(txHash)
}

func (b *AxisAPIBackend) GetOutByPKr(pkrs []c_type.PKr, start, end uint64) (br light.BlockOutResp, e error) {
	if b.axis.lightNode == nil {
		e = errors.New("not start light")
		return
	}
	return b.axis.lightNode.GetOutsByPKr(pkrs, start, end)
}

func (b *AxisAPIBackend) CheckNil(Nils []c_type.Uint256) (nilResps []light.NilValue, e error) {
	if b.axis.lightNode == nil {
		e = errors.New("not start light")
		return
	}
	return b.axis.lightNode.CheckNil(Nils)
}
