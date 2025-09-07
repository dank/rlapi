package rlapi

import "context"

type ShopID int

// Shop represents a shop in the game
type Shop struct {
	ID        ShopID  `json:"ID"`
	Type      string  `json:"Type"`
	StartDate int64   `json:"StartDate"`
	EndDate   *int64  `json:"EndDate"`
	LogoURL   *string `json:"LogoURL"`
	Name      *string `json:"Name"`
	Title     *string `json:"Title"`
}

// ShopCatalogue represents the catalogue for a given shop
type ShopCatalogue struct {
	ShopID    ShopID     `json:"ShopID"`
	ShopItems []ShopItem `json:"ShopItems"`
}

// ShopItem represents an item available for purchase in a shop
type ShopItem struct {
	ShopItemID             int                   `json:"ShopItemID"`
	StartDate              int64                 `json:"StartDate"`
	EndDate                *int64                `json:"EndDate"`
	MaxQuantityPerPlayer   *int                  `json:"MaxQuantityPerPlayer"`
	ImageURL               *string               `json:"ImageURL"`
	DeliverableProducts    []DeliverableProduct  `json:"DeliverableProducts"`
	DeliverableCurrencies  []DeliverableCurrency `json:"DeliverableCurrencies"`
	Costs                  []ShopItemCost        `json:"Costs"`
	ShopItemLocations      []int                 `json:"ShopItemLocations"`
	Title                  *string               `json:"Title"`
	Description            *string               `json:"Description"`
	FeaturedCollections    []interface{}         `json:"FeaturedCollections"`
	Attributes             []ProductAttribute    `json:"Attributes"`
	Disclaimer             *string               `json:"Disclaimer"`
	PurchasedQuantity      int                   `json:"PurchasedQuantity"`
	Purchasable            bool                  `json:"Purchasable"`
	MaxQuantityPerDay      *int                  `json:"MaxQuantityPerDay"`
	DailyPurchasedQuantity *int                  `json:"DailyPurchasedQuantity"`
}

// DeliverableProduct represents a product that can be delivered from a shop purchase
type DeliverableProduct struct {
	Count   int     `json:"Count"`
	Product Product `json:"Product"`
	SortID  *int    `json:"SortID"`
	IsOwned *bool   `json:"IsOwned,omitempty"`
}

// DeliverableCurrency represents currency that can be delivered from a shop purchase
type DeliverableCurrency struct {
	ID     int `json:"ID"`
	Amount int `json:"Amount"`
}

// Product represents a game product/item
type Product struct {
	ProductID        int                `json:"ProductID"`
	InstanceID       *string            `json:"InstanceID"`
	Attributes       []ProductAttribute `json:"Attributes"`
	SeriesID         int                `json:"SeriesID"`
	AddedTimestamp   *int64             `json:"AddedTimestamp"`
	UpdatedTimestamp *int64             `json:"UpdatedTimestamp"`
}

// ProductAttribute represents an attribute of a product
type ProductAttribute struct {
	Key   string      `json:"Key"`
	Value interface{} `json:"Value"`
}

// ShopItemCost represents the cost of a shop item
type ShopItemCost struct {
	ResetTime      *int64          `json:"ResetTime"`
	ShopItemCostID int             `json:"ShopItemCostID"`
	Discount       *interface{}    `json:"Discount"`
	BulkDiscounts  *interface{}    `json:"BulkDiscounts"`
	StartDate      int64           `json:"StartDate"`
	EndDate        *int64          `json:"EndDate"`
	Price          []CurrencyPrice `json:"Price"`
	SortID         int             `json:"SortID"`
	DisplayTypeID  int             `json:"DisplayTypeID"`
	ShopScaledCost interface{}     `json:"ShopScaledCost"`
}

// CurrencyPrice represents a price in the specified currency
type CurrencyPrice struct {
	ID     int `json:"ID"`
	Amount int `json:"Amount"`
}

type GetStandardShopsResponse struct {
	Shops []Shop `json:"Shops"`
}

type GetShopCatalogueRequest struct {
	ShopIDs []ShopID `json:"ShopIDs"`
}

type GetShopCatalogueResponse struct {
	Catalogues []ShopCatalogue `json:"Catalogues"`
}

type GetPlayerWalletRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
}

type GetPlayerWalletResponse struct {
	Currencies []struct {
		ID               int     `json:"ID"`
		Amount           int     `json:"Amount"`
		ExpirationTime   *string `json:"ExpirationTime"`
		UpdatedTimestamp int64   `json:"UpdatedTimestamp"`
		IsTradable       bool    `json:"IsTradable"`
		TradeHold        *string `json:"TradeHold"`
	} `json:"Currencies"`
}

type GetShopNotificationsResponse struct {
	ShopNotifications []struct {
		ShopNotificationID  int                  `json:"ShopNotificationID"`
		ShopItemCostID      int                  `json:"ShopItemCostID"`
		StartTime           int64                `json:"StartTime"`
		EndTime             int64                `json:"EndTime"`
		ImageURL            *string              `json:"ImageURL"`
		Title               string               `json:"Title"`
		DeliverableProducts []DeliverableProduct `json:"DeliverableProducts"`
	} `json:"ShopNotifications"`
}

// GetStandardShops retrieves the list of available shops.
func (p *PsyNetRPC) GetStandardShops(ctx context.Context) (*GetStandardShopsResponse, error) {
	var result GetStandardShopsResponse
	err := p.sendRequestSync(ctx, "Shops/GetStandardShops v1", map[string]interface{}{}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetShopCatalogue retrieves detailed information about items available in specified shops.
func (p *PsyNetRPC) GetShopCatalogue(ctx context.Context, shopIDs []ShopID) (*GetShopCatalogueResponse, error) {
	request := GetShopCatalogueRequest{
		ShopIDs: shopIDs,
	}

	var result GetShopCatalogueResponse
	err := p.sendRequestSync(ctx, "Shops/GetShopCatalogue v2", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPlayerWallet retrieves the authenticated player's wallet information.
func (p *PsyNetRPC) GetPlayerWallet(ctx context.Context, playerID PlayerID) (*GetPlayerWalletResponse, error) {
	var result GetPlayerWalletResponse
	err := p.sendRequestSync(ctx, "Shops/GetPlayerWallet v1", GetPlayerWalletRequest{PlayerID: playerID}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetShopNotifications retrieves current shop notifications.
func (p *PsyNetRPC) GetShopNotifications(ctx context.Context) (*GetShopNotificationsResponse, error) {
	var result GetShopNotificationsResponse
	err := p.sendRequestSync(ctx, "Shops/GetShopNotifications v1", map[string]interface{}{}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
