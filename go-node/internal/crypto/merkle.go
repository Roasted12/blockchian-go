package crypto

func MerkleRoot(txIDs []string) string {

	if len(txIDs) == 0 {
		return SHA256([]byte{})
	}

	hashes := make([]string, len(txIDs))
	copy(hashes, txIDs)

	for len(hashes) > 1 {

		var nextLevel []string

		for i := 0; i < len(hashes); i += 2 {

			if i+1 == len(hashes) {
				hashes = append(hashes, hashes[i])
			}

			combined := hashes[i] + hashes[i+1]
			parentHash := SHA256([]byte(combined))

			nextLevel = append(nextLevel, parentHash)
		}

		hashes = nextLevel
	}

	return hashes[0]
}
