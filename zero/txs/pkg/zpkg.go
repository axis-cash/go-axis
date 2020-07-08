package pkg

import (
	"github.com/axis-cash/go-axis-import/c_type"
	"github.com/axis-cash/go-axis/crypto/sha3"
	"github.com/axis-cash/go-axis/zero/utils"
)

type Pkg_Z struct {
	AssetCM c_type.Uint256
	EInfo   c_type.Einfo
}

func (this Pkg_Z) ToRef() (ret *Pkg_Z) {
	ret = &this
	return
}

func (self *Pkg_Z) ToHash() (ret c_type.Uint256) {
	d := sha3.NewKeccak256()
	d.Write(self.AssetCM[:])
	d.Write(self.EInfo[:])
	copy(ret[:], d.Sum(nil))
	return ret
}

func (self *Pkg_Z) Clone() (ret Pkg_Z) {
	utils.DeepCopy(&ret, self)
	return
}
