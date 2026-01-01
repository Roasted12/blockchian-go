#!/usr/bin/env python3
"""
Simple HTTP server for the blockchain web UI.

This serves the HTML files and enables CORS, allowing the UI
to communicate with the Go node and Java wallet APIs.

Usage:
    python server.py

Then open http://localhost:3000 in your browser.
"""

import http.server
import socketserver
import os
from pathlib import Path

PORT = 3000

class CORSRequestHandler(http.server.SimpleHTTPRequestHandler):
    def end_headers(self):
        # Add CORS headers
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type')
        super().end_headers()
    
    def do_OPTIONS(self):
        self.send_response(200)
        self.end_headers()

def main():
    # Change to the directory containing this script
    os.chdir(Path(__file__).parent)
    
    with socketserver.TCPServer(("", PORT), CORSRequestHandler) as httpd:
        print(f"🚀 Blockchain Web UI Server")
        print(f"📡 Serving on http://localhost:{PORT}")
        print(f"📁 Serving files from: {os.getcwd()}")
        print(f"\n✅ Open http://localhost:{PORT} in your browser")
        print(f"⏹️  Press Ctrl+C to stop\n")
        
        try:
            httpd.serve_forever()
        except KeyboardInterrupt:
            print("\n\n👋 Server stopped")

if __name__ == "__main__":
    main()

