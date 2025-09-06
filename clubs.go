package rlapi

import "context"

// ClubDetails represents detailed information about a club
type ClubDetails struct {
	ClubID              int64         `json:"ClubID"`
	ClubName            string        `json:"ClubName"`
	ClubTag             string        `json:"ClubTag"`
	PrimaryColor        int           `json:"PrimaryColor"`
	AccentColor         int           `json:"AccentColor"`
	EquippedTitle       string        `json:"EquippedTitle"`
	OwnerPlayerID       PlayerID      `json:"OwnerPlayerID"`
	Members             []ClubMember  `json:"Members"`
	Badges              []ClubBadge   `json:"Badges"`
	Flags               []interface{} `json:"Flags"`
	BVerified           bool          `json:"bVerified"`
	CreatedTime         int64         `json:"CreatedTime"`
	LastUpdatedTime     int64         `json:"LastUpdatedTime"`
	NameLastUpdatedTime int64         `json:"NameLastUpdatedTime"`
	DeletedTime         int64         `json:"DeletedTime"`
}

// ClubMember represents a member of a club
type ClubMember struct {
	PlayerID       PlayerID `json:"PlayerID"`
	PlayerName     string   `json:"PlayerName"`
	EpicPlayerID   PlayerID `json:"EpicPlayerID"`
	EpicPlayerName string   `json:"EpicPlayerName"`
	RoleID         int      `json:"RoleID"`
	CreatedTime    int64    `json:"CreatedTime"`
	DeletedTime    int64    `json:"DeletedTime"`
	PsyonixID      *string  `json:"PsyonixID"`
}

// ClubBadge represents a badge earned by a club
type ClubBadge struct {
	Stat  string `json:"Stat"`
	Badge int    `json:"Badge"`
}

// ClubInvite represents an invitation to join a club
type ClubInvite struct {
	ClubID     int64    `json:"ClubID"`
	ClubName   string   `json:"ClubName"`
	ClubTag    string   `json:"ClubTag"`
	InviterID  PlayerID `json:"InviterID"`
	InviteTime int64    `json:"InviteTime"`
}

// ClubStats represents statistics for a club
type ClubStats struct {
	ClubID int64                  `json:"ClubID"`
	Stats  map[string]interface{} `json:"Stats"`
}

// ClubTitleInstance represents a title instance for a club
type ClubTitleInstance struct {
	TitleID    string `json:"TitleID"`
	InstanceID string `json:"InstanceID"`
}

type GetClubDetailsRequest struct {
	ClubID int64 `json:"ClubID"`
}

type GetClubDetailsResponse struct {
	ClubDetails ClubDetails `json:"ClubDetails"`
}

type GetPlayerClubDetailsRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
}

type GetPlayerClubDetailsResponse struct {
	ClubDetails ClubDetails `json:"ClubDetails"`
}

type GetClubInvitesResponse struct {
	ClubInvites []ClubInvite `json:"ClubInvites"`
}

type CreateClubRequest struct {
	ClubName     string `json:"ClubName"`
	ClubTag      string `json:"ClubTag"`
	PrimaryColor int    `json:"PrimaryColor"`
	AccentColor  int    `json:"AccentColor"`
}

type CreateClubResponse struct {
	ClubID      int64       `json:"ClubID"`
	ClubDetails ClubDetails `json:"ClubDetails"`
}

type UpdateClubRequest struct {
	ClubID        int64   `json:"ClubID"`
	ClubName      *string `json:"ClubName,omitempty"`
	ClubTag       *string `json:"ClubTag,omitempty"`
	PrimaryColor  *int    `json:"PrimaryColor,omitempty"`
	AccentColor   *int    `json:"AccentColor,omitempty"`
	EquippedTitle *string `json:"EquippedTitle,omitempty"`
}

type UpdateClubResponse struct {
	Success     bool        `json:"Success"`
	ClubDetails ClubDetails `json:"ClubDetails"`
}

type InviteToClubRequest struct {
	ClubID   int64    `json:"ClubID"`
	PlayerID PlayerID `json:"PlayerID"`
}

type InviteToClubResponse struct {
	Success bool `json:"Success"`
}

type AcceptClubInviteRequest struct {
	ClubID int64 `json:"ClubID"`
}

type AcceptClubInviteResponse struct {
	Success     bool        `json:"Success"`
	ClubDetails ClubDetails `json:"ClubDetails"`
}

type LeaveClubRequest struct {
	ClubID int64 `json:"ClubID"`
}

type LeaveClubResponse struct {
	Success bool `json:"Success"`
}

type GetStatsRequest struct {
	ClubID int64 `json:"ClubID"`
}

