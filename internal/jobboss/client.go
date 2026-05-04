package jobboss

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	defaultHTTPTimeout           = 30 * time.Second
	defaultResponseBodySize      = 10 * 1024 * 1024
	defaultAPIServerURL          = "https://api-jb2.integrations.ecimanufacturing.com"
	defaultAuthServerURL         = "https://api-user.integrations.ecimanufacturing.com"
	defaultTokenEndpointPath     = "/oauth2/api-user/token"
	defaultTokenScope            = "openid"
	defaultTokenEarlyRefreshSkew = 5 * time.Minute
	defaultMaxAuthRetryOn401     = 2
)

// Config contains settings for the JobBOSS API client.
type Config struct {
	BaseURL      string
	AuthBaseURL  string
	TokenURL     string
	ClientID     string
	ClientSecret string
	Scope        string
	// TokenEarlyRefreshSkew controls how long before expiry a token is considered stale.
	TokenEarlyRefreshSkew time.Duration
	// MaxAuthRetryOn401 controls retries after the initial request when the API returns 401.
	MaxAuthRetryOn401 int
	HTTPClient        *http.Client
}

// Client is a low-level JobBOSS API client with OAuth2 bearer-token support.
type Client struct {
	baseURL       string
	tokenURL      string
	clientID      string
	clientSecret  string
	scope         string
	httpClient    *http.Client
	refreshSkew   time.Duration
	max401Retries int

	mu          sync.RWMutex
	accessToken string
	tokenExpiry time.Time
}

// TokenResponse is the OAuth2 token response payload.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// APIError is returned when the API responds with status >= 400.
type APIError struct {
	StatusCode int
	URL        string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("jobboss API error %d on %s: %s", e.StatusCode, e.URL, e.Body)
}

// NewClient creates a new low-level JobBOSS API client.
func NewClient(cfg Config) *Client {
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultHTTPTimeout}
	}

	baseURL := strings.TrimSpace(cfg.BaseURL)
	if baseURL == "" {
		baseURL = defaultAPIServerURL
	}

	tokenURL := strings.TrimSpace(cfg.TokenURL)
	if tokenURL == "" {
		authBaseURL := strings.TrimSpace(cfg.AuthBaseURL)
		if authBaseURL == "" {
			authBaseURL = defaultAuthServerURL
		}
		resolvedTokenURL, err := buildAPIURL(authBaseURL, defaultTokenEndpointPath)
		if err == nil {
			tokenURL = resolvedTokenURL
		}
	}

	scope := strings.TrimSpace(cfg.Scope)
	if scope == "" {
		scope = defaultTokenScope
	}

	refreshSkew := cfg.TokenEarlyRefreshSkew
	if refreshSkew <= 0 {
		refreshSkew = defaultTokenEarlyRefreshSkew
	}

	max401Retries := cfg.MaxAuthRetryOn401
	if max401Retries == 0 {
		max401Retries = defaultMaxAuthRetryOn401
	}
	if max401Retries < 0 {
		max401Retries = 0
	}

	return &Client{
		baseURL:       baseURL,
		tokenURL:      tokenURL,
		clientID:      strings.TrimSpace(cfg.ClientID),
		clientSecret:  strings.TrimSpace(cfg.ClientSecret),
		scope:         scope,
		httpClient:    httpClient,
		refreshSkew:   refreshSkew,
		max401Retries: max401Retries,
	}
}

// BaseURL returns the configured JobBOSS API base URL.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// TokenURL returns the configured auth token URL.
func (c *Client) TokenURL() string {
	return c.tokenURL
}

// AccessToken returns a valid bearer token, refreshing it when required.
func (c *Client) AccessToken(ctx context.Context) (string, error) {
	c.mu.RLock()
	if c.accessToken != "" && time.Now().Before(c.tokenExpiry) {
		tok := c.accessToken
		c.mu.RUnlock()
		return tok, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.tokenURL == "" || c.clientID == "" || c.clientSecret == "" {
		return "", fmt.Errorf("missing token configuration (token URL / client id / client secret)")
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", c.clientID)
	form.Set("client_secret", c.clientSecret)
	form.Set("scope", c.scope)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", &APIError{StatusCode: resp.StatusCode, URL: req.URL.String(), Body: strings.TrimSpace(string(b))}
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("decode token response: %w", err)
	}
	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("token response did not include access_token")
	}

	expiresIn := time.Duration(tokenResp.ExpiresIn) * time.Second
	if expiresIn <= 0 {
		expiresIn = 15 * time.Minute
	}
	expiry := time.Now().Add(expiresIn - c.refreshSkew)
	if expiry.Before(time.Now()) {
		expiry = time.Now().Add(expiresIn / 2)
	}

	c.accessToken = tokenResp.AccessToken
	c.tokenExpiry = expiry
	return c.accessToken, nil
}

func (c *Client) refreshAccessToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	c.accessToken = ""
	c.tokenExpiry = time.Time{}
	c.mu.Unlock()
	return c.AccessToken(ctx)
}

