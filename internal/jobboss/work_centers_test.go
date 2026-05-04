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

func TestWorkCentersGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/work-centers":
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(ListResponse[WorkCenter]{Data: []WorkCenter{{WorkCenter: "WC1", Description: "Laser"}}})
			case http.MethodPost:
				_ = json.NewEncoder(w).Encode(SingleResponse[WorkCenter]{Data: WorkCenter{WorkCenter: "WC2", Description: "Form"}})
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		case "/api/v1/work-centers/WC1":
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(SingleResponse[WorkCenter]{Data: WorkCenter{WorkCenter: "WC1", Active: true}})
			case http.MethodPatch:
				_ = json.NewEncoder(w).Encode(SingleResponse[WorkCenter]{Data: WorkCenter{WorkCenter: "WC1", Active: false}})
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		case "/api/v1/work-center-maintenance":
			_ = json.NewEncoder(w).Encode(ListResponse[WorkCenterMaintenance]{Data: []WorkCenterMaintenance{{WorkCenter: "WC1", Status: "DUE"}}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := NewClient(Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})
	q := url.Values{}
	q.Set("fields", "workCenter,description")

	workCenters, err := c.GetWorkCenters(context.Background(), q)
	if err != nil || len(workCenters.Data) != 1 || workCenters.Data[0].WorkCenter != "WC1" {
		t.Fatalf("GetWorkCenters failed: resp=%+v err=%v", workCenters, err)
	}
	workCenter, err := c.GetWorkCenterByCode(context.Background(), "WC1", nil)
	if err != nil || !workCenter.Data.Active {
		t.Fatalf("GetWorkCenterByCode failed: resp=%+v err=%v", workCenter, err)
	}
	_, err = c.GetWorkCenterByCode(context.Background(), "", nil)
	if err == nil || !strings.Contains(err.Error(), "work center is required") {
		t.Fatalf("expected work center validation error, got %v", err)
	}

	created, err := c.CreateWorkCenter(context.Background(), WorkCenter{WorkCenter: "WC2"})
	if err != nil || created.Data.WorkCenter != "WC2" {
		t.Fatalf("CreateWorkCenter failed: resp=%+v err=%v", created, err)
	}
	updated, err := c.UpdateWorkCenter(context.Background(), "WC1", WorkCenter{Active: false})
	if err != nil || updated.Data.Active {
		t.Fatalf("UpdateWorkCenter failed: resp=%+v err=%v", updated, err)
	}
	_, err = c.UpdateWorkCenter(context.Background(), "", WorkCenter{})
	if err == nil || !strings.Contains(err.Error(), "work center is required") {
		t.Fatalf("expected work center validation error, got %v", err)
	}

	maintenance, err := c.GetWorkCenterMaintenance(context.Background(), nil)
	if err != nil || len(maintenance.Data) != 1 || maintenance.Data[0].Status != "DUE" {
		t.Fatalf("GetWorkCenterMaintenance failed: resp=%+v err=%v", maintenance, err)
	}
}
