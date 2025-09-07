# _rlapi_
![GitHub Release](https://img.shields.io/github/v/release/dank/rlapi)
[![Go Reference](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/dank/rlapi)
![GitHub License](https://img.shields.io/github/license/dank/rlapi)

### [ITEM SHOP DEMO](https://rl.guac.net)

_rlapi_ is a reverse engineered collection of Rocket League's internal APIs with a Go SDK. It provides a full end-to-end flow, from authentication to accessing the item shop, player stats, inventory, match history, replays, and more. This repository also contains resources for reverse engineering and analyzing Rocket League network traffic, serving as a foundation for further exploration. Not all endpoints are fully documented—do not ask about specific ones, as I probably don't know.

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
All API requests and responses must include a `PsySig` header containing a Base64-encoded HMAC-SHA256 signature. The signing keys were reverse engineered from the game binary and are XOR'd with a 4-byte pattern. To decrypt:
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

#### Required Headers
_(Values may be outdated)_

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

- Challenges
    - [Challenges/CollectReward v1](REQUESTS.md#challengescollectreward-v1)
    - [Challenges/FTECheckpointComplete v1](REQUESTS.md#challengesftecheckpointcomplete-v1)
    - [Challenges/FTEGroupComplete v1](REQUESTS.md#challengesftegroupcomplete-v1)
    - [Challenges/GetActiveChallenges v1](REQUESTS.md#challengesgetactivechallenges-v1)
    - [Challenges/PlayerProgress v1](REQUESTS.md#challengesplayerprogress-v1)
- Clubs
    - [Clubs/AcceptClubInvite v2](REQUESTS.md#clubsacceptclubinvite-v2)
    - [Clubs/CreateClub v1](REQUESTS.md#clubscreateclub-v1)
    - [Clubs/GetClubDetails v1](REQUESTS.md#clubsgetclubdetails-v1)
    - [Clubs/GetClubInvites v1](REQUESTS.md#clubsgetclubinvites-v1)
    - [Clubs/GetClubTitleInstances v1](REQUESTS.md#clubsgetclubtitleinstances-v1)
    - [Clubs/GetPlayerClubDetails v2](REQUESTS.md#clubsgetplayerclubdetails-v2)
    - [Clubs/GetStats v1](REQUESTS.md#clubsgetstats-v1)
    - [Clubs/InviteToClub v4](REQUESTS.md#clubsinvitetoclub-v4)
    - [Clubs/LeaveClub v1](REQUESTS.md#clubsleaveclub-v1)
    - [Clubs/RejectClubInvite v1](REQUESTS.md#clubsrejectclubinvite-v1)
    - [Clubs/UpdateClub v2](REQUESTS.md#clubsupdateclub-v2)
- Drop
    - [Drop/GetTradeInFilters v1](REQUESTS.md#dropgettradeinfilters-v1)
- GameServer
    - [GameServer/GetClubPrivateMatches v1](REQUESTS.md#gameservergetclubprivatematches-v1)
    - [GameServer/GetGameServerPingList v2](REQUESTS.md#gameservergetgameserverpinglist-v2)
- Matches
    - [Matches/GetMatchHistory v1](REQUESTS.md#matchesgetmatchhistory-v1)
- Matchmaking
    - [Matchmaking/PlayerCancelMatchmaking v1](REQUESTS.md#matchmakingplayercancelmatchmaking-v1)
    - [Matchmaking/PlayerSearchPrivateMatch v1](REQUESTS.md#matchmakingplayersearchprivatematch-v1)
    - [Matchmaking/StartMatchmaking v2](REQUESTS.md#matchmakingstartmatchmaking-v2)
- Microtransaction
    - [Microtransaction/ClaimEntitlements v2](REQUESTS.md#microtransactionclaimentitlements-v2)
    - [Microtransaction/GetCatalog v1](REQUESTS.md#microtransactiongetcatalog-v1)
    - [Microtransaction/StartPurchase v1](REQUESTS.md#microtransactionstartpurchase-v1)
- Party
    - [Party/ChangePartyOwner v1](REQUESTS.md#partychangepartyowner-v1)
    - [Party/CreateParty v1](REQUESTS.md#partycreateparty-v1)
    - [Party/GetPlayerPartyInfo v1](REQUESTS.md#partygetplayerpartyinfo-v1)
    - [Party/JoinParty v1](REQUESTS.md#partyjoinparty-v1)
    - [Party/KickPartyMembers v1](REQUESTS.md#partykickpartymembers-v1)
    - [Party/LeaveParty v1](REQUESTS.md#partyleaveparty-v1)
    - [Party/SendPartyChatMessage v1](REQUESTS.md#partysendpartychatmessage-v1)
    - [Party/SendPartyInvite v2](REQUESTS.md#partysendpartyinvite-v2)
    - [Party/SendPartyJoinRequest v1](REQUESTS.md#partysendpartyjoinrequest-v1)
    - [Party/SendPartyMessage v1](REQUESTS.md#partysendpartymessage-v1)
- Players
    - [Players/GetBanStatus v3](REQUESTS.md#playersgetbanstatus-v3)
    - [Players/GetCreatorCode v1](REQUESTS.md#playersgetcreatorcode-v1)
    - [Players/GetProfile v1](REQUESTS.md#playersgetprofile-v1)
    - [Players/GetXP v1](REQUESTS.md#playersgetxp-v1)
    - [Players/Report v4](REQUESTS.md#playersreport-v4)
- Playlists
    - [Playlists/GetActivePlaylists v1](REQUESTS.md#playlistsgetactiveplaylists-v1)
- Population
    - [Population/GetPopulation v1](REQUESTS.md#populationgetpopulation-v1)
    - [Population/UpdatePlayerPlaylist v1](REQUESTS.md#populationupdateplayerplaylist-v1)
- Products
    - [Products/CrossEntitlement/GetProductStatus v1](REQUESTS.md#productscrossentitlementgetproductstatus-v1)
    - [Products/GetContainerDropTable v2](REQUESTS.md#productsgetcontainerdroptable-v2)
    - [Products/GetPlayerProducts v2](REQUESTS.md#productsgetplayerproducts-v2)
    - [Products/TradeIn v2](REQUESTS.md#productstradein-v2)
    - [Products/UnlockContainer v2](REQUESTS.md#productsunlockcontainer-v2)
- Regions
    - [Regions/GetSubRegions v1](REQUESTS.md#regionsgetsubregions-v1)
- Reservations
    - [Reservations/JoinMatch v1](REQUESTS.md#reservationsjoinmatch-v1)
- RocketPass
    - [RocketPass/GetPlayerInfo v2](REQUESTS.md#rocketpassgetplayerinfo-v2)
    - [RocketPass/GetPlayerPrestigeRewards v1](REQUESTS.md#rocketpassgetplayerprestigerewards-v1)
    - [RocketPass/GetRewardContent v1](REQUESTS.md#rocketpassgetrewardcontent-v1)
- Shops
    - [Shops/GetPlayerWallet v1](REQUESTS.md#shopsgetplayerwallet-v1)
    - [Shops/GetShopCatalogue v2](REQUESTS.md#shopsgetshopcatalogue-v2)
    - [Shops/GetShopNotifications v1](REQUESTS.md#shopsgetshopnotifications-v1)
    - [Shops/GetStandardShops v1](REQUESTS.md#shopsgetstandardshops-v1)
- Skills
    - [Skills/GetPlayerSkill v1](REQUESTS.md#skillsgetplayerskill-v1)
    - [Skills/GetPlayersSkills v1](REQUESTS.md#skillsgetplayersskills-v1)
    - [Skills/GetSkillLeaderboard v1](REQUESTS.md#skillsgetskillleaderboard-v1)
    - [Skills/GetSkillLeaderboardRankForUsers v1](REQUESTS.md#skillsgetskillleaderboardrankforusers-v1)
    - [Skills/GetSkillLeaderboardValueForUser v1](REQUESTS.md#skillsgetskillleaderboardvalueforuser-v1)
- Stats
    - [Stats/GetStatLeaderboard v1](REQUESTS.md#statsgetstatleaderboard-v1)
    - [Stats/GetStatLeaderboardRankForUsers v1](REQUESTS.md#statsgetstatleaderboardrankforusers-v1)
    - [Stats/GetStatLeaderboardValueForUser v1](REQUESTS.md#statsgetstatleaderboardvalueforuser-v1)
- Tournaments
    - [Tournaments/Registration/RegisterTournament v1](REQUESTS.md#tournamentsregistrationregistertournament-v1)
    - [Tournaments/Registration/UnsubscribeTournament v1](REQUESTS.md#tournamentsregistrationunsubscribetournament-v1)
    - [Tournaments/Search/GetPublicTournaments v1](REQUESTS.md#tournamentssearchgetpublictournaments-v1)
    - [Tournaments/Search/GetSchedule v1](REQUESTS.md#tournamentssearchgetschedule-v1)
    - [Tournaments/Status/GetCycleData v1](REQUESTS.md#tournamentsstatusgetcycledata-v1)
    - [Tournaments/Status/GetScheduleRegion v1](REQUESTS.md#tournamentsstatusgetscheduleregion-v1)
    - [Tournaments/Status/GetTournamentSubscriptions v1](REQUESTS.md#tournamentsstatusgettournamentsubscriptions-v1)
- Training
    - [Training/BrowseTrainingData v1](REQUESTS.md#trainingbrowsetrainingdata-v1)
    - [Training/GetTrainingMetadata v1](REQUESTS.md#traininggettrainingmetadata-v1)
- Users
    - [Users/CanShowAvatar v1](REQUESTS.md#userscanshowavatar-v1)
- Misc
    - [DSR/RelayToServer v1](REQUESTS.md#dsrrelaytoserver-v1)
    - [Filters/FilterContent v1](REQUESTS.md#filtersfiltercontent-v1)
    - [Metrics/RecordMetrics v1](REQUESTS.md#metricsrecordmetrics-v1)
    - [Party/System](REQUESTS.md#partysystem)
    - [PsyPing](REQUESTS.md#psyping)
