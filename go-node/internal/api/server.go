package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"ai-blockchain/go-node/internal/ai"
	"ai-blockchain/go-node/internal/chain"
	"ai-blockchain/go-node/internal/consensus"
	"ai-blockchain/go-node/internal/wallet"
)

/*
API SERVER – REST ENDPOINTS

This file implements the REST API for the blockchain node.

Endpoints:
- GET  /health          - Health check
- GET  /blocks          - Get all blocks
- GET  /blocks/:hash    - Get specific block
- GET  /chain           - Get blockchain info
- GET  /mempool         - Get pending transactions
- GET  /balance/:addr   - Get balance for address
- POST /transactions    - Submit new transaction
- POST /mine            - Mine a new block

The API is used by:
- Java wallet (submits transactions, queries balances)
- Block explorers (queries blocks, transactions)
- Other nodes (block propagation - future)
*/

//
// Server represents the API server.
//
type Server struct {
	blockchain *chain.Blockchain
	mempool    *chain.Mempool
	aiClient   *ai.Client
	difficulty int
	port       string
	walletStore *wallet.WalletStore
}

//
// NewServer creates a new API server.
//
// Parameters:
// - blockchain: The blockchain instance
// - mempool: The mempool instance
// - aiClient: AI scoring client (can be nil if AI is disabled)
// - difficulty: Mining difficulty
// - port: Server port (e.g., "8080")
//
func NewServer(
	blockchain *chain.Blockchain,
	mempool *chain.Mempool,
	aiClient *ai.Client,
	difficulty int,
	port string,
	walletStore *wallet.WalletStore,
) *Server {
	return &Server{
		blockchain: blockchain,
		mempool:    mempool,
		aiClient:   aiClient,
		difficulty: difficulty,
		port:       port,
		walletStore: walletStore,
	}
}

//
// corsMiddleware adds CORS headers to allow web UI access.
//
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// Call the next handler
		next(w, r)
	}
}

//
// Start starts the HTTP server.
//
// This function:
// 1. Registers all route handlers with CORS support
// 2. Starts listening on the specified port
// 3. Blocks until server stops
//
func (s *Server) Start() error {
	// Register routes with CORS middleware
	http.HandleFunc("/health", corsMiddleware(s.handleHealth))
	http.HandleFunc("/blocks", corsMiddleware(s.handleGetBlocks))
	http.HandleFunc("/chain", corsMiddleware(s.handleGetChain))
	http.HandleFunc("/mempool", corsMiddleware(s.handleGetMempool))
	http.HandleFunc("/transactions", corsMiddleware(s.handlePostTransaction))
	http.HandleFunc("/mine", corsMiddleware(s.handleMine))
	http.HandleFunc("/balance/", corsMiddleware(s.handleGetBalance))
	
	// Wallet routes
	http.HandleFunc("/api/wallet/generate", corsMiddleware(s.handleGenerateWallet))
	http.HandleFunc("/api/wallet/list", corsMiddleware(s.handleListWallets))
	http.HandleFunc("/api/wallet/transfer", corsMiddleware(s.handleTransfer))

	// Start server
	addr := ":" + s.port
	log.Printf("Starting API server on %s (CORS enabled)", addr)
	return http.ListenAndServe(addr, nil)
}

