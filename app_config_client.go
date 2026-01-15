package lembaas

import (
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
