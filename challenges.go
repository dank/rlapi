package rlapi

import (
	"context"
)

type ChallengeID int

// Challenge represents an in-game challenge (eg, season challenges, weekly challenges, etc.)
type Challenge struct {
	ID                 ChallengeID            `json:"ID"`
	Title              string                 `json:"Title"`
	Description        string                 `json:"Description"`
	Sort               int                    `json:"Sort"`
	GroupID            int                    `json:"GroupID"`
	XPUnlockLevel      int                    `json:"XPUnlockLevel"`
	IsRepeatable       bool                   `json:"bIsRepeatable"`
	RepeatLimit        int                    `json:"RepeatLimit"`
	IconURL            string                 `json:"IconURL"`
	BackgroundURL      *string                `json:"BackgroundURL"`
	BackgroundColor    int                    `json:"BackgroundColor"`
	Requirements       []ChallengeRequirement `json:"Requirements"`
	Rewards            ChallengeRewards       `json:"Rewards"`
	AutoClaimRewards   bool                   `json:"bAutoClaimRewards"`
	IsPremium          bool                   `json:"bIsPremium"`
	UnlockChallengeIDs []ChallengeID          `json:"UnlockChallengeIDs"`
}

type ChallengeRequirement struct {
	RequiredCount int `json:"RequiredCount"`
}

// ChallengeRewards represents rewards for completing a challenge
type ChallengeRewards struct {
	XP       int                      `json:"XP"`
	Currency []interface{}            `json:"Currency"`
	Products []ChallengeRewardProduct `json:"Products"`
	Pips     int                      `json:"Pips"`
}

// ChallengeRewardProduct represents a product reward from a challenge
type ChallengeRewardProduct struct {
	ID                 string             `json:"ID"`
	ChallengeID        ChallengeID        `json:"ChallengeID"`
	ProductID          int                `json:"ProductID"`
	InstanceID         *string            `json:"InstanceID"`
	OriginalInstanceID *string            `json:"OriginalInstanceID"`
	Attributes         []ProductAttribute `json:"Attributes"`
	SeriesID           int                `json:"SeriesID"`
	TradeHold          *string            `json:"TradeHold"`
	AddedTimestamp     *int64             `json:"AddedTimestamp"`
	UpdatedTimestamp   *int64             `json:"UpdatedTimestamp"`
	DeletedTimestamp   *int64             `json:"DeletedTimestamp"`
}

// ChallengeProgress represents a player's progress towards a challenge.
type ChallengeProgress struct {
	ID                  ChallengeID           `json:"ID"`
	CompleteCount       int                   `json:"CompleteCount"`
	IsHidden            bool                  `json:"bIsHidden"`
	NotifyCompleted     bool                  `json:"bNotifyCompleted"`
	NotifyAvailable     bool                  `json:"bNotifyAvailable"`
	NotifyNewInfo       bool                  `json:"bNotifyNewInfo"`
	RewardsAvailable    bool                  `json:"bRewardsAvailable"`
	IsComplete          bool                  `json:"bComplete"`
	RequirementProgress []RequirementProgress `json:"RequirementProgress"`
	ProgressResetTime   int64                 `json:"ProgressResetTimeUTC"`
}

type RequirementProgress struct {
	ProgressCount  int `json:"ProgressCount"`
	ProgressChange int `json:"ProgressChange"`
}

type GetActiveChallengesRequest struct {
	Challenges []interface{} `json:"Challenges"`
	Folders    []interface{} `json:"Folders"`
}

type GetActiveChallengesResponse struct {
	Challenges []Challenge `json:"Challenges"`
}

type PlayerProgressRequest struct {
	PlayerID     PlayerID      `json:"PlayerID"`
	ChallengeIDs []ChallengeID `json:"ChallengeIDs"`
}

type PlayerProgressResponse struct {
	ProgressData []ChallengeProgress `json:"ProgressData"`
}

type CollectRewardRequest struct {
	PlayerID    PlayerID    `json:"PlayerID"`
	ChallengeID ChallengeID `json:"ID"`
}

type FTECheckpointCompleteRequest struct {
	PlayerID       PlayerID `json:"PlayerID"`
	GroupName      string   `json:"GroupName"`
	CheckpointName string   `json:"CheckpointName"`
}

type FTEGroupCompleteRequest struct {
	PlayerID  PlayerID `json:"PlayerID"`
	GroupName string   `json:"GroupName"`
}

// GetActiveChallenges retrieves the list of currently active challenges.
func (p *PsyNetRPC) GetActiveChallenges(ctx context.Context) ([]Challenge, error) {
	request := GetActiveChallengesRequest{
		Challenges: []interface{}{},
		Folders:    []interface{}{},
	}

	var result GetActiveChallengesResponse
	err := p.sendRequestSync(ctx, "Challenges/GetActiveChallenges v1", request, &result)
	if err != nil {
		return nil, err
	}
	return result.Challenges, nil
}

// GetChallengeProgress retrieves a player's progression on challenges.
func (p *PsyNetRPC) GetChallengeProgress(ctx context.Context, playerID PlayerID, challengeIDs []ChallengeID) ([]ChallengeProgress, error) {
	request := PlayerProgressRequest{
		PlayerID:     playerID,
		ChallengeIDs: challengeIDs,
	}

	var result PlayerProgressResponse
	err := p.sendRequestSync(ctx, "Challenges/PlayerProgress v1", request, &result)
	if err != nil {
		return nil, err
	}
	return result.ProgressData, nil
}

// CollectChallengeReward collects rewards from a completed challenge.
func (p *PsyNetRPC) CollectChallengeReward(ctx context.Context, playerID PlayerID, challengeID ChallengeID) error {
	request := CollectRewardRequest{
		PlayerID:    playerID,
		ChallengeID: challengeID,
	}

	var result interface{}
	err := p.sendRequestSync(ctx, "Challenges/CollectReward v1", request, &result)
	if err != nil {
		return err
	}
	return nil
}

// FTECheckpointComplete marks a First Time Experience (FTE) checkpoint as complete.
func (p *PsyNetRPC) FTECheckpointComplete(ctx context.Context, playerID PlayerID, groupName string, checkpointName string) error {
	request := FTECheckpointCompleteRequest{
		PlayerID:       playerID,
		GroupName:      groupName,
		CheckpointName: checkpointName,
	}

	var result interface{}
	err := p.sendRequestSync(ctx, "Challenges/FTECheckpointComplete v1", request, &result)
	if err != nil {
		return err
	}
	return nil
}

// FTEGroupComplete marks a First Time Experience (FTE) group as complete.
func (p *PsyNetRPC) FTEGroupComplete(ctx context.Context, playerID PlayerID, groupName string) error {
	request := FTEGroupCompleteRequest{
		PlayerID:  playerID,
		GroupName: groupName,
	}

	var result interface{}
	err := p.sendRequestSync(ctx, "Challenges/FTEGroupComplete v1", request, &result)
	if err != nil {
		return err
	}
	return nil
}
