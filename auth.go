package rlapi

import "fmt"

type authPlayerRequest struct {
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

type authPlayerResponse struct {
	IsLastChanceAuthBan bool     `json:"IsLastChanceAuthBan"`
	SessionID           string   `json:"SessionID"`
	VerifiedPlayerName  string   `json:"VerifiedPlayerName"`
	UseWebSocket        bool     `json:"UseWebSocket"`
	PerConURL           string   `json:"PerConURL"`
	PerConURLv2         string   `json:"PerConURLv2"`
	PsyToken            string   `json:"PsyToken"`
	CountryRestrictions []string `json:"CountryRestrictions"`
}

// AuthPlayer authenticates with PsyNet and returns a WebSocket connection.
func (p *PsyNet) AuthPlayer(platform Platform, authToken string, accountID string, accountName string) (*PsyNetRPC, error) {
	localPlayerId := fmt.Sprintf("%s|%s|0", platform, accountID)
	req := &authPlayerRequest{
		Platform:            string(platform),
		PlayerName:          accountName,
		PlayerID:            accountID,
		Language:            "INT",
		AuthTicket:          authToken,
		BuildRegion:         "",
		FeatureSet:          featureSet,
		Device:              "PC",
		LocalFirstPlayerID:  localPlayerId,
		SkipAuth:            false,
		SetAsPrimaryAccount: true,
		EpicAuthTicket:      authToken,
		EpicAccountID:       accountID,
	}

	var res authPlayerResponse
	err := p.postJSON([]string{"Auth", "AuthPlayer", "v2"}, req, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate player: %w", err)
	}

	wsConn, err := p.establishSocket(res.PerConURLv2, res.PsyToken, res.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to establish websocket: %w", err)
	}

	rpc := newPsyNetRPC(wsConn, p.logger)
	go rpc.readMessages()
	go rpc.pingHandler()

	return rpc, nil
}
