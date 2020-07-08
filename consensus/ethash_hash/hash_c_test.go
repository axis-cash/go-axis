package ethash_hash

import (
	"fmt"
	"testing"

	"github.com/axis-cash/go-axis-import/c_type"
	"github.com/axis-cash/go-axis/common/hexutil"
)

func TestHash(t *testing.T) {
	k := c_type.RandUint256()
	o := Miner_Hash_0(k[:], 0)
	fmt.Print(hexutil.Encode(o))
}
