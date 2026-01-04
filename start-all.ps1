# Start All Blockchain Services
Write-Host "🚀 Starting Blockchain Services..." -ForegroundColor Green
Write-Host ""

# Start Go Node
Write-Host "1. Starting Go Node (port 8080)..." -ForegroundColor Yellow
Start-Process powershell -ArgumentList '-NoExit', '-Command', "cd '$PSScriptRoot\go-node'; Write-Host 'Starting Go Node...' -ForegroundColor Green; go run cmd/node/main.go"
Start-Sleep -Seconds 3

# Start Web UI
Write-Host "2. Starting Web UI Server (port 3000)..." -ForegroundColor Yellow
Start-Process powershell -ArgumentList '-NoExit', '-Command', "cd '$PSScriptRoot\web-ui'; Write-Host 'Starting Web UI Server...' -ForegroundColor Green; python server.py"
Start-Sleep -Seconds 2

# Optional: Start AI Scorer
Write-Host "3. (Optional) Starting AI Scorer (port 5000)..." -ForegroundColor Yellow
Write-Host "   Press Y to start AI scorer, or any other key to skip: " -NoNewline
$response = Read-Host
if ($response -eq 'Y' -or $response -eq 'y') {
    Start-Process powershell -ArgumentList '-NoExit', '-Command', "cd '$PSScriptRoot\ai-scorer'; Write-Host 'Starting AI Scorer...' -ForegroundColor Green; python app/api.py"
}

Write-Host ""
Write-Host "✅ All services started!" -ForegroundColor Green
Write-Host ""
Write-Host "📱 Open your browser and go to: http://localhost:3000" -ForegroundColor Cyan
Write-Host ""
Write-Host "Note: Check the PowerShell windows for any errors." -ForegroundColor Yellow
Write-Host "      Wait a few seconds for services to fully start." -ForegroundColor Yellow
