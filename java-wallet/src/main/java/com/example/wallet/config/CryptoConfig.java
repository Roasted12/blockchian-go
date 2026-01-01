package com.example.wallet.config;

import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.springframework.context.annotation.Configuration;

import jakarta.annotation.PostConstruct;
import java.security.Security;

/**
 * CRYPTO CONFIGURATION – BOUNCYCASTLE PROVIDER
 *
 * Registers BouncyCastle as a security provider.
 * Required for ECDSA operations.
 */
@Configuration
public class CryptoConfig {

    @PostConstruct
    public void init() {
        // Register BouncyCastle provider
        if (Security.getProvider("BC") == null) {
            Security.addProvider(new BouncyCastleProvider());
        }
    }
}

