package com.example.wallet.controller;

import com.example.wallet.client.GoNodeClient;
import com.example.wallet.model.Transaction;
import com.example.wallet.service.WalletService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.HashMap;
import java.util.Map;

/**
 * WALLET CONTROLLER – WALLET OPERATIONS
 *
 * Provides endpoints for:
 * - Generating wallets (with private key storage)
 * - Creating and signing transactions
 * - Submitting transactions to blockchain
 * - Querying balances
 * - Listing wallets
 */
@RestController
@RequestMapping("/api/wallet")
public class WalletController {

    @Autowired
    private WalletService walletService;

    @Autowired
    private GoNodeClient nodeClient;

    /**
     * Generate a new wallet (address + private key storage).
     *
     * GET /api/wallet/generate
     *
     * @return Wallet info (address and public key)
     */
    @GetMapping("/generate")
    public ResponseEntity<Map<String, Object>> generateWallet() {
        try {
            WalletService.WalletInfo wallet = walletService.generateWallet();
            
            Map<String, Object> response = new HashMap<>();
            response.put("address", wallet.getAddress());
            response.put("public_key", wallet.getPublicKey());
            response.put("message", "Wallet generated and stored successfully");
            response.put("note", "Private key is stored securely in wallet service");
            response.put("tip", "To receive coins, share your address or mine a block");
            
            return ResponseEntity.ok(response);
        } catch (Exception e) {
            Map<String, Object> error = new HashMap<>();
            error.put("error", e.getMessage());
            return ResponseEntity.internalServerError().body(error);
        }
    }

    /**
     * List all wallet addresses.
     *
     * GET /api/wallet/list
     *
     * @return List of wallet addresses
     */
    @GetMapping("/list")
    public ResponseEntity<Map<String, Object>> listWallets() {
        try {
            Map<String, Object> response = new HashMap<>();
            response.put("addresses", walletService.getAllAddresses());
            response.put("count", walletService.getAllAddresses().size());
            return ResponseEntity.ok(response);
        } catch (Exception e) {
            Map<String, Object> error = new HashMap<>();
            error.put("error", e.getMessage());
            return ResponseEntity.internalServerError().body(error);
        }
    }

    /**
     * Create and submit a transaction.
     *
     * POST /api/wallet/transfer
     *
     * Request body:
     * {
     *   "from": "address1",
     *   "to": "address2",
     *   "amount": 10.5
     * }
     *
     * The wallet service will:
     * 1. Look up the wallet by address
     * 2. Build the transaction with proper inputs/outputs
     * 3. Sign it with the stored private key
     * 4. Submit to Go node
     *
     * @param request Transaction request
     * @return Transaction result
     */
    @PostMapping("/transfer")
    public ResponseEntity<Map<String, Object>> transfer(@RequestBody Map<String, Object> request) {
        try {
            String fromAddress = (String) request.get("from");
            String toAddress = (String) request.get("to");
            Object amountObj = request.get("amount");
            
            if (fromAddress == null || toAddress == null || amountObj == null) {
                Map<String, Object> error = new HashMap<>();
                error.put("error", "Invalid request: from, to, and amount are required");
                return ResponseEntity.badRequest().body(error);
            }
            
            double amount;
            if (amountObj instanceof Number) {
                amount = ((Number) amountObj).doubleValue();
            } else {
                amount = Double.parseDouble(amountObj.toString());
            }
            
            if (amount <= 0) {
                Map<String, Object> error = new HashMap<>();
                error.put("error", "Amount must be positive");
                return ResponseEntity.badRequest().body(error);
            }
            
            // Build and sign transaction using wallet service
            Transaction tx;
            try {
                tx = walletService.buildAndSignTransaction(fromAddress, toAddress, amount);
            } catch (IllegalArgumentException e) {
                Map<String, Object> error = new HashMap<>();
                error.put("error", e.getMessage());
                error.put("hint", "Make sure you've generated a wallet first using /api/wallet/generate");
                return ResponseEntity.badRequest().body(error);
            }
            
            // Submit to Go node
            String result = nodeClient.submitTransaction(tx).block();
            
            Map<String, Object> response = new HashMap<>();
            response.put("status", "submitted");
            response.put("txid", tx.getId());
            response.put("result", result);
            response.put("message", "Transaction signed and submitted successfully");
            response.put("note", "Transaction may fail validation if UTXOs don't exist. For testing, use the genesis address or mine a block first.");
            
            return ResponseEntity.ok(response);
        } catch (IllegalArgumentException e) {
            Map<String, Object> error = new HashMap<>();
            error.put("error", e.getMessage());
            return ResponseEntity.badRequest().body(error);
        } catch (Exception e) {
            Map<String, Object> error = new HashMap<>();
            error.put("error", "Transaction creation failed: " + e.getMessage());
            error.put("details", e.getClass().getSimpleName());
            if (e.getCause() != null) {
                error.put("cause", e.getCause().getMessage());
            }
            return ResponseEntity.internalServerError().body(error);
        }
    }

    /**
     * Get balance for an address.
     *
     * GET /api/wallet/balance/:address
     *
     * @param address Address to query
     * @return Balance
     */
    @GetMapping("/balance/{address}")
    public ResponseEntity<Map<String, Object>> getBalance(@PathVariable String address) {
        try {
            String balance = nodeClient.getBalance(address).block();
            Map<String, Object> response = new HashMap<>();
            response.put("address", address);
            response.put("balance", balance);
            return ResponseEntity.ok(response);
        } catch (Exception e) {
            Map<String, Object> error = new HashMap<>();
            error.put("error", e.getMessage());
            return ResponseEntity.internalServerError().body(error);
        }
    }
}
