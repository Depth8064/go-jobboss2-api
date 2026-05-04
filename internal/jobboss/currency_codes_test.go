package jobboss

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCurrencyCodesGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/currency-codes":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(ListResponse[CurrencyCode]{
				Data: []CurrencyCode{{CurrencyCode: "USD", CurrencyDescription: "US Dollar", ExchangeRate: 1}},
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
	q.Set("fields", "currencyCode,currencyDescription,exchangeRate")
	resp, err := c.GetCurrencyCodes(context.Background(), q)
	if err != nil {
		t.Fatalf("GetCurrencyCodes failed: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].CurrencyCode != "USD" {
		t.Fatalf("unexpected currency code response: %+v", resp)
	}
}
