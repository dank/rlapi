package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/smtp"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/dank/rlapi"
	"github.com/dank/rlapi/examples/setup"
)

// Watchlist matches a shop item by any of its product IDs or by a keyword in its title (case-insensitive).
type Watchlist struct {
	ProductIDs []int    `json:"product_ids"`
	Keywords   []string `json:"keywords"`
}

// Polls the item shop every 5 minutes and sends a notification when an item from watchlist.json goes live.
// Notifications go to ntfy (NTFY_TOPIC) and/or email (SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASS, NOTIFY_EMAIL_TO).
// Run with debug logging to see the product IDs of every live item.
func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	data, err := os.ReadFile("watchlist.json")
	if err != nil {
		log.Fatalf("Failed to read watchlist.json: %v", err)
	}
	var watchlist Watchlist
	if err := json.Unmarshal(data, &watchlist); err != nil {
		log.Fatalf("Failed to parse watchlist.json: %v", err)
	}

	rpc, _ := setup.RPC()
	defer rpc.Close()

	notified := map[int]bool{}
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	if err := checkShops(rpc, watchlist, notified); err != nil {
		slog.Error("Initial shop check failed", slog.Any("error", err))
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			if err := checkShops(rpc, watchlist, notified); err != nil {
				slog.Error("Failed to check shops", slog.Any("error", err))
			}
		case <-sigChan:
			slog.Info("Shutting down shop watcher")
			return
		}
	}
}

func checkShops(rpc *rlapi.PsyNetRPC, watchlist Watchlist, notified map[int]bool) error {
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

	now := time.Now().Unix()
	for _, catalog := range catalogResp.Catalogues {
		for _, item := range catalog.ShopItems {
			if item.StartDate > now || (item.EndDate != nil && *item.EndDate <= now) || notified[item.ShopItemID] {
				continue
			}

			title := fmt.Sprintf("Shop item %d", item.ShopItemID)
			if item.Title != nil {
				title = *item.Title
			}
			var productIDs []int
			for _, p := range item.DeliverableProducts {
				productIDs = append(productIDs, p.Product.ProductID)
			}
			slog.Debug("Live item", slog.String("title", title), slog.Any("product_ids", productIDs))

			matched := slices.ContainsFunc(productIDs, func(id int) bool { return slices.Contains(watchlist.ProductIDs, id) }) ||
				slices.ContainsFunc(watchlist.Keywords, func(kw string) bool {
					return kw != "" && strings.Contains(strings.ToLower(title), strings.ToLower(kw))
				})
			if !matched {
				continue
			}

			slog.Info("Watched item is live", slog.String("title", title))
			if err := notify("Rocket League item shop: "+title, fmt.Sprintf("%s is live in the item shop (products %v).", title, productIDs)); err != nil {
				slog.Error("Failed to send notification", slog.Any("error", err))
				continue // retry on the next check
			}
			notified[item.ShopItemID] = true
		}
	}
	return nil
}

func notify(subject, body string) error {
	if topic := os.Getenv("NTFY_TOPIC"); topic != "" {
		req, err := http.NewRequest(http.MethodPost, "https://ntfy.sh/"+topic, strings.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Title", subject)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()
		if resp.StatusCode >= 300 {
			return fmt.Errorf("ntfy returned %s", resp.Status)
		}
	}

	if host, to := os.Getenv("SMTP_HOST"), os.Getenv("NOTIFY_EMAIL_TO"); host != "" && to != "" {
		user := os.Getenv("SMTP_USER")
		port := os.Getenv("SMTP_PORT")
		if port == "" {
			port = "587"
		}
		msg := "To: " + to + "\r\nSubject: " + subject + "\r\n\r\n" + body
		auth := smtp.PlainAuth("", user, os.Getenv("SMTP_PASS"), host)
		return smtp.SendMail(host+":"+port, auth, user, []string{to}, []byte(msg))
	}
	return nil
}
