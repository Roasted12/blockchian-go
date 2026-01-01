package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"ai-blockchain/go-node/internal/chain"
	"ai-blockchain/go-node/internal/wallet"
)

/*
WALLET API HANDLERS

These endpoints handle wallet operations:
- Generate wallets
- List wallets
- Create and sign transactions
- Check balances

All private key operations happen here, never exposed to clients.
*/

//
// handleGenerateWallet creates a new wallet.
//
// GET /api/wallet/generate
//
// Response:
// {
//   "address": "...",
//   "public_key": "...",
//   "message": "Wallet generated successfully"
// }
//
func (s *Server) handleGenerateWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Generate new wallet
	newWallet, err := s.walletStore.GenerateWallet()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate wallet: %v", err), http.StatusInternalServerError)
		return
	}

	// Encode public key
	publicKeyHex := wallet.EncodePublicKey(newWallet.PublicKey)

	response := map[string]interface{}{
		"address":    newWallet.Address,
		"public_key": publicKeyHex,
		"message":    "Wallet generated and stored successfully",
		"note":       "Private key is stored securely in wallet service",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//
// handleListWallets returns all wallet addresses.
//
// GET /api/wallet/list
//
// Response:
// {
//   "addresses": ["addr1", "addr2", ...],
//   "count": 2
// }
//
func (s *Server) handleListWallets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	addresses := s.walletStore.GetAllAddresses()

	response := map[string]interface{}{
		"addresses": addresses,
		"count":    len(addresses),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//
// handleTransfer creates and submits a transaction.
//
// POST /api/wallet/transfer
//
// Request:
// {
//   "from": "address1",
//   "to": "address2",
//   "amount": 10.5
// }
//
// Response:
// {
//   "status": "submitted",
//   "txid": "...",
//   "message": "Transaction signed and submitted successfully"
// }
//
func (s *Server) handleTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode request
	var request struct {
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Validate request
	if request.From == "" || request.To == "" || request.Amount <= 0 {
		http.Error(w, "Invalid request: from, to, and amount (positive) are required", http.StatusBadRequest)
		return
	}

	// Build and sign transaction
	tx, err := s.walletStore.BuildAndSignTransaction(
		request.From,
		request.To,
		request.Amount,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to build transaction: %v", err), http.StatusBadRequest)
		return
	}

	// Validate transaction (before submitting)
	if err := chain.VerifyTransaction(tx, s.blockchain.UTXO); err != nil {
		// Transaction might fail validation if UTXOs don't exist
		// This is expected for learning - user needs to have coins first
		response := map[string]interface{}{
			"error": fmt.Sprintf("Transaction validation failed: %v", err),
			"hint":  "Make sure you have coins. Try using genesis address or mine a block first.",
			"txid":  tx.ID,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Optional: Score transaction with AI
	if s.aiClient != nil {
		score, err := s.aiClient.ScoreTransaction(tx)
		if err != nil {
			// Log but don't fail
			log.Printf("AI scoring failed: %v (continuing anyway)", err)
		} else {
			log.Printf("Transaction %s scored: anomaly=%.2f, fee_adequacy=%.2f",
				tx.ID, score.AnomalyScore, score.FeeAdequacy)

			// Reject if anomaly score too high
			if score.AnomalyScore > 0.7 {
				response := map[string]interface{}{
					"error": "Transaction flagged as anomalous by AI",
					"score": score.AnomalyScore,
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
				return
			}
		}
	}

	// Add to mempool
	if err := s.mempool.AddTransaction(tx); err != nil {
		http.Error(w, fmt.Sprintf("Failed to add to mempool: %v", err), http.StatusConflict)
		return
	}

	// Return success
	response := map[string]interface{}{
		"status":  "submitted",
		"txid":    tx.ID,
		"message": "Transaction signed and submitted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

