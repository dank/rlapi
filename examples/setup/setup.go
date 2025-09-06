package setup

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/dank/rlapi"
)

// RPC return an authenticated PsyNetRPC client and a PlayerID for easy setup in examples.
func RPC() (*rlapi.PsyNetRPC, rlapi.PlayerID) {
	egs := rlapi.NewEGS()
	var auth *rlapi.TokenResponse
	var authToken *rlapi.EOSTokenResponse

	if refreshTokenData, err := os.ReadFile(".rlshops"); err == nil && len(strings.TrimSpace(string(refreshTokenData))) > 0 {
		refreshToken := strings.TrimSpace(string(refreshTokenData))
		auth, err = egs.AuthenticateWithRefreshToken(refreshToken)
		if err != nil {
			slog.Error("Failed to authenticate with refresh token", slog.Any("err", err))
			auth = authenticateWithCode(egs)
		}
	} else {
		auth = authenticateWithCode(egs)
	}

	err := os.WriteFile(".rlshops", []byte(auth.RefreshToken), 0644)
	if err != nil {
		slog.Error("Failed to save refresh token", slog.Any("err", err))
	}

	code, err := egs.GetExchangeCode(auth.AccessToken)
	if err != nil {
		log.Fatalf("Failed to get exchange code: %v", err)
	}

	authToken, err = egs.ExchangeEOSToken(code)
	if err != nil {
		log.Fatalf("Failed to exchange EOS token: %v", err)
	}

	psyNet := rlapi.NewPsyNet()
	rpc, err := psyNet.AuthPlayer(rlapi.PlatformEpic, authToken.AccessToken, authToken.AccountID, auth.DisplayName)
	if err != nil {
		log.Fatalf("Failed to authenticate player: %v", err)
	}

	return rpc, rlapi.NewPlayerID(rlapi.PlatformEpic, authToken.AccountID)
}

func authenticateWithCode(egs *rlapi.EGS) *rlapi.TokenResponse {
	authURL := egs.GetAuthURL()
	fmt.Println(authURL)

	fmt.Print("Auth code: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	authCode := strings.TrimSpace(scanner.Text())

	auth, err := egs.AuthenticateWithCode(authCode)
	if err != nil {
		log.Fatalf("Failed to authenticate with code: %v", err)
	}

	return auth
}
