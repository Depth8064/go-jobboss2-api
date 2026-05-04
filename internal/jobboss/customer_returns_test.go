package jobboss

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCustomerReturnsGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/customer-returns":
			_ = json.NewEncoder(w).Encode(ListResponse[CustomerReturn]{Data: []CustomerReturn{{CustomerReturnNumber: "CR-1", Status: "OPEN"}}})
		case "/api/v1/customer-return-line-items":
			_ = json.NewEncoder(w).Encode(ListResponse[CustomerReturnLineItem]{Data: []CustomerReturnLineItem{{CustomerReturnNumber: "CR-1", ItemNumber: 1}}})
		case "/api/v1/customer-return-releases":
			_ = json.NewEncoder(w).Encode(ListResponse[CustomerReturnRelease]{Data: []CustomerReturnRelease{{CustomerReturnNumber: "CR-1", ReleaseNumber: 1}}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := NewClient(Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})
	q := url.Values{}
	q.Set("fields", "customerReturnNumber,status")

	headers, err := c.GetCustomerReturns(context.Background(), q)
	if err != nil || len(headers.Data) != 1 || headers.Data[0].CustomerReturnNumber != "CR-1" {
		t.Fatalf("GetCustomerReturns failed: resp=%+v err=%v", headers, err)
	}
	items, err := c.GetCustomerReturnLineItems(context.Background(), nil)
	if err != nil || len(items.Data) != 1 || items.Data[0].ItemNumber != 1 {
		t.Fatalf("GetCustomerReturnLineItems failed: resp=%+v err=%v", items, err)
	}
	releases, err := c.GetCustomerReturnReleases(context.Background(), nil)
	if err != nil || len(releases.Data) != 1 || releases.Data[0].ReleaseNumber != 1 {
		t.Fatalf("GetCustomerReturnReleases failed: resp=%+v err=%v", releases, err)
	}
}
