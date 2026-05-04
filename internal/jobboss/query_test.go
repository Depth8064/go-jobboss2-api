package jobboss

import (
	"testing"
)

func TestQueryHelpers(t *testing.T) {
	q := NewQuery()
	if q == nil || len(q) != 0 {
		t.Fatalf("expected empty query map")
	}

	q = ApplyPaging(nil, 0, -1)
	if q.Get("take") != "200" || q.Get("skip") != "" {
		t.Fatalf("unexpected paging defaults: %v", q)
	}

	q = ApplyPaging(q, 500, 20)
	if q.Get("take") != "200" || q.Get("skip") != "20" {
		t.Fatalf("unexpected clamped paging: %v", q)
	}

	q = ApplyFields(nil, " a ", "", "b")
	if q.Get("fields") != "a,b" {
		t.Fatalf("unexpected fields: %q", q.Get("fields"))
	}

	q = ApplySort(nil, " +x", "", "-y ")
	if q.Get("sort") != "+x,-y" {
		t.Fatalf("unexpected sort: %q", q.Get("sort"))
	}

	q = ApplyFilter(nil, " ", FilterEqual, "v")
	if len(q) != 0 {
		t.Fatalf("expected no filter when field empty")
	}

	q = ApplyFilter(nil, "field", "", "value")
	if q.Get("field") != "value" {
		t.Fatalf("unexpected simple filter: %v", q)
	}

	q = ApplyFilter(nil, "field", FilterGreaterThan, "9")
	if q.Get("field[gt]") != "9" {
		t.Fatalf("unexpected op filter: %v", q)
	}

	q = ApplyInFilter(nil, "field", " ", "a", "b ")
	if q.Get("field[in]") != "a|b" {
		t.Fatalf("unexpected in filter: %v", q)
	}

	q = ApplyInFilter(nil, "field", " ")
	if q != nil {
		t.Fatalf("expected nil query when no values and q nil")
	}
}
