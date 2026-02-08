package chain

type UTXOKey struct {
	TxID  string // Transaction hash that created the output
	Index int    // Index of the output inside that transaction
}

type UTXOSet struct {
	store map[UTXOKey]TxOut
}

func NewUTXOSet() *UTXOSet {
	return &UTXOSet{
		store: make(map[UTXOKey]TxOut),
	}
}

func (u *UTXOSet) Get(key UTXOKey) (TxOut, bool) {
	out, ok := u.store[key]
	return out, ok
}

func (u *UTXOSet) Spend(key UTXOKey) {
	delete(u.store, key)
}

func (u *UTXOSet) Add(txid string, index int, out TxOut) {
	key := UTXOKey{
		TxID:  txid,
		Index: index,
	}
	u.store[key] = out
}

func (u *UTXOSet) ApplyTransaction(tx *Transaction) {

	for _, in := range tx.Inputs {
		key := UTXOKey{
			TxID:  in.TxID,
			Index: in.Index,
		}
		u.Spend(key)
	}

	for i, out := range tx.Outputs {
		u.Add(tx.ID, i, out)
	}
}

func (u *UTXOSet) BalanceOf(address string) float64 {
	var balance float64
	for _, out := range u.store {
		if out.Address == address {
			balance += out.Amount
		}
	}
	return balance
}

func (u *UTXOSet) FindSpendableOutputs(address string, amount float64) (float64, []UTXOKey) {
	var total float64
	var selected []UTXOKey

	for key, out := range u.store {
		if out.Address != address {
			continue
		}
		selected = append(selected, key)
		total += out.Amount
		if total >= amount {
			break
		}
	}

	return total, selected
}