package jobboss_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	jobboss "go-jobboss2-api"
)

func TestAPIErrorError(t *testing.T) {
	err := (&jobboss.APIError{StatusCode: 400, URL: "https://x", Body: "bad"}).Error()
	if !strings.Contains(err, "400") || !strings.Contains(err, "https://x") || !strings.Contains(err, "bad") {
		t.Fatalf("unexpected error string: %s", err)
	}
}

func TestNewClientDefaultsAndGetters(t *testing.T) {
	c := jobboss.NewClient(jobboss.Config{})
	if c.BaseURL() != "https://api-jb2.integrations.ecimanufacturing.com" {
		t.Fatalf("unexpected default base URL: %q", c.BaseURL())
	}
	if c.TokenURL() != "https://api-user.integrations.ecimanufacturing.com/oauth2/api-user/token" {
		t.Fatalf("unexpected default token URL: %q", c.TokenURL())
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
			_ = json.NewEncoder(w).Encode(jobboss.TokenResponse{AccessToken: tok, ExpiresIn: 3600})
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

	c := jobboss.NewClient(jobboss.Config{
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
	outG, err := jobboss.GetJSON[payload](context.Background(), c, "/g", nil)
	if err != nil || outG.Value != "g" {
		t.Fatalf("GetJSON failed: out=%+v err=%v", outG, err)
	}
	if _, err := jobboss.GetJSON[payload](context.Background(), c, "/bad", nil); err == nil {
		t.Fatalf("expected GetJSON decode error")
	}
	if _, err := jobboss.GetJSON[payload](context.Background(), c, "/missing", nil); err == nil {
		t.Fatalf("expected GetJSON request error")
	}

	outP, err := jobboss.PostJSON[map[string]string, payload](context.Background(), c, "/p", map[string]string{"x": "1"})
	if err != nil || outP.Value != "p" {
		t.Fatalf("PostJSON failed: out=%+v err=%v", outP, err)
	}
	if _, err := jobboss.PostJSON[map[string]string, payload](context.Background(), c, "/bad", map[string]string{}); err == nil {
		t.Fatalf("expected PostJSON decode error")
	}

	outPt, err := jobboss.PatchJSON[map[string]string, payload](context.Background(), c, "/pt", map[string]string{"x": "1"})
	if err != nil || outPt.Value != "pt" {
		t.Fatalf("PatchJSON failed: out=%+v err=%v", outPt, err)
	}
	if _, err := jobboss.PatchJSON[map[string]string, payload](context.Background(), c, "/bad", map[string]string{}); err == nil {
		t.Fatalf("expected PatchJSON decode error")
	}
	if _, err := jobboss.PatchJSON[map[string]string, payload](context.Background(), c, "/missing", map[string]string{}); err == nil {
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
