// copyright 2018 The axis.cash Authors
// This file is part of the go-axis library.
//
// The go-axis library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-axis library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-axis library. If not, see <http://www.gnu.org/licenses/>.

package test

import (
	"fmt"
	"testing"

	"github.com/axis-cash/go-axis/zero/txs/verify"
	"github.com/axis-cash/go-axis/zero/wallet/lstate/generate"

	"github.com/axis-cash/go-axis/zero/txs/zstate/txstate"

	"github.com/axis-cash/go-axis/common"
	"github.com/axis-cash/go-axis/zero/txs/zstate"

	"github.com/axis-cash/go-axis/zero/txs/tx"

	"github.com/axis-cash/go-axis-import/c_czero"
	"github.com/axis-cash/go-axis-import/c_type"
	"github.com/axis-cash/go-axis-import/superzk"

	"github.com/axis-cash/go-axis/zero/txs/assets"

	"github.com/axis-cash/go-axis/zero/txs/stx"
	"github.com/axis-cash/go-axis/zero/utils"

	"github.com/axis-cash/go-axis/core/state"

	"github.com/axis-cash/go-axis/axisdb"
)

type Blocks struct {
	ca  state.Database
	sd  *state.StateDB
	st  *zstate.ZState
	st0 *txstate.State
	st1 *state1.State1_storage
}

var g_blocks Blocks

func NewBlock() {
	if g_blocks.ca == nil {
		db := axisdb.NewMemDatabase()
		g_blocks.ca = state.NewDatabase(db)
		g_blocks.sd, _ = state.NewGenesis(common.Hash{}, g_blocks.ca)
		g_blocks.st = g_blocks.sd.GetZState()
		g_blocks.st0 = &g_blocks.st.State
	} else {
		g_blocks.st0.Block = txstate.StateBlock{}
	}
}

func EndBlock() {
	if g_blocks.st1 == nil {
		st1 := state1.LoadState(g_blocks.st, "")
		g_blocks.st1 = &st1
	} else {
		g_blocks.st1.State = g_blocks.st
	}
	g_blocks.st1.UpdateWitness(keys.Seeds2Tks(seeds))
	NewBlock()
}

type user struct {
	i    int
	seed c_type.Uint256
	addr c_type.Uint512
}

var seeds = []c_type.Uint256{}

func newUser(i int) (ret user) {
	fmt.Printf("\n\n===========new user(%v)============\n", i)
	ret = user{}
	ret.i = i
	ret.seed = c_type.Uint256{byte(i)}
	ret.addr = keys.Seed2Addr(&ret.seed)
	seeds = append(seeds, ret.seed)
	fmt.Printf("\nseed: ")
	ret.seed.LogOut()
	fmt.Printf("\naddr: ")
	ret.addr.LogOut()
	return
}

func (self *user) getAR() (pkr c_type.PKr) {
	pkr = superzk.Pk2PKr(&self.addr, nil)
	fmt.Printf("\nuser(%v):get pkr: ", self.i)
	pkr.LogOut()
	return
}

func (self *user) addOut(v int) {
	out := stx.Out_O{}
	out.Addr = self.getAR()
	out.Asset = assets.NewAsset(
		&assets.Token{
			utils.CurrencyToUint256("AXIS"),
			utils.NewU256(uint64(v)),
		},
		nil,
	)
	g_blocks.st.addOut_O(&out)
	g_blocks.st.Update()
	EndBlock()
}

func (self *user) addTkt(v int) {
	out := stx.Out_O{}
	out.Addr = self.getAR()
	out.Asset = assets.Asset{
		&assets.Token{
			utils.CurrencyToUint256("AXIS"),
			utils.NewU256(uint64(v)),
		},
		&assets.Ticket{
			utils.CurrencyToUint256("AXIS_TICKET"),
			cpt.Random(),
		},
	}
	g_blocks.st.addOut_O(&out)
	g_blocks.st.Update()
	EndBlock()
}

func (self *user) GetOuts() (outs []*state1.OutState) {
	if os, e := g_blocks.st1.GetOuts(keys.Seed2Tk(&self.seed).NewRef()); e != nil {
		panic(e)
		return
	} else {
		outs = os
		return
	}
}

func (self *user) GetPkgs(is_from bool) (ret []*state1.Pkg) {
	ret = g_blocks.st1.GetPkgs(keys.Seed2Tk(&self.seed).NewRef(), is_from)
	return
}

func (self *user) Gen(seed *c_type.Uint256, t *tx.T) (s stx.T, e error) {
	return generate.Gen_lstate(g_blocks.st1, seed, t)
}

func (self *user) Verify(t *stx.T) (e error) {
	return verify.Verify_state1(t, g_blocks.st1.State)
}

