package exchange

import (
	"github.com/axis-cash/go-axis-import/c_superzk"
	"github.com/axis-cash/go-axis-import/c_type"
	"github.com/axis-cash/go-axis/zero/txs/assets"
)

type Utxo struct {
	Pkr    c_type.PKr
	Root   c_type.Uint256
	TxHash c_type.Uint256
	Nil    c_type.Uint256
	Num    uint64
	Asset  assets.Asset
	IsZ    bool
	Ignore bool
	flag   int
}

func (utxo *Utxo) NilTxType() string {
	if c_superzk.IsSzkNil(&utxo.Nil) {
		return "SZK"
	} else {
		return "CZERO"
	}
}
