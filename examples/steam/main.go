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

// Steam authentication example using auth session ticket.
// First run the Node.js script to generate a Steam auth ticket, then paste it here.
// The Node.js script uses steam-user package to authenticate and create the session ticket.
func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	// !!! Replace with your actual Steam ID !!!
	steamID64 := ""
	if steamID64 == "" {
		log.Fatal("Steam ID must be set in code")
	}

	egs := rlapi.NewEGS()

	fmt.Println("Run the Node.js steam-user script to generate an auth session ticket, then paste the hex ticket here:")
	fmt.Print("Steam auth ticket: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	ticket := strings.TrimSpace(scanner.Text())

	if ticket == "" {
		log.Fatal("No ticket provided")
	}

	authToken, err := egs.ExchangeEOSTokenFromSteam(ticket)
	if err != nil {
		log.Fatalf("Failed to exchange Steam ticket: %v", err)
	}

	psyNet := rlapi.NewPsyNet()
	rpc, err := psyNet.AuthPlayerSteam(authToken.AccessToken, authToken.AccountID, steamID64, "")
	if err != nil {
		log.Fatalf("Failed to authenticate player: %v", err)
	}
	defer rpc.Close()

	slog.Info("Connected to PsyNet RPC")
	// ... do stuff with rpc

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	slog.Info("Shutting down")
}
