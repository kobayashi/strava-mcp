# strava-mcp

A MCP server that exposes your Strava data to AI assistants.

## Tools

| Tool | Description |
|------|-------------|
| `get_today_activities` | Fetch activities recorded today |
| `get_recent_activities` | Fetch activities from the last N days (default: 7) |
| `get_athlete_stats` | Fetch weekly, monthly, and yearly totals |

## Setup

### 1. Create a Strava API application

Go to https://www.strava.com/settings/api and create an app. Note the **Client ID** and **Client Secret**.

### 2. Obtain tokens

Use the Strava OAuth flow to get an access token and refresh token, then create `tokens.json` in the project root:

```json
{
  "client_id": "YOUR_CLIENT_ID",
  "client_secret": "YOUR_CLIENT_SECRET",
  "access_token": "YOUR_ACCESS_TOKEN",
  "refresh_token": "YOUR_REFRESH_TOKEN",
  "expires_at": 0
}
```

The server automatically refreshes the access token when it expires.

### 3. Build

```bash
go build -o strava-mcp .
```

### 4. Configure your MCP client

Add the server to your MCP client config (e.g. Claude Desktop `claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "strava": {
      "command": "/path/to/strava-mcp"
    }
  }
}
```

## Requirements

- Go 1.24+
- A Strava account with API access
