package rlapi

import "context"

// TournamentCycleData represents tournament cycle information
type TournamentCycleData struct {
	CycleID              int              `json:"CycleID"`
	CycleEndTime         int64            `json:"CycleEndTime"`
	WeekID               int              `json:"WeekID"`
	WeekEndTime          int64            `json:"WeekEndTime"`
	WeeklyCurrencies     []interface{}    `json:"WeeklyCurrencies"`
	Weeks                []TournamentWeek `json:"Weeks"`
	TournamentCurrencyID int              `json:"TournamentCurrencyID"`
}

// TournamentWeek represents a week within a tournament cycle
type TournamentWeek struct {
	Results []TournamentResult `json:"Results"`
}

// TournamentResult represents a tournament result
type TournamentResult struct {
	TournamentID string        `json:"TournamentID"`
	Rank         int           `json:"Rank"`
	Points       int           `json:"Points"`
	Rewards      []interface{} `json:"Rewards"`
}

// TournamentSchedule represents tournament schedule information
type TournamentSchedule struct {
	Tournaments []ScheduledTournament `json:"Tournaments"`
	Region      string                `json:"Region"`
}

// ScheduledTournament represents a scheduled tournament
type ScheduledTournament struct {
	TournamentID        string                 `json:"TournamentID"`
	Name                string                 `json:"Name"`
	Description         string                 `json:"Description"`
	StartTime           int64                  `json:"StartTime"`
	EndTime             int64                  `json:"EndTime"`
	RegistrationEnd     int64                  `json:"RegistrationEnd"`
	MaxParticipants     int                    `json:"MaxParticipants"`
	CurrentParticipants int                    `json:"CurrentParticipants"`
	Format              TournamentFormat       `json:"Format"`
	Requirements        TournamentRequirements `json:"Requirements"`
	Rewards             []TournamentReward     `json:"Rewards"`
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
	Products     []ProductData `json:"Products"`
	Currency     []interface{} `json:"Currency"`
	TournamentXP int           `json:"TournamentXP"`
}

// TournamentSubscription represents a player's tournament subscription
type TournamentSubscription struct {
	TournamentID string   `json:"TournamentID"`
	PlayerID     PlayerID `json:"PlayerID"`
	SubscribedAt int64    `json:"SubscribedAt"`
	Status       string   `json:"Status"`
}

// PublicTournament represents a public tournament listing
type PublicTournament struct {
	TournamentID     string                 `json:"TournamentID"`
	Name             string                 `json:"Name"`
	Description      string                 `json:"Description"`
	StartTime        int64                  `json:"StartTime"`
	EndTime          int64                  `json:"EndTime"`
	Format           TournamentFormat       `json:"Format"`
	Requirements     TournamentRequirements `json:"Requirements"`
	IsPublic         bool                   `json:"IsPublic"`
	ParticipantCount int                    `json:"ParticipantCount"`
	MaxParticipants  int                    `json:"MaxParticipants"`
}

// Request and Response types

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
	CycleID              int              `json:"CycleID"`
	CycleEndTime         int64            `json:"CycleEndTime"`
	WeekID               int              `json:"WeekID"`
	WeekEndTime          int64            `json:"WeekEndTime"`
	WeeklyCurrencies     []interface{}    `json:"WeeklyCurrencies"`
	Weeks                []TournamentWeek `json:"Weeks"`
	TournamentCurrencyID int              `json:"TournamentCurrencyID"`
}

type GetScheduleRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
	Region   string   `json:"Region"`
}

type GetScheduleResponse struct {
	Schedule TournamentSchedule `json:"Schedule"`
}

type GetPublicTournamentsRequest struct {
	Region string `json:"Region"`
	Limit  int    `json:"Limit"`
	Offset int    `json:"Offset"`
}

type GetPublicTournamentsResponse struct {
	Tournaments []PublicTournament `json:"Tournaments"`
	Total       int                `json:"Total"`
}

