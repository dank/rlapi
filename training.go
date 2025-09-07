package rlapi

import "context"

type TrainingPack struct {
	Code            string   `json:"Code"`
	TMName          string   `json:"TM_Name"`
	Type            int      `json:"Type"`
	Difficulty      int      `json:"Difficulty"`
	CreatorName     string   `json:"CreatorName"`
	CreatorPlayerID string   `json:"CreatorPlayerID,omitempty"`
	MapName         string   `json:"MapName"`
	Tags            []string `json:"Tags"`
	NumRounds       int      `json:"NumRounds"`
	TMGuid          string   `json:"TM_Guid"`
	CreatedAt       int64    `json:"CreatedAt"`
	UpdatedAt       int64    `json:"UpdatedAt"`
}

type BrowseTrainingDataRequest struct {
	FeaturedOnly bool `json:"bFeaturedOnly"`
}

type BrowseTrainingDataResponse struct {
	TrainingData []TrainingPack `json:"TrainingData"`
}

type GetTrainingMetadataRequest struct {
	Codes []string `json:"Codes"`
}

type GetTrainingMetadataResponse struct {
	TrainingData []TrainingPack `json:"TrainingData"`
}

// BrowseTrainingData retrieves training packs.
func (p *PsyNetRPC) BrowseTrainingData(ctx context.Context, featuredOnly bool) ([]TrainingPack, error) {
	request := BrowseTrainingDataRequest{
		FeaturedOnly: featuredOnly,
	}

	var result BrowseTrainingDataResponse
	err := p.sendRequestSync(ctx, "Training/BrowseTrainingData v1", request, &result)
	if err != nil {
		return nil, err
	}
	return result.TrainingData, nil
}

// GetTrainingMetadata retrieves training pack metadata for the given codes.
func (p *PsyNetRPC) GetTrainingMetadata(ctx context.Context, codes []string) ([]TrainingPack, error) {
	request := GetTrainingMetadataRequest{
		Codes: codes,
	}

	var result GetTrainingMetadataResponse
	err := p.sendRequestSync(ctx, "Training/GetTrainingMetadata v1", request, &result)
	if err != nil {
		return nil, err
	}
	return result.TrainingData, nil
}
