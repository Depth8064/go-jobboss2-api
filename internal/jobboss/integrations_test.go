package jobboss

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestIntegrationsGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/oauth2Exchange":
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[OAuth2ExchangeResponse]{Data: OAuth2ExchangeResponse{AccessToken: "abc", TokenType: "Bearer", ExpiresIn: 3600}})
		case "/api/v1/eci-aps/authenticate-user":
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[ECIAPSAuthenticateResponse]{Data: ECIAPSAuthenticateResponse{Authenticated: true, Token: "eci-token"}})
		case "/api/v1/eci-aps/get-schedule":
			_ = json.NewEncoder(w).Encode(ListResponse[ECIAPSScheduleItem]{Data: []ECIAPSScheduleItem{{JobNumber: "J1", WorkCenter: "WC1"}}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := NewClient(Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})

	oauthResp, err := c.ExchangeOAuth2(context.Background(), OAuth2ExchangeRequest{GrantType: "client_credentials"})
	if err != nil || oauthResp.Data.AccessToken != "abc" {
		t.Fatalf("ExchangeOAuth2 failed: resp=%+v err=%v", oauthResp, err)
	}
	authResp, err := c.AuthenticateECIAPS(context.Background(), ECIAPSAuthenticateRequest{Username: "u", Password: "p"})
	if err != nil || !authResp.Data.Authenticated {
		t.Fatalf("AuthenticateECIAPS failed: resp=%+v err=%v", authResp, err)
	}

	q := url.Values{}
	q.Set("fields", "jobNumber,workCenter")
	schedule, err := c.GetECIAPSSchedule(context.Background(), q)
	if err != nil || len(schedule.Data) != 1 || schedule.Data[0].JobNumber != "J1" {
		t.Fatalf("GetECIAPSSchedule failed: resp=%+v err=%v", schedule, err)
	}
}
