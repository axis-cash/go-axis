// Copyright 2017 The go-ethereum Authors
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

package ethash

import (
	"bytes"
	"errors"
	"fmt"
	// "github.com/axis-cash/go-axis/zero/stake"
	"math/big"
	"runtime"
	"time"

	"github.com/axis-cash/go-axis/zero/zconfig"

	"github.com/axis-cash/go-axis-import/axisparam"
	"github.com/axis-cash/go-axis-import/superzk"

	//"github.com/axis-cash/go-axis/crypto"

	"github.com/axis-cash/go-axis/zero/txs/assets"
	"github.com/axis-cash/go-axis/zero/utils"

	"github.com/axis-cash/go-axis/common"
	"github.com/axis-cash/go-axis/common/math"
	"github.com/axis-cash/go-axis/consensus"
	"github.com/axis-cash/go-axis/core/state"
	"github.com/axis-cash/go-axis/core/types"
	"github.com/axis-cash/go-axis/params"
)

// Ethash proof-of-work protocol constants.
var (
	allowedFutureBlockTime          = 15 * time.Second // Max time from current time allowed for blocks, before they're considered future blocks
	bigOne                 *big.Int = big.NewInt(1)
	big200W                *big.Int = big.NewInt(2000000)
)

// Various error messages to mark blocks invalid. These should be private to
// prevent engine specific errors from being referenced in the remainder of the
// codebase, inherently breaking if the engine is swapped out. Please put common
// error types into the consensus package.
var (
	errZeroBlockTime     = errors.New("timestamp equals parent's")
	errInvalidDifficulty = errors.New("non-positive difficulty")
	errInvalidMixDigest  = errors.New("invalid mix digest")
	errInvalidPoW        = errors.New("invalid proof-of-work")
)

// Author implements consensus.Engine, returning the header's coinbase as the
// proof-of-work verified author of the block.
func (ethash *Ethash) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules of the
// stock Ethereum ethash engine.
func (ethash *Ethash) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	// If we're running a full engine faking, accept any input as valid
	if ethash.config.PowMode == ModeFullFake {
		return nil
	}
	// Short circuit if the header is known, or it's parent not
	number := header.Number.Uint64()
	if chain.GetHeader(header.Hash(), number) != nil {
		return nil
	}
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	// Sanity checks passed, do a proper verification
	return ethash.verifyHeader(chain, header, parent, seal)
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications.
func (ethash *Ethash) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	// If we're running a full engine faking, accept any input as valid
	if ethash.config.PowMode == ModeFullFake || len(headers) == 0 {
		abort, results := make(chan struct{}), make(chan error, len(headers))
		for i := 0; i < len(headers); i++ {
			results <- nil
		}
		return abort, results
	}

	// Spawn as many workers as allowed threads
	workers := runtime.GOMAXPROCS(0)
	if len(headers) < workers {
		workers = len(headers)
	}

	// Create a task channel and spawn the verifiers
	var (
		inputs = make(chan int)
		done   = make(chan int, workers)
		errors = make([]error, len(headers))
		abort  = make(chan struct{})
	)
	for i := 0; i < workers; i++ {
		go func() {
			for index := range inputs {
				errors[index] = ethash.verifyHeaderWorker(chain, headers, seals, index)
				done <- index
			}
		}()
	}

	errorsOut := make(chan error, len(headers))
	go func() {
		defer close(inputs)
		var (
			in, out = 0, 0
			checked = make([]bool, len(headers))
			inputs  = inputs
		)
		for {
			select {
			case inputs <- in:
				if in++; in == len(headers) {
					// Reached end of headers. Stop sending to workers.
					inputs = nil
				}
			case index := <-done:
				for checked[index] = true; checked[out]; out++ {
					errorsOut <- errors[out]
					if out == len(headers)-1 {
						return
					}
				}
			case <-abort:
				return
			}
		}
	}()
	return abort, errorsOut
}

