package chain

import (
	"time"
)

type Transaction struct {
	ID        string   `json:"id"`        // Hash of canonical inputs+outputs
	Inputs    []TxIn   `json:"inputs"`   // UTXOs being spent
	Outputs   []TxOut  `json:"outputs"`  // New UTXOs being created
	Signature string   `json:"signature"` // ECDSA signature (hex-encoded)
	PubKey    string   `json:"pubkey"`    // Public key of signer (hex-encoded)
	Timestamp int64    `json:"timestamp"` // Creation time (Unix timestamp)
}

func NewTransaction(inputs []TxIn, outputs []TxOut) (*Transaction, error) {
	tx := &Transaction{
		Inputs:    inputs,
		Outputs:   outputs,
		Timestamp: time.Now().Unix(),
	}

	id, err := ComputeTxID(tx)
	if err != nil {
		return nil, err
	}
	tx.ID = id

	return tx, nil
}