package stx

import (
	"bufio"
	"bytes"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/axis-cash/go-axis-import/c_type"

	"github.com/axis-cash/go-axis/zero/utils"

	"github.com/axis-cash/go-axis/rlp"
)

func TestLoad(t *testing.T) {
	h := atomic.Value{}

	v, ok := h.Load().(c_type.Uint256)
	fmt.Println(v, ok)

	h.Store(c_type.RandUint256())
	v, ok = h.Load().(c_type.Uint256)
	fmt.Println(v, ok)
}

func TestRLP(t *testing.T) {
	buf := bytes.Buffer{}
	w := bufio.NewWriter(&buf)

	tx := T{}
	tx.Fee.Value = utils.NewU256(2)
	tx.Fee.Currency = utils.CurrencyToUint256("AXIS")
	tx.Desc_Cmd.RegistPool = &RegistPoolCmd{}
	tx.Desc_Cmd.RegistPool.Value = utils.NewU256(3)
	tx.Desc_Cmd.RegistPool.FeeRate = 10

	e := rlp.Encode(w, &tx)
	fmt.Println(e)
	w.Flush()

	dtx := T{}
	stream := rlp.NewStream(&buf, uint64(buf.Len()))
	_, size, _ := stream.Kind()
	fmt.Println(size)
	e = stream.Decode(&dtx)
	fmt.Println(e)
	fmt.Println(dtx)
}

func TestClose(t *testing.T) {
	buf := bytes.Buffer{}
	w := bufio.NewWriter(&buf)

	tx := T{}
	tx.Fee.Value = utils.NewU256(2)
	tx.Fee.Currency = utils.CurrencyToUint256("AXIS")
	tx.Desc_Cmd.ClosePool = &ClosePoolCmd{}

	e := rlp.Encode(w, &tx)
	fmt.Println(e)
	w.Flush()

	dtx := T{}
	stream := rlp.NewStream(&buf, uint64(buf.Len()))
	_, size, _ := stream.Kind()
	fmt.Println(size)
	e = stream.Decode(&dtx)
	fmt.Println(e)
	fmt.Println(dtx)
}
