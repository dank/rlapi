package rlapi

import "context"

type ContainerDrop struct {
	ProductID int       `json:"ProductID"`
	SeriesID  int       `json:"SeriesID"`
	Drops     []Product `json:"Drops"`
}

// UnlockResult represents the result of unlocking a container
type UnlockResult struct {
	UnlockedItems []Product `json:"UnlockedItems"`
	UsedKeys      []string  `json:"UsedKeys"`
	RemainingKeys []Product `json:"RemainingKeys"`
}

// TradeInResult represents the result of trading in items
type TradeInResult struct {
	ReceivedItems []Product `json:"ReceivedItems"`
	TradedItems   []string  `json:"TradedItems"`
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
	ProductData []Product `json:"ProductData"`
}

type GetContainerDropTableResponse struct {
	ContainerDrops []ContainerDrop `json:"ContainerDrops"`
}

type UnlockContainerRequest struct {
	PlayerID       PlayerID `json:"PlayerID"`
	InstanceIDs    []string `json:"InstanceIDs"`
	KeyInstanceIDs []string `json:"KeyInstanceIDs"`
}

type UnlockContainerResponse struct {
	Drops []Product `json:"Drops"`
}

type TradeInRequest struct {
	PlayerID         PlayerID `json:"PlayerID"`
	ProductInstances []string `json:"ProductInstances"`
}

type TradeInResponse struct {
	Drops []Product `json:"Drops"`
}

type GetProductStatusResponse struct {
	CrossEntitledProductIDs []int         `json:"CrossEntitledProductIDs"`
	LockedProductIDs        []interface{} `json:"LockedProductIDs"`
}

// GetPlayerProducts retrieves all products/items owned by the authenticated player.
func (p *PsyNetRPC) GetPlayerProducts(ctx context.Context, playerID PlayerID, updatedTimestamp string) ([]Product, error) {
	request := GetPlayerProductsRequest{
		PlayerID:         playerID,
		UpdatedTimestamp: updatedTimestamp,
	}

	var result GetPlayerProductsResponse
	err := p.sendRequestSync(ctx, "Products/GetPlayerProducts v2", request, &result)
	if err != nil {
		return nil, err
	}
	return result.ProductData, nil
}

// GetContainerDropTable retrieves the drop table for containers.
func (p *PsyNetRPC) GetContainerDropTable(ctx context.Context) ([]ContainerDrop, error) {
	var result GetContainerDropTableResponse
	err := p.sendRequestSync(ctx, "Products/GetContainerDropTable v2", emptyRequest{}, &result)
	if err != nil {
		return nil, err
	}
	return result.ContainerDrops, nil
}

// UnlockContainer unlocks containers returns the dropped items.
func (p *PsyNetRPC) UnlockContainer(ctx context.Context, playerID PlayerID, instanceIDs []string) ([]Product, error) {
	request := UnlockContainerRequest{
		PlayerID:       playerID,
		InstanceIDs:    instanceIDs,
		KeyInstanceIDs: []string{},
	}

	var result UnlockContainerResponse
	err := p.sendRequestSync(ctx, "Products/UnlockContainer v2", request, &result)
	if err != nil {
		return nil, err
	}
	return result.Drops, nil
}

// TradeIn trades in multiple items for new items.
func (p *PsyNetRPC) TradeIn(ctx context.Context, playerID PlayerID, productInstances []string) ([]Product, error) {
	request := TradeInRequest{
		PlayerID:         playerID,
		ProductInstances: productInstances,
	}

	var result TradeInResponse
	err := p.sendRequestSync(ctx, "Products/TradeIn v2", request, &result)
	if err != nil {
		return nil, err
	}
	return result.Drops, nil
}

func (p *PsyNetRPC) GetCrossEntitlementProductStatus(ctx context.Context) (*GetProductStatusResponse, error) {
	var result GetProductStatusResponse
	err := p.sendRequestSync(ctx, "Products/CrossEntitlement/GetProductStatus v1", emptyRequest{}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
