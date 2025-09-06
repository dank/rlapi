package rlapi

import (
	"context"
	"strings"
)

type ClubID int

// ClubDetails represents detailed information about a club
type ClubDetails struct {
	ClubID              ClubID        `json:"ClubID"`
	ClubName            string        `json:"ClubName"`
	ClubTag             string        `json:"ClubTag"`
	PrimaryColor        int           `json:"PrimaryColor"`
	AccentColor         int           `json:"AccentColor"`
	EquippedTitle       string        `json:"EquippedTitle"`
	OwnerPlayerID       PlayerID      `json:"OwnerPlayerID"`
	Members             []ClubMember  `json:"Members"`
	Badges              []ClubBadge   `json:"Badges"`
	Flags               []interface{} `json:"Flags"`
	Verified            bool          `json:"bVerified"`
	CreatedTime         int           `json:"CreatedTime"`
	LastUpdatedTime     int           `json:"LastUpdatedTime"`
	NameLastUpdatedTime int           `json:"NameLastUpdatedTime"`
	DeletedTime         int           `json:"DeletedTime"`
}

// ClubMember represents a member of a club
type ClubMember struct {
	PlayerID       PlayerID `json:"PlayerID"`
	PlayerName     string   `json:"PlayerName"`
	EpicPlayerID   PlayerID `json:"EpicPlayerID"`
	EpicPlayerName string   `json:"EpicPlayerName"`
	RoleID         int      `json:"RoleID"`
	CreatedTime    int      `json:"CreatedTime"`
	DeletedTime    int      `json:"DeletedTime"`
	PsyonixID      *string  `json:"PsyonixID"`
}

// ClubBadge represents a badge earned by a club
type ClubBadge struct {
	Stat  string `json:"Stat"`
	Badge int    `json:"Badge"`
}

// ClubInvite represents an invitation to join a club
type ClubInvite struct {
	ClubID        ClubID `json:"ClubID"`
	ClubName      string `json:"ClubName"`
	ClubTag       string `json:"ClubTag"`
	InvitedByID   string `json:"PlayerID"`
	InvitedByName string `json:"PlayerName"`
	EpicPlayerID  string `json:"EpicPlayerID"`
}

// ClubCareerStats represents career statistics for a club member
type ClubCareerStats struct {
	TimePlayed          int `json:"TimePlayed"`
	Goal                int `json:"Goal"`
	AerialGoal          int `json:"AerialGoal"`
	LongGoal            int `json:"LongGoal"`
	BackwardsGoal       int `json:"BackwardsGoal"`
	OvertimeGoal        int `json:"OvertimeGoal"`
	TurtleGoal          int `json:"TurtleGoal"`
	Assist              int `json:"Assist"`
	Playmaker           int `json:"Playmaker"`
	Save                int `json:"Save"`
	EpicSave            int `json:"EpicSave"`
	Savior              int `json:"Savior"`
	Shot                int `json:"Shot"`
	Center              int `json:"Center"`
	Clear               int `json:"Clear"`
	AerialHit           int `json:"AerialHit"`
	BicycleHit          int `json:"BicycleHit"`
	JuggleHit           int `json:"JuggleHit"`
	Demolish            int `json:"Demolish"`
	Demolition          int `json:"Demolition"`
	FirstTouch          int `json:"FirstTouch"`
	PoolShot            int `json:"PoolShot"`
	LowFive             int `json:"LowFive"`
	HighFive            int `json:"HighFive"`
	BreakoutDamage      int `json:"BreakoutDamage"`
	BreakoutDamageLarge int `json:"BreakoutDamageLarge"`
	HoopsSwishGoal      int `json:"HoopsSwishGoal"`
	MatchPlayed         int `json:"MatchPlayed"`
	Win                 int `json:"Win"`
}

// ClubSeasonalStat represents a seasonal statistic with milestones
type ClubSeasonalStat struct {
	Stat       string `json:"Stat"`
	Milestones []int  `json:"Milestones"`
	Value      int    `json:"Value"`
	Badge      int    `json:"Badge"`
}

// ClubSeasonalTitle represents a seasonal title
type ClubSeasonalTitle struct {
	Badge int    `json:"Badge"`
	Title string `json:"Title"`
}

// ClubStatsData represents the complete statistics data for a club member
type ClubStatsData struct {
	CareerStats            ClubCareerStats     `json:"CareerStats"`
	SeasonalStats          []ClubSeasonalStat  `json:"SeasonalStats"`
	PreviousSeasonalBadges []interface{}       `json:"PreviousSeasonalBadges"`
	SeasonalTitles         []ClubSeasonalTitle `json:"SeasonalTitles"`
}

type GetClubDetailsRequest struct {
	ClubID ClubID `json:"ClubID"`
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
	ClubDetails ClubDetails `json:"ClubDetails"`
}

type UpdateClubRequest struct {
	ClubName      *string `json:"ClubName,omitempty"`
	ClubTag       *string `json:"ClubTag,omitempty"`
	PrimaryColor  *int    `json:"PrimaryColor,omitempty"`
	AccentColor   *int    `json:"AccentColor,omitempty"`
	EquippedTitle *string `json:"EquippedTitle,omitempty"`
}

type UpdateClubResponse struct {
	ClubDetails ClubDetails `json:"ClubDetails"`
}

type InviteToClubRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
}

