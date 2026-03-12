package strava

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Activity struct {
    ID            int64   `json:"id"`
    Name          string  `json:"name"`
    Type          string  `json:"type"`
    StartDate     string  `json:"start_date_local"`
    Distance      float64 `json:"distance"`
    MovingTime    int     `json:"moving_time"`
    TotalElevGain float64 `json:"total_elevation_gain"`
    AverageSpeed  float64 `json:"average_speed"`
    MaxSpeed      float64 `json:"max_speed"`
    AverageHR     float64 `json:"average_heartrate"`
    MaxHR         float64 `json:"max_heartrate"`
    Kudos         int     `json:"kudos_count"`
}

type Client struct {
    token *TokenConfig
}

func NewClient() (*Client, error) {
    cfg, err := LoadToken()
    if err != nil {
        return nil, fmt.Errorf("Failed to load tokens.json: %w", err)
    }
    return &Client{token: cfg}, nil
}

func (c *Client) get(path string, params map[string]string) (*http.Response, error) {
    token, err := c.token.GetValidAccessToken()
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest("GET", "https://www.strava.com/api/v3"+path, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    req.Header.Set("Authorization", "Bearer "+token)

    q := req.URL.Query()
    for k, v := range params {
        q.Set(k, v)
    }
    req.URL.RawQuery = q.Encode()

    return http.DefaultClient.Do(req)
}

func (c *Client) GetTodayActivities() ([]Activity, error) {
    today := time.Now().Truncate(24 * time.Hour)
    return c.getActivitiesSince(today)
}

func (c *Client) GetRecentActivities(days int) ([]Activity, error) {
    since := time.Now().AddDate(0, 0, -days)
    return c.getActivitiesSince(since)
}

func (c *Client) getActivitiesSince(since time.Time) ([]Activity, error) {
    resp, err := c.get("/athlete/activities", map[string]string{
        "after":    fmt.Sprintf("%d", since.Unix()),
        "per_page": "30",
    })
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("strava API error: %s", resp.Status)
    }

    var activities []Activity
    return activities, json.NewDecoder(resp.Body).Decode(&activities)
}

func (c *Client) GetAthleteStats() (map[string]any, error) {
    resp, err := c.get("/athlete", nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("strava API error: %s", resp.Status)
    }

    var athlete struct {
        ID int64 `json:"id"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&athlete); err != nil {
        return nil, fmt.Errorf("failed to decode athlete: %w", err)
    }

    resp2, err := c.get(fmt.Sprintf("/athletes/%d/stats", athlete.ID), nil)
    if err != nil {
        return nil, err
    }
    defer resp2.Body.Close()

    if resp2.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("strava API error: %s", resp2.Status)
    }

    var stats map[string]any
    return stats, json.NewDecoder(resp2.Body).Decode(&stats)
}
