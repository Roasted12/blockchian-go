package chain

type Blockchain struct {
	Blocks []*Block // ordered list of blocks
	UTXO   *UTXOSet // current ledger state (derived)
}

func NewBlockchain(genesis *Block) *Blockchain {

	utxo := NewUTXOSet()

	for _, tx := range genesis.Transactions {
		if tx.ID == "" {
		}
		utxo.ApplyTransaction(&tx)
	}

	return &Blockchain{
		Blocks: []*Block{genesis},
		UTXO:   utxo,
	}
}

func (bc *Blockchain) Tip() *Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

func (bc *Blockchain) Height() int {
	return len(bc.Blocks)
}

func (bc *Blockchain) AddBlock(block *Block) {

	for _, tx := range block.Transactions {
		bc.UTXO.ApplyTransaction(&tx)
	}

	bc.Blocks = append(bc.Blocks, block)
}
