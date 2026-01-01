package com.example.wallet.client;

import com.example.wallet.model.Transaction;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Mono;

/**
 * GO NODE CLIENT – COMMUNICATION WITH BLOCKCHAIN NODE
 *
 * This client communicates with the Go blockchain node via REST API.
 *
 * Endpoints used:
 * - POST /transactions - Submit new transaction
 * - GET  /balance/:addr - Get balance for address
 * - GET  /blocks - Get all blocks
 * - GET  /chain - Get blockchain info
 */
@Component
public class GoNodeClient {

    private final WebClient webClient;
    private final String baseUrl;

    public GoNodeClient(@Value("${go.node.url:http://localhost:8080}") String baseUrl) {
        this.baseUrl = baseUrl;
        this.webClient = WebClient.builder()
            .baseUrl(baseUrl)
            .build();
    }

    /**
     * Submit a transaction to the blockchain node.
     *
     * @param transaction Transaction to submit
     * @return Response from node
     */
    public Mono<String> submitTransaction(Transaction transaction) {
        return webClient.post()
            .uri("/transactions")
            .bodyValue(transaction)
            .retrieve()
            .bodyToMono(String.class);
    }

    /**
     * Get balance for an address.
     *
     * @param address Address to query
     * @return Balance as string
     */
    public Mono<String> getBalance(String address) {
        return webClient.get()
            .uri("/balance/" + address)
            .retrieve()
            .bodyToMono(String.class);
    }

    /**
     * Get all blocks.
     *
     * @return Blocks as JSON string
     */
    public Mono<String> getBlocks() {
        return webClient.get()
            .uri("/blocks")
            .retrieve()
            .bodyToMono(String.class);
    }

    /**
     * Get blockchain info.
     *
     * @return Chain info as JSON string
     */
    public Mono<String> getChainInfo() {
        return webClient.get()
            .uri("/chain")
            .retrieve()
            .bodyToMono(String.class);
    }

    /**
     * Get all blocks (to find UTXOs).
     *
     * @return Blocks as JSON string
     */
    public Mono<String> getBlocks() {
        return webClient.get()
            .uri("/blocks")
            .retrieve()
            .bodyToMono(String.class);
    }
}

