package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"ai-blockchain/go-node/internal/chain"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	enabled    bool
}

type ScoreResponse struct {
	AnomalyScore float64 `json:"anomaly_score"`  // 0.0 = normal, 1.0 = highly anomalous
	FeeAdequacy  float64 `json:"fee_adequacy"`   // 0.0 = low fee, 1.0 = high fee
	Message      string  `json:"message,omitempty"`
}

func NewClient(baseURL string, timeout time.Duration, enabled bool) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		enabled: enabled,
	}
}

func (c *Client) ScoreTransaction(tx *chain.Transaction) (*ScoreResponse, error) {
	if !c.enabled {
		return &ScoreResponse{
			AnomalyScore: 0.0,
			FeeAdequacy:  0.5,
		}, nil
	}

	features := extractTxFeatures(tx)

	reqBody, err := json.Marshal(features)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal features: %w", err)
	}

	url := c.baseURL + "/score/tx"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &ScoreResponse{
			AnomalyScore: 0.0,
			FeeAdequacy:  0.5,
			Message:      "AI service unavailable",
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI service returned status %d: %s", resp.StatusCode, string(body))
	}

	var score ScoreResponse
	if err := json.NewDecoder(resp.Body).Decode(&score); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &score, nil
}

type TxFeatures struct {
	NumInputs    int     `json:"num_inputs"`
	NumOutputs   int     `json:"num_outputs"`
	TotalInput   float64 `json:"total_input"`
	TotalOutput  float64 `json:"total_output"`
	Fee          float64 `json:"fee"`
	FeeRate      float64 `json:"fee_rate"`      // Fee per byte (simplified)
	ChangeRatio   float64 `json:"change_ratio"` // Output / Input ratio
	InputDiversity int    `json:"input_diversity"` // Number of unique input addresses
}

func extractTxFeatures(tx *chain.Transaction) *TxFeatures {
	var totalInput float64
	inputAddresses := make(map[string]bool)
	for _, in := range tx.Inputs {
		inputAddresses[in.TxID] = true
	}

	var totalOutput float64
	for _, out := range tx.Outputs {
		totalOutput += out.Amount
	}

	fee := totalInput - totalOutput
	if fee < 0 {
		fee = 0
	}

	txSize := len(tx.ID) + len(tx.Signature) + len(tx.PubKey) // Rough estimate
	feeRate := 0.0
	if txSize > 0 {
		feeRate = fee / float64(txSize)
	}

	changeRatio := 0.0
	if totalInput > 0 {
		changeRatio = totalOutput / totalInput
	}

	return &TxFeatures{
		NumInputs:     len(tx.Inputs),
		NumOutputs:    len(tx.Outputs),
		TotalInput:    totalInput,
		TotalOutput:   totalOutput,
		Fee:           fee,
		FeeRate:       feeRate,
		ChangeRatio:   changeRatio,
		InputDiversity: len(inputAddresses),
	}
}

