package main

import (
	"context"
	"fmt"
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

	err := allChallenges(ctx, rpc, playerID)
	if err != nil {
		slog.Error("challenges error", slog.Any("err", err))
	}

	err = allClubs(ctx, rpc)
	if err != nil {
		slog.Error("clubs error", slog.Any("err", err))
	}
}

func allChallenges(ctx context.Context, rpc *rlapi.PsyNetRPC, playerID rlapi.PlayerID) error {
	challenges, err := rpc.GetActiveChallenges(ctx)
	if err != nil {
		return fmt.Errorf("GetActiveChallenges err: %w", err)
	}
	slog.Debug("GetActiveChallenges", slog.Any("challenges", challenges))

	var challengeIDs []rlapi.ChallengeID
	for _, challenge := range challenges {
		challengeIDs = append(challengeIDs, challenge.ID)
	}

	progress, err := rpc.GetChallengeProgress(ctx, playerID, challengeIDs)
	if err != nil {
		return fmt.Errorf("GetChallengeProgress err: %w", err)
	}
	slog.Debug("GetChallengeProgress", slog.Any("progress", progress))

	err = rpc.CollectChallengeReward(ctx, playerID, challengeIDs[0])
	if err != nil {
		return fmt.Errorf("CollectChallengeReward err: %w", err)
	}

	return nil
}

func allClubs(ctx context.Context, rpc *rlapi.PsyNetRPC) error {
	club, err := rpc.CreateClub(ctx, setup.RandString(10), setup.RandString(3), 0, 0)
	if err != nil {
		return fmt.Errorf("CreateClub err: %w", err)
	}
	slog.Debug("CreateClub", slog.Any("club", club))

	club, err = rpc.GetClubDetails(ctx, club.ClubID)
	if err != nil {
		return fmt.Errorf("GetClubDetails err: %w", err)
	}
	slog.Debug("GetClubDetails", slog.Any("club", club))

	invites, err := rpc.GetClubInvites(ctx)
	if err != nil {
		return fmt.Errorf("GetClubInvites err: %w", err)
	}
	slog.Debug("GetClubInvites", slog.Any("invites", invites))

	titles, err := rpc.GetClubTitleInstances(ctx)
	if err != nil {
		return fmt.Errorf("GetClubTitleInstances err: %w", err)
	}
	slog.Debug("GetClubTitleInstances", slog.Any("titles", titles))

	stats, err := rpc.GetClubStats(ctx)
	if err != nil {
		return fmt.Errorf("GetClubStats err: %w", err)
	}
	slog.Debug("GetClubStats", slog.Any("stats", stats))

	primary := -10879077
	accent := -1710619
	updatedClub, err := rpc.UpdateClub(ctx, &rlapi.UpdateClubRequest{PrimaryColor: &primary, AccentColor: &accent})
	if err != nil {
		return fmt.Errorf("UpdateClub err: %w", err)
	}
	slog.Debug("UpdateClub", slog.Any("updatedClub", updatedClub))

	err = rpc.InviteToClub(ctx, setup.RandomPlayerID)
	if err != nil {
		return fmt.Errorf("InviteToClub err: %w", err)
	}

	playerClub, err := rpc.GetPlayerClubDetails(ctx, setup.RandomPlayerID)
	if err != nil {
		return fmt.Errorf("GetPlayerClubDetails err: %w", err)
	}
	slog.Debug("GetPlayerClubDetails", slog.Any("playerClub", playerClub))

	err = rpc.LeaveClub(ctx)
	if err != nil {
		return fmt.Errorf("LeaveClub err: %w", err)
	}

	acceptedClub, err := rpc.AcceptClubInvite(ctx, playerClub.ClubID)
	slog.Debug("AcceptClub", slog.Any("acceptedClub", acceptedClub), slog.Any("err", err))

	err = rpc.RejectClubInvite(ctx, playerClub.ClubID)
	slog.Debug("RejectClub", slog.Any("err", err))

	return nil
}