type RegisterTournamentRequest struct {
	PlayerID     PlayerID `json:"PlayerID"`
	TournamentID string   `json:"TournamentID"`
	TeamID       *string  `json:"TeamID,omitempty"`
}

type RegisterTournamentResponse struct {
	Success        bool   `json:"Success"`
	RegistrationID string `json:"RegistrationID"`
}

type UnsubscribeTournamentRequest struct {
	PlayerID     PlayerID `json:"PlayerID"`
	TournamentID string   `json:"TournamentID"`
}

type UnsubscribeTournamentResponse struct {
	Success bool `json:"Success"`
}

type GetTournamentSubscriptionsRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
}

type GetTournamentSubscriptionsResponse struct {
	Subscriptions []TournamentSubscription `json:"Subscriptions"`
}

// GetScheduleRegion retrieves the tournament schedule region for a player.
func (p *PsyNetRPC) GetScheduleRegion(ctx context.Context, playerID PlayerID) (*GetScheduleRegionResponse, error) {
	request := GetScheduleRegionRequest{
		PlayerID: playerID,
	}

	var result GetScheduleRegionResponse
	err := p.sendRequestSync(ctx, "Tournaments/Status/GetScheduleRegion v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCycleData retrieves tournament cycle data for a player.
func (p *PsyNetRPC) GetCycleData(ctx context.Context, playerID PlayerID) (*GetCycleDataResponse, error) {
	request := GetCycleDataRequest{
		PlayerID: playerID,
	}

	var result GetCycleDataResponse
	err := p.sendRequestSync(ctx, "Tournaments/Status/GetCycleData v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSchedule retrieves the tournament schedule for a specific region.
func (p *PsyNetRPC) GetSchedule(ctx context.Context, playerID PlayerID, region string) (*GetScheduleResponse, error) {
	request := GetScheduleRequest{
		PlayerID: playerID,
		Region:   region,
	}

	var result GetScheduleResponse
	err := p.sendRequestSync(ctx, "Tournaments/Search/GetSchedule v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPublicTournaments retrieves a list of public tournaments.
func (p *PsyNetRPC) GetPublicTournaments(ctx context.Context, region string, limit, offset int) (*GetPublicTournamentsResponse, error) {
	request := GetPublicTournamentsRequest{
		Region: region,
		Limit:  limit,
		Offset: offset,
	}

	var result GetPublicTournamentsResponse
	err := p.sendRequestSync(ctx, "Tournaments/Search/GetPublicTournaments v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RegisterTournament registers a player for a tournament.
func (p *PsyNetRPC) RegisterTournament(ctx context.Context, playerID PlayerID, tournamentID string, teamID *string) (*RegisterTournamentResponse, error) {
	request := RegisterTournamentRequest{
		PlayerID:     playerID,
		TournamentID: tournamentID,
		TeamID:       teamID,
	}

	var result RegisterTournamentResponse
	err := p.sendRequestSync(ctx, "Tournaments/Registration/RegisterTournament v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UnsubscribeTournament unsubscribes a player from a tournament.
func (p *PsyNetRPC) UnsubscribeTournament(ctx context.Context, playerID PlayerID, tournamentID string) (*UnsubscribeTournamentResponse, error) {
	request := UnsubscribeTournamentRequest{
		PlayerID:     playerID,
		TournamentID: tournamentID,
	}

	var result UnsubscribeTournamentResponse
	err := p.sendRequestSync(ctx, "Tournaments/Registration/UnsubscribeTournament v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTournamentSubscriptions retrieves a player's tournament subscriptions.
func (p *PsyNetRPC) GetTournamentSubscriptions(ctx context.Context, playerID PlayerID) (*GetTournamentSubscriptionsResponse, error) {
	request := GetTournamentSubscriptionsRequest{
		PlayerID: playerID,
	}

	var result GetTournamentSubscriptionsResponse
	err := p.sendRequestSync(ctx, "Tournaments/Status/GetTournamentSubscriptions v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
