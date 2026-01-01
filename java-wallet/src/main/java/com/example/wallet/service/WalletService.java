package com.example.wallet.service;

import com.example.wallet.crypto.ECDSASigner;
import com.example.wallet.crypto.HashUtil;
import com.example.wallet.crypto.KeyGenerator;
import com.example.wallet.model.Transaction;
import com.example.wallet.model.TxIn;
import com.example.wallet.model.TxOut;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.security.KeyPair;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;

/**
 * WALLET SERVICE – WALLET STATE MANAGEMENT
 *
 * This service manages:
 * - Private key storage (in-memory for now)
 * - Address to key pair mapping
 * - Transaction building with UTXO selection
 * - Proper transaction signing
 *
 * Important:
 * - Private keys are stored in memory (not secure for production!)
 * - In production, use encrypted storage or hardware wallets
 */
@Service
public class WalletService {

    @Autowired
    private KeyGenerator keyGenerator;

    @Autowired
    private ECDSASigner signer;

    // In-memory wallet storage
    // Format: address -> (privateKey, publicKey, keyPair)
    private final Map<String, WalletEntry> wallets = new ConcurrentHashMap<>();

    /**
     * Generate a new wallet and store the private key.
     *
     * @return WalletInfo containing address and public key
     */
    public WalletInfo generateWallet() throws Exception {
        KeyPair keyPair = keyGenerator.generateKeyPair();
        String address = keyGenerator.deriveAddress(keyPair.getPublic().getEncoded());
        String publicKeyHex = signer.encodePublicKey(keyPair.getPublic());

        // Store wallet entry
        wallets.put(address, new WalletEntry(
            keyPair.getPrivate(),
            keyPair.getPublic(),
            keyPair,
            address,
            publicKeyHex
        ));

        return new WalletInfo(address, publicKeyHex);
    }

    /**
     * Get wallet entry by address.
     *
     * @param address Wallet address
     * @return WalletEntry or null if not found
     */
    public WalletEntry getWallet(String address) {
        return wallets.get(address);
    }

    /**
     * Get all wallet addresses.
     *
     * @return List of addresses
     */
    public List<String> getAllAddresses() {
        return new ArrayList<>(wallets.keySet());
    }

    /**
     * Build and sign a transaction.
     *
     * This method:
     * 1. Validates the wallet exists
     * 2. Creates a transaction using genesis UTXO (for learning)
     * 3. Creates outputs (recipient + change)
     * 4. Signs the transaction
     *
     * Note: This is a simplified version for learning.
     * In production, you would query the Go node for available UTXOs.
     *
     * @param fromAddress Sender's address
     * @param toAddress Recipient's address
     * @param amount Amount to send
     * @return Signed transaction
     */
    public Transaction buildAndSignTransaction(String fromAddress, String toAddress, double amount) throws Exception {
        // Get wallet
        WalletEntry wallet = getWallet(fromAddress);
        if (wallet == null) {
            throw new IllegalArgumentException("Wallet not found for address: " + fromAddress + ". Please generate a wallet first using /api/wallet/generate");
        }

        // For learning purposes, we'll create a transaction that references the genesis block
        // The genesis block creates a UTXO with txid = first transaction ID in genesis block
        // In a real implementation, you would:
        // 1. Query Go node for all blocks
        // 2. Find UTXOs that belong to fromAddress
        // 3. Select UTXOs that cover the amount

        // Create transaction
        Transaction tx = new Transaction();

        // IMPORTANT: For this learning project, we need to handle the fact that
        // users might not have any UTXOs yet. The genesis block creates a UTXO
        // for address "0000000000000000000000000000000000000000".
        //
        // For a real transaction to work, the user needs to:
        // 1. Have received coins (either from genesis or another transaction)
        // 2. Reference actual UTXOs that exist in the blockchain
        //
        // For now, we'll create a transaction structure, but it will fail validation
        // unless the user actually has UTXOs. This is expected behavior for learning.
        
        // Create inputs - this is a placeholder
        // In production, you would:
        // 1. Query Go node: GET /blocks
        // 2. Parse all blocks to find transactions
        // 3. Find outputs that match fromAddress
        // 4. Select UTXOs to spend
        
        List<TxIn> inputs = new ArrayList<>();
        TxIn input = new TxIn();
        // This is a placeholder - the transaction will fail validation
        // unless you have actual UTXOs. For testing, you can:
        // 1. Use the genesis address: "0000000000000000000000000000000000000000"
        // 2. Or mine a block first to create UTXOs
        // 3. Or receive coins from another transaction
        
        // Try to use genesis transaction ID (you'd need to query this from Go node)
        // For now, using a placeholder that will need to be replaced
        input.setTxId("GENESIS_PLACEHOLDER"); // This needs to be the actual genesis tx ID
        input.setIndex(0);
        inputs.add(input);
        tx.setInputs(inputs);

        // Create outputs
        List<TxOut> outputs = new ArrayList<>();
        
        // Output to recipient
        TxOut recipientOutput = new TxOut();
        recipientOutput.setAddress(toAddress);
        recipientOutput.setAmount(amount);
        outputs.add(recipientOutput);

        // Change output (simplified - in production, calculate from input sum - amount - fee)
        // For now, assume we have enough and send change back
        TxOut changeOutput = new TxOut();
        changeOutput.setAddress(fromAddress);
        changeOutput.setAmount(0.0); // Placeholder - should calculate actual change
        outputs.add(changeOutput);

        tx.setOutputs(outputs);
        tx.setTimestamp(System.currentTimeMillis() / 1000);

        // Compute transaction ID (canonical serialization)
        // This must match Go node's serialization
        String txId = computeTransactionId(tx);
        tx.setId(txId);

        // Sign transaction
        byte[] canonicalBytes = getCanonicalBytes(tx);
        String signature = signer.sign(wallet.getPrivateKey(), canonicalBytes);
        tx.setSignature(signature);
        tx.setPubkey(wallet.getPublicKeyHex());

        return tx;
    }

