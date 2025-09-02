# RL MITM

A local man-in-the-middle (MITM) proxy for intercepting Rocket League traffic. It forwards requests to the official servers while logging and modifying traffic as needed. Even with certificate pinning disabled, the game requires HTTPS and WSS connections; this tool handles them using self-signed certificates.

**Features:**
- Intercepts authentication requests to rewrite responses and redirect WebSocket connections through a local proxy.
- Logs all requests and responses and forwards them while acting as a seamless proxy.
- Properly re-signs HMAC-SHA256 signatures for modified responses.
- Optionally routes traffic through a Fiddler proxy (http://127.0.0.1:8888) to intercept and inspect it in a familiar interface.

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
