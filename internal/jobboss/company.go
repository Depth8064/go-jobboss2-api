package jobboss

import (
	"context"
	"net/url"
)

// Company models a company payload.
type Company struct {
	CompanyCode               string `json:"companyCode,omitempty"`
	CompanyDescription        string `json:"companyDescription,omitempty"`
	CompanyCountry            string `json:"companyCountry,omitempty"`
	CompanyTimeZone           string `json:"companyTimeZone,omitempty"`
	CurrencyCode              string `json:"currencyCode,omitempty"`
	ProductCode               string `json:"productCode,omitempty"`
	WorkCode                  string `json:"workCode,omitempty"`
	ReportDateFormat          string `json:"reportDateFormat,omitempty"`
	AccountsReceivableAccount string `json:"accountsReceivableAccount,omitempty"`
	AccountsPayableAccount    string `json:"accountsPayableAccount,omitempty"`
}

// CompanyCalendar models a company calendar payload.
type CompanyCalendar struct {
	UniqueID             int     `json:"uniqueID,omitempty"`
	CalendarType         string  `json:"calendarType,omitempty"`
	Code                 string  `json:"code,omitempty"`
	Date                 string  `json:"date,omitempty"`
	SearchDate           int     `json:"searchDate,omitempty"`
	CapacityFactor       float64 `json:"capacityFactor,omitempty"`
	HoursAvailable       float64 `json:"hoursAvailable,omitempty"`
	HoursAvailableShift1 float64 `json:"hoursAvailableShift1,omitempty"`
	HoursAvailableShift2 float64 `json:"hoursAvailableShift2,omitempty"`
	HoursAvailableShift3 float64 `json:"hoursAvailableShift3,omitempty"`
	Shift1Begin          string  `json:"shift1Begin,omitempty"`
	Shift1End            string  `json:"shift1End,omitempty"`
	Shift2Begin          string  `json:"shift2Begin,omitempty"`
	Shift2End            string  `json:"shift2End,omitempty"`
	Shift3Begin          string  `json:"shift3Begin,omitempty"`
	Shift3End            string  `json:"shift3End,omitempty"`
	LastModDate          string  `json:"lastModDate,omitempty"`
	LastModUser          string  `json:"lastModUser,omitempty"`
}

// GetCompany fetches company data from /api/v1/company.
func (c *Client) GetCompany(ctx context.Context, query url.Values) (ListResponse[Company], error) {
	return GetJSON[ListResponse[Company]](ctx, c, "/api/v1/company", query)
}

// GetCompanyCalendars fetches company calendars from /api/v1/company-calendars.
func (c *Client) GetCompanyCalendars(ctx context.Context, query url.Values) (ListResponse[CompanyCalendar], error) {
	return GetJSON[ListResponse[CompanyCalendar]](ctx, c, "/api/v1/company-calendars", query)
}
