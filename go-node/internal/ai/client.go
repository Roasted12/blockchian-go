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

/*
AI CLIENT – ADVISORY SCORING SERVICE

This package implements the client for the AI scoring service.

The AI service provides:
- Transaction anomaly detection (IsolationForest)
- Fee adequacy estimation
- Peer reliability scoring (future)

Important:
- AI scoring is ADVISORY ONLY
- It does NOT affect consensus
- It helps prioritize transactions and detect suspicious activity
- If AI service is down, node continues operating normally
*/

//
// Client represents the AI scoring service client.
//
type Client struct {
	baseURL    string
	httpClient *http.Client
	enabled    bool
}

//
// ScoreResponse represents the AI scoring response.
//
type ScoreResponse struct {
	AnomalyScore float64 `json:"anomaly_score"`  // 0.0 = normal, 1.0 = highly anomalous
	FeeAdequacy  float64 `json:"fee_adequacy"`   // 0.0 = low fee, 1.0 = high fee
	Message      string  `json:"message,omitempty"`
}

//
// NewClient creates a new AI scoring client.
//
// Parameters:
// - baseURL: Base URL of the AI service (e.g., "http://localhost:5000")
// - timeout: HTTP request timeout (e.g., 5 seconds)
// - enabled: Whether AI scoring is enabled
//
// If enabled=false, all scoring calls will return default scores.
//
func NewClient(baseURL string, timeout time.Duration, enabled bool) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		enabled: enabled,
	}
}

//
// ScoreTransaction scores a transaction using the AI service.
//
// What this does:
// 1. Extracts features from the transaction
// 2. Sends features to AI service
// 3. Returns anomaly score and fee adequacy
//
// If AI service is unavailable:
// - Returns default scores (anomaly=0.0, fee=0.5)
// - Logs error but doesn't fail
//
// This is called BEFORE adding transaction to mempool.
//
func (c *Client) ScoreTransaction(tx *chain.Transaction) (*ScoreResponse, error) {
	// If AI is disabled, return default scores
	if !c.enabled {
		return &ScoreResponse{
			AnomalyScore: 0.0,
			FeeAdequacy:  0.5,
		}, nil
	}

	// Extract features from transaction
	features := extractTxFeatures(tx)

	// Create request
	reqBody, err := json.Marshal(features)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal features: %w", err)
	}

	// Send request to AI service
	url := c.baseURL + "/score/tx"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute request with timeout
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// AI service is down - return default scores
		return &ScoreResponse{
			AnomalyScore: 0.0,
			FeeAdequacy:  0.5,
			Message:      "AI service unavailable",
		}, nil
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Decode response
	var score ScoreResponse
	if err := json.NewDecoder(resp.Body).Decode(&score); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &score, nil
}

//
// TxFeatures represents the features extracted from a transaction.
//
// These features are sent to the AI service for scoring.
//
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

//
// extractTxFeatures extracts features from a transaction.
//
// These features are used by the AI model to:
// - Detect anomalous patterns
// - Estimate fee adequacy
// - Classify transaction types
//
func extractTxFeatures(tx *chain.Transaction) *TxFeatures {
	// Calculate input sum
	var totalInput float64
	inputAddresses := make(map[string]bool)
	for _, in := range tx.Inputs {
		// Note: We can't get the actual input value here without UTXO lookup
		// For now, we'll use a placeholder
		// In production, you'd look up the UTXO to get the amount
		inputAddresses[in.TxID] = true
	}

	// Calculate output sum
	var totalOutput float64
	for _, out := range tx.Outputs {
		totalOutput += out.Amount
	}

	// Calculate fee (simplified: input - output)
	// In reality, we'd need to look up input values from UTXO set
	fee := totalInput - totalOutput
	if fee < 0 {
		fee = 0 // Can't have negative fee
	}

	// Calculate fee rate (simplified)
	// In production, you'd calculate actual transaction size in bytes
	txSize := len(tx.ID) + len(tx.Signature) + len(tx.PubKey) // Rough estimate
	feeRate := 0.0
	if txSize > 0 {
		feeRate = fee / float64(txSize)
	}

	// Calculate change ratio
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

