package types

import (
	"github.com/axis-cash/go-axis-import/c_type"
	"github.com/axis-cash/go-axis/common"
	"github.com/axis-cash/go-axis/crypto/sha3"
	"github.com/axis-cash/go-axis/rlp"
)

type Lottery struct {
	ParentHash common.Hash
	ParentNum  uint64
	PosHash    common.Hash
}

type Vote struct {
	Idx       uint32
	ParentNum uint64
	ShareId   common.Hash
	PosHash   common.Hash
	IsPool    bool
	Sign      c_type.Uint512
}

func (s Vote) Hash() common.Hash {

	hw := sha3.NewKeccak256()
	hash := common.Hash{}
	rlp.Encode(hw, []interface{}{
		s.Idx,
		s.ParentNum,
		s.ShareId,
		s.PosHash,
		s.IsPool,
		s.Sign,
	})
	hw.Sum(hash[:0])
	return hash
}

type HeaderVote struct {
	Id     common.Hash
	IsPool bool
	Sign   c_type.Uint512
}
