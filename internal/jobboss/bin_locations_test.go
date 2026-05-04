package jobboss

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestBinLocationsSection(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/bin-locations":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			if r.URL.Query().Get("filters[partNumber]") != "" {
				// no-op placeholder for custom filter patterns
			}
			_ = json.NewEncoder(w).Encode(ListResponse[BinLocation]{
				Data: []BinLocation{{BinLocation: "A-01", PartNumber: "PN-1", QuantityOnHand: 5}},
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
	q.Set("fields", "binLocation,partNumber,quantityOnHand")
	resp, err := c.GetBinLocations(context.Background(), q)
	if err != nil {
		t.Fatalf("GetBinLocations failed: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].BinLocation != "A-01" {
		t.Fatalf("unexpected bin locations response: %+v", resp)
	}
}