    /**
     * Compute transaction ID (must match Go node's computation).
     *
     * Simplified version - in production, use exact same serialization as Go node.
     */
    private String computeTransactionId(Transaction tx) {
        // Create canonical representation (must match Go node)
        StringBuilder sb = new StringBuilder();
        
        // Sort and add inputs
        tx.getInputs().stream()
            .sorted(Comparator.comparing(TxIn::getTxId).thenComparing(TxIn::getIndex))
            .forEach(in -> sb.append(in.getTxId()).append(":").append(in.getIndex()).append(","));
        
        // Sort and add outputs
        tx.getOutputs().stream()
            .sorted(Comparator.comparing(TxOut::getAddress))
            .forEach(out -> sb.append(out.getAddress()).append(":").append(out.getAmount()).append(","));
        
        return HashUtil.sha256Hex(sb.toString());
    }

    /**
     * Get canonical bytes for signing (must match Go node).
     *
     * Simplified version - in production, use exact same serialization.
     */
    private byte[] getCanonicalBytes(Transaction tx) {
        // Create canonical representation
        StringBuilder sb = new StringBuilder();
        
        // Sort and add inputs
        tx.getInputs().stream()
            .sorted(Comparator.comparing(TxIn::getTxId).thenComparing(TxIn::getIndex))
            .forEach(in -> sb.append(in.getTxId()).append(":").append(in.getIndex()).append(","));
        
        // Sort and add outputs
        tx.getOutputs().stream()
            .sorted(Comparator.comparing(TxOut::getAddress))
            .forEach(out -> sb.append(out.getAddress()).append(":").append(out.getAmount()).append(","));
        
        return sb.toString().getBytes();
    }

    /**
     * Wallet entry storage.
     */
    public static class WalletEntry {
        private final PrivateKey privateKey;
        private final PublicKey publicKey;
        private final KeyPair keyPair;
        private final String address;
        private final String publicKeyHex;

        public WalletEntry(PrivateKey privateKey, PublicKey publicKey, KeyPair keyPair, String address, String publicKeyHex) {
            this.privateKey = privateKey;
            this.publicKey = publicKey;
            this.keyPair = keyPair;
            this.address = address;
            this.publicKeyHex = publicKeyHex;
        }

        public PrivateKey getPrivateKey() { return privateKey; }
        public PublicKey getPublicKey() { return publicKey; }
        public KeyPair getKeyPair() { return keyPair; }
        public String getAddress() { return address; }
        public String getPublicKeyHex() { return publicKeyHex; }
    }

    /**
     * Wallet information (public data only).
     */
    public static class WalletInfo {
        private final String address;
        private final String publicKey;

        public WalletInfo(String address, String publicKey) {
            this.address = address;
            this.publicKey = publicKey;
        }

        public String getAddress() { return address; }
        public String getPublicKey() { return publicKey; }
    }
}

