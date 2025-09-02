# RL Frida

Frida hooks for bypassing certificate pinning and redirecting Rocket League API traffic to a local MITM server.

## How it works
### 1. Hook libcurl [`curl_easy_init`](https://curl.se/libcurl/c/curl_easy_init.html)
Disable SSL verification by overriding the following options:
- `CURLOPT_SSL_VERIFYPEER` - Disable peer certificate verification
- `CURLOPT_SSL_VERIFYHOST` - Disable hostname verification

This allows the game to accept a self-signed certificate from the local MITM server without errors.

### 2. Hook libcurl [`curl_easy_setopt`](https://curl.se/libcurl/c/curl_easy_setopt.html)
Intercept URL calls to redirect API traffic:
- Detect when `CURLOPT_URL` is set to `https://api.rlpp.psynet.gg`
- Replace the URL with `https://127.0.0.1` while preserving the path
- Force all HTTP API calls to route through the local MITM proxy

### 3. Hook OpenSSL [`X509_verify_cert`](https://docs.openssl.org/1.1.1/man3/X509_verify_cert/)
Bypass certificate validation for WebSocket connections:
- Replace the certificate verification function to always return success (1)
- Essential for WebSocket connections which use a separate validation path
- Prevent SSL handshake failures when connecting to the local proxy

## Prerequisites
- Python 3.x
- Frida (`pip install frida`)
- Node.js

## Usage

> [!NOTE]
> Steam and Epic builds have different offsets for each function. You may need to adjust them, check the code for details.

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