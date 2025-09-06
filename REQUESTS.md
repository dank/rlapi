## WebSocket Events
> [!NOTE]
> For complete request/response schema definitions, refer to the [godoc](https://pkg.go.dev/github.com/dank/rlapi) documentation instead.

Request and response bodies are example payloads, not full schemas.

### Challenges
#### Challenges/CollectReward v1
Collects rewards for a completed challenge.

###### Request
```json5
{
  "PlayerID": "<player id>",
  "ID": 3751                           
}
```

###### Response
```json5
{
  "Result": {}                          
}
```

#### Challenges/FTECheckpointComplete v1
#### Challenges/FTEGroupComplete v1
#### Challenges/GetActiveChallenges v1
Retrieves all available challenges for the authenticated player.

###### Request
```json5
{
  "Challenges": [],    
  "Folders": []        
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
        "GroupID": 8,           
        "XPUnlockLevel": 0,    
        "bIsRepeatable": false,
        "RepeatLimit": 0,
        "IconURL": "https://rl-cdn.psyonix.com/ChallengeIcons/Challenge_Play.jpg",
        "BackgroundURL": null,
        "BackgroundColor": 0,
        "Requirements": [
          {
            "RequiredCount": 1 
          }
        ],
        "Rewards": {
          "XP": 0,            
          "Currency": [],       
          "Products": [        
            {
              "ID": "861",
              "ChallengeID": 387,
              "ProductID": 29,  
              "InstanceID": null,
              "Attributes": [],
              "SeriesID": 861
            }
          ],
          "Pips": 0          
        },
        "bAutoClaimRewards": false,
        "bIsPremium": false,   
        "UnlockChallengeIDs": [] 
      }
      // ... additional challenges
    ]
  }
}
```

#### Challenges/PlayerProgress v1
Retrieves progress information for specific challenges.

###### Request
```json5
{
  "PlayerID": "<player id>",
  "ChallengeIDs": [387, 388, 389, 390]  
}
```

###### Response
```json5
{
  "Result": {
    "ProgressData": [
      {
        "ID": 387,                       
        "CompleteCount": 1,            
        "bIsHidden": false,             
        "bNotifyCompleted": false,      
        "bNotifyAvailable": false,       
        "bNotifyNewInfo": false,     
        "bRewardsAvailable": false,     
        "bComplete": true,             
        "RequirementProgress": [          
          {
            "ProgressCount": 0,          
            "ProgressChange": 0           
          }
        ],
        "ProgressResetTimeUTC": 0  
      }
      // ... additional challenge progress data
    ]
  }
}
```

### Clubs
#### Clubs/AcceptClubInvite v2
Accepts a club invitation and joins the club.

###### Request
```json5
{
  "ClubID": 123456
}
```

###### Response
```json5
{
  "Result": {
    "ClubDetails": {
      "ClubID": 123456,
      "ClubName": "<club name>",
      "ClubTag": "<tag>",
      "PrimaryColor": 0,
      "AccentColor": 0,
      "EquippedTitle": "Club_Supersonic_Acrobatic_Battle_Cars",
      "OwnerPlayerID": "<player id>",
      "Members": [
        // ... club member details (same structure as GetClubDetails)
      ],
      "Badges": [
        {"Stat": "Goal", "Badge": 2}
      ],
      "Flags": [],
      "bVerified": false,
      "CreatedTime": 1750299549,
      "LastUpdatedTime": 1750299580,
      "NameLastUpdatedTime": 0,
      "DeletedTime": 0
    }
  }
}
```

#### Clubs/CreateClub v1
Creates a new club with specified name, tag, and colors.

###### Request
```json5
{
  "ClubName": "My New Club",          
  "ClubTag": "TAG",              
  "PrimaryColor": 0,                 
  "AccentColor": 0                 
}
```

###### Response
```json5
{
  "Result": {
    "ClubDetails": {
      // ... same structure as GetPlayerClubDetails response
    }
  }
}
```

#### Clubs/GetClubDetails v1
Retrieves detailed information about a specific club by ID.

###### Request
```json5
{
  "ClubID": 123456 
}
```

###### Response
```json5
{
  "Result": {
    "ClubDetails": {
      // ... same structure as GetPlayerClubDetails response
    }
  }
}
```

#### Clubs/GetClubInvites v1
Retrieves pending club invitations for the authenticated player.

###### Request
```json5
{}
```

###### Response
```json5
{
  "Result": {
    "ClubInvites": []
  }
}
```

#### Clubs/GetClubTitleInstances v1
Retrieves all available club titles that can be equipped.

###### Request
```json5
{}
```

###### Response
```json5
{
  "Result": {
    "ClubTitles": [
      "Club_Supersonic_Acrobatic_Battle_Cars",
      // ... additional titles
    ]
  }
}
```
#### Clubs/GetPlayerClubDetails v2
Retrieves detailed information about the club that a specific player belongs to.

###### Request
```json5
{
  "PlayerID": "<player id>"
}
```

###### Response
```json5
{
  "Result": {
    "ClubDetails": {
      "ClubID": 30859631,
      "ClubName": "<club name>",
      "ClubTag": "<tag>",
      "PrimaryColor": 0,
      "AccentColor": 0,
      "EquippedTitle": "Club_Supersonic_Acrobatic_Battle_Cars",
      "OwnerPlayerID": "<player id>",
      "Members": [
        {
          "PlayerID": "<player id>",
          "PlayerName": "<player name>",
          "EpicPlayerID": "<player id>",
          "EpicPlayerName": "<epic name>",
          "RoleID": 1, 
          "CreatedTime": 1750299549,    
          "DeletedTime": 0,    
          "PsyonixID": null
        }
        // ... additional members
      ],
      "Badges": [ 
        {
          "Stat": "Goal",  
          "Badge": 2
        }
        // ... additional badges
      ],
      "Flags": [], 
      "bVerified": false, 
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
Invites a player to join the club.

###### Request
```json5
{
  "PlayerID": "<player id>"
}
```

###### Response
```json5
{
  "Result": {}
}
```

#### Clubs/LeaveClub v1
Leaves the current club.

###### Request
```json5
{}
```

###### Response
```json5
{
  "Result": {}
}
```

#### Clubs/RejectClubInvite v1
#### Clubs/UpdateClub v2
Updates club settings like colors and other properties.

###### Request
```json5
{
  "PrimaryColor": -10879077,
  "AccentColor": 0 
}
```

###### Response
```json5
{
  "Result": {
    "ClubDetails": {
      // ... same structure as GetPlayerClubDetails response
    }
  }
}
```
#### Clubs/GetStats v1
Retrieves club statistics and achievements for the authenticated player's club.

###### Request
```json5
{}
```

###### Response
```json5
{
  "Result": {
    "CareerStats": {
      "TimePlayed": 181826, 
      "Goal": 1736, 
      "AerialGoal": 190, 
      "LongGoal": 79,
      "BackwardsGoal": 9,   
      "OvertimeGoal": 60,  
      "TurtleGoal": 3, 
      "Assist": 931, 
      "Playmaker": 42,
      "Save": 1564,
      "EpicSave": 406, 
      "Savior": 144, 
      "Shot": 4372, 
      "Center": 5319, 
      "Clear": 3608,
      "AerialHit": 5578, 
      "BicycleHit": 162, 
      "JuggleHit": 33,
      "Demolish": 1048, 
      "Demolition": 0,
      "FirstTouch": 2121, 
      "PoolShot": 6,
      "LowFive": 20, 
      "HighFive": 1, 
      "BreakoutDamage": 0, 
      "BreakoutDamageLarge": 0,
      "HoopsSwishGoal": 4, 
      "MatchPlayed": 662, 
      "Win": 325
    },
    "SeasonalStats": [
      {
        "Stat": "Goal",
        "Milestones": [225, 675, 2025],
        "Value": 1736,
        "Badge": 2 
      }
    ],
    "PreviousSeasonalBadges": [],
    "SeasonalTitles": [
      {
        "Badge": 1,
        "Title": "Club_S19_M1"
      }
    ]
  }
}
```

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
        "Label": "Core Items",
        "SeriesIDs": [1, 47, 191, 207, 300, 443, 541, 542, 635, 902],
        "bBlueprint": false, 
        "TradeInQualities": [
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
        "SeriesIDs": [],
        "bBlueprint": true,
        "TradeInQualities": ["Rare", "VeryRare", "Import", "Exotic"]
      }
    ]
  }
}
```

### Filters
#### Filters/FilterContent v1

### GameServer
#### GameServer/GetClubPrivateMatches v1
Retrieves available private matches hosted by clubs.

###### Request
```json5
{}
```

###### Response
```json5
{
  "Result": {
    "Servers": []
  }
}
```

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
        "Region": "USE", 
        "Label": "US-East",
        "SubRegions": ["USE1", "USE3"]
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
Retrieves recent match history for a specific player including match details, player stats, and replay URLs.

###### Request
```json5
{
  "PlayerID": "<player id>"
}
```

###### Response
```json5
{
  "Result": {
    "Matches": [
      {
        "ReplayUrl": "http://api.rlpp.psynet.gg/Match.replay?MatchGUID=...&Timestamp=...&Expiration=...&Signature=...",
        "Match": {
          "MatchGUID": "<guid>",
          "RecordStartTimestamp": 1756529579,
          "MapName": "Wasteland_P", 
          "Playlist": 6,
          "SecondsPlayed": 14.9421,
          "OvertimeSecondsPlayed": 0,
          "WinningTeam": -1,
          "Team0Score": 0,
          "Team1Score": 0,
          "bOverTime": false,
          "bNoContest": false,
          "bForfeit": false,
          "CustomMatchCreatorPlayerID": "<player id>",
          "bClubVsClub": false, 
          "Mutators": [""],
          "Players": [
            {
              "PlayerID": "<player id>",
              "PlayerName": "<player name>",
              "ConnectTimestamp": 1756529580,
              "JoinTimestamp": 1756529581,
              "LeaveTimestamp": 1756529586,
              "PartyLeaderID": "<player id>",
              "InParty": false,
              "bAbandoned": false,
              "bMvp": false, 
              "LastTeam": 1,
              "TeamColor": "Orange",
              "SecondsPlayed": 0,
              "Score": 0,
              "Goals": 0,
              "Assists": 0,
              "Saves": 0,
              "Shots": 0, 
              "Demolishes": 0, 
              "OwnGoals": 0,
              "Skills": {
                "Mu": null,         
                "Sigma": null,           
                "Division": null,     
                "PrevMu": null,       
                "PrevSigma": null,   
                "PrevTier": null,       
                "PrevDivision": null,     
                "bValid": false   
              }
            }
          ]
        }
      }
    ]
  }
}
```

### Matchmaking
#### Matchmaking/PlayerCancelMatchmaking v1
Cancels the current matchmaking search.

###### Request
```json5
{}
```

###### Response
```json5
{
  "Result": {}
}
```

#### Matchmaking/PlayerSearchPrivateMatch v1
Searches for available private matches in a specific region and playlist.

###### Request
```json5
{
  "Region": "USE1", 
  "PlaylistID": 6
}
```

###### Response
```json5
{
  "Result": {}
}
```
#### Matchmaking/StartMatchmaking v2
Initiates matchmaking for specified playlists and regions.

###### Request
```json5
{
  "Regions": [
    {
      "Name": "USE1", 
      "Ping": 33  
    },
    {
      "Name": "USE3", 
      "Ping": 33
    }
  ],
  "Playlists": [11],     
  "SecondsSearching": 1,  
  "CurrentServerID": "",  
  "bDisableCrossplay": false,
  "PartyID": "<party id>", 
  "PartyMembers": [
    "<player id>"
  ]
}
```

###### Response
```json5
{
  "Result": {
    "EstimatedQueueTime": 32 
  }
}
```

### Metrics
#### Metrics/RecordMetrics v1

### Microtransaction
#### Microtransaction/ClaimEntitlements v2
#### Microtransaction/GetCatalog v1
Retrieves available starter packs for purchase.

###### Request
```json5
{
  "PlayerID": "<player id>",
  "Category": "StarterPack" 
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
        "PlatformProductID": "<product id>",
        "bIsOwned": false,
        "Items": [         
          {
            "ProductID": 4284,
            "InstanceID": null,
            "Attributes": [
              {
                "Key": "Painted",
                "Value": "9"  
              }
            ],
            "SeriesID": 8365
          }
          // ... additional items
        ],
        "Currencies": [
          {
            "ID": 13,
            "CurrencyID": 13,
            "Amount": 500
          }
        ]
      }
    ]
  }
}
```

#### Microtransaction/StartPurchase v1
Initiates a purchase transaction for catalog items.

###### Request
```json5
{
  "Language": "INT",  
  "PlayerID": "<player id>",
  "CartItems": [ 
    {
      "CatalogID": 13, 
      "Count": 1  
    }
  ]
}
```

### Party
#### Party/ChangePartyOwner v1
Changes the owner of a party to another member.

###### Request
```json5
{
  "NewOwnerID": "<player id>",
  "PartyID": "<party id>" 
}
```

###### Response
```json5
{
  "Result": {}
}
```

#### Party/CreateParty v1
Creates a new party.

###### Request
```json5
{
  "bForcePartyonix": true 
}
```

###### Response
```json5
{
  "Result": {
    "PartyID": "<party id>"
    // ... additional party details
  }
}
```

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
    "Invites": [] 
  }
}
```

#### Party/JoinParty v1
Joins an existing party.

###### Request
```json5
{
  "JoinID": "",            
  "PartyID": "<party id>" 
}
```

###### Response
```json5
{
  "Result": {
    // ... party details
  }
}
```

#### Party/KickPartyMembers v1
Kicks members from a party (owner only).

###### Request
```json5
{
  "Members": ["<player id>"], 
  "KickReason": 1,           
  "PartyID": "<party id>" 
}
```

###### Response
```json5
{
  "Result": {}                      
}
```

#### Party/LeaveParty v1
Leaves the current party.

###### Request
```json5
{
  "PartyID": "<party id>"   
}
```

###### Response
```json5
{
  "Result": {}  
}
```

#### Party/SendPartyChatMessage v1
Sends a text chat message to party members.

###### Request
```json5
{
  "Message": "among us", 
  "PartyID": "<party id>" 
}
```

###### Response
```json5
{
  "Result": {}
}
```

#### Party/SendPartyInvite v2
Sends a party invitation to another player.

###### Request
```json5
{
  "InviteeID": "<player id>", 
  "PartyID": "<party id>"
}
```

###### Response
```json5
{
  "Result": {} 
}
```

#### Party/SendPartyJoinRequest v1
Sends a request to join another player's party.

###### Request
```json5
{
  "PlayerID": "<player id>"
}
```

###### Response
```json5
{
  "Result": {}
}
```

#### Party/SendPartyMessage v1
Sends a voice message (?) to party members.

###### Request
```json5
{
  "Message": "",  // Base64-encoded message
  "PartyID": "<party id>"
}
```

###### Response
```json5
{
  "Result": {}
}
```

### Players
#### Players/GetBanStatus v3
Retrieves ban status information for specified players.

###### Request
```json5
{
  "Players": ["<player id>"] 
}
```

###### Response
```json5
{
  "Result": {
    "BanMessages": []
  }
}
```

#### Players/GetCreatorCode v1
Retrieves the player's creator code information.

###### Request
```json5
{}
```

###### Response
```json5
{}
```

#### Players/GetProfile v1
Retrieves basic profile information and presence status for multiple players.

###### Request
```json5
{
  "PlayerIDs": [
    "<player id>",
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
        "PlayerID": "<player id>",
        "PlayerName": "<player name>",
        "PresenceState": "Online",
        "PresenceInfo": ""
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
  "PlayerID": "<player id>"
}
```

###### Response
```json5
{
  "Result": {
    "XPInfoResponse": {
      "TotalXP": 12345,   
      "XPLevel": 123,  
      "XPTitle": "",                    
      "XPProgressInCurrentLevel": 123, 
      "XPRequiredForNextLevel": 123   
    }
  }
}
```
#### Players/Report v4
Reports a player for misconduct.

###### Request
```json5
{
  "Reports": [
    {
      "Reporter": "<player id>",
      "Offender": "<player id>",
      "ReasonIDs": [3],  
      "ReportTimeStamp": 0.0
    }
  ],
  "GameID": "" 
}
```

###### Response
```json5
{
  "Result": {}
}
```

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
        "Playlist": 1, 
        "Type": 1,
        "StartTime": null,
        "EndTime": null
      },
      {
        "NodeID": "ArcadeCasual1",
        "Playlist": 50, 
        "Type": 3,
        "StartTime": 1756310400,
        "EndTime": 1757001600
      }
      // ... additional casual playlists
    ],
    "RankedPlaylists": [
      {
        "NodeID": "OnesCompetitive",
        "Playlist": 10,
        "Type": 1,
        "StartTime": null,
        "EndTime": null
      }
      // ... additional ranked playlists
    ],
    "XPLevelUnlocked": 20 
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
        "Playlist": 10,  
        "PlayerCount": 10615
      },
      {
        "Playlist": 11,
        "PlayerCount": 90079
      },
      {
        "Playlist": 13,
        "PlayerCount": 29329
      }
      // ... additional playlists
    ]
  }
}
```

#### Population/UpdatePlayerPlaylist v1
Updates the player's current playlist for population tracking.

###### Request
```json5
{
  "Playlist": 0,
  "NumLocalPlayers": 1
}
```

###### Response
```json5
{
  "Result": {}
}
```

### Products
#### Products/CrossEntitlement/GetProductStatus v1
#### Products/GetContainerDropTable v2
#### Products/GetPlayerProducts v2
Retrieves a player's inventory including all owned items with their attributes and metadata.

###### Request
```json5
{
  "PlayerID": "<player id>",
  "UpdatedTimestamp": "<timestamp>" 
}
```

###### Response
```json5
{
  "Result": {
    "ProductData": [
      {
        "ProductID": 7076,
        "InstanceID": "<instance id>",
        "Attributes": [
          {
            "Key": "Painted",
            "Value": 11 
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
Trades in multiple items for a higher-tier item.

###### Request
```json5
{
  "PlayerID": "<player id>",
  "ProductInstances": [
    "62d1e4bc3f5b4076bba8ea044a14a36c",
    "b3f60779705f4c8a926a14d6899bad70",
    "e610e81c23904c0696a25e9537aff4ba",
    "a10e3fa58af141b3829fec050889c967",
    "7cf368745031457dbfa7ec8ae4ab316f"
  ]
}
```

###### Response
```json5
{
  "Result": {
    // TODO: Received item details from trade-in
  }
}
```

#### Products/UnlockContainer v2
Unlocks/opens loot containers to receive items.

###### Request
```json5
{
  "PlayerID": "<player id>",
  "InstanceIDs": ["90a79f045cad4556b95eea1270be0e76"],
  "KeyInstanceIDs": []    
}
```

###### Response
```json5
{
  "Result": {
    // TODO: Items received from opening containers
  }
}
```

### Regions
#### Regions/GetSubRegions v1
Retrieves all available server regions and their sub-regions.

###### Request
```json5
{
  "RequestRegions": [],
  "Regions": []  
}
```

###### Response
```json5
{
  "Result": {
    "Regions": [
      {
        "Region": "USE", 
        "Label": "US-East", 
        "SubRegions": ["USE1", "USE3"]
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
Join a private match by server name and password.

###### Request
```json5
{
  "JoinType": "JoinPrivate", 
  "ServerName": "<server name>", 
  "Password": "<password>" 
}
```

### Rocket Pass
#### RocketPass/GetPlayerInfo v2
Retrieves the authenticated player's Rocket Pass progress and available purchase options.

###### Request
```json5
{
  "PlayerID": "<player id>",
  "RocketPassID": 25, 
  "RocketPassInfo": {}, 
  "RocketPassStore": {}
}
```

###### Response
```json5
{
  "Result": {
    "StartTime": 1750258800,
    "EndTime": 1758031200, 
    "RocketPassInfo": {
      "TierLevel": 74,  
      "bOwnsPremium": false, 
      "XPMultiplier": 0, 
      "Pips": 730,   
      "PipsPerLevel": 10     
    },
    "RocketPassStore": {
      "Tiers": [ 
        {
          "PurchasableID": 144,
          "CurrencyID": 13,   
          "CurrencyCost": 200,  
          "OriginalCurrencyCost": null,
          "Tiers": 1,  
          "Savings": 0,    
          "ImageUrl": null
        }
        // ... additional tier skip options
      ],
      "Bundles": [  
        {
          "PurchasableID": 142,
          "CurrencyID": 13,
          "CurrencyCost": 1000, 
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
Retrieves prestige rewards for a player in the specified Rocket Pass.

###### Request
```json5
{
  "PlayerID": "<player id>",
  "RocketPassID": 25
}
```

###### Response
```json5
{
  "Result": {
    // TODO: Prestige reward information
  }
}
```

#### RocketPass/GetRewardContent v1
Retrieves reward content and tier information for a Rocket Pass.

###### Request
```json5
{
  "RocketPassID": 25,
  "TierCap": 0,
  "FreeMaxLevel": 0,
  "PremiumMaxLevel": 0 
}
```

###### Response
```json5
{
  "Result": {
    // TODO: Rocket Pass reward content and tiers
  }
}
```

### Shops
#### Shops/GetPlayerWallet v1
Retrieves the authenticated player's currency balances (credits, tokens, etc.).

###### Request
```json5
{
  "PlayerID": "<player id>"
}
```

###### Response
```json5
{
  "Result": {
    "Currencies": [
      {
        "ID": 13, 
        "Amount": 0,  
        "ExpirationTime": null,
        "UpdatedTimestamp": 1752883359,
        "IsTradable": false, 
        "TradeHold": null
      }
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
            "StartDate": 1756771200, 
            "EndDate": 1757203200,  
            "MaxQuantityPerPlayer": 1,
            "ImageURL": null,
            "DeliverableProducts": [ 
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
            ],
            "DeliverableCurrencies": [],
            "Costs": [
              {
                "ResetTime": null,
                "ShopItemCostID": 23839,
                "Discount": null,
                "BulkDiscounts": null,
                "StartDate": 1756771200,
                "EndDate": 1757203200,
                "Price": [
                  {
                    "ID": 13, 
                    "Amount": 1500  
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
Retrieves shop-related notifications and alerts.

###### Request
```json5
{}
```

###### Response
```json5
{
  "Result": {
    // TODO: Shop notification data
  }
}
```

#### Shops/GetStandardShops v1
Retrieves information about standard/permanent shop offerings.

###### Request
```json5
{}
```

###### Response
```json5
{
  "Result": {
    // TODO: Standard shop information
  }
}
```

### Skills
#### Skills/GetPlayerSkill v1
Retrieves skill data (rank, MMR, etc.) for a specific player across all playlists.

###### Request
```json5
{
  "PlayerID": "<player id>"
}
```

###### Response
```json5
{
  "Result": {
    "Skills": [
      {
        "Playlist": 10,
        "Mu": 30.4646,  
        "Sigma": 2.5,
        "Tier": 11,  
        "Division": 1,
        "MMR": 30.4646,
        "WinStreak": 3, 
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
Retrieves skill data for multiple players.

###### Request
```json5
{
  "PlayerIDs": ["<player id>"] 
}
```

###### Response
```json5
{
  "Result": {
    // TODO: Skills data for each requested player (same format as GetPlayerSkill)
  }
}
```

#### Skills/GetSkillLeaderboard v1
Retrieves skill-based leaderboard for a specific playlist.

###### Request
```json5
{
  "Playlist": 10,   
  "bDisableCrossplay": false  
}
```

###### Response
```json5
{
  "Result": {
    // TODO: Skill leaderboard data for the playlist
  }
}
```

#### Skills/GetSkillLeaderboardRankForUsers v1
Retrieves leaderboard ranks for specific players in a playlist.

###### Request
```json5
{
  "Playlist": 30, 
  "PlayerIDs": [
    "<player id>"
  ]
}
```

###### Response
```json5
{
  "Result": {
    // TODO: Leaderboard rank data for each requested player
  }
}
```

#### Skills/GetSkillLeaderboardValueForUser v1
Retrieves leaderboard value/position for a specific player.

###### Request
```json5
{
  "Playlist": 10,
  "PlayerID": "<player id>"
}
```

###### Response
```json5
{
  "Result": {
    // TODO: Player's leaderboard position and value
  }
}
```

### Stats
#### Stats/GetStatLeaderboard v1
Retrieves global leaderboard data for a specific statistic.

###### Request
```json5
{
  "Stat": "Wins", 
  "bDisableCrossplay": false 
}
```

###### Response
```json5
{
  "Result": {
    "LeaderboardID": "Wins",
    "Platforms": [
      {
        "Platform": "Epic", 
        "Stat": "Wins",
        "Players": [ 
          {
            "PlayerID": "<player id>",
            "PlayerName": "<player name>",
            "Value": 35128
          }
          // ... additional players
        ]
      },
      {
        "Platform": "PS4",
        "Stat": "Wins", 
        "Players": [
          // ...
        ]
      }
      // ... additional platforms (Steam, XboxOne, Switch)
    ]
  }
}
```

#### Stats/GetStatLeaderboardRankForUsers v1
Retrieves leaderboard ranks for specific players in a stat category.

###### Request
```json5
{
  "Stat": "Wins", 
  "PlayerIDs": [     
    "<player id>"
  ],
  "bDisableCrossplay": false              // Whether to exclude crossplay results
}
```

###### Response
```json5
{
  "Result": {
    // TODO: Stat leaderboard rank data for each requested player
  }
}
```

#### Stats/GetStatLeaderboardValueForUser v1
Retrieves leaderboard value/position for a specific player in a stat category.

###### Request
```json5
{
  "Stat": "Wins",                   
  "PlayerID": "<player id>",
  "bDisableCrossplay": false  
}
```

###### Response
```json5
{
  "Result": {
    // TODO: Player's stat leaderboard position and value
  }
}
```

### Tournaments
#### Tournaments/Registration/RegisterTournament v1
Registers the authenticated player for a tournament.

###### Request
```json5
{
  "PlayerID": "<player id>",
  "TournamentID": "44289515", 
  "Credentials": {
    "Title": "",  
    "Password": ""     
  }
}
```

###### Response
```json5
{
  "Result": {
    "Tournament": {
      "ID": 44289515,
      "Title": "2v2",
      "CreatorName": "Psyonix",
      "CreatorPlayerID": "Steam|0|0",
      "StartTime": 1756850400,
      "GenerateBracketTime": null,
      "MaxBracketSize": 32, 
      "RankMin": 0,
      "RankMax": 22, 
      "Region": "USC", 
      "Platforms": ["Steam", "PS4", "XboxOne", "Switch", "Epic"],
      "GameTags": "",
      "GameMode": 0, 
      "GameModes": [], 
      "TeamSize": 2,
      "MapSetName": null,
      "DisabledMaps": [
        "ARC_P",
        "EuroStadium_SnowNight_P",
        "labs_circlepillars_p"
        // ... additional disabled maps
      ],
      "SeriesLength": 1,
      "FinalSeriesLength": 3,
      "SeriesRoundLengths": [3, 3, 1],
      "SeedingType": 2, 
      "TieBreaker": 0, 
      "bPublic": false,
      "TeamsRegistered": 0,
      "ScheduleID": 39147,
      "IsSchedulingTournament": true
    }
  }
}
```

#### Tournaments/Registration/UnsubscribeTournament v1
#### Tournaments/Search/GetPublicTournaments v1
#### Tournaments/Search/GetSchedule v1
Retrieves scheduled tournaments for a specific region.

###### Request
```json5
{
  "PlayerID": "<player id>",
  "Region": "USE" 
}
```

###### Response
```json5
{
  "Result": {
    "Schedules": [
      {
        "Time": 1756843200, 
        "ScheduleID": 39143,
        "bUpdateSkill": false,
        "Tournaments": [
          {
            "ID": 44287528,
            "Title": "2v2 Pentathlon",
            "CreatorName": "Psyonix",
            "CreatorPlayerID": "Steam|0|0", 
            "StartTime": 1756843200,
            "GenerateBracketTime": null,
            "MaxBracketSize": 32,
            "RankMin": 0, 
            "RankMax": 22,
            "Region": "USC",
            "Platforms": ["Steam", "PS4", "XboxOne", "Switch", "Epic"],
            "GameTags": "",
            "GameMode": 27, 
            "GameModes": [12, 25, 6, 8, 0], 
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
Retrieves tournament cycle information including weekly and seasonal data.

###### Request
```json5
{
  "PlayerID": "<player id>"
}
```

###### Response
```json5
{
  "Result": {
    "CycleID": 19, 
    "CycleEndTime": 1757916000,  
    "WeekID": 11,  
    "WeekEndTime": 1757311200,  
    "WeeklyCurrencies": [], 
    "Weeks": [        
      {"Results": []},
      {"Results": []}
      // ... additional weeks
    ],
    "TournamentCurrencyID": 35 
  }
}
```

#### Tournaments/Status/GetScheduleRegion v1
Retrieves the player's current tournament schedule region.

###### Request
```json5
{
  "PlayerID": "<player id>"
}
```

###### Response
```json5
{
  "Result": {
    "ScheduleRegion": "USE"
  }
}
```

#### Tournaments/Status/GetTournamentSubscriptions v1

### Training
#### Training/BrowseTrainingData v1
Browses available training packs with filtering options.

###### Request
```json5
{
  // TODO: Request structure varies - can include filters for difficulty, creator, etc.
}
```

###### Response
```json5
{
  "Result": {
    "TrainingData": [
      {
        "Code": "4CA7-FADD-0DF1-AEC2", 
        "TM_Name": "Diamond Pack May 2023", 
        "Type": 3,
        "Difficulty": 2, 
        "CreatorName": "Psyonix", 
        "MapName": "cs_p", 
        "Tags": [],     
        "NumRounds": 9,
        "TM_Guid": "<guid>", 
        "CreatedAt": 1683788495, 
        "UpdatedAt": 1756388418
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
  "Codes": ["2BFC-F8D6-22AC-2AFE"]
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
Checks which players from a list are allowed to display avatars.

###### Request
```json5
{
  "PlayerIDs": [
    "<player id>",
    "<player id>"
    // ... additional player IDs to check
  ]
}
```

###### Response
```json5
{
  "Result": {
    "AllowedPlayerIDs": [
      "<player id>",
      "<player id>"
    ],
    "HiddenPlayerIDs": []
  }
}
```

### Misc
These requests have non-standard message schemas and I don't really know what they do.
#### DSR/RelayToServer v1
Sent when joining a server.
#### Party/System
Non-standard schema.
#### PsyPing
Sent every 20 seconds. `PsyPing` header with an empty body.