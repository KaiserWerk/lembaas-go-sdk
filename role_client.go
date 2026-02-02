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
		Error string     `json:"error,omitempty"`
		Count int        `json:"count"`
		Roles []*AppRole `json:"roles"`
	}
	AppRoleResponse struct {
		Error string `json:"error,omitempty"`
		AppRole
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

func (c *RoleClient) ListRoles(ctx context.Context) (*AllAppRolesResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/roles", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var roles AllAppRolesResponse
	if err = json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("expected status '%s', got '%s' (%s)", http.StatusText(http.StatusOK), resp.Status, roles.Error)
	}

	return &roles, err
}

func (c *RoleClient) CreateRole(ctx context.Context, role *CreateAppRoleRequest) (*AppRoleResponse, error) {
	reqBody, err := json.Marshal(role)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/roles/create", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var roles AppRoleResponse
	if err = json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("expected status '%s', got '%s' (%s)", http.StatusText(http.StatusCreated), resp.Status, roles.Error)
	}

	return &roles, err
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
