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

type Server struct {
	blockchain *chain.Blockchain
	mempool    *chain.Mempool
	aiClient   *ai.Client
	difficulty int
	port       string
	walletStore *wallet.WalletStore
}

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

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next(w, r)
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/health", corsMiddleware(s.handleHealth))
	http.HandleFunc("/blocks", corsMiddleware(s.handleGetBlocks))
	http.HandleFunc("/chain", corsMiddleware(s.handleGetChain))
	http.HandleFunc("/mempool", corsMiddleware(s.handleGetMempool))
	http.HandleFunc("/transactions", corsMiddleware(s.handlePostTransaction))
	http.HandleFunc("/mine", corsMiddleware(s.handleMine))
	http.HandleFunc("/balance/", corsMiddleware(s.handleGetBalance))
	
	http.HandleFunc("/api/wallet/generate", corsMiddleware(s.handleGenerateWallet))
	http.HandleFunc("/api/wallet/list", corsMiddleware(s.handleListWallets))
	http.HandleFunc("/api/wallet/transfer", corsMiddleware(s.handleTransfer))

	addr := ":" + s.port
	log.Printf("Starting API server on %s (CORS enabled)", addr)
	return http.ListenAndServe(addr, nil)
}

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

func (s *Server) handleGetBlocks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	blocks := s.blockchain.Blocks

	response := map[string]interface{}{
		"blocks": blocks,
		"count":  len(blocks),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

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

func (s *Server) handlePostTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var tx chain.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if err := chain.VerifyTransaction(&tx, s.blockchain.UTXO); err != nil {
		http.Error(w, fmt.Sprintf("Invalid transaction: %v", err), http.StatusBadRequest)
		return
	}

	if s.aiClient != nil {
		score, err := s.aiClient.ScoreTransaction(&tx)
		if err != nil {
			log.Printf("AI scoring failed: %v (continuing anyway)", err)
		} else {
			log.Printf("Transaction %s scored: anomaly=%.2f, fee_adequacy=%.2f",
				tx.ID, score.AnomalyScore, score.FeeAdequacy)
			
			if score.AnomalyScore > 0.7 {
				http.Error(w, "Transaction flagged as anomalous by AI", http.StatusBadRequest)
				return
			}
		}
	}

	if err := s.mempool.AddTransaction(&tx); err != nil {
		http.Error(w, fmt.Sprintf("Failed to add transaction: %v", err), http.StatusConflict)
		return
	}

	response := map[string]interface{}{
		"status":  "accepted",
		"txid":    tx.ID,
		"message": "Transaction added to mempool",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleMine(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	txs := s.mempool.GetTransactions()
	if len(txs) == 0 {
		http.Error(w, "No transactions in mempool", http.StatusBadRequest)
		return
	}

	txSlice := make([]chain.Transaction, len(txs))
	for i, tx := range txs {
		txSlice[i] = *tx
	}

	tip := s.blockchain.Tip()

	block := chain.NewBlock(
		tip.Index+1,
		tip.Hash,
		txSlice,
	)

	log.Printf("Mining block %d with difficulty %d...", block.Index, s.difficulty)
	startTime := time.Now()
	
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
	
	block.Hash = hash
	block.Nonce = nonce

	duration := time.Since(startTime)
	log.Printf("Block %d mined in %v (hash: %s)", block.Index, duration, block.Hash)

	s.blockchain.AddBlock(block)

	for _, tx := range txs {
		s.mempool.RemoveTransaction(tx.ID)
	}

	response := map[string]interface{}{
		"block":   block,
		"message": "Block mined successfully",
		"time":    duration.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleGetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	address := r.URL.Path[len("/balance/"):]
	if address == "" {
		http.Error(w, "Address required", http.StatusBadRequest)
		return
	}

	balance := s.blockchain.UTXO.BalanceOf(address)

	response := map[string]interface{}{
		"address": address,
		"balance": balance,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

