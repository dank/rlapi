package rlapi

import "context"

// MatchmakingSettings represents matchmaking configuration
type MatchmakingSettings struct {
	Playlists         []int  `json:"Playlists"`
	Region            string `json:"Region"`
	PartyID           string `json:"PartyID"`
	BDisableCrossplay bool   `json:"bDisableCrossplay"`
}

// PrivateMatchSearch represents search parameters for private matches
type PrivateMatchSearch struct {
	Name     string `json:"Name"`
	Password string `json:"Password"`
	Region   string `json:"Region"`
}

// MatchmakingResult represents the result of a matchmaking request
type MatchmakingResult struct {
	MatchID           string `json:"MatchID"`
	Status            string `json:"Status"`
	EstimatedWaitTime int    `json:"EstimatedWaitTime"`
}

// Request and Response types

type StartMatchmakingRequest struct {
	Playlists         []int  `json:"Playlists"`
	Region            string `json:"Region"`
	PartyID           string `json:"PartyID"`
	BDisableCrossplay bool   `json:"bDisableCrossplay"`
}

type StartMatchmakingResponse struct {
	Success           bool   `json:"Success"`
	MatchmakingID     string `json:"MatchmakingID"`
	EstimatedWaitTime int    `json:"EstimatedWaitTime"`
}

type PlayerCancelMatchmakingRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
}

type PlayerCancelMatchmakingResponse struct {
	Success bool `json:"Success"`
}

type PlayerSearchPrivateMatchRequest struct {
	Name     string `json:"Name"`
	Password string `json:"Password"`
	Region   string `json:"Region"`
}

type PlayerSearchPrivateMatchResponse struct {
	Matches []PrivateMatch `json:"Matches"`
}

type PrivateMatch struct {
	MatchID     string                 `json:"MatchID"`
	Name        string                 `json:"Name"`
	Region      string                 `json:"Region"`
	PlayerCount int                    `json:"PlayerCount"`
	MaxPlayers  int                    `json:"MaxPlayers"`
	Settings    map[string]interface{} `json:"Settings"`
}

// StartMatchmaking starts matchmaking for specified playlists.
func (p *PsyNetRPC) StartMatchmaking(ctx context.Context, playlists []int, region, partyID string, disableCrossplay bool) (*StartMatchmakingResponse, error) {
	request := StartMatchmakingRequest{
		Playlists:         playlists,
		Region:            region,
		PartyID:           partyID,
		BDisableCrossplay: disableCrossplay,
	}

	var result StartMatchmakingResponse
	err := p.sendRequestSync(ctx, "Matchmaking/StartMatchmaking v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// PlayerCancelMatchmaking cancels ongoing matchmaking for a player.
func (p *PsyNetRPC) PlayerCancelMatchmaking(ctx context.Context, playerID PlayerID) (*PlayerCancelMatchmakingResponse, error) {
	request := PlayerCancelMatchmakingRequest{
		PlayerID: playerID,
	}

	var result PlayerCancelMatchmakingResponse
	err := p.sendRequestSync(ctx, "Matchmaking/PlayerCancelMatchmaking v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// PlayerSearchPrivateMatch searches for private matches.
func (p *PsyNetRPC) PlayerSearchPrivateMatch(ctx context.Context, name, password, region string) (*PlayerSearchPrivateMatchResponse, error) {
	request := PlayerSearchPrivateMatchRequest{
		Name:     name,
		Password: password,
		Region:   region,
	}

	var result PlayerSearchPrivateMatchResponse
	err := p.sendRequestSync(ctx, "Matchmaking/PlayerSearchPrivateMatch v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
