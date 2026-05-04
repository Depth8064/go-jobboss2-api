package jobboss

import (
	"context"
	"net/url"
)

// DocumentControl models a document control payload.
type DocumentControl struct {
	UniqueID       int    `json:"uniqueID,omitempty"`
	DocumentNumber string `json:"documentNumber,omitempty"`
	DocumentType   string `json:"documentType,omitempty"`
	DocumentStatus string `json:"documentStatus,omitempty"`
	Description    string `json:"description,omitempty"`
	Revision       string `json:"revision,omitempty"`
	RevisionDate   string `json:"revisionDate,omitempty"`
	DocumentDate   string `json:"documentDate,omitempty"`
	DateEntered    string `json:"dateEntered,omitempty"`
	EnteredBy      string `json:"enteredBy,omitempty"`
	FileLocation   string `json:"fileLocation,omitempty"`
	RepositoryID   int    `json:"repositoryID,omitempty"`
	LastModDate    string `json:"lastModDate,omitempty"`
	LastModUser    string `json:"lastModUser,omitempty"`
	Printed        string `json:"printed,omitempty"`
}

// DocumentHistory models a document history payload.
type DocumentHistory struct {
	UniqueID       int    `json:"uniqueID,omitempty"`
	DocumentNumber string `json:"documentNumber,omitempty"`
	Revision       string `json:"revision,omitempty"`
	RevisionDate   string `json:"revisionDate,omitempty"`
	Comments       string `json:"comments,omitempty"`
	FileLocation   string `json:"fileLocation,omitempty"`
	RepositoryID   int    `json:"repositoryID,omitempty"`
	SpunOffFile    string `json:"spunOffFile,omitempty"`
	UpdatedJobs    string `json:"updatedJobs,omitempty"`
	UpdatedParts   string `json:"updatedParts,omitempty"`
	UserID         string `json:"userID,omitempty"`
}

// DocumentReview models a document review payload.
type DocumentReview struct {
	UniqueID       int     `json:"uniqueID,omitempty"`
	DocumentNumber string  `json:"documentNumber,omitempty"`
	ReviewCode     string  `json:"reviewCode,omitempty"`
	Description    string  `json:"description,omitempty"`
	EmployeeCode   int     `json:"employeeCode,omitempty"`
	VendorCode     string  `json:"vendorCode,omitempty"`
	JobNumber      string  `json:"jobNumber,omitempty"`
	InvoiceNumber  string  `json:"invoiceNumber,omitempty"`
	StartDate      string  `json:"startDate,omitempty"`
	EndDate        string  `json:"endDate,omitempty"`
	Completed      string  `json:"completed,omitempty"`
	Cost           float64 `json:"cost,omitempty"`
	LastModDate    string  `json:"lastModDate,omitempty"`
	LastModUser    string  `json:"lastModUser,omitempty"`
}

// GetDocumentControls fetches document controls from /api/v1/document-controls.
func (c *Client) GetDocumentControls(ctx context.Context, query url.Values) (ListResponse[DocumentControl], error) {
	return GetJSON[ListResponse[DocumentControl]](ctx, c, "/api/v1/document-controls", query)
}

// GetDocumentHistories fetches document histories from /api/v1/document-histories.
func (c *Client) GetDocumentHistories(ctx context.Context, query url.Values) (ListResponse[DocumentHistory], error) {
	return GetJSON[ListResponse[DocumentHistory]](ctx, c, "/api/v1/document-histories", query)
}

// GetDocumentReviews fetches document review records from /api/v1/document-review.
func (c *Client) GetDocumentReviews(ctx context.Context, query url.Values) (ListResponse[DocumentReview], error) {
	return GetJSON[ListResponse[DocumentReview]](ctx, c, "/api/v1/document-review", query)
}
