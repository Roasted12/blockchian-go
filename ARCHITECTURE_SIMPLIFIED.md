# Simplified Architecture (Go-Only Wallet)

## Overview

The wallet functionality has been moved into the Go node, simplifying the architecture significantly.

## New Architecture

```
┌─────────────┐
│  Web UI     │  ← User interface
└──────┬──────┘
       │
       │ REST API
       ▼
┌─────────────────────┐     ┌─────────────┐
│     Go Node         │────▶│ AI Scorer   │
│   (Port 8080)       │     │  (Port 5000)│
│                     │     └─────────────┘
│  Components:        │
│  - Blockchain       │
│  - Wallet Service   │  ← Now built-in!
│  - Mining           │
│  - API Server       │
└─────────────────────┘
```

## What Changed

### Before (Multi-Language)
- Go Node: Blockchain operations
- Java Wallet: Private key operations
- Python AI: Scoring
- **3 services to run**

### After (Simplified)
- Go Node: Blockchain + Wallet operations
- Python AI: Scoring
- **2 services to run**

## Benefits

1. **Simpler Setup**
   - One less service to run
   - No Java/Maven dependencies
   - Faster startup

2. **Easier Development**
   - All code in Go
   - No cross-language debugging
   - Single codebase

3. **Still Secure**
   - Private keys still isolated in wallet package
   - API endpoints separate from core blockchain
   - Same security model

## Wallet Package Structure

```
go-node/internal/wallet/
├── wallet.go          # Wallet storage and operations
└── (private keys managed here)

go-node/internal/api/
├── wallet_handler.go   # Wallet API endpoints
└── server.go           # Main API server
```

## API Endpoints (All on Go Node)

### Blockchain Operations
- `GET /blocks` - Get all blocks
- `GET /chain` - Get blockchain info
- `GET /mempool` - Get pending transactions
- `POST /mine` - Mine a block
- `GET /balance/:addr` - Get balance

### Wallet Operations
- `GET /api/wallet/generate` - Generate new wallet
- `GET /api/wallet/list` - List all wallets
- `POST /api/wallet/transfer` - Create and submit transaction

## Migration Notes

- **Java wallet is now optional** - can be removed if desired
- **Web UI updated** - now calls Go node directly
- **Same functionality** - wallet features work the same way
- **Simpler deployment** - one less service to manage

## Why This Makes Sense

For a learning project:
- ✅ Simpler architecture
- ✅ Less moving parts
- ✅ Easier to understand
- ✅ Still demonstrates key concepts

For production:
- You might still want separate wallet service
- But for learning, built-in is fine
- Private keys are still properly isolated

