package lembaas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AppClient struct {
	baseURL   string
	authToken string

	httpClient *http.Client
}

func NewAppClient(baseURL string, apiVersion int) *AppClient {
	return &AppClient{
		baseURL:    fmt.Sprintf("%s/api/v%d", baseURL, apiVersion),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *AppClient) GetAuthToken(ctx context.Context, clientID, clientSecret string) (AppTokenResponse, error) {
	reqBody := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
	}
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return AppTokenResponse{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/token", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return AppTokenResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return AppTokenResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AppTokenResponse{}, fmt.Errorf("expected status '%s', got '%s'", http.StatusText(http.StatusOK), resp.Status)
	}

	var tokenResp AppTokenResponse
	if err = json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return AppTokenResponse{}, err
	}

	return tokenResp, nil
}

func (c *AppClient) GetAppInfo(ctx context.Context, token string) (AppInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/app", nil)
	if err != nil {
		return AppInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return AppInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AppInfo{}, fmt.Errorf("expected status '%s', got '%s'", http.StatusText(http.StatusOK), resp.Status)
	}

	var appInfo AppInfo
	if err = json.NewDecoder(resp.Body).Decode(&appInfo); err != nil {
		return AppInfo{}, err
	}

	return appInfo, nil
}
