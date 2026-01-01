package chain

import (
	"time"
)

/*
TRANSACTION – VALUE TRANSFER UNIT

A transaction represents a transfer of value from inputs to outputs.

Structure:
- Inputs: References to UTXOs being spent
- Outputs: New UTXOs being created
- Signature: Cryptographic proof of ownership
- PubKey: Public key of the signer (revealed when spending)

Lifecycle:
1. Created by wallet (Java side)
2. Signed with private key
3. Broadcast to network
4. Validated by nodes
5. Added to mempool
6. Included in block
7. Confirmed on chain
*/

//
// Transaction represents a single value transfer operation.
//
// What it does:
// - Consumes inputs (destroys old UTXOs)
// - Creates outputs (creates new UTXOs)
// - Proves ownership via signature
//
// Important fields:
// - ID: Hash of inputs+outputs (computed before signing)
// - Signature: Signs the transaction ID
// - PubKey: Public key that corresponds to the private key used to sign
//
type Transaction struct {
	ID        string   `json:"id"`        // Hash of canonical inputs+outputs
	Inputs    []TxIn   `json:"inputs"`   // UTXOs being spent
	Outputs   []TxOut  `json:"outputs"`  // New UTXOs being created
	Signature string   `json:"signature"` // ECDSA signature (hex-encoded)
	PubKey    string   `json:"pubkey"`    // Public key of signer (hex-encoded)
	Timestamp int64    `json:"timestamp"` // Creation time (Unix timestamp)
}

//
// NewTransaction creates a new transaction from inputs and outputs.
//
// This function:
// 1. Creates the transaction structure
// 2. Computes the transaction ID (hash of inputs+outputs)
// 3. Sets timestamp
//
// Note: Transaction is NOT signed yet. Signing happens separately.
//
func NewTransaction(inputs []TxIn, outputs []TxOut) (*Transaction, error) {
	tx := &Transaction{
		Inputs:    inputs,
		Outputs:   outputs,
		Timestamp: time.Now().Unix(),
	}

	// Compute transaction ID (must be done before signing)
	id, err := ComputeTxID(tx)
	if err != nil {
		return nil, err
	}
	tx.ID = id

	return tx, nil
}