func (ethash *Ethash) verifyHeaderWorker(chain consensus.ChainReader, headers []*types.Header, seals []bool, index int) error {
	var parent *types.Header
	if index == 0 {
		parent = chain.GetHeader(headers[0].ParentHash, headers[0].Number.Uint64()-1)
	} else if headers[index-1].Hash() == headers[index].ParentHash {
		parent = headers[index-1]
	}
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	if chain.GetHeader(headers[index].Hash(), headers[index].Number.Uint64()) != nil {
		return nil // known block
	}
	return ethash.verifyHeader(chain, headers[index], parent, seals[index])
}

// verifyHeader checks whether a header conforms to the consensus rules of the
// stock Ethereum ethash engine.
// See YP section 4.3.4. "Block Header Validity"
func (ethash *Ethash) verifyHeader(chain consensus.ChainReader, header, parent *types.Header, seal bool) error {
	if !superzk.CheckLICr(header.Coinbase.ToPKr(), &header.Licr, header.Number.Uint64()) {
		return fmt.Errorf("invalid Licr : pkr %v, licr %v", header.Coinbase, header.Licr)
	}

	if !header.Valid() {
		return fmt.Errorf("invalid Licr : pkr %v, licr %v, disable", header.Coinbase, header.Licr)
	}
	// Ensure that the header's extra-data section is of a reasonable size
	if uint64(len(header.Extra)) > params.MaximumExtraDataSize {
		return fmt.Errorf("extra-data too long: %d > %d", len(header.Extra), params.MaximumExtraDataSize)
	}
	// Verify the header's timestamp
	if header.Time.Cmp(big.NewInt(time.Now().Add(allowedFutureBlockTime).Unix())) > 0 {
		return consensus.ErrFutureBlock
	}

	if header.Time.Cmp(parent.Time) <= 0 {
		return errZeroBlockTime
	}
	// Verify the block's difficulty based in it's timestamp and parent's difficulty
	expected := ethash.CalcDifficulty(chain, header.Time.Uint64(), parent)

	if expected.Cmp(header.Difficulty) != 0 {
		return fmt.Errorf("invalid difficulty: have %v, want %v", header.Difficulty, expected)
	}
	// Verify that the gas limit is <= 2^63-1
	cap := uint64(0x7fffffffffffffff)
	if header.GasLimit > cap {
		return fmt.Errorf("invalid gasLimit: have %v, max %v", header.GasLimit, cap)
	}
	// Verify that the gasUsed is <= gasLimit
	if header.GasUsed > header.GasLimit {
		return fmt.Errorf("invalid gasUsed: have %d, gasLimit %d", header.GasUsed, header.GasLimit)
	}

	// Verify that the gas limit remains within allowed bounds
	diff := int64(parent.GasLimit) - int64(header.GasLimit)
	divisor := uint64(1024)
	if diff < 0 {
		diff *= -1
		divisor = uint64(128)
	}
	limit := parent.GasLimit / divisor

	if uint64(diff) >= limit || header.GasLimit < params.MinGasLimit {
		return fmt.Errorf("invalid gas limit: have %d, want %d += %d", header.GasLimit, parent.GasLimit, limit)
	}
	// Verify that the block number is parent's +1
	if diff := new(big.Int).Sub(header.Number, parent.Number); diff.Cmp(big.NewInt(1)) != 0 {
		return consensus.ErrInvalidNumber
	}
	// Verify the engine specific seal securing the block
	if seal {
		if err := ethash.VerifySeal(chain, header); err != nil {
			return err
		}
	}
	return nil
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns
// the difficulty that a new block should have when created at time
// given the parent block's time and difficulty.
func (ethash *Ethash) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	return CalcDifficulty(chain.Config(), time, parent)
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns
// the difficulty that a new block should have when created at time
// given the parent block's time and difficulty.
func CalcDifficulty(config *params.ChainConfig, time uint64, parent *types.Header) *big.Int {
	if zconfig.IsTestFork() {
		if parent.Number.Uint64() >= zconfig.GetTestForkStartBlock() {
			return big.NewInt(10000)
		}
	}

	return calcDifficultyAutumnTwilight(time, parent)
}

// Some weird constants to avoid constant memory allocs for them.
var (
	expDiffPeriod = big.NewInt(100000)
	big1          = big.NewInt(1)
	big2          = big.NewInt(2)
	big6          = big.NewInt(6)
	big9          = big.NewInt(9)
	bigMinus99    = big.NewInt(-99)
)

// calcDifficultyAutumnTwilight is the difficulty adjustment algorithm. It returns
// the difficulty that a new block should have when created at time given the
// parent block's time and difficulty. The calculation uses the AutumnTwilight rules.
func calcDifficultyAutumnTwilight(time uint64, parent *types.Header) *big.Int {
	// https://github.com/ethereum/EIPs/issues/100.
	// algorithm:
	// diff = (parent_diff +
	//         (parent_diff / 2048 * max((1 - ((timestamp - parent.timestamp) // 9), -99))
	//        )
	if axisparam.Is_Dev() {
		return  big.NewInt(1);
	}
	bigTime := new(big.Int).SetUint64(time)
	bigParentTime := new(big.Int).Set(parent.Time)

	// holds intermediate values to make the algo easier to read & audit
	x := new(big.Int)
	y := new(big.Int)

	// 1 - (block_timestamp - parent_timestamp) // 9
	x.Sub(bigTime, bigParentTime)
	x.Div(x, big9)
	x.Sub(big1, x)
	// max(1 - (block_timestamp - parent_timestamp) // 9, -99)
	if x.Cmp(bigMinus99) < 0 {
		x.Set(bigMinus99)
	}
	// parent_diff + (parent_diff / 2048 * max(1 - ((timestamp - parent.timestamp) // 9), -99))
	y.Div(parent.Difficulty, params.DifficultyBoundDivisor)
	x.Mul(y, x)
	x.Add(parent.Difficulty, x)

	// minimum difficulty can ever be (before exponential factor)
	if x.Cmp(params.MinimumDifficulty) < 0 {
		x.Set(params.MinimumDifficulty)
	}
	return x
}

// calcDifficultyFrontier is the difficulty adjustment algorithm. It returns the
// difficulty that a new block should have when created at time given the parent
// block's time and difficulty. The calculation uses the Frontier rules.
func calcDifficultyFrontier(time uint64, parent *types.Header) *big.Int {
	diff := new(big.Int)
	adjust := new(big.Int).Div(parent.Difficulty, params.DifficultyBoundDivisor)
	bigTime := new(big.Int)
	bigParentTime := new(big.Int)

	bigTime.SetUint64(time)
	bigParentTime.Set(parent.Time)

	if bigTime.Sub(bigTime, bigParentTime).Cmp(params.DurationLimit) < 0 {
		diff.Add(parent.Difficulty, adjust)
	} else {
		diff.Sub(parent.Difficulty, adjust)
	}
	if diff.Cmp(params.MinimumDifficulty) < 0 {
		diff.Set(params.MinimumDifficulty)
	}

	periodCount := new(big.Int).Add(parent.Number, big1)
	periodCount.Div(periodCount, expDiffPeriod)
	if periodCount.Cmp(big1) > 0 {
		// diff = diff + 2^(periodCount - 2)
		expDiff := periodCount.Sub(periodCount, big2)
		expDiff.Exp(big2, expDiff, nil)
		diff.Add(diff, expDiff)
		diff = math.BigMax(diff, params.MinimumDifficulty)
	}
	return diff
}

// VerifySeal implements consensus.Engine, checking whether the given block satisfies
// the PoW difficulty requirements.
func (ethash *Ethash) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	// If we're running a fake PoW, accept any seal as valid
	if ethash.config.PowMode == ModeFake || ethash.config.PowMode == ModeFullFake {
		time.Sleep(ethash.fakeDelay)
		if ethash.fakeFail == header.Number.Uint64() {
			return errInvalidPoW
		}
		return nil
	}
	// If we're running a shared PoW, delegate verification to it
	if ethash.shared != nil {
		return ethash.shared.VerifySeal(chain, header)
	}
	// Ensure that we have a valid difficulty for the block
	if header.Difficulty.Sign() <= 0 {
		return errInvalidDifficulty
	}
	// Recompute the digest and PoW value and verify against the header
	number := header.Number.Uint64()

	cache := ethash.cache(number)
	size := datasetSize(number)
	if ethash.config.PowMode == ModeTest {
		size = 32 * 1024
	}

	var digest []byte
	var result []byte
	if number >= axisparam.XIP1() {
		// dataset := ethash.dataset_async(number)
		// if dataset.generated() {
		//	digest, result = progpowFull(dataset.dataset, header.HashPow().Bytes(), header.Nonce.Uint64(), number)
		// } else {
		digest, result = progpowLightWithoutCDag(size, cache.cache, cache.cdag, header.HashPow().Bytes(), header.Nonce.Uint64(), number)
		// }
	} else {
		digest, result = hashimotoLight(size, cache.cache, header.HashPow().Bytes(), header.Nonce.Uint64(), number)
	}
	// Caches are unmapped in a finalizer. Ensure that the cache stays live
	// until after the call to hashimotoLight so it's not unmapped while being used.
	runtime.KeepAlive(cache)

	if !bytes.Equal(header.MixDigest[:], digest) {
		return errInvalidMixDigest
	}
	target := new(big.Int).Div(maxUint256, header.ActualDifficulty())
	if new(big.Int).SetBytes(result).Cmp(target) > 0 {
		return errInvalidPoW
	}
	return nil
}

