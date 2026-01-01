package chain

/*
BLOCKCHAIN – STRUCTURAL CONTAINER

This file defines:
- how blocks are stored
- how they are linked
- how the current chain head is tracked

It does NOT:
- validate blocks
- resolve forks
- enforce consensus
*/

//
// Blockchain represents a simple linear chain of blocks.
//
type Blockchain struct {
	Blocks []*Block // ordered list of blocks
	UTXO   *UTXOSet // current ledger state (derived)
}

//
// NewBlockchain creates a new blockchain with a genesis block.
//
func NewBlockchain(genesis *Block) *Blockchain {

	utxo := NewUTXOSet()

	// Apply genesis transactions to UTXO set
	for _, tx := range genesis.Transactions {
		utxo.ApplyTransaction(&tx)
	}

	return &Blockchain{
		Blocks: []*Block{genesis},
		UTXO:   utxo,
	}
}

//
// Tip returns the latest block in the chain.
//
func (bc *Blockchain) Tip() *Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

//
// Height returns the number of blocks in the chain.
//
func (bc *Blockchain) Height() int {
	return len(bc.Blocks)
}

//
// AddBlock appends a block to the chain.
//
// IMPORTANT:
// - This function does NOT validate the block
// - Validation will be added later
//
func (bc *Blockchain) AddBlock(block *Block) {

	// Apply all transactions to UTXO state
	for _, tx := range block.Transactions {
		bc.UTXO.ApplyTransaction(&tx)
	}

	bc.Blocks = append(bc.Blocks, block)
}
