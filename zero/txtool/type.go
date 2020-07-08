package txtool

import (
	"github.com/axis-cash/go-axis-import/c_type"
	"github.com/axis-cash/go-axis/common/hexutil"
	"github.com/axis-cash/go-axis/zero/localdb"
	"github.com/axis-cash/go-axis/zero/txs/assets"
	"github.com/axis-cash/go-axis/zero/txs/stx"
)

type Kr struct {
	SKr c_type.PKr
	PKr c_type.PKr
}

type Out struct {
	Root  c_type.Uint256
	State localdb.RootState
}

type TDOut struct {
	Asset assets.Asset
	Memo  c_type.Uint512
	Nils  []c_type.Uint256
}

type DOut struct {
	Asset assets.Asset
	Memo  c_type.Uint512
	Nil   c_type.Uint256
}

type Block struct {
	Num  hexutil.Uint64
	Hash c_type.Uint256
	Outs []Out
	Nils []c_type.Uint256
	Pkgs []localdb.ZPkg
}

type Witness struct {
	Pos    hexutil.Uint64
	Paths  [c_type.DEPTH]c_type.Uint256
	Anchor c_type.Uint256
}

type Tx struct {
	Hash c_type.Uint256
	Tx   stx.T
}
