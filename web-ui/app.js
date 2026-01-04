/**
 * BLOCKCHAIN UI - JavaScript Application
 * 
 * Communicates with:
 * - Go Node (port 8080) - Blockchain operations
 * - Java Wallet (port 8081) - Wallet operations
 */

// Configuration
const GO_NODE_URL = 'http://localhost:8080';
const JAVA_WALLET_URL = 'http://localhost:8081';

// Tab management
function showTab(tabName) {
    // Hide all tabs
    document.querySelectorAll('.tab-content').forEach(tab => {
        tab.classList.remove('active');
    });
    
    // Remove active class from all tab buttons
    document.querySelectorAll('.tab').forEach(btn => {
        btn.classList.remove('active');
    });
    
    // Show selected tab
    document.getElementById(tabName).classList.add('active');
    
    // Activate corresponding button
    event.target.classList.add('active');
    
    // Auto-refresh when switching to certain tabs
    if (tabName === 'dashboard') {
        refreshDashboard();
    } else if (tabName === 'blocks') {
        loadBlocks();
    } else if (tabName === 'transactions') {
        loadMempool();
    }
}

// Status messages
function showStatus(message, type = 'info') {
    const statusEl = document.getElementById('status');
    statusEl.textContent = message;
    statusEl.className = `status-message show ${type}`;
    
    setTimeout(() => {
        statusEl.classList.remove('show');
    }, 3000);
}

// Dashboard functions
async function refreshDashboard() {
    try {
        const response = await fetch(`${GO_NODE_URL}/chain`);
        const data = await response.json();
        
        document.getElementById('chain-height').textContent = data.height || '0';
        document.getElementById('difficulty').textContent = data.difficulty || '4';
        
        if (data.tip) {
            const hash = data.tip.hash || 'N/A';
            document.getElementById('last-hash').textContent = 
                hash.length > 20 ? hash.substring(0, 20) + '...' : hash;
        }
        
        // Get mempool size
        const mempoolRes = await fetch(`${GO_NODE_URL}/mempool`);
        const mempoolData = await mempoolRes.json();
        document.getElementById('mempool-size').textContent = mempoolData.count || '0';
        
        showStatus('Dashboard refreshed', 'success');
    } catch (error) {
        showStatus('Error refreshing dashboard: ' + error.message, 'error');
        console.error('Dashboard error:', error);
    }
}

// Blocks functions
async function loadBlocks() {
    try {
        const response = await fetch(`${GO_NODE_URL}/blocks`);
        const data = await response.json();
        
        const blocksList = document.getElementById('blocks-list');
        blocksList.innerHTML = '';
        
        if (!data.blocks || data.blocks.length === 0) {
            blocksList.innerHTML = '<p>No blocks found. Mine the first block!</p>';
            return;
        }
        
        // Reverse to show newest first
        data.blocks.reverse().forEach(block => {
            const blockCard = createBlockCard(block);
            blocksList.appendChild(blockCard);
        });
        
        showStatus(`Loaded ${data.blocks.length} blocks`, 'success');
    } catch (error) {
        showStatus('Error loading blocks: ' + error.message, 'error');
        console.error('Blocks error:', error);
    }
}

function createBlockCard(block) {
    const div = document.createElement('div');
    div.className = 'block-card';
    
    const date = new Date(block.timestamp * 1000).toLocaleString();
    const txCount = block.transactions ? block.transactions.length : 0;
    
    div.innerHTML = `
        <div class="block-header">
            <div class="block-index">Block #${block.index}</div>
            <div class="block-hash">${block.hash.substring(0, 20)}...</div>
        </div>
        <div class="block-info">
            <div class="block-info-item">
                <span class="block-info-label">Hash:</span>
                <span style="font-family: monospace; font-size: 0.8em;">${block.hash}</span>
            </div>
            <div class="block-info-item">
                <span class="block-info-label">Previous:</span>
                <span style="font-family: monospace; font-size: 0.8em;">${block.prevHash.substring(0, 20)}...</span>
            </div>
            <div class="block-info-item">
                <span class="block-info-label">Transactions:</span>
                ${txCount}
            </div>
            <div class="block-info-item">
                <span class="block-info-label">Timestamp:</span>
                ${date}
            </div>
            <div class="block-info-item">
                <span class="block-info-label">Nonce:</span>
                ${block.nonce}
            </div>
        </div>
    `;
    
    return div;
}

