package rlapi

import "context"

type PartyID string

type PartyInfo struct {
	PartyID         string `json:"PartyID"`
	CreatedAt       int64  `json:"CreatedAt"`
	CreatedByUserID string `json:"CreatedByUserID"`
	JoinID          string `json:"JoinID"`
}

type PartyMember struct {
	PartyID  string `json:"PartyID"`
	UserID   string `json:"UserID"`
	UserName string `json:"UserName"`
	JoinedAt int64  `json:"JoinedAt"`
	Role     string `json:"Role"`
}

type PartyResponse struct {
	Info    PartyInfo     `json:"Info"`
	Members []PartyMember `json:"Members"`
}

type GetPlayerPartyInfoResponse struct {
	Invites []interface{} `json:"Invites"`
}

type CreatePartyRequest struct {
	ForcePartyonix bool `json:"bForcePartyonix"`
}

type JoinPartyRequest struct {
	JoinID  string  `json:"JoinID"`
	PartyID PartyID `json:"PartyID"`
}

type LeavePartyRequest struct {
	PartyID PartyID `json:"PartyID"`
}

type SendPartyInviteRequest struct {
	InviteeID PlayerID `json:"InviteeID"`
	PartyID   PartyID  `json:"PartyID"`
}

type SendPartyJoinRequestRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
}

type ChangePartyOwnerRequest struct {
	NewOwnerID PlayerID `json:"NewOwnerID"`
	PartyID    PartyID  `json:"PartyID"`
}

type KickPartyMembersRequest struct {
	Members    []PlayerID `json:"Members"`
	KickReason int        `json:"KickReason"`
	PartyID    PartyID    `json:"PartyID"`
}

type SendPartyChatMessageRequest struct {
	Message string  `json:"Message"`
	PartyID PartyID `json:"PartyID"`
}

type SendPartyChatMessageResponse struct {
	Success   bool   `json:"Success"`
	MessageID string `json:"MessageID"`
}

type SendPartyMessageRequest struct {
	Message string  `json:"Message"`
	PartyID PartyID `json:"PartyID"`
}

type SendPartyMessageResponse struct {
	Success   bool   `json:"Success"`
	MessageID string `json:"MessageID"`
}

func (p *PsyNetRPC) GetPlayerPartyInfo(ctx context.Context) ([]interface{}, error) {
	var result GetPlayerPartyInfoResponse
	err := p.sendRequestSync(ctx, "Party/GetPlayerPartyInfo v1", emptyRequest{}, &result)
	if err != nil {
		return nil, err
	}
	return result.Invites, nil
}

// CreateParty creates a new party.
func (p *PsyNetRPC) CreateParty(ctx context.Context) (*PartyResponse, error) {
	request := CreatePartyRequest{
		ForcePartyonix: true,
	}

	var result PartyResponse
	err := p.sendRequestSync(ctx, "Party/CreateParty v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// JoinParty joins an existing party by join ID or party ID.
func (p *PsyNetRPC) JoinParty(ctx context.Context, joinID string, partyID PartyID) (*PartyResponse, error) {
	request := JoinPartyRequest{
		JoinID:  joinID,
		PartyID: partyID,
	}

	var result PartyResponse
	err := p.sendRequestSync(ctx, "Party/JoinParty v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// LeaveParty leaves a party by ID.
func (p *PsyNetRPC) LeaveParty(ctx context.Context, partyID PartyID) error {
	request := LeavePartyRequest{
		PartyID: partyID,
	}

	var result interface{}
	err := p.sendRequestSync(ctx, "Party/LeaveParty v1", request, &result)
	if err != nil {
		return err
	}
	return nil
}

// SendPartyInvite sends an invitation to join the party.
func (p *PsyNetRPC) SendPartyInvite(ctx context.Context, inviteeID PlayerID, partyID PartyID) error {
	request := SendPartyInviteRequest{
		InviteeID: inviteeID,
		PartyID:   partyID,
	}

	var result interface{}
	err := p.sendRequestSync(ctx, "Party/SendPartyInvite v2", request, &result)
	if err != nil {
		return err
	}
	return nil
}

// SendPartyJoinRequest sends a request to join another player's party.
func (p *PsyNetRPC) SendPartyJoinRequest(ctx context.Context, playerID PlayerID) error {
	request := SendPartyJoinRequestRequest{
		PlayerID: playerID,
	}

	var result interface{}
	err := p.sendRequestSync(ctx, "Party/SendPartyJoinRequest v1", request, &result)
	if err != nil {
		return err
	}
	return nil
}

// ChangePartyOwner transfers party ownership to another party member.
func (p *PsyNetRPC) ChangePartyOwner(ctx context.Context, newOwnerID PlayerID, partyID PartyID) (*PartyResponse, error) {
	request := ChangePartyOwnerRequest{
		NewOwnerID: newOwnerID,
		PartyID:    partyID,
	}

	var result PartyResponse
	err := p.sendRequestSync(ctx, "Party/ChangePartyOwner v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// KickPartyMembers removes members from the party.
func (p *PsyNetRPC) KickPartyMembers(ctx context.Context, members []PlayerID, kickReason int, partyID PartyID) error {
	request := KickPartyMembersRequest{
		Members:    members,
		KickReason: kickReason,
		PartyID:    partyID,
	}

	var result interface{}
	err := p.sendRequestSync(ctx, "Party/KickPartyMembers v1", request, &result)
	if err != nil {
		return err
	}
	return nil
}

// SendPartyChatMessage sends a chat message to the party.
func (p *PsyNetRPC) SendPartyChatMessage(ctx context.Context, message string, partyID PartyID) error {
	request := SendPartyMessageRequest{
		Message: message,
		PartyID: partyID,
	}

	var result SendPartyChatMessageResponse
	err := p.sendRequestSync(ctx, "Party/SendPartyChatMessage v1", request, &result)
	if err != nil {
		return err
	}
	return nil
}

// SendPartyMessage sends an encoded message to the party.
func (p *PsyNetRPC) SendPartyMessage(ctx context.Context, message string, partyID PartyID) error {
	request := SendPartyMessageRequest{
		Message: message,
		PartyID: partyID,
	}

	var result SendPartyMessageResponse
	err := p.sendRequestSync(ctx, "Party/SendPartyMessage v1", request, &result)
	if err != nil {
		return err
	}
	return nil
}
