package rlapi

import "context"

// BanMessage represents a ban message for a player
type BanMessage struct {
	PlayerID    PlayerID `json:"PlayerID"`
	BanType     string   `json:"BanType"`
	BanReason   string   `json:"BanReason"`
	StartTime   int64    `json:"StartTime"`
	EndTime     *int64   `json:"EndTime"`
	Description string   `json:"Description"`
}

// PlayerProfile represents a player's profile information
type PlayerProfile struct {
	PlayerID          PlayerID               `json:"PlayerID"`
	PlayerName        string                 `json:"PlayerName"`
	Avatar            PlayerAvatar           `json:"Avatar"`
	Platform          string                 `json:"Platform"`
	PrimaryTitle      *string                `json:"PrimaryTitle"`
	SeasonRewardLevel int                    `json:"SeasonRewardLevel"`
	XPLevel           int                    `json:"XPLevel"`
	TotalXP           int64                  `json:"TotalXP"`
	Stats             map[string]interface{} `json:"Stats"`
	CreatedTime       int64                  `json:"CreatedTime"`
	LastSeenTime      *int64                 `json:"LastSeenTime"`
}

// PlayerAvatar represents a player's avatar information
type PlayerAvatar struct {
	AvatarURL *string `json:"AvatarURL"`
	BorderURL *string `json:"BorderURL"`
	AvatarID  *string `json:"AvatarID"`
	BorderID  *string `json:"BorderID"`
}

// PlayerXP represents a player's XP information
type PlayerXP struct {
	PlayerID      PlayerID `json:"PlayerID"`
	Level         int      `json:"Level"`
	CurrentXP     int      `json:"CurrentXP"`
	XPToNextLevel int      `json:"XPToNextLevel"`
	TotalXP       int64    `json:"TotalXP"`
}

// CreatorCode represents a creator code information
type CreatorCode struct {
	Code        string `json:"Code"`
	CreatorName string `json:"CreatorName"`
	IsActive    bool   `json:"IsActive"`
}

// ReportReason represents reasons for reporting a player
type ReportReason struct {
	ReasonID    int    `json:"ReasonID"`
	Description string `json:"Description"`
}

// Request and Response types

type GetBanStatusRequest struct {
	Players []PlayerID `json:"Players"`
}

type GetBanStatusResponse struct {
	BanMessages []BanMessage `json:"BanMessages"`
}

type GetProfileRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
}

type GetProfileResponse struct {
	Profile PlayerProfile `json:"Profile"`
}

type GetXPRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
}

type GetXPResponse struct {
	XP PlayerXP `json:"XP"`
}

type GetCreatorCodeRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
}

type GetCreatorCodeResponse struct {
	CreatorCode *CreatorCode `json:"CreatorCode"`
}

type ReportRequest struct {
	ReporterID  PlayerID `json:"ReporterID"`
	ReportedID  PlayerID `json:"ReportedID"`
	ReasonID    int      `json:"ReasonID"`
	Description string   `json:"Description"`
	GameID      *string  `json:"GameID,omitempty"`
	MatchID     *string  `json:"MatchID,omitempty"`
}

type ReportResponse struct {
	Success  bool   `json:"Success"`
	ReportID string `json:"ReportID"`
	Message  string `json:"Message"`
}

// GetBanStatus retrieves ban status information for specified players.
func (p *PsyNetRPC) GetBanStatus(ctx context.Context, playerIDs []PlayerID) (*GetBanStatusResponse, error) {
	request := GetBanStatusRequest{
		Players: playerIDs,
	}

	var result GetBanStatusResponse
	err := p.sendRequestSync(ctx, "Players/GetBanStatus v3", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetProfile retrieves profile information for a specific player.
func (p *PsyNetRPC) GetProfile(ctx context.Context, playerID PlayerID) (*GetProfileResponse, error) {
	request := GetProfileRequest{
		PlayerID: playerID,
	}

	var result GetProfileResponse
	err := p.sendRequestSync(ctx, "Players/GetProfile v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetXP retrieves XP information for a specific player.
func (p *PsyNetRPC) GetXP(ctx context.Context, playerID PlayerID) (*GetXPResponse, error) {
	request := GetXPRequest{
		PlayerID: playerID,
	}

	var result GetXPResponse
	err := p.sendRequestSync(ctx, "Players/GetXP v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCreatorCode retrieves creator code information for a specific player.
func (p *PsyNetRPC) GetCreatorCode(ctx context.Context, playerID PlayerID) (*GetCreatorCodeResponse, error) {
	request := GetCreatorCodeRequest{
		PlayerID: playerID,
	}

	var result GetCreatorCodeResponse
	err := p.sendRequestSync(ctx, "Players/GetCreatorCode v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Report reports a player for inappropriate behavior.
func (p *PsyNetRPC) Report(ctx context.Context, reporterID, reportedID PlayerID, reasonID int, description string, gameID, matchID *string) (*ReportResponse, error) {
	request := ReportRequest{
		ReporterID:  reporterID,
		ReportedID:  reportedID,
		ReasonID:    reasonID,
		Description: description,
		GameID:      gameID,
		MatchID:     matchID,
	}

	var result ReportResponse
	err := p.sendRequestSync(ctx, "Players/Report v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
