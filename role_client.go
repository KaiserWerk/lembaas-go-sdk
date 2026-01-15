package lembaas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RoleClient struct {
	baseURL   string
	authToken string

	httpClient *http.Client
}

func NewRoleClient(baseURL string, apiVersion int, token string) *RoleClient {
	return &RoleClient{
		baseURL:   fmt.Sprintf("%s/api/v%d", baseURL, apiVersion),
		authToken: token,

		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *RoleClient) ListRoles(ctx context.Context) (AppRoleCollection, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/roles", nil)
	if err != nil {
		return AppRoleCollection{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return AppRoleCollection{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AppRoleCollection{}, fmt.Errorf("expected status '%s', got '%s'", http.StatusText(http.StatusOK), resp.Status)
	}

	var roles AppRoleCollection
	if err = json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return AppRoleCollection{}, err
	}

	return roles, nil
}

func (c *RoleClient) CreateRole(ctx context.Context, role AppRole) (AppRole, error) {
	reqBody, err := json.Marshal(role)
	if err != nil {
		return AppRole{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/roles/create", bytes.NewBuffer(reqBody))
	if err != nil {
		return AppRole{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return AppRole{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return AppRole{}, fmt.Errorf("expected status '%s', got '%s'", http.StatusText(http.StatusCreated), resp.Status)
	}

	var roles AppRole
	if err = json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return AppRole{}, err
	}

	return roles, nil
}

func (c *RoleClient) DeleteRole(ctx context.Context, roleID int64) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/roles/%d/delete", c.baseURL, roleID), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status '%s', got '%s'", http.StatusText(http.StatusOK), resp.Status)
	}

	return nil
}
