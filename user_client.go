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

// Requests
type (
	CreateAppUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		RoleID   int64  `json:"role_id"`
		IsActive bool   `json:"is_active"`
	}
	UpdateAppUserRequest struct {
		ID       int64  `json:"id"`
		Email    string `json:"email"`
		Password string `json:"password"`
		RoleID   int64  `json:"role_id"`
		IsActive bool   `json:"is_active"`
	}
	AppUserAuthRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	TOTPLoginRequest struct {
		LoginCode string `json:"login_code"`
		TOTPCode  string `json:"totp_code"`
	}
)

// Responses
type (
	TOTPEnableConfirmRequest struct {
		TOTPCode string `json:"totp_code"`
	}
	CreateAppUserResponse struct {
		Error string  `json:"error,omitempty"`
		User  AppUser `json:"user"`
	}

	AppUserCollectionResponse struct {
		Error string     `json:"error,omitempty"`
		Count int        `json:"count"`
		Users []*AppUser `json:"users"`
	}
	AppUserAuthResponse struct {
		Error        string    `json:"error,omitempty"`
		SessionToken string    `json:"session_token"`
		UserID       int64     `json:"user_id"`
		Email        string    `json:"email"`
		RoleID       int64     `json:"role_id"`
		ExpiresAt    time.Time `json:"expires_at"`
		ExpiresIn    int       `json:"expires_in"`

		/* For 2FA */
		LoginCode           string    `json:"login_code"`
		LoginCodeValidUntil time.Time `json:"login_code_valid_until"`
	}
	TOTPEnableResponse struct {
		Error  string `json:"error,omitempty"`
		QRCode []byte `json:"qr_code"`
	}
)

type UserClient struct {
	baseURL    string
	authToken  string
	httpClient *http.Client
}

func NewUserClient(baseURL, token string, apiVersion int) *UserClient {
	return &UserClient{
		baseURL:   fmt.Sprintf("%s/api/v%d", baseURL, apiVersion),
		authToken: token,

		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *UserClient) ListUsers(ctx context.Context) (*AppUserCollectionResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/users", nil)
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

	var users AppUserCollectionResponse
	if err = json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return &users, nil
}

func (c *UserClient) GetUser(ctx context.Context, userID int64) (*AppUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/users/%d/get", c.baseURL, userID), nil)
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

	var user AppUser
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *UserClient) GetUserByEmail(ctx context.Context, email string) (*AppUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/users/email/%s/get", c.baseURL, url.PathEscape(email)), nil)
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

	var user AppUser
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *UserClient) RegisterUser(ctx context.Context, request *CreateAppUserRequest) (*CreateAppUserResponse, error) {
	reqBody, err := json.Marshal(request)
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

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("expected status '%s', got '%s'", http.StatusText(http.StatusCreated), resp.Status)
	}

	var updatedUser CreateAppUserResponse
	if err = json.NewDecoder(resp.Body).Decode(&updatedUser); err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (c *UserClient) UpdateUser(ctx context.Context, request *UpdateAppUserRequest) (*AppUser, error) {
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/users/update", bytes.NewBuffer(reqBody))
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

func (c *UserClient) EnableTOTPForUser(ctx context.Context, userID int64) (*TOTPEnableResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/users/%d/totp/enable", c.baseURL, userID), nil)
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

	var totpResp TOTPEnableResponse
	if err = json.NewDecoder(resp.Body).Decode(&totpResp); err != nil {
		return nil, err
	}

	return &totpResp, nil
}

func (c *UserClient) ConfirmEnableTOTPForUser(ctx context.Context, userID int64, code string) (*TOTPEnableResponse, error) {
	r := TOTPEnableConfirmRequest{TOTPCode: code}
	reqBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/users/%d/totp/enable/confirm", c.baseURL, userID), bytes.NewBuffer(reqBody))
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

	var totpResp TOTPEnableResponse
	if err = json.NewDecoder(resp.Body).Decode(&totpResp); err != nil {
		return nil, err
	}

	return &totpResp, nil
}

func (c *UserClient) LoginUser(ctx context.Context, request *AppUserAuthRequest) (*AppUserAuthResponse, error) {
	reqBody, err := json.Marshal(request)
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

func (c *UserClient) LoginUserWithTOTP(ctx context.Context, request *TOTPLoginRequest) (*AppUserAuthResponse, error) {
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/users/login/totp", bytes.NewBuffer(reqBody))
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