//
// handleHealth returns server health status.
//
// Used for:
// - Load balancer health checks
// - Monitoring systems
// - Debugging
//
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"height":    s.blockchain.Height(),
		"mempool":   s.mempool.Size(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//
// handleGetBlocks returns all blocks in the blockchain.
//
// Response format:
// {
//   "blocks": [...],
//   "count": 10
// }
//
func (s *Server) handleGetBlocks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get all blocks
	blocks := s.blockchain.Blocks

	response := map[string]interface{}{
		"blocks": blocks,
		"count":  len(blocks),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//
// handleGetChain returns blockchain information.
//
// Response format:
// {
//   "height": 10,
//   "tip": {...},
//   "difficulty": 4
// }
//
func (s *Server) handleGetChain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tip := s.blockchain.Tip()

	response := map[string]interface{}{
		"height":    s.blockchain.Height(),
		"tip":       tip,
		"difficulty": s.difficulty,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//
// handleGetMempool returns all pending transactions.
//
// Response format:
// {
//   "transactions": [...],
//   "count": 5
// }
//
func (s *Server) handleGetMempool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	txs := s.mempool.GetTransactions()

	response := map[string]interface{}{
		"transactions": txs,
		"count":        len(txs),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//
// handlePostTransaction accepts a new transaction.
//
// Request body:
// {
//   "id": "...",
//   "inputs": [...],
//   "outputs": [...],
//   "signature": "...",
//   "pubkey": "...",
//   "timestamp": 1234567890
// }
//
// This endpoint:
// 1. Validates the transaction
// 2. Optionally scores it with AI
// 3. Adds it to mempool if valid
//
func (s *Server) handlePostTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode transaction from request body
	var tx chain.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Validate transaction
	if err := chain.VerifyTransaction(&tx, s.blockchain.UTXO); err != nil {
		http.Error(w, fmt.Sprintf("Invalid transaction: %v", err), http.StatusBadRequest)
		return
	}

	// Optional: Score transaction with AI
	if s.aiClient != nil {
		score, err := s.aiClient.ScoreTransaction(&tx)
		if err != nil {
			log.Printf("AI scoring failed: %v (continuing anyway)", err)
		} else {
			log.Printf("Transaction %s scored: anomaly=%.2f, fee_adequacy=%.2f",
				tx.ID, score.AnomalyScore, score.FeeAdequacy)
			
			// If anomaly score is too high, reject transaction
			// Threshold: 0.7 (higher = more anomalous)
			if score.AnomalyScore > 0.7 {
				http.Error(w, "Transaction flagged as anomalous by AI", http.StatusBadRequest)
				return
			}
		}
	}

	// Add to mempool
	if err := s.mempool.AddTransaction(&tx); err != nil {
		http.Error(w, fmt.Sprintf("Failed to add transaction: %v", err), http.StatusConflict)
		return
	}

	// Return success
	response := map[string]interface{}{
		"status":  "accepted",
		"txid":    tx.ID,
		"message": "Transaction added to mempool",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

//
// handleMine mines a new block from mempool transactions.
//
// This endpoint:
// 1. Gets transactions from mempool
// 2. Creates a new block
// 3. Mines the block (Proof-of-Work)
// 4. Adds block to blockchain
// 5. Removes transactions from mempool
//
// Response format:
// {
//   "block": {...},
//   "message": "Block mined successfully"
// }
//
func (s *Server) handleMine(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get transactions from mempool
	txs := s.mempool.GetTransactions()
	if len(txs) == 0 {
		http.Error(w, "No transactions in mempool", http.StatusBadRequest)
		return
	}

	// Convert []*Transaction to []Transaction
	txSlice := make([]chain.Transaction, len(txs))
	for i, tx := range txs {
		txSlice[i] = *tx
	}

	// Get current tip
	tip := s.blockchain.Tip()

	// Create new block
	block := chain.NewBlock(
		tip.Index+1,
		tip.Hash,
		txSlice,
	)

	// Mine the block (Proof-of-Work)
	log.Printf("Mining block %d with difficulty %d...", block.Index, s.difficulty)
	startTime := time.Now()
	
	// Create functions for mining (avoids import cycle)
	computeHashFunc := func(nonce int64) string {
		block.Nonce = nonce
		return block.ComputeHash()
	}
	setNonceFunc := func(nonce int64) {
		block.Nonce = nonce
	}
	
	hash, nonce := consensus.MineBlock(computeHashFunc, setNonceFunc, s.difficulty)
	if hash == "" {
		http.Error(w, "Failed to mine block", http.StatusInternalServerError)
		return
	}
	
	// Set the final hash and nonce
	block.Hash = hash
	block.Nonce = nonce

	duration := time.Since(startTime)
	log.Printf("Block %d mined in %v (hash: %s)", block.Index, duration, block.Hash)

	// Add block to blockchain
	s.blockchain.AddBlock(block)

	// Remove transactions from mempool
	for _, tx := range txs {
		s.mempool.RemoveTransaction(tx.ID)
	}

	// Return success
	response := map[string]interface{}{
		"block":   block,
		"message": "Block mined successfully",
		"time":    duration.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//
// handleGetBalance returns the balance for an address.
//
// URL format: /balance/:address
//
// Response format:
// {
//   "address": "...",
//   "balance": 100.5
// }
//
func (s *Server) handleGetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract address from URL path
	// Path format: /balance/address_here
	address := r.URL.Path[len("/balance/"):]
	if address == "" {
		http.Error(w, "Address required", http.StatusBadRequest)
		return
	}

	// Get balance from UTXO set
	balance := s.blockchain.UTXO.BalanceOf(address)

	response := map[string]interface{}{
		"address": address,
		"balance": balance,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

