package rlapi

import "context"

// ProductData represents a player's product/item data
type ProductData struct {
	ProductID        int                `json:"ProductID"`
	InstanceID       string             `json:"InstanceID"`
	Attributes       []ProductAttribute `json:"Attributes"`
	SeriesID         int                `json:"SeriesID"`
	AddedTimestamp   int64              `json:"AddedTimestamp"`
	UpdatedTimestamp int64              `json:"UpdatedTimestamp"`
	DeletedTimestamp *int64             `json:"DeletedTimestamp"`
}

// ContainerDropTable represents the drop table for containers
type ContainerDropTable struct {
	ContainerID int                      `json:"ContainerID"`
	DropRates   []ContainerDropRate      `json:"DropRates"`
	Items       []ContainerDropTableItem `json:"Items"`
}

// ContainerDropRate represents drop rates for different item rarities
type ContainerDropRate struct {
	Rarity string  `json:"Rarity"`
	Rate   float64 `json:"Rate"`
}

// ContainerDropTableItem represents an item that can drop from a container
type ContainerDropTableItem struct {
	ProductID  int                `json:"ProductID"`
	Rarity     string             `json:"Rarity"`
	Attributes []ProductAttribute `json:"Attributes"`
	Weight     int                `json:"Weight"`
}

// UnlockResult represents the result of unlocking a container
type UnlockResult struct {
	UnlockedItems []ProductData `json:"UnlockedItems"`
	UsedKeys      []string      `json:"UsedKeys"`
	RemainingKeys []ProductData `json:"RemainingKeys"`
}

// TradeInResult represents the result of trading in items
type TradeInResult struct {
	ReceivedItems []ProductData `json:"ReceivedItems"`
	TradedItems   []string      `json:"TradedItems"`
}

// CrossEntitlementStatus represents cross-platform entitlement status
type CrossEntitlementStatus struct {
	CrossEntitledProductIDs []int `json:"CrossEntitledProductIDs"`
	LockedProductIDs        []int `json:"LockedProductIDs"`
}

type GetPlayerProductsRequest struct {
	PlayerID         PlayerID `json:"PlayerID"`
	UpdatedTimestamp string   `json:"UpdatedTimestamp"`
}

type GetPlayerProductsResponse struct {
	ProductData []ProductData `json:"ProductData"`
}

type GetContainerDropTableResponse struct {
	DropTables []ContainerDropTable `json:"DropTables"`
}

type UnlockContainerRequest struct {
	PlayerID       PlayerID `json:"PlayerID"`
	InstanceIDs    []string `json:"InstanceIDs"`
	KeyInstanceIDs []string `json:"KeyInstanceIDs"`
}

type UnlockContainerResponse struct {
	Results []UnlockResult `json:"Results"`
}

type TradeInRequest struct {
	PlayerID         PlayerID `json:"PlayerID"`
	ProductInstances []string `json:"ProductInstances"`
}

type TradeInResponse struct {
	TradeInResults []TradeInResult `json:"TradeInResults"`
}

type GetProductStatusResponse struct {
	CrossEntitledProductIDs []int `json:"CrossEntitledProductIDs"`
	LockedProductIDs        []int `json:"LockedProductIDs"`
}

// GetPlayerProducts retrieves all products owned by a specific player.
func (p *PsyNetRPC) GetPlayerProducts(ctx context.Context, playerID PlayerID, updatedTimestamp string) (*GetPlayerProductsResponse, error) {
	request := GetPlayerProductsRequest{
		PlayerID:         playerID,
		UpdatedTimestamp: updatedTimestamp,
	}

	var result GetPlayerProductsResponse
	err := p.sendRequestSync(ctx, "Products/GetPlayerProducts v2", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetContainerDropTable retrieves the drop table for containers.
func (p *PsyNetRPC) GetContainerDropTable(ctx context.Context) (*GetContainerDropTableResponse, error) {
	var result GetContainerDropTableResponse
	err := p.sendRequestSync(ctx, "Products/GetContainerDropTable v2", map[string]interface{}{}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UnlockContainer unlocks containers using keys and returns the unlocked items.
func (p *PsyNetRPC) UnlockContainer(ctx context.Context, playerID PlayerID, instanceIDs, keyInstanceIDs []string) (*UnlockContainerResponse, error) {
	request := UnlockContainerRequest{
		PlayerID:       playerID,
		InstanceIDs:    instanceIDs,
		KeyInstanceIDs: keyInstanceIDs,
	}

	var result UnlockContainerResponse
	err := p.sendRequestSync(ctx, "Products/UnlockContainer v2", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// TradeIn trades in multiple items for new items of higher rarity.
func (p *PsyNetRPC) TradeIn(ctx context.Context, playerID PlayerID, productInstances []string) (*TradeInResponse, error) {
	request := TradeInRequest{
		PlayerID:         playerID,
		ProductInstances: productInstances,
	}

	var result TradeInResponse
	err := p.sendRequestSync(ctx, "Products/TradeIn v2", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCrossEntitlementProductStatus retrieves cross-platform product entitlement status.
func (p *PsyNetRPC) GetCrossEntitlementProductStatus(ctx context.Context) (*GetProductStatusResponse, error) {
	var result GetProductStatusResponse
	err := p.sendRequestSync(ctx, "Products/CrossEntitlement/GetProductStatus v1", map[string]interface{}{}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