// Mempool functions
async function loadMempool() {
    try {
        const response = await fetch(`${GO_NODE_URL}/mempool`);
        const data = await response.json();
        
        const mempoolList = document.getElementById('mempool-list');
        mempoolList.innerHTML = '';
        
        if (!data.transactions || data.transactions.length === 0) {
            mempoolList.innerHTML = '<p>No pending transactions in mempool.</p>';
            return;
        }
        
        data.transactions.forEach(tx => {
            const txCard = createTransactionCard(tx);
            mempoolList.appendChild(txCard);
        });
        
        showStatus(`Loaded ${data.transactions.length} pending transactions`, 'success');
    } catch (error) {
        showStatus('Error loading mempool: ' + error.message, 'error');
        console.error('Mempool error:', error);
    }
}

function createTransactionCard(tx) {
    const div = document.createElement('div');
    div.className = 'transaction-card';
    
    const inputCount = tx.inputs ? tx.inputs.length : 0;
    const outputCount = tx.outputs ? tx.outputs.length : 0;
    const totalOutput = tx.outputs ? tx.outputs.reduce((sum, out) => sum + (out.amount || 0), 0) : 0;
    
    div.innerHTML = `
        <div class="tx-id">TX ID: ${tx.id}</div>
        <div class="tx-details">
            <div><strong>Inputs:</strong> ${inputCount}</div>
            <div><strong>Outputs:</strong> ${outputCount}</div>
            <div><strong>Total Amount:</strong> ${totalOutput.toFixed(2)}</div>
            <div><strong>Timestamp:</strong> ${new Date(tx.timestamp * 1000).toLocaleString()}</div>
        </div>
    `;
    
    return div;
}

// Wallet storage (in browser)
let wallets = [];
let walletBalances = {};

// Wallet functions
async function generateAddress() {
    try {
        const response = await fetch(`${JAVA_WALLET_URL}/api/wallet/generate`);
        const data = await response.json();
        
        if (data.error) {
            showStatus('Error: ' + data.error, 'error');
        } else {
            showStatus('✅ Wallet generated successfully!', 'success');
            // Reload wallets list
            await loadWallets();
        }
    } catch (error) {
        showStatus('Error generating address: ' + error.message, 'error');
        console.error('Generate address error:', error);
    }
}

async function loadWallets() {
    try {
        const response = await fetch(`${JAVA_WALLET_URL}/api/wallet/list`);
        const data = await response.json();
        
        wallets = data.addresses || [];
        const walletsList = document.getElementById('wallets-list');
        const sendFromSelect = document.getElementById('send-from');
        const receiveSelect = document.getElementById('receive-address');
        
        // Clear existing options (keep first option)
        sendFromSelect.innerHTML = '<option value="">Select a wallet...</option>';
        receiveSelect.innerHTML = '<option value="">Select a wallet to receive...</option>';
        
        if (wallets.length === 0) {
            walletsList.innerHTML = '<p style="color: #666; padding: 20px; text-align: center;">No wallets yet. Generate one to get started!</p>';
            return;
        }
        
        // Display wallets
        walletsList.innerHTML = '';
        wallets.forEach(async (address, index) => {
            // Create wallet card
            const walletCard = document.createElement('div');
            walletCard.className = 'wallet-card';
            walletCard.style.cssText = 'background: white; border: 2px solid #e0e0e0; border-radius: 8px; padding: 15px; margin-bottom: 10px;';
            
            // Get balance
            let balance = 'Loading...';
            try {
                const balanceRes = await fetch(`${GO_NODE_URL}/balance/${address}`);
                const balanceData = await balanceRes.json();
                balance = balanceData.balance || '0';
                walletBalances[address] = balance;
            } catch {
                balance = '0';
                walletBalances[address] = '0';
            }
            
            walletCard.innerHTML = `
                <div style="display: flex; justify-content: space-between; align-items: start;">
                    <div style="flex: 1;">
                        <div style="font-weight: bold; color: #667eea; margin-bottom: 5px;">Wallet #${index + 1}</div>
                        <div style="font-family: monospace; font-size: 0.85em; color: #666; word-break: break-all; margin-bottom: 10px;">${address}</div>
                        <div style="font-size: 1.2em; font-weight: bold; color: #4caf50;">💰 ${balance} coins</div>
                    </div>
                    <div>
                        <button onclick="copyAddress('${address}')" class="btn-primary" style="padding: 5px 10px; font-size: 0.9em; margin: 2px;">📋 Copy</button>
                    </div>
                </div>
            `;
            
            walletsList.appendChild(walletCard);
            
            // Add to selects
            const sendOption = document.createElement('option');
            sendOption.value = address;
            sendOption.textContent = `Wallet #${index + 1} (${balance} coins)`;
            sendFromSelect.appendChild(sendOption);
            
            const receiveOption = document.createElement('option');
            receiveOption.value = address;
            receiveOption.textContent = `Wallet #${index + 1} (${balance} coins)`;
            receiveSelect.appendChild(receiveOption);
        });
        
        showStatus(`Loaded ${wallets.length} wallet(s)`, 'success');
    } catch (error) {
        showStatus('Error loading wallets: ' + error.message, 'error');
        console.error('Load wallets error:', error);
    }
}