func (self *user) Logout() (ret uint64) {
	fmt.Printf("\n\n===========user(%v)============\n", self.i)
	outs := self.GetOuts()
	for _, out := range outs {
		if out.Out_O.Asset.Tkn != nil {
			fmt.Printf("TKN: (%v:%v)---%v-----%v\n", out.Root[1], out.OutIndex, out.Out_O.Asset.Tkn.Currency[0], out.Out_O.Asset.Tkn.Value.ToIntRef().Int64())
			ret += out.Out_O.Asset.Tkn.Value.ToIntRef().Uint64()
		}
		if out.Out_O.Asset.Tkt != nil {
			fmt.Printf("TKT: (%v:%v)---%v-----%v\n", out.Root[1], out.OutIndex, out.Out_O.Asset.Tkt.Category[0], out.Out_O.Asset.Tkt.Value)
		}
	}
	fmt.Printf("===========user(%v)============\n\n", self.i)
	return
}

func (self *user) Close(id *c_type.Uint256, v int, key *c_type.Uint256) {
	fmt.Printf("user(%v) close pkg %v\n", self.i, id)
	t := tx.T{}
	t.Fee = assets.Token{
		utils.CurrencyToUint256("AXIS"),
		utils.NewU256(uint64(0)),
	}
	t.PkgClose = &tx.PkgClose{}
	t.PkgClose.Key = *key
	t.PkgClose.Id = *id

	out1 := tx.Out{}
	out1.Asset = assets.Asset{
		&assets.Token{
			utils.CurrencyToUint256("AXIS"),
			utils.NewU256(uint64(v)),
		},
		nil,
	}
	out1.IsZ = true
	out1.Addr = self.getAR()
	t.Outs = append(t.Outs, out1)

	s, e := self.Gen(&self.seed, &t)
	if e != nil {
		fmt.Printf("user(%v) send gen error: %v", self.i, e)
	}

	if e := self.Verify(&s); e != nil {
		fmt.Printf("user(%v) send verify error: %v", self.i, e)
	}

	g_blocks.st.AddStx(&s)
	g_blocks.st.Update()
	EndBlock()
	return
}

func (self *user) Package(v int, fee int, u user) (ret c_type.PKr) {
	fmt.Printf("user(%v) send %v:%v to user(%v)\n", self.i, v, fee, u.i)
	outs := self.GetOuts()
	in := tx.In{}
	in.Root = outs[0].Root
	out0 := tx.PkgCreate{}
	out0.PKr = u.getAR()
	ret = out0.PKr
	out0.Pkg.Asset = assets.Asset{
		&assets.Token{
			utils.CurrencyToUint256("AXIS"),
			utils.NewU256(uint64(v)),
		},
		nil,
	}

	out1 := tx.Out{}
	out1.Addr = self.getAR()
	out1.Asset = outs[0].Out_O.Asset.Clone()
	out1.Asset.Tkn.Value.SubU(utils.NewU256(uint64(v)).ToRef())
	out1.Asset.Tkn.Value.SubU(utils.NewU256(uint64(fee)).ToRef())

	out1.IsZ = true

	t := tx.T{}
	t.Fee = assets.Token{
		utils.CurrencyToUint256("AXIS"),
		utils.NewU256(uint64(fee)),
	}
	t.Ins = append(t.Ins, in)
	t.Outs = append(t.Outs, out1)
	t.PkgCreate = &out0

	s, e := self.Gen(&self.seed, &t)
	if e != nil {
		fmt.Printf("user(%v) send gen error: %v", self.i, e)
	}

	if e := self.Verify(&s); e != nil {
		fmt.Printf("user(%v) send verify error: %v", self.i, e)
	}

	g_blocks.st.AddStx(&s)
	g_blocks.st.Update()
	EndBlock()
	return
}

func (self *user) Send(v int, fee int, u user, z bool) {
	fmt.Printf("user(%v) send %v:%v to user(%v)\n", self.i, v, fee, u.i)
	outs := self.GetOuts()
	in := tx.In{}
	in.Root = outs[0].Root
	out0 := tx.Out{}
	out0.Addr = u.getAR()
	out0.Asset = assets.Asset{
		&assets.Token{
			utils.CurrencyToUint256("AXIS"),
			utils.NewU256(uint64(v)),
		},
		nil,
	}
	out0.IsZ = z

	out1 := tx.Out{}
	out1.Addr = self.getAR()
	out1.Asset = outs[0].Out_O.Asset.Clone()
	out1.Asset.Tkn.Value.SubU(utils.NewU256(uint64(v)).ToRef())
	out1.Asset.Tkn.Value.SubU(utils.NewU256(uint64(fee)).ToRef())

	out1.IsZ = z

	t := tx.T{}
	t.Fee = assets.Token{
		utils.CurrencyToUint256("AXIS"),
		utils.NewU256(uint64(fee)),
	}
	t.Ins = append(t.Ins, in)
	t.Outs = append(t.Outs, out0)
	t.Outs = append(t.Outs, out1)

	s, e := self.Gen(&self.seed, &t)
	if e != nil {
		fmt.Printf("user(%v) send gen error: %v", self.i, e)
	}

	if e := self.Verify(&s); e != nil {
		fmt.Printf("user(%v) send verify error: %v", self.i, e)
	}

	g_blocks.st.AddStx(&s)
	g_blocks.st.Update()
	EndBlock()
}

