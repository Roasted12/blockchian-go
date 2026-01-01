package chain

import (
	"errors"
	"sync"
)

/*
MEMPOOL

The mempool stores transactions that:
- are structurally valid
- are signed correctly
- are NOT yet included in a block

It is:
- temporary
- local to a node
- NOT consensus-critical
*/

//
// Mempool represents an in-memory pool of pending transactions.
//
type Mempool struct {
	mu  sync.Mutex
	txs map[string]*Transaction // txID → transaction
}

//
// NewMempool creates an empty mempool.
//
func NewMempool() *Mempool {
	return &Mempool{
		txs: make(map[string]*Transaction),
	}
}

//
// AddTransaction inserts a transaction into the mempool.
//
// IMPORTANT:
// - This function assumes the transaction has already been validated
// - Validation logic stays outside the mempool
//
func (mp *Mempool) AddTransaction(tx *Transaction) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if _, exists := mp.txs[tx.ID]; exists {
		return errors.New("transaction already in mempool")
	}

	mp.txs[tx.ID] = tx
	return nil
}

//
// RemoveTransaction removes a transaction from the mempool.
//
func (mp *Mempool) RemoveTransaction(txID string) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	delete(mp.txs, txID)
}

//
// GetTransactions returns all pending transactions.
//
func (mp *Mempool) GetTransactions() []*Transaction {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	var result []*Transaction
	for _, tx := range mp.txs {
		result = append(result, tx)
	}
	return result
}

//
// Size returns the number of transactions in the mempool.
//
func (mp *Mempool) Size() int {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	return len(mp.txs)
}

//
// Clear removes all transactions from the mempool.
//
func (mp *Mempool) Clear() {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	mp.txs = make(map[string]*Transaction)
}