// Prepare implements consensus.Engine, initializing the difficulty field of a
// header to conform to the ethash protocol. The changes are done inline.
func (ethash *Ethash) Prepare(chain consensus.ChainReader, header *types.Header) error {
	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	header.Difficulty = ethash.CalcDifficulty(chain, header.Time.Uint64(), parent)
	return nil
}


// Finalize implements consensus.Engine, accumulating the block rewards,
// setting the final state and assembling the block.
func (ethash *Ethash) Finalize(chain consensus.ChainReader, header *types.Header, stateDB *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt, gasReward uint64) (*types.Block, error) {
	if header.Number.Uint64() == axisparam.XIP1() {
		stateDB.SetBalance(state.EmptyAddress, "AXIS", new(big.Int))
	}
	// Accumulate any block rewards and commit the final state root
	accumulateRewards(chain.Config(), stateDB, header, gasReward)

	stateDB.NextZState().PreGenerateRoot(header, chain)

	header.Root = stateDB.IntermediateRoot(true)

	// Header seems complete, assemble into a block and return
	return types.NewBlock(header, txs, receipts), nil
}

var (
	base    = big.NewInt(1e+17)
	big100  = big.NewInt(100)
	oneAxis = new(big.Int).Mul(big.NewInt(10), base)

	oriReward    = new(big.Int).Mul(big.NewInt(66773505743), big.NewInt(1000000000))
	interval     = big.NewInt(8760000)
	halveNimber  = big.NewInt(7908000)
	difficultyL1 = big.NewInt(340000000)
	difficultyL2 = big.NewInt(1700000000)
	difficultyL3 = big.NewInt(4000000000)
	difficultyL4 = big.NewInt(17000000000)

	lReward   = new(big.Int).Mul(big.NewInt(175), base)
	hReward = new(big.Int).Mul(big.NewInt(285), base)

	argA, _ = new(big.Int).SetString("985347985347985", 10)
	argB, _ = new(big.Int).SetString("16910256410256400000", 10)

	//teamRewardPool      = common.BytesToAddress(crypto.Keccak512([]byte{1}))
	//teamAddress      = common.Base58ToAddress( "NGsyb5qeqLJs8TCvAeWSMmPGRn85H2SRWn4Y9Gq8sMDeqN5jRJx86akFTCLunNiWLNYgAcdHn4jGg4KtG8Pbbanp8sRNKHAXCcBhNZJVp5yadBzZ5wHQiePnx1TC6i5bPDk")
	techContractAddress = common.Base58ToAddress("3nDfqP1rcYcDktu7ewSATm9uvfcF1D6DruVdrF7jTv9RQnCDBKKkyT7EZRgFZihso2x6XwcceBBtWxLyJH285Ah3")
	marketContractAddress = common.Base58ToAddress("TF8RFwwCYYweU7NrKSfwFb746x8mcCuF7UGcJnnuoJsPgE2jELf65RGDW34k7zWqGcGLC46ZePjRUEUMDbs7atL")
	foundContractAddress = common.Base58ToAddress("pRs6cEZSe82N8ho17mJuhodTzCcoXhsnMfbjY6aP44pCWhLeUx8wopwRUrXjepz6Dw2VeKKvHxFyzK15mXEjDZd")
)

