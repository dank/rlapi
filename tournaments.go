package rlapi

import "context"

type TournamentID string

type Tournament struct {
	// Same as TournamentID but this is returned as an int for some reason
	ID                     int      `json:"ID"`
	Title                  string   `json:"Title"`
	CreatorName            string   `json:"CreatorName"`
	CreatorPlayerID        string   `json:"CreatorPlayerID"`
	StartTime              int      `json:"StartTime"`
	GenerateBracketTime    *int     `json:"GenerateBracketTime"`
	MaxBracketSize         int      `json:"MaxBracketSize"`
	RankMin                int      `json:"RankMin"`
	RankMax                int      `json:"RankMax"`
	Region                 string   `json:"Region"`
	Platforms              []string `json:"Platforms"`
	GameTags               string   `json:"GameTags"`
	GameMode               int      `json:"GameMode"`
	GameModes              []int    `json:"GameModes"`
	TeamSize               int      `json:"TeamSize"`
	MapSetName             *string  `json:"MapSetName"`
	DisabledMaps           []string `json:"DisabledMaps"`
	SeriesLength           int      `json:"SeriesLength"`
	FinalSeriesLength      int      `json:"FinalSeriesLength"`
	SeriesRoundLengths     []int    `json:"SeriesRoundLengths"`
	SeedingType            int      `json:"SeedingType"`
	TieBreaker             int      `json:"TieBreaker"`
	Public                 bool     `json:"bPublic"`
	TeamsRegistered        int      `json:"TeamsRegistered"`
	ScheduleID             *int64   `json:"ScheduleID"`
	IsSchedulingTournament bool     `json:"IsSchedulingTournament"`
}

// TournamentSchedule represents tournament schedule information
type TournamentSchedule struct {
	Time        int          `json:"Time"`
	ScheduleID  int          `json:"ScheduleID"`
	UpdateSkill bool         `json:"bUpdateSkill"`
	Tournaments []Tournament `json:"Tournaments"`
}

// TournamentFormat represents tournament format settings
type TournamentFormat struct {
	Type        string `json:"Type"`
	TeamSize    int    `json:"TeamSize"`
	MaxRounds   int    `json:"MaxRounds"`
	BracketType string `json:"BracketType"`
}

// TournamentRequirements represents requirements to join a tournament
type TournamentRequirements struct {
	MinRank        int     `json:"MinRank"`
	MaxRank        int     `json:"MaxRank"`
	MinLevel       int     `json:"MinLevel"`
	RequiredRegion *string `json:"RequiredRegion"`
}

// TournamentReward represents a reward for tournament participation
type TournamentReward struct {
	Rank         int           `json:"Rank"`
	Products     []Product     `json:"Products"`
	Currency     []interface{} `json:"Currency"`
	TournamentXP int           `json:"TournamentXP"`
}

// TournamentSubscription represents a player's tournament subscription
type TournamentSubscription struct {
	TournamentID TournamentID `json:"TournamentID"`
	PlayerID     PlayerID     `json:"PlayerID"`
	SubscribedAt int          `json:"SubscribedAt"`
	Status       string       `json:"Status"`
}

type TournamentCredentials struct {
	Title    string `json:"Title"`
	Password string `json:"Password"`
}

type GetScheduleRegionRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
}

type GetScheduleRegionResponse struct {
	ScheduleRegion string `json:"ScheduleRegion"`
}

type GetCycleDataRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
}

type GetCycleDataResponse struct {
	CycleID              int           `json:"CycleID"`
	CycleEndTime         int64         `json:"CycleEndTime"`
	WeekID               int           `json:"WeekID"`
	WeekEndTime          int64         `json:"WeekEndTime"`
	WeeklyCurrencies     []interface{} `json:"WeeklyCurrencies"`
	Weeks                []interface{} `json:"Weeks"`
	TournamentCurrencyID int           `json:"TournamentCurrencyID"`
}

type GetScheduleRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
	Region   string   `json:"Region"`
}

type GetScheduleResponse struct {
	Schedules []TournamentSchedule `json:"Schedules"`
}

type GetPublicTournamentsRequest struct {
	PlayerID    PlayerID             `json:"PlayerID"`
	Search      TournamentSearchInfo `json:"Search"`
	TeamMembers []PlayerID           `json:"TeamMembers"`
}

type TournamentSearchInfo struct {
	Text               string   `json:"Text"`
	RankMin            int      `json:"RankMin"`
	RankMax            int      `json:"RankMax"`
	GameModes          []int    `json:"GameModes"`
	Regions            []string `json:"Regions"`
	TeamSize           int      `json:"TeamSize"`
	BracketSize        int      `json:"BracketSize"`
	EnableCrossplay    bool     `json:"bEnableCrossplay"`
	StartTime          string   `json:"StartTime"`
	EndTime            string   `json:"EndTime"`
	ShowFull           bool     `json:"bShowFull"`
	ShowIneligibleRank bool     `json:"bShowIneligibleRank"`
}

