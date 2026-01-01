package com.example.wallet.model;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.List;

/**
 * TRANSACTION MODEL
 *
 * Represents a blockchain transaction.
 * Must match Go node's transaction structure.
 */
public class Transaction {
    @JsonProperty("id")
    private String id;

    @JsonProperty("inputs")
    private List<TxIn> inputs;

    @JsonProperty("outputs")
    private List<TxOut> outputs;

    @JsonProperty("signature")
    private String signature;

    @JsonProperty("pubkey")
    private String pubkey;

    @JsonProperty("timestamp")
    private long timestamp;

    // Getters and setters
    public String getId() { return id; }
    public void setId(String id) { this.id = id; }

    public List<TxIn> getInputs() { return inputs; }
    public void setInputs(List<TxIn> inputs) { this.inputs = inputs; }

    public List<TxOut> getOutputs() { return outputs; }
    public void setOutputs(List<TxOut> outputs) { this.outputs = outputs; }

    public String getSignature() { return signature; }
    public void setSignature(String signature) { this.signature = signature; }

    public String getPubkey() { return pubkey; }
    public void setPubkey(String pubkey) { this.pubkey = pubkey; }

    public long getTimestamp() { return timestamp; }
    public void setTimestamp(long timestamp) { this.timestamp = timestamp; }
}

