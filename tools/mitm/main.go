package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	targetURL = "https://api.rlpp.psynet.gg"
	psyKey    = "3b932153785842ac927744b292e40e52"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func logHTTP(isResp bool, method, path, status string, headers http.Header, body []byte) {
	if isResp {
		log.Printf("<<< HTTP: %s", status)
	} else {
		log.Printf(">>> HTTP: %s %s", method, path)
	}
	if len(headers) > 0 {
		log.Printf("    Headers:")
		for k, v := range headers {
			log.Printf("      %s: %s", k, strings.Join(v, ", "))
		}
	}
	if len(body) > 0 {
		log.Printf("    Body: %s", string(body))
	}
}

func logWS(isResp bool, msg []byte) {
	if isResp {
		log.Printf("<<< WS: %s", string(msg))
	} else {
		log.Printf(">>> WS: %s", string(msg))
	}
}

func main() {
	// logFile, err := os.OpenFile("ws.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	// if err != nil {
	// 	log.Fatalf("Failed to open log file: %v", err)
	// }
	// defer logFile.Close()
	// log.SetOutput(logFile)

	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("Failed to load TLS certificate: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		ClientAuth:   tls.NoClientCert,
	}

	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("Failed to parse target URL: %v", err)
	}

	// proxyURL, _ := url.Parse("http://127.0.0.1:8888")
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.Host = target.Host

			var body []byte
			if req.Body != nil && (req.Method != http.MethodGet && req.Method != http.MethodHead) {
				body, _ = io.ReadAll(req.Body)
				req.Body = io.NopCloser(bytes.NewReader(body))
			}
			logHTTP(false, req.Method, req.URL.Path, "", req.Header, body)
		},
		Transport: &http.Transport{
			// Proxy: http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				MinVersion:         tls.VersionTLS12,
			},
		},
		ModifyResponse: func(resp *http.Response) error {
			bodyBytes, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			var data map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &data); err == nil {
				if result, ok := data["Result"].(map[string]interface{}); ok {
					// result["UseWebSocket"] = false
					result["PerConURLv2"] = "wss://127.0.0.1/ws"
				}
				bodyBytes, _ = json.Marshal(data)
			}

			psyTime := resp.Header.Get("PsyTime")
			h := hmac.New(sha256.New, []byte(psyKey))
			h.Write([]byte(psyTime + "-"))
			h.Write(bodyBytes)
			signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
			resp.Header.Set("PsySig", signature)

			for k, v := range resp.Header {
				if strings.ToLower(k) != "content-length" && strings.ToLower(k) != "psysig" {
					for _, vv := range v {
						resp.Header.Set(k, vv)
					}
				}
			}

			resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			resp.ContentLength = int64(len(bodyBytes))
			resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(bodyBytes)))

			logHTTP(true, "", "", resp.Status, resp.Header, bodyBytes)
			return nil
		},
	}

	mux := http.NewServeMux()
	mux.Handle("/", proxy)

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		clientConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WS upgrade error: %v", err)
			return
		}
		defer clientConn.Close()

		// proxyURL, _ := url.Parse("http://127.0.0.1:8888")
		dialer := websocket.Dialer{
			// Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		serverURL := url.URL{Scheme: "wss", Host: "ws.rlpp.psynet.gg", Path: "/ws/gc2"}
		headers := http.Header{
			"PsyBuildID":     []string{r.Header.Get("PsyBuildID")},
			"User-Agent":     []string{"RL Win/250811.43331.492665 gzip"},
			"PsyEnvironment": []string{"Prod"},
			"PsyToken":       []string{r.Header.Get("PsyToken")},
			"PsySessionID":   []string{r.Header.Get("PsySessionID")},
		}

		serverConn, _, err := dialer.Dial(serverURL.String(), headers)
		if err != nil {
			log.Printf("WS upstream dial error: %v", err)
			return
		}
		defer serverConn.Close()

		done := make(chan struct{})

		go func() {
			defer close(done)
			for {
				mt, msg, err := serverConn.ReadMessage()
				if err != nil {
					log.Printf("WS upstream read error: %v", err)
					return
				}

				logWS(true, msg)

				if err := clientConn.WriteMessage(mt, msg); err != nil {
					log.Printf("WS client write error: %v", err)
					return
				}
			}
		}()

		for {
			mt, msg, err := clientConn.ReadMessage()
			if err != nil {
				log.Printf("WS client read error: %v", err)
				break
			}

			logWS(false, msg)

			if err := serverConn.WriteMessage(mt, msg); err != nil {
				log.Printf("WS upstream write error: %v", err)
				break
			}
		}

		<-done
	})

	server := &http.Server{
		Addr:      ":443",
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Starting HTTPS/WS MITM server")
	if err := server.ListenAndServeTLS("server.crt", "server.key"); err != nil {
		log.Fatal(err)
	}
}
