package ethapi

import (
	"fmt"
	"os"
	"testing"

	"github.com/axis-cash/go-axis-import/c_type"
	"github.com/axis-cash/go-axis-import/cpt"
	"github.com/axis-cash/go-axis-import/superzk"
	"github.com/axis-cash/go-axis/common/address"
	"github.com/axis-cash/go-axis/common/hexutil"
	"github.com/axis-cash/go-axis/crypto"
)

func TestMain(m *testing.M) {
	cpt.ZeroInit_NoCircuit()
	os.Exit(m.Run())
}

func Test_getPoolId(t *testing.T) {
	tk := address.Base58ToAccount("3fCJhSjsGJPPB3tSqbycBbwyTahv1WAz8RJY7fpVBqr3mNTLL7NfejjtEywp7jvN3r4isHrh16hrvV8exqGYW4FM")
	pk := address.Base58ToAccount("3fCJhSjsGJPPB3tSqbycBbwyTahv1WAz8RJY7fpVBqr44A7foQAZjWssGXHjc7uVofYCx5cNkmV3k2kEJWU97nKY")
	randHash := crypto.Keccak256Hash(tk[:])
	var rand c_type.Uint256
	copy(rand[:], randHash[:])
	pkr := superzk.Pk2PKr(pk.ToUint512(), &rand)
	id := crypto.Keccak256Hash(pkr[:])
	fmt.Println(hexutil.Encode(id[:]))
}