func Halve(blockNumber *big.Int) *big.Int {
	i := new(big.Int).Add(new(big.Int).Div(new(big.Int).Sub(blockNumber, halveNimber), interval), big1)
	return new(big.Int).Exp(big2, i, nil)
}

// AccumulateRewards credits the coinbase of the given block with the mining
// reward. The total reward consists of the static block reward .
func accumulateRewards(config *params.ChainConfig, statedb *state.StateDB, header *types.Header, gasReward uint64) {

	var reward *big.Int
	//if header.Number.Uint64() >= axisparam.XIP4() {
	reward = accumulateRewardsV1(statedb, header)
	//}

	if axisparam.Is_Dev() {
		reward = new(big.Int).Set(new(big.Int).Mul(big.NewInt(10000), oneAxis))
	}
	// log.Info(fmt.Sprintf("BlockNumber = %v, gasLimie = %v, gasUsed = %v, reward = %v", header.Number.Uint64(), header.GasLimit, header.GasUsed, reward))
	reward.Add(reward, new(big.Int).SetUint64(gasReward))

	asset := assets.Asset{Tkn: &assets.Token{
		Currency: *common.BytesToHash(common.LeftPadBytes([]byte("AXIS"), 32)).HashToUint256(),
		Value:    utils.U256(*reward),
	},
	}
	statedb.NextZState().AddTxOut(header.Coinbase, asset, common.BytesToHash([]byte{1}))
}

