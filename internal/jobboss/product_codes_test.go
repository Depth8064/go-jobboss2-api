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

func TestProductCodesGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/product-codes":
			_ = json.NewEncoder(w).Encode(ListResponse[ProductCode]{Data: []ProductCode{{ProductCode: "FG", Description: "Finished Good"}}})
		case "/api/v1/product-codes/FG":
			_ = json.NewEncoder(w).Encode(SingleResponse[ProductCode]{Data: ProductCode{ProductCode: "FG", Description: "Finished Good", Active: true}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := NewClient(Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})
	q := url.Values{}
	q.Set("fields", "productCode,description")
	list, err := c.GetProductCodes(context.Background(), q)
	if err != nil || len(list.Data) != 1 || list.Data[0].ProductCode != "FG" {
		t.Fatalf("GetProductCodes failed: resp=%+v err=%v", list, err)
	}
	one, err := c.GetProductCodeByCode(context.Background(), "FG", nil)
	if err != nil || !one.Data.Active {
		t.Fatalf("GetProductCodeByCode failed: resp=%+v err=%v", one, err)
	}
	_, err = c.GetProductCodeByCode(context.Background(), "", nil)
	if err == nil || !strings.Contains(err.Error(), "product code is required") {
		t.Fatalf("expected product code validation error, got %v", err)
	}
}
