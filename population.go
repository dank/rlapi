package rlapi

import "context"

type PlaylistID int

type PlaylistPopulation struct {
	PlaylistID PlaylistID `json:"PlaylistID"`
	Population int        `json:"Population"`
}

type GetPopulationResponse struct {
	Playlists []PlaylistPopulation `json:"Playlists"`
	Timestamp int                  `json:"Timestamp"`
}

type UpdatePlayerPlaylistRequest struct {
	Playlist        PlaylistID `json:"Playlist"`
	NumLocalPlayers int        `json:"NumLocalPlayers"`
}

// GetPopulation retrieves current game population (player count) by playlists.
func (p *PsyNetRPC) GetPopulation(ctx context.Context) ([]PlaylistPopulation, error) {
	var result GetPopulationResponse
	err := p.sendRequestSync(ctx, "Population/GetPopulation v1", emptyRequest{}, &result)
	if err != nil {
		return nil, err
	}
	return result.Playlists, nil
}

func (p *PsyNetRPC) UpdatePlayerPlaylist(ctx context.Context, playlistID PlaylistID, numLocalPlayers int) error {
	request := UpdatePlayerPlaylistRequest{
		Playlist:        playlistID,
		NumLocalPlayers: numLocalPlayers,
	}

	var result interface{}
	err := p.sendRequestSync(ctx, "Population/UpdatePlayerPlaylist v1", request, &result)
	if err != nil {
		return err
	}
	return nil
}
