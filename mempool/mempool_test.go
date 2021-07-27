package mempool_test

import (
	"reflect"
	"testing"

	"github.com/pmatseykanets/mempool/mempool"
	"github.com/pmatseykanets/mempool/types"
	"github.com/shopspring/decimal"
)

func TestMempool(t *testing.T) {
	tx1 := &types.Transaction{Gas: 729000, FeePerGas: decimal.New(11134106816568039, -17)} // Fee  81167.63869278100431
	tx2 := &types.Transaction{Gas: 834000, FeePerGas: decimal.New(27503931836911927, -17)} // Fee 229382.79151984547118
	tx3 := &types.Transaction{Gas: 303000, FeePerGas: decimal.New(13954792852514392, -17)} // Fee  42283.02234311860776
	tx4 := &types.Transaction{Gas: 179000, FeePerGas: decimal.New(6767094874229113, -16)}  // Fee 121130.9982487011227
	tx5 := &types.Transaction{Gas: 988000, FeePerGas: decimal.New(5453765337483497, -16)}  // Fee 538832.0153433695036
	tx6 := &types.Transaction{Gas: 988000, FeePerGas: decimal.New(5453765337483497, -16)}  // Fee 538832.0153433695036
	tx7 := &types.Transaction{Gas: 113000, FeePerGas: decimal.New(1053652881033623, -16)}  // Fee  11906.2775556799399
	tx8 := &types.Transaction{Gas: 834000, FeePerGas: decimal.New(27503931836911927, -17)} // Fee 229382.79151984547118 same as tx2
	tx1.Hash.SetString("40E10C7CF56A738C0B8AD4EE30EA8008C7B2334B3ADA195083F8CB18BD3911A0")
	tx2.Hash.SetString("4B2B252899DC689106C8FCEA3E24E4AFFC597D2B4E701F99EB8CD909217D323F")
	tx3.Hash.SetString("F75F133F149FDA7DEB391B2446C5196E7C704F45456E69312C310C72893F5B6A")
	tx4.Hash.SetString("16633D0A25ECA886F100A34BA5C43366732836E6E7B140159298C71CF78309F9")
	tx5.Hash.SetString("34CCBBCD977F868B1F46DB18697D7688ED0053C77F52DF643CA4DB5C3982D1FF")
	tx6.Hash.SetString("34CCBBCD977F868B1F46DB18697D7688ED0053C77F52DF643CA4DB5C3982D1FF")
	tx7.Hash.SetString("30F5571C3D010129DEBEE6317A5C0ECDF5AEC74A310065298AE47AA95A177682")
	tx8.Hash.SetString("F906FEA8E835B88635BECB73FF2E3FC628062931D814393AFCDD1FBCB043D77E")

	pool := mempool.New(3)

	checkPrioritized := func(t *testing.T, want []*types.Transaction, desc string) {
		var got []*types.Transaction
		pool.GetPrioritized(func(tx *types.Transaction) bool {
			got = append(got, tx)
			return true
		})

		if !reflect.DeepEqual(want, got) {
			t.Fatalf("%s Expected %v got %v", desc, want, got)
		}
	}

	pool.Push(tx1)
	pool.Push(tx2)
	if want, got := 2, pool.Len(); want != got {
		t.Fatalf("Expected Len %d got %d", want, got)
	}
	checkPrioritized(t, []*types.Transaction{tx2, tx1}, "[1]")

	pool.Push(tx3)
	if want, got := 3, pool.Len(); want != got {
		t.Fatalf("Expected Len %d got %d", want, got)
	}
	checkPrioritized(t, []*types.Transaction{tx2, tx1, tx3}, "[2]")

	pool.Push(tx4)
	if want, got := 3, pool.Len(); want != got {
		t.Fatalf("Expected Len %d got %d", want, got)
	}
	checkPrioritized(t, []*types.Transaction{tx2, tx4, tx1}, "[3]")

	pool.Push(tx5)
	if want, got := 3, pool.Len(); want != got {
		t.Fatalf("Expected Len %d got %d", want, got)
	}
	checkPrioritized(t, []*types.Transaction{tx5, tx2, tx4}, "[4]")

	// Already seen transaction should be rejected.
	pool.Push(tx6)
	if want, got := 3, pool.Len(); want != got {
		t.Fatalf("Expected Len %d got %d", want, got)
	}
	checkPrioritized(t, []*types.Transaction{tx5, tx2, tx4}, "[5]")

	// A transaction with a fee less than the current lowest in the pool should be rejected.
	pool.Push(tx7)
	if want, got := 3, pool.Len(); want != got {
		t.Fatalf("Expected Len %d got %d", want, got)
	}
	checkPrioritized(t, []*types.Transaction{tx5, tx2, tx4}, "[6]")

	// A newer transaction with the same fee as an existing one gets a lower priority.
	pool.Push(tx8)
	if want, got := 3, pool.Len(); want != got {
		t.Fatalf("Expected Len %d got %d", want, got)
	}
	checkPrioritized(t, []*types.Transaction{tx5, tx2, tx8}, "[7]")
}
