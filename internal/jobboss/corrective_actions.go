package jobboss

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// CorrectiveAction models a corrective-preventive action payload.
type CorrectiveAction struct {
	UniqueID                    int     `json:"uniqueID,omitempty"`
	CorrectiveActionNumber      string  `json:"correctiveActionNumber,omitempty"`
	CorrectiveActionCode        string  `json:"correctiveActionCode,omitempty"`
	Status                      string  `json:"status,omitempty"`
	Description                 string  `json:"description,omitempty"`
	CustomerCode                string  `json:"customerCode,omitempty"`
	OrderNumber                 string  `json:"orderNumber,omitempty"`
	PartNumber                  string  `json:"partNumber,omitempty"`
	PONumber                    string  `json:"PONumber,omitempty"`
	QuantityReturned            float64 `json:"quantityReturned,omitempty"`
	QuantityToRework            float64 `json:"quantityToRework,omitempty"`
	EnterBy                     string  `json:"enterBy,omitempty"`
	EnterDate                   string  `json:"enterDate,omitempty"`
	LastModDate                 string  `json:"lastModDate,omitempty"`
	LastModUser                 string  `json:"lastModUser,omitempty"`
	RootCause                   string  `json:"rootCause,omitempty"`
	PermanentAction             string  `json:"permanentAction,omitempty"`
	VerificationOfEffectiveness string  `json:"verificationOfEffectiveness,omitempty"`
}

// GetCorrectiveActions fetches corrective-preventive actions.
func (c *Client) GetCorrectiveActions(ctx context.Context, query url.Values) (ListResponse[CorrectiveAction], error) {
	return GetJSON[ListResponse[CorrectiveAction]](ctx, c, "/api/v1/corrective-preventive-actions", query)
}

// GetCorrectiveActionByNumber fetches one corrective-preventive action by number.
func (c *Client) GetCorrectiveActionByNumber(ctx context.Context, correctiveActionNumber string, query url.Values) (SingleResponse[CorrectiveAction], error) {
	correctiveActionNumber = strings.TrimSpace(correctiveActionNumber)
	if correctiveActionNumber == "" {
		return SingleResponse[CorrectiveAction]{}, fmt.Errorf("corrective action number is required")
	}
	endpoint := "/api/v1/corrective-preventive-actions/" + url.PathEscape(correctiveActionNumber)
	return GetJSON[SingleResponse[CorrectiveAction]](ctx, c, endpoint, query)
}
