package flight

import (
	"github.com/axis-cash/go-axis-import/c_type"
	"github.com/axis-cash/go-axis/zero/txtool"
)

type PreTxParam struct {
	Gas      uint64
	GasPrice uint64
	From     c_type.PKr
	Ins      []c_type.Uint256
	Outs     []txtool.GOut
}
