package rlapi

import "fmt"

type AuthPlayerRequest struct {
	Platform            string `json:"Platform"`
	PlayerName          string `json:"PlayerName"`
	PlayerID            string `json:"PlayerID"`
	Language            string `json:"Language"`
	AuthTicket          string `json:"AuthTicket"`
	BuildRegion         string `json:"BuildRegion"`
	FeatureSet          string `json:"FeatureSet"`
	Device              string `json:"Device"`
	LocalFirstPlayerID  string `json:"LocalFirstPlayerID"`
	SkipAuth            bool   `json:"bSkipAuth"`
	SetAsPrimaryAccount bool   `json:"bSetAsPrimaryAccount"`
	EpicAuthTicket      string `json:"EpicAuthTicket"`
	EpicAccountID       string `json:"EpicAccountID"`
}

type AuthPlayerResponse struct {
	IsLastChanceAuthBan bool     `json:"IsLastChanceAuthBan"`
	SessionID           string   `json:"SessionID"`
	VerifiedPlayerName  string   `json:"VerifiedPlayerName"`
	UseWebSocket        bool     `json:"UseWebSocket"`
	PerConURL           string   `json:"PerConURL"`
	PerConURLv2         string   `json:"PerConURLv2"`
	PsyToken            string   `json:"PsyToken"`
	CountryRestrictions []string `json:"CountryRestrictions"`
}

// AuthPlayer authenticates with PsyNet via EGS and returns a WebSocket connection.
func (p *PsyNet) AuthPlayer(authToken string, accountID string, accountName string) (*PsyNetRPC, error) {
	localPlayerId := NewPlayerID(PlatformEpic, accountID)
	req := &AuthPlayerRequest{
		Platform:            string(PlatformEpic),
		PlayerName:          accountName,
		PlayerID:            accountID,
		Language:            "INT",
		AuthTicket:          authToken,
		BuildRegion:         "",
		FeatureSet:          featureSet,
		Device:              "PC",
		LocalFirstPlayerID:  localPlayerId.String(),
		SkipAuth:            false,
		SetAsPrimaryAccount: true,
		EpicAuthTicket:      authToken,
		EpicAccountID:       accountID,
	}

	var res AuthPlayerResponse
	err := p.postJSON([]string{"Auth", "AuthPlayer", "v2"}, req, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate player: %w", err)
	}

	rpc, err := p.establishSocket(res.PerConURLv2, localPlayerId, res.PsyToken, res.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to establish websocket: %w", err)
	}

	go rpc.readMessages()
	rpc.schedulePing()

	return rpc, nil
}

// AuthPlayerSteam authenticates with PsyNet via Steam session ticket and returns a WebSocket connection.
func (p *PsyNet) AuthPlayerSteam(authToken string, epicAccountID string, steamAccountID string, accountName string) (*PsyNetRPC, error) {
	localPlayerId := NewPlayerID(PlatformSteam, steamAccountID)
	req := &AuthPlayerRequest{
		Platform:            string(PlatformSteam),
		PlayerName:          accountName,
		PlayerID:            steamAccountID,
		Language:            "INT",
		AuthTicket:          authToken,
		BuildRegion:         "",
		FeatureSet:          featureSet,
		Device:              "PC",
		SkipAuth:            false,
		SetAsPrimaryAccount: true,
		EpicAuthTicket:      authToken,
		EpicAccountID:       epicAccountID,
	}

	var res AuthPlayerResponse
	err := p.postJSON([]string{"Auth", "AuthPlayer", "v2"}, req, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate player: %w", err)
	}

	rpc, err := p.establishSocket(res.PerConURLv2, localPlayerId, res.PsyToken, res.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to establish websocket: %w", err)
	}

	go rpc.readMessages()
	rpc.schedulePing()

	return rpc, nil
}
