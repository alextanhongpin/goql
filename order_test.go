package goql_test

import (
	"testing"

	"github.com/alextanhongpin/goql"
	"github.com/google/go-cmp/cmp"
)

func TestOrder(t *testing.T) {

	tests := []struct {
		order string
		exp   *goql.Order
	}{
		{"", nil},
		{"name", &goql.Order{"name", "asc", "nullslast"}},
		{"name.asc", &goql.Order{"name", "asc", "nullslast"}},
		{"name.desc", &goql.Order{"name", "desc", "nullsfirst"}},
		{"name.asc.nullsfirst", &goql.Order{"name", "asc", "nullsfirst"}},
		{"name.asc.nullslast", &goql.Order{"name", "asc", "nullslast"}},
		{"name.desc.nullsfirst", &goql.Order{"name", "desc", "nullsfirst"}},
		{"name.desc.nullslast", &goql.Order{"name", "desc", "nullslast"}},
	}

	for _, tt := range tests {
		t.Run(tt.order, func(t *testing.T) {
			ord, err := goql.NewOrder(tt.order)
			if err != nil {
				t.FailNow()
			}

			if diff := cmp.Diff(tt.exp, ord); diff != "" {
				t.Fatalf("exp+, got-, %s", diff)
			}
		})
	}
}
