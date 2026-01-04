# AI-Blockchain Learning Project

A multi-language blockchain implementation for learning purposes, featuring:
- **Go**: Core blockchain node with Proof-of-Work consensus
- **Java**: Wallet and explorer service (Spring Boot)
- **Python**: AI scoring service for transaction anomaly detection

## Project Structure

```
ai-blockchain/
├── go-node/              # Core blockchain node (Go)
├── java-wallet/          # Wallet & explorer service (Java Spring Boot)
├── ai-scorer/            # AI scoring service (Python Flask)
└── schemas/              # Shared data schemas
```

## Features

### Go Node
- ✅ UTXO-based transaction model
- ✅ Proof-of-Work consensus
- ✅ Transaction validation
- ✅ Mempool management
- ✅ REST API
- ✅ AI integration (advisory scoring)

### Java Wallet
- ✅ Key pair generation (ECDSA)
- ✅ Transaction creation and signing
- ✅ Balance queries
- ✅ Integration with Go node

### Python AI Scorer
- ✅ Transaction anomaly detection (IsolationForest)
- ✅ Fee adequacy estimation
- ✅ REST API for scoring

## Getting Started

### Prerequisites
- Go 1.21+
- Java 17+
- Python 3.9+
- Maven 3.8+

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

### Running the Java Wallet

```bash
cd java-wallet
mvn spring-boot:run
```

## API Endpoints

### Go Node (port 8080)
- `GET /health` - Health check
- `GET /blocks` - Get all blocks
- `GET /chain` - Get blockchain info
- `GET /mempool` - Get pending transactions
- `GET /balance/:addr` - Get balance for address
- `POST /transactions` - Submit new transaction
- `POST /mine` - Mine a new block

### Java Wallet (port 8081)
- `GET /api/wallet/generate` - Generate new key pair
- `GET /api/wallet/balance/:address` - Get balance
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

2. **Cross-Language Integration**
   - Go ↔ Java communication
   - Shared data schemas
   - Cryptographic compatibility

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

