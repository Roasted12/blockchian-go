# Why Java Wallet Exists

## Current State (What It Does Now)
The Java wallet currently:
- ✅ Generates key pairs (ECDSA)
- ✅ Signs transactions (stub/incomplete)
- ⚠️ Just proxies balance queries to Go node
- ⚠️ Doesn't store private keys
- ⚠️ Doesn't properly build transactions

## The Problem
You're right - **the Go node can do most of this**. So why have Java?

## The Real Purpose (What It SHOULD Do)

### 1. **Cryptographic Operations - Private Key Management**
- **Go node**: Only VERIFIES signatures (public key operations)
- **Java wallet**: GENERATES and STORES private keys (private key operations)
- **Why**: Private keys should NEVER be on the blockchain node (security risk)
- **Separation**: Node = public, Wallet = private

### 2. **Transaction Building**
- **Go node**: Validates transactions
- **Java wallet**: BUILDS transactions by:
  - Selecting which UTXOs to spend
  - Calculating change outputs
  - Properly constructing inputs/outputs
  - Signing with private keys

### 3. **Wallet State Management**
- Store multiple addresses per user
- Track balances locally (cached)
- Manage transaction history
- Handle key derivation (HD wallets)

### 4. **Cross-Language Learning**
- Demonstrates microservices architecture
- Shows how different languages communicate
- Proves cryptographic compatibility (Go ↔ Java)

## What Should Be Enhanced

1. **Private Key Storage** (in-memory for now, encrypted file later)
2. **Proper Transaction Signing** (complete the stub code)
3. **UTXO Selection** (query Go node, select inputs, build transaction)
4. **Wallet Management** (multiple addresses, key derivation)

## Architecture

```
┌─────────────┐
│  Web UI     │  ← User interface
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Java Wallet │  ← PRIVATE KEY OPERATIONS
│  (Port 8081)│     - Key generation
│             │     - Transaction signing
│             │     - Wallet management
└──────┬──────┘
       │
       │ REST API (public operations only)
       ▼
┌─────────────┐
│  Go Node    │  ← PUBLIC OPERATIONS
│  (Port 8080)│     - Validation
│             │     - Mining
│             │     - Block storage
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ AI Scorer   │  ← Advisory scoring
│  (Port 5000)│
└─────────────┘
```

## Key Principle
**Never send private keys to the blockchain node!**

- Java wallet: Has private keys, signs transactions
- Go node: Never sees private keys, only verifies signatures

This is how real blockchains work (Bitcoin, Ethereum, etc.)

