const http = require('http');
const WebSocket = require('ws');
// Initialize base system constants
const PORT = 8080;
const server = http.createServer((req, res) => {
    res.writeHead(200, { 'Content-Type': 'text/plain' });
    res.end('UNBOUND NET STREAMING NODE PROXY');
});
// Attach WebSocket layer directly onto the server instance
const wss = new WebSocket.Server({ server });
// Track client registration arrays
wss.on('connection', (ws) => {
    console.log('[NODE PROXY] Transmit channel established with client.');
    
    // Create an isolated sub-interval loop for this connection channel
    let cycleCount = 0;
    const streamInterval = setInterval(() => {
        cycleCount++;
        const isAnomaly = Math.random() < 0.08;
        
        const packet = {
            timestamp: new Date().toISOString(),
            heartbeat: cycleCount,
            anomaly: isAnomaly,
            resonanceHz: isAnomaly ? (32.0 + (Math.random() * 33.0)) : (32.0 + (Math.random() * 4.0 - 2.0))
        };
        
        // Transmit packet down the socket wire in string format
        if (ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify(packet));
        }
    }, 120); // 120ms tick matching the terminal refresh matrix rate

    ws.on('close', () => {
        console.log('[NODE PROXY] Client disconnected. Clearing channel loop.');
        clearInterval(streamInterval);
    });
});

server.listen(PORT, () => {
    console.log(`[INIT] Streaming proxy active on ws://localhost:${PORT}`);
});
