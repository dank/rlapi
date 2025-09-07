package rlapi

import "context"

// StatPlatformLeaderboard represents leaderboard data for a given platform
type StatPlatformLeaderboard struct {
	Platform string                  `json:"Platform"`
	Players  []StatLeaderboardPlayer `json:"Players"`
}

// StatLeaderboardPlayer represents a player entry in a stat leaderboard
type StatLeaderboardPlayer struct {
	PlayerID   PlayerID `json:"PlayerID"`
	PlayerName string   `json:"PlayerName"`
	Value      float64  `json:"Value"`
	Rank       int      `json:"Rank"`
}

type StatLeaderboardRankPlayer struct {
	PlayerID   PlayerID `json:"PlayerID"`
	PlayerName string   `json:"PlayerName"`
	Value      float64  `json:"Value"`
	Rank       int      `json:"Rank"`
}

type GetStatLeaderboardRequest struct {
	Stat             string `json:"Stat"`
	DisableCrossplay bool   `json:"bDisableCrossplay"`
}

type GetStatLeaderboardResponse struct {
	LeaderboardID string                    `json:"LeaderboardID"`
	Platforms     []StatPlatformLeaderboard `json:"Platforms"`
}

type GetStatLeaderboardValueForUserRequest struct {
	Stat     string   `json:"Stat"`
	PlayerID PlayerID `json:"PlayerID"`
}

type GetStatLeaderboardValueForUserResponse struct {
	LeaderboardID string `json:"LeaderboardID"`
	HasStat       bool   `json:"bHasStat"`
	Value         string `json:"Value"`
	Rank          int    `json:"Rank"`
}

type GetStatLeaderboardRankForUsersRequest struct {
	Stat      string     `json:"Stat"`
	PlayerIDs []PlayerID `json:"PlayerIDs"`
}

type GetStatLeaderboardRankForUsersResponse struct {
	LeaderboardID string                      `json:"LeaderboardID"`
	Players       []StatLeaderboardRankPlayer `json:"Players"`
}

// GetStatLeaderboard retrieves the stats leaderboard for a given stat.
func (p *PsyNetRPC) GetStatLeaderboard(ctx context.Context, statName string, disableCrossplay bool) (*GetStatLeaderboardResponse, error) {
	request := GetStatLeaderboardRequest{
		Stat:             statName,
		DisableCrossplay: disableCrossplay,
	}

	var result GetStatLeaderboardResponse
	err := p.sendRequestSync(ctx, "Stats/GetStatLeaderboard v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetStatLeaderboardValueForUser retrieves a player's position and data on a stat leaderboard.
func (p *PsyNetRPC) GetStatLeaderboardValueForUser(ctx context.Context, statName string, playerID PlayerID) (*GetStatLeaderboardValueForUserResponse, error) {
	request := GetStatLeaderboardValueForUserRequest{
		Stat:     statName,
		PlayerID: playerID,
	}

	var result GetStatLeaderboardValueForUserResponse
	err := p.sendRequestSync(ctx, "Stats/GetStatLeaderboardValueForUser v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *PsyNetRPC) GetStatLeaderboardRankForUsers(ctx context.Context, statName string, playerIDs []PlayerID) (*GetStatLeaderboardRankForUsersResponse, error) {
	request := GetStatLeaderboardRankForUsersRequest{
		Stat:      statName,
		PlayerIDs: playerIDs,
	}

	var result GetStatLeaderboardRankForUsersResponse
	err := p.sendRequestSync(ctx, "Stats/GetStatLeaderboardRankForUsers v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
