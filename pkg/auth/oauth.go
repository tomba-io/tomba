package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/tomba-io/tomba/pkg/config"
)

const (
	// DefaultAPIBaseURL is the default base URL for the Tomba API
	DefaultAPIBaseURL = "https://api.tomba.io/v1"
	// ClientID for the Tomba CLI (public client)
	ClientID = "tomba-cli"
	// DefaultScopes requested by the CLI
	DefaultScopes = "read search verify enrich leads account"
)

// GetAPIBaseURL returns the API base URL, allowing override via TOMBA_API_URL env var.
func GetAPIBaseURL() string {
	if u := os.Getenv("TOMBA_API_URL"); u != "" {
		return u
	}
	return DefaultAPIBaseURL
}

// DeviceCodeResponse is the response from POST /oauth/device
type DeviceCodeResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

// TokenResponse is the response from POST /oauth/token
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// TokenErrorResponse is the error response from POST /oauth/token
type TokenErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// RequestDeviceCode initiates the device authorization flow.
func RequestDeviceCode() (*DeviceCodeResponse, error) {
	data := url.Values{}
	data.Set("client_id", ClientID)
	data.Set("scope", DefaultScopes)

	resp, err := http.PostForm(GetAPIBaseURL()+"/oauth/device", data)
	if err != nil {
		return nil, fmt.Errorf("failed to request device code: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("device code request failed (%d): %s", resp.StatusCode, string(body))
	}

	var dcResp DeviceCodeResponse
	if err := json.Unmarshal(body, &dcResp); err != nil {
		return nil, fmt.Errorf("failed to parse device code response: %w", err)
	}

	return &dcResp, nil
}

// PollForToken polls the token endpoint until the user authorizes or the code expires.
func PollForToken(deviceCode string, interval, expiresIn int) (*TokenResponse, error) {
	deadline := time.Now().Add(time.Duration(expiresIn) * time.Second)
	pollInterval := time.Duration(interval) * time.Second

	for time.Now().Before(deadline) {
		time.Sleep(pollInterval)

		data := url.Values{}
		data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
		data.Set("device_code", deviceCode)
		data.Set("client_id", ClientID)

		resp, err := http.PostForm(GetAPIBaseURL()+"/oauth/token", data)
		if err != nil {
			continue // Network error, retry
		}

		body, _ := io.ReadAll(resp.Body)
		func() { _ = resp.Body.Close() }()

		if resp.StatusCode == 200 {
			var tokenResp TokenResponse
			if err := json.Unmarshal(body, &tokenResp); err != nil {
				return nil, fmt.Errorf("failed to parse token response: %w", err)
			}
			return &tokenResp, nil
		}

		var errResp TokenErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			continue
		}

		switch errResp.Error {
		case "authorization_pending":
			continue
		case "slow_down":
			pollInterval += 5 * time.Second
			continue
		case "access_denied":
			return nil, fmt.Errorf("access denied by user")
		case "expired_token":
			return nil, fmt.Errorf("device code expired")
		default:
			return nil, fmt.Errorf("token error: %s - %s", errResp.Error, errResp.ErrorDescription)
		}
	}

	return nil, fmt.Errorf("device code expired (timeout)")
}

// ErrSessionExpired indicates the refresh token was revoked or expired.
var ErrSessionExpired = fmt.Errorf("session expired")

// RefreshAccessToken uses a refresh token to get a new access token.
func RefreshAccessToken(refreshToken string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", ClientID)

	resp, err := http.PostForm(GetAPIBaseURL()+"/oauth/token", data)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		var errResp TokenErrorResponse
		if json.Unmarshal(body, &errResp) == nil && errResp.Error == "invalid_grant" {
			return nil, ErrSessionExpired
		}
		return nil, fmt.Errorf("refresh failed (%d): %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse refresh response: %w", err)
	}

	return &tokenResp, nil
}

