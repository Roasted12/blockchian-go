# System Architecture

## Overview

This is a multi-language blockchain learning project demonstrating:
- **Go**: Core blockchain node (consensus, validation, storage)
- **Java**: Wallet service (private key management, transaction signing)
- **Python**: AI scoring service (advisory transaction analysis)
- **Web UI**: User interface (HTML/JavaScript)

## Component Responsibilities

### Go Node (Port 8080)
**Role**: Blockchain core - public operations only

**Responsibilities**:
- ✅ Block creation and validation
- ✅ Transaction validation (signature verification)
- ✅ Proof-of-Work mining
- ✅ UTXO set management
- ✅ Mempool management
- ✅ Block storage
- ✅ REST API for blockchain operations

**What it does NOT do**:
- ❌ Store private keys
- ❌ Sign transactions
- ❌ Build transactions

### Java Wallet (Port 8081)
**Role**: Wallet service - private key operations

**Responsibilities**:
- ✅ Generate and store private keys
- ✅ Build transactions (UTXO selection, change calculation)
- ✅ Sign transactions with private keys
- ✅ Manage multiple wallets
- ✅ Provide wallet API

**What it does NOT do**:
- ❌ Validate blocks
- ❌ Mine blocks
- ❌ Store blockchain state

### Python AI Scorer (Port 5000)
**Role**: Advisory ML service

**Responsibilities**:
- ✅ Transaction anomaly detection
- ✅ Fee adequacy estimation
- ✅ Peer reliability scoring (future)

**Important**: AI scoring is **advisory only** - does not affect consensus

### Web UI
**Role**: User interface

**Responsibilities**:
- ✅ Display blockchain status
- ✅ Browse blocks and transactions
- ✅ Generate wallets
- ✅ Create and submit transactions
- ✅ Mine blocks

## Data Flow

### Creating a Transaction

```
1. User (Web UI)
   ↓
2. Java Wallet: Build transaction
   - Select UTXOs
   - Calculate change
   - Sign with private key
   ↓
3. Go Node: Validate transaction
   - Verify signature (public key)
   - Check UTXO availability
   - Validate value conservation
   ↓
4. AI Scorer (optional): Score transaction
   - Extract features
   - Detect anomalies
   - Return advisory score
   ↓
5. Go Node: Add to mempool
   ↓
6. Miner: Mine block with transactions
```

### Mining a Block

```
1. User/Node: Request to mine
   ↓
2. Go Node: Get transactions from mempool
   ↓
3. Go Node: Create block
   ↓
4. Go Node: Mine block (Proof-of-Work)
   ↓
5. Go Node: Add block to chain
   ↓
6. Go Node: Update UTXO set
   ↓
7. Go Node: Remove transactions from mempool
```

## Security Model

### Private Key Security
- **Java Wallet**: Stores private keys in memory (in-memory for learning)
- **Go Node**: Never receives private keys
- **Web UI**: Never receives private keys

### Transaction Signing
- Private keys stay in Java wallet
- Only signed transactions are sent to Go node
- Go node verifies signatures (public key operations only)

## API Communication

```
Web UI ↔ Go Node (8080)
  - GET /blocks
  - GET /chain
  - GET /mempool
  - GET /balance/:addr
  - POST /mine

Web UI ↔ Java Wallet (8081)
  - GET /api/wallet/generate
  - GET /api/wallet/list
  - POST /api/wallet/transfer

Go Node ↔ Java Wallet
  - Java wallet queries Go node for UTXOs
  - Java wallet submits signed transactions

Go Node ↔ AI Scorer (5000)
  - POST /score/tx (advisory scoring)
```

## Key Design Principles

1. **Separation of Concerns**
   - Node = public operations
   - Wallet = private operations
   - AI = advisory only

2. **Never Trust, Always Verify**
   - Go node validates everything
   - Signatures verified on every transaction

3. **Private Keys Never Leave Wallet**
   - Java wallet signs locally
   - Only signed transactions sent to node

4. **AI is Advisory**
   - Scoring doesn't affect consensus
   - Can be disabled without breaking blockchain

## Future Enhancements

- [ ] Persistent wallet storage (encrypted)
- [ ] HD wallet support (key derivation)
- [ ] P2P networking between nodes
- [ ] Block explorer with transaction history
- [ ] Multi-signature support
- [ ] Smart contracts (scripting)

