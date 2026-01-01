# Java Wallet - Complete Explanation

## Overview

The Java wallet is a **Spring Boot microservice** that handles **private key operations** - the things the Go blockchain node should NEVER do. It's the "wallet" part of the system, similar to how Bitcoin Core (node) and Bitcoin wallets work separately.

## Architecture

```
┌─────────────┐
│  Web UI     │  ← User interface
└──────┬──────┘
       │
       │ REST API
       ▼
┌─────────────┐
│ Java Wallet │  ← PRIVATE KEY OPERATIONS
│  (Port 8081)│     - Stores private keys
│             │     - Signs transactions
│             │     - Builds transactions
└──────┬──────┘
       │
       │ REST API (signed transactions only)
       ▼
┌─────────────┐
│  Go Node    │  ← PUBLIC OPERATIONS ONLY
│  (Port 8080)│     - Validates signatures
│             │     - Never sees private keys
└─────────────┘
```

## Key Principle

**Private keys NEVER leave the Java wallet service!**

- ✅ Java wallet: Has private keys, signs transactions
- ✅ Go node: Only verifies signatures (public key operations)
- ❌ Go node: NEVER receives private keys

This is how real blockchains work (Bitcoin, Ethereum, etc.)

---

## Components Breakdown

### 1. WalletService.java - Core Wallet Logic

**Purpose**: Manages wallets and transaction building

#### What it does:

**A. Wallet Storage**
```java
private final Map<String, WalletEntry> wallets = new ConcurrentHashMap<>();
```
- Stores wallets in memory (address → private key mapping)
- Each wallet contains:
  - Private key (for signing)
  - Public key (for verification)
  - Address (derived from public key)
  - Key pair (for convenience)

**B. Generate Wallet**
```java
public WalletInfo generateWallet()
```
1. Generates ECDSA key pair (using KeyGenerator)
2. Derives address from public key (SHA256 hash)
3. Stores private key securely in memory
4. Returns public info (address, public key)

**C. Build and Sign Transaction**
```java
public Transaction buildAndSignTransaction(fromAddress, toAddress, amount)
```

**Step-by-step process:**

1. **Lookup Wallet**
   - Finds wallet by address
   - Retrieves stored private key
   - Throws error if wallet doesn't exist

2. **Build Transaction Structure**
   - Creates `Transaction` object
   - Sets inputs (UTXOs to spend) - currently simplified
   - Sets outputs (recipient + change)
   - Sets timestamp

3. **Compute Transaction ID**
   - Creates canonical serialization (must match Go node!)
   - Sorts inputs and outputs deterministically
   - Hashes the canonical bytes
   - This ID is what gets signed

4. **Sign Transaction**
   - Gets canonical bytes (same as for ID)
   - Hashes the bytes (SHA-256)
   - Signs hash with private key (ECDSA)
   - Encodes signature as hex string

5. **Return Signed Transaction**
   - Transaction now has:
     - ID (hash of inputs+outputs)
     - Signature (proof of ownership)
     - Public key (for verification)

**Important Note**: The current implementation uses placeholder UTXOs. In production, you would:
- Query Go node for available UTXOs
- Select UTXOs that belong to `fromAddress`
- Calculate proper change output
- Build real transaction inputs

---

### 2. WalletController.java - REST API

**Purpose**: HTTP endpoints for wallet operations

#### Endpoints:

**A. `GET /api/wallet/generate`**
- Generates a new wallet
- Stores private key in WalletService
- Returns address and public key (NOT private key!)

**B. `GET /api/wallet/list`**
- Lists all wallet addresses
- Does NOT return private keys (security!)

**C. `POST /api/wallet/transfer`**
- Request: `{ "from": "addr1", "to": "addr2", "amount": 10.5 }`
- Process:
  1. Validates request
  2. Calls `WalletService.buildAndSignTransaction()`
  3. Gets signed transaction
  4. Submits to Go node via `GoNodeClient`
  5. Returns result

