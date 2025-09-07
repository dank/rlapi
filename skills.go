package rlapi

import "context"

// Skill represents a player's skill data for the specified playlist
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

// LeaderboardPlayer represents a player entry in a skill leaderboard
type LeaderboardPlayer struct {
	PlayerID   PlayerID `json:"PlayerID"`
	PlayerName string   `json:"PlayerName"`
	MMR        float64  `json:"MMR"`
	Value      int      `json:"Value"`
}

// PlatformLeaderboard represents leaderboard data for a given platform
type PlatformLeaderboard struct {
	Platform string              `json:"Platform"`
	Players  []LeaderboardPlayer `json:"Players"`
}

type LeaderboardRankPlayer struct {
	PlayerID   string `json:"PlayerID"`
	PlayerName string `json:"PlayerName"`
	Value      int    `json:"Value"`
}

type PlayerWithSkills struct {
	PlayerID PlayerID `json:"PlayerID"`
	Skills   []Skill  `json:"Skills"`
}

type GetPlayerSkillRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
}

type GetPlayerSkillResponse struct {
	Skills       []Skill      `json:"Skills"`
	RewardLevels RewardLevels `json:"RewardLevels"`
}

type GetSkillLeaderboardRequest struct {
	Playlist         PlaylistID `json:"Playlist"`
	DisableCrossplay bool       `json:"bDisableCrossplay"`
}

type GetSkillLeaderboardResponse struct {
	LeaderboardID string                `json:"LeaderboardID"`
	Platforms     []PlatformLeaderboard `json:"Platforms"`
}

type GetSkillLeaderboardValueForUserRequest struct {
	Playlist PlaylistID `json:"Playlist"`
	PlayerID PlayerID   `json:"PlayerID"`
}

type GetSkillLeaderboardValueForUserResponse struct {
	LeaderboardID string  `json:"LeaderboardID"`
	HasSkill      bool    `json:"bHasSkill"`
	MMR           float64 `json:"MMR"`
	Value         int     `json:"Value"`
}

type GetSkillLeaderboardRankForUsersRequest struct {
	Playlist  PlaylistID `json:"Playlist"`
	PlayerIDs []PlayerID `json:"PlayerIDs"`
}

type GetSkillLeaderboardRankForUsersResponse struct {
	LeaderboardID string                  `json:"LeaderboardID"`
	Players       []LeaderboardRankPlayer `json:"Players"`
}

type GetPlayersSkillsRequest struct {
	PlayerIDs []PlayerID `json:"PlayerIDs"`
}

type GetPlayersSkillsResponse struct {
	Players []PlayerWithSkills `json:"Players"`
}

// GetPlayerSkill retrieves skill data for a given player.
func (p *PsyNetRPC) GetPlayerSkill(ctx context.Context, playerID PlayerID) (*GetPlayerSkillResponse, error) {
	request := GetPlayerSkillRequest{
		PlayerID: playerID,
	}

	var result GetPlayerSkillResponse
	err := p.sendRequestSync(ctx, "Skills/GetPlayerSkill v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPlayersSkills retrieves skill data for given players.
func (p *PsyNetRPC) GetPlayersSkills(ctx context.Context, playerIDs []PlayerID) ([]PlayerWithSkills, error) {
	request := GetPlayersSkillsRequest{
		PlayerIDs: playerIDs,
	}

	var result GetPlayersSkillsResponse
	err := p.sendRequestSync(ctx, "Skills/GetPlayersSkills v1", request, &result)
	if err != nil {
		return nil, err
	}
	return result.Players, nil
}

// GetSkillLeaderboard retrieves the skill leaderboard on all platforms for a given playlist.
func (p *PsyNetRPC) GetSkillLeaderboard(ctx context.Context, playlist PlaylistID, disableCrossplay bool) (*GetSkillLeaderboardResponse, error) {
	request := GetSkillLeaderboardRequest{
		Playlist:         playlist,
		DisableCrossplay: disableCrossplay,
	}

	var result GetSkillLeaderboardResponse
	err := p.sendRequestSync(ctx, "Skills/GetSkillLeaderboard v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSkillLeaderboardValueForUser retrieves a player's position and data on the skill leaderboard for a given playlist.
func (p *PsyNetRPC) GetSkillLeaderboardValueForUser(ctx context.Context, playlist PlaylistID, playerID PlayerID) (*GetSkillLeaderboardValueForUserResponse, error) {
	request := GetSkillLeaderboardValueForUserRequest{
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

// GetSkillLeaderboardRankForUsers retrieves rank information for multiple players on the skill leaderboard.
func (p *PsyNetRPC) GetSkillLeaderboardRankForUsers(ctx context.Context, playlist PlaylistID, playerIDs []PlayerID) (*GetSkillLeaderboardRankForUsersResponse, error) {
	request := GetSkillLeaderboardRankForUsersRequest{
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
