package localdb

import (
	"time"

	"github.com/axis-cash/go-axis-import/c_type"
	"github.com/axis-cash/go-axis/rlp"
	"github.com/axis-cash/go-axis/axisdb"
	"github.com/axis-cash/go-axis/zero/txs/zstate/tri"
	"github.com/axis-cash/go-axis/zero/utils"
)

type OutStat struct {
	Time  int64
	Value utils.U256
	Z     bool
}

func (self *OutStat) Serial() (ret []byte, e error) {
	if self != nil {
		return rlp.EncodeToBytes(self)
	} else {
		return
	}
}

type OutStatGet struct {
	out *OutStat
}

func (self *OutStatGet) Unserial(v []byte) (e error) {
	if len(v) == 0 {
		self.out = nil
		return
	} else {
		self.out = &OutStat{}
		if err := rlp.DecodeBytes(v, self.out); err != nil {
			e = err
			return
		} else {
			return
		}
	}
}

func outStatName(root *c_type.Uint256) (ret []byte) {
	ret = []byte("$ZSTATE_OUT_STAT$")
	ret = append(ret, root[:]...)
	return
}

func UpdateOutStat(db axisdb.Putter, root *c_type.Uint256, os *OutStat) {
	os.Time = time.Now().UnixNano()
	tri.UpdateDBObj(db, outStatName(root), os)
}

func GetOutStat(db axisdb.Getter, root *c_type.Uint256) (ret *OutStat) {
	get := OutStatGet{}
	tri.GetDBObj(db, outStatName(root), &get)
	if get.out != nil {
		ret = get.out
	}
	return
}
