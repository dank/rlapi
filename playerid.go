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

	return Platform(parts[0]), parts[1], nil
}

// String returns the string representation of PlayerID
func (p PlayerID) String() string {
	return string(p)
}
