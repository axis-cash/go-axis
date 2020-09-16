package zconfig

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/axis-cash/go-axis-import/c_type"
	"github.com/axis-cash/go-axis-import/axisparam"
)

/*******
var current=axis.blockNumber
for(var i=0;i<current;i+=10000) {
	console.log('{"Num":'+i+',"Root":"'+axis.getBlock(i).stateRoot+'"},')
}
********/

var checkpoints_json = `[
]`

type checkPoint struct {
	Num  uint64
	Root c_type.Uint256
}

type checkPoints struct {
	maxNum uint64
	points map[uint64][]byte
}

func (self *checkPoints) MaxNum() (ret uint64) {
	if axisparam.Is_Dev() {
		return 0
	} else {
		return self.maxNum
	}
}

func (self *checkPoints) Check(num uint64, root []byte) (e error) {
	if num > self.maxNum {
		panic(fmt.Errorf("check points error: the num > maxNum %d-%s", num, hex.EncodeToString(root)))
	}
	if (num > 0) && (num%10000 == 0) {
		if rt, ok := self.points[num]; !ok {
			return fmt.Errorf("check points error: can not find the point %d-%s", num, hex.EncodeToString(root))
		} else {
			if bytes.Compare(rt, root) != 0 {
				return fmt.Errorf("check points error: the point are not match %d-%s (%s)", num, hex.EncodeToString(root), hex.EncodeToString(rt))
			} else {
				return nil
			}
		}
	} else {
		return nil
	}
}

func newCheckPoints() (ret checkPoints) {
	ret.points = make(map[uint64][]byte)
	var cps []*checkPoint
	json.Unmarshal([]byte(checkpoints_json), &cps)

	for _, cp := range cps {
		ret.points[cp.Num] = cp.Root[:]
		if cp.Num > ret.maxNum {
			ret.maxNum = cp.Num
		}
	}
	return
}

var CheckPoints = newCheckPoints()
