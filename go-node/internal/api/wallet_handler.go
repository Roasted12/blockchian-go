package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"ai-blockchain/go-node/internal/chain"
	"ai-blockchain/go-node/internal/wallet"
)

func (s *Server) handleGenerateWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	newWallet, err := s.walletStore.GenerateWallet()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate wallet: %v", err), http.StatusInternalServerError)
		return
	}

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

func (s *Server) handleTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if request.From == "" || request.To == "" || request.Amount <= 0 {
		http.Error(w, "Invalid request: from, to, and amount (positive) are required", http.StatusBadRequest)
		return
	}

	tx, err := s.walletStore.BuildAndSignTransaction(
		request.From,
		request.To,
		request.Amount,
		s.blockchain.UTXO,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to build transaction: %v", err), http.StatusBadRequest)
		return
	}

	if err := chain.VerifyTransaction(tx, s.blockchain.UTXO); err != nil {
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

	if s.aiClient != nil {
		score, err := s.aiClient.ScoreTransaction(tx)
		if err != nil {
			log.Printf("AI scoring failed: %v (continuing anyway)", err)
		} else {
			log.Printf("Transaction %s scored: anomaly=%.2f, fee_adequacy=%.2f",
				tx.ID, score.AnomalyScore, score.FeeAdequacy)

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

	if err := s.mempool.AddTransaction(tx); err != nil {
		http.Error(w, fmt.Sprintf("Failed to add to mempool: %v", err), http.StatusConflict)
		return
	}

	response := map[string]interface{}{
		"status":  "submitted",
		"txid":    tx.ID,
		"message": "Transaction signed and submitted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

