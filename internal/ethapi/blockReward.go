package ethapi

import (
	"math/big"

	"github.com/axis-cash/go-axis/core/types"

	//"github.com/axis-cash/go-axis-import/axisparam"
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

	lReward   = new(big.Int).Mul(big.NewInt(175), base)
	hReward = new(big.Int).Mul(big.NewInt(285), base)

	argA, _ = new(big.Int).SetString("985347985347985", 10)
	argB, _ = new(big.Int).SetString("16910256410256400000", 10)

	oriReward    = new(big.Int).Mul(big.NewInt(66773505743), big.NewInt(1000000000))
	interval     = big.NewInt(8760000)
	halveNimber  = big.NewInt(7908000)
	difficultyL1 = big.NewInt(340000000)
	difficultyL2 = big.NewInt(1700000000)
	difficultyL3 = big.NewInt(4000000000)
	difficultyL4 = big.NewInt(17000000000)
)


func accumulateRewardsV1(number, bdiff *big.Int) [2]*big.Int {
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
	hRewardOrg :=  new(big.Int).Set(hReward)
	hRewardOrg.Div(hRewardOrg, new(big.Int).Exp(big2, i, nil))

	res[0] = reward
	res[1] = new(big.Int).Sub(hRewardOrg, reward)
	return res

}

/**
  [0] block reward
  [1] community reward
  [2] team reward
*/
func GetBlockReward(block *types.Block) [2]*big.Int {
	number := block.Number()
	diff := block.Difficulty()
	//gasUsed := block.GasUsed()
	//gasLimit := block.GasLimit()
	//if number.Uint64() >= axisparam.XIP1() {
	return accumulateRewardsV1(number, diff)
	//}
}
