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

func TestOrdersSection(t *testing.T) {
	salesID := "S2"
	terms := "NET30"
	notes := "Rush"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/orders":
			switch r.Method {
			case http.MethodGet:
				if r.URL.Query().Get("fields") == "" {
					t.Fatalf("expected fields query")
				}
				_ = json.NewEncoder(w).Encode(ListResponse[Order]{
					Data: []Order{{OrderNumber: "SO-100", CustomerCode: "ACME", Status: "Open", OrderTotal: 1250}},
				})
			case http.MethodPost:
				_ = json.NewEncoder(w).Encode(SingleResponse[Order]{
					Data: Order{OrderNumber: "SO-100", CustomerCode: "ACME", Status: "Open"},
				})
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		case "/api/v1/orders/SO-100":
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(SingleResponse[Order]{
					Data: Order{OrderNumber: "SO-100", CustomerCode: "ACME", Status: "Open"},
				})
			case http.MethodPatch:
				_ = json.NewEncoder(w).Encode(SingleResponse[Order]{
					Data: Order{OrderNumber: "SO-100", CustomerCode: "ACME", Status: "Open", SalesID: "S2", TermsCode: "NET30"},
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
	q.Set("fields", "orderNumber,customerCode,status")
	listResp, err := c.GetOrders(context.Background(), q)
	if err != nil {
		t.Fatalf("GetOrders failed: %v", err)
	}
	if len(listResp.Data) != 1 || listResp.Data[0].OrderNumber != "SO-100" {
		t.Fatalf("unexpected order list response: %+v", listResp)
	}

	singleResp, err := c.GetOrderByNumber(context.Background(), "SO-100", nil)
	if err != nil {
		t.Fatalf("GetOrderByNumber failed: %v", err)
	}
	if singleResp.Data.OrderNumber != "SO-100" {
		t.Fatalf("unexpected order response: %+v", singleResp)
	}

	_, err = c.GetOrderByNumber(context.Background(), "   ", nil)
	if err == nil || !strings.Contains(err.Error(), "order number is required") {
		t.Fatalf("expected empty order number error, got %v", err)
	}

	createResp, err := c.CreateOrder(context.Background(), OrderCreate{
		OrderNumber:  "SO-100",
		CustomerCode: "ACME",
		PONumber:     "PO-1",
	})
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}
	if createResp.Data.OrderNumber != "SO-100" {
		t.Fatalf("unexpected created order: %+v", createResp)
	}

	updateResp, err := c.UpdateOrder(context.Background(), "SO-100", OrderUpdate{
		SalesID:         &salesID,
		TermsCode:       &terms,
		NotesToCustomer: &notes,
	})
	if err != nil {
		t.Fatalf("UpdateOrder failed: %v", err)
	}
	if updateResp.StatusCode != http.StatusOK || updateResp.NoContent {
		t.Fatalf("unexpected updated order: %+v", updateResp)
	}

	_, err = c.UpdateOrder(context.Background(), "", OrderUpdate{})
	if err == nil || !strings.Contains(err.Error(), "order number is required") {
		t.Fatalf("expected empty order number update error, got %v", err)
	}
}

func TestOrderLineItemsRoutingsAndReleases(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/order-line-items":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(ListResponse[OrderLineItem]{
				Data: []OrderLineItem{{OrderNumber: "SO-100", ItemNumber: 1, PartNumber: "PN-1", Status: "Open"}},
			})
		case "/api/v1/order-routings":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(ListResponse[OrderRouting]{
				Data: []OrderRouting{{OrderNumber: "SO-100", ItemNumber: 1, StepNumber: 10, OperationCode: "CUT"}},
			})
		case "/api/v1/orders/SO-100/order-line-items":
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[OrderLineItem]{
				Data: OrderLineItem{OrderNumber: "SO-100", ItemNumber: 1, PartNumber: "PN-1", Status: "Open"},
			})
		case "/api/v1/orders/SO-100/order-line-items/1":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			if r.URL.Query().Get("fields") == "" {
				t.Fatalf("expected fields query")
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[OrderLineItem]{
				Data: OrderLineItem{OrderNumber: "SO-100", ItemNumber: 1, PartNumber: "PN-1", Status: "Open"},
			})
		case "/api/v1/orders/SO-100/order-line-items/1/order-routings":
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[OrderRouting]{
				Data: OrderRouting{OrderNumber: "SO-100", ItemNumber: 1, StepNumber: 10, OperationCode: "CUT"},
			})
		case "/api/v1/orders/SO-100/order-line-items/1/order-routings/10":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[OrderRouting]{
				Data: OrderRouting{OrderNumber: "SO-100", ItemNumber: 1, StepNumber: 10, OperationCode: "CUT"},
			})
		case "/api/v1/orders/SO-100/order-line-items/1/releases":
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[Release]{
				Data: Release{OrderNumber: "SO-100", ItemNumber: 1, UniqueID: 55, Quantity: 2},
			})
		case "/api/v1/orders/SO-100/order-line-items/1/releases/55":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[Release]{
				Data: Release{OrderNumber: "SO-100", ItemNumber: 1, UniqueID: 55, Quantity: 2},
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
	q.Set("fields", "orderNumber,itemNumber,partNumber")
	allItems, err := c.GetAllOrderLineItems(context.Background(), q)
	if err != nil || len(allItems.Data) != 1 || allItems.Data[0].PartNumber != "PN-1" {
		t.Fatalf("GetAllOrderLineItems failed: resp=%+v err=%v", allItems, err)
	}
	allRoutings, err := c.GetAllOrderRoutings(context.Background(), nil)
	if err != nil || len(allRoutings.Data) != 1 || allRoutings.Data[0].StepNumber != 10 {
		t.Fatalf("GetAllOrderRoutings failed: resp=%+v err=%v", allRoutings, err)
	}

	liCreate, err := c.CreateOrderLineItem(context.Background(), "SO-100", OrderLineItemCreate{
		PartNumber:      "PN-1",
		PartDescription: "Widget",
		QuantityOrdered: 2,
	})
	if err != nil || liCreate.Data.ItemNumber != 1 {
		t.Fatalf("CreateOrderLineItem failed: resp=%+v err=%v", liCreate, err)
	}

	liGet, err := c.GetOrderLineItem(context.Background(), "SO-100", "1", q)
	if err != nil || liGet.Data.PartNumber != "PN-1" {
		t.Fatalf("GetOrderLineItem failed: resp=%+v err=%v", liGet, err)
	}

	rCreate, err := c.CreateOrderRouting(context.Background(), "SO-100", "1", OrderRoutingCreate{
		StepNumber:    10,
		OperationCode: "CUT",
	})
	if err != nil || rCreate.Data.StepNumber != 10 {
		t.Fatalf("CreateOrderRouting failed: resp=%+v err=%v", rCreate, err)
	}

	rGet, err := c.GetOrderRouting(context.Background(), "SO-100", "1", "10", nil)
	if err != nil || rGet.Data.OperationCode != "CUT" {
		t.Fatalf("GetOrderRouting failed: resp=%+v err=%v", rGet, err)
	}

	relCreate, err := c.CreateOrderRelease(context.Background(), "SO-100", "1", ReleaseCreate{Quantity: 2})
	if err != nil || relCreate.Data.UniqueID != 55 {
		t.Fatalf("CreateOrderRelease failed: resp=%+v err=%v", relCreate, err)
	}

	relGet, err := c.GetOrderRelease(context.Background(), "SO-100", "1", "55", nil)
	if err != nil || relGet.Data.UniqueID != 55 {
		t.Fatalf("GetOrderRelease failed: resp=%+v err=%v", relGet, err)
	}

	_, err = c.CreateOrderLineItem(context.Background(), " ", OrderLineItemCreate{})
	if err == nil || !strings.Contains(err.Error(), "order number is required") {
		t.Fatalf("expected order number validation error, got %v", err)
	}

	_, err = c.GetOrderLineItem(context.Background(), "SO-100", "", nil)
	if err == nil || !strings.Contains(err.Error(), "item number is required") {
		t.Fatalf("expected item number validation error, got %v", err)
	}

	_, err = c.GetOrderLineItem(context.Background(), "", "1", nil)
	if err == nil || !strings.Contains(err.Error(), "order number is required") {
		t.Fatalf("expected order number validation error, got %v", err)
	}

	_, err = c.CreateOrderRouting(context.Background(), "", "1", OrderRoutingCreate{})
	if err == nil || !strings.Contains(err.Error(), "order number is required") {
		t.Fatalf("expected order number validation error, got %v", err)
	}

	_, err = c.CreateOrderRouting(context.Background(), "SO-100", " ", OrderRoutingCreate{})
	if err == nil || !strings.Contains(err.Error(), "item number is required") {
		t.Fatalf("expected item number validation error, got %v", err)
	}

	_, err = c.GetOrderRouting(context.Background(), "", "1", "10", nil)
	if err == nil || !strings.Contains(err.Error(), "order number is required") {
		t.Fatalf("expected order number validation error, got %v", err)
	}

	_, err = c.GetOrderRouting(context.Background(), "SO-100", "", "10", nil)
	if err == nil || !strings.Contains(err.Error(), "item number is required") {
		t.Fatalf("expected item number validation error, got %v", err)
	}

	_, err = c.GetOrderRouting(context.Background(), "SO-100", "1", "", nil)
	if err == nil || !strings.Contains(err.Error(), "step number is required") {
		t.Fatalf("expected step number validation error, got %v", err)
	}

	_, err = c.CreateOrderRelease(context.Background(), "", "1", ReleaseCreate{})
	if err == nil || !strings.Contains(err.Error(), "order number is required") {
		t.Fatalf("expected order number validation error, got %v", err)
	}

	_, err = c.CreateOrderRelease(context.Background(), "SO-100", "", ReleaseCreate{})
	if err == nil || !strings.Contains(err.Error(), "item number is required") {
		t.Fatalf("expected item number validation error, got %v", err)
	}

	_, err = c.GetOrderRelease(context.Background(), "", "1", "55", nil)
	if err == nil || !strings.Contains(err.Error(), "order number is required") {
		t.Fatalf("expected order number validation error, got %v", err)
	}

	_, err = c.GetOrderRelease(context.Background(), "SO-100", "", "55", nil)
	if err == nil || !strings.Contains(err.Error(), "item number is required") {
		t.Fatalf("expected item number validation error, got %v", err)
	}

	_, err = c.GetOrderRelease(context.Background(), "SO-100", "1", "", nil)
	if err == nil || !strings.Contains(err.Error(), "unique id is required") {
		t.Fatalf("expected unique id validation error, got %v", err)
	}
}
