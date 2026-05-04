package jobboss

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// AttendanceTicketDetail models an attendance ticket detail payload.
type AttendanceTicketDetail struct {
	UniqueID             int    `json:"uniqueID,omitempty"`
	EmployeeCode         int    `json:"employeeCode,omitempty"`
	EmployeeName         string `json:"employeeName,omitempty"`
	TicketDate           string `json:"ticketDate,omitempty"`
	AttendanceCode       int    `json:"attendanceCode,omitempty"`
	ActualClockInTime    string `json:"actualClockInTime,omitempty"`
	ActualClockOutTime   string `json:"actualClockOutTime,omitempty"`
	AdjustedClockInTime  string `json:"adjustedClockInTime,omitempty"`
	AdjustedClockOutTime string `json:"adjustedClockOutTime,omitempty"`
	ClockOutDate         string `json:"clockOutDate,omitempty"`
	Comments             string `json:"comments,omitempty"`
	GLCode               string `json:"GLCode,omitempty"`
	IsHoliday            bool   `json:"isHoliday,omitempty"`
	IsOvertime           bool   `json:"isOvertime,omitempty"`
	PayRateCode          int    `json:"payRateCode,omitempty"`
	Shift                int    `json:"shift,omitempty"`
	LastModDate          string `json:"lastModDate,omitempty"`
}

// AttendanceTicketDetailCreate is a payload for creating attendance ticket details.
type AttendanceTicketDetailCreate struct {
	ActualClockInTime    string `json:"actualClockInTime,omitempty"`
	ActualClockOutTime   string `json:"actualClockOutTime,omitempty"`
	AdjustedClockInTime  string `json:"adjustedClockInTime,omitempty"`
	AdjustedClockOutTime string `json:"adjustedClockOutTime,omitempty"`
	AttendanceCode       int    `json:"attendanceCode,omitempty"`
	ClockOutDate         string `json:"clockOutDate,omitempty"`
	Comments             string `json:"comments,omitempty"`
	GLCode               string `json:"GLCode,omitempty"`
	IsHoliday            bool   `json:"isHoliday,omitempty"`
	IsOvertime           bool   `json:"isOvertime,omitempty"`
	PayRateCode          int    `json:"payRateCode,omitempty"`
	Shift                int    `json:"shift,omitempty"`
}

// AttendanceTicketDetailUpdate is a payload for patching attendance ticket details.
type AttendanceTicketDetailUpdate struct {
	ActualClockInTime    *string `json:"actualClockInTime,omitempty"`
	ActualClockOutTime   *string `json:"actualClockOutTime,omitempty"`
	AdjustedClockInTime  *string `json:"adjustedClockInTime,omitempty"`
	AdjustedClockOutTime *string `json:"adjustedClockOutTime,omitempty"`
	AttendanceCode       *int    `json:"attendanceCode,omitempty"`
	ClockOutDate         *string `json:"clockOutDate,omitempty"`
	Comments             *string `json:"comments,omitempty"`
	GLCode               *string `json:"GLCode,omitempty"`
	IsHoliday            *bool   `json:"isHoliday,omitempty"`
	IsOvertime           *bool   `json:"isOvertime,omitempty"`
	PayRateCode          *int    `json:"payRateCode,omitempty"`
	Shift                *int    `json:"shift,omitempty"`
}

// AttendanceTicket models an attendance ticket payload.
type AttendanceTicket struct {
	UniqueID     int    `json:"uniqueID,omitempty"`
	EmployeeCode int    `json:"employeeCode,omitempty"`
	EmployeeName string `json:"employeeName,omitempty"`
	TicketDate   string `json:"ticketDate,omitempty"`
	EnteredBy    string `json:"enteredBy,omitempty"`
	EnteredDate  string `json:"enteredDate,omitempty"`
	LastModDate  string `json:"lastModDate,omitempty"`
	LastModUser  string `json:"lastModUser,omitempty"`
	SearchDate   int    `json:"searchDate,omitempty"`
	IsExported   bool   `json:"isExported,omitempty"`
}

