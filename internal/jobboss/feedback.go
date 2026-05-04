package jobboss

import (
	"context"
	"net/url"
)

// Feedback models a feedback payload.
type Feedback struct {
	UniqueID             int    `json:"uniqueID,omitempty"`
	FeedbackNo           string `json:"feedbackNo,omitempty"`
	FeedbackCode         string `json:"feedbackCode,omitempty"`
	FeedbackType         string `json:"feedbackType,omitempty"`
	Status               string `json:"status,omitempty"`
	CustomerCode         string `json:"customerCode,omitempty"`
	VendorCode           string `json:"vendorCode,omitempty"`
	OrderNumber          string `json:"orderNumber,omitempty"`
	Comments             string `json:"comments,omitempty"`
	Resolution           string `json:"resolution,omitempty"`
	EnterBy              string `json:"enterBy,omitempty"`
	EnterDate            string `json:"enterDate,omitempty"`
	LastModDate          string `json:"lastModDate,omitempty"`
	LastModUser          string `json:"lastModUser,omitempty"`
	CorrectiveActionCode string `json:"correctiveActionCode,omitempty"`
}

// GetFeedback fetches feedback records.
func (c *Client) GetFeedback(ctx context.Context, query url.Values) (ListResponse[Feedback], error) {
	return GetJSON[ListResponse[Feedback]](ctx, c, "/api/v1/feedback", query)
}
