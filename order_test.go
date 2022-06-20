package goql_test

import (
	"testing"

	"github.com/alextanhongpin/goql"
)

func TestOrder(t *testing.T) {
	sortBy, err := goql.ParseOrder("name.asc")
	if err != nil {
		t.Fatalf("ParseOrder: %s", err)
	}

	if len(sortBy) != 1 {
		t.FailNow()
	}

	if exp, got := "name", sortBy[0].Column; exp != got {
		t.Fatalf("expected %s, got %s", exp, got)
	}

	if exp, got := goql.DirectionAscending, sortBy[0].Direction; exp != got {
		t.Fatalf("expected %s, got %s", exp, got)
	}
}
