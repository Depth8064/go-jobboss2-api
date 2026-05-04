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
	defaultHTTPTimeout      = 30 * time.Second
	defaultResponseBodySize = 10 * 1024 * 1024
	tokenEarlyRefreshBuffer = 60 * time.Second
)

// Config contains settings for the JobBOSS API client.
type Config struct {
	BaseURL      string
	TokenURL     string
	ClientID     string
	ClientSecret string
	HTTPClient   *http.Client
}

// Client is a low-level JobBOSS API client with OAuth2 bearer-token support.
type Client struct {
	baseURL      string
	tokenURL     string
	clientID     string
	clientSecret string
	httpClient   *http.Client

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

	return &Client{
		baseURL:      strings.TrimSpace(cfg.BaseURL),
		tokenURL:     strings.TrimSpace(cfg.TokenURL),
		clientID:     strings.TrimSpace(cfg.ClientID),
		clientSecret: strings.TrimSpace(cfg.ClientSecret),
		httpClient:   httpClient,
	}
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

	if c.accessToken != "" && time.Now().Before(c.tokenExpiry) {
		return c.accessToken, nil
	}

	if c.tokenURL == "" || c.clientID == "" || c.clientSecret == "" {
		return "", fmt.Errorf("missing token configuration (token URL / client id / client secret)")
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", c.clientID)
	form.Set("client_secret", c.clientSecret)
	form.Set("scope", "openid")

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
	expiry := time.Now().Add(expiresIn - tokenEarlyRefreshBuffer)
	if expiry.Before(time.Now()) {
		expiry = time.Now().Add(expiresIn / 2)
	}

	c.accessToken = tokenResp.AccessToken
	c.tokenExpiry = expiry
	return c.accessToken, nil
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

	token, err := c.AccessToken(ctx)
	if err != nil {
		return nil, err
	}

	requestURL, err := buildAPIURL(c.baseURL, endpoint)
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(requestURL)
	if err != nil {
		return nil, fmt.Errorf("parse request URL: %w", err)
	}
	if query != nil {
		q := parsedURL.Query()
		for key, values := range query {
			for _, v := range values {
				q.Add(key, v)
			}
		}
		parsedURL.RawQuery = q.Encode()
	}

	var requestBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		requestBody = bytes.NewReader(jsonBody)
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
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, defaultResponseBodySize))
	if err != nil {
		return nil, fmt.Errorf("read API response: %w", err)
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
