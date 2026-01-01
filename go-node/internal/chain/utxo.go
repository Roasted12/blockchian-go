package chain

/*
UTXO SET – CORE LEDGER STATE

In a UTXO-based blockchain, there is:
- no account table
- no balance map
- no global "who owns how much" state

The ONLY economic state is:
  "Which transaction outputs are still unspent?"

This file implements that state.
*/

//
// UTXOKey uniquely identifies a single transaction output.
//
// Why this is needed:
// - One transaction can create multiple outputs
// - Each output must be spendable exactly once
// - (txid, output_index) is globally unique
//
type UTXOKey struct {
	TxID  string // Transaction hash that created the output
	Index int    // Index of the output inside that transaction
}

//
// UTXOSet represents the entire spendable state of the blockchain.
//
// Conceptually:
//   UTXOSet = map[(txid, index)] -> TxOut
//
// If this map is correct:
// - balances are correct
// - double spending is impossible
// - transaction validation is simple
//
type UTXOSet struct {
	store map[UTXOKey]TxOut
}

//
// NewUTXOSet creates an empty UTXO set.
//
// Why empty?
// - At genesis, nothing has been created yet
// - Genesis transactions will explicitly add initial UTXOs
//
func NewUTXOSet() *UTXOSet {
	return &UTXOSet{
		store: make(map[UTXOKey]TxOut),
	}
}

//
// Get retrieves an unspent output if it exists.
//
// This answers the critical validation question:
//   "Does this input reference a real, unspent output?"
//
// If ok == false:
// - the output never existed
// - OR it has already been spent
// - either case makes the transaction invalid
//
func (u *UTXOSet) Get(key UTXOKey) (TxOut, bool) {
	out, ok := u.store[key]
	return out, ok
}

//
// Spend removes a UTXO from the set.
//
// This is how double spending is prevented.
//
// Important:
// - We do NOT mark outputs as spent
// - We do NOT decrement balances
// - We literally remove the output from existence
//
// After this:
// - Any transaction trying to spend it again will fail
//
func (u *UTXOSet) Spend(key UTXOKey) {
	delete(u.store, key)
}

//
// Add inserts a new unspent output into the set.
//
// This is called ONLY AFTER a transaction has been:
// - fully validated
// - signature-verified
// - value-conserving
//
// New outputs are the ONLY way value enters the UTXO set.
//
func (u *UTXOSet) Add(txid string, index int, out TxOut) {
	key := UTXOKey{
		TxID:  txid,
		Index: index,
	}
	u.store[key] = out
}

//
// ApplyTransaction updates the UTXO set using a valid transaction.
//
// What this function does conceptually:
//
// 1. Consumes (destroys) all input UTXOs
// 2. Creates new UTXOs from the transaction outputs
//
// This transforms the ledger state from:
//   OLD_STATE -> NEW_STATE
//
func (u *UTXOSet) ApplyTransaction(tx *Transaction) {

	// Step 1: Spend all referenced inputs
	for _, in := range tx.Inputs {
		key := UTXOKey{
			TxID:  in.TxID,
			Index: in.Index,
		}
		u.Spend(key)
	}

	// Step 2: Add newly created outputs
	for i, out := range tx.Outputs {
		u.Add(tx.ID, i, out)
	}
}

//
// BalanceOf computes the balance for an address.
//
// Important:
// - Balances are NOT stored
// - They are derived by scanning the UTXO set
//
// This function is for:
// - wallets
// - explorers
// - debugging
//
// It is NOT used during consensus validation.
//
func (u *UTXOSet) BalanceOf(address string) float64 {
	var balance float64
	for _, out := range u.store {
		if out.Address == address {
			balance += out.Amount
		}
	}
	return balance
}
