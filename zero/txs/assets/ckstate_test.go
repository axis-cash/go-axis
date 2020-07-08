package assets

import (
	"testing"

	"github.com/axis-cash/go-axis-import/c_type"

	"github.com/axis-cash/go-axis/zero/utils"
)

var axis_token = Token{
	utils.CurrencyToUint256("AXIS"),
	utils.NewU256(100),
}

var tk_ticket = Ticket{
	utils.CurrencyToUint256("TK"),
	c_type.RandUint256(),
}

var token_asset = Asset{&axis_token, nil}

var ticket_asset = Asset{nil, &tk_ticket}

var asset = Asset{&axis_token, &tk_ticket}

func TestCkState_OutPlus(t *testing.T) {
	ck := NewCKState(true, &axis_token)
	if ck.Check() == nil {
		t.Fail()
	}
	ck.AddIn(&token_asset)
	if ck.Check() != nil {
		t.Fail()
	}
	ck.AddIn(&ticket_asset)
	if ck.Check() == nil {
		t.Fail()
	}

	ck.AddOut(&asset)
	if ck.Check() == nil {
		t.Fail()
	}

	tkns, tkts := ck.GetList()
	if len(tkns) != 1 {
		t.Fail()
	}
	if len(tkts) != 0 {
		t.Fail()
	}
	ck.AddIn(&token_asset)
	if ck.Check() != nil {
		t.Fail()
	}
	tkns, tkts = ck.GetList()
	if len(tkns) != 0 {
		t.Fail()
	}
	if len(tkts) != 0 {
		t.Fail()
	}

}

func TestCkState_InPlus(t *testing.T) {
	ck := NewCKState(false, &axis_token)
	if ck.Check() == nil {
		t.Fail()
	}
	ck.AddIn(&token_asset)
	if ck.Check() != nil {
		t.Fail()
	}
	ck.AddIn(&ticket_asset)
	if ck.Check() == nil {
		t.Fail()
	}

	tkns, tkts := ck.GetList()
	if len(tkns) != 0 {
		t.Fail()
	}
	if len(tkts) != 1 {
		t.Fail()
	}

	ck.AddOut(&asset)
	if ck.Check() == nil {
		t.Fail()
	}

	ck.AddIn(&token_asset)
	if ck.Check() != nil {
		t.Fail()
	}
	tkns, tkts = ck.GetList()
	if len(tkns) != 0 {
		t.Fail()
	}
	if len(tkts) != 0 {
		t.Fail()
	}

}
