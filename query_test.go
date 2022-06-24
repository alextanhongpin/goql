package goql_test

import (
	"net/url"
	"testing"

	"github.com/alextanhongpin/goql"
)

func TestQuery(t *testing.T) {
	v := url.Values{
		"name":    []string{"eq:john", "neq:jane"},
		"age":     []string{"gt:10", "lt:100"},
		"married": []string{"is:true", "bad value"},
		"and":     []string{"(age.is:true,or(age.eq:13, age.eq.17))"},
		"or":      []string{"married.isnot:true"},
	}

	queries := goql.ParseQuery(v)

	if exp, got := 5, len(queries); exp != got {
		t.Fatalf("expected %d, got %d: %v", exp, got, queries)
	}
	t.Logf("queries: %+v", queries)
}
