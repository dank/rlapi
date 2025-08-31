# RL MITM

HTTPS/WebSocket man-in-the-middle proxy server for intercepting Rocket League API traffic. Acts as a local proxy that forwards requests to the actual Rocket League servers while logging and modifying traffic as needed.

## Overview

This tool provides:
- HTTPS and WSS proxy server with self-signed certificates
- Authentication response modification to redirect WebSocket connections to local proxy
- Proper HMAC-SHA256 signature verification and re-signing

## Usage

1. Generate self-signed certificates (if not already present):
```bash
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr
openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt
```

2. Build and run the proxy:
```bash
go mod tidy
go run main.go
```

The server will start on port 443.
