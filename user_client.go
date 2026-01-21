package lembaas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type UserClient struct {
	baseURL   string
	authToken string

	httpClient *http.Client
}

func NewUserClient(baseURL, token string, apiVersion int) *UserClient {
	return &UserClient{
		baseURL:   fmt.Sprintf("%s/api/v%d", baseURL, apiVersion),
		authToken: token,

		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *UserClient) ListUsers(ctx context.Context) (AppUserCollection, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/users", nil)
	if err != nil {
		return AppUserCollection{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return AppUserCollection{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AppUserCollection{}, fmt.Errorf("expected status '%s', got '%s'", http.StatusText(http.StatusOK), resp.Status)
	}

	var users AppUserCollection
	if err = json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return AppUserCollection{}, err
	}

	return users, nil
}

func (c *UserClient) GetUser(ctx context.Context, userID int64) (AppUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/users/%d/get", c.baseURL, userID), nil)
	if err != nil {
		return AppUser{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return AppUser{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AppUser{}, fmt.Errorf("expected status '%s', got '%s'", http.StatusText(http.StatusOK), resp.Status)
	}

	var user AppUser
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return AppUser{}, err
	}

	return user, nil
}

func (c *UserClient) GetUserByEmail(ctx context.Context, email string) (AppUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/users/email/%s/get", c.baseURL, url.PathEscape(email)), nil)
	if err != nil {
		return AppUser{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return AppUser{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AppUser{}, fmt.Errorf("expected status '%s', got '%s'", http.StatusText(http.StatusOK), resp.Status)
	}

	var user AppUser
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return AppUser{}, err
	}

	return user, nil
}

func (c *UserClient) RegisterUser(ctx context.Context, user *AppUser) (AppUser, error) {
	reqBody, err := json.Marshal(user)
	if err != nil {
		return AppUser{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/users/register", bytes.NewBuffer(reqBody))
	if err != nil {
		return AppUser{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return AppUser{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return AppUser{}, fmt.Errorf("expected status '%s', got '%s'", http.StatusText(http.StatusCreated), resp.Status)
	}

	var updatedUser AppUser
	if err = json.NewDecoder(resp.Body).Decode(&updatedUser); err != nil {
		return AppUser{}, err
	}

	return updatedUser, nil
}

func (c *UserClient) UpdateUser(ctx context.Context, user *AppUser) (*AppUser, error) {
	reqBody, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/users/register", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status '%s', got '%s'", http.StatusText(http.StatusOK), resp.Status)
	}

	var createdUser AppUser
	if err = json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (c *UserClient) DeleteUser(ctx context.Context, userID int64) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/users/%d/delete", c.baseURL, userID), nil)
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

func (c *UserClient) LoginUser(ctx context.Context, auth AppUserAuthRequest) (*AppUserAuthResponse, error) {
	reqBody, err := json.Marshal(auth)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/users/login", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status '%s', got '%s'", http.StatusText(http.StatusOK), resp.Status)
	}

	var authResp AppUserAuthResponse
	if err = json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}

func (c *UserClient) LoginUserWithTOTP(ctx context.Context, auth TOTPRequest) (*AppUserAuthResponse, error) {
	reqBody, err := json.Marshal(auth)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/users/totp", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status '%s', got '%s'", http.StatusText(http.StatusOK), resp.Status)
	}

	var authResp AppUserAuthResponse
	if err = json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}
