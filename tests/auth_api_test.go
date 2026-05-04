package jobboss_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	jobboss "go-jobboss2-api"
)

func TestCreateDevice(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/oauth2/api-user/token" {
				_ = json.NewEncoder(w).Encode(jobboss.TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
				return
			}
			if r.URL.Path != "/api/devices" || r.Method != http.MethodPost {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			_ = json.NewEncoder(w).Encode(jobboss.DeviceCreateResponse{ID: "abc", Code: "123"})
		}))
		defer ts.Close()

		c := jobboss.NewClient(jobboss.Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})
		resp, err := c.CreateDevice(context.Background(), jobboss.DeviceCreateModel{ClientID: "id"})
		if err != nil || resp == nil || resp.ID != "abc" {
			t.Fatalf("expected created device response, got resp=%+v err=%v", resp, err)
		}
	})

	t.Run("error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/oauth2/api-user/token" {
				_ = json.NewEncoder(w).Encode(jobboss.TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("bad request"))
		}))
		defer ts.Close()

		c := jobboss.NewClient(jobboss.Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})
		_, err := c.CreateDevice(context.Background(), jobboss.DeviceCreateModel{ClientID: "id"})
		if err == nil {
			t.Fatalf("expected error")
		}
	})
}

func TestTryActivateDevice(t *testing.T) {
	t.Run("missing id", func(t *testing.T) {
		c := jobboss.NewClient(jobboss.Config{})
		_, activated, err := c.TryActivateDevice(context.Background(), "  ", jobboss.DeviceActivateRequestModel{Code: "x"})
		if err == nil || activated {
			t.Fatalf("expected missing id error")
		}
	})

	t.Run("request error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/oauth2/api-user/token" {
				_ = json.NewEncoder(w).Encode(jobboss.TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("nope"))
		}))
		defer ts.Close()
		c := jobboss.NewClient(jobboss.Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})
		_, _, err := c.TryActivateDevice(context.Background(), "id", jobboss.DeviceActivateRequestModel{Code: "x"})
		if err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("not activated 204 empty body", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/oauth2/api-user/token" {
				_ = json.NewEncoder(w).Encode(jobboss.TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
				return
			}
			w.WriteHeader(http.StatusNoContent)
		}))
		defer ts.Close()
		c := jobboss.NewClient(jobboss.Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})
		resp, activated, err := c.TryActivateDevice(context.Background(), "id", jobboss.DeviceActivateRequestModel{Code: "x"})
		if err != nil || activated || resp != nil {
			t.Fatalf("expected no activation without error, got resp=%+v activated=%v err=%v", resp, activated, err)
		}
	})

	t.Run("decode error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/oauth2/api-user/token" {
				_ = json.NewEncoder(w).Encode(jobboss.TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{"))
		}))
		defer ts.Close()
		c := jobboss.NewClient(jobboss.Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})
		_, _, err := c.TryActivateDevice(context.Background(), "id", jobboss.DeviceActivateRequestModel{Code: "x"})
		if err == nil {
			t.Fatalf("expected decode error")
		}
	})

	t.Run("activated", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/oauth2/api-user/token" {
				_ = json.NewEncoder(w).Encode(jobboss.TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
				return
			}
			_ = json.NewEncoder(w).Encode(jobboss.DeviceActivatedResponse{ID: "id", IsActivated: true})
		}))
		defer ts.Close()
		c := jobboss.NewClient(jobboss.Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})
		resp, activated, err := c.TryActivateDevice(context.Background(), "id", jobboss.DeviceActivateRequestModel{Code: "x"})
		if err != nil || !activated || resp == nil || resp.ID != "id" {
			t.Fatalf("expected activated response, got resp=%+v activated=%v err=%v", resp, activated, err)
		}
	})
}