**D. `GET /api/wallet/balance/{address}`**
- Queries Go node for balance
- Just a proxy (balance is stored in Go node's UTXO set)

---

### 3. KeyGenerator.java - Cryptography

**Purpose**: Generate ECDSA key pairs

#### What it does:

**A. Generate Key Pair**
```java
public KeyPair generateKeyPair()
```
- Uses BouncyCastle library
- Curve: secp256r1 (P-256) - same as Go node
- Generates random private/public key pair
- Returns Java `KeyPair` object

**B. Derive Address**
```java
public String deriveAddress(byte[] publicKeyBytes)
```
- Takes public key bytes
- Hashes with SHA-256
- Returns hex-encoded hash (this is the address)
- Simplified version (production might use RIPEMD160)

**Why this matters:**
- Address = hash of public key
- You can share address without revealing public key
- Public key is only revealed when spending (in transaction)

---

### 4. ECDSASigner.java - Transaction Signing

**Purpose**: Sign transactions with private keys

#### What it does:

**A. Sign Data**
```java
public String sign(PrivateKey privateKey, byte[] data)
```
1. Hashes the data (SHA-256)
2. Signs the hash with private key (ECDSA)
3. Encodes signature as hex string
4. Returns signature

**B. Encode Public Key**
```java
public String encodePublicKey(PublicKey publicKey)
```
- Extracts X and Y coordinates from elliptic curve point
- Concatenates them: `x || y`
- Encodes as hex string
- Must match Go node's encoding format!

**Important**: The canonical serialization MUST match Go node exactly, or signatures won't verify!

---

### 5. GoNodeClient.java - Communication

**Purpose**: Talk to Go blockchain node

#### What it does:

- Uses Spring WebFlux (reactive HTTP client)
- Sends signed transactions to Go node
- Queries balances, blocks, chain info
- Never sends private keys (only signed transactions!)

---

## Data Flow: Creating a Transaction

### Step-by-Step:

1. **User Request** (Web UI)
   ```
   POST /api/wallet/transfer
   {
     "from": "address1",
     "to": "address2", 
     "amount": 10.0
   }
   ```

2. **WalletController** receives request
   - Validates input
   - Calls `WalletService.buildAndSignTransaction()`

3. **WalletService** builds transaction
   - Looks up wallet by `fromAddress`
   - Gets private key from storage
   - Creates transaction structure:
     ```java
     Transaction {
       inputs: [UTXO references],
       outputs: [
         { address: "address2", amount: 10.0 },
         { address: "address1", amount: change }
       ]
     }
     ```
   - Computes transaction ID (hash of inputs+outputs)
   - Signs transaction with private key
   - Returns signed transaction

4. **WalletController** submits to Go node
   - Uses `GoNodeClient.submitTransaction()`
   - Sends signed transaction (NO private key!)
   - Go node receives: `{ id, inputs, outputs, signature, pubkey }`

5. **Go Node** validates
   - Verifies signature (public key operation)
   - Checks UTXOs exist
   - Validates value conservation
   - Adds to mempool if valid

---

## Security Model

### What Java Wallet Stores:
- ✅ Private keys (in memory)
- ✅ Public keys
- ✅ Addresses
- ✅ Key pairs

### What Java Wallet NEVER Sends:
- ❌ Private keys (never!)
- ❌ Raw key material

### What Java Wallet Sends:
- ✅ Signed transactions (signature + public key)
- ✅ Public information (addresses, balances)

### What Go Node Receives:
- ✅ Signed transactions
- ✅ Public keys (for verification)
- ✅ Addresses

### What Go Node NEVER Receives:
- ❌ Private keys
- ❌ Unencrypted key material

---

## Current Limitations (Learning Project)

### 1. UTXO Selection
- Currently uses placeholder UTXOs
- In production: Query Go node for actual UTXOs
- Select UTXOs that belong to sender
- Calculate proper change

### 2. Storage
- Private keys stored in memory
- Lost when service restarts
- Production: Use encrypted file storage or hardware wallets

### 3. Transaction Building
- Simplified transaction structure
- Doesn't query actual blockchain state
- Production: Full UTXO query and selection

### 4. Change Calculation
- Currently set to 0.0
- Production: Calculate from input sum - amount - fee

---

## Why This Architecture?

### Separation of Concerns:
- **Go Node**: Consensus, validation, storage (public operations)
- **Java Wallet**: Key management, signing (private operations)

### Security:
- Private keys isolated from blockchain node
- Node compromise doesn't expose keys
- Wallet can be on different machine

### Scalability:
- Multiple wallets can connect to one node
- Wallet can connect to multiple nodes
- Microservices architecture

### Real-World Example:
- **Bitcoin Core**: Like Go node (validates, stores)
- **Electrum Wallet**: Like Java wallet (signs, manages keys)
- They communicate via RPC/API

---

## API Endpoints Summary

| Endpoint | Method | Purpose | Returns |
|----------|--------|---------|---------|
| `/api/wallet/generate` | GET | Create new wallet | Address, public key |
| `/api/wallet/list` | GET | List all wallets | Array of addresses |
| `/api/wallet/transfer` | POST | Send coins | Transaction result |
| `/api/wallet/balance/{addr}` | GET | Check balance | Balance amount |

---

## Key Takeaways

1. **Java wallet = Private key operations**
   - Generates keys
   - Stores keys
   - Signs transactions

2. **Go node = Public operations**
   - Validates signatures
   - Never sees private keys

3. **Communication = Signed transactions only**
   - Java wallet sends signed transactions
   - Go node verifies and stores

4. **Security = Separation**
   - Private keys stay in wallet
   - Node can't access keys
   - Even if node is compromised, keys are safe

This is exactly how real cryptocurrency wallets work!

