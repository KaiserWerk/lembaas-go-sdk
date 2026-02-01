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
	CreateAppRoleRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Permissions string `json:"permissions"`
		IsDefault   bool   `json:"is_default"`
	}
	AllAppRolesResponse struct {
		Message string     `json:"message,omitempty"`
		Count   int        `json:"count"`
		Roles   []*AppRole `json:"roles"`
	}
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

func (c *RoleClient) ListRoles(ctx context.Context) (AllAppRolesResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/roles", nil)
	if err != nil {
		return AllAppRolesResponse{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return AllAppRolesResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AllAppRolesResponse{}, fmt.Errorf("expected status '%s', got '%s'", http.StatusText(http.StatusOK), resp.Status)
	}

	var roles AllAppRolesResponse
	if err = json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return AllAppRolesResponse{}, err
	}

	return roles, nil
}

func (c *RoleClient) CreateRole(ctx context.Context, role CreateAppRoleRequest) (AppRole, error) {
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

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("expected status '%s', got '%s'", http.StatusText(http.StatusNoContent), resp.Status)
	}

	return nil
}
