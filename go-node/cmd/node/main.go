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

func main() {
	port := flag.String("port", "8080", "API server port")
	difficulty := flag.Int("difficulty", consensus.DefaultDifficulty, "Mining difficulty")
	aiURL := flag.String("ai-url", "", "AI service URL (empty = disabled)")
	aiTimeout := flag.Int("ai-timeout", 5, "AI service timeout in seconds")
	flag.Parse()

	log.Println("Starting blockchain node...")
	log.Printf("Port: %s, Difficulty: %d", *port, *difficulty)

	walletStore := wallet.NewWalletStore()
	log.Println("Wallet store initialized")

	defaultWallet, err := walletStore.GenerateWallet()
	if err != nil {
		log.Fatalf("Failed to create default wallet for genesis: %v", err)
	}
	log.Printf("Default wallet created for genesis: %s", defaultWallet.Address)

	genesisOutput := chain.TxOut{
		Address: defaultWallet.Address,
		Amount:  1000.0,
	}
	
	genesisTx, err := chain.NewTransaction(
		[]chain.TxIn{}, // No inputs (genesis creates coins)
		[]chain.TxOut{genesisOutput},
	)
	if err != nil {
		log.Fatalf("Failed to create genesis transaction: %v", err)
	}
	
	genesisTx.Signature = "genesis"
	genesisTx.PubKey = "genesis"

	genesisBlock := chain.NewBlock(
		0,
		"0",
		[]chain.Transaction{*genesisTx},
	)

	blockchain := chain.NewBlockchain(genesisBlock)
	log.Printf("Genesis block created: %s", genesisBlock.Hash)

	genesisBalance := blockchain.UTXO.BalanceOf(defaultWallet.Address)
	log.Printf("Default wallet (genesis recipient) balance: %.2f coins", genesisBalance)
	if genesisBalance == 0 {
		log.Printf("WARNING: Genesis coins not found in UTXO set!")
	}

	mempool := chain.NewMempool()
	log.Println("Mempool initialized")

	var aiClient *ai.Client
	if *aiURL != "" {
		timeout := time.Duration(*aiTimeout) * time.Second
		aiClient = ai.NewClient(*aiURL, timeout, true)
		log.Printf("AI scoring enabled: %s (timeout: %v)", *aiURL, timeout)
	} else {
		aiClient = ai.NewClient("", 0, false)
		log.Println("AI scoring disabled")
	}

	server := api.NewServer(blockchain, mempool, aiClient, *difficulty, *port, walletStore)

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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("\nShutting down gracefully...")
	log.Println("Node stopped")
}

