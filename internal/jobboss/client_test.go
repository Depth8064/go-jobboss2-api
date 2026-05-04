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

	t.Run("exhausted transient retry loop fallback", func(t *testing.T) {
		c := &Client{baseURL: "https://api.test", httpClient: &http.Client{}, maxRetries: -1}
		if _, err := c.Do(context.Background(), http.MethodGet, "/x", nil, nil); err == nil || !strings.Contains(err.Error(), "retries") {
			t.Fatalf("expected transient retries fallback error, got %v", err)
		}
	})
}

func TestPatchNoContentJSONAndPatchJSONEmptyBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/patch-204":
			if r.Method != http.MethodPatch {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		case "/patch-200":
			if r.Method != http.MethodPatch {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_, _ = w.Write([]byte(`{"Data":{"value":"ok"}}`))
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
		MaxAuthRetryOn401: 0,
	})

	meta, err := PatchNoContentJSON[map[string]string](context.Background(), c, "/patch-204", map[string]string{"x": "1"})
	if err != nil {
		t.Fatalf("PatchNoContentJSON failed: %v", err)
	}
	if meta.StatusCode != http.StatusNoContent || !meta.NoContent {
		t.Fatalf("unexpected PatchNoContentJSON metadata: %+v", meta)
	}

	metaWithBody, err := PatchNoContentJSON[map[string]string](context.Background(), c, "/patch-200", map[string]string{"x": "1"})
	if err != nil {
		t.Fatalf("PatchNoContentJSON with body failed: %v", err)
	}
	if metaWithBody.StatusCode != http.StatusOK || metaWithBody.NoContent {
		t.Fatalf("unexpected PatchNoContentJSON body metadata: %+v", metaWithBody)
	}

	if _, err := PatchNoContentJSON[map[string]string](context.Background(), c, "/missing", map[string]string{"x": "1"}); err == nil {
		t.Fatalf("expected PatchNoContentJSON request error")
	}

	type payload struct {
		Data struct {
			Value string `json:"value"`
		} `json:"Data"`
	}
	if _, err := PatchJSON[map[string]string, payload](context.Background(), c, "/patch-204", map[string]string{"x": "1"}); err == nil {
		t.Fatalf("expected PatchJSON empty-body error")
	}

	parsed, err := PatchJSON[map[string]string, payload](context.Background(), c, "/patch-200", map[string]string{"x": "1"})
	if err != nil {
		t.Fatalf("PatchJSON with body failed: %v", err)
	}
	if parsed.Data.Value != "ok" {
		t.Fatalf("unexpected PatchJSON parsed value: %+v", parsed)
	}
}

func TestNewClientRetryConfig(t *testing.T) {
	// Non-default values should be respected.
	c := NewClient(Config{
		MaxRetries:       5,
		RetryBackoffBase: 200 * time.Millisecond,
		RetryBackoffMax:  10 * time.Second,
	})
	if c.maxRetries != 5 {
		t.Fatalf("expected maxRetries=5, got %d", c.maxRetries)
	}
	if c.retryBackoffBase != 200*time.Millisecond {
		t.Fatalf("expected retryBackoffBase=200ms, got %v", c.retryBackoffBase)
	}
	if c.retryBackoffMax != 10*time.Second {
		t.Fatalf("expected retryBackoffMax=10s, got %v", c.retryBackoffMax)
	}

	// MaxRetries: 0 → uses default.
	c0 := NewClient(Config{})
	if c0.maxRetries != defaultMaxRetries {
		t.Fatalf("expected maxRetries=%d for zero config, got %d", defaultMaxRetries, c0.maxRetries)
	}

	// MaxRetries: -1 → disabled (0 retries, 1 attempt total).
	cNeg := NewClient(Config{MaxRetries: -1})
	if cNeg.maxRetries != 0 {
		t.Fatalf("expected maxRetries=0 for negative input, got %d", cNeg.maxRetries)
	}
}

