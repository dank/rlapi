package main

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/dank/rlapi"
	"github.com/dank/rlapi/examples/setup"
)

// Executes all available RPC requests to verify that each endpoint returns a valid response
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

	err = allMatchmaking(ctx, rpc, playerID)
	if err != nil {
		slog.Error("matchmaking error", slog.Any("err", err))
	}

	err = allParty(ctx, rpc)
	if err != nil {
		slog.Error("party error", slog.Any("err", err))
	}

	err = allPlayers(ctx, rpc, playerID)
	if err != nil {
		slog.Error("players error", slog.Any("err", err))
	}

	err = allSkills(ctx, rpc, playerID)
	if err != nil {
		slog.Error("skills error", slog.Any("err", err))
	}

	err = allStats(ctx, rpc, playerID)
	if err != nil {
		slog.Error("stats error", slog.Any("err", err))
	}

	err = allShops(ctx, rpc, playerID)
	if err != nil {
		slog.Error("shops error", slog.Any("err", err))
	}

	err = allProducts(ctx, rpc, playerID)
	if err != nil {
		slog.Error("products error", slog.Any("err", err))
	}

	err = allMTX(ctx, rpc)
	if err != nil {
		slog.Error("microtransactions error", slog.Any("err", err))
	}

	err = allRocketPass(ctx, rpc, playerID)
	if err != nil {
		slog.Error("rocketpass error", slog.Any("err", err))
	}

	err = allTournaments(ctx, rpc, playerID)
	if err != nil {
		slog.Error("tournaments error", slog.Any("err", err))
	}

	err = allTraining(ctx, rpc)
	if err != nil {
		slog.Error("training error", slog.Any("err", err))
	}

	err = allMisc(ctx, rpc, playerID)
	if err != nil {
		slog.Error("miscellaneous error", slog.Any("err", err))
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

	err = rpc.FTECheckpointComplete(ctx, playerID, "test", "test")
	slog.Debug("FTECheckpointComplete", slog.Any("err", err))

	err = rpc.FTEGroupComplete(ctx, playerID, "test")
	slog.Debug("FTEGroupComplete", slog.Any("err", err))

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

func allMatchmaking(ctx context.Context, rpc *rlapi.PsyNetRPC, playerID rlapi.PlayerID) error {
	regions := []rlapi.MatchmakingRegion{
		{Name: "USE1", Ping: 33},
		{Name: "USE3", Ping: 33},
	}
	playlists := []int{10, 11}

	estimatedQueueTime, err := rpc.StartMatchmaking(ctx, playlists, regions, false, "", []rlapi.PlayerID{playerID})
	if err != nil {
		return fmt.Errorf("StartMatchmaking err: %w", err)
	}
	slog.Debug("StartMatchmaking", slog.Any("estimatedQueueTime", estimatedQueueTime))

	err = rpc.PlayerCancelMatchmaking(ctx)
	if err != nil {
		return fmt.Errorf("PlayerCancelMatchmaking err: %w", err)
	}

	err = rpc.PlayerSearchPrivateMatch(ctx, "USE1", 6)
	if err != nil {
		return fmt.Errorf("PlayerSearchPrivateMatch err: %w", err)
	}

	return nil
}

func allParty(ctx context.Context, rpc *rlapi.PsyNetRPC) error {
	party, err := rpc.CreateParty(ctx)
	if err != nil {
		return fmt.Errorf("CreateParty err: %w", err)
	}
	slog.Debug("CreateParty", slog.Any("party", party))

	partyInfo, err := rpc.GetPlayerPartyInfo(ctx)
	if err != nil {
		return fmt.Errorf("GetPlayerPartyInfo err: %w", err)
	}
	slog.Debug("GetPlayerPartyInfo", slog.Any("partyInfo", partyInfo))

	partyID := rlapi.PartyID(party.Info.PartyID)

	err = rpc.SendPartyInvite(ctx, setup.RandomPlayerID, partyID)
	slog.Debug("SendPartyInvite", slog.Any("err", err))

	err = rpc.SendPartyJoinRequest(ctx, setup.RandomPlayerID)
	slog.Debug("SendPartyJoinRequest", slog.Any("err", err))

	updatedParty, err := rpc.ChangePartyOwner(ctx, setup.RandomPlayerID, partyID)
	slog.Debug("ChangePartyOwner", slog.Any("updatedParty", updatedParty), slog.Any("err", err))

	err = rpc.KickPartyMembers(ctx, []rlapi.PlayerID{setup.RandomPlayerID}, 1, partyID)
	slog.Debug("KickPartyMembers", slog.Any("err", err))

	err = rpc.SendPartyChatMessage(ctx, "hello", partyID)
	slog.Debug("SendPartyChatMessage", slog.Any("err", err))

	err = rpc.SendPartyMessage(ctx, "hello", partyID)
	slog.Debug("SendPartyMessage", slog.Any("err", err))

	newParty, err := rpc.JoinParty(ctx, "", partyID)
	slog.Debug("JoinParty", slog.Any("newParty", newParty), slog.Any("err", err))

	err = rpc.LeaveParty(ctx, partyID)
	slog.Debug("LeaveParty", slog.Any("err", err))

	return nil
}

func allPlayers(ctx context.Context, rpc *rlapi.PsyNetRPC, playerID rlapi.PlayerID) error {
	banStatus, err := rpc.GetBanStatus(ctx, []rlapi.PlayerID{playerID, setup.RandomPlayerID})
	if err != nil {
		return fmt.Errorf("GetBanStatus err: %w", err)
	}
	slog.Debug("GetBanStatus", slog.Any("banStatus", banStatus))

	profiles, err := rpc.GetProfiles(ctx, []rlapi.PlayerID{playerID, setup.RandomPlayerID})
	if err != nil {
		return fmt.Errorf("GetProfiles err: %w", err)
	}
	slog.Debug("GetProfiles", slog.Any("profiles", profiles))

	xp, err := rpc.GetXP(ctx, playerID)
	if err != nil {
		return fmt.Errorf("GetXP err: %w", err)
	}
	slog.Debug("GetXP", slog.Any("xp", xp))

	creatorCode, err := rpc.GetCreatorCode(ctx)
	slog.Debug("GetCreatorCode", slog.Any("creatorCode", creatorCode), slog.Any("err", err))

	err = rpc.ReportPlayer(ctx, []rlapi.Report{{Reporter: playerID, Offender: setup.RandomPlayerID, ReasonIDs: []int{3}, ReportTimestamp: 0.0}}, "")
	slog.Debug("ReportPlayer", slog.Any("err", err))

	return nil
}

func allSkills(ctx context.Context, rpc *rlapi.PsyNetRPC, playerID rlapi.PlayerID) error {
	playerSkill, err := rpc.GetPlayerSkill(ctx, setup.RandomPlayerID)
	if err != nil {
		return fmt.Errorf("GetPlayerSkill err: %w", err)
	}
	slog.Debug("GetPlayerSkill", slog.Any("playerSkill", playerSkill))

	playersSkills, err := rpc.GetPlayersSkills(ctx, []rlapi.PlayerID{playerID, setup.RandomPlayerID})
	if err != nil {
		return fmt.Errorf("GetPlayersSkills err: %w", err)
	}
	slog.Debug("GetPlayersSkills", slog.Any("playersSkills", playersSkills))

	leaderboard, err := rpc.GetSkillLeaderboard(ctx, 10, false)
	if err != nil {
		return fmt.Errorf("GetSkillLeaderboard err: %w", err)
	}
	slog.Debug("GetSkillLeaderboard", slog.Any("leaderboard", leaderboard))

	userValue, err := rpc.GetSkillLeaderboardValueForUser(ctx, 10, playerID)
	if err != nil {
		return fmt.Errorf("GetSkillLeaderboardValueForUser err: %w", err)
	}
	slog.Debug("GetSkillLeaderboardValueForUser", slog.Any("userValue", userValue))

	userRanks, err := rpc.GetSkillLeaderboardRankForUsers(ctx, 10, []rlapi.PlayerID{playerID})
	if err != nil {
		return fmt.Errorf("GetSkillLeaderboardRankForUsers err: %w", err)
	}
	slog.Debug("GetSkillLeaderboardRankForUsers", slog.Any("userRanks", userRanks))

	return nil
}

func allStats(ctx context.Context, rpc *rlapi.PsyNetRPC, playerID rlapi.PlayerID) error {
	statLeaderboard, err := rpc.GetStatLeaderboard(ctx, "Wins", false)
	if err != nil {
		return fmt.Errorf("GetStatLeaderboard err: %w", err)
	}
	slog.Debug("GetStatLeaderboard", slog.Any("statLeaderboard", statLeaderboard))

	statValue, err := rpc.GetStatLeaderboardValueForUser(ctx, "Wins", playerID)
	if err != nil {
		return fmt.Errorf("GetStatLeaderboardValueForUser err: %w", err)
	}
	slog.Debug("GetStatLeaderboardValueForUser", slog.Any("statValue", statValue))

	statRanks, err := rpc.GetStatLeaderboardRankForUsers(ctx, "Wins", []rlapi.PlayerID{playerID})
	if err != nil {
		return fmt.Errorf("GetStatLeaderboardRankForUsers err: %w", err)
	}
	slog.Debug("GetStatLeaderboardRankForUsers", slog.Any("statRanks", statRanks))

	return nil
}

func allShops(ctx context.Context, rpc *rlapi.PsyNetRPC, playerID rlapi.PlayerID) error {
	standardShops, err := rpc.GetStandardShops(ctx)
	if err != nil {
		return fmt.Errorf("GetStandardShops err: %w", err)
	}
	slog.Debug("GetStandardShops", slog.Any("standardShops", standardShops))

	if len(standardShops.Shops) > 0 {
		shopCatalogue, err := rpc.GetShopCatalogue(ctx, []rlapi.ShopID{standardShops.Shops[0].ID})
		if err != nil {
			return fmt.Errorf("GetShopCatalogue err: %w", err)
		}
		slog.Debug("GetShopCatalogue", slog.Any("shopCatalogue", shopCatalogue))
	}

	wallet, err := rpc.GetPlayerWallet(ctx, playerID)
	if err != nil {
		return fmt.Errorf("GetPlayerWallet err: %w", err)
	}
	slog.Debug("GetPlayerWallet", slog.Any("wallet", wallet))

	notifications, err := rpc.GetShopNotifications(ctx)
	if err != nil {
		return fmt.Errorf("GetShopNotifications err: %w", err)
	}
	slog.Debug("GetShopNotifications", slog.Any("notifications", notifications))

	return nil
}

func allProducts(ctx context.Context, rpc *rlapi.PsyNetRPC, playerID rlapi.PlayerID) error {
	products, err := rpc.GetPlayerProducts(ctx, playerID, "1756577989")
	if err != nil {
		return fmt.Errorf("GetPlayerProducts err: %w", err)
	}
	slog.Debug("GetPlayerProducts", slog.Any("products", products))

	dropTable, err := rpc.GetContainerDropTable(ctx)
	if err != nil {
		return fmt.Errorf("GetContainerDropTable err: %w", err)
	}
	slog.Debug("GetContainerDropTable", slog.Any("dropTable", dropTable))

	unlocked, err := rpc.UnlockContainer(ctx, playerID, []string{"90a79f045cad4556b95eea1270be0e76"})
	slog.Debug("UnlockContainer", slog.Any("unlocked", unlocked), slog.Any("err", err))

	traded, err := rpc.TradeIn(ctx, playerID, []string{"62d1e4bc3f5b4076bba8ea044a14a36c", "b3f60779705f4c8a926a14d6899bad70", "e610e81c23904c0696a25e9537aff4ba", "a10e3fa58af141b3829fec050889c967", "7cf368745031457dbfa7ec8ae4ab316f"})
	slog.Debug("TradeIn", slog.Any("traded", traded), slog.Any("err", err))

	crossStatus, err := rpc.GetCrossEntitlementProductStatus(ctx)
	if err != nil {
		return fmt.Errorf("GetCrossEntitlementProductStatus err: %w", err)
	}
	slog.Debug("GetCrossEntitlementProductStatus", slog.Any("crossStatus", crossStatus))

	return nil
}

func allMTX(ctx context.Context, rpc *rlapi.PsyNetRPC) error {
	catalog, err := rpc.GetMTXCatalog(ctx, "StarterPack")
	if err != nil {
		return fmt.Errorf("GetMTXCatalog err: %w", err)
	}
	slog.Debug("GetMTXCatalog", slog.Any("catalog", catalog))

	err = rpc.StartMTXPurchase(ctx, []rlapi.MTXCartItem{{CatalogID: 13, Count: 1}})
	slog.Debug("StartMTXPurchase", slog.Any("err", err))

	entitlements, err := rpc.ClaimMTXEntitlements(ctx, "test-auth-code")
	slog.Debug("ClaimMTXEntitlements", slog.Any("entitlements", entitlements), slog.Any("err", err))

	return nil
}

func allRocketPass(ctx context.Context, rpc *rlapi.PsyNetRPC, playerID rlapi.PlayerID) error {
	playerInfo, err := rpc.GetRocketPassPlayerInfo(ctx, playerID, 25)
	if err != nil {
		return fmt.Errorf("GetRocketPassPlayerInfo err: %w", err)
	}
	slog.Debug("GetRocketPassPlayerInfo", slog.Any("playerInfo", playerInfo))

	rewardContent, err := rpc.GetRocketPassRewardContent(ctx, 25, 0, 0, 0)
	if err != nil {
		return fmt.Errorf("GetRocketPassRewardContent err: %w", err)
	}
	slog.Debug("GetRocketPassRewardContent", slog.Any("rewardContent", rewardContent))

	prestigeRewards, err := rpc.GetRocketPassPrestigeRewards(ctx, playerID, 25)
	if err != nil {
		return fmt.Errorf("GetRocketPassPrestigeRewards err: %w", err)
	}
	slog.Debug("GetRocketPassPrestigeRewards", slog.Any("prestigeRewards", prestigeRewards))

	return nil
}

func allTournaments(ctx context.Context, rpc *rlapi.PsyNetRPC, playerID rlapi.PlayerID) error {
	scheduleRegion, err := rpc.GetTournamentScheduleRegion(ctx, playerID)
	if err != nil {
		return fmt.Errorf("GetTournamentScheduleRegion err: %w", err)
	}
	slog.Debug("GetTournamentScheduleRegion", slog.Any("scheduleRegion", scheduleRegion))

	cycleData, err := rpc.GetTournamentCycleData(ctx, playerID)
	if err != nil {
		return fmt.Errorf("GetTournamentCycleData err: %w", err)
	}
	slog.Debug("GetTournamentCycleData", slog.Any("cycleData", cycleData))

	schedule, err := rpc.GetTournamentSchedule(ctx, playerID, scheduleRegion)
	if err != nil {
		return fmt.Errorf("GetTournamentSchedule err: %w", err)
	}
	slog.Debug("GetTournamentSchedule", slog.Any("schedule", schedule))

	publicTournaments, err := rpc.GetPublicTournaments(ctx, playerID, rlapi.TournamentSearchInfo{Text: "", RankMin: -1, RankMax: 22, TeamSize: 1, BracketSize: 0, EnableCrossplay: true, StartTime: "1757216536", EndTime: "0", ShowFull: false, ShowIneligibleRank: false}, []rlapi.PlayerID{})
	slog.Debug("GetPublicTournaments", slog.Any("publicTournaments", publicTournaments), slog.Any("err", err))

	tournamentId := rlapi.TournamentID(strconv.Itoa(schedule[0].Tournaments[0].ID))
	tournament, err := rpc.RegisterTournament(ctx, playerID, tournamentId, rlapi.TournamentCredentials{})
	slog.Debug("RegisterTournament", slog.Any("tournament", tournament), slog.Any("err", err))

	err = rpc.UnsubscribeTournament(ctx, playerID, tournamentId, true, []rlapi.PlayerID{})
	slog.Debug("UnsubscribeTournament", slog.Any("err", err))

	subscriptions, err := rpc.GetTournamentSubscriptions(ctx, playerID)
	slog.Debug("GetTournamentSubscriptions", slog.Any("subscriptions", subscriptions), slog.Any("err", err))

	return nil
}

func allTraining(ctx context.Context, rpc *rlapi.PsyNetRPC) error {
	trainingPacks, err := rpc.BrowseTrainingData(ctx, true)
	if err != nil {
		return fmt.Errorf("BrowseTrainingData err: %w", err)
	}
	slog.Debug("BrowseTrainingData", slog.Any("trainingPacks", trainingPacks))

	metadata, err := rpc.GetTrainingMetadata(ctx, []string{"2BFC-F8D6-22AC-2AFE"})
	slog.Debug("GetTrainingMetadata", slog.Any("metadata", metadata), slog.Any("err", err))

	return nil
}

func allMisc(ctx context.Context, rpc *rlapi.PsyNetRPC, playerID rlapi.PlayerID) error {
	tradeIns, err := rpc.GetTradeInFilters(ctx)
	slog.Debug("GetTradeInFilters", slog.Any("tradeIns", tradeIns), slog.Any("err", err))

	private, err := rpc.GetClubPrivateMatches(ctx)
	slog.Debug("GetClubPrivateMatches", slog.Any("private", private), slog.Any("err", err))

	ping, err := rpc.GetGameServerPingList(ctx)
	slog.Debug("GetGameServerPingList", slog.Any("ping", ping), slog.Any("err", err))

	matches, err := rpc.GetMatchHistory(ctx, playerID)
	slog.Debug("GetMatchHistory", slog.Any("matches", matches), slog.Any("err", err))

	regions, err := rpc.GetSubRegions(ctx)
	if err != nil {
		return fmt.Errorf("GetSubRegions err: %w", err)
	}
	slog.Debug("GetSubRegions", slog.Any("regions", regions))

	playlists, err := rpc.GetActivePlaylists(ctx)
	if err != nil {
		return fmt.Errorf("GetActivePlaylists err: %w", err)
	}
	slog.Debug("GetActivePlaylists", slog.Any("playlists", playlists))

	population, err := rpc.GetPopulation(ctx)
	if err != nil {
		return fmt.Errorf("GetPopulation err: %w", err)
	}
	slog.Debug("GetPopulation", slog.Any("population", population))

	err = rpc.UpdatePlayerPlaylist(ctx, 10, 1)
	if err != nil {
		return fmt.Errorf("UpdatePlayerPlaylist err: %w", err)
	}

	joinResult, err := rpc.JoinMatch(ctx, "JoinPrivate", "server", "password")
	slog.Debug("JoinMatch", slog.Any("joinResult", joinResult), slog.Any("err", err))

	filtered, err := rpc.FilterContent(ctx, []string{"test", "test2"}, "Content")
	slog.Debug("FilterContent", slog.Any("filtered", filtered), slog.Any("err", err))

	avatarStatus, err := rpc.CanShowAvatar(ctx, []rlapi.PlayerID{playerID})
	slog.Debug("CanShowAvatar", slog.Any("avatarStatus", avatarStatus), slog.Any("err", err))

	return nil
}
