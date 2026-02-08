# AI-Blockchain

A multi-service blockchain prototype with a Go node, Java wallet service, and Python AI scoring service. The Go node provides a UTXO-based ledger, Proof-of-Work mining, and a REST API. The wallet service handles key management and transaction creation, while the AI service provides advisory transaction scoring.

## Components

- **Go Node**: blockchain core, mempool, PoW mining, REST API
- **Java Wallet**: key generation, transaction creation/signing
- **Python AI Scorer**: anomaly and fee adequacy scoring

## Project Layout

```
ai-blockchain/
├── go-node/
├── java-wallet/
├── ai-scorer/
└── schemas/
```

## Requirements

- Go 1.21+
- Java 17+
- Python 3.9+
- Maven 3.8+

## Quick Start

### Go Node
```bash
cd go-node
go mod tidy
go run cmd/node/main.go -port 8080 -difficulty 4
```

With AI scoring:
```bash
go run cmd/node/main.go -port 8080 -difficulty 4 -ai-url http://localhost:5000 -ai-timeout 5
```

### Python AI Scorer
```bash
cd ai-scorer
pip install -r requirements.txt
python app/api.py
```

### Java Wallet
```bash
cd java-wallet
mvn spring-boot:run
```

## API Endpoints

### Go Node (8080)
- `GET /health`
- `GET /blocks`
- `GET /chain`
- `GET /mempool`
- `GET /balance/:addr`
- `POST /transactions`
- `POST /mine`

### Java Wallet (8081)
- `GET /api/wallet/generate`
- `GET /api/wallet/balance/:address`
- `POST /api/wallet/transfer`

### Python AI Scorer (5000)
- `GET /health`
- `POST /score/tx`
