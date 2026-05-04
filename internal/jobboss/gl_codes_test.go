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

func TestGLCodesGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/gl-codes":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(ListResponse[GLCode]{
				Data: []GLCode{{GLAccountNumber: "4000", Description: "Sales"}},
			})
		case "/api/v1/gl-codes/4000":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[GLCode]{
				Data: GLCode{GLAccountNumber: "4000", Description: "Sales", Active: true},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := NewClient(Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})

	q := url.Values{}
	q.Set("fields", "GLAccountNumber,description")
	list, err := c.GetGLCodes(context.Background(), q)
	if err != nil || len(list.Data) != 1 || list.Data[0].GLAccountNumber != "4000" {
		t.Fatalf("GetGLCodes failed: resp=%+v err=%v", list, err)
	}

	one, err := c.GetGLCodeByAccountNumber(context.Background(), "4000", nil)
	if err != nil || one.Data.Description != "Sales" {
		t.Fatalf("GetGLCodeByAccountNumber failed: resp=%+v err=%v", one, err)
	}

	_, err = c.GetGLCodeByAccountNumber(context.Background(), "", nil)
	if err == nil || !strings.Contains(err.Error(), "gl account number is required") {
		t.Fatalf("expected gl account number validation error, got %v", err)
	}
}
