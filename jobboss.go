package jobboss

import (
	"context"
	"net/url"
	"time"

	internal "github.com/Depth8064/go-jobboss2-api/internal/jobboss"
)

// Public type aliases re-exported from internal implementation.
type (
	Config                       = internal.Config
	Client                       = internal.Client
	TokenResponse                = internal.TokenResponse
	APIError                     = internal.APIError
	FilterOp                     = internal.FilterOp
	DeviceCredentialType         = internal.DeviceCredentialType
	DeviceCreateModel            = internal.DeviceCreateModel
	DeviceCreateResponse         = internal.DeviceCreateResponse
	DeviceActivateRequestModel   = internal.DeviceActivateRequestModel
	DeviceCredentialModel        = internal.DeviceCredentialModel
	DeviceActivatedResponse      = internal.DeviceActivatedResponse
	Customer                     = internal.Customer
	CustomerCreate               = internal.CustomerCreate
	CustomerUpdate               = internal.CustomerUpdate
	Order                        = internal.Order
	OrderCreate                  = internal.OrderCreate
	OrderUpdate                  = internal.OrderUpdate
	OrderLineItem                = internal.OrderLineItem
	OrderLineItemCreate          = internal.OrderLineItemCreate
	OrderRouting                 = internal.OrderRouting
	OrderRoutingCreate           = internal.OrderRoutingCreate
	Release                      = internal.Release
	ReleaseCreate                = internal.ReleaseCreate
	ARInvoiceDetail              = internal.ARInvoiceDetail
	ARInvoice                    = internal.ARInvoice
	AttendanceTicketDetail       = internal.AttendanceTicketDetail
	AttendanceTicketDetailCreate = internal.AttendanceTicketDetailCreate
	AttendanceTicketDetailUpdate = internal.AttendanceTicketDetailUpdate
	AttendanceTicket             = internal.AttendanceTicket
	AttendanceTicketCreate       = internal.AttendanceTicketCreate
	BinLocation                  = internal.BinLocation
	Company                      = internal.Company
	CompanyCalendar              = internal.CompanyCalendar
	Contact                      = internal.Contact
	ContactCreate                = internal.ContactCreate
	ContactUpdate                = internal.ContactUpdate
	CorrectiveAction             = internal.CorrectiveAction
	CurrencyCode                 = internal.CurrencyCode
	Department                   = internal.Department
	DocumentControl              = internal.DocumentControl
	DocumentHistory              = internal.DocumentHistory
	DocumentReview               = internal.DocumentReview
	EmployeeTraining             = internal.EmployeeTraining
	Employee                     = internal.Employee
	EmployeeUpdate               = internal.EmployeeUpdate
	Estimate                     = internal.Estimate
	EstimateCreate               = internal.EstimateCreate
	EstimateUpdate               = internal.EstimateUpdate
	Feedback                     = internal.Feedback
	GLCode                       = internal.GLCode
	JobMaterial                  = internal.JobMaterial
	JobRequirement               = internal.JobRequirement
	Material                     = internal.Material
	NonConformance               = internal.NonConformance
	OperationCode                = internal.OperationCode
	PackingList                  = internal.PackingList
	PackingListLineItem          = internal.PackingListLineItem
	ProblemResponse              = internal.ProblemResponse
	ProductCode                  = internal.ProductCode
	PurchaseOrder                = internal.PurchaseOrder
	PurchaseOrderLineItem        = internal.PurchaseOrderLineItem
	PurchaseOrderRelease         = internal.PurchaseOrderRelease
	Quote                        = internal.Quote
	QuoteLineItem                = internal.QuoteLineItem
)

type ListResponse[T any] struct {
	Data []T `json:"Data"`
}

type SingleResponse[T any] struct {
	Data T `json:"Data"`
}

const (
	DefaultTake = internal.DefaultTake
	MaxTake     = internal.MaxTake

	FilterEqual            = internal.FilterEqual
	FilterNotEqual         = internal.FilterNotEqual
	FilterGreaterThan      = internal.FilterGreaterThan
	FilterGreaterThanEqual = internal.FilterGreaterThanEqual
	FilterLessThan         = internal.FilterLessThan
	FilterLessThanEqual    = internal.FilterLessThanEqual
	FilterIn               = internal.FilterIn
	FilterNotIn            = internal.FilterNotIn
	FilterNull             = internal.FilterNull

	DeviceCredentialNone           = internal.DeviceCredentialNone
	DeviceCredentialToken          = internal.DeviceCredentialToken
	DeviceCredentialSecret         = internal.DeviceCredentialSecret
	DeviceCredentialTokenAndSecret = internal.DeviceCredentialTokenAndSecret
)

func NewClient(cfg Config) *Client { return internal.NewClient(cfg) }

func NewQuery() url.Values { return internal.NewQuery() }

func ApplyPaging(q url.Values, take, skip int) url.Values { return internal.ApplyPaging(q, take, skip) }

func ApplyFields(q url.Values, fields ...string) url.Values {
	return internal.ApplyFields(q, fields...)
}

func ApplySort(q url.Values, expressions ...string) url.Values {
	return internal.ApplySort(q, expressions...)
}

func ApplyFilter(q url.Values, field string, op FilterOp, value string) url.Values {
	return internal.ApplyFilter(q, field, op, value)
}

func ApplyInFilter(q url.Values, field string, values ...string) url.Values {
	return internal.ApplyInFilter(q, field, values...)
}

func FormatUTCDate(t time.Time) string { return internal.FormatUTCDate(t) }

func FormatUTCDateTime(t time.Time) string { return internal.FormatUTCDateTime(t) }

func GetJSON[T any](ctx context.Context, c *Client, endpoint string, query url.Values) (T, error) {
	return internal.GetJSON[T](ctx, c, endpoint, query)
}

func PostJSON[Req any, Resp any](ctx context.Context, c *Client, endpoint string, payload Req) (Resp, error) {
	return internal.PostJSON[Req, Resp](ctx, c, endpoint, payload)
}

func PatchJSON[Req any, Resp any](ctx context.Context, c *Client, endpoint string, payload Req) (Resp, error) {
	return internal.PatchJSON[Req, Resp](ctx, c, endpoint, payload)
}
