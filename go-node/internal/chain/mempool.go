package chain

import (
	"errors"
	"sync"
)

type Mempool struct {
	mu  sync.Mutex
	txs map[string]*Transaction // txID → transaction
}

func NewMempool() *Mempool {
	return &Mempool{
		txs: make(map[string]*Transaction),
	}
}

func (mp *Mempool) AddTransaction(tx *Transaction) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if _, exists := mp.txs[tx.ID]; exists {
		return errors.New("transaction already in mempool")
	}

	mp.txs[tx.ID] = tx
	return nil
}

func (mp *Mempool) RemoveTransaction(txID string) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	delete(mp.txs, txID)
}

func (mp *Mempool) GetTransactions() []*Transaction {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	var result []*Transaction
	for _, tx := range mp.txs {
		result = append(result, tx)
	}
	return result
}

func (mp *Mempool) Size() int {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	return len(mp.txs)
}

func (mp *Mempool) Clear() {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	mp.txs = make(map[string]*Transaction)
}
