# Wallet Usage Guide

## How to Send and Receive Coins

### Step 1: Generate a Wallet
1. Go to the **Wallet** tab
2. Click **"➕ Generate New Wallet"**
3. Your wallet address will appear in the list
4. The wallet is stored in the Java wallet service (private keys are secure)

### Step 2: Receive Coins

**Option A: From Genesis (Easiest for Testing)**
1. In the **Receive** section, select your wallet
2. Copy your address
3. Go to **Send** section
4. **From:** Select or paste: `0000000000000000000000000000000000000000` (genesis address - has 1000 coins)
5. **To:** Paste your wallet address
6. **Amount:** Enter amount (e.g., 100)
7. Click **"🚀 Send Transaction"**

**Option B: From Another Wallet**
1. Generate multiple wallets
2. Send from one wallet to another using the Send section

### Step 3: Send Coins
1. Go to **Send** section
2. **From Wallet:** Select the wallet you want to send from
3. **To Address:** Enter recipient address (or use "📋 Use Another Wallet" button)
4. **Amount:** Enter the amount
5. Click **"🚀 Send Transaction"**

## Understanding the UI

### Your Wallets Section
- Shows all wallets you've generated
- Displays balance for each wallet
- Click **"📋 Copy"** to copy an address

### Send Section
- **From Wallet:** Dropdown of all your wallets with balances
- **To Address:** Recipient's address
- **Amount:** How much to send
- Shows available balance when you select a wallet

### Receive Section
- **Your Address:** Select a wallet to see its receive address
- **Copy Address:** Share this address to receive coins
- Shows current balance

## Quick Tips

1. **Genesis Address:** `0000000000000000000000000000000000000000` has 1000 coins - use it for testing!

2. **Check Balances:** Click **"💰 Check All Balances"** to refresh all wallet balances

3. **Copy Genesis:** Click **"📋 Copy Genesis Address"** for quick access

4. **Transaction Status:** After sending, check the **Transactions** tab to see pending transactions

5. **Mine Blocks:** Go to **Mine** tab to confirm transactions by mining blocks

## Example Flow

1. **Generate Wallet 1**
   - Click "Generate New Wallet"
   - Copy the address

2. **Send from Genesis to Wallet 1**
   - Send section: From = genesis address, To = Wallet 1 address, Amount = 100
   - Click Send

3. **Mine a Block**
   - Go to Mine tab
   - Click "Mine Block"
   - This confirms the transaction

4. **Check Balance**
   - Wallet 1 should now show 100 coins

5. **Send from Wallet 1 to Wallet 2**
   - Generate Wallet 2
   - Send section: From = Wallet 1, To = Wallet 2, Amount = 50
   - Click Send
   - Mine another block to confirm

## Troubleshooting

**"Wallet not found" error:**
- Make sure you've generated a wallet first
- Use the "Generate New Wallet" button

**"UTXO not found" error:**
- The wallet doesn't have any coins yet
- Send coins from genesis address first
- Or mine a block to create new UTXOs

**Transaction pending:**
- Transactions go to mempool first
- Mine a block to confirm them
- Check Transactions tab to see pending transactions

