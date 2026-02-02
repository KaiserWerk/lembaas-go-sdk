package lembaas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type (
	AppConfigClient struct {
		baseURL    string
		authToken  string
		httpClient *http.Client
	}

	AppConfigValueResponse struct {
		Error       string `json:"error,omitempty"`
		ConfigKey   string `json:"config_key"`
		ConfigValue string `json:"config_value"`
		Enabled     bool   `json:"enabled"`
	}
	AllAppConfigValuesResponse struct {
		Error        string            `json:"error,omitempty"`
		Count        int               `json:"count"`
		ConfigValues []*AppConfigValue `json:"config_values"`
	}
)

func NewConfigClient(baseURL string, apiVersion int, token string) *AppConfigClient {
	return &AppConfigClient{
		baseURL:   fmt.Sprintf("%s/api/v%d", baseURL, apiVersion),
		authToken: token,

		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *AppConfigClient) ListCustomConfigValues(ctx context.Context) (*AllAppConfigValuesResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/config/all", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var configValues AllAppConfigValuesResponse
	if err = json.NewDecoder(resp.Body).Decode(&configValues); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("expected status '%s', got '%s' (%s)", http.StatusText(http.StatusOK), resp.Status, configValues.Error)
	}

	return &configValues, err
}

func (c *AppConfigClient) GetCustomConfigValue(ctx context.Context, configKey string) (*AppConfigValueResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/config/%s/get", c.baseURL, configKey), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var configResponse AppConfigValueResponse
	if err = json.NewDecoder(resp.Body).Decode(&configResponse); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("expected status '%s', got '%s' (%s)", http.StatusText(http.StatusOK), resp.Status, configResponse.Error)
	}

	return &configResponse, err
}

func (c *AppConfigClient) SetCustomConfigValue(ctx context.Context, key, value string) (*AppConfigValueResponse, error) {
	config := AppConfigValue{
		ConfigKey:   key,
		ConfigValue: value,
	}
	reqBody, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/config/set", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var configResponse AppConfigValueResponse
	if err = json.NewDecoder(resp.Body).Decode(&configResponse); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("expected status '%s', got '%s' (%s)", http.StatusText(http.StatusOK), resp.Status, configResponse.Error)
	}

	return &configResponse, err
}

func (c *AppConfigClient) DeleteCustomConfigValue(ctx context.Context, configKey string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/config/%s/delete", c.baseURL, configKey), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("expected status '%s', got '%s'", http.StatusText(http.StatusNoContent), resp.Status)
	}

	return nil
}