function copyAddress(address) {
    navigator.clipboard.writeText(address).then(() => {
        showStatus('Address copied to clipboard!', 'success');
    }).catch(() => {
        showStatus('Failed to copy address', 'error');
    });
}

function updateSendFromBalance() {
    const address = document.getElementById('send-from').value;
    const balanceDiv = document.getElementById('send-from-balance');
    
    if (address) {
        const balance = walletBalances[address] || '0';
        balanceDiv.innerHTML = `<strong>Available:</strong> ${balance} coins`;
    } else {
        balanceDiv.innerHTML = '';
    }
}

function updateReceiveInfo() {
    const address = document.getElementById('receive-address').value;
    const receiveInfo = document.getElementById('receive-info');
    const balanceDiv = document.getElementById('receive-balance');
    
    if (address) {
        const balance = walletBalances[address] || '0';
        receiveInfo.innerHTML = `
            <div style="margin-bottom: 10px;">
                <strong>Your Receive Address:</strong>
            </div>
            <div style="font-family: monospace; font-size: 0.9em; word-break: break-all; background: #f5f5f5; padding: 10px; border-radius: 4px; margin-bottom: 10px;">
                ${address}
            </div>
            <button onclick="copyAddress('${address}')" class="btn-primary" style="width: 100%;">📋 Copy Address</button>
            <div style="margin-top: 10px; padding: 10px; background: #e3f2fd; border-radius: 4px; font-size: 0.9em;">
                💡 Share this address to receive coins from others
            </div>
        `;
        balanceDiv.style.display = 'block';
        document.getElementById('receive-balance-value').textContent = balance;
    } else {
        receiveInfo.innerHTML = '<p style="color: #666;">Select a wallet above to see your receive address</p>';
        balanceDiv.style.display = 'none';
    }
}

function fillRecipientFromWallets() {
    const receiveSelect = document.getElementById('receive-address');
    const sendToInput = document.getElementById('send-to');
    
    if (receiveSelect.value) {
        sendToInput.value = receiveSelect.value;
        showStatus('Recipient address filled from wallet list', 'success');
    } else {
        showStatus('Please select a wallet in the Receive section first', 'error');
    }
}

async function sendTransaction() {
    const from = document.getElementById('send-from').value;
    const to = document.getElementById('send-to').value.trim();
    const amount = parseFloat(document.getElementById('send-amount').value);
    
    if (!from) {
        showStatus('Please select a wallet to send from', 'error');
        return;
    }
    
    if (!to) {
        showStatus('Please enter a recipient address', 'error');
        return;
    }
    
    if (!amount || amount <= 0) {
        showStatus('Please enter a valid amount', 'error');
        return;
    }
    
    const resultBox = document.getElementById('send-result');
    resultBox.innerHTML = '<div>Creating transaction... <span class="loading"></span></div>';
    
    try {
        const txData = {
            from: from,
            to: to,
            amount: amount
        };
        
        const response = await fetch(`${JAVA_WALLET_URL}/api/wallet/transfer`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(txData)
        });
        
        if (!response.ok) {
            const errorText = await response.text();
            let errorData;
            try {
                errorData = JSON.parse(errorText);
            } catch {
                errorData = { error: errorText || `HTTP ${response.status}` };
            }
            
            resultBox.innerHTML = `<div style="color: red;"><strong>❌ Error:</strong> ${errorData.error || 'Transaction failed'}</div>`;
            if (errorData.hint) {
                resultBox.innerHTML += `<div style="margin-top: 10px; color: #666;">💡 ${errorData.hint}</div>`;
            }
            showStatus('Transaction failed', 'error');
            return;
        }
        
        const data = await response.json();
        
        if (data.error) {
            resultBox.innerHTML = `<div style="color: red;"><strong>❌ Error:</strong> ${data.error}</div>`;
            showStatus('Transaction failed', 'error');
        } else {
            resultBox.innerHTML = `
                <div style="color: green;"><strong>✅ Transaction Submitted!</strong></div>
                <div style="margin-top: 10px;"><strong>TX ID:</strong> <code style="font-size: 0.9em;">${data.txid || 'N/A'}</code></div>
                <div style="margin-top: 10px; padding: 10px; background: #e8f5e9; border-radius: 4px; font-size: 0.9em;">
                    ${data.message || 'Transaction is being processed'}
                </div>
            `;
            showStatus('✅ Transaction sent successfully!', 'success');
            
            // Clear form
            document.getElementById('send-amount').value = '';
            document.getElementById('send-to').value = '';
            
            // Refresh wallets and mempool
            setTimeout(() => {
                loadWallets();
                loadMempool();
            }, 1000);
        }
    } catch (error) {
        resultBox.innerHTML = `<div style="color: red;"><strong>❌ Error:</strong> ${error.message}</div>`;
        showStatus('Error: ' + error.message, 'error');
        console.error('Send transaction error:', error);
    }
}

