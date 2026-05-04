package jobboss_test

import (
	"testing"

	jobboss "go-jobboss2-api"
)

func TestQueryHelpers(t *testing.T) {
	q := jobboss.NewQuery()
	if q == nil || len(q) != 0 {
		t.Fatalf("expected empty query map")
	}

	q = jobboss.ApplyPaging(nil, 0, -1)
	if q.Get("take") != "200" || q.Get("skip") != "" {
		t.Fatalf("unexpected paging defaults: %v", q)
	}

	q = jobboss.ApplyPaging(q, 500, 20)
	if q.Get("take") != "200" || q.Get("skip") != "20" {
		t.Fatalf("unexpected clamped paging: %v", q)
	}

	q = jobboss.ApplyFields(nil, " a ", "", "b")
	if q.Get("fields") != "a,b" {
		t.Fatalf("unexpected fields: %q", q.Get("fields"))
	}

	q = jobboss.ApplySort(nil, " +x", "", "-y ")
	if q.Get("sort") != "+x,-y" {
		t.Fatalf("unexpected sort: %q", q.Get("sort"))
	}

	q = jobboss.ApplyFilter(nil, " ", jobboss.FilterEqual, "v")
	if len(q) != 0 {
		t.Fatalf("expected no filter when field empty")
	}

	q = jobboss.ApplyFilter(nil, "field", "", "value")
	if q.Get("field") != "value" {
		t.Fatalf("unexpected simple filter: %v", q)
	}

	q = jobboss.ApplyFilter(nil, "field", jobboss.FilterGreaterThan, "9")
	if q.Get("field[gt]") != "9" {
		t.Fatalf("unexpected op filter: %v", q)
	}

	q = jobboss.ApplyInFilter(nil, "field", " ", "a", "b ")
	if q.Get("field[in]") != "a|b" {
		t.Fatalf("unexpected in filter: %v", q)
	}

	q = jobboss.ApplyInFilter(nil, "field", " ")
	if q != nil {
		t.Fatalf("expected nil query when no values and q nil")
	}
}
