package com.example.wallet.crypto;

import org.springframework.stereotype.Component;

import java.nio.charset.StandardCharsets;
import java.security.*;
import java.security.interfaces.ECPrivateKey;
import java.security.interfaces.ECPublicKey;
import java.security.spec.ECParameterSpec;
import java.util.Base64;

/**
 * ECDSA SIGNER – TRANSACTION SIGNING
 *
 * This class signs transactions using ECDSA.
 *
 * Process:
 * 1. Create canonical transaction bytes (must match Go node)
 * 2. Hash the bytes (SHA-256)
 * 3. Sign the hash with private key
 * 4. Encode signature as hex string
 *
 * Important:
 * - Signature format must match Go node's expectations
 * - Canonical serialization must be identical to Go node
 */
@Component
public class ECDSASigner {

    /**
     * Sign transaction data with private key.
     *
     * @param privateKey Private key for signing
     * @param data Transaction data to sign (canonical bytes)
     * @return Hex-encoded signature (r || s)
     * @throws Exception if signing fails
     */
    public String sign(PrivateKey privateKey, byte[] data) throws Exception {
        // Hash the data (ECDSA signs hashes, not raw data)
        byte[] hash = HashUtil.sha256Hex(data).getBytes(StandardCharsets.UTF_8);
        
        // Create signature
        Signature signature = Signature.getInstance("SHA256withECDSA");
        signature.initSign(privateKey);
        signature.update(hash);
        byte[] signatureBytes = signature.sign();
        
        // Encode as hex
        return bytesToHex(signatureBytes);
    }

    /**
     * Encode public key to hex string.
     *
     * Format: x || y (concatenated coordinates)
     * Must match Go node's public key encoding
     *
     * @param publicKey Public key to encode
     * @return Hex-encoded public key
     */
    public String encodePublicKey(PublicKey publicKey) {
        if (publicKey instanceof ECPublicKey) {
            ECPublicKey ecPublicKey = (ECPublicKey) publicKey;
            byte[] xBytes = ecPublicKey.getW().getAffineX().toByteArray();
            byte[] yBytes = ecPublicKey.getW().getAffineY().toByteArray();
            
            // Concatenate x and y
            byte[] combined = new byte[xBytes.length + yBytes.length];
            System.arraycopy(xBytes, 0, combined, 0, xBytes.length);
            System.arraycopy(yBytes, 0, combined, xBytes.length, yBytes.length);
            
            return bytesToHex(combined);
        }
        throw new IllegalArgumentException("Public key must be ECPublicKey");
    }

    /**
     * Convert bytes to hex string.
     */
    private String bytesToHex(byte[] bytes) {
        StringBuilder result = new StringBuilder();
        for (byte b : bytes) {
            result.append(String.format("%02x", b));
        }
        return result.toString();
    }
}