type GetStatsResponse struct {
	Stats ClubStats `json:"Stats"`
}

type GetClubTitleInstancesRequest struct {
	ClubID int64 `json:"ClubID"`
}

type GetClubTitleInstancesResponse struct {
	TitleInstances []ClubTitleInstance `json:"TitleInstances"`
}

type RejectClubInviteRequest struct {
	ClubID int64 `json:"ClubID"`
}

type RejectClubInviteResponse struct {
	Success bool `json:"Success"`
}

// GetClubDetails retrieves detailed information about a specific club.
func (p *PsyNetRPC) GetClubDetails(ctx context.Context, clubID int64) (*GetClubDetailsResponse, error) {
	request := GetClubDetailsRequest{
		ClubID: clubID,
	}

	var result GetClubDetailsResponse
	err := p.sendRequestSync(ctx, "Clubs/GetClubDetails v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPlayerClubDetails retrieves club details for a specific player.
func (p *PsyNetRPC) GetPlayerClubDetails(ctx context.Context, playerID PlayerID) (*GetPlayerClubDetailsResponse, error) {
	request := GetPlayerClubDetailsRequest{
		PlayerID: playerID,
	}

	var result GetPlayerClubDetailsResponse
	err := p.sendRequestSync(ctx, "Clubs/GetPlayerClubDetails v2", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetClubInvites retrieves pending club invitations for the current player.
func (p *PsyNetRPC) GetClubInvites(ctx context.Context) (*GetClubInvitesResponse, error) {
	var result GetClubInvitesResponse
	err := p.sendRequestSync(ctx, "Clubs/GetClubInvites v1", map[string]interface{}{}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateClub creates a new club with the specified details.
func (p *PsyNetRPC) CreateClub(ctx context.Context, clubName, clubTag string, primaryColor, accentColor int) (*CreateClubResponse, error) {
	request := CreateClubRequest{
		ClubName:     clubName,
		ClubTag:      clubTag,
		PrimaryColor: primaryColor,
		AccentColor:  accentColor,
	}

	var result CreateClubResponse
	err := p.sendRequestSync(ctx, "Clubs/CreateClub v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateClub updates the details of an existing club.
func (p *PsyNetRPC) UpdateClub(ctx context.Context, request UpdateClubRequest) (*UpdateClubResponse, error) {
	var result UpdateClubResponse
	err := p.sendRequestSync(ctx, "Clubs/UpdateClub v2", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// InviteToClub invites a player to join a club.
func (p *PsyNetRPC) InviteToClub(ctx context.Context, clubID int64, playerID PlayerID) (*InviteToClubResponse, error) {
	request := InviteToClubRequest{
		ClubID:   clubID,
		PlayerID: playerID,
	}

	var result InviteToClubResponse
	err := p.sendRequestSync(ctx, "Clubs/InviteToClub v4", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AcceptClubInvite accepts an invitation to join a club.
func (p *PsyNetRPC) AcceptClubInvite(ctx context.Context, clubID int64) (*AcceptClubInviteResponse, error) {
	request := AcceptClubInviteRequest{
		ClubID: clubID,
	}

	var result AcceptClubInviteResponse
	err := p.sendRequestSync(ctx, "Clubs/AcceptClubInvite v2", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// LeaveClub leaves a club.
func (p *PsyNetRPC) LeaveClub(ctx context.Context, clubID int64) (*LeaveClubResponse, error) {
	request := LeaveClubRequest{
		ClubID: clubID,
	}

	var result LeaveClubResponse
	err := p.sendRequestSync(ctx, "Clubs/LeaveClub v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetClubStats retrieves statistics for a club.
func (p *PsyNetRPC) GetClubStats(ctx context.Context, clubID int64) (*GetStatsResponse, error) {
	request := GetStatsRequest{
		ClubID: clubID,
	}

	var result GetStatsResponse
	err := p.sendRequestSync(ctx, "Clubs/GetStats v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetClubTitleInstances retrieves title instances for a club.
func (p *PsyNetRPC) GetClubTitleInstances(ctx context.Context, clubID int64) (*GetClubTitleInstancesResponse, error) {
	request := GetClubTitleInstancesRequest{
		ClubID: clubID,
	}

	var result GetClubTitleInstancesResponse
	err := p.sendRequestSync(ctx, "Clubs/GetClubTitleInstances v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RejectClubInvite rejects an invitation to join a club.
func (p *PsyNetRPC) RejectClubInvite(ctx context.Context, clubID int64) (*RejectClubInviteResponse, error) {
	request := RejectClubInviteRequest{
		ClubID: clubID,
	}

	var result RejectClubInviteResponse
	err := p.sendRequestSync(ctx, "Clubs/RejectClubInvite v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
