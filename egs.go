package rlapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	egsUserAgent    = "UELauncher/11.0.1-14907503+++Portal+Release-Live Windows/10.0.19041.1.256.64bit"
	egsClientID     = "34a02cf8f4414e29b15921876da36f9a"
	egsClientSecret = "daafbccc737745039dffe53d94fc76cf"
	egsOAuthURL     = "account-public-service-prod03.ol.epicgames.com"
	eosAuthHeader   = "eHl6YTc4OTFwNUQ3czlSNkdtNm1vVEhXR2xvZXJwN0I6S25oMThkdTROVmxGcyszdVErWlBwRENWdG8wV1lmNHlYUDgrT2N3VnQxbw=="
	eosDeploymentID = "da32ae9c12ae40e8a112c52e1f17f3ba" // Rocket League
)

type TokenResponse struct {
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	ExpiresIn      int    `json:"expires_in"`
	ExpiresAt      string `json:"expires_at"`
	TokenType      string `json:"token_type"`
	ClientID       string `json:"client_id"`
	InternalClient bool   `json:"internal_client"`
	ClientService  string `json:"client_service"`
	AccountID      string `json:"account_id"`
	DisplayName    string `json:"displayName"`
	App            string `json:"app"`
	InAppID        string `json:"in_app_id"`
	DeviceID       string `json:"device_id"`
}

type EOSTokenResponse struct {
	AccessToken       string   `json:"access_token"`
	RefreshToken      string   `json:"refresh_token"`
	IDToken           string   `json:"id_token"`
	ExpiresIn         int      `json:"expires_in"`
	ExpiresAt         string   `json:"expires_at"`
	RefreshExpiresIn  int      `json:"refresh_expires_in"`
	RefreshExpiresAt  string   `json:"refresh_expires_at"`
	TokenType         string   `json:"token_type"`
	Scope             string   `json:"scope"`
	ClientID          string   `json:"client_id"`
	ApplicationID     string   `json:"application_id"`
	AccountID         string   `json:"account_id"`
	SelectedAccountID string   `json:"selected_account_id"`
	MergedAccounts    []string `json:"merged_accounts"`
	ACR               string   `json:"acr"`
	AuthTime          string   `json:"auth_time"`
}

// EGS provides an authentication layer for Epic Games Store -- largely adapted from https://github.com/derrod/legendary
type EGS struct {
	client *http.Client
}

// NewEGS creates a new Epic Games Store client
func NewEGS() *EGS {
	return &EGS{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// GetAuthURL returns the EGS login URL for manual browser authentication
func (e *EGS) GetAuthURL() string {
	loginURL := "https://www.epicgames.com/id/login?redirectUrl="
	redirectURL := fmt.Sprintf("https://www.epicgames.com/id/api/redirect?clientId=%s&responseType=code", egsClientID)
	return loginURL + url.QueryEscape(redirectURL)
}

// AuthenticateWithCode authenticates with EGS using an authorization code
func (e *EGS) AuthenticateWithCode(authCode string) (*TokenResponse, error) {
	return e.requestToken(map[string]string{
		"grant_type": "authorization_code",
		"code":       authCode,
		"token_type": "eg1",
	})
}

// AuthenticateWithRefreshToken authenticates with EGS using a refresh token
func (e *EGS) AuthenticateWithRefreshToken(refreshToken string) (*TokenResponse, error) {
	return e.requestToken(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"token_type":    "eg1",
	})
}

func (e *EGS) requestToken(params map[string]string) (*TokenResponse, error) {
	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/account/api/oauth/token", egsOAuthURL), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", egsUserAgent)

	auth := base64.StdEncoding.EncodeToString([]byte(egsClientID + ":" + egsClientSecret))
	req.Header.Set("Authorization", "Basic "+auth)

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			ErrorCode    string `json:"errorCode"`
			ErrorMessage string `json:"errorMessage"`
		}
		json.Unmarshal(body, &errorResp)
		return nil, fmt.Errorf("authentication failed: %s - %s", errorResp.ErrorCode, errorResp.ErrorMessage)
	}

	return &tokenResp, nil
}

// GetExchangeCode converts an EGS access token into an exchange code for EOS
func (e *EGS) GetExchangeCode(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/account/api/oauth/exchange", egsOAuthURL), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "bearer "+accessToken)
	req.Header.Set("User-Agent", egsUserAgent)

	resp, err := e.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	var tokenResp struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return tokenResp.Code, nil
}

// ExchangeEOSToken exchanges an exchange code for an EOS authentication token
func (e *EGS) ExchangeEOSToken(exchangeCode string) (*EOSTokenResponse, error) {
	return e.requestEOSToken(map[string]string{
		"grant_type":    "exchange_code",
		"exchange_code": exchangeCode,
	})
}

// RefreshEOSToken refreshes an EOS authentication token using a refresh token
func (e *EGS) RefreshEOSToken(refreshToken string) (*EOSTokenResponse, error) {
	return e.requestEOSToken(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	})
}

func (e *EGS) requestEOSToken(params map[string]string) (*EOSTokenResponse, error) {
	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}
	form.Set("deployment_id", eosDeploymentID)
	form.Set("scope", "basic_profile+friends_list+presence+friends_management+openid")

	req, err := http.NewRequest("POST", "https://api.epicgames.dev/epic/oauth/v2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+eosAuthHeader)
	req.Header.Set("User-Agent", egsUserAgent)

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)

	}

	var tokenResp EOSTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &tokenResp, nil
}
