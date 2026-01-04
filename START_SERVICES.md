# How to Start All Services

## Step 1: Start Go Node
Open a PowerShell window and run:
```powershell
cd ai-blockchain/go-node
go run cmd/node/main.go
```

Wait until you see: `Starting API server on :8080 (CORS enabled)`

## Step 2: Start Web UI Server
Open another PowerShell window and run:
```powershell
cd ai-blockchain/web-ui
python server.py
```

Wait until you see: `✅ Open http://localhost:3000 in your browser`

## Step 3: (Optional) Start AI Scorer
Open another PowerShell window and run:
```powershell
cd ai-blockchain/ai-scorer
python app/api.py
```

## Step 4: Open Browser
Go to: **http://localhost:3000**

## Troubleshooting

### "Not connected" message in UI
- Make sure Go node is running (Step 1)
- Check that port 8080 is not in use
- Look for errors in the Go node window

### "Failed to fetch" errors
- Make sure Go node is running
- Check browser console (F12) for CORS errors
- Make sure you're accessing via http://localhost:3000 (not file://)

### Wallet operations fail
- Make sure Go node is running
- Check that wallet routes are registered (should see wallet endpoints in Go node startup)

