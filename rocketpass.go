package rlapi

import "context"

// RocketPassInfo represents information about a player's Rocket Pass progress
type RocketPassInfo struct {
	TierLevel    int     `json:"TierLevel"`
	BOwnsPremium bool    `json:"bOwnsPremium"`
	XPMultiplier float64 `json:"XPMultiplier"`
	Pips         int     `json:"Pips"`
	PipsPerLevel int     `json:"PipsPerLevel"`
}

// RocketPassStore represents purchasable items in the Rocket Pass store
type RocketPassStore struct {
	Tiers   []RocketPassTier   `json:"Tiers"`
	Bundles []RocketPassBundle `json:"Bundles"`
}

// RocketPassTier represents a purchasable tier in the Rocket Pass
type RocketPassTier struct {
	PurchasableID        int     `json:"PurchasableID"`
	CurrencyID           int     `json:"CurrencyID"`
	CurrencyCost         int     `json:"CurrencyCost"`
	OriginalCurrencyCost *int    `json:"OriginalCurrencyCost"`
	Tiers                int     `json:"Tiers"`
	Savings              int     `json:"Savings"`
	ImageURL             *string `json:"ImageUrl"`
}

// RocketPassBundle represents a purchasable bundle in the Rocket Pass
type RocketPassBundle struct {
	PurchasableID        int     `json:"PurchasableID"`
	CurrencyID           int     `json:"CurrencyID"`
	CurrencyCost         int     `json:"CurrencyCost"`
	OriginalCurrencyCost *int    `json:"OriginalCurrencyCost"`
	Tiers                int     `json:"Tiers"`
	Savings              int     `json:"Savings"`
	ImageURL             *string `json:"ImageUrl"`
}

// RocketPassReward represents rewards available at specific tiers
type RocketPassReward struct {
	Tier          int            `json:"Tier"`
	ProductData   []ProductData  `json:"ProductData"`
	XPRewards     []XPReward     `json:"XPRewards"`
	CurrencyDrops []CurrencyDrop `json:"CurrencyDrops"`
}

// XPReward represents an XP-based reward
type XPReward struct {
	Name   string  `json:"Name"`
	Amount float64 `json:"Amount"`
}

// CurrencyDrop represents a currency reward
type CurrencyDrop struct {
	ID         int `json:"ID"`
	CurrencyID int `json:"CurrencyID"`
	Amount     int `json:"Amount"`
}

// RocketPassPlayerInfo represents a player's complete Rocket Pass information
type RocketPassPlayerInfo struct {
	StartTime       int64           `json:"StartTime"`
	EndTime         int64           `json:"EndTime"`
	RocketPassInfo  RocketPassInfo  `json:"RocketPassInfo"`
	RocketPassStore RocketPassStore `json:"RocketPassStore"`
}

// RocketPassRewardContent represents the reward content for a Rocket Pass
type RocketPassRewardContent struct {
	TierCap         int                `json:"TierCap"`
	FreeMaxLevel    int                `json:"FreeMaxLevel"`
	PremiumMaxLevel int                `json:"PremiumMaxLevel"`
	FreeRewards     []RocketPassReward `json:"FreeRewards"`
	PremiumRewards  []RocketPassReward `json:"PremiumRewards"`
}

// PrestigeReward represents a prestige reward in Rocket Pass
type PrestigeReward struct {
	Level       int           `json:"Level"`
	ProductData []ProductData `json:"ProductData"`
	Currency    []interface{} `json:"Currency"`
}

type GetPlayerInfoRequest struct {
	PlayerID        PlayerID    `json:"PlayerID"`
	RocketPassID    int         `json:"RocketPassID"`
	RocketPassInfo  interface{} `json:"RocketPassInfo"`
	RocketPassStore interface{} `json:"RocketPassStore"`
}

type GetPlayerInfoResponse struct {
	StartTime       int64           `json:"StartTime"`
	EndTime         int64           `json:"EndTime"`
	RocketPassInfo  RocketPassInfo  `json:"RocketPassInfo"`
	RocketPassStore RocketPassStore `json:"RocketPassStore"`
}

type GetRewardContentRequest struct {
	RocketPassID    int `json:"RocketPassID"`
	TierCap         int `json:"TierCap"`
	FreeMaxLevel    int `json:"FreeMaxLevel"`
	PremiumMaxLevel int `json:"PremiumMaxLevel"`
}

type GetRewardContentResponse struct {
	TierCap         int                `json:"TierCap"`
	FreeMaxLevel    int                `json:"FreeMaxLevel"`
	PremiumMaxLevel int                `json:"PremiumMaxLevel"`
	FreeRewards     []RocketPassReward `json:"FreeRewards"`
	PremiumRewards  []RocketPassReward `json:"PremiumRewards"`
}

type GetPlayerPrestigeRewardsRequest struct {
	PlayerID     PlayerID `json:"PlayerID"`
	RocketPassID int      `json:"RocketPassID"`
}

type GetPlayerPrestigeRewardsResponse struct {
	PrestigeRewards []PrestigeReward `json:"PrestigeRewards"`
}

// GetPlayerInfo retrieves Rocket Pass information for a specific player.
func (p *PsyNetRPC) GetRocketPassPlayerInfo(ctx context.Context, playerID PlayerID, rocketPassID int) (*GetPlayerInfoResponse, error) {
	request := GetPlayerInfoRequest{
		PlayerID:        playerID,
		RocketPassID:    rocketPassID,
		RocketPassInfo:  map[string]interface{}{},
		RocketPassStore: map[string]interface{}{},
	}

	var result GetPlayerInfoResponse
	err := p.sendRequestSync(ctx, "RocketPass/GetPlayerInfo v2", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRewardContent retrieves the reward content for a specific Rocket Pass.
func (p *PsyNetRPC) GetRocketPassRewardContent(ctx context.Context, rocketPassID, tierCap, freeMaxLevel, premiumMaxLevel int) (*GetRewardContentResponse, error) {
	request := GetRewardContentRequest{
		RocketPassID:    rocketPassID,
		TierCap:         tierCap,
		FreeMaxLevel:    freeMaxLevel,
		PremiumMaxLevel: premiumMaxLevel,
	}

	var result GetRewardContentResponse
	err := p.sendRequestSync(ctx, "RocketPass/GetRewardContent v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPlayerPrestigeRewards retrieves prestige rewards for a player's Rocket Pass.
func (p *PsyNetRPC) GetPlayerPrestigeRewards(ctx context.Context, playerID PlayerID, rocketPassID int) (*GetPlayerPrestigeRewardsResponse, error) {
	request := GetPlayerPrestigeRewardsRequest{
		PlayerID:     playerID,
		RocketPassID: rocketPassID,
	}

	var result GetPlayerPrestigeRewardsResponse
	err := p.sendRequestSync(ctx, "RocketPass/GetPlayerPrestigeRewards v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
