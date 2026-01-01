# Troubleshooting Guide

## Transaction Creation Errors

### "400 Bad Request" when creating transaction

**Common Causes:**

1. **Wallet doesn't exist**
   - Error: "Wallet not found for address: ..."
   - Solution: Generate a wallet first using `/api/wallet/generate`

2. **No UTXOs available**
   - Error: "referenced UTXO not found"
   - Solution: You need to have coins to spend. Options:
     - Use the genesis address: `0000000000000000000000000000000000000000` (has 1000 coins)
     - Mine a block first to create new UTXOs
     - Receive coins from another transaction

3. **Invalid transaction format**
   - Error: "Invalid transaction: ..."
   - Solution: The transaction structure is invalid. Check that:
     - All required fields are present
     - Signature is valid
     - Transaction ID matches content

### How to Test Transactions

1. **Use Genesis Address (Easiest)**
   ```bash
   # The genesis address has 1000 coins
   From: 0000000000000000000000000000000000000000
   To: <your-wallet-address>
   Amount: 10.0
   ```

2. **Mine a Block First**
   ```bash
   # Mine a block to create UTXOs
   curl -X POST http://localhost:8080/mine
   ```

3. **Check Balance First**
   ```bash
   # Verify you have coins
   curl http://localhost:8080/balance/<your-address>
   ```

## CORS Errors

### "Failed to fetch" in browser console

**Cause**: Opening HTML file directly (`file://`) or CORS not enabled

**Solution**: Use a web server:
```bash
cd web-ui
python server.py
# Then open http://localhost:3000
```

## Service Connection Issues

### "Connection refused"

**Cause**: Service not running

**Solution**: Start all services:
1. Python AI Scorer (port 5000)
2. Go Node (port 8080)
3. Java Wallet (port 8081)

### "Cannot connect to Go node"

**Cause**: Go node not running or wrong URL

**Solution**: 
- Check Go node is running: `curl http://localhost:8080/health`
- Verify Java wallet config: `application.yml` has correct `go.node.url`

## Java Compilation Errors

### "javax.annotation does not exist"

**Cause**: Spring Boot 3.x uses Jakarta EE

**Solution**: Already fixed - use `jakarta.annotation.PostConstruct`

### "ECCurve cannot be converted to EllipticCurve"

**Cause**: BouncyCastle API mismatch

**Solution**: Already fixed - use `ECNamedCurveParameterSpec` directly

## Go Build Errors

### "import cycle not allowed"

**Cause**: Circular dependencies between packages

**Solution**: Already fixed - refactored to break cycles

## General Tips

1. **Check service logs** - They often show the exact error
2. **Verify all services are running** - Use `/health` endpoints
3. **Check ports are available** - No other services using 8080, 8081, 5000
4. **Restart services** - After code changes, restart all services

## Getting Help

1. Check service logs for detailed error messages
2. Verify all prerequisites are installed
3. Ensure all services are running
4. Check network connectivity between services

