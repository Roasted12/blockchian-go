package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ai-blockchain/go-node/internal/ai"
	"ai-blockchain/go-node/internal/api"
	"ai-blockchain/go-node/internal/chain"
	"ai-blockchain/go-node/internal/consensus"
	"ai-blockchain/go-node/internal/wallet"
)

/*
MAIN ENTRY POINT – BLOCKCHAIN NODE

This is the entry point for the blockchain node.

What it does:
1. Creates genesis block
2. Initializes blockchain
3. Initializes mempool
4. Optionally connects to AI service
5. Starts API server
6. Handles graceful shutdown

Command-line flags:
- -port: API server port (default: 8080)
- -difficulty: Mining difficulty (default: 4)
- -ai-url: AI service URL (default: "", disabled)
- -ai-timeout: AI service timeout in seconds (default: 5)
*/

func main() {
	// Parse command-line flags
	port := flag.String("port", "8080", "API server port")
	difficulty := flag.Int("difficulty", consensus.DefaultDifficulty, "Mining difficulty")
	aiURL := flag.String("ai-url", "", "AI service URL (empty = disabled)")
	aiTimeout := flag.Int("ai-timeout", 5, "AI service timeout in seconds")
	flag.Parse()

	log.Println("Starting blockchain node...")
	log.Printf("Port: %s, Difficulty: %d", *port, *difficulty)

	// Create genesis block
	// Genesis block is special: it has no previous block
	// It typically contains initial coin distribution
	genesisTx, err := createGenesisTransaction()
	if err != nil {
		log.Fatalf("Failed to create genesis transaction: %v", err)
	}

	genesisBlock := chain.NewBlock(
		0,                    // Index 0 (first block)
		"0",                  // Previous hash (none, so use "0")
		[]chain.Transaction{genesisTx}, // Genesis transaction
	)

	// Initialize blockchain with genesis block
	blockchain := chain.NewBlockchain(genesisBlock)
	log.Printf("Genesis block created: %s", genesisBlock.Hash)

	// Initialize mempool
	mempool := chain.NewMempool()
	log.Println("Mempool initialized")

	// Initialize wallet store
	walletStore := wallet.NewWalletStore()
	log.Println("Wallet store initialized")

	// Initialize AI client (optional)
	var aiClient *ai.Client
	if *aiURL != "" {
		timeout := time.Duration(*aiTimeout) * time.Second
		aiClient = ai.NewClient(*aiURL, timeout, true)
		log.Printf("AI scoring enabled: %s (timeout: %v)", *aiURL, timeout)
	} else {
		aiClient = ai.NewClient("", 0, false)
		log.Println("AI scoring disabled")
	}

	// Create and start API server
	server := api.NewServer(blockchain, mempool, aiClient, walletStore, *difficulty, *port)

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("Blockchain node is running!")
	log.Println("API endpoints:")
	log.Println("  GET  /health          - Health check")
	log.Println("  GET  /blocks          - Get all blocks")
	log.Println("  GET  /chain           - Get blockchain info")
	log.Println("  GET  /mempool         - Get pending transactions")
	log.Println("  GET  /balance/:addr  - Get balance for address")
	log.Println("  POST /transactions    - Submit new transaction")
	log.Println("  POST /mine            - Mine a new block")
	log.Println("")
	log.Println("Wallet endpoints:")
	log.Println("  GET  /api/wallet/generate - Generate new wallet")
	log.Println("  GET  /api/wallet/list    - List all wallets")
	log.Println("  POST /api/wallet/transfer - Create and submit transaction")

	// Wait for interrupt signal (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("\nShutting down gracefully...")
	log.Println("Node stopped")
}

//
// createGenesisTransaction creates the initial transaction for the genesis block.
//
// Genesis transaction:
// - Has no inputs (creates coins from nothing)
// - Has one output (initial coin distribution)
// - Is not signed (genesis is special)
//
// In a real blockchain:
// - Genesis block might have multiple outputs
// - Outputs might go to founders, developers, etc.
// - This is the ONLY way new coins are created
//
func createGenesisTransaction() (chain.Transaction, error) {
	// Create a dummy address for genesis
	// In production, this would be a real address
	genesisAddress := "0000000000000000000000000000000000000000"

	// Create genesis output (initial coin distribution)
	// Amount: 1000 coins (arbitrary for learning)
	genesisOutput := chain.TxOut{
		Address: genesisAddress,
		Amount:  1000.0,
	}

	// Create transaction with no inputs (genesis creates coins)
	tx, err := chain.NewTransaction(
		[]chain.TxIn{}, // No inputs
		[]chain.TxOut{genesisOutput},
	)
	if err != nil {
		return chain.Transaction{}, err
	}

	// Genesis transaction doesn't need a signature
	// (it's the first transaction, so there's no previous owner)
	tx.Signature = "genesis"
	tx.PubKey = "genesis"

	return *tx, nil
}

