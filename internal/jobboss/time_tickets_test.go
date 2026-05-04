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

func TestTimeTicketsGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/time-tickets":
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(ListResponse[TimeTicket]{Data: []TimeTicket{{TicketDate: "2024-01-01", EmployeeCode: "E1"}}})
			case http.MethodPost:
				_ = json.NewEncoder(w).Encode(SingleResponse[TimeTicket]{Data: TimeTicket{TicketDate: "2024-01-02", EmployeeCode: "E2"}})
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		case "/api/v1/time-tickets/2024-01-01/employees/E1":
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(SingleResponse[TimeTicket]{Data: TimeTicket{TicketDate: "2024-01-01", EmployeeCode: "E1", Status: "OPEN"}})
			case http.MethodPatch:
				_ = json.NewEncoder(w).Encode(SingleResponse[TimeTicket]{Data: TimeTicket{TicketDate: "2024-01-01", EmployeeCode: "E1", Status: "CLOSED"}})
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		case "/api/v1/time-ticket-details":
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(ListResponse[TimeTicketDetail]{Data: []TimeTicketDetail{{TimeTicketGUID: "GUID-1", OrderNumber: "SO-1"}}})
			case http.MethodPost:
				_ = json.NewEncoder(w).Encode(SingleResponse[TimeTicketDetail]{Data: TimeTicketDetail{TimeTicketGUID: "GUID-2", OrderNumber: "SO-2"}})
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		case "/api/v1/time-ticket-details/GUID-1":
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(SingleResponse[TimeTicketDetail]{Data: TimeTicketDetail{TimeTicketGUID: "GUID-1", Status: "OPEN"}})
			case http.MethodPatch:
				_ = json.NewEncoder(w).Encode(SingleResponse[TimeTicketDetail]{Data: TimeTicketDetail{TimeTicketGUID: "GUID-1", Status: "CLOSED"}})
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	c := NewClient(Config{BaseURL: ts.URL, TokenURL: ts.URL + "/oauth2/api-user/token", ClientID: "id", ClientSecret: "secret", HTTPClient: ts.Client(), MaxAuthRetryOn401: 0})
	q := url.Values{}
	q.Set("fields", "ticketDate,employeeCode")

	tickets, err := c.GetTimeTickets(context.Background(), q)
	if err != nil || len(tickets.Data) != 1 || tickets.Data[0].EmployeeCode != "E1" {
		t.Fatalf("GetTimeTickets failed: resp=%+v err=%v", tickets, err)
	}
	createdTicket, err := c.CreateTimeTicket(context.Background(), TimeTicket{TicketDate: "2024-01-02", EmployeeCode: "E2"})
	if err != nil || createdTicket.Data.EmployeeCode != "E2" {
		t.Fatalf("CreateTimeTicket failed: resp=%+v err=%v", createdTicket, err)
	}
	oneTicket, err := c.GetTimeTicketByDateAndEmployee(context.Background(), "2024-01-01", "E1", nil)
	if err != nil || oneTicket.Data.Status != "OPEN" {
		t.Fatalf("GetTimeTicketByDateAndEmployee failed: resp=%+v err=%v", oneTicket, err)
	}
	updatedTicket, err := c.UpdateTimeTicketByDateAndEmployee(context.Background(), "2024-01-01", "E1", TimeTicket{Status: "CLOSED"})
	if err != nil || updatedTicket.StatusCode != http.StatusOK || updatedTicket.NoContent {
		t.Fatalf("UpdateTimeTicketByDateAndEmployee failed: resp=%+v err=%v", updatedTicket, err)
	}
	_, err = c.GetTimeTicketByDateAndEmployee(context.Background(), "", "E1", nil)
	if err == nil || !strings.Contains(err.Error(), "ticket date is required") {
		t.Fatalf("expected ticket date validation error, got %v", err)
	}
	_, err = c.GetTimeTicketByDateAndEmployee(context.Background(), "2024-01-01", "", nil)
	if err == nil || !strings.Contains(err.Error(), "employee code is required") {
		t.Fatalf("expected employee code validation error, got %v", err)
	}
	_, err = c.UpdateTimeTicketByDateAndEmployee(context.Background(), "", "E1", TimeTicket{})
	if err == nil || !strings.Contains(err.Error(), "ticket date is required") {
		t.Fatalf("expected ticket date validation error, got %v", err)
	}
	_, err = c.UpdateTimeTicketByDateAndEmployee(context.Background(), "2024-01-01", "", TimeTicket{})
	if err == nil || !strings.Contains(err.Error(), "employee code is required") {
		t.Fatalf("expected employee code validation error, got %v", err)
	}

	details, err := c.GetTimeTicketDetails(context.Background(), nil)
	if err != nil || len(details.Data) != 1 || details.Data[0].TimeTicketGUID != "GUID-1" {
		t.Fatalf("GetTimeTicketDetails failed: resp=%+v err=%v", details, err)
	}
	createdDetail, err := c.CreateTimeTicketDetail(context.Background(), TimeTicketDetail{TimeTicketGUID: "GUID-2"})
	if err != nil || createdDetail.Data.TimeTicketGUID != "GUID-2" {
		t.Fatalf("CreateTimeTicketDetail failed: resp=%+v err=%v", createdDetail, err)
	}
	oneDetail, err := c.GetTimeTicketDetailByGUID(context.Background(), "GUID-1", nil)
	if err != nil || oneDetail.Data.Status != "OPEN" {
		t.Fatalf("GetTimeTicketDetailByGUID failed: resp=%+v err=%v", oneDetail, err)
	}
	updatedDetail, err := c.UpdateTimeTicketDetailByGUID(context.Background(), "GUID-1", TimeTicketDetail{Status: "CLOSED"})
	if err != nil || updatedDetail.StatusCode != http.StatusOK || updatedDetail.NoContent {
		t.Fatalf("UpdateTimeTicketDetailByGUID failed: resp=%+v err=%v", updatedDetail, err)
	}
	_, err = c.GetTimeTicketDetailByGUID(context.Background(), "", nil)
	if err == nil || !strings.Contains(err.Error(), "time ticket guid is required") {
		t.Fatalf("expected time ticket guid validation error, got %v", err)
	}
	_, err = c.UpdateTimeTicketDetailByGUID(context.Background(), "", TimeTicketDetail{})
	if err == nil || !strings.Contains(err.Error(), "time ticket guid is required") {
		t.Fatalf("expected time ticket guid validation error, got %v", err)
	}
}
