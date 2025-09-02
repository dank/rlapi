package rlapi

import "context"

// PartyInfo represents information about a player's party
type PartyInfo struct {
	PartyID   string        `json:"PartyID"`
	OwnerID   PlayerID      `json:"OwnerID"`
	Members   []PartyMember `json:"Members"`
	Invites   []PartyInvite `json:"Invites"`
	CreatedAt int64         `json:"CreatedAt"`
	Settings  PartySettings `json:"Settings"`
}

// PartyMember represents a member of a party
type PartyMember struct {
	PlayerID   PlayerID `json:"PlayerID"`
	PlayerName string   `json:"PlayerName"`
	Platform   string   `json:"Platform"`
	JoinedAt   int64    `json:"JoinedAt"`
	IsReady    bool     `json:"IsReady"`
}

// PartyInvite represents an invitation to join a party
type PartyInvite struct {
	InviterID   PlayerID `json:"InviterID"`
	InviterName string   `json:"InviterName"`
	InviteeID   PlayerID `json:"InviteeID"`
	PartyID     string   `json:"PartyID"`
	SentAt      int64    `json:"SentAt"`
}

// PartySettings represents party configuration settings
type PartySettings struct {
	MaxMembers   int  `json:"MaxMembers"`
	IsPrivate    bool `json:"IsPrivate"`
	AllowInvites bool `json:"AllowInvites"`
}

// Request and Response types

type GetPlayerPartyInfoResponse struct {
	Invites []PartyInvite `json:"Invites"`
}

type createPartyRequest struct {
	BForcePartyonix bool `json:"bForcePartyonix"`
}

type CreatePartyResponse struct {
	PartyID   string    `json:"PartyID"`
	PartyInfo PartyInfo `json:"PartyInfo"`
}

type joinPartyRequest struct {
	JoinID  string `json:"JoinID"`
	PartyID string `json:"PartyID"`
}

type JoinPartyResponse struct {
	Success   bool      `json:"Success"`
	PartyInfo PartyInfo `json:"PartyInfo"`
}

type leavePartyRequest struct {
	PartyID string `json:"PartyID"`
}

type LeavePartyResponse struct {
	Success bool `json:"Success"`
}

type sendPartyInviteRequest struct {
	InviteeID PlayerID `json:"InviteeID"`
	PartyID   string   `json:"PartyID"`
}

type SendPartyInviteResponse struct {
	Success bool `json:"Success"`
}

type sendPartyJoinRequestRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
}

type SendPartyJoinRequestResponse struct {
	Success bool `json:"Success"`
}

type changePartyOwnerRequest struct {
	NewOwnerID PlayerID `json:"NewOwnerID"`
	PartyID    string   `json:"PartyID"`
}

type ChangePartyOwnerResponse struct {
	Success   bool      `json:"Success"`
	PartyInfo PartyInfo `json:"PartyInfo"`
}

type kickPartyMembersRequest struct {
	Members    []PlayerID `json:"Members"`
	KickReason int        `json:"KickReason"`
	PartyID    string     `json:"PartyID"`
}

type KickPartyMembersResponse struct {
	Success   bool      `json:"Success"`
	PartyInfo PartyInfo `json:"PartyInfo"`
}

type sendPartyChatMessageRequest struct {
	Message string `json:"Message"`
	PartyID string `json:"PartyID"`
}

type SendPartyChatMessageResponse struct {
	Success   bool   `json:"Success"`
	MessageID string `json:"MessageID"`
}

type sendPartyMessageRequest struct {
	Message string `json:"Message"`
	PartyID string `json:"PartyID"`
}

type SendPartyMessageResponse struct {
	Success   bool   `json:"Success"`
	MessageID string `json:"MessageID"`
}