type AcceptClubInviteRequest struct {
	ClubID ClubID `json:"ClubID"`
}

type AcceptClubInviteResponse struct {
	Success     bool        `json:"Success"`
	ClubDetails ClubDetails `json:"ClubDetails"`
}

type LeaveClubRequest struct {
	ClubID ClubID `json:"ClubID"`
}

type GetStatsResponse struct {
	Stats ClubStatsData `json:"Stats"`
}

type GetClubTitleInstancesResponse struct {
	ClubTitles []string `json:"ClubTitles"`
}

type RejectClubInviteRequest struct {
	ClubID ClubID `json:"ClubID"`
}

// GetClubDetails retrieves detailed information about a specific club.
func (p *PsyNetRPC) GetClubDetails(ctx context.Context, clubID ClubID) (*ClubDetails, error) {
	request := GetClubDetailsRequest{
		ClubID: clubID,
	}

	var result GetClubDetailsResponse
	err := p.sendRequestSync(ctx, "Clubs/GetClubDetails v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result.ClubDetails, nil
}

// GetPlayerClubDetails retrieves club details for a specific player.
func (p *PsyNetRPC) GetPlayerClubDetails(ctx context.Context, playerID PlayerID) (*ClubDetails, error) {
	request := GetPlayerClubDetailsRequest{
		PlayerID: playerID,
	}

	var result GetPlayerClubDetailsResponse
	err := p.sendRequestSync(ctx, "Clubs/GetPlayerClubDetails v2", request, &result)
	if err != nil {
		return nil, err
	}
	return &result.ClubDetails, nil
}

// CreateClub creates a new club with the specified details.
func (p *PsyNetRPC) CreateClub(ctx context.Context, clubName, clubTag string, primaryColor, accentColor int) (*ClubDetails, error) {
	request := CreateClubRequest{
		ClubName:     clubName,
		ClubTag:      strings.ToUpper(clubTag),
		PrimaryColor: primaryColor,
		AccentColor:  accentColor,
	}

	var result CreateClubResponse
	err := p.sendRequestSync(ctx, "Clubs/CreateClub v1", request, &result)
	if err != nil {
		return nil, err
	}
	return &result.ClubDetails, nil
}

// UpdateClub updates the details of an existing club.
func (p *PsyNetRPC) UpdateClub(ctx context.Context, request *UpdateClubRequest) (*ClubDetails, error) {
	var result UpdateClubResponse
	err := p.sendRequestSync(ctx, "Clubs/UpdateClub v2", request, &result)
	if err != nil {
		return nil, err
	}
	return &result.ClubDetails, nil
}

// GetClubStats retrieves statistics for the current player's club.
func (p *PsyNetRPC) GetClubStats(ctx context.Context) (*ClubStatsData, error) {
	var result GetStatsResponse
	err := p.sendRequestSync(ctx, "Clubs/GetStats v1", &emptyRequest{}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Stats, nil
}

// GetClubTitleInstances retrieves title instances for a club.
func (p *PsyNetRPC) GetClubTitleInstances(ctx context.Context) ([]string, error) {
	var result GetClubTitleInstancesResponse
	err := p.sendRequestSync(ctx, "Clubs/GetClubTitleInstances v1", &emptyRequest{}, &result)
	if err != nil {
		return nil, err
	}
	return result.ClubTitles, nil
}

// LeaveClub leaves a club.
func (p *PsyNetRPC) LeaveClub(ctx context.Context) error {
	var result interface{}
	err := p.sendRequestSync(ctx, "Clubs/LeaveClub v1", &emptyRequest{}, &result)
	if err != nil {
		return err
	}
	return nil
}

// GetClubInvites retrieves pending club invitations for the current player.
func (p *PsyNetRPC) GetClubInvites(ctx context.Context) ([]ClubInvite, error) {
	var result GetClubInvitesResponse
	err := p.sendRequestSync(ctx, "Clubs/GetClubInvites v1", &emptyRequest{}, &result)
	if err != nil {
		return nil, err
	}
	return result.ClubInvites, nil
}

// InviteToClub invites a player to join a club.
// Note: ClubID is inferred from the current player's club context
func (p *PsyNetRPC) InviteToClub(ctx context.Context, playerID PlayerID) error {
	request := InviteToClubRequest{
		PlayerID: playerID,
	}

	var result interface{}
	err := p.sendRequestSync(ctx, "Clubs/InviteToClub v4", request, &result)
	if err != nil {
		return err
	}
	return nil
}

// AcceptClubInvite accepts an invitation to join a club.
func (p *PsyNetRPC) AcceptClubInvite(ctx context.Context, clubID ClubID) (*ClubDetails, error) {
	request := AcceptClubInviteRequest{
		ClubID: clubID,
	}

	var result AcceptClubInviteResponse
	err := p.sendRequestSync(ctx, "Clubs/AcceptClubInvite v2", request, &result)
	if err != nil {
		return nil, err
	}
	return &result.ClubDetails, nil
}

// RejectClubInvite rejects an invitation to join a club.
func (p *PsyNetRPC) RejectClubInvite(ctx context.Context, clubID ClubID) error {
	request := RejectClubInviteRequest{
		ClubID: clubID,
	}

	var result interface{}
	err := p.sendRequestSync(ctx, "Clubs/RejectClubInvite v1", request, &result)
	if err != nil {
		return err
	}
	return nil
}