// SaveTokens saves OAuth tokens to the config file.
func SaveTokens(tokens *TokenResponse) error {
	expiry := time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)
	return config.UpdateConfig(config.Config{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenExpiry:  expiry.Format(time.RFC3339),
		AuthMethod:   "oauth",
	})
}

// IsTokenExpired checks if the stored access token has expired.
func IsTokenExpired(expiryStr string) bool {
	if expiryStr == "" {
		return true
	}
	expiry, err := time.Parse(time.RFC3339, expiryStr)
	if err != nil {
		return true
	}
	return time.Now().After(expiry.Add(-30 * time.Second)) // 30s buffer
}

// EnsureValidToken checks if the access token is valid, refreshing if needed.
// Returns the valid access token or an error.
func EnsureValidToken(conf *config.Config) (string, error) {
	if conf.AuthMethod != "oauth" || conf.AccessToken == "" {
		return "", fmt.Errorf("not using OAuth authentication")
	}

	if !IsTokenExpired(conf.TokenExpiry) {
		return conf.AccessToken, nil
	}

	// Token expired, try refresh
	if conf.RefreshToken == "" {
		return "", fmt.Errorf("access token expired and no refresh token available")
	}

	tokens, err := RefreshAccessToken(conf.RefreshToken)
	if err != nil {
		if err == ErrSessionExpired {
			// Clear saved tokens so next command prompts re-login
			_ = config.UpdateConfig(config.Config{AuthMethod: ""})
			return "", fmt.Errorf("your session has expired or been revoked. Please run 'tomba login' to re-authenticate")
		}
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}

	// Save new tokens
	if err := SaveTokens(tokens); err != nil {
		return "", fmt.Errorf("failed to save refreshed tokens: %w", err)
	}

	return tokens.AccessToken, nil
}

// FetchAPICredentials uses an OAuth access token to retrieve the user's API key and secret,
// allowing the existing Go SDK to work with API key auth.
func FetchAPICredentials(accessToken string) (apiKey, apiSecret string, err error) {
	// Fetch secret_token from /me
	meReq, err := http.NewRequest("GET", GetAPIBaseURL()+"/me", nil)
	if err != nil {
		return "", "", err
	}
	meReq.Header.Set("Authorization", "Bearer "+accessToken)

	meResp, err := http.DefaultClient.Do(meReq)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch account: %w", err)
	}
	defer func() { _ = meResp.Body.Close() }()

	meBody, _ := io.ReadAll(meResp.Body)
	if meResp.StatusCode != 200 {
		return "", "", fmt.Errorf("account request failed (%d): %s", meResp.StatusCode, string(meBody))
	}

	var meResult struct {
		Data struct {
			SecretToken string `json:"secret_token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(meBody, &meResult); err != nil {
		return "", "", fmt.Errorf("failed to parse account response: %w", err)
	}

	// Fetch API key from /keys
	keysReq, err := http.NewRequest("GET", GetAPIBaseURL()+"/keys", nil)
	if err != nil {
		return "", "", err
	}
	keysReq.Header.Set("Authorization", "Bearer "+accessToken)

	keysResp, err := http.DefaultClient.Do(keysReq)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch keys: %w", err)
	}
	defer func() { _ = keysResp.Body.Close() }()

	keysBody, _ := io.ReadAll(keysResp.Body)
	if keysResp.StatusCode != 200 {
		return "", "", fmt.Errorf("keys request failed (%d): %s", keysResp.StatusCode, string(keysBody))
	}

	var keysResult struct {
		Data struct {
			Keys []struct {
				Key string `json:"key"`
			} `json:"keys"`
		} `json:"data"`
	}
	if err := json.Unmarshal(keysBody, &keysResult); err != nil {
		return "", "", fmt.Errorf("failed to parse keys response: %w", err)
	}

	if len(keysResult.Data.Keys) == 0 {
		return "", "", fmt.Errorf("no API keys found")
	}

	return keysResult.Data.Keys[0].Key, meResult.Data.SecretToken, nil
}

// OpenBrowser opens the given URL in the default browser.
func OpenBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", strings.ReplaceAll(url, "&", "^&"))
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}
