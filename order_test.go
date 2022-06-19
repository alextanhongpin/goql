package goql_test

import (
	"net/url"
	"testing"

	"github.com/alextanhongpin/goql"
)

func TestOrder(t *testing.T) {
	u := make(url.Values)
	u.Set("sort_by", "name.asc")
	sortBy, err := goql.ParseSortBy(u)
	if err != nil {
		t.Fatalf("ParseSortBy: %s", err)
	}

	if len(sortBy) != 1 {
		t.FailNow()
	}

	if exp, got := "name", sortBy[0].Field; exp != got {
		t.Fatalf("expected %s, got %s", exp, got)
	}

	if exp, got := goql.SortOrderAscending, sortBy[0].Order; exp != got {
		t.Fatalf("expected %s, got %s", exp, got)
	}
}
