package rlapi

import "context"

type MatchEntry struct {
	ReplayUrl string `json:"ReplayUrl"`
	Match     Match  `json:"Match"`
}

type Match struct {
	MatchGUID                  string        `json:"MatchGUID"`
	RecordStartTimestamp       int64         `json:"RecordStartTimestamp"`
	MapName                    string        `json:"MapName"`
	Playlist                   int           `json:"Playlist"`
	SecondsPlayed              float64       `json:"SecondsPlayed"`
	OvertimeSecondsPlayed      float64       `json:"OvertimeSecondsPlayed"`
	WinningTeam                int           `json:"WinningTeam"`
	Team0Score                 int           `json:"Team0Score"`
	Team1Score                 int           `json:"Team1Score"`
	OverTime                   bool          `json:"bOverTime"`
	NoContest                  bool          `json:"bNoContest"`
	Forfeit                    bool          `json:"bForfeit"`
	CustomMatchCreatorPlayerID string        `json:"CustomMatchCreatorPlayerID,omitempty"`
	ClubVsClub                 bool          `json:"bClubVsClub"`
	Mutators                   []string      `json:"Mutators"`
	Players                    []MatchPlayer `json:"Players"`
}

type MatchPlayer struct {
	PlayerID         string      `json:"PlayerID"`
	PlayerName       string      `json:"PlayerName"`
	ConnectTimestamp int64       `json:"ConnectTimestamp"`
	JoinTimestamp    int64       `json:"JoinTimestamp"`
	LeaveTimestamp   int64       `json:"LeaveTimestamp"`
	PartyLeaderID    string      `json:"PartyLeaderID"`
	InParty          bool        `json:"InParty"`
	Abandoned        bool        `json:"bAbandoned"`
	MVP              bool        `json:"bMvp"`
	LastTeam         int         `json:"LastTeam"`
	TeamColor        string      `json:"TeamColor"`
	SecondsPlayed    float64     `json:"SecondsPlayed"`
	Score            int         `json:"Score"`
	Goals            int         `json:"Goals"`
	Assists          int         `json:"Assists"`
	Saves            int         `json:"Saves"`
	Shots            int         `json:"Shots"`
	Demolishes       int         `json:"Demolishes"`
	OwnGoals         int         `json:"OwnGoals"`
	Skills           MatchSkills `json:"Skills"`
}

type MatchSkills struct {
	Mu           float64 `json:"Mu"`
	Sigma        float64 `json:"Sigma"`
	Tier         int     `json:"Tier"`
	Division     int     `json:"Division"`
	PrevMu       float64 `json:"PrevMu"`
	PrevSigma    float64 `json:"PrevSigma"`
	PrevTier     int     `json:"PrevTier"`
	PrevDivision int     `json:"PrevDivision"`
	Valid        bool    `json:"bValid"`
}

type GetMatchHistoryRequest struct {
	PlayerID PlayerID `json:"PlayerID"`
}

type GetMatchHistoryResponse struct {
	Matches []MatchEntry `json:"Matches"`
}

// GetMatchHistory retrieves match history for the authenticated player.
func (p *PsyNetRPC) GetMatchHistory(ctx context.Context) ([]MatchEntry, error) {
	request := GetMatchHistoryRequest{
		PlayerID: p.localPlayerID,
	}

	var result GetMatchHistoryResponse
	err := p.sendRequestSync(ctx, "Matches/GetMatchHistory v1", request, &result)
	if err != nil {
		return nil, err
	}
	return result.Matches, nil
}
