package lembaas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Responses
type (
	AppTokenResponse struct {
		Error string `json:"error,omitempty"`
		Token string `json:"token"`
		// ExpiresIn is in seconds
		ExpiresIn int    `json:"expires_in"`
		TokenType string `json:"token_type"`
	}
	AppInfoResponse struct {
		Error       string    `json:"error,omitempty"`
		ID          int64     `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		ClientID    string    `json:"client_id"`
		IconURL     string    `json:"icon_url"`
		CreatedAt   time.Time `json:"created_at"`
	}
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

func (c *AppClient) GetAuthToken(ctx context.Context, clientID, clientSecret string) (*AppTokenResponse, error) {
	reqBody := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
	}
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/token", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResponse AppTokenResponse
	if err = json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("expected status '%s', got '%s' (%s)", http.StatusText(http.StatusOK), resp.Status, tokenResponse.Error)
	}

	return &tokenResponse, err
}

func (c *AppClient) GetAppInfo(ctx context.Context, token string) (*AppInfoResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/app", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var appInfoResponse AppInfoResponse
	if err = json.NewDecoder(resp.Body).Decode(&appInfoResponse); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("expected status '%s', got '%s' (%s)", http.StatusText(http.StatusOK), resp.Status, appInfoResponse.Error)
	}

	return &appInfoResponse, err
}
