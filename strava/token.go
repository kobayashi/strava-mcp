package strava

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type TokenConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

func tokenFilePath() string {
	if p := os.Getenv("STRAVA_TOKENS_PATH"); p != "" {
		return p
	}
	exe, err := os.Executable()
	if err != nil {
		return "./tokens.json"
	}
	return filepath.Join(filepath.Dir(exe), "tokens.json")
}

func LoadToken() (*TokenConfig, error) {
	data, err := os.ReadFile(tokenFilePath())
	if err != nil {
		return nil, err
	}
	var cfg TokenConfig
	return &cfg, json.Unmarshal(data, &cfg)
}

func (cfg *TokenConfig) GetValidAccessToken() (string, error) {
	if time.Now().Unix() < cfg.ExpiresAt-60 {
		return cfg.AccessToken, nil
	}
	return cfg.refresh()
}

func (cfg *TokenConfig) refresh() (string, error) {
	resp, err := http.PostForm("https://www.strava.com/oauth/token", url.Values{
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
		"refresh_token": {cfg.RefreshToken},
		"grant_type":    {"refresh_token"},
	})
	if err != nil {
		return "", fmt.Errorf("token refresh failed: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresAt    int64  `json:"expires_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	cfg.AccessToken = result.AccessToken
	cfg.RefreshToken = result.RefreshToken
	cfg.ExpiresAt = result.ExpiresAt

	data, _ := json.MarshalIndent(cfg, "", "  ")
	if err := os.WriteFile(tokenFilePath(), data, 0600); err != nil {
		return "", fmt.Errorf("failed to persist refreshed token: %w", err)
	}

	return cfg.AccessToken, nil
}
