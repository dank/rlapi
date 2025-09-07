package rlapi

import "context"

// MatchmakingSettings represents matchmaking configuration.
type MatchmakingSettings struct {
	Playlists         []int  `json:"Playlists"`
	Region            string `json:"Region"`
	PartyID           string `json:"PartyID"`
	BDisableCrossplay bool   `json:"bDisableCrossplay"`
}

// PrivateMatchSearch represents search parameters for private matches.
type PrivateMatchSearch struct {
	Name     string `json:"Name"`
	Password string `json:"Password"`
	Region   string `json:"Region"`
}

type MatchmakingRegion struct {
	Name string `json:"Name"`
	Ping int    `json:"Ping"`
}

type StartMatchmakingRequest struct {
	Regions          []MatchmakingRegion `json:"Regions"`
	Playlists        []int               `json:"Playlists"`
	SecondsSearching int                 `json:"SecondsSearching"`
	CurrentServerID  string              `json:"CurrentServerID"`
	DisableCrossplay bool                `json:"bDisableCrossplay"`
	PartyID          PartyID             `json:"PartyID"`
	PartyMembers     []PlayerID          `json:"PartyMembers"`
}

type StartMatchmakingResponse struct {
	EstimatedQueueTime int `json:"EstimatedQueueTime"`
}

type PlayerSearchPrivateMatchRequest struct {
	Region     string     `json:"Region"`
	PlaylistID PlaylistID `json:"PlaylistID"`
}

// StartMatchmaking starts matchmaking for given playlists, returns the estimated queue time.
func (p *PsyNetRPC) StartMatchmaking(ctx context.Context, playlists []int, region []MatchmakingRegion, disableCrossplay bool, partyID PartyID, partyMembers []PlayerID) (int, error) {
	request := StartMatchmakingRequest{
		Regions:          region,
		Playlists:        playlists,
		SecondsSearching: 1,
		CurrentServerID:  "",
		DisableCrossplay: disableCrossplay,
		PartyID:          partyID,
		PartyMembers:     partyMembers,
	}

	var result StartMatchmakingResponse
	err := p.sendRequestSync(ctx, "Matchmaking/StartMatchmaking v2", request, &result)
	if err != nil {
		return 0, err
	}
	return result.EstimatedQueueTime, nil
}

// PlayerCancelMatchmaking cancels ongoing matchmaking for the authenticated player.
func (p *PsyNetRPC) PlayerCancelMatchmaking(ctx context.Context) error {
	var result interface{}
	err := p.sendRequestSync(ctx, "Matchmaking/PlayerCancelMatchmaking v1", emptyRequest{}, &result)
	if err != nil {
		return err
	}
	return nil
}

// PlayerSearchPrivateMatch searches for private matches.
func (p *PsyNetRPC) PlayerSearchPrivateMatch(ctx context.Context, region string, playlistID PlaylistID) error {
	request := PlayerSearchPrivateMatchRequest{
		Region:     region,
		PlaylistID: playlistID,
	}

	var result interface{}
	err := p.sendRequestSync(ctx, "Matchmaking/PlayerSearchPrivateMatch v1", request, &result)
	if err != nil {
		return err
	}
	return nil
}