func TestParseRetryAfter(t *testing.T) {
	if parseRetryAfter("") != 0 {
		t.Fatal("expected 0 for empty header")
	}
	if parseRetryAfter("   ") != 0 {
		t.Fatal("expected 0 for whitespace-only header")
	}
	if parseRetryAfter("0") != 0 {
		t.Fatal("expected 0 for zero seconds (not positive)")
	}
	if d := parseRetryAfter("30"); d != 30*time.Second {
		t.Fatalf("expected 30s for integer header, got %v", d)
	}
	// HTTP-date in the future.
	future := time.Now().Add(10 * time.Minute).UTC().Format(http.TimeFormat)
	if d := parseRetryAfter(future); d <= 0 {
		t.Fatalf("expected positive duration for future HTTP-date, got %v", d)
	}
	// HTTP-date in the past.
	past := time.Now().Add(-10 * time.Minute).UTC().Format(http.TimeFormat)
	if d := parseRetryAfter(past); d != 0 {
		t.Fatalf("expected 0 for past HTTP-date, got %v", d)
	}
	// Unparseable string.
	if parseRetryAfter("not-a-date") != 0 {
		t.Fatal("expected 0 for unparseable header")
	}
}

func TestRetryableStatus(t *testing.T) {
	notRetryable := []int{200, 201, 204, 400, 401, 403, 404, 600}
	for _, code := range notRetryable {
		if retryableStatus(code) {
			t.Errorf("expected status %d to not be retryable", code)
		}
	}
	isRetryable := []int{429, 500, 502, 503, 504, 599}
	for _, code := range isRetryable {
		if !retryableStatus(code) {
			t.Errorf("expected status %d to be retryable", code)
		}
	}
}

func TestRetryDelay(t *testing.T) {
	// 1a: Retry-After header takes precedence over backoff.
	c := &Client{retryBackoffBase: 100 * time.Millisecond, retryBackoffMax: 500 * time.Millisecond}
	if d := c.retryDelay(1, "5"); d != 5*time.Second {
		t.Fatalf("expected 5s from Retry-After header, got %v", d)
	}

	// 1b + 3b (loop skipped) + 5b (base < max) + 6b (room > maxJitter) + 7b (maxJitter > 0):
	// base=100ms, max=500ms, attempt=1 → backoff=100ms, room=400ms, maxJitter=50ms → [100ms, 150ms].
	for i := 0; i < 20; i++ {
		d := c.retryDelay(1, "")
		if d < 100*time.Millisecond || d > 150*time.Millisecond {
			t.Fatalf("attempt 1 delay %v outside [100ms, 150ms]", d)
		}
	}

	// 2a: attempt <= 0 clamped to 1 — same result range.
	for i := 0; i < 20; i++ {
		d := c.retryDelay(0, "")
		if d < 100*time.Millisecond || d > 150*time.Millisecond {
			t.Fatalf("attempt 0 (clamped) delay %v outside [100ms, 150ms]", d)
		}
	}

	// 3a + 4b: loop entered, backoff doubles but stays below max.
	// base=100ms, max=500ms, attempt=2 → backoff=200ms, room=300ms, maxJitter=100ms → [200ms, 300ms].
	for i := 0; i < 20; i++ {
		d := c.retryDelay(2, "")
		if d < 200*time.Millisecond || d > 300*time.Millisecond {
			t.Fatalf("attempt 2 delay %v outside [200ms, 300ms]", d)
		}
	}

	// 3a + 4a + 6a + 7a: loop hits max cap → room=0 → maxJitter=0 → return exactly backoff.
	// base=100ms, max=200ms, attempt=3: 100→200→break. room=0 → maxJitter=0 → 200ms.
	c2 := &Client{retryBackoffBase: 100 * time.Millisecond, retryBackoffMax: 200 * time.Millisecond}
	if d := c2.retryDelay(3, ""); d != 200*time.Millisecond {
		t.Fatalf("expected exactly 200ms when backoff saturates max, got %v", d)
	}

	// 5a: base > max → post-loop cap (branch 5a), then room=0 → 7a → return max.
	// base=300ms, max=200ms, attempt=1: backoff=300ms > 200ms → cap to 200ms → return 200ms.
	c3 := &Client{retryBackoffBase: 300 * time.Millisecond, retryBackoffMax: 200 * time.Millisecond}
	if d := c3.retryDelay(1, ""); d != 200*time.Millisecond {
		t.Fatalf("expected 200ms when base > max, got %v", d)
	}

	// 6a (room < maxJitter) + 7b (maxJitter > 0):
	// base=180ms, max=200ms, attempt=1 → backoff=180ms, room=20ms, maxJitter=min(90ms,20ms)=20ms → [180ms, 200ms].
	c4 := &Client{retryBackoffBase: 180 * time.Millisecond, retryBackoffMax: 200 * time.Millisecond}
	for i := 0; i < 20; i++ {
		d := c4.retryDelay(1, "")
		if d < 180*time.Millisecond || d > 200*time.Millisecond {
			t.Fatalf("close-to-max delay %v outside [180ms, 200ms]", d)
		}
	}
}

