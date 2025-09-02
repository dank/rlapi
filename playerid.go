package rlapi

import (
	"fmt"
	"strings"
)

// PlayerID represents a unique player identifier in the format "Platform|ID|0"
type PlayerID string

// Platform represents a gaming platform
type Platform string

// Platform constants
const (
	PlatformEpic   Platform = "Epic"
	PlatformSteam  Platform = "Steam"
	PlatformPS4    Platform = "PS4"
	PlatformXbox   Platform = "XboxOne"
	PlatformSwitch Platform = "Switch"
)

// String returns the string representation of Platform
func (p Platform) String() string {
	return string(p)
}

// NewPlayerID creates a PlayerID for the specified platform and ID
func NewPlayerID(platform Platform, id string) PlayerID {
	return PlayerID(fmt.Sprintf("%s|%s|0", platform, id))
}

// ParsePlayerID parses a PlayerID string and returns its components
func ParsePlayerID(playerID string) (platform Platform, id string, err error) {
	parts := strings.Split(playerID, "|")
	if len(parts) != 3 {
		return "", "", fmt.Errorf("invalid PlayerID format: %s", playerID)
	}

	if parts[2] != "0" {
		return "", "", fmt.Errorf("invalid PlayerID suffix: expected '0', got '%s'", parts[2])
	}

	return Platform(parts[0]), parts[1], nil
}

// String returns the string representation of PlayerID
func (p PlayerID) String() string {
	return string(p)
}

// Platform returns the platform component of the PlayerID
func (p PlayerID) Platform() (Platform, error) {
	platform, _, err := ParsePlayerID(string(p))
	return platform, err
}

// ID returns the ID component of the PlayerID
func (p PlayerID) ID() (string, error) {
	_, id, err := ParsePlayerID(string(p))
	return id, err
}

// IsValid checks if the PlayerID has a valid format
func (p PlayerID) IsValid() bool {
	_, _, err := ParsePlayerID(string(p))
	return err == nil
}

// IsEpic returns true if this is an Epic Games PlayerID
func (p PlayerID) IsEpic() bool {
	platform, _ := p.Platform()
	return platform == PlatformEpic
}

// IsSteam returns true if this is a Steam PlayerID
func (p PlayerID) IsSteam() bool {
	platform, _ := p.Platform()
	return platform == PlatformSteam
}

// IsPS4 returns true if this is a PlayStation 4 PlayerID
func (p PlayerID) IsPS4() bool {
	platform, _ := p.Platform()
	return platform == PlatformPS4
}

// IsXbox returns true if this is an Xbox One PlayerID
func (p PlayerID) IsXbox() bool {
	platform, _ := p.Platform()
	return platform == PlatformXbox
}

// IsSwitch returns true if this is a Nintendo Switch PlayerID
func (p PlayerID) IsSwitch() bool {
	platform, _ := p.Platform()
	return platform == PlatformSwitch
}
