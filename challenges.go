package rlapi

import (
	"context"
	"fmt"
)

// Challenge represents a game challenge
type Challenge struct {
	ID                 int                    `json:"ID"`
	Title              string                 `json:"Title"`
	Description        string                 `json:"Description"`
	Sort               int                    `json:"Sort"`
	GroupID            int                    `json:"GroupID"`
	XPUnlockLevel      int                    `json:"XPUnlockLevel"`
	BIsRepeatable      bool                   `json:"bIsRepeatable"`
	RepeatLimit        int                    `json:"RepeatLimit"`
	IconURL            string                 `json:"IconURL"`
	BackgroundURL      *string                `json:"BackgroundURL"`
	BackgroundColor    int                    `json:"BackgroundColor"`
	Requirements       []ChallengeRequirement `json:"Requirements"`
	Rewards            ChallengeRewards       `json:"Rewards"`
	BAutoClaimRewards  bool                   `json:"bAutoClaimRewards"`
	BIsPremium         bool                   `json:"bIsPremium"`
	UnlockChallengeIDs []int                  `json:"UnlockChallengeIDs"`
}

// ChallengeRequirement represents a requirement for completing a challenge
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
	ChallengeID        int                `json:"ChallengeID"`
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

// PlayerProgress represents a player's progress on challenges
type PlayerProgress struct {
	PlayerID        PlayerID         `json:"PlayerID"`
	ChallengeStates []ChallengeState `json:"ChallengeStates"`
	Pips            int              `json:"Pips"`
	StarLevel       int              `json:"StarLevel"`
	StarLevelXP     int              `json:"StarLevelXP"`
}

// ChallengeState represents the state of a challenge for a player
type ChallengeState struct {
	ID              int                         `json:"ID"`
	CompletedCount  int                         `json:"CompletedCount"`
	Requirements    []ChallengeRequirementState `json:"Requirements"`
	Completed       bool                        `json:"Completed"`
	ClaimedRewards  bool                        `json:"ClaimedRewards"`
	ClaimedProducts []interface{}               `json:"ClaimedProducts"`
}

// ChallengeRequirementState represents the state of a challenge requirement
type ChallengeRequirementState struct {
	Progress int `json:"Progress"`
}

type GetActiveChallengesRequest struct {
	Challenges []interface{} `json:"Challenges"`
	Folders    []interface{} `json:"Folders"`
}

type GetActiveChallengesResponse struct {
	Challenges []Challenge `json:"Challenges"`
}

type PlayerProgressRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
}

type PlayerProgressResponse struct {
	PlayerProgress PlayerProgress `json:"PlayerProgress"`
}

type CollectRewardRequest struct {
	PlayerID    PlayerID `json:"PlayerID"`
	ChallengeID int      `json:"ChallengeID"`
}

type CollectRewardResponse struct {
	CollectedRewards ChallengeRewards `json:"CollectedRewards"`
	UpdatedProgress  PlayerProgress   `json:"UpdatedProgress"`
}

type FTECheckpointCompleteRequest struct {
	PlayerID     PlayerID `json:"PlayerID"`
	CheckpointID int      `json:"CheckpointID"`
}

type FTECheckpointCompleteResponse struct {
	Success bool `json:"Success"`
}

type FTEGroupCompleteRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
	GroupID  int      `json:"GroupID"`
}

type FTEGroupCompleteResponse struct {
	Success bool `json:"Success"`
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

// PlayerProgress retrieves a player's progress on challenges.
func (p *PsyNetRPC) PlayerProgress(ctx context.Context, playerID PlayerID) (*PlayerProgress, error) {
	request := PlayerProgressRequest{
		PlayerID: playerID,
	}

	var result PlayerProgressResponse
	err := p.sendRequestSync(ctx, "Challenges/PlayerProgress v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result.PlayerProgress, nil
}

// CollectRewardResult represents the result of collecting a reward
type CollectRewardResult struct {
	CollectedRewards ChallengeRewards `json:"CollectedRewards"`
	UpdatedProgress  PlayerProgress   `json:"UpdatedProgress"`
}

// CollectReward collects rewards from a completed challenge.
func (p *PsyNetRPC) CollectReward(ctx context.Context, playerID PlayerID, challengeID int) (*CollectRewardResult, error) {
	request := CollectRewardRequest{
		PlayerID:    playerID,
		ChallengeID: challengeID,
	}

	var result CollectRewardResponse
	err := p.sendRequestSync(ctx, "Challenges/CollectReward v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &CollectRewardResult{
		CollectedRewards: result.CollectedRewards,
		UpdatedProgress:  result.UpdatedProgress,
	}, nil
}

// FTECheckpointComplete marks a First Time Experience (FTE) checkpoint as complete.
func (p *PsyNetRPC) FTECheckpointComplete(ctx context.Context, playerID PlayerID, checkpointID int) error {
	request := FTECheckpointCompleteRequest{
		PlayerID:     playerID,
		CheckpointID: checkpointID,
	}

	var result FTECheckpointCompleteResponse
	err := p.sendRequestSync(ctx, "Challenges/FTECheckpointComplete v1", request, &result)
	if err != nil {
		return err
	}
	if !result.Success {
		return fmt.Errorf("failed to complete FTE checkpoint %d", checkpointID)
	}
	return nil
}

// FTEGroupComplete marks a First Time Experience (FTE) group as complete.
func (p *PsyNetRPC) FTEGroupComplete(ctx context.Context, playerID PlayerID, groupID int) error {
	request := FTEGroupCompleteRequest{
		PlayerID: playerID,
		GroupID:  groupID,
	}

	var result FTEGroupCompleteResponse
	err := p.sendRequestSync(ctx, "Challenges/FTEGroupComplete v1", request, &result)
	if err != nil {
		return err
	}
	if !result.Success {
		return fmt.Errorf("failed to complete FTE group %d", groupID)
	}
	return nil
}
