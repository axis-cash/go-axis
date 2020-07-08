package localdb

import (
	"math/big"

	"github.com/axis-cash/go-axis-import/c_type"
	"github.com/axis-cash/go-axis/rlp"
	"github.com/axis-cash/go-axis/axisdb"
	"github.com/axis-cash/go-axis/zero/txs/zstate/tri"
)

type Block struct {
	Roots []c_type.Uint256
	Dels  []c_type.Uint256
	Pkgs  []c_type.Uint256
}

func (self *Block) Serial() (ret []byte, e error) {
	if self != nil {
		if bytes, err := rlp.EncodeToBytes(self); err != nil {
			e = err
			return
		} else {
			ret = bytes
			return
		}
	} else {
		return
	}
}

type BlockGet struct {
	Out *Block
}

func (self *BlockGet) Unserial(v []byte) (e error) {
	if len(v) == 0 {
		return
	} else {
		out := Block{}
		if err := rlp.DecodeBytes(v, &out); err != nil {
			return
		} else {
			self.Out = &out
			return
		}
	}
}
func BlockKey(num uint64, hash *c_type.Uint256) []byte {
	block_key := []byte("$AXIS_ZSTATE_BLOCK_SHOOTCUT$")
	block_key = append(block_key, big.NewInt(int64(num)).Bytes()...)
	block_key = append(block_key, []byte("$")...)
	block_key = append(block_key, hash[:]...)
	return block_key
}

func PutBlock(db axisdb.Putter, num uint64, hash *c_type.Uint256, block *Block) {
	blockkey := BlockKey(num, hash)
	tri.UpdateDBObj(db, blockkey, block)
}

func GetBlock(db axisdb.Database, num uint64, hash *c_type.Uint256) (ret *Block) {
	blockkey := BlockKey(num, hash)
	blockget := BlockGet{}
	tri.GetDBObj(db, blockkey, &blockget)
	ret = blockget.Out
	return
}
