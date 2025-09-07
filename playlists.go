package rlapi

import "context"

// Playlist represents a game playlist
type Playlist struct {
	NodeID    string `json:"NodeID"`
	Playlist  int    `json:"Playlist"`
	Type      int    `json:"Type"`
	StartTime *int   `json:"StartTime"`
	EndTime   *int   `json:"EndTime"`
}

// ActivePlaylists represents all active playlists
type ActivePlaylists struct {
	CasualPlaylists []Playlist `json:"CasualPlaylists"`
	RankedPlaylists []Playlist `json:"RankedPlaylists"`
	XPLevelUnlocked int        `json:"XPLevelUnlocked"`
}

type GetActivePlaylistsResponse struct {
	CasualPlaylists []Playlist `json:"CasualPlaylists"`
	RankedPlaylists []Playlist `json:"RankedPlaylists"`
	XPLevelUnlocked int        `json:"XPLevelUnlocked"`
}

// GetActivePlaylists retrieves all currently active playlists.
func (p *PsyNetRPC) GetActivePlaylists(ctx context.Context) (*GetActivePlaylistsResponse, error) {
	var result GetActivePlaylistsResponse
	err := p.sendRequestSync(ctx, "Playlists/GetActivePlaylists v1", emptyRequest{}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
