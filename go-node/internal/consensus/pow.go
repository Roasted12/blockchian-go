package consensus

import (
	"encoding/hex"
	"math/big"
)

const (
	DefaultDifficulty = 4 // Start with difficulty 4 for learning
)

func MineBlock(computeHashFunc func(int64) string, setNonceFunc func(int64), difficulty int) (string, int64) {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	nonce := int64(0)
	maxNonce := int64(^uint64(0) >> 1) // Max int64 value (safety limit)

	for nonce < maxNonce {
		setNonceFunc(nonce)

		hash := computeHashFunc(nonce)

		hashInt := new(big.Int)
		hashBytes, err := hex.DecodeString(hash)
		if err != nil {
			return "", 0
		}
		hashInt.SetBytes(hashBytes)

		if hashInt.Cmp(target) == -1 {
			return hash, nonce
		}

		nonce++
	}

	return "", 0
}

func ValidateProofOfWork(hash string, difficulty int) bool {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	hashInt := new(big.Int)
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return false
	}
	hashInt.SetBytes(hashBytes)

	return hashInt.Cmp(target) == -1
}

func GetDifficultyFromHash(hash string) int {
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

func AdjustDifficulty(currentDifficulty int, targetBlockTime, actualBlockTime int64) int {
	if actualBlockTime < targetBlockTime/2 {
		return currentDifficulty + 1
	}

	if actualBlockTime > targetBlockTime*2 {
		if currentDifficulty > 1 {
			return currentDifficulty - 1
		}
		return 1 // Minimum difficulty is 1
	}

	return currentDifficulty
}

func HashToBigInt(hash string) (*big.Int, error) {
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}
	hashInt := new(big.Int)
	hashInt.SetBytes(hashBytes)
	return hashInt, nil
}

