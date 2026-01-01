package consensus

import (
	"encoding/hex"
	"math/big"
)

/*
PROOF-OF-WORK CONSENSUS

Proof-of-Work (PoW) is a consensus mechanism where:
- Miners compete to find a nonce that makes block hash < target
- Difficulty adjusts to maintain consistent block times
- First miner to find valid nonce gets to create the block

How it works:
1. Block is created with transactions
2. Nonce starts at 0
3. Hash is computed: SHA256(block data + nonce)
4. If hash < target: block is valid, mining succeeds
5. If hash >= target: increment nonce, repeat

Difficulty:
- Target = 2^(256 - difficulty)
- Lower target = harder to find valid hash
- Difficulty adjusts based on block time
*/

//
// Difficulty represents the mining difficulty.
//
// Higher difficulty = harder to mine = more secure
// Lower difficulty = easier to mine = faster blocks
//
// Typical values:
// - Difficulty 1-4: Testing/development
// - Difficulty 5-10: Small networks
// - Difficulty 20+: Production networks (Bitcoin uses ~20-30)
//
const (
	DefaultDifficulty = 4 // Start with difficulty 4 for learning
)

//
// MineBlock attempts to mine a block using Proof-of-Work.
//
// What this function does:
// 1. Sets initial nonce to 0
// 2. Computes block hash with current nonce (using computeHashFunc)
// 3. Checks if hash meets difficulty target
// 4. If not, increments nonce and repeats
// 5. Returns when valid hash is found
//
// Parameters:
// - computeHashFunc: Function that computes hash given a nonce
// - setNonceFunc: Function that sets the nonce on the block
// - difficulty: Number of leading zeros required in hash
//
// Returns:
// - (hash, nonce) if mining succeeded
// - ("", 0) if mining failed (should not happen in normal operation)
//
// This design avoids importing the chain package, breaking the import cycle.
//
func MineBlock(computeHashFunc func(int64) string, setNonceFunc func(int64), difficulty int) (string, int64) {
	// Compute target: hash must be less than this value
	// Target = 2^(256 - difficulty)
	// Example: difficulty 4 means hash must start with 4 zeros
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	// Try nonces starting from 0
	// In practice, miners use random starting points to avoid collisions
	nonce := int64(0)
	maxNonce := int64(^uint64(0) >> 1) // Max int64 value (safety limit)

	for nonce < maxNonce {
		// Set nonce on block
		setNonceFunc(nonce)

		// Compute hash with current nonce
		hash := computeHashFunc(nonce)

		// Convert hash to big integer for comparison
		hashInt := new(big.Int)
		hashBytes, err := hex.DecodeString(hash)
		if err != nil {
			return "", 0
		}
		hashInt.SetBytes(hashBytes)

		// Check if hash meets target (hash < target)
		if hashInt.Cmp(target) == -1 {
			// Valid hash found! Mining succeeded
			return hash, nonce
		}

		// Hash doesn't meet target, try next nonce
		nonce++
	}

	// Exceeded max nonce (should not happen in practice)
	return "", 0
}

//
// ValidateProofOfWork checks if a hash meets the difficulty target.
//
// This is called when receiving blocks from other nodes.
// We verify that the block was actually mined (not just created).
//
// Parameters:
// - hash: The block hash to validate (hex-encoded)
// - difficulty: The difficulty level
//
// Returns:
// - true if hash is valid (meets difficulty)
// - false if hash is invalid (doesn't meet difficulty)
//
// Note: This function only checks if the hash meets the difficulty.
// The caller (chain package) should verify that the hash matches the block data.
//
func ValidateProofOfWork(hash string, difficulty int) bool {
	// Compute target (same as in mining)
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	// Convert hash to big integer
	hashInt := new(big.Int)
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return false
	}
	hashInt.SetBytes(hashBytes)

	// Check if hash meets target (hash < target)
	return hashInt.Cmp(target) == -1
}

//
// GetDifficultyFromHash returns the number of leading zeros in a hash.
//
// This is useful for:
// - Displaying mining progress
// - Adjusting difficulty dynamically
// - Debugging mining issues
//
func GetDifficultyFromHash(hash string) int {
	// Count leading zeros in hex representation
	// Each hex digit represents 4 bits
	// Leading zero = 4 bits of zeros
	leadingZeros := 0
	for _, char := range hash {
		if char == '0' {
			leadingZeros++
		} else {
			break
		}
	}
	return leadingZeros
}

//
// AdjustDifficulty adjusts difficulty based on block time.
//
// Goal: Maintain consistent block time (e.g., 10 minutes)
//
// Algorithm:
// - If blocks are too fast: increase difficulty
// - If blocks are too slow: decrease difficulty
//
// Parameters:
// - currentDifficulty: Current mining difficulty
// - targetBlockTime: Desired time between blocks (seconds)
// - actualBlockTime: Actual time since last block (seconds)
//
// Returns:
// - New difficulty value
//
func AdjustDifficulty(currentDifficulty int, targetBlockTime, actualBlockTime int64) int {
	// If blocks are coming too fast, increase difficulty
	if actualBlockTime < targetBlockTime/2 {
		return currentDifficulty + 1
	}

	// If blocks are coming too slow, decrease difficulty
	if actualBlockTime > targetBlockTime*2 {
		if currentDifficulty > 1 {
			return currentDifficulty - 1
		}
		return 1 // Minimum difficulty is 1
	}

	// Block time is acceptable, keep current difficulty
	return currentDifficulty
}

//
// HashToBigInt converts a hex-encoded hash to a big integer.
//
// Helper function for difficulty calculations.
//
func HashToBigInt(hash string) (*big.Int, error) {
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}
	hashInt := new(big.Int)
	hashInt.SetBytes(hashBytes)
	return hashInt, nil
}

