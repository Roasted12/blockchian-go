package chain

import (
	"errors"
	"fmt"

	"ai-blockchain/go-node/internal/consensus"
	"ai-blockchain/go-node/internal/crypto"
)

//
// VerifyBlock validates a block before adding it to the blockchain.
//
// This function checks:
// 1. Block structure is valid
// 2. Previous block exists (except genesis)
// 3. Block index is correct
// 4. Block hash is valid
// 5. Proof-of-Work is valid
// 6. Merkle root matches transactions
// 7. All transactions are valid
//
// This is called when:
// - Receiving blocks from other nodes
// - Mining new blocks (before adding to chain)
//
func VerifyBlock(block *Block, blockchain *Blockchain, difficulty int) error {
	// Check 1: Block must have at least one transaction
	if len(block.Transactions) == 0 {
		return errors.New("block must contain at least one transaction")
	}

	// Check 2: Verify block hash matches block data
	computedHash := block.ComputeHash()
	if computedHash != block.Hash {
		return errors.New("block hash does not match block data")
	}

	// Check 3: Verify Merkle root matches transactions
	computedMerkleRoot := block.computeMerkleRoot()
	if computedMerkleRoot != block.MerkleRoot {
		return errors.New("merkle root does not match transactions")
	}

	// Check 4: Verify Proof-of-Work
	// Note: We've already verified that block.Hash matches block data (Check 2)
	// Now we just need to verify the hash meets the difficulty target
	if !consensus.ValidateProofOfWork(block.Hash, difficulty) {
		return errors.New("block does not meet proof-of-work requirement")
	}

	// Check 5: Verify previous block (except genesis)
	if block.Index > 0 {
		if blockchain.Height() < block.Index {
			return errors.New("previous block not found")
		}

		prevBlock := blockchain.Blocks[block.Index-1]
		if prevBlock.Hash != block.PrevHash {
			return errors.New("previous hash mismatch")
		}

		// Check index is sequential
		if block.Index != prevBlock.Index+1 {
			return errors.New("block index is not sequential")
		}
	} else {
		// Genesis block: previous hash should be "0"
		if block.PrevHash != "0" {
			return errors.New("genesis block must have previous hash '0'")
		}
	}

	// Check 6: Verify all transactions
	// Create a temporary UTXO set to validate transactions
	// (we can't modify the real UTXO set until block is confirmed)
	tempUTXO := NewUTXOSet()

	// For each transaction, verify it and apply it to temp UTXO
	for i, tx := range block.Transactions {
		// Verify transaction
		if err := VerifyTransaction(&tx, tempUTXO); err != nil {
			return fmt.Errorf("transaction %d invalid: %w", i, err)
		}

		// Apply transaction to temp UTXO (for next transaction validation)
		tempUTXO.ApplyTransaction(&tx)
	}

	// All checks passed
	return nil
}

//
// VerifyTransaction validates a transaction against the current UTXO set.
//
// Order matters. Each check prevents a specific class of attack.
//
func VerifyTransaction(tx *Transaction, utxo *UTXOSet) error {

	// ------------------------------------------------------------
	// 1️⃣ Recompute transaction ID
	// ------------------------------------------------------------
	//
	// Why this check exists:
	// - Prevents tampering with inputs/outputs after signing
	// - Ensures tx.ID actually represents tx content
	//
	// What breaks if removed:
	// - Attacker can change outputs without changing tx.ID
	//
	computedID, err := ComputeTxID(tx)
	if err != nil {
		return err
	}

	if computedID != tx.ID {
		return errors.New("transaction ID mismatch")
	}

	// ------------------------------------------------------------
	// 2️⃣ Check for duplicate inputs
	// ------------------------------------------------------------
	//
	// Why this check exists:
	// - Prevents spending the SAME UTXO twice in ONE transaction
	//
	// Example attack:
	//   Inputs: [(txA,0), (txA,0)]
	//
	// Without this check:
	// - input sum would be double-counted
	//
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

	// ------------------------------------------------------------
	// 3️⃣ Verify inputs exist and sum input value
	// ------------------------------------------------------------
	//
	// Why this check exists:
	// - Ensures inputs refer to real, unspent outputs
	// - Prevents double spending
	//
	// What breaks if removed:
	// - Users could spend nonexistent coins
	//
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

	// ------------------------------------------------------------
	// 4️⃣ Sum output values
	// ------------------------------------------------------------
	//
	// Why this exists:
	// - Needed for value conservation check
	//
	var outputSum float64
	for _, out := range tx.Outputs {
		if out.Amount <= 0 {
			return errors.New("output amount must be positive")
		}
		outputSum += out.Amount
	}

	// ------------------------------------------------------------
	// 5️⃣ Value conservation check
	// ------------------------------------------------------------
	//
	// Rule:
	//   sum(outputs) <= sum(inputs)
	//
	// Why this exists:
	// - Prevents inflation
	//
	// What breaks if removed:
	// - Users can create money out of thin air
	//
	if outputSum > inputSum {
		return errors.New("output value exceeds input value")
	}

	// ------------------------------------------------------------
	// 6️⃣ Signature verification
	// ------------------------------------------------------------
	//
	// Verify that the transaction was signed by the owner of the public key.
	// We compute canonical bytes and pass them to the crypto package.
	//
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

	// ------------------------------------------------------------
	// ✅ All checks passed
	// ------------------------------------------------------------
	return nil
}
