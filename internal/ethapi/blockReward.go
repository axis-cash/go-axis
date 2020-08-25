package ethapi

import (
	"github.com/axis-cash/go-axis/common"
	"github.com/axis-cash/go-axis/core/state"
	"github.com/axis-cash/go-axis/zero/txs/assets"
	"github.com/axis-cash/go-axis/zero/utils"
	"math/big"

	"github.com/axis-cash/go-axis/core/types"

	"github.com/axis-cash/go-axis-import/axisparam"
)

var (
	big1                = big.NewInt(1)
	big2                = big.NewInt(2)
	big6                = big.NewInt(6)
	big9                = big.NewInt(9)
	bigMinus99          = big.NewInt(-99)
	bigOne     *big.Int = big.NewInt(1)
	big200W    *big.Int = big.NewInt(2000000)
	base                = big.NewInt(1e+17)
	big100              = big.NewInt(100)
	oneAxis             = new(big.Int).Mul(big.NewInt(10), base)

	lReward   = new(big.Int).Mul(big.NewInt(176), base)
	hReward   = new(big.Int).Mul(big.NewInt(445), base)
	hRewardV4 = new(big.Int).Mul(big.NewInt(356), base)

	argA, _ = new(big.Int).SetString("985347985347985", 10)
	argB, _ = new(big.Int).SetString("16910256410256400000", 10)

	oriReward    = new(big.Int).Mul(big.NewInt(66773505743), big.NewInt(1000000000))
	interval     = big.NewInt(8294400)
	halveNimber  = big.NewInt(3057600)
	difficultyL1 = big.NewInt(340000000)
	difficultyL2 = big.NewInt(1700000000)
	difficultyL3 = big.NewInt(4000000000)
	difficultyL4 = big.NewInt(17000000000)
)

func accumulateRewardsV1(diff *big.Int, gasUsed uint64, gasLimit uint64) *big.Int {

	reward := new(big.Int).Mul(big.NewInt(350), base)

	difficulty := big.NewInt(1717986918)

	if diff.Cmp(difficulty) < 0 {
		ratio := new(big.Int).Div(new(big.Int).Mul(diff, big100), difficulty).Uint64()
		if ratio >= 80 {
			reward = reward.Mul(reward, big.NewInt(4)).Div(reward, big.NewInt(5))
		} else if ratio >= 60 {
			reward = reward.Mul(reward, big.NewInt(3)).Div(reward, big.NewInt(5))
		} else if ratio >= 40 {
			reward = reward.Mul(reward, big.NewInt(2)).Div(reward, big.NewInt(5))
		} else if ratio >= 20 {
			reward = reward.Mul(reward, big.NewInt(1)).Div(reward, big.NewInt(5))
		} else {
			reward = big.NewInt(0).Set(oneAxis)
		}
	}

	ratio := new(big.Int).Div(new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), big100), new(big.Int).SetUint64(gasLimit)).Uint64()
	if ratio >= 80 {
		reward = new(big.Int).Div(new(big.Int).Mul(reward, big6), big.NewInt(5))
	} else {
		reward = reward.Mul(reward, big.NewInt(4)).Div(reward, big.NewInt(5))
	}

	if reward.Cmp(oneAxis) < 0 {
		reward = big.NewInt(0).Set(oneAxis)
	}
	return reward
}

func accumulateRewardsV2(number, diff *big.Int) [2]*big.Int {
	var res [2]*big.Int
	rewardStd := new(big.Int).Set(oriReward)
	if number.Uint64() >= halveNimber.Uint64() {
		i := new(big.Int).Add(new(big.Int).Div(new(big.Int).Sub(number, halveNimber), interval), big1)
		rewardStd.Div(rewardStd, new(big.Int).Exp(big2, i, nil))
	}

	var reward *big.Int
	if diff.Cmp(difficultyL1) < 0 { //<3.4
		reward = new(big.Int).Mul(big.NewInt(10), base)
	} else if diff.Cmp(difficultyL2) < 0 { //<17
		ratio := new(big.Int).Add(new(big.Int).Mul(big.NewInt(56), base), new(big.Int).Mul(big.NewInt(16470000000), new(big.Int).Sub(diff, difficultyL1)))
		reward = new(big.Int).Div(new(big.Int).Mul(rewardStd, ratio), oriReward)
	} else if diff.Cmp(difficultyL3) < 0 { //<40
		ratio := new(big.Int).Add(new(big.Int).Mul(big.NewInt(280), base), new(big.Int).Mul(big.NewInt(2170000000), new(big.Int).Sub(diff, difficultyL2)))
		reward = new(big.Int).Div(new(big.Int).Mul(rewardStd, ratio), oriReward)
	} else if diff.Cmp(difficultyL4) < 0 { //<170
		ratio := new(big.Int).Add(new(big.Int).Mul(big.NewInt(330), base), new(big.Int).Mul(big.NewInt(2590000000), new(big.Int).Sub(diff, difficultyL3)))
		reward = new(big.Int).Div(new(big.Int).Mul(rewardStd, ratio), oriReward)
	} else {
		reward = rewardStd
	}
	res[0] = reward
	res[1] = new(big.Int).Div(rewardStd, big.NewInt(15))
	return res
}

