package chain

import (
	"bytes"
	"encoding/json"
	"sort"

	"ai-blockchain/go-node/internal/crypto"
)

/*
SERIALIZATION – CANONICAL FORMATS

This file ensures that:
- Transaction hashing is deterministic
- Block hashing is deterministic
- Cross-language compatibility (Java ↔ Go ↔ Python)

Why canonicalization matters:
- Same transaction must always produce same hash
- Different serialization = different hash = broken consensus
*/

//
// txForHash is a helper struct for hashing transactions.
//
// Why separate struct?
// - We exclude signature, pubkey, timestamp from hash
// - Hash is computed BEFORE signing
// - Signature signs the hash, not the full transaction
//
type txForHash struct {
	Inputs  []TxIn  `json:"inputs"`
	Outputs []TxOut `json:"outputs"`
}

//
// CanonicalTxBytes serializes a transaction in a deterministic way.
//
// What this function does:
// 1. Sorts inputs by (TxID, Index) - ensures same order always
// 2. Sorts outputs by Address - ensures same order always
// 3. Serializes to JSON without HTML escaping
// 4. Returns bytes ready for hashing
//
// Why sorting matters:
// - Inputs: [(txA,1), (txA,0)] vs [(txA,0), (txA,1)] must hash the same
// - Outputs: [addrB, addrA] vs [addrA, addrB] must hash the same
//
// This is called:
// - When computing tx.ID (before signing)
// - When verifying signatures (must match exactly)
//
func CanonicalTxBytes(tx *Transaction) ([]byte, error) {
	// Create a copy to avoid mutating the original transaction
	// (sorting modifies the slice in place)
	inputsCopy := make([]TxIn, len(tx.Inputs))
	copy(inputsCopy, tx.Inputs)
	outputsCopy := make([]TxOut, len(tx.Outputs))
	copy(outputsCopy, tx.Outputs)

	// Sort inputs: first by TxID, then by Index
	// This ensures deterministic ordering
	sort.Slice(inputsCopy, func(i, j int) bool {
		if inputsCopy[i].TxID == inputsCopy[j].TxID {
			return inputsCopy[i].Index < inputsCopy[j].Index
		}
		return inputsCopy[i].TxID < inputsCopy[j].TxID
	})

	// Sort outputs by Address (lexicographic)
	// This ensures deterministic ordering
	sort.Slice(outputsCopy, func(i, j int) bool {
		return outputsCopy[i].Address < outputsCopy[j].Address
	})

	// Create the hash-only struct (excludes signature, pubkey, timestamp)
	tmp := txForHash{
		Inputs:  inputsCopy,
		Outputs: outputsCopy,
	}

	// Serialize to JSON
	// SetEscapeHTML(false) ensures consistent output
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	err := enc.Encode(tmp)
	if err != nil {
		return nil, err
	}

	// Remove trailing newline that json.Encoder adds
	data := buf.Bytes()
	if len(data) > 0 && data[len(data)-1] == '\n' {
		data = data[:len(data)-1]
	}

	return data, nil
}

//
// ComputeTxID computes the SHA-256 hash of canonical transaction bytes.
//
// This is the transaction's unique identifier.
//
// Important:
// - ID is computed BEFORE signing
// - Signature signs the hash, not the raw bytes
// - Same inputs/outputs = same ID (deterministic)
//
func ComputeTxID(tx *Transaction) (string, error) {
	canonical, err := CanonicalTxBytes(tx)
	if err != nil {
		return "", err
	}
	return crypto.SHA256(canonical), nil
}