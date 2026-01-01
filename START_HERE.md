# Quick Start Guide

## Step 1: Set Up Python Environment (Conda)

### Option A: Using Conda (Recommended)
```bash
cd ai-scorer
conda env create -f environment.yml
conda activate ai-blockchain-scorer
python app/api.py
```

### Option B: Using pip
```bash
cd ai-scorer
pip install -r requirements.txt
python app/api.py
```

**Note:** See `ai-scorer/SETUP_CONDA.md` for detailed Conda setup instructions.

The AI scorer will start on port 5000.

## Step 2: Start the Go Blockchain Node

```bash
cd go-node
go mod tidy
go run cmd/node/main.go -port 8080 -difficulty 4 -ai-url http://localhost:5000
```

The node will:
- Create a genesis block
- Start the REST API on port 8080
- Connect to AI scorer for transaction scoring

## Step 3: Open the Web UI

**Note**: Wallet functionality is now built into the Go node! No Java needed.

### ⚠️ Important: Use a Web Server
**Do NOT open `index.html` directly!** Use a web server to avoid CORS errors.

### Option A: Use the Python Server (Recommended)
```bash
cd web-ui
python server.py
# Then open http://localhost:3000
```

### Option B: Use Python's Built-in Server
```bash
cd web-ui
python -m http.server 3000
# Then open http://localhost:3000
```

### Option C: Use Node.js (if installed)
```bash
cd web-ui
npx http-server -p 3000 --cors
# Then open http://localhost:3000
```

The web UI provides a complete interface for:
- 📊 Viewing blockchain status
- 🔗 Browsing blocks
- 💼 Managing wallets
- 💸 Creating transactions
- ⛏️ Mining blocks

**Note**: CORS is now enabled on both Go node and Java wallet APIs.

## Testing the System

### 1. Check Go Node Health
```bash
curl http://localhost:8080/health
```

### 2. Generate a Wallet Address
```bash
curl http://localhost:8080/api/wallet/generate
```

### 3. Check Balance
```bash
curl http://localhost:8080/balance/0000000000000000000000000000000000000000
```

### 4. Mine a Block
```bash
curl -X POST http://localhost:8080/mine
```

### 5. View All Blocks
```bash
curl http://localhost:8080/blocks
```

## Architecture

```
┌─────────────┐
│  Web UI     │  ← User interface
└──────┬──────┘
       │
       │ REST API
       ▼
┌─────────────┐     ┌─────────────┐
│  Go Node    │────▶│ AI Scorer   │  ← Scores transactions
│  (Port 8080)│     │  (Port 5000)│
│             │     └─────────────┘
│  - Blockchain
│  - Wallet   │  ← Wallet functionality built-in!
│  - Mining
└─────────────┘
       │
       │ Validates, mines blocks
       ▼
  Blockchain
```

## Next Steps

1. Create transactions from Java wallet
2. Submit them to Go node
3. Mine blocks to confirm transactions
4. Observe AI scoring in action

## Troubleshooting

- **Go node fails to start**: Check if port 8080 is available
- **Java wallet can't connect**: Verify Go node is running on port 8080
- **AI scorer not responding**: Check Python dependencies are installed
- **Transactions rejected**: Check transaction format matches Go node expectations

