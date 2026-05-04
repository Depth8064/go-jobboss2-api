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

func TestOperationCodesGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/operation-codes":
			_ = json.NewEncoder(w).Encode(ListResponse[OperationCode]{Data: []OperationCode{{OperationCode: "CUT", Description: "Cutting"}}})
		case "/api/v1/operation-codes/CUT":
			_ = json.NewEncoder(w).Encode(SingleResponse[OperationCode]{Data: OperationCode{OperationCode: "CUT", Description: "Cutting", Active: true}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := NewClient(Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})
	q := url.Values{}
	q.Set("fields", "operationCode,description")
	list, err := c.GetOperationCodes(context.Background(), q)
	if err != nil || len(list.Data) != 1 || list.Data[0].OperationCode != "CUT" {
		t.Fatalf("GetOperationCodes failed: resp=%+v err=%v", list, err)
	}
	one, err := c.GetOperationCodeByCode(context.Background(), "CUT", nil)
	if err != nil || !one.Data.Active {
		t.Fatalf("GetOperationCodeByCode failed: resp=%+v err=%v", one, err)
	}
	_, err = c.GetOperationCodeByCode(context.Background(), "", nil)
	if err == nil || !strings.Contains(err.Error(), "operation code is required") {
		t.Fatalf("expected operation code validation error, got %v", err)
	}
}
