package types_test

import (
	"strings"
	"testing"

	"github.com/pmatseykanets/mempool/types"
)

func TestTxHash(t *testing.T) {
	src := "005A6D600AA58ED5C2DCF508CC1CC116C24A4901FAD34383B1D87C44F981B2C8"
	hash := new(types.TxHash)

	err := hash.SetString(src)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := src, hash.String(); !strings.EqualFold(want, got) {
		t.Errorf("Expected %s got %s", want, got)
	}
}

func TestTxSignature(t *testing.T) {
	src := "00F212D7BD6AB33E95951CC7B49AA7F62D046A5FEDAF34E9421B28034787D617976E2F9128138DC26EEEB826F088CDC8461E9F329175C9595D9267FF7D44AB50"
	sig := new(types.TxSignature)

	err := sig.SetString(src)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := src, sig.String(); !strings.EqualFold(want, got) {
		t.Errorf("Expected %s got %s", want, got)
	}
}
