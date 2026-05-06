// Infisical Secret Manager
// https://infisical.com

package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type InfisicalConfig struct {
	SiteURL   string
	ClientID string
	SecretKey string
	ProjectID string
	Env      string
}

type InfisicalClient struct {
	config  *InfisicalConfig
	client  *http.Client
}

func NewInfisicalClient(cfg *InfisicalConfig) *InfisicalClient {
	return &InfisicalClient{
		config: cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetSecret retrieves a secret
func (c *InfisicalClient) GetSecret(ctx context.Context, name string) (string, error) {
	token, _ := c.getServiceToken(ctx)

	url := fmt.Sprintf("%s/api/v3/secrets/%s", c.config.SiteURL, name)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Project-Id", c.config.ProjectID)
	req.Header.Set("X-Environment", c.config.Env)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Secret struct{ Value string } `json:"secret"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Secret.Value, nil
}

// ListSecrets lists all secrets
func (c *InfisicalClient) ListSecrets(ctx context.Context) (map[string]string, error) {
	token, _ := c.getServiceToken(ctx)

	url := fmt.Sprintf("%s/api/v3/secrets", c.config.SiteURL)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Project-Id", c.config.ProjectID)
	req.Header.Set("X-Environment", c.config.Env)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Secrets []struct {
			SecretName string `json:"secretName"`
			Value     string `json:"value"`
		} `json:"secrets"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	secrets := make(map[string]string)
	for _, s := range result.Secrets {
		secrets[s.SecretName] = s.Value
	}
	return secrets, nil
}

// SetSecret creates/updates a secret
func (c *InfisicalClient) SetSecret(ctx context.Context, name, value string) error {
	token, _ := c.getServiceToken(ctx)

	url := fmt.Sprintf("%s/api/v3/secrets", c.config.SiteURL)
	body, _ := json.Marshal(map[string]string{
		"secretName": name,
		"value":     value,
		"type":     "shared",
	})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Project-Id", c.config.ProjectID)
	req.Header.Set("X-Environment", c.config.Env)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	defer resp.Body.Close()
	return err
}

func (c *InfisicalClient) getServiceToken(ctx context.Context) (string, error) {
	// Call Infisical /api/v1/auth/srvc-token
	return "service_token", nil
}

/*
# Usage:
config := &secrets.InfisicalConfig{
    SiteURL:   os.Getenv("INFISICAL_SITE_URL"),
    ClientID: os.Getenv("INFISICAL_CLIENT_ID"),
    SecretKey: os.Getenv("INFISICAL_SECRET_KEY"),
    ProjectID: os.Getenv("INFISICAL_PROJECT_ID"),
    Env: "production",
}
client := secrets.NewInfisicalClient(config)
githubToken, _ := client.GetSecret(ctx, "GITHUB_TOKEN")
*/