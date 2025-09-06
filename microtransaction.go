package rlapi

import "context"

// Microtransaction represents a microtransaction item
type Microtransaction struct {
	TransactionID string                 `json:"TransactionID"`
	ProductID     int                    `json:"ProductID"`
	Price         float64                `json:"Price"`
	Currency      string                 `json:"Currency"`
	Status        string                 `json:"Status"`
	Metadata      map[string]interface{} `json:"Metadata"`
}

// CatalogItem represents an item in the microtransaction catalog
type CatalogItem struct {
	ItemID       string                 `json:"ItemID"`
	Name         string                 `json:"Name"`
	Description  string                 `json:"Description"`
	Price        float64                `json:"Price"`
	Currency     string                 `json:"Currency"`
	ImageURL     string                 `json:"ImageURL"`
	Category     string                 `json:"Category"`
	IsAvailable  bool                   `json:"IsAvailable"`
	ProductData  []ProductData          `json:"ProductData"`
	Requirements map[string]interface{} `json:"Requirements"`
}

// Entitlement represents a claimed entitlement
type Entitlement struct {
	EntitlementID string        `json:"EntitlementID"`
	ProductID     int           `json:"ProductID"`
	ClaimedAt     int64         `json:"ClaimedAt"`
	Products      []ProductData `json:"Products"`
	Status        string        `json:"Status"`
}

// Purchase represents a microtransaction purchase
type Purchase struct {
	PurchaseID    string                 `json:"PurchaseID"`
	TransactionID string                 `json:"TransactionID"`
	ItemID        string                 `json:"ItemID"`
	Status        string                 `json:"Status"`
	CreatedAt     int64                  `json:"CreatedAt"`
	CompletedAt   *int64                 `json:"CompletedAt"`
	Metadata      map[string]interface{} `json:"Metadata"`
}

type GetCatalogRequest struct {
	Category string `json:"Category,omitempty"`
	Region   string `json:"Region,omitempty"`
}

type GetCatalogResponse struct {
	Items     []CatalogItem `json:"Items"`
	Timestamp int64         `json:"Timestamp"`
}

type StartPurchaseRequest struct {
	ItemID   string                 `json:"ItemID"`
	Quantity int                    `json:"Quantity"`
	Metadata map[string]interface{} `json:"Metadata,omitempty"`
}

type StartPurchaseResponse struct {
	PurchaseID    string `json:"PurchaseID"`
	TransactionID string `json:"TransactionID"`
	Status        string `json:"Status"`
	RedirectURL   string `json:"RedirectURL,omitempty"`
}

type ClaimEntitlementsRequest struct {
	PlayerID       PlayerID `json:"PlayerID"`
	EntitlementIDs []string `json:"EntitlementIDs,omitempty"`
}

type ClaimEntitlementsResponse struct {
	ClaimedEntitlements []Entitlement `json:"ClaimedEntitlements"`
	NewProducts         []ProductData `json:"NewProducts"`
}

// GetCatalog retrieves the microtransaction catalog.
func (p *PsyNetRPC) GetCatalog(ctx context.Context, category, region string) (*GetCatalogResponse, error) {
	request := GetCatalogRequest{
		Category: category,
		Region:   region,
	}

	var result GetCatalogResponse
	err := p.sendRequestSync(ctx, "Microtransaction/GetCatalog v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// StartPurchase initiates a microtransaction purchase.
func (p *PsyNetRPC) StartPurchase(ctx context.Context, itemID string, quantity int, metadata map[string]interface{}) (*StartPurchaseResponse, error) {
	request := StartPurchaseRequest{
		ItemID:   itemID,
		Quantity: quantity,
		Metadata: metadata,
	}

	var result StartPurchaseResponse
	err := p.sendRequestSync(ctx, "Microtransaction/StartPurchase v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ClaimEntitlements claims available entitlements for a player.
func (p *PsyNetRPC) ClaimEntitlements(ctx context.Context, playerID PlayerID, entitlementIDs []string) (*ClaimEntitlementsResponse, error) {
	request := ClaimEntitlementsRequest{
		PlayerID:       playerID,
		EntitlementIDs: entitlementIDs,
	}

	var result ClaimEntitlementsResponse
	err := p.sendRequestSync(ctx, "Microtransaction/ClaimEntitlements v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
