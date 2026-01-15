package lembaas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AppConfigClient struct {
	baseURL   string
	authToken string

	httpClient *http.Client
}

func NewAppConfigClient(baseURL string, apiVersion int, token string) *AppConfigClient {
	return &AppConfigClient{
		baseURL:   fmt.Sprintf("%s/api/v%d", baseURL, apiVersion),
		authToken: token,

		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *AppConfigClient) ListCustomConfigValues(ctx context.Context) (AppConfigValueCollection, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/config/all", nil)
	if err != nil {
		return AppConfigValueCollection{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return AppConfigValueCollection{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AppConfigValueCollection{}, fmt.Errorf("expected status '%s', got '%s'", http.StatusText(http.StatusOK), resp.Status)
	}

	var configValues AppConfigValueCollection
	if err = json.NewDecoder(resp.Body).Decode(&configValues); err != nil {
		return AppConfigValueCollection{}, err
	}

	return configValues, nil
}

func (c *AppConfigClient) GetCustomConfigValue(ctx context.Context, configKey string) (AppConfigValue, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/config/%s/get", c.baseURL, configKey), nil)
	if err != nil {
		return AppConfigValue{}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return AppConfigValue{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return AppConfigValue{}, fmt.Errorf("expected status '%s', got '%s'", http.StatusText(http.StatusOK), resp.Status)
	}

	var configValue AppConfigValue
	if err = json.NewDecoder(resp.Body).Decode(&configValue); err != nil {
		return AppConfigValue{}, err
	}

	return configValue, nil
}

func (c *AppConfigClient) SetCustomConfigValue(ctx context.Context, config AppConfigValue) (AppConfigValue, error) {
	reqBody, err := json.Marshal(config)
	if err != nil {
		return AppConfigValue{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/config/set", bytes.NewBuffer(reqBody))
	if err != nil {
		return AppConfigValue{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return AppConfigValue{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AppConfigValue{}, fmt.Errorf("expected status '%s', got '%s'", http.StatusText(http.StatusOK), resp.Status)
	}

	var updatedConfig AppConfigValue
	if err = json.NewDecoder(resp.Body).Decode(&updatedConfig); err != nil {
		return AppConfigValue{}, err
	}

	return updatedConfig, nil
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
