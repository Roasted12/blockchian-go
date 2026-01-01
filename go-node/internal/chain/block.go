package chain

import (
	"encoding/json"
	"time"

	"ai-blockchain/go-node/internal/crypto"
)

/*
BLOCK – CONSENSUS CONTAINER

A block does NOT:
- decide if transactions are valid (that’s validation.go)
- decide ownership (that’s UTXO)
- decide signatures (that’s crypto)

A block ONLY:
- groups transactions
- commits to them cryptographically
- links to the previous block
*/

//
// Block represents a single block in the blockchain.
//
type Block struct {
	Index       int           `json:"index"`        // position in the chain
	Timestamp   int64         `json:"timestamp"`    // block creation time
	PrevHash    string        `json:"prevHash"`     // hash of previous block
	MerkleRoot  string        `json:"merkleRoot"`   // commitment to transactions
	Transactions []Transaction `json:"transactions"`
	Hash        string        `json:"hash"`         // hash of this block
	Nonce       int64         `json:"nonce"`        // used later for PoW / PoA
}

//
// NewBlock creates a new block from transactions.
//
func NewBlock(
	index int,
	prevHash string,
	txs []Transaction,
) *Block {

	block := &Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		PrevHash:     prevHash,
		Transactions: txs,
		Nonce:        0, // will matter when we add consensus
	}

	// Step 1: compute Merkle root of transactions
	block.MerkleRoot = block.computeMerkleRoot()

	// Step 2: compute block hash
	block.Hash = block.ComputeHash()

	return block
}

//
// ComputeHash computes the block hash (public method).
//
// This is a public wrapper around computeHash().
// It's used by:
// - Mining functions (need to recompute hash with different nonces)
// - Validation functions (need to verify hash matches block data)
//
func (b *Block) ComputeHash() string {
	return b.computeHash()
}

//
// computeMerkleRoot commits to all transactions in the block.
//
// Why this exists:
// - Allows efficient verification of tx inclusion
// - Any tx modification changes the root
//
func (b *Block) computeMerkleRoot() string {

	var txIDs []string
	for _, tx := range b.Transactions {
		txIDs = append(txIDs, tx.ID)
	}

	return crypto.MerkleRoot(txIDs)
}

//
// computeHash computes the block hash.
//
// What we hash:
// - index (block height)
// - timestamp (block creation time)
// - previous hash (chain linkage)
// - merkle root (commitment to all transactions)
// - nonce (for Proof-of-Work)
//
// Why transactions themselves are NOT hashed here:
// - Merkle root already commits to them
// - Including full transactions would make hashing expensive
// - Merkle root allows efficient inclusion proofs
//
// Serialization format:
// - We use JSON for deterministic serialization
// - This ensures same block always produces same hash
// - JSON is human-readable and cross-language compatible
//
func (b *Block) computeHash() string {
	// Create a struct for hashing (excludes transactions and hash itself)
	// This ensures deterministic serialization
	hashData := struct {
		Index      int    `json:"index"`
		Timestamp  int64  `json:"timestamp"`
		PrevHash   string `json:"prevHash"`
		MerkleRoot string `json:"merkleRoot"`
		Nonce      int64  `json:"nonce"`
	}{
		Index:      b.Index,
		Timestamp:  b.Timestamp,
		PrevHash:   b.PrevHash,
		MerkleRoot: b.MerkleRoot,
		Nonce:      b.Nonce,
	}

	// Serialize to JSON
	// JSON encoding is deterministic when fields are ordered
	// This ensures same block always produces same hash
	data, err := json.Marshal(hashData)
	if err != nil {
		// This should never happen, but handle it gracefully
		// Fallback to simple string concatenation
		data = []byte(
			"index:" + string(rune(b.Index)) +
				"timestamp:" + string(rune(b.Timestamp)) +
				"prevHash:" + b.PrevHash +
				"merkleRoot:" + b.MerkleRoot +
				"nonce:" + string(rune(b.Nonce)),
		)
	}

	return crypto.SHA256(data)
}