type ArrayObj struct {
	p0 int
	p1 uint64
}

func TestArrayObj(t *testing.T) {
	a := [6]ArrayObj{}
	a[1].p0 = 0
	fmt.Printf("%v\n", a)
}

func TestMain(m *testing.M) {
	cpt.ZeroInit("", c_czero.NET_Dev)
	NewBlock()
	m.Run()
}

func TestTx(t *testing.T) {
	user_m := newUser(1)
	user_a := newUser(2)
	user_m.addOut(100)
	user_m.Send(50, 10, user_a, true)
	if user_m.Logout() != 40 {
		t.Fail()
	}
}

func TestPkg(t *testing.T) {
	user_m := newUser(1)
	user_a := newUser(2)
	user_m.addOut(100)
	user_m.Package(50, 10, user_a)
	if user_m.Logout() != 40 {
		t.FailNow()
	}
	user_a.Close(&c_type.Uint256{}, 50, &g_blocks.st1.GetPkgs(nil, true)[0].Key)
	if user_a.Logout() != 50 {
		t.FailNow()
	}
}

func TestTxs(t *testing.T) {
	user_m := newUser(1)
	user_a := newUser(2)
	user_b := newUser(3)
	user_c := newUser(4)

	user_m.addTkt(100)
	user_m.addOut(100)
	user_m.addOut(100)
	user_m.addOut(100)

	if user_m.Logout() != 400 {
		t.FailNow()
	}

	pkg_pkr := user_m.Package(50, 10, user_a)
	if g_blocks.st.Pkgs.GetPkg(&c_type.Uint256{}) == nil {
		t.FailNow()
	}

	var key c_type.Uint256
	if pkgs := user_m.GetPkgs(true); len(pkgs) == 0 {
		t.FailNow()
	} else {
		key = pkgs[0].Key
	}
	if pkgs := user_a.GetPkgs(false); len(pkgs) == 0 {
		t.FailNow()
	}

	g_blocks.st.Pkgs.Close(&c_type.Uint256{}, &pkg_pkr, &key)
	g_blocks.st.Update()
	EndBlock()

	if g_blocks.st.Pkgs.GetPkg(&c_type.Uint256{}) != nil {
		t.FailNow()
	}
	if pkgs := user_m.GetPkgs(true); len(pkgs) != 0 {
		t.FailNow()
	}
	if pkgs := user_a.GetPkgs(false); len(pkgs) != 0 {
		t.FailNow()
	}

	EndBlock()

	user_a.addOut(50)

	if user_m.Logout() != 340 {
		t.FailNow()
	}
	if user_a.Logout() != 50 {
		t.FailNow()
	}

	user_m.addOut(100)

	if user_m.Logout() != 440 {
		t.FailNow()
	}
	if user_a.Logout() != 50 {
		t.FailNow()
	}

	user_a.Send(20, 5, user_b, true)

	if user_a.Logout() != 25 {
		t.FailNow()
	}
	if user_b.Logout() != 20 {
		t.FailNow()
	}

	user_b.Send(10, 5, user_c, true)

	if user_b.Logout() != 5 {
		t.FailNow()
	}
	if user_c.Logout() != 10 {
		t.FailNow()
	}
}

func TestStrTree(t *testing.T) {

	outState := txstate.NewMerkleTree(g_blocks.st.Tri)

	cm1 := c_type.Uint256{1}
	outState.AppendLeaf(cm1)

	cm2 := c_type.Uint256{2}
	outState.AppendLeaf(cm2)

	cm3 := c_type.Uint256{3}
	outState.AppendLeaf(cm3)

	cm4 := c_type.Uint256{4}
	rt4 := outState.AppendLeaf(cm4)

	pos, path, anchor := outState.GetPaths(cm3)
	rt := txstate.CalcRoot(&cm3, pos, &path)

	if rt != rt4 {
		t.FailNow()
	}
	if rt != anchor {
		t.FailNow()
	}

	fmt.Print(pos, path)
}
