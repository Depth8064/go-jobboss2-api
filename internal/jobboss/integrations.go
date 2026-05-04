package jobboss

import (
	"context"
	"net/url"
)

// OAuth2ExchangeRequest models an OAuth2 exchange request payload.
type OAuth2ExchangeRequest struct {
	GrantType    string `json:"grant_type,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// OAuth2ExchangeResponse models an OAuth2 exchange response payload.
type OAuth2ExchangeResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	Scope       string `json:"scope,omitempty"`
}

// ECIAPSAuthenticateRequest models an ECI APS authentication request payload.
type ECIAPSAuthenticateRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// ECIAPSAuthenticateResponse models an ECI APS authentication response payload.
type ECIAPSAuthenticateResponse struct {
	Authenticated bool   `json:"authenticated,omitempty"`
	Token         string `json:"token,omitempty"`
	Message       string `json:"message,omitempty"`
}

// ECIAPSScheduleItem models an ECI APS schedule item payload.
type ECIAPSScheduleItem struct {
	JobNumber     string `json:"jobNumber,omitempty"`
	OperationCode string `json:"operationCode,omitempty"`
	WorkCenter    string `json:"workCenter,omitempty"`
	StartDate     string `json:"startDate,omitempty"`
	DueDate       string `json:"dueDate,omitempty"`
	Status        string `json:"status,omitempty"`
}

// ExchangeOAuth2 exchanges credentials for an OAuth2 token response.
func (c *Client) ExchangeOAuth2(ctx context.Context, req OAuth2ExchangeRequest) (SingleResponse[OAuth2ExchangeResponse], error) {
	return PostJSON[OAuth2ExchangeRequest, SingleResponse[OAuth2ExchangeResponse]](ctx, c, "/api/v1/oauth2Exchange", req)
}

// AuthenticateECIAPS authenticates against ECI APS.
func (c *Client) AuthenticateECIAPS(ctx context.Context, req ECIAPSAuthenticateRequest) (SingleResponse[ECIAPSAuthenticateResponse], error) {
	return PostJSON[ECIAPSAuthenticateRequest, SingleResponse[ECIAPSAuthenticateResponse]](ctx, c, "/api/v1/eci-aps/authenticate-user", req)
}

// GetECIAPSSchedule fetches ECI APS schedule data.
func (c *Client) GetECIAPSSchedule(ctx context.Context, query url.Values) (ListResponse[ECIAPSScheduleItem], error) {
	return GetJSON[ListResponse[ECIAPSScheduleItem]](ctx, c, "/api/v1/eci-aps/get-schedule", query)
}
