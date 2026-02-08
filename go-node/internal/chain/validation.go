package chain

import (
	"errors"
	"fmt"

	"ai-blockchain/go-node/internal/consensus"
	"ai-blockchain/go-node/internal/crypto"
)

func VerifyBlock(block *Block, blockchain *Blockchain, difficulty int) error {
	if len(block.Transactions) == 0 {
		return errors.New("block must contain at least one transaction")
	}

	computedHash := block.ComputeHash()
	if computedHash != block.Hash {
		return errors.New("block hash does not match block data")
	}

	computedMerkleRoot := block.computeMerkleRoot()
	if computedMerkleRoot != block.MerkleRoot {
		return errors.New("merkle root does not match transactions")
	}

	if !consensus.ValidateProofOfWork(block.Hash, difficulty) {
		return errors.New("block does not meet proof-of-work requirement")
	}

	if block.Index > 0 {
		if blockchain.Height() < block.Index {
			return errors.New("previous block not found")
		}

		prevBlock := blockchain.Blocks[block.Index-1]
		if prevBlock.Hash != block.PrevHash {
			return errors.New("previous hash mismatch")
		}

		if block.Index != prevBlock.Index+1 {
			return errors.New("block index is not sequential")
		}
	} else {
		if block.PrevHash != "0" {
			return errors.New("genesis block must have previous hash '0'")
		}
	}

	tempUTXO := NewUTXOSet()

	for i, tx := range block.Transactions {
		if err := VerifyTransaction(&tx, tempUTXO); err != nil {
			return fmt.Errorf("transaction %d invalid: %w", i, err)
		}

		tempUTXO.ApplyTransaction(&tx)
	}

	return nil
}

func VerifyTransaction(tx *Transaction, utxo *UTXOSet) error {

	computedID, err := ComputeTxID(tx)
	if err != nil {
		return err
	}

	if computedID != tx.ID {
		return errors.New("transaction ID mismatch")
	}

	seenInputs := make(map[UTXOKey]bool)

	for _, in := range tx.Inputs {
		key := UTXOKey{
			TxID:  in.TxID,
			Index: in.Index,
		}

		if seenInputs[key] {
			return fmt.Errorf("duplicate input detected: %+v", key)
		}
		seenInputs[key] = true
	}

	var inputSum float64

	for _, in := range tx.Inputs {
		key := UTXOKey{
			TxID:  in.TxID,
			Index: in.Index,
		}

		out, ok := utxo.Get(key)
		if !ok {
			return fmt.Errorf("referenced UTXO not found: %+v", key)
		}

		inputSum += out.Amount
	}

	var outputSum float64
	for _, out := range tx.Outputs {
		if out.Amount <= 0 {
			return errors.New("output amount must be positive")
		}
		outputSum += out.Amount
	}

	if outputSum > inputSum {
		return errors.New("output value exceeds input value")
	}

	canonicalBytes, err := CanonicalTxBytes(tx)
	if err != nil {
		return fmt.Errorf("failed to compute canonical bytes: %w", err)
	}

	ok, err := crypto.VerifySignature(canonicalBytes, tx.Signature, tx.PubKey)
	if err != nil {
		return fmt.Errorf("signature verification error: %w", err)
	}
	if !ok {
		return errors.New("invalid transaction signature")
	}

	return nil
}
