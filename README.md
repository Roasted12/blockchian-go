# AI-Blockchain Learning Project

A multi-language blockchain implementation for learning purposes, featuring:
- **Go**: Core blockchain node with Proof-of-Work consensus + built-in wallet
- **Python**: AI scoring service for transaction anomaly detection

## Project Structure

```
ai-blockchain/
├── go-node/              # Core blockchain node + wallet (Go)
├── ai-scorer/            # AI scoring service (Python Flask)
├── web-ui/               # Web interface (HTML/JS)
└── schemas/              # Shared data schemas
```

## Features

### Go Node
- ✅ UTXO-based transaction model
- ✅ Proof-of-Work consensus
- ✅ Transaction validation
- ✅ Mempool management
- ✅ REST API
- ✅ **Built-in wallet service** (key generation, transaction signing)
- ✅ AI integration (advisory scoring)

### Python AI Scorer
- ✅ Transaction anomaly detection (IsolationForest)
- ✅ Fee adequacy estimation
- ✅ REST API for scoring

## Getting Started

### Prerequisites
- Go 1.21+
- Python 3.9+

### Running the Go Node

```bash
cd go-node
go mod tidy
go run cmd/node/main.go -port 8080 -difficulty 4
```

With AI scoring:
```bash
go run cmd/node/main.go -port 8080 -difficulty 4 -ai-url http://localhost:5000 -ai-timeout 5
```

### Running the Python AI Scorer

```bash
cd ai-scorer
pip install -r requirements.txt
python app/api.py
```

**Note**: Wallet functionality is now built into the Go node! No Java needed.

## API Endpoints

### Go Node (port 8080)
- `GET /health` - Health check
- `GET /blocks` - Get all blocks
- `GET /chain` - Get blockchain info
- `GET /mempool` - Get pending transactions
- `GET /balance/:addr` - Get balance for address
- `POST /transactions` - Submit new transaction
- `POST /mine` - Mine a new block
- `GET /api/wallet/generate` - Generate new wallet
- `GET /api/wallet/list` - List all wallets
- `POST /api/wallet/transfer` - Create and submit transaction

### Python AI Scorer (port 5000)
- `GET /health` - Health check
- `POST /score/tx` - Score transaction

## Learning Objectives

1. **Blockchain Fundamentals**
   - Block structure and chaining
   - Transaction validation
   - UTXO model
   - Proof-of-Work consensus

2. **Cryptographic Operations**
   - ECDSA key generation and signing
   - Transaction signing with private keys
   - Address derivation from public keys

3. **AI Integration**
   - Feature extraction
   - Anomaly detection
   - Advisory scoring (non-consensus)

4. **System Design**
   - Microservices architecture
   - REST API design
   - Error handling and resilience

## Notes

- This is a **learning project** - not production-ready
- AI scoring is **advisory only** - does not affect consensus
- Simplified implementations for educational purposes
- Extensive comments throughout codebase

## License

Educational use only.

