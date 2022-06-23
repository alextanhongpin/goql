package goql_test

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/alextanhongpin/goql"
)

func TestTypes(t *testing.T) {
	now := time.Now()

	tests := []struct {
		exp         string
		val         any
		null, array bool
	}{
		{"text", "hello world", false, false},
		{"bigint", 1, false, false},
		{"text", []string{"hello", "world"}, false, true},
		{"timestamptz", now, false, false},
		{"timestamptz", &now, true, false},
		{"timestamptz", sql.NullTime{}, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.exp, func(t *testing.T) {
			typ, null, array := goql.GetSQLType(reflect.TypeOf(tt.val))
			if exp, got := tt.exp, typ; exp != got {
				t.Fatalf("expected %v, got %v", exp, got)
			}
			if exp, got := tt.null, null; exp != got {
				t.Fatalf("expected %v, got %v", exp, got)
			}
			if exp, got := tt.array, array; exp != got {
				t.Fatalf("expected %v, got %v", exp, got)
			}
		})
	}
}
