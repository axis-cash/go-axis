package tx

import (
	"fmt"
	"testing"

	"github.com/axis-cash/go-axis-import/c_type"
	"github.com/axis-cash/go-axis/zero/utils"
)

func TestT_TokenCost(t *testing.T) {
	axisCy := utils.CurrencyToUint256("AXIS")
	fmt.Printf("%t\n", axisCy)
	cy := utils.CurrencyToUint256("d")
	ret := make(map[c_type.Uint256]utils.U256)
	ret[axisCy] = utils.NewU256(24)
	if cost, ok := ret[axisCy]; ok {
		add := utils.NewU256(12)
		cost.AddU(&add)
		ret[axisCy] = cost
	} else {
		cost := utils.NewU256(48)
		ret[cy] = cost
	}

	fmt.Printf("%t", ret)

}

func Test_ReservedTree(t *testing.T) {
	reserveds := NewReserveds(10240)

	reserveds.Insert(1025)
	reserveds.Insert(1023)
	reserveds.Insert(1000)
	reserveds.Insert(900)
	reserveds.Insert(500)

}
