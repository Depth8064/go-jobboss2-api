package jobboss

import "time"

const (
	jobbossDateFormat     = "2006-01-02"
	jobbossDateTimeFormat = "2006-01-02T15:04:05Z"
)

// FormatUTCDate formats a date as yyyy-MM-dd in UTC.
func FormatUTCDate(t time.Time) string {
	return t.UTC().Format(jobbossDateFormat)
}

// FormatUTCDateTime formats a datetime as yyyy-MM-ddTHH:mm:ssZ in UTC.
func FormatUTCDateTime(t time.Time) string {
	return t.UTC().Format(jobbossDateTimeFormat)
}
