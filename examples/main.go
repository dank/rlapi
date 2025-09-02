package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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

	//auth, err := egs.AuthenticateWithRefreshToken("eg1~eyJraWQiOiJnX19WS2pTU21xSjB4WmoxUllrTEdLUTdkbkhpTTlNTGhGVndLUHlTREI0IiwiYWxnIjoiUFMyNTYifQ.eyJzdWIiOiI4YTAyMTM0OTY1M2E0ZDY0OWZkZTlkYmM1MmE0NjdkMyIsImR2aWQiOiIxYTUyYTZhN2FkMWQ0MGMwODhkOWFlNGM0OTc3YzRmZiIsInQiOiJyIiwiY2xpZCI6IjM0YTAyY2Y4ZjQ0MTRlMjliMTU5MjE4NzZkYTM2ZjlhIiwiZXhwIjoxNzY0MzkyNzMxLCJhbSI6ImF1dGhvcml6YXRpb25fY29kZSIsImp0aSI6ImI0NDlkY2Y3MDhhYjQ1ZTc4MWIwNDZjZGIyMmQ4NGMyIn0.wepg8UsqqRAcDZ9k5jGF7yDpj4kUwq01QxJbvVMdOVNQC_vV63Q0VsxdJ4x1Q3OiFZJPubFfT9UeKsAG1227ttfp1OHi1W1m1MyHFUW-aWct19GWM7cXe2WA2QLphKFWhd0yyCIa8jQ1b4nVS5tqRiZPRTegw4OwETWWepLf-YlLzlrCAkEZXobf5k-thZp_fvwZD2hb5GfIZbFCMARolprH6XMonYkYn8AgzbJY66ZyVbMCE9jaHmHPpoqxiRjn_g8FFbjHqPZpCh4Kw4_adX_i4SjnyS9kPCTxqYmP9gRigB_Ft2cHd3Mszl00TMDxuT0wWSKfHx7pdtBiQ6lmt0OPK06LFsz7R3svDJzXCt5lgFjxlXTRGZeYHljiUlJZ8ozy7AD6AWTud4UWomm0NripFO3Mm_quBjUOsBrzgp5rYfEcKJ_-6DY6MMD61pBWhLFPENJjKcqchQt8Umpv-scmIkD1K12f2rGvwkZVjZskL4iPeel7HmSLBDlojEUq3fSgPxNk5CDoMaMMpCWSaDnc6fGtW9DbtOBeKyES_73I6NwvtwMGbb4Wv8Uru2-KRvszbcJ-FLe3ByHWS2dSTjCb_kftXeeuS2OJuOxHxTkqKG9GlMgO7WEzyjfTdldzllFX3G5WEQOPY0PIHx4zE3q8oKNqyxcq9j6IqO5fgzM")
	//if err != nil {
	//	log.Fatal(err)
	//}

	code, err := egs.GetExchangeCode(auth.AccessToken)
	if err != nil {
		log.Fatal(err)
	}

	authToken, err := egs.ExchangeEOSToken(code)
	if err != nil {
		log.Fatal(err)
	}

	//refreshToken, err := egs.RefreshEOSToken(authToken.AuthenticateWithRefreshToken)

	psyNet := rlapi.NewPsyNet()
	rpc, err := psyNet.AuthPlayer(rlapi.PlatformEpic, authToken.AccessToken, auth.AccountID, auth.DisplayName)
	if err != nil {
		log.Fatal(err)
	}
	defer rpc.Close()

	// Create context for API requests
	apiCtx, apiCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer apiCancel()

	shopsResp, err := rpc.GetStandardShops(apiCtx)
	if err != nil {
		log.Fatal(err)
	}

	// Pretty print shops
	shopsJSON, err := json.MarshalIndent(shopsResp, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Available Shops:")
	fmt.Println(string(shopsJSON))
	fmt.Println()

	var shopIDs []rlapi.ShopID
	for _, shop := range shopsResp.Shops {
		shopIDs = append(shopIDs, shop.ID)
	}

	catalogResp, err := rpc.GetShopCatalogue(apiCtx, shopIDs)
	if err != nil {
		log.Fatal(err)
	}

	// Pretty print catalogs
	catalogJSON, err := json.MarshalIndent(catalogResp, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Shop Catalogs:")
	fmt.Println(string(catalogJSON))
	fmt.Println()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}
