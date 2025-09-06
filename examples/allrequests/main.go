package main

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/dank/rlapi"
	"github.com/dank/rlapi/examples/setup"
)

// Run all RPC requests to verify all of them return a valid response
func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	rpc, playerID := setup.RPC()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	challenges, err := rpc.GetActiveChallenges(ctx)
	if err != nil {
		log.Fatalf("GetActiveChallenges err: %v", err)
	}
	slog.Debug("GetActiveChallenges", slog.Any("challenges", challenges))

	var challengeIDs []rlapi.ChallengeID
	for _, challenge := range challenges {
		challengeIDs = append(challengeIDs, challenge.ID)
	}

	progress, err := rpc.GetChallengeProgress(ctx, playerID, challengeIDs)
	if err != nil {
		log.Fatalf("GetChallengeProgress err: %v", err)
	}
	slog.Debug("GetChallengeProgress", slog.Any("progress", progress))

	err = rpc.CollectChallengeReward(ctx, playerID, challengeIDs[0])
	if err != nil {
		log.Fatalf("CollectChallengeReward err: %v", err)
	}
}
