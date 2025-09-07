package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dank/rlapi"
	"github.com/dank/rlapi/examples/setup"
)

// Polls the item shop every 5 minutes and logs when changes are detected.
// For more durability auth tokens should be automatically refreshed, and connections re-established on disconnect.
func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	rpc, _ := setup.RPC()
	defer rpc.Close()

	var lastHash string
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	if err := pollShops(rpc, &lastHash); err != nil {
		slog.Error("Initial shop poll failed", slog.Any("error", err))
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			if err := pollShops(rpc, &lastHash); err != nil {
				slog.Error("Failed to poll shops", slog.Any("error", err))
			}
		case <-sigChan:
			slog.Info("Shutting down item shop monitor")
			return
		}
	}
}

func pollShops(rpc *rlapi.PsyNetRPC, lastHash *string) error {
	apiCtx, apiCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer apiCancel()

	shopsResp, err := rpc.GetStandardShops(apiCtx)
	if err != nil {
		return fmt.Errorf("failed to get shops: %w", err)
	}

	var shopIDs []rlapi.ShopID
	for _, shop := range shopsResp.Shops {
		shopIDs = append(shopIDs, shop.ID)
	}

	catalogResp, err := rpc.GetShopCatalogue(apiCtx, shopIDs)
	if err != nil {
		return fmt.Errorf("failed to get shop catalogue: %w", err)
	}

	// Create hash of catalog content for change detection
	catalogJSON, err := json.Marshal(catalogResp)
	if err != nil {
		return fmt.Errorf("failed to marshal catalog: %w", err)
	}

	hash := fmt.Sprintf("%x", sha256.Sum256(catalogJSON))

	if *lastHash == "" {
		slog.Info("Initial item shop data loaded")
		*lastHash = hash
	} else if *lastHash != hash {
		slog.Info("CHANGE DETECTED: Item shop has been updated!")
		*lastHash = hash

		// Pretty print the updated catalog
		catalogPretty, err := json.MarshalIndent(catalogResp, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format catalog: %w", err)
		}

		fmt.Println("Updated Item Shop Catalog:")
		fmt.Println(string(catalogPretty))
		fmt.Println()
	} else {
		slog.Debug("No changes detected", slog.String("checked_at", time.Now().Format("15:04:05")))
	}
	return nil
}
