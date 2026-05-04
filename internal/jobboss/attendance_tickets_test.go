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

func TestAttendanceTicketGroups(t *testing.T) {
	comment := "fixed"
	isOT := true

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/api-user/token":
			_ = json.NewEncoder(w).Encode(TokenResponse{AccessToken: "tok", ExpiresIn: 3600})
		case "/api/v1/attendance-ticket-details":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(ListResponse[AttendanceTicketDetail]{
				Data: []AttendanceTicketDetail{{UniqueID: 1, EmployeeCode: 101, EmployeeName: "John", TicketDate: "2026-03-01"}},
			})
		case "/api/v1/attendance-ticket-details/1":
			if r.Method != http.MethodPatch {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[AttendanceTicketDetail]{
				Data: AttendanceTicketDetail{UniqueID: 1, Comments: "fixed", IsOvertime: true},
			})
		case "/api/v1/attendance-tickets/2026-03-01/employees/101/attendance-ticket-details":
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[AttendanceTicketDetail]{
				Data: AttendanceTicketDetail{UniqueID: 2, EmployeeCode: 101, TicketDate: "2026-03-01"},
			})
		case "/api/v1/attendance-tickets":
			switch r.Method {
			case http.MethodGet:
				if r.URL.Query().Get("fields") == "" {
					t.Fatalf("expected fields query")
				}
				_ = json.NewEncoder(w).Encode(ListResponse[AttendanceTicket]{
					Data: []AttendanceTicket{{UniqueID: 3, EmployeeCode: 101, TicketDate: "2026-03-01"}},
				})
			case http.MethodPost:
				_ = json.NewEncoder(w).Encode(SingleResponse[AttendanceTicket]{
					Data: AttendanceTicket{UniqueID: 4, EmployeeCode: 101, TicketDate: "2026-03-02"},
				})
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		case "/api/v1/attendance-tickets/2026-03-02/employees/101":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewEncoder(w).Encode(SingleResponse[AttendanceTicket]{
				Data: AttendanceTicket{UniqueID: 4, EmployeeCode: 101, TicketDate: "2026-03-02"},
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

	details, err := c.GetAttendanceTicketDetails(context.Background(), nil)
	if err != nil || len(details.Data) != 1 || details.Data[0].UniqueID != 1 {
		t.Fatalf("GetAttendanceTicketDetails failed: resp=%+v err=%v", details, err)
	}

	upd, err := c.UpdateAttendanceTicketDetail(context.Background(), "1", AttendanceTicketDetailUpdate{Comments: &comment, IsOvertime: &isOT})
	if err != nil || upd.Data.Comments != "fixed" || !upd.Data.IsOvertime {
		t.Fatalf("UpdateAttendanceTicketDetail failed: resp=%+v err=%v", upd, err)
	}

	createdDetail, err := c.CreateAttendanceTicketDetail(context.Background(), "2026-03-01", "101", AttendanceTicketDetailCreate{Comments: "new"})
	if err != nil || createdDetail.Data.UniqueID != 2 {
		t.Fatalf("CreateAttendanceTicketDetail failed: resp=%+v err=%v", createdDetail, err)
	}

	q := url.Values{}
	q.Set("fields", "employeeCode,ticketDate")
	tickets, err := c.GetAttendanceTickets(context.Background(), q)
	if err != nil || len(tickets.Data) != 1 || tickets.Data[0].UniqueID != 3 {
		t.Fatalf("GetAttendanceTickets failed: resp=%+v err=%v", tickets, err)
	}

	createdTicket, err := c.CreateAttendanceTicket(context.Background(), AttendanceTicketCreate{EmployeeCode: 101, TicketDate: "2026-03-02"})
	if err != nil || createdTicket.Data.UniqueID != 4 {
		t.Fatalf("CreateAttendanceTicket failed: resp=%+v err=%v", createdTicket, err)
	}

	ticket, err := c.GetAttendanceTicket(context.Background(), "2026-03-02", "101", nil)
	if err != nil || ticket.Data.UniqueID != 4 {
		t.Fatalf("GetAttendanceTicket failed: resp=%+v err=%v", ticket, err)
	}

	_, err = c.UpdateAttendanceTicketDetail(context.Background(), "", AttendanceTicketDetailUpdate{})
	if err == nil || !strings.Contains(err.Error(), "id is required") {
		t.Fatalf("expected id validation error, got %v", err)
	}

	_, err = c.CreateAttendanceTicketDetail(context.Background(), "", "101", AttendanceTicketDetailCreate{})
	if err == nil || !strings.Contains(err.Error(), "ticket date is required") {
		t.Fatalf("expected ticket date validation error, got %v", err)
	}

	_, err = c.CreateAttendanceTicketDetail(context.Background(), "2026-03-01", "", AttendanceTicketDetailCreate{})
	if err == nil || !strings.Contains(err.Error(), "employee code is required") {
		t.Fatalf("expected employee code validation error, got %v", err)
	}

	_, err = c.GetAttendanceTicket(context.Background(), "", "101", nil)
	if err == nil || !strings.Contains(err.Error(), "ticket date is required") {
		t.Fatalf("expected ticket date validation error, got %v", err)
	}

	_, err = c.GetAttendanceTicket(context.Background(), "2026-03-02", "", nil)
	if err == nil || !strings.Contains(err.Error(), "employee code is required") {
		t.Fatalf("expected employee code validation error, got %v", err)
	}
}
