package jobboss

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestLookupGroups(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/routings":
			_ = json.NewEncoder(w).Encode(ListResponse[Routing]{Data: []Routing{{OperationCode: "CUT", Description: "Cutting"}}})
		case "/api/v1/tooling-maintenance":
			_ = json.NewEncoder(w).Encode(ListResponse[ToolingMaintenance]{Data: []ToolingMaintenance{{ToolNumber: "T-1", Status: "DUE"}}})
		case "/api/v1/user-labels":
			_ = json.NewEncoder(w).Encode(ListResponse[UserLabel]{Data: []UserLabel{{LabelName: "LBL1", DisplayText: "Label One"}}})
		case "/api/v1/user-transactions":
			_ = json.NewEncoder(w).Encode(ListResponse[UserTransaction]{Data: []UserTransaction{{UserCode: "USR1", TransactionType: "LOGIN"}}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := NewClient(Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})
	q := url.Values{}
	q.Set("fields", "description")

	routings, err := c.GetRoutings(context.Background(), q)
	if err != nil || len(routings.Data) != 1 || routings.Data[0].OperationCode != "CUT" {
		t.Fatalf("GetRoutings failed: resp=%+v err=%v", routings, err)
	}
	tooling, err := c.GetToolingMaintenance(context.Background(), nil)
	if err != nil || len(tooling.Data) != 1 || tooling.Data[0].ToolNumber != "T-1" {
		t.Fatalf("GetToolingMaintenance failed: resp=%+v err=%v", tooling, err)
	}
	labels, err := c.GetUserLabels(context.Background(), nil)
	if err != nil || len(labels.Data) != 1 || labels.Data[0].LabelName != "LBL1" {
		t.Fatalf("GetUserLabels failed: resp=%+v err=%v", labels, err)
	}
	transactions, err := c.GetUserTransactions(context.Background(), nil)
	if err != nil || len(transactions.Data) != 1 || transactions.Data[0].UserCode != "USR1" {
		t.Fatalf("GetUserTransactions failed: resp=%+v err=%v", transactions, err)
	}
}
