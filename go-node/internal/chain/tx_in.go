/*
Each input points to exactly one UTXO

Inputs are validated independently

Multiple inputs can fund one transaction

This design:

Prevents double-spending

Makes balances emergent (not stored explicitly)

*/
package chain

type TxIn struct {
	TxID string `json:"tx_id"`
	Index int `json:"index"`
}