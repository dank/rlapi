package rlapi

import "context"

type TradeInFilter struct {
	ID               int      `json:"ID"`
	Label            string   `json:"Label"`
	SeriesIDs        []int    `json:"SeriesIDs"`
	Blueprint        bool     `json:"bBlueprint"`
	TradeInQualities []string `json:"TradeInQualities"`
}

type Server struct {
	Region    string `json:"Region"`
	Host      string `json:"Host"`
	Port      string `json:"Port"`
	SubRegion string `json:"SubRegion"`
}

type PopulationData struct {
	Playlists []PlaylistPopulation `json:"Playlists"`
	Timestamp int64                `json:"Timestamp"`
}

type PlaylistPopulation struct {
	PlaylistID int `json:"PlaylistID"`
	Population int `json:"Population"`
}

type SubRegion struct {
	RegionID   string   `json:"RegionID"`
	RegionName string   `json:"RegionName"`
	SubRegions []string `json:"SubRegions"`
}

type TrainingData struct {
	Packs []TrainingPack `json:"Packs"`
}

type TrainingPack struct {
	PackID      string `json:"PackID"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	CreatorID   string `json:"CreatorID"`
	Difficulty  int    `json:"Difficulty"`
}

type TrainingMetadata struct {
	Categories []TrainingCategory `json:"Categories"`
}

type TrainingCategory struct {
	CategoryID   int    `json:"CategoryID"`
	CategoryName string `json:"CategoryName"`
}

type GetPopulationResponse struct {
	Playlists []PlaylistPopulation `json:"Playlists"`
	Timestamp int64                `json:"Timestamp"`
}

type UpdatePlayerPlaylistRequest struct {
	PlayerID   PlayerID `json:"PlayerID"`
	PlaylistID int      `json:"PlaylistID"`
}

type UpdatePlayerPlaylistResponse struct {
	Success bool `json:"Success"`
}

type GetSubRegionsResponse struct {
	Regions []SubRegion `json:"Regions"`
}

type GetGameServerPingListRequest struct {
	Regions []interface{} `json:"Regions"`
}

type GetGameServerPingListResponse struct {
	Servers []Server `json:"Servers"`
}

type GetClubPrivateMatchesResponse struct {
	Servers []Server `json:"Servers"`
}

type BrowseTrainingDataRequest struct {
	CategoryID *int    `json:"CategoryID,omitempty"`
	Search     *string `json:"Search,omitempty"`
	Limit      int     `json:"Limit"`
	Offset     int     `json:"Offset"`
}

type BrowseTrainingDataResponse struct {
	Packs []TrainingPack `json:"Packs"`
	Total int            `json:"Total"`
}

type GetTrainingMetadataResponse struct {
	Categories []TrainingCategory `json:"Categories"`
}

type JoinMatchRequest struct {
	MatchID  string   `json:"MatchID"`
	PlayerID PlayerID `json:"PlayerID"`
}

type JoinMatchResponse struct {
	Success    bool   `json:"Success"`
	ServerInfo string `json:"ServerInfo"`
}

type FilterContentRequest struct {
	Content []string `json:"Content"`
	Policy  string   `json:"Policy"`
}

type FilterContentResponse struct {
	FilteredContent []string `json:"FilteredContent"`
}

type RecordMetricsRequest struct {
	PlayerID PlayerID               `json:"PlayerID"`
	Metrics  map[string]interface{} `json:"Metrics"`
}

type RecordMetricsResponse struct {
	Success bool `json:"Success"`
}

type GetTradeInFiltersResponse struct {
	TradeInFilters []TradeInFilter `json:"TradeInFilters"`
}

type RelayToServerRequest struct {
	ServerID string      `json:"ServerID"`
	Message  interface{} `json:"Message"`
}

type RelayToServerResponse struct {
	Success  bool        `json:"Success"`
	Response interface{} `json:"Response"`
}

type CanShowAvatarRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
}

type CanShowAvatarResponse struct {
	CanShow bool `json:"CanShow"`
}

// GetPopulation retrieves current playlist population data.
func (p *PsyNetRPC) GetPopulation(ctx context.Context) (*GetPopulationResponse, error) {
	var result GetPopulationResponse
	err := p.sendRequestSync(ctx, "Population/GetPopulation v1", emptyRequest{}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdatePlayerPlaylist updates a player's playlist preference.
func (p *PsyNetRPC) UpdatePlayerPlaylist(ctx context.Context, playerID PlayerID, playlistID int) (*UpdatePlayerPlaylistResponse, error) {
	request := UpdatePlayerPlaylistRequest{
		PlayerID:   playerID,
		PlaylistID: playlistID,
	}

	var result UpdatePlayerPlaylistResponse
	err := p.sendRequestSync(ctx, "Population/UpdatePlayerPlaylist v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSubRegions retrieves available sub-regions.
func (p *PsyNetRPC) GetSubRegions(ctx context.Context) (*GetSubRegionsResponse, error) {
	var result GetSubRegionsResponse
	err := p.sendRequestSync(ctx, "Regions/GetSubRegions v1", emptyRequest{}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetGameServerPingList retrieves ping information for game servers.
func (p *PsyNetRPC) GetGameServerPingList(ctx context.Context) ([]Server, error) {
	var result GetGameServerPingListResponse
	err := p.sendRequestSync(ctx, "GameServer/GetGameServerPingList v2", GetGameServerPingListRequest{}, &result)
	if err != nil {
		return nil, err
	}
	return result.Servers, nil
}

// GetClubPrivateMatches retrieves private matches for a club.
func (p *PsyNetRPC) GetClubPrivateMatches(ctx context.Context) ([]Server, error) {
	var result GetClubPrivateMatchesResponse
	err := p.sendRequestSync(ctx, "GameServer/GetClubPrivateMatches v1", emptyRequest{}, &result)
	if err != nil {
		return nil, err
	}
	return result.Servers, nil
}

// BrowseTrainingData searches for training packs.
func (p *PsyNetRPC) BrowseTrainingData(ctx context.Context, categoryID *int, search *string, limit, offset int) (*BrowseTrainingDataResponse, error) {
	request := BrowseTrainingDataRequest{
		CategoryID: categoryID,
		Search:     search,
		Limit:      limit,
		Offset:     offset,
	}

	var result BrowseTrainingDataResponse
	err := p.sendRequestSync(ctx, "Training/BrowseTrainingData v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTrainingMetadata retrieves training metadata including categories.
func (p *PsyNetRPC) GetTrainingMetadata(ctx context.Context) (*GetTrainingMetadataResponse, error) {
	var result GetTrainingMetadataResponse
	err := p.sendRequestSync(ctx, "Training/GetTrainingMetadata v1", emptyRequest{}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// JoinMatch joins a specific match.
func (p *PsyNetRPC) JoinMatch(ctx context.Context, matchID string, playerID PlayerID) (*JoinMatchResponse, error) {
	request := JoinMatchRequest{
		MatchID:  matchID,
		PlayerID: playerID,
	}

	var result JoinMatchResponse
	err := p.sendRequestSync(ctx, "Reservations/JoinMatch v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *PsyNetRPC) FilterContent(ctx context.Context, content []string, policy string) ([]string, error) {
	request := FilterContentRequest{
		Content: content,
		Policy:  policy,
	}

	var result FilterContentResponse
	err := p.sendRequestSync(ctx, "Filters/FilterContent v1", request, &result)
	if err != nil {
		return nil, err
	}
	return result.FilteredContent, nil
}

// RecordMetrics records player metrics.
func (p *PsyNetRPC) RecordMetrics(ctx context.Context, playerID PlayerID, metrics map[string]interface{}) (*RecordMetricsResponse, error) {
	request := RecordMetricsRequest{
		PlayerID: playerID,
		Metrics:  metrics,
	}

	var result RecordMetricsResponse
	err := p.sendRequestSync(ctx, "Metrics/RecordMetrics v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTradeInFilters retrieves trade-in filters.
func (p *PsyNetRPC) GetTradeInFilters(ctx context.Context) ([]TradeInFilter, error) {
	var result GetTradeInFiltersResponse
	err := p.sendRequestSync(ctx, "Drop/GetTradeInFilters v1", emptyRequest{}, &result)
	if err != nil {
		return nil, err
	}
	return result.TradeInFilters, nil
}

// RelayToServer relays a message to a game server.
func (p *PsyNetRPC) RelayToServer(ctx context.Context, serverID string, message interface{}) (*RelayToServerResponse, error) {
	request := RelayToServerRequest{
		ServerID: serverID,
		Message:  message,
	}

	var result RelayToServerResponse
	err := p.sendRequestSync(ctx, "DSR/RelayToServer v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CanShowAvatar checks if a player's avatar can be shown.
func (p *PsyNetRPC) CanShowAvatar(ctx context.Context, playerID PlayerID) (*CanShowAvatarResponse, error) {
	request := CanShowAvatarRequest{
		PlayerID: playerID,
	}

	var result CanShowAvatarResponse
	err := p.sendRequestSync(ctx, "Users/CanShowAvatar v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
