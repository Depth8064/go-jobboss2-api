package jobboss

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, errors.New("forced read error") }

func TestAPIErrorError(t *testing.T) {
	err := (&APIError{StatusCode: 400, URL: "https://x", Body: "bad"}).Error()
	if !strings.Contains(err, "400") || !strings.Contains(err, "https://x") || !strings.Contains(err, "bad") {
		t.Fatalf("unexpected error string: %s", err)
	}
}

func TestNewClientDefaultsAndGetters(t *testing.T) {
	c := NewClient(Config{})
	if c.BaseURL() != "https://api-jb2.integrations.ecimanufacturing.com" {
		t.Fatalf("unexpected default base URL: %q", c.BaseURL())
	}
	if c.TokenURL() != "https://api-user.integrations.ecimanufacturing.com/oauth2/api-user/token" {
		t.Fatalf("unexpected default token URL: %q", c.TokenURL())
	}
}

func TestNewClientNegativeRetriesClamped(t *testing.T) {
	c := NewClient(Config{MaxAuthRetryOn401: -1})
	if c.max401Retries != 0 {
		t.Fatalf("expected negative retry count to clamp to 0, got %d", c.max401Retries)
	}
}

func TestClientMethodsAndJSONHelpers(t *testing.T) {
	tokenCalls := 0
	apiCalls := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			tokenCalls++
			tok := "token-1"
			if tokenCalls > 1 {
				tok = "token-2"
			}
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: tok, ExpiresIn: 3600})
		case "/retry":
			apiCalls++
			if r.Header.Get("Authorization") != "Bearer token-2" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			_, _ = w.Write([]byte(`{"value":"ok"}`))
		case "/g":
			_, _ = w.Write([]byte(`{"value":"g"}`))
		case "/p":
			if r.Method != http.MethodPost {
				t.Fatalf("expected POST")
			}
			_, _ = w.Write([]byte(`{"value":"p"}`))
		case "/pt":
			if r.Method != http.MethodPatch {
				t.Fatalf("expected PATCH")
			}
			_, _ = w.Write([]byte(`{"value":"pt"}`))
		case "/d":
			if r.Method != http.MethodDelete {
				t.Fatalf("expected DELETE")
			}
			_, _ = w.Write([]byte(`{"value":"d"}`))
		case "/bad":
			_, _ = w.Write([]byte("{"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := NewClient(Config{
		BaseURL:           ts.URL,
		TokenURL:          ts.URL + "/oauth2/api-user/token",
		ClientID:          "id",
		ClientSecret:      "secret",
		HTTPClient:        ts.Client(),
		MaxAuthRetryOn401: 2,
	})

	if _, err := c.Get(context.Background(), "/retry", nil); err != nil {
		t.Fatalf("expected retry success: %v", err)
	}
	if tokenCalls < 2 || apiCalls < 2 {
		t.Fatalf("expected auth refresh path to run, tokenCalls=%d apiCalls=%d", tokenCalls, apiCalls)
	}

	if _, err := c.Get(context.Background(), "/g", nil); err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if _, err := c.Post(context.Background(), "/p", map[string]string{"a": "b"}); err != nil {
		t.Fatalf("Post failed: %v", err)
	}
	if _, err := c.Patch(context.Background(), "/pt", map[string]string{"a": "b"}); err != nil {
		t.Fatalf("Patch failed: %v", err)
	}
	if _, err := c.Delete(context.Background(), "/d", nil); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	type payload struct {
		Value string `json:"value"`
	}
	outG, err := GetJSON[payload](context.Background(), c, "/g", nil)
	if err != nil || outG.Value != "g" {
		t.Fatalf("GetJSON failed: out=%+v err=%v", outG, err)
	}
	if _, err := GetJSON[payload](context.Background(), c, "/bad", nil); err == nil {
		t.Fatalf("expected GetJSON decode error")
	}
	if _, err := GetJSON[payload](context.Background(), c, "/missing", nil); err == nil {
		t.Fatalf("expected GetJSON request error")
	}

	outP, err := PostJSON[map[string]string, payload](context.Background(), c, "/p", map[string]string{"x": "1"})
	if err != nil || outP.Value != "p" {
		t.Fatalf("PostJSON failed: out=%+v err=%v", outP, err)
	}
	if _, err := PostJSON[map[string]string, payload](context.Background(), c, "/bad", map[string]string{}); err == nil {
		t.Fatalf("expected PostJSON decode error")
	}

	outPt, err := PatchJSON[map[string]string, payload](context.Background(), c, "/pt", map[string]string{"x": "1"})
	if err != nil || outPt.Value != "pt" {
		t.Fatalf("PatchJSON failed: out=%+v err=%v", outPt, err)
	}
	if _, err := PatchJSON[map[string]string, payload](context.Background(), c, "/bad", map[string]string{}); err == nil {
		t.Fatalf("expected PatchJSON decode error")
	}
	if _, err := PatchJSON[map[string]string, payload](context.Background(), c, "/missing", map[string]string{}); err == nil {
		t.Fatalf("expected PatchJSON request error")
	}

	if _, err := c.Do(context.Background(), http.MethodPost, "/p", nil, map[string]any{"bad": make(chan int)}); err == nil {
		t.Fatalf("expected marshal error")
	}

	q := url.Values{"a": {"1"}}
	if _, err := c.Do(context.Background(), "BAD\nMETHOD", "/g", q, nil); err == nil {
		t.Fatalf("expected request creation error")
	}
}

func TestBuildAPIURLBranches(t *testing.T) {
	if _, err := buildAPIURL("", "/x"); err == nil {
		t.Fatalf("expected empty base URL error")
	}

	if _, err := buildAPIURL("://bad", "/x"); err == nil {
		t.Fatalf("expected invalid base URL error")
	}

	if _, err := buildAPIURL("https://example.test", "http://[::1"); err == nil {
		t.Fatalf("expected invalid endpoint error")
	}

	abs, err := buildAPIURL("https://example.test/api", "https://other.test/path?q=1#f")
	if err != nil || abs != "https://other.test/path?q=1#f" {
		t.Fatalf("expected absolute endpoint passthrough, got %q err=%v", abs, err)
	}

	root, err := buildAPIURL("https://example.test", "")
	if err != nil || root != "https://example.test/" {
		t.Fatalf("expected root URL, got %q err=%v", root, err)
	}

	v1, err := buildAPIURL("https://example.test/api/v1", "/api/v1/customers")
	if err != nil || v1 != "https://example.test/api/v1/customers" {
		t.Fatalf("expected deduplicated /api/v1 path, got %q err=%v", v1, err)
	}

	api, err := buildAPIURL("https://example.test/api", "/api/devices")
	if err != nil || api != "https://example.test/api/devices" {
		t.Fatalf("expected deduplicated /api path, got %q err=%v", api, err)
	}

	onlyQuery, err := buildAPIURL("https://example.test", "?q=1")
	if err != nil || onlyQuery != "https://example.test/?q=1" {
		t.Fatalf("expected query-only endpoint handling, got %q err=%v", onlyQuery, err)
	}

	relativePath, err := buildAPIURL("https://example.test/api", "customers")
	if err != nil || relativePath != "https://example.test/api/customers" {
		t.Fatalf("expected missing-leading-slash normalization, got %q err=%v", relativePath, err)
	}
}

func TestAccessTokenBranches(t *testing.T) {
	t.Run("missing token config", func(t *testing.T) {
		c := NewClient(Config{})
		if _, err := c.AccessToken(context.Background()); err == nil {
			t.Fatalf("expected missing token configuration error")
		}
	})

	t.Run("token request failed", func(t *testing.T) {
		c := NewClient(Config{
			BaseURL:      "https://api.test",
			TokenURL:     "https://auth.test/token",
			ClientID:     "id",
			ClientSecret: "secret",
			HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return nil, errors.New("network down")
			})},
		})

		if _, err := c.AccessToken(context.Background()); err == nil {
			t.Fatalf("expected token request transport error")
		}
	})

	t.Run("token request creation failed", func(t *testing.T) {
		c := &Client{
			tokenURL:     "://bad",
			clientID:     "id",
			clientSecret: "secret",
			httpClient:   &http.Client{},
		}

		if _, err := c.AccessToken(context.Background()); err == nil {
			t.Fatalf("expected token request creation error")
		}
	})

	t.Run("token status error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("oops"))
		}))
		defer ts.Close()

		c := NewClient(Config{BaseURL: ts.URL, TokenURL: ts.URL, ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client()})
		if _, err := c.AccessToken(context.Background()); err == nil {
			t.Fatalf("expected API error for non-200 token response")
		}
	})

	t.Run("decode and token value errors", func(t *testing.T) {
		t.Run("decode error", func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("{"))
			}))
			defer ts.Close()

			c := NewClient(Config{BaseURL: ts.URL, TokenURL: ts.URL, ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client()})
			if _, err := c.AccessToken(context.Background()); err == nil {
				t.Fatalf("expected decode error")
			}
		})

		t.Run("missing access token", func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(TokenResponse{ExpiresIn: 60})
			}))
			defer ts.Close()

			c := NewClient(Config{BaseURL: ts.URL, TokenURL: ts.URL, ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client()})
			if _, err := c.AccessToken(context.Background()); err == nil {
				t.Fatalf("expected missing access_token error")
			}
		})
	})

	t.Run("cached token with fallback expiry", func(t *testing.T) {
		calls := 0
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 0})
		}))
		defer ts.Close()

		c := NewClient(Config{
			BaseURL:               ts.URL,
			TokenURL:              ts.URL,
			ClientID:              "id",
			ClientSecret:          "secret",
			HTTPClient:            ts.Client(),
			TokenEarlyRefreshSkew: time.Hour,
		})

		t1, err := c.AccessToken(context.Background())
		if err != nil || t1 != "tok" {
			t.Fatalf("expected token, got %q err=%v", t1, err)
		}

		t2, err := c.AccessToken(context.Background())
		if err != nil || t2 != "tok" {
			t.Fatalf("expected cached token, got %q err=%v", t2, err)
		}

		if calls != 1 {
			t.Fatalf("expected token endpoint to be called once, got %d", calls)
		}
	})
}

