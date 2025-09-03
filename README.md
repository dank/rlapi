# _rlapi_
![GitHub Release](https://img.shields.io/github/v/release/dank/rlapi)
[![Go Reference](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/dank/rlapi)
![GitHub License](https://img.shields.io/github/license/dank/rlapi)

### [ITEM SHOP DEMO](https://rl.guac.net)

_rlapi_ is a reverse-engineered collection of Rocket League's internal APIs with a Go SDK. It provides a full end-to-end flow, from authentication to accessing the item shop, player stats, inventory, match history, replays, and more. This repository also contains resources for reverse engineering and analyzing Rocket League network traffic, serving as a foundation for further exploration. Not all endpoints are fully documented—do not ask about specific ones, as I probably don't know.

### Contributions
All contributions are welcome! If you discover new endpoints, extend the Go SDK, or add additional functionality, please submit a PR.

## Getting Started
Refer to the [godoc](https://pkg.go.dev/github.com/dank/rlapi) for detailed documentation on the Go SDK.

Comprehensive examples are available in the [`examples`](examples) directory.

### Usage
```bash
go get github.com/dank/rlapi
```

### Authentication
Rocket League authentication always goes through [Epic Online Services (EOS)](https://dev.epicgames.com/docs/web-api-ref/authentication), either via the Epic Games Store (EGS) or by exchanging a Steam session ticket for an EOS token.

This library provides full end-to-end authentication via EGS. Steam login and ticket generation are out of scope, but a method is provided to exchange a Steam session ticket for an EOS token, and users can leverage external libraries such as [steam-user](https://github.com/DoctorMcKay/node-steam-user/wiki/Steam-App-Auth) or [Steamworks](https://partner.steamgames.com/doc/api/ISteamUser#GetAuthSessionTicket) to obtain the ticket.

## Intercepting Requests
Traditional proxy tools like Fiddler don't work with Rocket League due to certificate pinning.

To intercept traffic, we use Frida dynamic instrumentation to hook curl functions at runtime, disabling SSL verification and redirecting API traffic to a local MITM server.

### MITM Server
The MITM proxy forwards both HTTP and WebSocket traffic to the official servers while intercepting and logging requests and responses. Authentication responses are rewritten so the WebSocket URL points to the local server.

> [!NOTE]
> Refer to the respective READMEs in [`tools/frida/`](tools/frida) and [`tools/mitm/`](tools/mitm) directories for more details and usage instructions.

## Reconstructing Requests
### Authentication
Rocket League initially establishes a connection with the HTTP API before transitioning to a WebSocket connection. The client sends an EOS access token via HTTP and receives session credentials, the WebSocket endpoint URL, and any tokens required for further communication. The client then connects to the WebSocket using these tokens, allowing all subsequent API calls to occur over a persistent WebSocket connection.

### Signing
All API requests and responses must include a `PsySig` header containing a Base64-encoded HMAC-SHA256 signature. The signing keys were reverse-engineered from the game binary and are XOR'd with a 4-byte pattern. To decrypt:
```python
# Raw IDA dump
data = [0x36, 0xEA, 0x37, 0x0C, ...]  # 36 bytes total

key_bytes = [data[i] ^ data[(i % 4) + 32] for i in range(32)]
```

- **Request Key**: `c338bd36fb8c42b1a431d30add939fc7`
  - Format: `HMAC-SHA256(key, "-" + <request body>)`
- **Response Key**: `3b932153785842ac927744b292e40e52`  
  - Format: `HMAC-SHA256(key, PsyTime + "-" + <response body>)`


### Request Protocol
#### WebSocket Schema
WebSocket messages require a custom HTTP-like schema with headers and JSON body:

```json5
PsyService: Matchmaking/StartMatchmaking v2
PsyRequestID: PsyNetMessage_X_1
PsyToken: authentication-token
PsySessionID: session-identifier
PsySig: request-signature
PsyBuildID: 151471783
User-Agent: RL Win/250811.43331.492665 gzip
PsyEnvironment: Prod

{"playlist_id": 10, "region": "USE"}
```

The message format is: headers (each ending with `\r\n`) followed by `\r\n\r\n` separator, then JSON body.

**Required Headers:** _(Values may be outdated)_

| Name           | HTTP | WS | Value                                                                  |                                                                  |
|----------------|------|----|------------------------------------------------------------------------|------------------------------------------------------------------|
| PsyService     |      | ✅  |                                                                        | WS event name                                                    |
| PsyRequestID   | ✅    | ✅  | PsyNetMessage_X_0                                                      | Incrementing idempotency key (also for request/response matching |
| PsyToken       |      | ✅  |                                                                        | WS auth token                                                    |
| PsySessionID   |      | ✅  |                                                                        | WS session identifier                                            |
| PsySig         | ✅    | ✅  |                                                                        | Base64-encoded HMAC signature of the body                        |
| PsyBuildID     | ✅    | ✅  | 151471783                                                              | Varies by build                                                  |
| PsyEnvironment | ✅    | ✅  | Prod                                                                   | Varies by build                                                  |
| FeatureSet     | ✅    |    | PrimeUpdate55_1                                                        | Varies by build                                                  |
| User-Agent     |      | ✅  | RL Win/250811.43331.492665 gzip                                        | Varies by build                                                  |
| User-Agent     | ✅    |    | RL Win/250811.43331.492665 gzip (x86_64-pc-win32) curl-7.67.0 Schannel | Varies by build                                                  |
| Content-Type   | ✅    |    | application/x-www-form-urlencoded                                      | JSON body but form type                                          |

## HTTP Endpoints
The HTTP API is used only for authentication and to bootstrap the WebSocket connection. Unlike before, there is no way to "downgrade" the WebSocket connection to HTTP (AFAIK).

The base URL for HTTP requests is: `https://api.rlpp.psynet.gg/rpc/`.

### POST Auth/AuthPlayer/v2

> [!NOTE]
> When the API refers to a "Player ID", it typically expects the following format:
>  ```
> <platform>|<platform-specific account ID>|0
> ```
> For example, on Steam: `Steam|76561197960287930|0`. The final 0 is always included, though its purpose is unknown.

###### Request
```json5
PsyRequestID: PsyNetMessage_X_0 
PsyBuildID: 151471783
PsyEnvironment: Prod
User-Agent: User-Agent: RL Win/250811.43331.492665 gzip (x86_64-pc-win32) curl-7.67.0 Schannel
PsySig: <HMAC signature>
Content-Type: application/x-www-form-urlencoded
        
{
  "Platform": "<platform>",
  "PlayerName": "<player name>",
  "PlayerID": "<platform account ID>",
  "Language": "INT",
  "AuthTicket": "<EOS access token>",
  "BuildRegion": "",
  "FeatureSet": "PrimeUpdate55_1", // varies by build
  "Device": "PC",
  "LocalFirstPlayerID": "<player ID>",
  "bSkipAuth": false,
  "bSetAsPrimaryAccount": true,
  "EpicAuthTicket": "<EOS auth ticket>",
  "EpicAccountID": "<Epic account ID>"
}
```

###### Response
```json5
{
  "Result": {
    "IsLastChanceAuthBan": false,
    "VerifiedPlayerName": "<player name>",
    "UseWebSocket": true, // forcing this to false doesn't do anything
    "PerConURL": "wss://...", // unused / legacy?
    "PerConURLv2": "wss://...", // url for ws connection
    "PsyToken": "<auth token>", // used by subsequent ws requests
    "SessionID": "<session id>", // used by subsequent ws requests
    "CountryRestrictions": []
  }
}
```

## WebSocket Events
> [!NOTE]
> For complete request/response schema definitions, refer to the [godoc](https://pkg.go.dev/github.com/dank/rlapi) documentation instead.

The following is an incomplete list of WebSocket events, some events may be undocumented or partially understood.

The WebSocket endpoint URL is returned by the authentication endpoint, it is currently: `wss://ws.rlpp.psynet.gg/ws/gc2`.

### Challenges
#### Challenges/CollectReward v1
#### Challenges/FTECheckpointComplete v1
#### Challenges/FTEGroupComplete v1
#### Challenges/GetActiveChallenges v1
Retrieves all available challenges (quests/objectives) for the authenticated player.

###### Request
```json5
{
  "Challenges": [],    // Optional: specific challenge IDs to fetch
  "Folders": []        // Optional: specific challenge folders to fetch
}
```

###### Response
```json5
{
  "Result": {
    "Challenges": [
      {
        "ID": 387,
        "Title": "New Driver Challenge",
        "Description": "Complete the Basic Tutorial in the Training Playlist",
        "Sort": 0,
        "GroupID": 8,           // Challenge group/category
        "XPUnlockLevel": 0,     // Level required to unlock
        "bIsRepeatable": false,
        "RepeatLimit": 0,
        "IconURL": "https://rl-cdn.psyonix.com/ChallengeIcons/Challenge_Play.jpg",
        "BackgroundURL": null,
        "BackgroundColor": 0,
        "Requirements": [
          {
            "RequiredCount": 1   // How many times to complete the objective
          }
        ],
        "Rewards": {
          "XP": 0,              // XP reward amount
          "Currency": [],       // Currency rewards (credits, etc.)
          "Products": [         // Item rewards
            {
              "ID": "861",
              "ChallengeID": 387,
              "ProductID": 29,    // Item product ID
              "InstanceID": null,
              "Attributes": [],
              "SeriesID": 861
            }
          ],
          "Pips": 0            // Battle Pass tier progress
        },
        "bAutoClaimRewards": false,
        "bIsPremium": false,   // Requires Rocket Pass Premium
        "UnlockChallengeIDs": []  // Challenge IDs unlocked by completing this
      }
      // ... additional challenges
    ]
  }
}
```

#### Challenges/PlayerProgress v1

### Clubs
#### Clubs/AcceptClubInvite v2
#### Clubs/CreateClub v1
#### Clubs/GetClubDetails v1
#### Clubs/GetClubInvites v1
#### Clubs/GetClubTitleInstances v1
#### Clubs/GetPlayerClubDetails v2
Retrieves detailed information about the club that a specific player belongs to.

###### Request
```json5
{
  "PlayerID": "<platform>|<player_id>|0"
}
```

###### Response
```json5
{
  "Result": {
    "ClubDetails": {
      "ClubID": 30859631,
      "ClubName": "<club_name>",
      "ClubTag": "<tag>",                  // Club tag (3-4 characters)
      "PrimaryColor": 0,                   // Club primary color ID
      "AccentColor": 0,                    // Club accent color ID
      "EquippedTitle": "Club_Supersonic_Acrobatic_Battle_Cars",  // Club title
      "OwnerPlayerID": "<platform>|<player_id>|0",
      "Members": [
        {
          "PlayerID": "<platform>|<player_id>|0",
          "PlayerName": "<player_name>",
          "EpicPlayerID": "<platform>|<player_id>|0",
          "EpicPlayerName": "<epic_name>",
          "RoleID": 1,                     // Role (0=Member, 1=Owner, 2=Officer)
          "CreatedTime": 1750299549,       // When player joined
          "DeletedTime": 0,                // When player left (0 = still member)
          "PsyonixID": null
        }
        // ... additional members
      ],
      "Badges": [                         // Club achievement badges
        {
          "Stat": "Goal",                 // Badge type
          "Badge": 2                      // Badge level
        }
        // ... additional badges
      ],
      "Flags": [],                        // Club flags/moderation status
      "bVerified": false,                 // Official verification status
      "CreatedTime": 1750299549,
      "LastUpdatedTime": 1750299580,
      "NameLastUpdatedTime": 0,
      "DeletedTime": 0
    }
  }
}
```
#### Clubs/GetStats v1
#### Clubs/InviteToClub v4
#### Clubs/LeaveClub v1
#### Clubs/RejectClubInvite v1
#### Clubs/UpdateClub v2

### Drop
#### Drop/GetTradeInFilters v1
Retrieves available trade-in categories and their eligible item series.

###### Request
```json5
{}
```

###### Response
```json5
{
  "Result": {
    "TradeInFilters": [
      {
        "ID": 1,
        "Label": "Core Items",              // Category name
        "SeriesIDs": [1, 47, 191, 207, 300, 443, 541, 542, 635, 902],  // Eligible series
        "bBlueprint": false,               // Whether this is for blueprints
        "TradeInQualities": [              // Eligible item rarities
          "Uncommon",
          "Rare", 
          "VeryRare",
          "Import",
          "Exotic"
        ]
      },
      {
        "ID": 2,
        "Label": "Tournament Items",
        "SeriesIDs": [855, 1147, 1204, 1761, 2281, 2717],
        "bBlueprint": false,
        "TradeInQualities": ["Uncommon", "Rare", "VeryRare", "Import", "Exotic"]
      },
      {
        "ID": 3,
        "Label": "Blueprints",
        "SeriesIDs": [],                   // All series eligible for blueprint trade-ins
        "bBlueprint": true,
        "TradeInQualities": ["Rare", "VeryRare", "Import", "Exotic"]
      }
    ]
  }
}
```

### Filters
#### Filters/FilterContent v1
Filters text content for profanity and inappropriate language using Rocket League's content policy.

###### Request
```json5
{
  "Content": ["text1", "text2", "text3"],  // Array of strings to filter
  "Policy": "Content"                       // Filter policy type
}
```

###### Response
```json5
{
  "Result": {
    "FilteredContent": ["text1", "text2", "text3"]  // Filtered/censored content
  }
}
```

### GameServer
#### GameServer/GetClubPrivateMatches v1
#### GameServer/GetGameServerPingList v2
Retrieves ping measurements to available game server regions.

###### Request
```json5
{}
```

###### Response
```json5
{
  "Result": {
    "Regions": [
      {
        "Region": "USE",             // Region code
        "Label": "US-East",         // Display name
        "SubRegions": ["USE1", "USE3"]  // Available server clusters
      },
      {
        "Region": "EU",
        "Label": "Europe",
        "SubRegions": ["EU5", "EU1", "EU3", "EU7", "EU9"]
      },
      {
        "Region": "OCE",
        "Label": "Oceania",
        "SubRegions": ["OCE1"]
      }
      // ... additional regions
    ]
  }
}
```

### Matches
#### Matches/GetMatchHistory v1

### Matchmaking
#### Matchmaking/PlayerCancelMatchmaking v1
#### Matchmaking/PlayerSearchPrivateMatch v1
Searches for available private matches in a specific region and playlist.

###### Request
```json5
{
  "Region": "USE1",                   // Server region
  "PlaylistID": 6                     // Playlist ID to search within
}
```

###### Response
```json5
{
  "Result": {
    // Private match results would be returned here
    // Empty result indicates no matches found
  }
}
```
#### Matchmaking/StartMatchmaking v2
Initiates matchmaking for specified playlists and regions.

###### Request
```json5
{
  "Regions": [
    {
      "Name": "USE1",     // US East region
      "Ping": 33          // Ping to region in ms
    },
    {
      "Name": "USE3", 
      "Ping": 33
    }
  ],
  "Playlists": [11],              // Playlist IDs (11 = Ranked 2v2)
  "SecondsSearching": 1,          // How long already searching
  "CurrentServerID": "",          // Current server if reconnecting
  "bDisableCrossplay": false,
  "PartyID": "<party_id>",         // Party identifier
  "PartyMembers": [              // All party member player IDs
    "<platform>|<player_id>|0"
  ]
}
```

###### Response
```json5
{
  "Result": {
    "EstimatedQueueTime": 32    // Estimated wait time in seconds
  }
}
```

### Metrics
#### Metrics/RecordMetrics v1

### Microtransaction
#### Microtransaction/ClaimEntitlements v2
#### Microtransaction/GetCatalog v1
Retrieves available DLC/starter pack products for purchase.

###### Request
```json5
{
  "PlayerID": "<platform>|<player_id>|0",
  "Category": "StarterPack"    // Category filter ("StarterPack", etc.)
}
```

###### Response
```json5
{
  "Result": {
    "MTXProducts": [
      {
        "ID": 139,
        "Title": "Season 19 Veteran Pack",
        "Description": "LIMITED TIME",
        "TabTitle": "Season 19 Veteran Pack",
        "PriceDescription": "",
        "ImageURL": "",
        "PlatformProductID": "<platform_product_id>",  // Platform store ID
        "bIsOwned": false,
        "Items": [              // Items included in the pack
          {
            "ProductID": 4284,
            "InstanceID": null,
            "Attributes": [
              {
                "Key": "Painted",
                "Value": "9"      // Paint color ID
              }
            ],
            "SeriesID": 8365
          }
          // ... additional items
        ],
        "Currencies": [         // Currency rewards included
          {
            "ID": 13,           // Currency ID (13 = Credits)
            "CurrencyID": 13,
            "Amount": 500        // Amount granted
          }
        ]
      }
    ]
  }
}
```

#### Microtransaction/StartPurchase v1

### Party
#### Party/ChangePartyOwner v1
#### Party/CreateParty v1
#### Party/GetPlayerPartyInfo v1
Retrieves current party information and pending invitations for the authenticated player.

###### Request
```json5
{}
```

###### Response
```json5
{
  "Result": {
    "Invites": []    // Array of pending party invitations
  }
}
```

#### Party/JoinParty v1
#### Party/KickPartyMembers v1
#### Party/LeaveParty v1
#### Party/SendPartyChatMessage v1
#### Party/SendPartyInvite v2
#### Party/SendPartyJoinRequest v1
#### Party/SendPartyMessage v1

### Players
#### Players/GetBanStatus v3
#### Players/GetCreatorCode v1
#### Players/GetProfile v1
Retrieves basic profile information and presence status for multiple players.

###### Request
```json5
{
  "PlayerIDs": [
    "<platform>|<player_id>|0",
    "<platform>|<player_id>|0"
    // ... additional player IDs
  ]
}
```

###### Response
```json5
{
  "Result": {
    "PlayerData": [
      {
        "PlayerID": "<platform>|<player_id>|0",
        "PlayerName": "<player_name>",
        "PresenceState": "Online",    // "Online", "Offline", "Away", etc.
        "PresenceInfo": ""           // Additional presence details
      }
      // ... additional players
    ]
  }
}
```

#### Players/GetXP v1
Retrieves the authenticated player's XP level and progress information.

###### Request
```json5
{
  "PlayerID": "<platform>|<player_id>|0"
}
```

###### Response
```json5
{
  "Result": {
    "XPInfoResponse": {
      "TotalXP": 3255694,                    // Total XP earned across all time
      "XPLevel": 173,                       // Current XP level
      "XPTitle": "",                        // Title unlocked at this level
      "XPProgressInCurrentLevel": 5694,     // XP progress within current level
      "XPRequiredForNextLevel": 20000       // XP needed to reach next level
    }
  }
}
```
#### Players/Report v4

### Playlists
#### Playlists/GetActivePlaylists v1
Retrieves all currently available playlists (casual and ranked) with their availability windows.

###### Request
```json5
{}
```

###### Response
```json5
{
  "Result": {
    "CasualPlaylists": [
      {
        "NodeID": "OnesCasual",
        "Playlist": 1,           // Playlist ID (1 = Casual 1v1)
        "Type": 1,              // Playlist type
        "StartTime": null,      // Unix timestamp or null if always available
        "EndTime": null
      },
      {
        "NodeID": "ArcadeCasual1",
        "Playlist": 50,         // Limited time playlist
        "Type": 3,
        "StartTime": 1756310400,
        "EndTime": 1757001600
      }
      // ... additional casual playlists
    ],
    "RankedPlaylists": [
      {
        "NodeID": "OnesCompetitive",
        "Playlist": 10,         // Playlist ID (10 = Ranked 1v1)
        "Type": 1,
        "StartTime": null,
        "EndTime": null
      }
      // ... additional ranked playlists
    ],
    "XPLevelUnlocked": 20      // Level required to unlock ranked
  }
}
```

### Population
#### Population/GetPopulation v1
Retrieves current player counts across all playlists.

###### Request
```json5
{}
```

###### Response
```json5
{
  "Result": {
    "Playlists": [
      {
        "Playlist": 10,        // Playlist ID (10 = Ranked 1v1)
        "PlayerCount": 10615    // Current players in queue/matches
      },
      {
        "Playlist": 11,        // Ranked 2v2
        "PlayerCount": 90079
      },
      {
        "Playlist": 13,        // Ranked 3v3
        "PlayerCount": 29329
      }
      // ... additional playlists
    ]
  }
}
```

#### Population/UpdatePlayerPlaylist v1

### Products
#### Products/CrossEntitlement/GetProductStatus v1
#### Products/GetContainerDropTable v2
#### Products/GetPlayerProducts v2
Retrieves a player's inventory including all owned items with their attributes and metadata.

###### Request
```json5
{
  "PlayerID": "<platform>|<player_id>|0",
  "UpdatedTimestamp": "<timestamp>"    // Optional: only return items updated after this timestamp
}
```

###### Response
```json5
{
  "Result": {
    "ProductData": [
      {
        "ProductID": 11854,                    // Item type ID
        "InstanceID": "<instance_id>", // Unique item instance
        "Attributes": [],                      // Item modifiers (painted, certified, etc.)
        "SeriesID": 8350,                     // Item series/collection
        "AddedTimestamp": 1755057379,         // When item was obtained
        "UpdatedTimestamp": 1755387884,
        "DeletedTimestamp": 1755387884        // When item was removed/traded
      },
      {
        "ProductID": 7173,
        "InstanceID": "<instance_id>",
        "Attributes": [
          {
            "Key": "Quality",                  // Blueprint quality
            "Value": "Rare"
          },
          {
            "Key": "Blueprint",                // Blueprint product ID
            "Value": 6127
          },
          {
            "Key": "BlueprintCost",           // Cost to build in credits
            "Value": "100"
          }
        ],
        "SeriesID": 4,
        "AddedTimestamp": 1751433293,
        "UpdatedTimestamp": 1756577976,
        "DeletedTimestamp": 1756577976
      },
      {
        "ProductID": 7076,
        "InstanceID": "<instance_id>",
        "Attributes": [
          {
            "Key": "Painted",                 // Paint color ID
            "Value": 11                       // 11 = Purple
          },
          {
            "Key": "Quality",
            "Value": "Import"
          },
          {
            "Key": "Blueprint",
            "Value": 7073
          },
          {
            "Key": "BlueprintCost",
            "Value": "500"
          }
        ],
        "SeriesID": 4,
        "AddedTimestamp": 1755399374,
        "UpdatedTimestamp": 1755399374
      }
      // ... additional inventory items
    ]
  }
}
```

#### Products/TradeIn v2
#### Products/UnlockContainer v2

### Regions
#### Regions/GetSubRegions v1
Retrieves all available server regions and their sub-regions.

###### Request
```json5
{
  "RequestRegions": [],    // Optional: specific regions to query
  "Regions": []            // Optional: region filter
}
```

###### Response
```json5
{
  "Result": {
    "Regions": [
      {
        "Region": "USE",                    // Region code
        "Label": "US-East",                 // Display name
        "SubRegions": ["USE1", "USE3"]      // Available server clusters
      },
      {
        "Region": "EU",
        "Label": "Europe",
        "SubRegions": ["EU5", "EU1", "EU3", "EU7", "EU9"]
      },
      {
        "Region": "OCE",
        "Label": "Oceania",
        "SubRegions": ["OCE1"]
      }
      // ... additional regions
    ]
  }
}
```

### Reservations
#### Reservations/JoinMatch v1
Attempts to join a private match by server name and password.

###### Request
```json5
{
  "JoinType": "JoinPrivate",           // Join type ("JoinPrivate", etc.)
  "ServerName": "<server_name>",       // Private match server name
  "Password": "<password>"             // Server password
}
```

###### Response (Success)
```json5
{
  "Result": {
    // Match connection details would be returned here
  }
}
```

###### Response (Error)
```json5
{
  "Error": {
    "Type": "ServerNotFound",          // Error type
    "Message": ""                      // Error message
  }
}
```

### Rocket Pass
#### RocketPass/GetPlayerInfo v2
Retrieves the authenticated player's Rocket Pass progress and available purchase options.

###### Request
```json5
{
  "PlayerID": "<platform>|<player_id>|0",
  "RocketPassID": 25,             // Current season's Rocket Pass ID
  "RocketPassInfo": {},           // Additional info filters
  "RocketPassStore": {}           // Store info filters
}
```

###### Response
```json5
{
  "Result": {
    "StartTime": 1750258800,        // Season start timestamp
    "EndTime": 1758031200,          // Season end timestamp
    "RocketPassInfo": {
      "TierLevel": 74,              // Current tier level
      "bOwnsPremium": false,       // Whether player owns premium pass
      "XPMultiplier": 0,           // XP boost multiplier
      "Pips": 730,                 // Progress within current tier
      "PipsPerLevel": 10            // Pips required per tier
    },
    "RocketPassStore": {
      "Tiers": [                   // Tier skip purchase options
        {
          "PurchasableID": 144,
          "CurrencyID": 13,         // Currency type (13 = Credits)
          "CurrencyCost": 200,      // Cost in credits
          "OriginalCurrencyCost": null,
          "Tiers": 1,               // Number of tiers to skip
          "Savings": 0,             // Discount percentage
          "ImageUrl": null
        }
        // ... additional tier skip options
      ],
      "Bundles": [                 // Premium pass purchase options
        {
          "PurchasableID": 142,
          "CurrencyID": 13,
          "CurrencyCost": 1000,     // Premium pass cost
          "OriginalCurrencyCost": null,
          "Tiers": 0,
          "Savings": 0,
          "ImageUrl": "https://rl-cdn.psyonix.com/RocketPass/Images/S19/..."
        }
        // ... additional bundle options
      ]
    }
  }
}
```

#### RocketPass/GetPlayerPrestigeRewards v1
#### RocketPass/GetRewardContent v1

### Shops
#### Shops/GetPlayerWallet v1
Retrieves the authenticated player's currency balances (credits, tokens, etc.).

###### Request
```json5
{
  "PlayerID": "<platform>|<player_id>|0"
}
```

###### Response
```json5
{
  "Result": {
    "Currencies": [
      {
        "ID": 13,                    // Currency ID (13 = Credits)
        "Amount": 0,                 // Current balance
        "ExpirationTime": null,     // Expiry timestamp (null = no expiry)
        "UpdatedTimestamp": 1752883359,
        "IsTradable": false,        // Can be traded to other players
        "TradeHold": null           // Trade restriction timestamp
      }
      // ... additional currencies
    ]
  }
}
```
#### Shops/GetShopCatalogue v2
Retrieves available items and their prices from specified shop catalogues.

###### Request
```json5
{
  "ShopIDs": [52, 397, 354, 382, 220, 51, 55, 357, 358, 359, 360, 361, 362, 363, 364, 365, 366, 367, 368]
}
```

###### Response
```json5
{
  "Result": {
    "Catalogues": [
      {
        "ShopID": 354,
        "ShopItems": [
          {
            "ShopItemID": 12387,
            "StartDate": 1756771200,        // Unix timestamp
            "EndDate": 1757203200,          // Unix timestamp
            "MaxQuantityPerPlayer": 1,
            "ImageURL": null,
            "DeliverableProducts": [        // Items included in purchase
              {
                "Count": 1,
                "Product": {
                  "ProductID": 11499,
                  "InstanceID": null,
                  "Attributes": [],
                  "SeriesID": 1
                },
                "SortID": 1,
                "IsOwned": false
              }
              // ... additional products
            ],
            "DeliverableCurrencies": [],    // Currency rewards if any
            "Costs": [                      // Price options
              {
                "ResetTime": null,
                "ShopItemCostID": 23839,
                "Discount": null,
                "BulkDiscounts": null,
                "StartDate": 1756771200,
                "EndDate": 1757203200,
                "Price": [
                  {
                    "ID": 13,              // Currency ID (13 = Credits)
                    "Amount": 1500         // Price in credits
                  }
                ],
                "SortID": 1,
                "DisplayTypeID": 0
              }
            ],
            "Title": "ADVENTURE TIME + MAMBA",
            "Description": "BUNDLE",
            "Purchasable": true,
            "PurchasedQuantity": 0
          }
          // ... additional shop items
        ]
      }
      // ... additional catalogues
    ]
  }
}
```

#### Shops/GetShopNotifications v1
#### Shops/GetStandardShops v1

### Skills
#### Skills/GetPlayerSkill v1
Retrieves skill data (rank, MMR, etc.) for a specific player across all playlists.

###### Request
```json5
{
  "PlayerID": "<platform>|<player_id>|0"
}
```

###### Response
```json5
{
  "Result": {
    "Skills": [
      {
        "Playlist": 10,        // 1v1 Duel
        "Mu": 30.4646,        // TrueSkill Mu value
        "Sigma": 2.5,         // TrueSkill Sigma value
        "Tier": 11,           // Rank tier (0-22)
        "Division": 1,        // Division within tier (0-3)
        "MMR": 30.4646,       // Matchmaking Rating
        "WinStreak": 3,       // Current win/loss streak
        "MatchesPlayed": 81,
        "PlacementMatchesPlayed": 10
      }
      // ... additional playlists
    ],
    "RewardLevels": {
      "SeasonLevel": 6,
      "SeasonLevelWins": 0
    }
  }
}
```

#### Skills/GetPlayersSkills v1
Query arbitrary player's skills.
#### Skills/GetSkillLeaderboard v1
#### Skills/GetSkillLeaderboardRankForUsers v1
#### Skills/GetSkillLeaderboardValueForUser v1

### Stats
#### Stats/GetStatLeaderboard v1
#### Stats/GetStatLeaderboardRankForUsers v1
#### Stats/GetStatLeaderboardValueForUser v1

### Tournaments
#### Tournaments/Registration/RegisterTournament v1
#### Tournaments/Registration/UnsubscribeTournament v1
#### Tournaments/Search/GetPublicTournaments v1
#### Tournaments/Search/GetSchedule v1
Retrieves scheduled tournaments for a specific region.

###### Request
```json5
{
  "PlayerID": "<platform>|<player_id>|0",
  "Region": "USE"    // Region code (USE, EU, etc.)
}
```

###### Response
```json5
{
  "Result": {
    "Schedules": [
      {
        "Time": 1756843200,        // Tournament start time (Unix timestamp)
        "ScheduleID": 39143,
        "bUpdateSkill": false,     // Whether tournament affects MMR
        "Tournaments": [
          {
            "ID": 44287528,
            "Title": "2v2 Pentathlon",
            "CreatorName": "Psyonix",
            "CreatorPlayerID": "Steam|0|0",      // System/official creator
            "StartTime": 1756843200,
            "GenerateBracketTime": null,
            "MaxBracketSize": 32,
            "RankMin": 0,            // Minimum rank requirement
            "RankMax": 22,           // Maximum rank requirement
            "Region": "USC",
            "Platforms": ["Steam", "PS4", "XboxOne", "Switch", "Epic"],
            "GameTags": "",
            "GameMode": 27,          // Game mode ID (27 = Pentathlon)
            "GameModes": [12, 25, 6, 8, 0],  // Specific game modes for pentathlon
            "TeamSize": 2,
            "MapSetName": null,
            "DisabledMaps": [],
            "SeriesLength": 1,
            "FinalSeriesLength": 3,
            "SeriesRoundLengths": [3, 3, 1],
            "SeedingType": 2,
            "TieBreaker": 0,
            "bPublic": false,
            "TeamsRegistered": 0,
            "ScheduleID": 39143,
            "IsSchedulingTournament": true
          }
          // ... additional tournaments in this time slot
        ]
      }
      // ... additional schedules
    ]
  }
}
```

#### Tournaments/Status/GetCycleData v1
#### Tournaments/Status/GetScheduleRegion v1
#### Tournaments/Status/GetTournamentSubscriptions v1

### Training
#### Training/BrowseTrainingData v1
Browses available training packs with filtering options.

###### Request
```json5
{
  // Request structure varies - can include filters for difficulty, creator, etc.
}
```

###### Response
```json5
{
  "Result": {
    "TrainingData": [
      {
        "Code": "4CA7-FADD-0DF1-AEC2",      // Training pack code
        "TM_Name": "Diamond Pack May 2023",   // Pack name
        "Type": 3,                           // Pack type
        "Difficulty": 2,                     // Difficulty level (0=Rookie, 1=Pro, 2=All-Star)
        "CreatorName": "Psyonix",           // Pack creator
        "MapName": "cs_p",                   // Map used for training
        "Tags": [],                          // Pack tags/categories
        "NumRounds": 9,                      // Number of shots in pack
        "TM_Guid": "<guid>",                 // Internal pack identifier
        "CreatedAt": 1683788495,             // Creation timestamp
        "UpdatedAt": 1756388418              // Last update timestamp
      }
      // ... additional training packs
    ]
  }
}
```

#### Training/GetTrainingMetadata v1
Retrieves metadata for specific training packs by their codes.

###### Request
```json5
{
  "Codes": ["2BFC-F8D6-22AC-2AFE"]    // Array of training pack codes
}
```

###### Response
```json5
{
  "Result": {
    "TrainingData": [
      {
        "Code": "2BFC-F8D6-22AC-2AFE",
        "TM_Name": "Diamond Pack Nov 2024",
        "Type": 3,
        "Difficulty": 1,
        "CreatorName": "Psyonix",
        "CreatorPlayerID": "",
        "MapName": "cs_p",
        "Tags": [],
        "NumRounds": 13,
        "TM_Guid": "<guid>",
        "CreatedAt": 1732585868,
        "UpdatedAt": 1749745526
      }
    ]
  }
}
```

### Users
#### Users/CanShowAvatar v1
Checks which players from a list are allowed to display avatars (based on privacy settings).

###### Request
```json5
{
  "PlayerIDs": [
    "<platform>|<player_id>|0",
    "<platform>|<player_id>|0"
    // ... additional player IDs to check
  ]
}
```

###### Response
```json5
{
  "Result": {
    "AllowedPlayerIDs": [               // Players who allow avatar display
      "<platform>|<player_id>|0",
      "<platform>|<player_id>|0"
    ],
    "HiddenPlayerIDs": []               // Players who have disabled avatar display
  }
}
```

### Misc
These requests have non-standard message schemas and I don't really know what they do.
#### DSR/RelayToServer v1
Sent when joining a match.
#### Party/System
Related to parties but uses a non-standard schema.
#### PsyPing
Sent every 20 seconds. `PsyPing` header with an empty body.