func accumulateRewardsV3(number, bdiff *big.Int) [2]*big.Int {
	var res [2]*big.Int
	diff := new(big.Int).Div(bdiff, big.NewInt(1000000000))
	reward := new(big.Int).Add(new(big.Int).Mul(argA, diff), argB)

	if reward.Cmp(lReward) < 0 {
		reward = new(big.Int).Set(lReward)
	} else if reward.Cmp(hReward) > 0 {
		reward = new(big.Int).Set(hReward)
	}

	i := new(big.Int).Add(new(big.Int).Div(new(big.Int).Sub(number, halveNimber), interval), big1)
	reward.Div(reward, new(big.Int).Exp(big2, i, nil))

	res[0] = reward
	res[1] = big.NewInt(0)
	return res
}

func accumulateRewardsV4(number, bdiff *big.Int) [2]*big.Int {
	var res [2]*big.Int
	diff := new(big.Int).Div(bdiff, big.NewInt(1000000000))
	reward := new(big.Int).Add(new(big.Int).Mul(argA, diff), argB)

	if reward.Cmp(lReward) < 0 {
		reward = new(big.Int).Set(lReward)
	} else if reward.Cmp(hRewardV4) > 0 {
		reward = new(big.Int).Set(hRewardV4)
	}

	i := new(big.Int).Add(new(big.Int).Div(new(big.Int).Sub(number, halveNimber), interval), big1)
	reward.Div(reward, new(big.Int).Exp(big2, i, nil))

	res[0] = reward
	res[1] = big.NewInt(0)
	return res

}

func accumulateRewardsV5(statedb *state.StateDB, header *types.Header) *big.Int {
	diff := new(big.Int).Div(header.Difficulty, big.NewInt(1000000000))
	reward := new(big.Int).Add(new(big.Int).Mul(argA, diff), argB)

	if reward.Cmp(lReward) < 0 {
		reward = new(big.Int).Set(lReward)
	} else if reward.Cmp(hRewardV4) > 0 {
		reward = new(big.Int).Set(hRewardV4)
	}
	reward.Div(reward, big.NewInt(2))
	i := new(big.Int).Add(new(big.Int).Div(new(big.Int).Sub(header.Number, halveNimber), interval), big1)
	reward.Div(reward, new(big.Int).Exp(big2, i, nil))

	teamReward := new(big.Int).Div(hRewardV4, big.NewInt(4))
	teamReward = new(big.Int).Div(teamReward, new(big.Int).Exp(big2, i, nil))
	statedb.AddBalance(teamRewardPool, "AXIS", teamReward)

	if header.Number.Uint64()%5 == 0 {
		balance := statedb.GetBalance(teamRewardPool, "AXIS")
		statedb.SubBalance(teamRewardPool, "AXIS", balance)
		assetTeam := assets.Asset{Tkn: &assets.Token{
			Currency: *common.BytesToHash(common.LeftPadBytes([]byte("AXIS"), 32)).HashToUint256(),
			Value:    utils.U256(*balance),
		},
		}
		statedb.NextZState().AddTxOut(teamAddress, assetTeam, common.Hash{})
	}
	return reward
}

/**
  [0] block reward
  [1] community reward
  [2] team reward
*/
func GetBlockReward(block *types.Block) [2]*big.Int {
	number := block.Number()
	diff := block.Difficulty()
	gasUsed := block.GasUsed()
	gasLimit := block.GasLimit()
	if number.Uint64() >= axisparam.XIP7() {
		return accumulateRewardsV5(number, diff)
	} else if number.Uint64() >= axisparam.XIP4() {
		return accumulateRewardsV4(number, diff)
	} else if number.Uint64() >= axisparam.XIP3() {
		return accumulateRewardsV3(number, diff)
	} else if number.Uint64() >= axisparam.XIP1() {
		return accumulateRewardsV2(number, diff)
	} else {
		var res [2]*big.Int
		res[0] = accumulateRewardsV1(diff, gasUsed, gasLimit)
		res[1] = big.NewInt(0)
		return res
	}
}
