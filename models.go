package lembaas

import "time"

type Error struct {
	Code    *int    `json:"code,omitzero"`
	Message *string `json:"message,omitzero"`
}

type AppTokenResponse struct {
	Error
	Token string `json:"token"`
	// ExpiresIn is in seconds
	ExpiresIn int    `json:"expires_in"`
	TokenType string `json:"token_type"`
}

type AppInfo struct {
	Error
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ClientID    string    `json:"client_id"`
	IconURL     string    `json:"icon_url"`
	CreatedAt   time.Time `json:"created_at"`
}

// AppUser represents a user within a specific application (stored in main database)
type AppUser struct {
	Error
	ID        int64     `json:"id"`
	AppID     int64     `json:"app_id"`
	Email     string    `json:"email"`
	RoleID    int64     `json:"role_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AppUserCollection struct {
	Error
	Count int        `json:"count"`
	Users []*AppUser `json:"users"`
}

// AppUserSession represents a user session for an app user (stored in main database)
type AppUserSession struct {
	Error
	ID        string    `json:"id"`
	AppID     int64     `json:"app_id"`
	UserID    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// AppRole represents a role within a specific application (stored in main database)
type AppRole struct {
	Error
	ID          int64     `json:"id"`
	AppID       int64     `json:"app_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Permissions string    `json:"permissions"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
}

type AppRoleCollection struct {
	Error
	Count int        `json:"count"`
	Roles []*AppRole `json:"roles"`
}

// AppCustomConfig represents a custom configuration key-value pair for a tenant app
type AppCustomConfig struct {
	ID          int64  `json:"id"`
	AppID       int64  `json:"app_id"`
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"`
	Enabled     bool   `json:"enabled"`
}

type AppCustomConfigCollection struct {
	Error
	Count        int                `json:"count"`
	ConfigValues []*AppCustomConfig `json:"config_values"`
}
