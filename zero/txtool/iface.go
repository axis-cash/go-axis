package txtool

import (
	"errors"
	"math/big"

	"github.com/axis-cash/go-axis-import/axisparam"

	"github.com/axis-cash/go-axis/zero/utils"

	"github.com/axis-cash/go-axis/common/hexutil"

	"github.com/axis-cash/go-axis-import/c_type"
	"github.com/axis-cash/go-axis/zero/txs/assets"
	"github.com/axis-cash/go-axis/zero/txs/stx"
)

type GIn struct {
	SKr     c_type.PKr
	Out     Out
	Witness Witness
	A       *c_type.Uint256
	Ar      *c_type.Uint256
	ArOld   *c_type.Uint256
	Vskr    *c_type.Uint256
	CC      *c_type.Uint256
}

type GOut struct {
	PKr   c_type.PKr
	Asset assets.Asset
	Memo  c_type.Uint512
	Ar    *c_type.Uint256
}

type GTx struct {
	Gas      hexutil.Uint64
	GasPrice hexutil.Big
	Tx       stx.T
	Hash     c_type.Uint256
	Roots    []c_type.Uint256
	Keys     []c_type.Uint256
	Bases    []c_type.Uint256
}

type GPkgCloseCmd struct {
	Id      c_type.Uint256
	Owner   c_type.PKr
	AssetCM c_type.Uint256
	Ar      c_type.Uint256
}

type GPkgTransferCmd struct {
	Id    c_type.Uint256
	Owner c_type.PKr
	PKr   c_type.PKr
}

type GPkgCreateCmd struct {
	Id    c_type.Uint256
	PKr   c_type.PKr
	Asset assets.Asset
	Memo  c_type.Uint512
	Ar    c_type.Uint256
}

type Cmds struct {
	//Share
	BuyShare *stx.BuyShareCmd
	//Pool
	RegistPool *stx.RegistPoolCmd
	ClosePool  *stx.ClosePoolCmd
	//Contract
	Contract *stx.ContractCmd
	//Package
	PkgCreate   *GPkgCreateCmd
	PkgTransfer *GPkgTransferCmd
	PkgClose    *GPkgCloseCmd
}

type GTxParam struct {
	Gas      uint64
	GasPrice *big.Int
	Fee      assets.Token
	From     Kr
	Ins      []GIn
	Outs     []GOut
	Cmds     Cmds
	Z        *bool
	Num      *uint64
}

func (self *GTxParam) IsSzk() (ret bool) {
	check := utils.NewPKrChecker()
	check.AddPKr(&self.From.PKr)
	for _, in := range self.Ins {
		check.AddPKr(in.Out.State.OS.ToPKr())
	}
	for _, out := range self.Outs {
		check.AddPKr(&out.PKr)
	}
	return check.IsSzk()
}

func (self *GTxParam) GenZ() (e error) {
	Z := true
	self.Z = &Z
	if Ref_inst.Bc != nil {
		num := Ref_inst.Bc.GetCurrenHeader().Number.Uint64()
		if num < axisparam.XIP6() {
			if len(self.Outs) > 9 {
				e = errors.New("only 9 outs allowed before sip6")
				return
			}
		}
		self.Num = &num
	}
	return
}
