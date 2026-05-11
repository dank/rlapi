package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dank/rlapi"
)

// Device authorization flow for EGS. For durability save the refresh token in a persistent store, and restart with egs.AuthenticateWithRefreshToken.
// See examples/setup/setup.go for an example on persistent recovery.
func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	egs := rlapi.NewEGS()

	device, err := egs.AuthenticateWithDevice()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Visit %s and enter code: %s\n", device.VerificationURI, device.UserCode)

	token, err := egs.WaitForDeviceAuthorization(device)
	if err != nil {
		log.Fatal(err)
	}

	psyNet := rlapi.NewPsyNet()
	rpc, err := psyNet.AuthPlayer(token.AccessToken, token.AccountID, "")
	if err != nil {
		log.Fatal(err)
	}
	defer rpc.Close()

	// ... do stuff with rpc

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}
