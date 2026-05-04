package jobboss

// ListResponse is the common JobBOSS list envelope.
type ListResponse[T any] struct {
	Data []T `json:"Data"`
}

// SingleResponse is the common JobBOSS single-item envelope.
type SingleResponse[T any] struct {
	Data T `json:"Data"`
}

// ProblemResponse models typical API error payload fields.
type ProblemResponse struct {
	Title   string `json:"Title,omitempty"`
	Status  int    `json:"Status,omitempty"`
	Detail  string `json:"Detail,omitempty"`
	TraceID string `json:"TraceId,omitempty"`
}
