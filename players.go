package rlapi

import "context"

type PlayerData struct {
	PlayerID      string `json:"PlayerID"`
	PlayerName    string `json:"PlayerName"`
	PresenceState string `json:"PresenceState"`
	PresenceInfo  string `json:"PresenceInfo"`
}

type PlayerXPInfo struct {
	TotalXP                  int    `json:"TotalXP"`
	XPLevel                  int    `json:"XPLevel"`
	XPTitle                  string `json:"XPTitle"`
	XPProgressInCurrentLevel int    `json:"XPProgressInCurrentLevel"`
	XPRequiredForNextLevel   int    `json:"XPRequiredForNextLevel"`
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

type GetBanStatusRequest struct {
	Players []PlayerID `json:"Players"`
}

type GetBanStatusResponse struct {
	BanMessages []interface{} `json:"BanMessages"`
}

type GetProfileRequest struct {
	PlayerIDs []PlayerID `json:"PlayerIDs"`
}

type GetProfileResponse struct {
	PlayerData []PlayerData `json:"PlayerData"`
}

type GetXPRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
}

type GetXPResponse struct {
	XPInfoResponse PlayerXPInfo `json:"XPInfoResponse"`
}

type GetCreatorCodeRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
}

type GetCreatorCodeResponse struct {
	CreatorCode interface{} `json:"CreatorCode"`
}

type ReportRequest struct {
	Reports []Report `json:"Reports"`
	GameID  string   `json:"GameID"`
}

type Report struct {
	Reporter        PlayerID `json:"Reporter"`
	Offender        PlayerID `json:"Offender"`
	ReasonIDs       []int    `json:"ReasonIDs"`
	ReportTimestamp float64  `json:"ReportTimestamp"`
}

type ReportResponse struct {
	Success  bool   `json:"Success"`
	ReportID string `json:"ReportID"`
	Message  string `json:"Message"`
}

// GetBanStatus retrieves ban status information for given players.
func (p *PsyNetRPC) GetBanStatus(ctx context.Context, playerIDs []PlayerID) ([]interface{}, error) {
	request := GetBanStatusRequest{
		Players: playerIDs,
	}

	var result GetBanStatusResponse
	err := p.sendRequestSync(ctx, "Players/GetBanStatus v3", request, &result)
	if err != nil {
		return nil, err
	}
	return result.BanMessages, nil
}

// GetProfiles retrieves profile information for given players.
func (p *PsyNetRPC) GetProfiles(ctx context.Context, playerIDs []PlayerID) ([]PlayerData, error) {
	request := GetProfileRequest{
		PlayerIDs: playerIDs,
	}

	var result GetProfileResponse
	err := p.sendRequestSync(ctx, "Players/GetProfile v1", request, &result)
	if err != nil {
		return nil, err
	}
	return result.PlayerData, nil
}

// GetXP retrieves XP information for the authenticated player.
func (p *PsyNetRPC) GetXP(ctx context.Context, playerID PlayerID) (*PlayerXPInfo, error) {
	request := GetXPRequest{
		PlayerID: playerID,
	}

	var result GetXPResponse
	err := p.sendRequestSync(ctx, "Players/GetXP v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result.XPInfoResponse, nil
}

// GetCreatorCode retrieves creator code information for the authenticated player.
func (p *PsyNetRPC) GetCreatorCode(ctx context.Context) (interface{}, error) {
	var result GetCreatorCodeResponse
	err := p.sendRequestSync(ctx, "Players/GetCreatorCode v1", emptyRequest{}, &result)
	if err != nil {
		return nil, err
	}
	return &result.CreatorCode, nil
}

// ReportPlayer reports a player.
func (p *PsyNetRPC) ReportPlayer(ctx context.Context, reports []Report, gameID string) error {
	request := ReportRequest{
		Reports: reports,
		GameID:  gameID,
	}

	var result interface{}
	err := p.sendRequestSync(ctx, "Players/Report v4", request, &result)
	if err != nil {
		return err
	}
	return nil
}
