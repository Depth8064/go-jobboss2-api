package jobboss

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCompanyGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/company":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(ListResponse[Company]{
				Data: []Company{{CompanyCode: "ABC", CompanyDescription: "Acme", CurrencyCode: "USD"}},
			})
		case "/api/v1/company-calendars":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(ListResponse[CompanyCalendar]{
				Data: []CompanyCalendar{{UniqueID: 1, CalendarType: "WORK", Date: "2026-03-01", HoursAvailable: 8}},
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

	cq := url.Values{}
	cq.Set("fields", "companyCode,companyDescription,currencyCode")
	companyResp, err := c.GetCompany(context.Background(), cq)
	if err != nil {
		t.Fatalf("GetCompany failed: %v", err)
	}
	if len(companyResp.Data) != 1 || companyResp.Data[0].CompanyCode != "ABC" {
		t.Fatalf("unexpected company response: %+v", companyResp)
	}

	calQ := url.Values{}
	calQ.Set("fields", "calendarType,date,hoursAvailable")
	calendarResp, err := c.GetCompanyCalendars(context.Background(), calQ)
	if err != nil {
		t.Fatalf("GetCompanyCalendars failed: %v", err)
	}
	if len(calendarResp.Data) != 1 || calendarResp.Data[0].UniqueID != 1 {
		t.Fatalf("unexpected company calendar response: %+v", calendarResp)
	}
}