// AttendanceTicketCreate is a payload for creating attendance tickets.
type AttendanceTicketCreate struct {
	EmployeeCode int    `json:"employeeCode,omitempty"`
	TicketDate   string `json:"ticketDate,omitempty"`
}

// GetAttendanceTicketDetails fetches attendance ticket details.
func (c *Client) GetAttendanceTicketDetails(ctx context.Context, query url.Values) (ListResponse[AttendanceTicketDetail], error) {
	return GetJSON[ListResponse[AttendanceTicketDetail]](ctx, c, "/api/v1/attendance-ticket-details", query)
}

// UpdateAttendanceTicketDetail updates an attendance ticket detail by ID.
func (c *Client) UpdateAttendanceTicketDetail(ctx context.Context, id string, payload AttendanceTicketDetailUpdate) (SingleResponse[AttendanceTicketDetail], error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return SingleResponse[AttendanceTicketDetail]{}, fmt.Errorf("id is required")
	}
	endpoint := "/api/v1/attendance-ticket-details/" + url.PathEscape(id)
	return PatchJSON[AttendanceTicketDetailUpdate, SingleResponse[AttendanceTicketDetail]](ctx, c, endpoint, payload)
}

// CreateAttendanceTicketDetail creates an attendance ticket detail for employee and date.
func (c *Client) CreateAttendanceTicketDetail(ctx context.Context, ticketDate, employeeCode string, payload AttendanceTicketDetailCreate) (SingleResponse[AttendanceTicketDetail], error) {
	ticketDate = strings.TrimSpace(ticketDate)
	if ticketDate == "" {
		return SingleResponse[AttendanceTicketDetail]{}, fmt.Errorf("ticket date is required")
	}
	employeeCode = strings.TrimSpace(employeeCode)
	if employeeCode == "" {
		return SingleResponse[AttendanceTicketDetail]{}, fmt.Errorf("employee code is required")
	}
	endpoint := "/api/v1/attendance-tickets/" + url.PathEscape(ticketDate) + "/employees/" + url.PathEscape(employeeCode) + "/attendance-ticket-details"
	return PostJSON[AttendanceTicketDetailCreate, SingleResponse[AttendanceTicketDetail]](ctx, c, endpoint, payload)
}

// GetAttendanceTickets fetches attendance tickets.
func (c *Client) GetAttendanceTickets(ctx context.Context, query url.Values) (ListResponse[AttendanceTicket], error) {
	return GetJSON[ListResponse[AttendanceTicket]](ctx, c, "/api/v1/attendance-tickets", query)
}

// CreateAttendanceTicket creates an attendance ticket.
func (c *Client) CreateAttendanceTicket(ctx context.Context, payload AttendanceTicketCreate) (SingleResponse[AttendanceTicket], error) {
	return PostJSON[AttendanceTicketCreate, SingleResponse[AttendanceTicket]](ctx, c, "/api/v1/attendance-tickets", payload)
}

// GetAttendanceTicket fetches a single attendance ticket by ticket date and employee code.
func (c *Client) GetAttendanceTicket(ctx context.Context, ticketDate, employeeCode string, query url.Values) (SingleResponse[AttendanceTicket], error) {
	ticketDate = strings.TrimSpace(ticketDate)
	if ticketDate == "" {
		return SingleResponse[AttendanceTicket]{}, fmt.Errorf("ticket date is required")
	}
	employeeCode = strings.TrimSpace(employeeCode)
	if employeeCode == "" {
		return SingleResponse[AttendanceTicket]{}, fmt.Errorf("employee code is required")
	}
	endpoint := "/api/v1/attendance-tickets/" + url.PathEscape(ticketDate) + "/employees/" + url.PathEscape(employeeCode)
	return GetJSON[SingleResponse[AttendanceTicket]](ctx, c, endpoint, query)
}
