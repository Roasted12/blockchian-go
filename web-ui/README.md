# Blockchain Web UI

A simple, modern web interface for interacting with the AI blockchain.

## Features

- **Dashboard**: View blockchain status, height, mempool size, difficulty
- **Blocks**: Browse all blocks in the chain
- **Transactions**: View pending transactions in mempool
- **Wallet**: Generate addresses, check balances, create transactions
- **Mine**: Mine new blocks from pending transactions

## Usage

### Important: CORS Issue
**Do NOT open `index.html` directly from the file system!** 
Browsers block fetch requests from `file://` URLs due to CORS security.

### Option 1: Use the Python Server (Recommended)
```bash
cd web-ui
python server.py
# Then open http://localhost:3000
```

### Option 2: Use Python's Built-in Server
```bash
cd web-ui
python -m http.server 3000
# Then open http://localhost:3000
```

### Option 3: Use Node.js (if you have it)
```bash
cd web-ui
npx http-server -p 3000 --cors
# Then open http://localhost:3000
```

## Prerequisites

1. **Go node** running on port 8080
2. **Java wallet** running on port 8081
3. **Python AI scorer** running on port 5000 (optional)

## Troubleshooting

### "Failed to fetch" Error
- **Cause**: Opening HTML file directly (file://) or CORS not enabled
- **Solution**: Use a web server (see options above)

### "Connection refused" Error
- **Cause**: Go node or Java wallet not running
- **Solution**: Start all services first (see START_HERE.md)

## Architecture

```
Web UI (Browser)
    │
    ├──→ Go Node (8080) - Blockchain operations
    │    - View blocks
    │    - View mempool
    │    - Mine blocks
    │    - Check balances
    │
    └──→ Java Wallet (8081) - Wallet operations
         - Generate wallets
         - Create transactions
         - Sign transactions
```

## API Endpoints Used

### Go Node
- `GET /chain` - Get blockchain info
- `GET /blocks` - Get all blocks
- `GET /mempool` - Get pending transactions
- `GET /balance/:addr` - Get balance
- `POST /mine` - Mine block

### Java Wallet
- `GET /api/wallet/generate` - Generate wallet
- `GET /api/wallet/list` - List wallets
- `POST /api/wallet/transfer` - Create and submit transaction

