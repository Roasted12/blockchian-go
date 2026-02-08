package chain

import (
	"encoding/json"
	"time"

	"ai-blockchain/go-node/internal/crypto"
)

type Block struct {
	Index       int           `json:"index"`        // position in the chain
	Timestamp   int64         `json:"timestamp"`    // block creation time
	PrevHash    string        `json:"prevHash"`     // hash of previous block
	MerkleRoot  string        `json:"merkleRoot"`   // commitment to transactions
	Transactions []Transaction `json:"transactions"`
	Hash        string        `json:"hash"`         // hash of this block
	Nonce       int64         `json:"nonce"`        // used later for PoW / PoA
}

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

	block.MerkleRoot = block.computeMerkleRoot()

	block.Hash = block.ComputeHash()

	return block
}

func (b *Block) ComputeHash() string {
	return b.computeHash()
}

func (b *Block) computeMerkleRoot() string {

	var txIDs []string
	for _, tx := range b.Transactions {
		txIDs = append(txIDs, tx.ID)
	}

	return crypto.MerkleRoot(txIDs)
}

func (b *Block) computeHash() string {
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

	data, err := json.Marshal(hashData)
	if err != nil {
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