func buildAPIURL(baseURL, endpoint string) (string, error) {
	baseURL = strings.TrimSpace(baseURL)
	endpoint = strings.TrimSpace(endpoint)
	if baseURL == "" {
		return "", fmt.Errorf("base URL is empty")
	}
	if endpoint == "" {
		endpoint = "/"
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL %q: %w", baseURL, err)
	}
	rel, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("invalid endpoint %q: %w", endpoint, err)
	}
	if rel.IsAbs() {
		return rel.String(), nil
	}

	basePath := strings.TrimSuffix(base.EscapedPath(), "/")
	relPath := rel.EscapedPath()
	if relPath == "" {
		relPath = "/"
	}
	if !strings.HasPrefix(relPath, "/") {
		relPath = "/" + relPath
	}

	if strings.HasSuffix(basePath, "/api/v1") && strings.HasPrefix(relPath, "/api/v1/") {
		relPath = strings.TrimPrefix(relPath, "/api/v1")
	} else if strings.HasSuffix(basePath, "/api") && strings.HasPrefix(relPath, "/api/") {
		relPath = strings.TrimPrefix(relPath, "/api")
	}

	if basePath == "" {
		base.Path = relPath
	} else {
		base.Path = basePath + relPath
	}
	base.RawQuery = rel.RawQuery
	base.Fragment = rel.Fragment
	return base.String(), nil
}

// Do sends a raw HTTP request to the JobBOSS API and returns response bytes.
func (c *Client) Do(ctx context.Context, method, endpoint string, query url.Values, body any) ([]byte, error) {
	if c.baseURL == "" {
		return nil, fmt.Errorf("base URL is not configured")
	}

	requestURL, err := buildAPIURL(c.baseURL, endpoint)
	if err != nil {
		return nil, err
	}

	// buildAPIURL already validates and normalizes URL input.
	parsedURL, _ := url.Parse(requestURL)
	if query != nil {
		q := parsedURL.Query()
		for key, values := range query {
			for _, v := range values {
				q.Add(key, v)
			}
		}
		parsedURL.RawQuery = q.Encode()
	}

	var bodyBytes []byte
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
	}

	var token string
	for attempt := 0; attempt <= c.max401Retries; attempt++ {
		if attempt == 0 {
			token, err = c.AccessToken(ctx)
		} else {
			token, err = c.refreshAccessToken(ctx)
		}
		if err != nil {
			return nil, err
		}

		var requestBody io.Reader
		if len(bodyBytes) > 0 {
			requestBody = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, parsedURL.String(), requestBody)
		if err != nil {
			return nil, fmt.Errorf("create API request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "go-jobboss2-api")
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("API request failed: %w", err)
		}

		respBody, readErr := io.ReadAll(io.LimitReader(resp.Body, defaultResponseBodySize))
		resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("read API response: %w", readErr)
		}

		if resp.StatusCode == http.StatusUnauthorized && attempt < c.max401Retries {
			continue
		}

		if resp.StatusCode >= 400 {
			return nil, &APIError{
				StatusCode: resp.StatusCode,
				URL:        req.URL.String(),
				Body:       strings.TrimSpace(string(respBody)),
			}
		}

		return respBody, nil
	}

	return nil, fmt.Errorf("request failed after authentication retries")
}

func (c *Client) Get(ctx context.Context, endpoint string, query url.Values) ([]byte, error) {
	return c.Do(ctx, http.MethodGet, endpoint, query, nil)
}

func (c *Client) Post(ctx context.Context, endpoint string, body any) ([]byte, error) {
	return c.Do(ctx, http.MethodPost, endpoint, nil, body)
}

func (c *Client) Patch(ctx context.Context, endpoint string, body any) ([]byte, error) {
	return c.Do(ctx, http.MethodPatch, endpoint, nil, body)
}

func (c *Client) Delete(ctx context.Context, endpoint string, query url.Values) ([]byte, error) {
	return c.Do(ctx, http.MethodDelete, endpoint, query, nil)
}

// GetJSON sends a GET request and unmarshals the JSON response into out.
func GetJSON[T any](ctx context.Context, c *Client, endpoint string, query url.Values) (T, error) {
	var out T
	body, err := c.Get(ctx, endpoint, query)
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return out, fmt.Errorf("decode GET response: %w", err)
	}
	return out, nil
}

// PostJSON sends a POST request and unmarshals the JSON response into out.
func PostJSON[Req any, Resp any](ctx context.Context, c *Client, endpoint string, payload Req) (Resp, error) {
	var out Resp
	body, err := c.Post(ctx, endpoint, payload)
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return out, fmt.Errorf("decode POST response: %w", err)
	}
	return out, nil
}

// PatchJSON sends a PATCH request and unmarshals the JSON response into out.
func PatchJSON[Req any, Resp any](ctx context.Context, c *Client, endpoint string, payload Req) (Resp, error) {
	var out Resp
	body, err := c.Patch(ctx, endpoint, payload)
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return out, fmt.Errorf("decode PATCH response: %w", err)
	}
	return out, nil
}