// GetPlayerPartyInfo retrieves party information for the current player.
func (p *PsyNetRPC) GetPlayerPartyInfo(ctx context.Context) (*GetPlayerPartyInfoResponse, error) {
	var result GetPlayerPartyInfoResponse
	err := p.sendRequestSync(ctx, "Party/GetPlayerPartyInfo v1", map[string]interface{}{}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateParty creates a new party.
func (p *PsyNetRPC) CreateParty(ctx context.Context, forcePartyonix bool) (*CreatePartyResponse, error) {
	request := createPartyRequest{
		BForcePartyonix: forcePartyonix,
	}

	var result CreatePartyResponse
	err := p.sendRequestSync(ctx, "Party/CreateParty v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// JoinParty joins an existing party.
func (p *PsyNetRPC) JoinParty(ctx context.Context, joinID, partyID string) (*JoinPartyResponse, error) {
	request := joinPartyRequest{
		JoinID:  joinID,
		PartyID: partyID,
	}

	var result JoinPartyResponse
	err := p.sendRequestSync(ctx, "Party/JoinParty v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// LeaveParty leaves the current party.
func (p *PsyNetRPC) LeaveParty(ctx context.Context, partyID string) (*LeavePartyResponse, error) {
	request := leavePartyRequest{
		PartyID: partyID,
	}

	var result LeavePartyResponse
	err := p.sendRequestSync(ctx, "Party/LeaveParty v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SendPartyInvite sends an invitation to join the party.
func (p *PsyNetRPC) SendPartyInvite(ctx context.Context, inviteeID PlayerID, partyID string) (*SendPartyInviteResponse, error) {
	request := sendPartyInviteRequest{
		InviteeID: inviteeID,
		PartyID:   partyID,
	}

	var result SendPartyInviteResponse
	err := p.sendRequestSync(ctx, "Party/SendPartyInvite v2", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SendPartyJoinRequest sends a request to join another player's party.
func (p *PsyNetRPC) SendPartyJoinRequest(ctx context.Context, playerID PlayerID) (*SendPartyJoinRequestResponse, error) {
	request := sendPartyJoinRequestRequest{
		PlayerID: playerID,
	}

	var result SendPartyJoinRequestResponse
	err := p.sendRequestSync(ctx, "Party/SendPartyJoinRequest v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ChangePartyOwner transfers party ownership to another member.
func (p *PsyNetRPC) ChangePartyOwner(ctx context.Context, newOwnerID PlayerID, partyID string) (*ChangePartyOwnerResponse, error) {
	request := changePartyOwnerRequest{
		NewOwnerID: newOwnerID,
		PartyID:    partyID,
	}

	var result ChangePartyOwnerResponse
	err := p.sendRequestSync(ctx, "Party/ChangePartyOwner v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// KickPartyMembers removes members from the party.
func (p *PsyNetRPC) KickPartyMembers(ctx context.Context, members []PlayerID, kickReason int, partyID string) (*KickPartyMembersResponse, error) {
	request := kickPartyMembersRequest{
		Members:    members,
		KickReason: kickReason,
		PartyID:    partyID,
	}

	var result KickPartyMembersResponse
	err := p.sendRequestSync(ctx, "Party/KickPartyMembers v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SendPartyChatMessage sends a chat message to the party.
func (p *PsyNetRPC) SendPartyChatMessage(ctx context.Context, message, partyID string) (*SendPartyChatMessageResponse, error) {
	request := sendPartyChatMessageRequest{
		Message: message,
		PartyID: partyID,
	}

	var result SendPartyChatMessageResponse
	err := p.sendRequestSync(ctx, "Party/SendPartyChatMessage v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SendPartyMessage sends a system message to the party.
func (p *PsyNetRPC) SendPartyMessage(ctx context.Context, message, partyID string) (*SendPartyMessageResponse, error) {
	request := sendPartyMessageRequest{
		Message: message,
		PartyID: partyID,
	}

	var result SendPartyMessageResponse
	err := p.sendRequestSync(ctx, "Party/SendPartyMessage v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
