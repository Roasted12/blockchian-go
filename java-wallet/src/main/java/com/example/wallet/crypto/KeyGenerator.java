package com.example.wallet.crypto;

import org.bouncycastle.jce.ECNamedCurveTable;
import org.bouncycastle.jce.spec.ECNamedCurveParameterSpec;
import org.springframework.stereotype.Component;

import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.SecureRandom;

/**
 * KEY GENERATOR – ECDSA KEY PAIR CREATION
 *
 * This class generates ECDSA key pairs for wallet addresses.
 *
 * Algorithm: ECDSA with secp256r1 (P-256) curve
 * - Same curve as Go node uses
 * - Ensures cross-language compatibility
 *
 * Process:
 * 1. Generate private/public key pair
 * 2. Derive address from public key (hash of public key)
 * 3. Return key pair and address
 */
@Component
public class KeyGenerator {

    /**
     * Generate a new ECDSA key pair.
     *
     * @return KeyPair containing private and public keys
     * @throws Exception if key generation fails
     */
    public KeyPair generateKeyPair() throws Exception {
        // Use secp256r1 (P-256) curve - same as Go node
        // Get curve parameters from BouncyCastle
        ECNamedCurveParameterSpec curveParams = ECNamedCurveTable.getParameterSpec("secp256r1");
        
        // Create key pair generator using BouncyCastle provider
        KeyPairGenerator keyGen = KeyPairGenerator.getInstance("EC", "BC");
        
        // Initialize with BouncyCastle's curve parameter spec
        // This works because BouncyCastle's KeyPairGenerator accepts ECNamedCurveParameterSpec
        keyGen.initialize(curveParams, new SecureRandom());
        
        // Generate key pair
        return keyGen.generateKeyPair();
    }

    /**
     * Derive address from public key.
     *
     * Address = SHA256(public key) (simplified)
     * In production, you might use RIPEMD160(SHA256(pubkey))
     *
     * @param publicKeyBytes Public key bytes
     * @return Address (hex-encoded hash)
     */
    public String deriveAddress(byte[] publicKeyBytes) {
        // Hash public key to get address
        // This is a simplified version - in production, use proper address encoding
        return HashUtil.sha256Hex(publicKeyBytes);
    }
}

