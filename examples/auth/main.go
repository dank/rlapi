package main

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dank/rlapi"
)

// End-to-end authentication flow for EGS. For durability save the refresh token in a persistent store, and restart with egs.AuthenticateWithRefreshToken.
// See examples/setup/setup.go for an example on persistent recovery.
func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	egs := rlapi.NewEGS()

	// go to this link, authenticate with Epic, and paste in your auth code
	authURL := egs.GetAuthURL()
	fmt.Println(authURL)

	fmt.Print("auth code: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	authCode := strings.TrimSpace(scanner.Text())

	auth, err := egs.AuthenticateWithCode(authCode)
	if err != nil {
		log.Fatal(err)
	}

	code, err := egs.GetExchangeCode(auth.AccessToken)
	if err != nil {
		log.Fatal(err)
	}

	authToken, err := egs.ExchangeEOSToken(code)
	if err != nil {
		log.Fatal(err)
	}

	psyNet := rlapi.NewPsyNet()
	rpc, err := psyNet.AuthPlayer(authToken.AccessToken, auth.AccountID, auth.DisplayName)
	if err != nil {
		log.Fatal(err)
	}
	defer rpc.Close()

	// ... do stuff with rpc

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}
