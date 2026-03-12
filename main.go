package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "strconv"

    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
    "strava-mcp/strava"
)

func main() {
    s := server.NewMCPServer(
        "Strava Coach MCP",
        "1.0.0",
        server.WithToolCapabilities(false),
    )

    s.AddTool(
        mcp.NewTool("get_today_activities",
            mcp.WithDescription("Get today's activities"),
        ),
        func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
            client, err := strava.NewClient()
            if err != nil {
                return mcp.NewToolResultError(err.Error()), nil
            }
            activities, err := client.GetTodayActivities()
            if err != nil {
                return mcp.NewToolResultError(err.Error()), nil
            }
            if len(activities) == 0 {
                return mcp.NewToolResultText("No activity today"), nil
            }
            data, err := json.MarshalIndent(activities, "", "  ")
            if err != nil {
                return mcp.NewToolResultError(err.Error()), nil
            }
            return mcp.NewToolResultText(string(data)), nil
        },
    )

    s.AddTool(
        mcp.NewTool("get_recent_activities",
            mcp.WithDescription("Get recent N days activities"),
            mcp.WithString("days",
							mcp.Description("days (default 7)"),
            ),
        ),
        func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
            days := 7
            if d := req.GetString("days", ""); d != "" {
                if n, err := strconv.Atoi(d); err == nil {
                    days = n
                }
            }
            client, err := strava.NewClient()
            if err != nil {
                return mcp.NewToolResultError(err.Error()), nil
            }
            activities, err := client.GetRecentActivities(days)
            if err != nil {
                return mcp.NewToolResultError(err.Error()), nil
            }
            data, err := json.MarshalIndent(activities, "", "  ")
            if err != nil {
                return mcp.NewToolResultError(err.Error()), nil
            }
            return mcp.NewToolResultText(fmt.Sprintf("%d days (%d activities):\n%s", days, len(activities), data)), nil
        },
    )

    s.AddTool(
        mcp.NewTool("get_athlete_stats",
            mcp.WithDescription("Get weekly, monthly, and yearly statistics"),
        ),
        func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
            client, err := strava.NewClient()
            if err != nil {
                return mcp.NewToolResultError(err.Error()), nil
            }
            stats, err := client.GetAthleteStats()
            if err != nil {
                return mcp.NewToolResultError(err.Error()), nil
            }
            data, err := json.MarshalIndent(stats, "", "  ")
            if err != nil {
                return mcp.NewToolResultError(err.Error()), nil
            }
            return mcp.NewToolResultText(string(data)), nil
        },
    )

    if err := server.NewStdioServer(s).Listen(context.Background(), nil, nil); err != nil {
        log.Fatal(err)
    }
}
