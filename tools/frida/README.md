# RL Frida

Frida dynamic instrumentation to bypass certificate pinning and redirect Rocket League API traffic to a local MITM server.

## Overview
### Hook `curl_easy_init`
Disables SSL verification by setting:
- `CURLOPT_SSL_VERIFYPEER=0` - Disables peer certificate verification
- `CURLOPT_SSL_VERIFYHOST=0` - Disables hostname verification

This allows the game to accept the self-signed certificate from the local MITM server without validation errors.

### Hook `curl_easy_setopt`
Intercepts URL setting calls and redirects API traffic:
- Detects when `CURLOPT_URL` is set to `https://api.rlpp.psynet.gg`
- Replaces the URL with `https://127.0.0.1` while preserving the path
- Forces all HTTP API calls to route through the local MITM proxy

### Hook `X509_verify_cert`
Bypasses certificate validation for WebSocket connections:
- Replaces the certificate verification function to always return success (1)
- Essential for WebSocket connections which use a separate validation path
- Prevents SSL handshake failures when connecting to the local proxy

## Prerequisites
- Python 3.x
- Frida (`pip install frida`)
- Node.js

## Usage
1. Build the TypeScript Frida script:
```bash
npm install
npm run build
```

2. Start the MITM server (see `../mitm/` directory)

3. Launch Rocket League

4. **Immediately** run the Frida hook:
```bash
npm start
# or directly:
python main.py
```