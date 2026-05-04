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

func TestCorrectiveActionsGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/corrective-preventive-actions":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(ListResponse[CorrectiveAction]{
				Data: []CorrectiveAction{{CorrectiveActionNumber: "1001", Status: "OPEN", CustomerCode: "ACME"}},
			})
		case "/api/v1/corrective-preventive-actions/1001":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[CorrectiveAction]{
				Data: CorrectiveAction{CorrectiveActionNumber: "1001", Status: "OPEN", Description: "Inspection issue"},
			})
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

	q := url.Values{}
	q.Set("fields", "correctiveActionNumber,status,customerCode")
	list, err := c.GetCorrectiveActions(context.Background(), q)
	if err != nil || len(list.Data) != 1 || list.Data[0].CorrectiveActionNumber != "1001" {
		t.Fatalf("GetCorrectiveActions failed: resp=%+v err=%v", list, err)
	}

	one, err := c.GetCorrectiveActionByNumber(context.Background(), "1001", nil)
	if err != nil || one.Data.Description != "Inspection issue" {
		t.Fatalf("GetCorrectiveActionByNumber failed: resp=%+v err=%v", one, err)
	}

	_, err = c.GetCorrectiveActionByNumber(context.Background(), "", nil)
	if err == nil || !strings.Contains(err.Error(), "corrective action number is required") {
		t.Fatalf("expected corrective action number validation error, got %v", err)
	}
}