func TestTransientRetry(t *testing.T) {
	t.Run("retries on 5xx and succeeds", func(t *testing.T) {
		callCount := 0
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/oauth2/api-user/token" {
				_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
				return
			}
			callCount++
			if callCount == 1 {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			_, _ = w.Write([]byte(`"ok"`))
		}))
		defer ts.Close()

		var slept []time.Duration
		c := NewClient(Config{
			BaseURL:          ts.URL,
			TokenURL:         ts.URL + "/oauth2/api-user/token",
			ClientID:         "id",
			ClientSecret:     "secret",
			HTTPClient:       ts.Client(),
			MaxRetries:       2,
			RetryBackoffBase: time.Nanosecond,
			RetryBackoffMax:  time.Millisecond,
		})
		c.sleepFn = func(d time.Duration) { slept = append(slept, d) }

		if _, err := c.Get(context.Background(), "/data", nil); err != nil {
			t.Fatalf("expected success after retry, got %v", err)
		}
		if len(slept) != 1 {
			t.Fatalf("expected 1 sleep call, got %d", len(slept))
		}
	})

	t.Run("retries on 429 with Retry-After and succeeds", func(t *testing.T) {
		callCount := 0
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/oauth2/api-user/token" {
				_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
				return
			}
			callCount++
			if callCount == 1 {
				w.Header().Set("Retry-After", "7")
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			_, _ = w.Write([]byte(`"done"`))
		}))
		defer ts.Close()

		var slept []time.Duration
		c := NewClient(Config{
			BaseURL:      ts.URL,
			TokenURL:     ts.URL + "/oauth2/api-user/token",
			ClientID:     "id",
			ClientSecret: "secret",
			HTTPClient:   ts.Client(),
			MaxRetries:   2,
		})
		c.sleepFn = func(d time.Duration) { slept = append(slept, d) }

		if _, err := c.Get(context.Background(), "/data", nil); err != nil {
			t.Fatalf("expected success after 429 retry, got %v", err)
		}
		if len(slept) != 1 || slept[0] != 7*time.Second {
			t.Fatalf("expected one 7s sleep from Retry-After, got %v", slept)
		}
	})

	t.Run("retries on network error and succeeds", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/oauth2/api-user/token" {
				_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
				return
			}
			_, _ = w.Write([]byte(`"ok"`))
		}))
		defer ts.Close()

		apiCalls := 0
		var slept []time.Duration
		transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.Path == "/oauth2/api-user/token" {
				return ts.Client().Transport.RoundTrip(r)
			}
			apiCalls++
			if apiCalls <= 2 {
				return nil, errors.New("network down")
			}
			return ts.Client().Transport.RoundTrip(r)
		})

		c := NewClient(Config{
			BaseURL:          ts.URL,
			TokenURL:         ts.URL + "/oauth2/api-user/token",
			ClientID:         "id",
			ClientSecret:     "secret",
			HTTPClient:       &http.Client{Transport: transport},
			MaxRetries:       3,
			RetryBackoffBase: time.Nanosecond,
			RetryBackoffMax:  time.Millisecond,
		})
		c.sleepFn = func(d time.Duration) { slept = append(slept, d) }

		if _, err := c.Get(context.Background(), "/data", nil); err != nil {
			t.Fatalf("expected success after network retries, got %v", err)
		}
		if apiCalls != 3 {
			t.Fatalf("expected 3 API calls, got %d", apiCalls)
		}
		if len(slept) != 2 {
			t.Fatalf("expected 2 sleep calls, got %d", len(slept))
		}
	})

	t.Run("exhausts retries on persistent 5xx", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/oauth2/api-user/token" {
				_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		sleepCalls := 0
		c := NewClient(Config{
			BaseURL:          ts.URL,
			TokenURL:         ts.URL + "/oauth2/api-user/token",
			ClientID:         "id",
			ClientSecret:     "secret",
			HTTPClient:       ts.Client(),
			MaxRetries:       2,
			RetryBackoffBase: time.Nanosecond,
			RetryBackoffMax:  time.Millisecond,
		})
		c.sleepFn = func(d time.Duration) { sleepCalls++ }

		if _, err := c.Get(context.Background(), "/data", nil); err == nil {
			t.Fatal("expected error after exhausted retries")
		}
		if sleepCalls != 2 {
			t.Fatalf("expected 2 sleep calls for 2 retries, got %d", sleepCalls)
		}
	})

	t.Run("exhausts retries on persistent network error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/oauth2/api-user/token" {
				_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
				return
			}
		}))
		defer ts.Close()

		apiCalls := 0
		transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.Path == "/oauth2/api-user/token" {
				return ts.Client().Transport.RoundTrip(r)
			}
			apiCalls++
			return nil, errors.New("always down")
		})

		c := NewClient(Config{
			BaseURL:          ts.URL,
			TokenURL:         ts.URL + "/oauth2/api-user/token",
			ClientID:         "id",
			ClientSecret:     "secret",
			HTTPClient:       &http.Client{Transport: transport},
			MaxRetries:       2,
			RetryBackoffBase: time.Nanosecond,
			RetryBackoffMax:  time.Millisecond,
		})
		sleepCalls := 0
		c.sleepFn = func(d time.Duration) { sleepCalls++ }

		if _, err := c.Get(context.Background(), "/data", nil); err == nil {
			t.Fatal("expected error after exhausted network retries")
		}
		if apiCalls != 3 {
			t.Fatalf("expected 3 API calls (1 initial + 2 retries), got %d", apiCalls)
		}
		if sleepCalls != 2 {
			t.Fatalf("expected 2 sleep calls, got %d", sleepCalls)
		}
	})

	t.Run("429 RetryAfter field set on final error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/oauth2/api-user/token" {
				_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
				return
			}
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusTooManyRequests)
		}))
		defer ts.Close()

		c := NewClient(Config{
			BaseURL:      ts.URL,
			TokenURL:     ts.URL + "/oauth2/api-user/token",
			ClientID:     "id",
			ClientSecret: "secret",
			HTTPClient:   ts.Client(),
			MaxRetries:   -1, // disabled → 0 retries, 1 attempt
		})
		c.sleepFn = func(d time.Duration) {}

		_, err := c.Get(context.Background(), "/data", nil)
		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected *APIError, got %T: %v", err, err)
		}
		if apiErr.RetryAfter != "60" {
			t.Fatalf("expected RetryAfter=%q, got %q", "60", apiErr.RetryAfter)
		}
	})
}
