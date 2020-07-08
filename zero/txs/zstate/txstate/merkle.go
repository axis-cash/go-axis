package txstate

import (
	"github.com/axis-cash/go-axis-import/c_superzk"
	"github.com/axis-cash/go-axis-import/c_type"
	"github.com/axis-cash/go-axis/crypto"
	"github.com/axis-cash/go-axis/zero/txs/zstate/merkle"
)

var CzeroAddress = c_type.NewPKrByBytes(crypto.Keccak512(nil))
var CzeroMerkleParam = merkle.NewParam(&CzeroAddress, c_superzk.Czero_combine)

var SzkAddress = c_type.NewPKrByBytes(crypto.Keccak256([]byte("$SuperZK$MerkleTree")))
var SzkMerkleParam = merkle.NewParam(&SzkAddress, c_superzk.Combine)
