package types

import (
	"encoding/hex"
	"fmt"

	"github.com/shopspring/decimal"
)

// Transaction holds details of a transaction.
type Transaction struct {
	Hash      TxHash
	Signature TxSignature
	FeePerGas decimal.Decimal
	Gas       uint64
}

// Fee retuns the fee accociated with the transaction.
func (t *Transaction) Fee() decimal.Decimal {
	return t.FeePerGas.Mul(decimal.NewFromInt(int64(t.Gas)))
}

const TxHashSize = 32

// TxHash holds Transaction hash.
type TxHash [TxHashSize]byte

func (h *TxHash) SetBytes(src []byte) error {
	if n := len(src); n != TxHashSize {
		return fmt.Errorf("invalid hash size %d", n)
	}

	copy(h[:], src)

	return nil
}

func (h *TxHash) SetString(src string) error {
	if n := len(src); n != TxHashSize*2 {
		return fmt.Errorf("invalid hash string len %d", n)
	}

	b, err := hex.DecodeString(src)
	if err != nil {
		return err
	}

	return h.SetBytes(b)
}

func (h *TxHash) String() string {
	return hex.EncodeToString(h[:])
}

const TxSignatureSize = 64

// TxSignature holds transaction signature.
type TxSignature [TxSignatureSize]byte

func (s *TxSignature) SetBytes(src []byte) error {
	if n := len(src); n != TxSignatureSize {
		return fmt.Errorf("invalid signature size %d", n)
	}

	copy(s[:], src)

	return nil
}

func (s *TxSignature) SetString(src string) error {
	if n := len(src); n != TxSignatureSize*2 {
		return fmt.Errorf("invalid signature string len %d", n)
	}

	b, err := hex.DecodeString(src)
	if err != nil {
		return err
	}

	return s.SetBytes(b)
}

func (s *TxSignature) String() string {
	return hex.EncodeToString(s[:])
}
