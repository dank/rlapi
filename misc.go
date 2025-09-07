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

type Region struct {
	Region     string   `json:"Region"`
	Label      string   `json:"Label"`
	SubRegions []string `json:"SubRegions"`
}

type GetSubRegionsRequest struct {
	RequestRegions []interface{} `json:"RequestRegions"`
	Regions        []interface{} `json:"Regions"`
}

type GetSubRegionsResponse struct {
	Regions []Region `json:"Regions"`
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

type JoinMatchRequest struct {
	JoinType   string `json:"JoinType"`
	ServerName string `json:"ServerName"`
	Password   string `json:"Password"`
}

type FilterContentRequest struct {
	Content []string `json:"Content"`
	Policy  string   `json:"Policy"`
}

type FilterContentResponse struct {
	FilteredContent []string `json:"FilteredContent"`
}

type RecordMetricsRequest struct {
	AppSessionID       string  `json:"AppSessionID"`
	LevelSessionID     string  `json:"LevelSessionID"`
	CurrentTimeSeconds float64 `json:"CurrentTimeSeconds"`
	FirstEventIndex    int     `json:"FirstEventIndex"`
	Events             []struct {
		PlayerID    PlayerID `json:"PlayerID,omitempty"`
		TimeSeconds float64  `json:"TimeSeconds"`
		Version     int      `json:"Version"`
		EventName   string   `json:"EventName"`
		EventData   string   `json:"EventData"`
	} `json:"Events"`
}

type GetTradeInFiltersResponse struct {
	TradeInFilters []TradeInFilter `json:"TradeInFilters"`
}

type RelayToServerRequest struct {
	DSConnectToken string `json:"DSConnectToken"`
	ReservationID  string `json:"ReservationID"`
	MessageType    string `json:"MessageType"`
	MessagePayload struct {
		Settings struct {
			MatchType        int    `json:"MatchType"`
			PlaylistID       int    `json:"PlaylistID"`
			BFriendJoin      bool   `json:"bFriendJoin"`
			BMigration       bool   `json:"bMigration"`
			BRankedReconnect bool   `json:"bRankedReconnect"`
			Password         string `json:"Password"`
		} `json:"Settings"`
		MapPrefs []struct {
			PlayerID    string `json:"PlayerID"`
			MapLikes    []int  `json:"MapLikes"`
			MapDislikes []int  `json:"MapDislikes"`
		} `json:"MapPrefs"`
		Players []struct {
			PlayerID      string  `json:"PlayerID"`
			PlayerName    string  `json:"PlayerName"`
			SkillMu       float64 `json:"SkillMu"`
			SkillSigma    float64 `json:"SkillSigma"`
			Tier          int     `json:"Tier"`
			BRemotePlayer bool    `json:"bRemotePlayer"`
			Loadout       []int   `json:"Loadout"`
			MapLikes      []int   `json:"MapLikes"`
			MapDislikes   []int   `json:"MapDislikes"`
			ClubID        int     `json:"ClubID"`
		} `json:"Players"`
		PartyLeaderID     string `json:"PartyLeaderID"`
		ReservationID     string `json:"ReservationID"`
		BDisableCrossPlay bool   `json:"bDisableCrossPlay"`
	} `json:"MessagePayload"`
}

type CanShowAvatarRequest struct {
	PlayerIDs []PlayerID `json:"PlayerIDs"`
}

type CanShowAvatarResponse struct {
	AllowedPlayerIDs []PlayerID `json:"AllowedPlayerIDs"`
	HiddenPlayerIDs  []PlayerID `json:"HiddenPlayerIDs"`
}

// GetSubRegions retrieves available server regions.
func (p *PsyNetRPC) GetSubRegions(ctx context.Context) ([]Region, error) {
	request := GetSubRegionsRequest{
		RequestRegions: []interface{}{},
		Regions:        []interface{}{},
	}

	var result GetSubRegionsResponse
	err := p.sendRequestSync(ctx, "Regions/GetSubRegions v1", request, &result)
	if err != nil {
		return nil, err
	}
	return result.Regions, nil
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

// JoinMatch joins a private match.
func (p *PsyNetRPC) JoinMatch(ctx context.Context, joinType, serverName, password string) (interface{}, error) {
	request := JoinMatchRequest{
		JoinType:   joinType,
		ServerName: serverName,
		Password:   password,
	}

	var result interface{}
	err := p.sendRequestSync(ctx, "Reservations/JoinMatch v1", request, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
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

func (p *PsyNetRPC) RecordMetrics(ctx context.Context, request *RecordMetricsRequest) error {
	var result interface{}
	err := p.sendRequestSync(ctx, "Metrics/RecordMetrics v1", request, &result)
	if err != nil {
		return err
	}
	return nil
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

func (p *PsyNetRPC) RelayToServer(ctx context.Context, request *RelayToServerRequest) error {
	var result interface{}
	err := p.sendRequestSync(ctx, "DSR/RelayToServer v1", request, &result)
	if err != nil {
		return err
	}
	return nil
}

// CanShowAvatar checks if players' avatars can be shown.
func (p *PsyNetRPC) CanShowAvatar(ctx context.Context, playerIDs []PlayerID) (*CanShowAvatarResponse, error) {
	request := CanShowAvatarRequest{
		PlayerIDs: playerIDs,
	}

	var result CanShowAvatarResponse
	err := p.sendRequestSync(ctx, "Users/CanShowAvatar v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
