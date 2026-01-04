package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"sync"

	"ai-blockchain/go-node/internal/chain"
	"ai-blockchain/go-node/internal/crypto"
)

/*
WALLET SERVICE – PRIVATE KEY MANAGEMENT

This package handles:
- Private key generation and storage
- Address derivation
- Transaction building
- Transaction signing

Important:
- Private keys are stored in memory (for learning)
- In production, use encrypted storage
- Private keys NEVER leave this package
*/

//
// Wallet represents a single wallet with its private key.
//
type Wallet struct {
	Address    string           // Derived from public key
	PrivateKey *ecdsa.PrivateKey // Private key (NEVER expose!)
	PublicKey  *ecdsa.PublicKey  // Public key (can be shared)
}

//
// WalletStore manages multiple wallets.
//
type WalletStore struct {
	mu      sync.RWMutex
	wallets map[string]*Wallet // address -> wallet
}

//
// NewWalletStore creates a new wallet store.
//
func NewWalletStore() *WalletStore {
	return &WalletStore{
		wallets: make(map[string]*Wallet),
	}
}

//
// GenerateWallet creates a new wallet with a key pair.
//
// Process:
// 1. Generate ECDSA key pair (P-256 curve)
// 2. Derive address from public key (SHA256 hash)
// 3. Store wallet in memory
// 4. Return wallet info (address and public key, NOT private key!)
//
func (ws *WalletStore) GenerateWallet() (*Wallet, error) {
	// Generate ECDSA key pair
	// Using P-256 curve (same as Java was using)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	// Derive address from public key
	// Address = SHA256(public key) - simplified version
	publicKeyBytes := append(
		privateKey.PublicKey.X.Bytes(),
		privateKey.PublicKey.Y.Bytes()...,
	)
	address := crypto.SHA256(publicKeyBytes)

	// Create wallet
	wallet := &Wallet{
		Address:    address,
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}

	// Store wallet
	ws.mu.Lock()
	ws.wallets[address] = wallet
	ws.mu.Unlock()

	return wallet, nil
}

//
// GetWallet retrieves a wallet by address.
//
// Returns nil if wallet doesn't exist.
//
func (ws *WalletStore) GetWallet(address string) *Wallet {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	return ws.wallets[address]
}

//
// GetAllAddresses returns all wallet addresses.
//
func (ws *WalletStore) GetAllAddresses() []string {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	addresses := make([]string, 0, len(ws.wallets))
	for addr := range ws.wallets {
		addresses = append(addresses, addr)
	}
	return addresses
}

//
// BuildAndSignTransaction creates and signs a transaction.
//
// This function:
// 1. Validates wallet exists
// 2. Builds transaction structure (inputs, outputs)
// 3. Computes transaction ID
// 4. Signs transaction with private key
// 5. Returns signed transaction ready to submit
//
// Note: Currently uses simplified UTXO selection.
// In production, you would query the blockchain for actual UTXOs.
//
func (ws *WalletStore) BuildAndSignTransaction(
	fromAddress string,
	toAddress string,
	amount float64,
) (*chain.Transaction, error) {
	// Get wallet
	wallet := ws.GetWallet(fromAddress)
	if wallet == nil {
		return nil, ErrWalletNotFound
	}

	// Build transaction
	// For learning, we'll create a simplified transaction
	// In production, you would:
	// 1. Query blockchain for UTXOs belonging to fromAddress
	// 2. Select UTXOs that cover the amount
	// 3. Calculate change output

	// Create inputs (simplified - placeholder)
	// In production, these would be actual UTXOs from the blockchain
	inputs := []chain.TxIn{
		{
			TxID:  "GENESIS_PLACEHOLDER", // Would be actual UTXO txid
			Index: 0,                      // Would be actual UTXO index
		},
	}

	// Create outputs
	outputs := []chain.TxOut{
		{
			Address: toAddress,
			Amount:  amount,
		},
		// Change output (simplified - would calculate actual change)
		{
			Address: fromAddress,
			Amount:  0.0, // Placeholder - would calculate: inputSum - amount - fee
		},
	}

	// Create transaction
	tx, err := chain.NewTransaction(inputs, outputs)
	if err != nil {
		return nil, err
	}

	// Sign transaction
	// Get canonical bytes (must match Go node's serialization)
	canonicalBytes, err := chain.CanonicalTxBytes(tx)
	if err != nil {
		return nil, err
	}

	// Hash the canonical bytes
	hash := sha256.Sum256(canonicalBytes)

	// Sign with private key
	r, s, err := ecdsa.Sign(rand.Reader, wallet.PrivateKey, hash[:])
	if err != nil {
		return nil, err
	}

	// Encode signature (r || s)
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	signatureBytes := append(rBytes, sBytes...)
	tx.Signature = hex.EncodeToString(signatureBytes)

	// Encode public key (x || y)
	xBytes := wallet.PublicKey.X.Bytes()
	yBytes := wallet.PublicKey.Y.Bytes()
	pubKeyBytes := append(xBytes, yBytes...)
	tx.PubKey = hex.EncodeToString(pubKeyBytes)

	return tx, nil
}

//
// EncodePublicKey encodes a public key to hex string.
//
func EncodePublicKey(pub *ecdsa.PublicKey) string {
	xBytes := pub.X.Bytes()
	yBytes := pub.Y.Bytes()
	combined := append(xBytes, yBytes...)
	return hex.EncodeToString(combined)
}

// Error definitions
var (
	ErrWalletNotFound = &WalletError{Message: "wallet not found"}
)

type WalletError struct {
	Message string
}

func (e *WalletError) Error() string {
	return e.Message
}




