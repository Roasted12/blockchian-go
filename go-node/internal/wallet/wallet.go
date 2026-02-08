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

type Wallet struct {
	Address    string           // Derived from public key
	PrivateKey *ecdsa.PrivateKey // Private key (NEVER expose!)
	PublicKey  *ecdsa.PublicKey  // Public key (can be shared)
}

type WalletStore struct {
	mu      sync.RWMutex
	wallets map[string]*Wallet // address -> wallet
}

func NewWalletStore() *WalletStore {
	return &WalletStore{
		wallets: make(map[string]*Wallet),
	}
}

func (ws *WalletStore) GenerateWallet() (*Wallet, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	publicKeyBytes := append(
		privateKey.PublicKey.X.Bytes(),
		privateKey.PublicKey.Y.Bytes()...,
	)
	address := crypto.SHA256(publicKeyBytes)

	wallet := &Wallet{
		Address:    address,
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}

	ws.mu.Lock()
	ws.wallets[address] = wallet
	ws.mu.Unlock()

	return wallet, nil
}

func (ws *WalletStore) GetWallet(address string) *Wallet {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	return ws.wallets[address]
}

func (ws *WalletStore) GetAllAddresses() []string {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	addresses := make([]string, 0, len(ws.wallets))
	for addr := range ws.wallets {
		addresses = append(addresses, addr)
	}
	return addresses
}

func (ws *WalletStore) BuildAndSignTransaction(
	fromAddress string,
	toAddress string,
	amount float64,
	utxo *chain.UTXOSet,
) (*chain.Transaction, error) {
	wallet := ws.GetWallet(fromAddress)
	if wallet == nil {
		return nil, ErrWalletNotFound
	}

	total, selected := utxo.FindSpendableOutputs(fromAddress, amount)
	if total < amount {
		return nil, ErrInsufficientFunds
	}

	inputs := make([]chain.TxIn, 0, len(selected))
	for _, key := range selected {
		inputs = append(inputs, chain.TxIn{
			TxID:  key.TxID,
			Index: key.Index,
		})
	}

	outputs := []chain.TxOut{
		{
			Address: toAddress,
			Amount:  amount,
		},
	}

	change := total - amount
	if change > 0 {
		outputs = append(outputs, chain.TxOut{
			Address: fromAddress,
			Amount:  change,
		})
	}

	tx, err := chain.NewTransaction(inputs, outputs)
	if err != nil {
		return nil, err
	}

	canonicalBytes, err := chain.CanonicalTxBytes(tx)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(canonicalBytes)

	r, s, err := ecdsa.Sign(rand.Reader, wallet.PrivateKey, hash[:])
	if err != nil {
		return nil, err
	}

	rBytes := r.Bytes()
	sBytes := s.Bytes()
	signatureBytes := append(rBytes, sBytes...)
	tx.Signature = hex.EncodeToString(signatureBytes)

	xBytes := wallet.PublicKey.X.Bytes()
	yBytes := wallet.PublicKey.Y.Bytes()
	pubKeyBytes := append(xBytes, yBytes...)
	tx.PubKey = hex.EncodeToString(pubKeyBytes)

	return tx, nil
}

func EncodePublicKey(pub *ecdsa.PublicKey) string {
	xBytes := pub.X.Bytes()
	yBytes := pub.Y.Bytes()
	combined := append(xBytes, yBytes...)
	return hex.EncodeToString(combined)
}

var (
	ErrWalletNotFound = &WalletError{Message: "wallet not found"}
	ErrInsufficientFunds = &WalletError{Message: "insufficient funds"}
)

type WalletError struct {
	Message string
}

func (e *WalletError) Error() string {
	return e.Message
}




