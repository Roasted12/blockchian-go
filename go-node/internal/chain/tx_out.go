package chain

/*
TRANSACTION OUTPUT – VALUE TRANSFER

Design choice:
- Address = hash of public key
- Smaller than storing full public key
- Safer (you don't reveal pubkey until spending)
- Matches real blockchains (Bitcoin, etc.)

Later:
- When spending, pubkey is revealed in transaction
- Hash(pubkey) must match address in the UTXO being spent
*/

//
// TxOut represents a single output in a transaction.
//
// What it means:
// - "Send Amount coins to Address"
// - Address is a hash of the recipient's public key
// - Amount is the value being transferred
//
// Important:
// - Amount must be positive (enforced in validation)
// - Address must be valid (enforced in validation)
// - This output becomes a UTXO after the transaction is confirmed
//
type TxOut struct {
	Address string  `json:"address"` // Hash of recipient's public key
	Amount  float64 `json:"amount"`  // Value in coins (using float64 for precision)
}