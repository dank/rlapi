package rlapi

import "context"

// Skill represents a player's skill data for a specific playlist
type Skill struct {
	Playlist               int     `json:"Playlist"`
	Mu                     float64 `json:"Mu"`
	Sigma                  float64 `json:"Sigma"`
	Tier                   int     `json:"Tier"`
	Division               int     `json:"Division"`
	MMR                    float64 `json:"MMR"`
	WinStreak              int     `json:"WinStreak"`
	MatchesPlayed          int     `json:"MatchesPlayed"`
	PlacementMatchesPlayed int     `json:"PlacementMatchesPlayed"`
}

// RewardLevels represents seasonal reward level information
type RewardLevels struct {
	SeasonLevel     int `json:"SeasonLevel"`
	SeasonLevelWins int `json:"SeasonLevelWins"`
}

// PlayerSkills represents a player's skills across all playlists
type PlayerSkills struct {
	Skills       []Skill      `json:"Skills"`
	RewardLevels RewardLevels `json:"RewardLevels"`
}

// LeaderboardPlayer represents a player entry in a skill leaderboard
type LeaderboardPlayer struct {
	PlayerID   PlayerID `json:"PlayerID"`
	PlayerName string   `json:"PlayerName"`
	MMR        float64  `json:"MMR"`
	Value      int      `json:"Value"`
}

// PlatformLeaderboard represents leaderboard data for a specific platform
type PlatformLeaderboard struct {
	Platform string              `json:"Platform"`
	Players  []LeaderboardPlayer `json:"Players"`
}

// SkillLeaderboard represents the complete leaderboard for a playlist
type SkillLeaderboard struct {
	LeaderboardID string                `json:"LeaderboardID"`
	Platforms     []PlatformLeaderboard `json:"Platforms"`
}

// LeaderboardRankPlayer represents a player's rank data
type LeaderboardRankPlayer struct {
	PlayerID   string `json:"PlayerID"`
	PlayerName string `json:"PlayerName"`
	Value      int    `json:"Value"`
}

// PlayerWithSkills represents a player with their complete skill set
type PlayerWithSkills struct {
	PlayerID PlayerID `json:"PlayerID"`
	Skills   []Skill  `json:"Skills"`
}

// Request and Response types

type getPlayerSkillRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
}

type GetPlayerSkillResponse struct {
	Skills       []Skill      `json:"Skills"`
	RewardLevels RewardLevels `json:"RewardLevels"`
}

type getSkillLeaderboardRequest struct {
	Playlist          int  `json:"Playlist"`
	BDisableCrossplay bool `json:"bDisableCrossplay"`
}

type GetSkillLeaderboardResponse struct {
	LeaderboardID string                `json:"LeaderboardID"`
	Platforms     []PlatformLeaderboard `json:"Platforms"`
}

type getSkillLeaderboardValueForUserRequest struct {
	Playlist int      `json:"Playlist"`
	PlayerID PlayerID `json:"PlayerID"`
}

type GetSkillLeaderboardValueForUserResponse struct {
	LeaderboardID string  `json:"LeaderboardID"`
	BHasSkill     bool    `json:"bHasSkill"`
	MMR           float64 `json:"MMR"`
	Value         int     `json:"Value"`
}

type getSkillLeaderboardRankForUsersRequest struct {
	Playlist  int        `json:"Playlist"`
	PlayerIDs []PlayerID `json:"PlayerIDs"`
}

type GetSkillLeaderboardRankForUsersResponse struct {
	LeaderboardID string                  `json:"LeaderboardID"`
	Players       []LeaderboardRankPlayer `json:"Players"`
}

type getPlayersSkillsRequest struct {
	PlayerIDs []PlayerID `json:"PlayerIDs"`
}

type GetPlayersSkillsResponse struct {
	Players []PlayerWithSkills `json:"Players"`
}

// GetPlayerSkill retrieves skill data for a specific player.
func (p *PsyNetRPC) GetPlayerSkill(ctx context.Context, playerID PlayerID) (*GetPlayerSkillResponse, error) {
	request := getPlayerSkillRequest{
		PlayerID: playerID,
	}

	var result GetPlayerSkillResponse
	err := p.sendRequestSync(ctx, "Skills/GetPlayerSkill v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSkillLeaderboard retrieves the skill leaderboard for a specific playlist.
func (p *PsyNetRPC) GetSkillLeaderboard(ctx context.Context, playlist int, disableCrossplay bool) (*GetSkillLeaderboardResponse, error) {
	request := getSkillLeaderboardRequest{
		Playlist:          playlist,
		BDisableCrossplay: disableCrossplay,
	}

	var result GetSkillLeaderboardResponse
	err := p.sendRequestSync(ctx, "Skills/GetSkillLeaderboard v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSkillLeaderboardValueForUser retrieves a specific player's position and data on a skill leaderboard.
func (p *PsyNetRPC) GetSkillLeaderboardValueForUser(ctx context.Context, playlist int, playerID PlayerID) (*GetSkillLeaderboardValueForUserResponse, error) {
	request := getSkillLeaderboardValueForUserRequest{
		Playlist: playlist,
		PlayerID: playerID,
	}

	var result GetSkillLeaderboardValueForUserResponse
	err := p.sendRequestSync(ctx, "Skills/GetSkillLeaderboardValueForUser v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSkillLeaderboardRankForUsers retrieves rank information for multiple players on a skill leaderboard.
func (p *PsyNetRPC) GetSkillLeaderboardRankForUsers(ctx context.Context, playlist int, playerIDs []PlayerID) (*GetSkillLeaderboardRankForUsersResponse, error) {
	request := getSkillLeaderboardRankForUsersRequest{
		Playlist:  playlist,
		PlayerIDs: playerIDs,
	}

	var result GetSkillLeaderboardRankForUsersResponse
	err := p.sendRequestSync(ctx, "Skills/GetSkillLeaderboardRankForUsers v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPlayersSkills retrieves skill data for multiple players.
func (p *PsyNetRPC) GetPlayersSkills(ctx context.Context, playerIDs []PlayerID) (*GetPlayersSkillsResponse, error) {
	request := getPlayersSkillsRequest{
		PlayerIDs: playerIDs,
	}

	var result GetPlayersSkillsResponse
	err := p.sendRequestSync(ctx, "Skills/GetPlayersSkills v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
