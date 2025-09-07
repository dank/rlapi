package rlapi

import "context"

type MTXProduct struct {
	ID                int           `json:"ID"`
	Title             string        `json:"Title"`
	Description       string        `json:"Description"`
	TabTitle          string        `json:"TabTitle"`
	PriceDescription  string        `json:"PriceDescription"`
	ImageURL          string        `json:"ImageURL"`
	PlatformProductID string        `json:"PlatformProductID"`
	IsOwned           bool          `json:"bIsOwned"`
	Items             []Product     `json:"Items"`
	Currencies        []MTXCurrency `json:"Currencies"`
}

type MTXCurrency struct {
	ID         int `json:"ID"`
	CurrencyID int `json:"CurrencyID"`
	Amount     int `json:"Amount"`
}

type MTXCartItem struct {
	CatalogID int `json:"CatalogID"`
	Count     int `json:"Count"`
}

type GetCatalogRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
	Category string   `json:"Category"`
}

type GetCatalogResponse struct {
	MTXProducts []MTXProduct `json:"MTXProducts"`
}

type StartPurchaseRequest struct {
	Language  string        `json:"Language"`
	PlayerID  PlayerID      `json:"PlayerID"`
	CartItems []MTXCartItem `json:"CartItems"`
}

type ClaimEntitlementsRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
	AuthCode string   `json:"AuthCode"`
}

type ClaimEntitlementsResponse struct {
	Products []interface{} `json:"Products"`
}

// GetMTXCatalog retrieves the DLC catalog (eg, starter packs).
func (p *PsyNetRPC) GetMTXCatalog(ctx context.Context, category string) (*GetCatalogResponse, error) {
	request := GetCatalogRequest{
		PlayerID: p.localPlayerID,
		Category: category,
	}

	var result GetCatalogResponse
	err := p.sendRequestSync(ctx, "Microtransaction/GetCatalog v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// StartMTXPurchase initiates a DLC purchase via EGS.
func (p *PsyNetRPC) StartMTXPurchase(ctx context.Context, cartItems []MTXCartItem) error {
	request := StartPurchaseRequest{
		Language:  "INT",
		PlayerID:  p.localPlayerID,
		CartItems: cartItems,
	}

	var result interface{}
	err := p.sendRequestSync(ctx, "Microtransaction/StartPurchase v1", request, &result)
	if err != nil {
		return err
	}
	return nil
}

func (p *PsyNetRPC) ClaimMTXEntitlements(ctx context.Context, authCode string) ([]interface{}, error) {
	request := ClaimEntitlementsRequest{
		PlayerID: p.localPlayerID,
		AuthCode: authCode,
	}

	var result ClaimEntitlementsResponse
	err := p.sendRequestSync(ctx, "Microtransaction/ClaimEntitlements v2", request, &result)
	if err != nil {
		return nil, err
	}
	return result.Products, nil
}
