package crypto

/*
MERKLE TREE

A Merkle tree allows us to:
- Commit to a list of transactions
- Detect tampering efficiently
- Prove inclusion with logarithmic data

This implementation:
- Uses SHA-256
- Operates on transaction IDs (already hashes)
*/

//
// MerkleRoot computes the Merkle root of a list of transaction IDs.
//
func MerkleRoot(txIDs []string) string {

	// Special case: no transactions
	if len(txIDs) == 0 {
		return SHA256([]byte{})
	}

	// Start with leaf hashes
	hashes := make([]string, len(txIDs))
	copy(hashes, txIDs)

	// Build tree upwards
	for len(hashes) > 1 {

		var nextLevel []string

		for i := 0; i < len(hashes); i += 2 {

			// If odd number of nodes, duplicate last
			if i+1 == len(hashes) {
				hashes = append(hashes, hashes[i])
			}

			combined := hashes[i] + hashes[i+1]
			parentHash := SHA256([]byte(combined))

			nextLevel = append(nextLevel, parentHash)
		}

		hashes = nextLevel
	}

	// Final root
	return hashes[0]
}
