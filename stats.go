package rlapi

import "context"

// StatLeaderboard represents a statistical leaderboard
type StatLeaderboard struct {
	LeaderboardID string                    `json:"LeaderboardID"`
	Platforms     []StatPlatformLeaderboard `json:"Platforms"`
}

// StatPlatformLeaderboard represents leaderboard data for a specific platform
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

// StatLeaderboardRankPlayer represents a player's rank data
type StatLeaderboardRankPlayer struct {
	PlayerID   PlayerID `json:"PlayerID"`
	PlayerName string   `json:"PlayerName"`
	Value      float64  `json:"Value"`
	Rank       int      `json:"Rank"`
}

// Request and Response types

type GetStatLeaderboardRequest struct {
	StatName          string `json:"StatName"`
	BDisableCrossplay bool   `json:"bDisableCrossplay"`
}

type GetStatLeaderboardResponse struct {
	LeaderboardID string                    `json:"LeaderboardID"`
	Platforms     []StatPlatformLeaderboard `json:"Platforms"`
}

type GetStatLeaderboardValueForUserRequest struct {
	StatName string   `json:"StatName"`
	PlayerID PlayerID `json:"PlayerID"`
}

type GetStatLeaderboardValueForUserResponse struct {
	LeaderboardID string  `json:"LeaderboardID"`
	BHasStat      bool    `json:"bHasStat"`
	Value         float64 `json:"Value"`
	Rank          int     `json:"Rank"`
}

type GetStatLeaderboardRankForUsersRequest struct {
	StatName  string     `json:"StatName"`
	PlayerIDs []PlayerID `json:"PlayerIDs"`
}

type GetStatLeaderboardRankForUsersResponse struct {
	LeaderboardID string                      `json:"LeaderboardID"`
	Players       []StatLeaderboardRankPlayer `json:"Players"`
}

// GetStatLeaderboard retrieves the statistical leaderboard for a specific stat.
func (p *PsyNetRPC) GetStatLeaderboard(ctx context.Context, statName string, disableCrossplay bool) (*GetStatLeaderboardResponse, error) {
	request := GetStatLeaderboardRequest{
		StatName:          statName,
		BDisableCrossplay: disableCrossplay,
	}

	var result GetStatLeaderboardResponse
	err := p.sendRequestSync(ctx, "Stats/GetStatLeaderboard v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetStatLeaderboardValueForUser retrieves a specific player's position and data on a stat leaderboard.
func (p *PsyNetRPC) GetStatLeaderboardValueForUser(ctx context.Context, statName string, playerID PlayerID) (*GetStatLeaderboardValueForUserResponse, error) {
	request := GetStatLeaderboardValueForUserRequest{
		StatName: statName,
		PlayerID: playerID,
	}

	var result GetStatLeaderboardValueForUserResponse
	err := p.sendRequestSync(ctx, "Stats/GetStatLeaderboardValueForUser v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetStatLeaderboardRankForUsers retrieves rank information for multiple players on a stat leaderboard.
func (p *PsyNetRPC) GetStatLeaderboardRankForUsers(ctx context.Context, statName string, playerIDs []PlayerID) (*GetStatLeaderboardRankForUsersResponse, error) {
	request := GetStatLeaderboardRankForUsersRequest{
		StatName:  statName,
		PlayerIDs: playerIDs,
	}

	var result GetStatLeaderboardRankForUsersResponse
	err := p.sendRequestSync(ctx, "Stats/GetStatLeaderboardRankForUsers v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
