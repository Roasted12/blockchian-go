package com.example.wallet.crypto;

import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;

/**
 * HASH UTILITIES – CRYPTOGRAPHIC HASHING
 *
 * Provides SHA-256 hashing functionality.
 * Must match Go node's hashing implementation.
 */
public class HashUtil {

    /**
     * Compute SHA-256 hash of bytes.
     *
     * @param data Input data
     * @return Hex-encoded hash
     */
    public static String sha256Hex(byte[] data) {
        try {
            MessageDigest digest = MessageDigest.getInstance("SHA-256");
            byte[] hash = digest.digest(data);
            return bytesToHex(hash);
        } catch (NoSuchAlgorithmException e) {
            throw new RuntimeException("SHA-256 algorithm not available", e);
        }
    }

    /**
     * Compute SHA-256 hash of string.
     *
     * @param data Input string
     * @return Hex-encoded hash
     */
    public static String sha256Hex(String data) {
        return sha256Hex(data.getBytes(StandardCharsets.UTF_8));
    }

    /**
     * Convert bytes to hex string.
     *
     * @param bytes Input bytes
     * @return Hex string
     */
    private static String bytesToHex(byte[] bytes) {
        StringBuilder result = new StringBuilder();
        for (byte b : bytes) {
            result.append(String.format("%02x", b));
        }
        return result.toString();
    }
}

