package jobboss

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// TimeTicket models a time ticket payload.
type TimeTicket struct {
	UniqueID     int     `json:"uniqueID,omitempty"`
	TicketDate   string  `json:"ticketDate,omitempty"`
	EmployeeCode string  `json:"employeeCode,omitempty"`
	Department   string  `json:"department,omitempty"`
	Status       string  `json:"status,omitempty"`
	HoursWorked  float64 `json:"hoursWorked,omitempty"`
	Comments     string  `json:"comments,omitempty"`
	LastModDate  string  `json:"lastModDate,omitempty"`
	LastModUser  string  `json:"lastModUser,omitempty"`
}

// TimeTicketDetail models a time ticket detail payload.
type TimeTicketDetail struct {
	UniqueID       int     `json:"uniqueID,omitempty"`
	TimeTicketGUID string  `json:"timeTicketGUID,omitempty"`
	OrderNumber    string  `json:"orderNumber,omitempty"`
	ItemNumber     int     `json:"itemNumber,omitempty"`
	OperationCode  string  `json:"operationCode,omitempty"`
	HoursWorked    float64 `json:"hoursWorked,omitempty"`
	Status         string  `json:"status,omitempty"`
	LastModDate    string  `json:"lastModDate,omitempty"`
	LastModUser    string  `json:"lastModUser,omitempty"`
}

// GetTimeTickets fetches time tickets.
func (c *Client) GetTimeTickets(ctx context.Context, query url.Values) (ListResponse[TimeTicket], error) {
	return GetJSON[ListResponse[TimeTicket]](ctx, c, "/api/v1/time-tickets", query)
}

// CreateTimeTicket creates a time ticket.
func (c *Client) CreateTimeTicket(ctx context.Context, ticket TimeTicket) (SingleResponse[TimeTicket], error) {
	return PostJSON[TimeTicket, SingleResponse[TimeTicket]](ctx, c, "/api/v1/time-tickets", ticket)
}

// GetTimeTicketByDateAndEmployee fetches a time ticket by ticket date and employee code.
func (c *Client) GetTimeTicketByDateAndEmployee(ctx context.Context, ticketDate, employeeCode string, query url.Values) (SingleResponse[TimeTicket], error) {
	ticketDate = strings.TrimSpace(ticketDate)
	employeeCode = strings.TrimSpace(employeeCode)
	if ticketDate == "" {
		return SingleResponse[TimeTicket]{}, fmt.Errorf("ticket date is required")
	}
	if employeeCode == "" {
		return SingleResponse[TimeTicket]{}, fmt.Errorf("employee code is required")
	}
	endpoint := "/api/v1/time-tickets/" + url.PathEscape(ticketDate) + "/employees/" + url.PathEscape(employeeCode)
	return GetJSON[SingleResponse[TimeTicket]](ctx, c, endpoint, query)
}

// UpdateTimeTicketByDateAndEmployee updates a time ticket by ticket date and employee code.
func (c *Client) UpdateTimeTicketByDateAndEmployee(ctx context.Context, ticketDate, employeeCode string, ticket TimeTicket) (UpdateResult, error) {
	ticketDate = strings.TrimSpace(ticketDate)
	employeeCode = strings.TrimSpace(employeeCode)
	if ticketDate == "" {
		return UpdateResult{}, fmt.Errorf("ticket date is required")
	}
	if employeeCode == "" {
		return UpdateResult{}, fmt.Errorf("employee code is required")
	}
	endpoint := "/api/v1/time-tickets/" + url.PathEscape(ticketDate) + "/employees/" + url.PathEscape(employeeCode)
	return PatchNoContentJSON[TimeTicket](ctx, c, endpoint, ticket)
}

// GetTimeTicketDetails fetches time ticket details.
func (c *Client) GetTimeTicketDetails(ctx context.Context, query url.Values) (ListResponse[TimeTicketDetail], error) {
	return GetJSON[ListResponse[TimeTicketDetail]](ctx, c, "/api/v1/time-ticket-details", query)
}

// CreateTimeTicketDetail creates a time ticket detail.
func (c *Client) CreateTimeTicketDetail(ctx context.Context, detail TimeTicketDetail) (SingleResponse[TimeTicketDetail], error) {
	return PostJSON[TimeTicketDetail, SingleResponse[TimeTicketDetail]](ctx, c, "/api/v1/time-ticket-details", detail)
}

// GetTimeTicketDetailByGUID fetches a time ticket detail by GUID.
func (c *Client) GetTimeTicketDetailByGUID(ctx context.Context, timeTicketGUID string, query url.Values) (SingleResponse[TimeTicketDetail], error) {
	timeTicketGUID = strings.TrimSpace(timeTicketGUID)
	if timeTicketGUID == "" {
		return SingleResponse[TimeTicketDetail]{}, fmt.Errorf("time ticket guid is required")
	}
	endpoint := "/api/v1/time-ticket-details/" + url.PathEscape(timeTicketGUID)
	return GetJSON[SingleResponse[TimeTicketDetail]](ctx, c, endpoint, query)
}

// UpdateTimeTicketDetailByGUID updates a time ticket detail by GUID.
func (c *Client) UpdateTimeTicketDetailByGUID(ctx context.Context, timeTicketGUID string, detail TimeTicketDetail) (UpdateResult, error) {
	timeTicketGUID = strings.TrimSpace(timeTicketGUID)
	if timeTicketGUID == "" {
		return UpdateResult{}, fmt.Errorf("time ticket guid is required")
	}
	endpoint := "/api/v1/time-ticket-details/" + url.PathEscape(timeTicketGUID)
	return PatchNoContentJSON[TimeTicketDetail](ctx, c, endpoint, detail)
}