async function checkAllBalances() {
    showStatus('Checking all balances...', 'info');
    await loadWallets();
    showStatus('All balances updated!', 'success');
}

function copyGenesisAddress() {
    const genesisAddress = '0000000000000000000000000000000000000000';
    navigator.clipboard.writeText(genesisAddress).then(() => {
        showStatus('Genesis address copied! (Has 1000 coins)', 'success');
    });
}

// Legacy functions (kept for compatibility, but wallet tab uses new functions)
async function checkBalance() {
    const address = document.getElementById('balance-address')?.value.trim();
    
    if (!address) {
        showStatus('Please enter an address', 'error');
        return;
    }
    
    try {
        const response = await fetch(`${GO_NODE_URL}/balance/${address}`);
        const data = await response.json();
        
        const resultBox = document.getElementById('balance-result');
        if (resultBox) {
            if (data.error) {
                resultBox.innerHTML = `<div style="color: red;">Error: ${data.error}</div>`;
            } else {
                resultBox.innerHTML = `
                    <div><strong>Address:</strong> ${data.address}</div>
                    <div><strong>Balance:</strong> ${data.balance} coins</div>
                `;
                showStatus('Balance retrieved', 'success');
            }
        }
    } catch (error) {
        showStatus('Error checking balance: ' + error.message, 'error');
        console.error('Balance error:', error);
    }
}

// Mining functions
async function mineBlock() {
    const resultBox = document.getElementById('mine-result');
    resultBox.innerHTML = '<div>Mining block... <span class="loading"></span></div>';
    
    try {
        const response = await fetch(`${GO_NODE_URL}/mine`, {
            method: 'POST'
        });
        
        const data = await response.json();
        
        if (data.error) {
            resultBox.innerHTML = `<div style="color: red;">Error: ${data.error}</div>`;
            showStatus('Mining failed: ' + data.error, 'error');
        } else {
            resultBox.innerHTML = `
                <div style="color: green;"><strong>Block mined successfully!</strong></div>
                <div><strong>Block Index:</strong> ${data.block.index}</div>
                <div><strong>Hash:</strong> ${data.block.hash.substring(0, 40)}...</div>
                <div><strong>Transactions:</strong> ${data.block.transactions.length}</div>
                <div><strong>Mining Time:</strong> ${data.time}</div>
            `;
            showStatus('Block mined successfully!', 'success');
            
            // Refresh dashboard and blocks
            setTimeout(() => {
                refreshDashboard();
                loadBlocks();
                loadMempool();
            }, 1000);
        }
    } catch (error) {
        resultBox.innerHTML = `<div style="color: red;">Error: ${error.message}</div>`;
        showStatus('Mining error: ' + error.message, 'error');
        console.error('Mining error:', error);
    }
}

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    // Load dashboard on startup
    refreshDashboard();
    
    // Load wallets on startup
    loadWallets();
    
    // Set up tab click handlers
    document.querySelectorAll('.tab').forEach(tab => {
        tab.addEventListener('click', (e) => {
            const tabName = e.target.textContent.toLowerCase();
            const tabMap = {
                'dashboard': 'dashboard',
                'blocks': 'blocks',
                'transactions': 'transactions',
                'wallet': 'wallet',
                'mine': 'mine'
            };
            const selectedTab = tabMap[tabName] || 'dashboard';
            showTab(selectedTab);
            
            // Auto-load wallets when wallet tab is opened
            if (selectedTab === 'wallet') {
                loadWallets();
            }
        });
    });
});

