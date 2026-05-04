package jobboss_test

import (
	"testing"
	"time"

	jobboss "go-jobboss2-api"
)

func TestDateFormatting(t *testing.T) {
	tm := time.Date(2026, 5, 4, 15, 16, 17, 0, time.FixedZone("X", -7*3600))
	if got := jobboss.FormatUTCDate(tm); got != "2026-05-04" {
		t.Fatalf("unexpected date: %s", got)
	}
	if got := jobboss.FormatUTCDateTime(tm); got != "2026-05-04T22:16:17Z" {
		t.Fatalf("unexpected datetime: %s", got)
	}
}
