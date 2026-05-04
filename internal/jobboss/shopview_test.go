package jobboss

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestShopViewGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/shopview/authenticate-user":
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[ShopViewAuthenticateResponse]{Data: ShopViewAuthenticateResponse{Authenticated: true, Token: "sv-token"}})
		case "/api/v1/shopview/filters":
			_ = json.NewEncoder(w).Encode(ListResponse[ShopViewFilter]{Data: []ShopViewFilter{{Name: "workCenter", Type: "select", Options: []string{"WC1"}}}})
		case "/api/v1/shopview/get-jobs":
			_ = json.NewEncoder(w).Encode(ListResponse[ShopViewJob]{Data: []ShopViewJob{{JobNumber: "J1", Status: "IN_PROGRESS"}}})
		case "/api/v1/shopview/grid-option":
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[ShopViewGridOption]{Data: ShopViewGridOption{OptionName: "columns", OptionValue: "default"}})
		case "/api/v1/shopview/grid-options":
			_ = json.NewEncoder(w).Encode(ListResponse[ShopViewGridOption]{Data: []ShopViewGridOption{{OptionName: "columns", OptionValue: "default"}}})
		case "/api/v1/shopview/kpi-definitions":
			_ = json.NewEncoder(w).Encode(ListResponse[ShopViewKPIDefinition]{Data: []ShopViewKPIDefinition{{Key: "jobs-closed", Title: "Jobs Closed"}}})
		case "/api/v1/shopview/kpi/jobs-closed":
			_ = json.NewEncoder(w).Encode(SingleResponse[ShopViewKPIValue]{Data: ShopViewKPIValue{Value: 5}})
		case "/api/v1/shopview/kpi/jobs-in-progress":
			_ = json.NewEncoder(w).Encode(SingleResponse[ShopViewKPIValue]{Data: ShopViewKPIValue{Value: 12}})
		case "/api/v1/shopview/kpi/jobs-on-hold":
			_ = json.NewEncoder(w).Encode(SingleResponse[ShopViewKPIValue]{Data: ShopViewKPIValue{Value: 2}})
		case "/api/v1/shopview/kpi/jobs-past-due":
			_ = json.NewEncoder(w).Encode(SingleResponse[ShopViewKPIValue]{Data: ShopViewKPIValue{Value: 3}})
		case "/api/v1/shopview/reset-grid-options":
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[ShopViewResetResponse]{Data: ShopViewResetResponse{Success: true, Message: "reset"}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := NewClient(Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})
	q := url.Values{}
	q.Set("fields", "name,type")

	auth, err := c.AuthenticateShopView(context.Background(), ShopViewAuthenticateRequest{Username: "u", Password: "p"})
	if err != nil || !auth.Data.Authenticated {
		t.Fatalf("AuthenticateShopView failed: resp=%+v err=%v", auth, err)
	}
	filters, err := c.GetShopViewFilters(context.Background(), q)
	if err != nil || len(filters.Data) != 1 || filters.Data[0].Name != "workCenter" {
		t.Fatalf("GetShopViewFilters failed: resp=%+v err=%v", filters, err)
	}
	jobs, err := c.GetShopViewJobs(context.Background(), nil)
	if err != nil || len(jobs.Data) != 1 || jobs.Data[0].JobNumber != "J1" {
		t.Fatalf("GetShopViewJobs failed: resp=%+v err=%v", jobs, err)
	}
	createdOption, err := c.CreateShopViewGridOption(context.Background(), ShopViewGridOption{OptionName: "columns"})
	if err != nil || createdOption.Data.OptionName != "columns" {
		t.Fatalf("CreateShopViewGridOption failed: resp=%+v err=%v", createdOption, err)
	}
	options, err := c.GetShopViewGridOptions(context.Background(), nil)
	if err != nil || len(options.Data) != 1 || options.Data[0].OptionName != "columns" {
		t.Fatalf("GetShopViewGridOptions failed: resp=%+v err=%v", options, err)
	}
	kpiDefs, err := c.GetShopViewKPIDefinitions(context.Background(), nil)
	if err != nil || len(kpiDefs.Data) != 1 || kpiDefs.Data[0].Key != "jobs-closed" {
		t.Fatalf("GetShopViewKPIDefinitions failed: resp=%+v err=%v", kpiDefs, err)
	}
	closed, err := c.GetShopViewKPIJobsClosed(context.Background(), nil)
	if err != nil || closed.Data.Value != 5 {
		t.Fatalf("GetShopViewKPIJobsClosed failed: resp=%+v err=%v", closed, err)
	}
	inProgress, err := c.GetShopViewKPIJobsInProgress(context.Background(), nil)
	if err != nil || inProgress.Data.Value != 12 {
		t.Fatalf("GetShopViewKPIJobsInProgress failed: resp=%+v err=%v", inProgress, err)
	}
	onHold, err := c.GetShopViewKPIJobsOnHold(context.Background(), nil)
	if err != nil || onHold.Data.Value != 2 {
		t.Fatalf("GetShopViewKPIJobsOnHold failed: resp=%+v err=%v", onHold, err)
	}
	pastDue, err := c.GetShopViewKPIJobsPastDue(context.Background(), nil)
	if err != nil || pastDue.Data.Value != 3 {
		t.Fatalf("GetShopViewKPIJobsPastDue failed: resp=%+v err=%v", pastDue, err)
	}
	reset, err := c.ResetShopViewGridOptions(context.Background())
	if err != nil || !reset.Data.Success {
		t.Fatalf("ResetShopViewGridOptions failed: resp=%+v err=%v", reset, err)
	}
}
