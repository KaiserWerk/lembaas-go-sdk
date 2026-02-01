package lembaas

import "time"

// AppUser represents a user within a specific application
type AppUser struct {
	ID        int64     `json:"id"`
	AppID     int64     `json:"app_id"`
	Email     string    `json:"email"`
	RoleID    int64     `json:"role_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AppUserSession represents a user session for an app user
type AppUserSession struct {
	Message   string    `json:"message,omitempty"`
	ID        string    `json:"id"`
	AppID     int64     `json:"app_id"`
	UserID    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// AppRole represents a role within a specific application
type AppRole struct {
	Message     string    `json:"message,omitempty"`
	ID          int64     `json:"id"`
	AppID       int64     `json:"app_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Permissions string    `json:"permissions"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
}

type AppConfigValue struct {
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"`
	Enabled     bool   `json:"enabled"`
}
