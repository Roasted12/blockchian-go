package com.example.wallet;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

/**
 * WALLET APPLICATION – ENTRY POINT
 *
 * This is a Spring Boot application that provides:
 * - Wallet functionality (key generation, transaction signing)
 * - Block explorer (query blocks, transactions, balances)
 * - Integration with Go blockchain node
 *
 * The service communicates with the Go node via REST API.
 */
@SpringBootApplication
public class WalletApplication {
    public static void main(String[] args) {
        SpringApplication.run(WalletApplication.class, args);
    }
}

