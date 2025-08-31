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

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	egs := rlapi.NewEGS()
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

	//refreshToken, err := egs.RefreshEOSToken(authToken.RefreshToken)

	psyNet := rlapi.NewPsyNet()
	defer psyNet.Close()
	err = psyNet.AuthPlayer(rlapi.Epic, authToken.AccessToken, auth.AccountID, auth.DisplayName)
	if err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}
