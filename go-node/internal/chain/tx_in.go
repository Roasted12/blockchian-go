package chain

type TxIn struct {
	TxID string `json:"tx_id"`
	Index int `json:"index"`
}