func accumulateRewardsV1(statedb *state.StateDB, header *types.Header) *big.Int {
	diff := new(big.Int).Div(header.Difficulty, big.NewInt(1000000000))
	reward := new(big.Int).Add(new(big.Int).Mul(argA, diff), argB)

	if reward.Cmp(lReward) < 0 {
		reward = new(big.Int).Set(lReward)
	} else if reward.Cmp(hReward) > 0 {
		reward = new(big.Int).Set(hReward)
	}

	i := new(big.Int).Add(new(big.Int).Div(new(big.Int).Sub(header.Number, halveNimber), interval), big1)
	reward.Div(reward, new(big.Int).Exp(big2, i, nil))

	teamReward := new(big.Int).Div(hReward, big.NewInt(5))
	teamReward = new(big.Int).Div(teamReward, new(big.Int).Exp(big2, i, nil))
	if header.Number.Uint64() >= axisparam.VP0() {
		statedb.AddBalance(techContractAddress, "AXIS", new(big.Int).Div(new(big.Int).Mul(teamReward,
			big.NewInt(35)), big.NewInt(100)))
		statedb.AddBalance(marketContractAddress, "AXIS", new(big.Int).Div(new(big.Int).Mul(teamReward,
			big.NewInt(35)), big.NewInt(100)))
		statedb.AddBalance(foundContractAddress, "AXIS", new(big.Int).Div(new(big.Int).Mul(teamReward,
			big.NewInt(30)), big.NewInt(100)))
	}
	

	/*
	if header.Number.Uint64()%5000 == 0 {
		balance := statedb.GetBalance(teamRewardPool, "AXIS")
		statedb.SubBalance(teamRewardPool, "AXIS", balance)
		assetTeam := assets.Asset{Tkn: &assets.Token{
			Currency: *common.BytesToHash(common.LeftPadBytes([]byte("AXIS"), 32)).HashToUint256(),
			Value:    utils.U256(*balance),
		},
		}
		statedb.NextZState().AddTxOut(teamAddress, assetTeam, common.Hash{})
	}*/
	return reward
}

