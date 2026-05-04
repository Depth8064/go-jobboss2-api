package jobboss

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestRootFacadeWrappers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/g":
			_ = json.NewEncoder(w).Encode(map[string]string{"value": "g"})
		case "/p":
			_ = json.NewEncoder(w).Encode(map[string]string{"value": "p"})
		case "/pt":
			_ = json.NewEncoder(w).Encode(map[string]string{"value": "pt"})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := NewClient(Config{
		BaseURL:      ts.URL,
		TokenURL:     ts.URL + "/oauth2/api-user/token",
		ClientID:     "id",
		ClientSecret: "secret",
		HTTPClient:   ts.Client(),
	})

	q := NewQuery()
	q = ApplyPaging(q, 25, 10)
	q = ApplyFields(q, "a", "b")
	q = ApplySort(q, "+a", "-b")
	q = ApplyFilter(q, "x", FilterEqual, "1")
	q = ApplyInFilter(q, "id", "1", "2")
	if q.Get("take") == "" || q.Get("fields") == "" || q.Get("sort") == "" || q.Get("x[eq]") == "" || q.Get("id[in]") == "" {
		t.Fatalf("expected query wrapper output, got %v", q)
	}

	if got := FormatUTCDate(time.Date(2026, 5, 4, 12, 0, 0, 0, time.FixedZone("X", -7*3600))); got != "2026-05-04" {
		t.Fatalf("unexpected date wrapper output: %q", got)
	}
	if got := FormatUTCDateTime(time.Date(2026, 5, 4, 12, 0, 0, 0, time.FixedZone("X", -7*3600))); got != "2026-05-04T19:00:00Z" {
		t.Fatalf("unexpected datetime wrapper output: %q", got)
	}

	type payload struct {
		Value string `json:"value"`
	}

	outG, err := GetJSON[payload](context.Background(), c, "/g", url.Values{"z": {"1"}})
	if err != nil || outG.Value != "g" {
		t.Fatalf("GetJSON wrapper failed: out=%+v err=%v", outG, err)
	}

	outP, err := PostJSON[map[string]string, payload](context.Background(), c, "/p", map[string]string{"x": "1"})
	if err != nil || outP.Value != "p" {
		t.Fatalf("PostJSON wrapper failed: out=%+v err=%v", outP, err)
	}

	outPt, err := PatchJSON[map[string]string, payload](context.Background(), c, "/pt", map[string]string{"x": "1"})
	if err != nil || outPt.Value != "pt" {
		t.Fatalf("PatchJSON wrapper failed: out=%+v err=%v", outPt, err)
	}
}
