package pkgstate

import (
	"fmt"
	"sync"

	"github.com/axis-cash/go-axis-import/c_superzk"

	"github.com/axis-cash/go-axis/axisdb"

	"github.com/axis-cash/go-axis/common/hexutil"
	"github.com/axis-cash/go-axis/zero/localdb"

	"github.com/axis-cash/go-axis-import/c_type"
	"github.com/axis-cash/go-axis/zero/txs/pkg"
	"github.com/axis-cash/go-axis/zero/txs/stx"
	"github.com/axis-cash/go-axis/zero/txs/zstate/pkgstate/data"

	"github.com/axis-cash/go-axis/zero/utils"

	"github.com/axis-cash/go-axis/zero/txs/zstate/tri"
)

type PkgState struct {
	tri tri.Tri
	rw  *sync.RWMutex
	num uint64

	data      data.Data
	snapshots utils.Snapshots
}

func NewPkgState(tri tri.Tri, num uint64) (state PkgState) {
	state = PkgState{tri: tri, num: num}
	state.data = *data.NewData()
	state.rw = new(sync.RWMutex)
	state.data.Clear()
	state.load()
	return
}

func (self *PkgState) Snapshot(revid int) {
	self.snapshots.Push(revid, &self.data)
}
func (self *PkgState) Revert(revid int) {
	self.data.Clear()
	self.data = *self.snapshots.Revert(revid).(*data.Data)
	return
}

func (self *PkgState) load() {
}

func (self *PkgState) Update() {
	self.data.SaveState(self.tri)
	return
}

func (self *PkgState) RecordState(putter axisdb.Putter, hash *c_type.Uint256) {
	self.data.RecordState(putter, hash)
}

func (self *PkgState) GetPkgByHash(hash *c_type.Uint256) (ret *localdb.ZPkg) {
	ret = self.data.GetPkgByHash(self.tri, hash)
	return
}

func (self *PkgState) GetPkgById(id *c_type.Uint256) (ret *localdb.ZPkg) {
	ret = self.data.GetPkgById(self.tri, id)
	return
}

func (state *PkgState) GetPkgHashes() (ret []c_type.Uint256) {
	return state.data.GetHashes()
}

func (self *PkgState) Force_del(hash *c_type.Uint256, close *stx.PkgClose) (e error) {
	self.rw.Lock()
	defer self.rw.Unlock()
	if pg := self.data.GetPkgById(self.tri, &close.Id); pg == nil || pg.Closed {
		e = fmt.Errorf("Close Pkg is nil: %v", hexutil.Encode(close.Id[:]))
		return
	} else {
		if c_superzk.VerifyPKr_X(hash, &close.Sign, &pg.Pack.PKr) {
			pg.Closed = true
			self.data.Add(pg)
		} else {
			e = fmt.Errorf("Close Pkg signed error: %v", hexutil.Encode(close.Id[:]))
			return
		}
		return
	}
}

func (self *PkgState) Force_add(from *c_type.PKr, pack *stx.PkgCreate) (e error) {
	self.rw.Lock()
	defer self.rw.Unlock()

	if pg := self.data.GetPkgById(self.tri, &pack.Id); pg != nil {
		e = fmt.Errorf("Create Pkg is not nil: %v", hexutil.Encode(pack.Id[:]))
		return
	} else {
		zpkg := localdb.ZPkg{
			self.num,
			*from,
			pack.Clone(),
			false,
		}
		self.data.Add(&zpkg)
		return
	}

}

func (self *PkgState) Force_transfer(hash *c_type.Uint256, trans *stx.PkgTransfer) (e error) {
	self.rw.Lock()
	defer self.rw.Unlock()
	if pg := self.data.GetPkgById(self.tri, &trans.Id); pg == nil || pg.Closed {
		e = fmt.Errorf("Transfer Pkg is nil: %v", hexutil.Encode(trans.Id[:]))
		return
	} else {
		if c_superzk.VerifyPKr_X(hash, &trans.Sign, &pg.Pack.PKr) {
			pg.Pack.PKr = trans.PKr
			self.data.Add(pg)
		} else {
			e = fmt.Errorf("Transfer Pkg signed error: %v", hexutil.Encode(trans.Id[:]))
			return
		}
		return
	}
}

type OPkg struct {
	Z localdb.ZPkg
	O pkg.Pkg_O
}

func (self *PkgState) Close(id *c_type.Uint256, pkr *c_type.PKr, key *c_type.Uint256) (ret OPkg, e error) {
	self.rw.Lock()
	defer self.rw.Unlock()
	if pg := self.data.GetPkgById(self.tri, id); pg == nil || pg.Closed {
		e = fmt.Errorf("Close Pkg is nil: %v", hexutil.Encode(id[:]))
		return
	} else {
		if pg.Pack.PKr != *pkr {
			e = fmt.Errorf("Close Pkg Owner Check Failed: %v", hexutil.Encode(id[:]))
			return
		} else {
			if ret.O, e = pkg.DePkg(key, &pg.Pack.Pkg); e != nil {
				return
			} else {
				ret.Z = *pg
				if e = pkg.ConfirmPkg(&ret.O, &ret.Z.Pack.Pkg); e != nil {
					return
				} else {
					pg.Closed = true
					self.data.Add(pg)
					return
				}
			}
		}
	}
}

func (self *PkgState) Transfer(id *c_type.Uint256, pkr *c_type.PKr, to *c_type.PKr) (e error) {
	self.rw.Lock()
	defer self.rw.Unlock()
	if pg := self.data.GetPkgById(self.tri, id); pg == nil || pg.Closed {
		e = fmt.Errorf("Transfer Pkg is nil: %v", hexutil.Encode(id[:]))
		return
	} else {
		if pg.Pack.PKr != *pkr {
			e = fmt.Errorf("Transfer Pkg Owner Check Failed: %v", hexutil.Encode(id[:]))
			return
		} else {
			pg.Pack.PKr = *to
			self.data.Add(pg)
			return
		}
	}
}
