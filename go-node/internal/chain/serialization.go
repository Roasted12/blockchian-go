package chain

import (
	"bytes"
	"encoding/json"
	"sort"

	"ai-blockchain/go-node/internal/crypto"
)

type txForHash struct {
	Inputs  []TxIn  `json:"inputs"`
	Outputs []TxOut `json:"outputs"`
}

func CanonicalTxBytes(tx *Transaction) ([]byte, error) {
	inputsCopy := make([]TxIn, len(tx.Inputs))
	copy(inputsCopy, tx.Inputs)
	outputsCopy := make([]TxOut, len(tx.Outputs))
	copy(outputsCopy, tx.Outputs)

	sort.Slice(inputsCopy, func(i, j int) bool {
		if inputsCopy[i].TxID == inputsCopy[j].TxID {
			return inputsCopy[i].Index < inputsCopy[j].Index
		}
		return inputsCopy[i].TxID < inputsCopy[j].TxID
	})

	sort.Slice(outputsCopy, func(i, j int) bool {
		return outputsCopy[i].Address < outputsCopy[j].Address
	})

	tmp := txForHash{
		Inputs:  inputsCopy,
		Outputs: outputsCopy,
	}

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	err := enc.Encode(tmp)
	if err != nil {
		return nil, err
	}

	data := buf.Bytes()
	if len(data) > 0 && data[len(data)-1] == '\n' {
		data = data[:len(data)-1]
	}

	return data, nil
}

func ComputeTxID(tx *Transaction) (string, error) {
	canonical, err := CanonicalTxBytes(tx)
	if err != nil {
		return "", err
	}
	return crypto.SHA256(canonical), nil
}