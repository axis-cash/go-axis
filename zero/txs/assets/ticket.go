package assets

import (
	"github.com/axis-cash/go-axis-import/c_type"
	"github.com/axis-cash/go-axis/crypto"
	"github.com/axis-cash/go-axis/zero/utils"
)

type Ticket struct {
	Category c_type.Uint256
	Value    c_type.Uint256
}

func (self *Ticket) Clone() (ret Ticket) {
	utils.DeepCopy(&ret, self)
	return
}

func (this Ticket) ToRef() (ret *Ticket) {
	ret = &this
	return
}

func (self *Ticket) ToHash() (ret c_type.Uint256) {
	if self == nil {
		return
	} else {
		hash := crypto.Keccak256(
			self.Category[:],
			self.Value[:],
		)
		copy(ret[:], hash)
		return
	}
}