func TestDoErrorBranches(t *testing.T) {
	t.Run("base URL not configured", func(t *testing.T) {
		c := &Client{httpClient: &http.Client{}}
		if _, err := c.Do(context.Background(), http.MethodGet, "/x", nil, nil); err == nil {
			t.Fatalf("expected base URL error")
		}
	})

	t.Run("invalid API URL", func(t *testing.T) {
		c := &Client{baseURL: "://bad", httpClient: &http.Client{}}
		if _, err := c.Do(context.Background(), http.MethodGet, "/x", nil, nil); err == nil {
			t.Fatalf("expected URL build error")
		}
	})

	t.Run("API request failed", func(t *testing.T) {
		c := &Client{
			baseURL:     "https://api.test",
			httpClient:  &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, errors.New("boom") })},
			accessToken: "tok",
			tokenExpiry: time.Now().Add(time.Hour),
		}

		if _, err := c.Do(context.Background(), http.MethodGet, "/x", nil, nil); err == nil {
			t.Fatalf("expected API request transport error")
		}
	})

	t.Run("token acquisition fails inside Do", func(t *testing.T) {
		c := &Client{baseURL: "https://api.test", httpClient: &http.Client{}}
		if _, err := c.Do(context.Background(), http.MethodGet, "/x", nil, nil); err == nil {
			t.Fatalf("expected Do to return AccessToken error")
		}
	})

	t.Run("read API response error", func(t *testing.T) {
		c := &Client{
			baseURL: "https://api.test",
			httpClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(errReader{}),
					Header:     make(http.Header),
				}, nil
			})},
			accessToken: "tok",
			tokenExpiry: time.Now().Add(time.Hour),
		}

		if _, err := c.Do(context.Background(), http.MethodGet, "/x", nil, nil); err == nil {
			t.Fatalf("expected response read error")
		}
	})

	t.Run("exhausted loop fallback", func(t *testing.T) {
		c := &Client{baseURL: "https://api.test", httpClient: &http.Client{}, max401Retries: -1}
		if _, err := c.Do(context.Background(), http.MethodGet, "/x", nil, nil); err == nil || !strings.Contains(err.Error(), "authentication retries") {
			t.Fatalf("expected authentication retries fallback error, got %v", err)
		}
	})
}
