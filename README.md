# rlapi
Reverse engineered Rocket League internal HTTP & WebSocket API (with a partial Go wrapper).

**THIS PROJECT IS PROVIDED AS-IS** and is a compilation of research and intercepted requests for Rocket League's internal APIs. While not all endpoints are fully documented or implemented, this repository contains all relevant resources and serves as a foundation for further exploration. Do not ask me about specific endpoints, I probably don’t know.

## Contributions
All contributions are welcome! If you discover new endpoints, extend the Go wrapper, or add additional functionality, please submit a PR.

## Getting Started
See the godoc for detailed documentation on the Go wrapper. The rest of this README contains resources on reverse engineering and API endpoints.

Comprehensive examples are available in the `examples` directory.

### Usage
```bash
go get github.com/dank/rlapi
```

### Authentication
Rocket League authentication always goes through Epic Online Services (EOS), either via the Epic Games Store (EGS) or by exchanging a Steam session ticket for an EOS token.

This library provides full end-to-end authentication via EGS. Steam login and ticket generation are out of scope, but a method is provided to exchange a Steam session ticket for an EOS token, and users can leverage external libraries such as [steam-user](https://github.com/DoctorMcKay/node-steam-user) or [SteamKit](https://partner.steamgames.com/doc/api/ISteamUser#GetAuthSessionTicket) to obtain the ticket.

## Intercepting Requests

Traditional proxy tools like Fiddler don’t work with Rocket League due to certificate pinning.

To intercept traffic, we use Frida dynamic instrumentation to hook curl functions at runtime, disabling SSL verification and redirecting API calls from `api.rlpp.psynet.gg` to a local MITM server.

### MITM Server

Even with certificate pinning disabled, Rocket League still requires HTTPS and WSS connections. To handle this, the MITM proxy uses self-signed certificates, acting as a forwarder proxy for both HTTP and WebSocket traffic.
The authentication endpoint response is intercepted and rewritten so that the WebSocket URL returned by the server points to the local MITM WebSocket server.

Optionally, the MITM server can be configured to route through a Fiddler proxy (`http://127.0.0.1:8888`) to intercept the traffic in a more familiar interface.

See the respective READMEs in `tools/frida/` and `tools/mitm/` directories for usage instructions and setup details.

## Reconstructing Requests

### Authentication

Rocket League initially establishes a connection with the HTTP API before transitioning to a WebSocket connection. The client sends an EOS access token via HTTP and receives session credentials, the WebSocket endpoint URL, and any tokens required for further communication. The client then connects to the WebSocket using these tokens, allowing all subsequent API calls to occur over a persistent WebSocket connection.

### Signing
All API requests and responses must be signed using `PsySig` headers with HMAC-SHA256. The signing keys were reverse engineered from the game binary and are XOR-encrypted with a 4-byte pattern. For example, in Python:
```python
# Raw data from IDA dump
data = [0x36, 0xEA, 0x37, 0x0C, ...]  # 36 bytes total

key_bytes = [data[i] ^ data[(i % 4) + 32] for i in range(32)]
```

- **Requests**: `c338bd36fb8c42b1a431d30add939fc7`
  - Format: `HMAC-SHA256(key, "-" + request_body)`
- **Responses**: `3b932153785842ac927744b292e40e52`  
  - Format: `HMAC-SHA256(key, PsyTime + "-" + response_body)`

All signatures are base64-encoded.

### Schema

Rocket League WebSocket messages use a custom HTTP-like schema with headers and JSON body:

```
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

**Required Headers:**
- `PsyService`: API endpoint (e.g., `Matchmaking/StartMatchmaking v2`)
- `PsyRequestID`: Idempotency key for request/response matching (`PsyNetMessage_X_`)
- `PsyToken`: Authentication token
- `PsySessionID`: Session identifier
- `PsySig`: HMAC signature of the body
- `PsyBuildID`: Varies by build
- `User-Agent`: Varies by build
- `PsyEnvironment`: Environment (`Prod` for production)

The message format is: headers (each ending with `\r\n`) followed by `\r\n\r\n` separator, then JSON body.

## HTTP Endpoints
### Auth/AuthPlayer/v2
## WebSocket Endpoints
### Challenges/FTECheckpointComplete v1
### Challenges/FTEGroupComplete v1
### Challenges/GetActiveChallenges v1
### Challenges/PlayerProgress v1
### Clubs/GetClubInvites v1
### Clubs/GetClubTitleInstances v1
### Clubs/GetPlayerClubDetails v2
### Clubs/GetStats v1
### Clubs/RejectClubInvite v1
### Drop/GetTradeInFilters v1
### Filters/FilterContent v1
### GameServer/GetGameServerPingList v2
### Matches/GetMatchHistory v1
### Matchmaking/PlayerCancelMatchmaking v1
### Matchmaking/StartMatchmaking v2
### Metrics/RecordMetrics v1
### Microtransaction/ClaimEntitlements v2
### Microtransaction/GetCatalog v1
### Microtransaction/StartPurchase v1
### Party/CreateParty v1
### Party/GetPlayerPartyInfo v1
### Party/LeaveParty v1
### Party/SendPartyMessage v1
### Party/System
### Players/GetBanStatus v3
### Players/GetCreatorCode v1
### Players/GetProfile v1
### Players/GetXP v1
### Playlists/GetActivePlaylists v1
### Population/GetPopulation v1
### Population/UpdatePlayerPlaylist v1
### Products/CrossEntitlement/GetProductStatus v1
### Products/GetContainerDropTable v2
### Products/GetPlayerProducts v2
### Products/TradeIn v2
### Products/UnlockContainer v2
### Regions/GetSubRegions v1
### RocketPass/GetPlayerInfo v2
### RocketPass/GetPlayerPrestigeRewards v1
### RocketPass/GetRewardContent v1
### Shops/GetPlayerWallet v1
### Shops/GetShopCatalogue v2
### Shops/GetShopNotifications v1
### Shops/GetStandardShops v1
### Skills/GetPlayerSkill v1
### Skills/GetSkillLeaderboard v1
### Skills/GetSkillLeaderboardRankForUsers v1
### Skills/GetSkillLeaderboardValueForUser v1
### Stats/GetStatLeaderboard v1
### Stats/GetStatLeaderboardValueForUser v1
### Tournaments/Search/GetSchedule v1
### Tournaments/Status/GetCycleData v1
### Tournaments/Status/GetScheduleRegion v1
### Tournaments/Status/GetTournamentSubscriptions v1
### Users/CanShowAvatar v1
