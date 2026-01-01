package com.example.wallet.model;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * TRANSACTION OUTPUT MODEL
 *
 * Represents a new UTXO being created.
 */
public class TxOut {
    @JsonProperty("address")
    private String address;

    @JsonProperty("amount")
    private double amount;

    // Getters and setters
    public String getAddress() { return address; }
    public void setAddress(String address) { this.address = address; }

    public double getAmount() { return amount; }
    public void setAmount(double amount) { this.amount = amount; }
}

