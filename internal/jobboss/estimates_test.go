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

func TestEstimatesGroup(t *testing.T) {
	desc := "Updated desc"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/estimates":
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(ListResponse[Estimate]{
					Data: []Estimate{{PartNumber: "PN-100", Description: "Widget", ProductCode: "JOBSHOP"}},
				})
			case http.MethodPost:
				_ = json.NewEncoder(w).Encode(SingleResponse[Estimate]{
					Data: Estimate{PartNumber: "PN-100", Description: "Widget"},
				})
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		case "/api/v1/estimates/PN-100":
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(SingleResponse[Estimate]{
					Data: Estimate{PartNumber: "PN-100", Description: "Widget", Revision: "A"},
				})
			case http.MethodPatch:
				_ = json.NewEncoder(w).Encode(SingleResponse[Estimate]{
					Data: Estimate{PartNumber: "PN-100", Description: "Updated desc", Revision: "B"},
				})
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
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
	q.Set("fields", "partNumber,description,productCode")
	list, err := c.GetEstimates(context.Background(), q)
	if err != nil || len(list.Data) != 1 || list.Data[0].PartNumber != "PN-100" {
		t.Fatalf("GetEstimates failed: resp=%+v err=%v", list, err)
	}

	one, err := c.GetEstimateByPartNumber(context.Background(), "PN-100", nil)
	if err != nil || one.Data.Revision != "A" {
		t.Fatalf("GetEstimateByPartNumber failed: resp=%+v err=%v", one, err)
	}

	created, err := c.CreateEstimate(context.Background(), EstimateCreate{PartNumber: "PN-100", Description: "Widget"})
	if err != nil || created.Data.PartNumber != "PN-100" {
		t.Fatalf("CreateEstimate failed: resp=%+v err=%v", created, err)
	}

	updated, err := c.UpdateEstimate(context.Background(), "PN-100", EstimateUpdate{Description: &desc})
	if err != nil || updated.Data.Description != "Updated desc" {
		t.Fatalf("UpdateEstimate failed: resp=%+v err=%v", updated, err)
	}

	_, err = c.GetEstimateByPartNumber(context.Background(), "", nil)
	if err == nil || !strings.Contains(err.Error(), "part number is required") {
		t.Fatalf("expected part number validation error, got %v", err)
	}

	_, err = c.UpdateEstimate(context.Background(), "", EstimateUpdate{})
	if err == nil || !strings.Contains(err.Error(), "part number is required") {
		t.Fatalf("expected part number validation error, got %v", err)
	}
}
