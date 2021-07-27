package mempool

import (
	"github.com/google/btree"
	"github.com/pmatseykanets/mempool/types"
	"github.com/shopspring/decimal"
)

// trn wraps type.Transaction in order to disambiguate transactions
// with the same fee by considering newer transaction being less then older ones.
// trn implements btree.Item.
type trn struct {
	id uint64 // In real world this should be something like tx timestamp.
	tx *types.Transaction
}

func (t *trn) Less(than btree.Item) bool {
	thanTx := than.(*trn)
	return t.tx.Fee().LessThan(thanTx.tx.Fee()) ||
		t.tx.Fee().Equal(thanTx.tx.Fee()) && t.id > thanTx.id
}

// Mempool implements a fixed capacity/size priority (by fee) mempool.
// It discards already seen transactions and if at capacity
// transactions with the fee less than the lowest seen fee.
type Mempool struct {
	cap       int                       // Capacity.
	seen      map[types.TxHash]struct{} // A cache of seen transactions.
	index     *btree.BTree              // A btree index of transactions by the fee.
	lowestFee decimal.Decimal           // The lowest seen fee.
	lastID    uint64                    // The monotonically increasing internal trn ID counter.
}

// New creates a instance of Mempool with a given capacity.
// It panics if cap < 1.
func New(cap int) *Mempool {
	if cap < 1 {
		panic("mempool: capacity can't be less than 1")
	}

	return &Mempool{
		cap:       cap,
		seen:      make(map[types.TxHash]struct{}),
		index:     btree.New(2),
		lowestFee: decimal.New(0, 0),
	}
}

// Push attempts to add a transaction to the mempool.
func (p *Mempool) Push(tx *types.Transaction) {
	// Reject if we already seen the transaction.
	if _, ok := p.seen[tx.Hash]; ok {
		return
	}
	p.seen[tx.Hash] = struct{}{}

	txFee := tx.Fee()

	if p.Len() >= p.cap {
		// If we're at capacity reject the transaction
		// if it's priority (fee) is lower than or equal to the current lowest in the pool.
		if txFee.LessThanOrEqual(p.lowestFee) {
			return
		}
		// Otherwise drop the transaction with the lowest priority.
		_ = p.index.DeleteMin()
	}
	// Keep track of the lowest priority.
	if p.lowestFee.Equal(decimal.New(0, 0)) || txFee.LessThan(p.lowestFee) {
		p.lowestFee = txFee
	}

	// Add the accepted transaction.
	_ = p.index.ReplaceOrInsert(&trn{id: p.lastID, tx: tx})
	p.lastID++
}

// Len returns the number of transactions in the mempool.
func (p *Mempool) Len() int {
	return p.index.Len()
}

// Cap returns the capacity of the mempool.
func (p *Mempool) Cap() int {
	return p.cap
}

// Iterator allows iterate in-order over transactions in the mempool.
type Iterator func(tx *types.Transaction) bool

// GetPrioritized calls the Iterator for every transaction in the mempool
// ordered by the priority (fee) descending until the iterator returns false.
func (p *Mempool) GetPrioritized(itr Iterator) {
	p.index.Descend(func(i btree.Item) bool {
		return itr(i.(*trn).tx)
	})
}