type GetPublicTournamentsResponse struct {
	Tournaments []Tournament `json:"Tournaments"`
}

type RegisterTournamentRequest struct {
	PlayerID     PlayerID              `json:"PlayerID"`
	TournamentID TournamentID          `json:"TournamentID"`
	Credentials  TournamentCredentials `json:"Credentials,omitempty"`
}

type RegisterTournamentResponse struct {
	Tournament Tournament `json:"Tournament"`
}

type UnsubscribeTournamentRequest struct {
	PlayerID                           PlayerID     `json:"PlayerID"`
	TournamentID                       TournamentID `json:"TournamentID"`
	UnsubscribeAnyRegisteredTournament bool         `json:"bUnsubscribeAnyRegisteredTournament"`
	TeamMembers                        []PlayerID   `json:"TeamMembers"`
}

type GetTournamentSubscriptionsRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
}

// GetTournamentScheduleRegion retrieves the tournament schedule region for the authenticated player.
func (p *PsyNetRPC) GetTournamentScheduleRegion(ctx context.Context) (string, error) {
	request := GetScheduleRegionRequest{
		PlayerID: p.localPlayerID,
	}

	var result GetScheduleRegionResponse
	err := p.sendRequestSync(ctx, "Tournaments/Status/GetScheduleRegion v1", request, &result)
	if err != nil {
		return "", err
	}
	return result.ScheduleRegion, nil
}

// GetTournamentCycleData retrieves tournament cycle data for the authenticated player.
func (p *PsyNetRPC) GetTournamentCycleData(ctx context.Context) (*GetCycleDataResponse, error) {
	request := GetCycleDataRequest{
		PlayerID: p.localPlayerID,
	}

	var result GetCycleDataResponse
	err := p.sendRequestSync(ctx, "Tournaments/Status/GetCycleData v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTournamentSchedule retrieves the tournament schedule for a given region.
func (p *PsyNetRPC) GetTournamentSchedule(ctx context.Context, region string) ([]TournamentSchedule, error) {
	request := GetScheduleRequest{
		PlayerID: p.localPlayerID,
		Region:   region,
	}

	var result GetScheduleResponse
	err := p.sendRequestSync(ctx, "Tournaments/Search/GetSchedule v1", request, &result)
	if err != nil {
		return nil, err
	}
	return result.Schedules, nil
}

// GetPublicTournaments retrieves a list of public tournaments for the authenticated player.
func (p *PsyNetRPC) GetPublicTournaments(ctx context.Context, searchInfo TournamentSearchInfo, teamMembers []PlayerID) ([]Tournament, error) {
	request := GetPublicTournamentsRequest{
		PlayerID:    p.localPlayerID,
		Search:      searchInfo,
		TeamMembers: teamMembers,
	}

	var result GetPublicTournamentsResponse
	err := p.sendRequestSync(ctx, "Tournaments/Search/GetPublicTournaments v1", request, &result)
	if err != nil {
		return nil, err
	}
	return result.Tournaments, nil
}

// RegisterTournament registers the authenticated player for a tournament.
func (p *PsyNetRPC) RegisterTournament(ctx context.Context, tournamentID TournamentID, credentials TournamentCredentials) (*Tournament, error) {
	request := RegisterTournamentRequest{
		PlayerID:     p.localPlayerID,
		TournamentID: tournamentID,
		Credentials:  credentials,
	}

	var result RegisterTournamentResponse
	err := p.sendRequestSync(ctx, "Tournaments/Registration/RegisterTournament v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result.Tournament, nil
}

// UnsubscribeTournament unsubscribes the authenticated player from a tournament.
func (p *PsyNetRPC) UnsubscribeTournament(ctx context.Context, tournamentID TournamentID, unsubscribeAnyRegisteredTournament bool, teamMembers []PlayerID) error {
	request := UnsubscribeTournamentRequest{
		PlayerID:                           p.localPlayerID,
		TournamentID:                       tournamentID,
		UnsubscribeAnyRegisteredTournament: unsubscribeAnyRegisteredTournament,
		TeamMembers:                        teamMembers,
	}

	var result interface{}
	err := p.sendRequestSync(ctx, "Tournaments/Registration/UnsubscribeTournament v1", request, &result)
	if err != nil {
		return err
	}
	return nil
}

// GetTournamentSubscriptions retrieves the authenticated player's tournament subscriptions.
func (p *PsyNetRPC) GetTournamentSubscriptions(ctx context.Context) (interface{}, error) {
	request := GetTournamentSubscriptionsRequest{
		PlayerID: p.localPlayerID,
	}

	var result interface{}
	err := p.sendRequestSync(ctx, "Tournaments/Status/GetTournamentSubscriptions v1", request, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
