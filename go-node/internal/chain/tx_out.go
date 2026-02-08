package chain

type TxOut struct {
	Address string  `json:"address"` // Hash of recipient's public key
	Amount  float64 `json:"amount"`  // Value in coins (using float64 for precision)
}