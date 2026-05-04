package jobboss

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestNonConformancesGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/non-conformances":
			_ = json.NewEncoder(w).Encode(ListResponse[NonConformance]{Data: []NonConformance{{NCNumber: "NC-1", Status: "OPEN"}}})
		case "/api/v1/non-conformances/NC-1":
			_ = json.NewEncoder(w).Encode(SingleResponse[NonConformance]{Data: NonConformance{NCNumber: "NC-1", Status: "OPEN", NCDescription: "Scratch"}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := NewClient(Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})
	q := url.Values{}
	q.Set("fields", "ncNumber,status")
	list, err := c.GetNonConformances(context.Background(), q)
	if err != nil || len(list.Data) != 1 || list.Data[0].NCNumber != "NC-1" {
		t.Fatalf("GetNonConformances failed: resp=%+v err=%v", list, err)
	}
	one, err := c.GetNonConformanceByNumber(context.Background(), "NC-1", nil)
	if err != nil || one.Data.NCDescription != "Scratch" {
		t.Fatalf("GetNonConformanceByNumber failed: resp=%+v err=%v", one, err)
	}
	_, err = c.GetNonConformanceByNumber(context.Background(), "", nil)
	if err == nil || !strings.Contains(err.Error(), "nc number is required") {
		t.Fatalf("expected nc number validation error, got %v", err)
	}
}
