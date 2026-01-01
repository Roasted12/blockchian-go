package com.example.wallet.model;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * TRANSACTION INPUT MODEL
 *
 * Represents a reference to a UTXO being spent.
 */
public class TxIn {
    @JsonProperty("tx_id")
    private String txId;

    @JsonProperty("index")
    private int index;

    // Getters and setters
    public String getTxId() { return txId; }
    public void setTxId(String txId) { this.txId = txId; }

    public int getIndex() { return index; }
    public void setIndex(int index) { this.index = index; }
}

