package goql_test

import (
	"net/url"
	"testing"

	"github.com/alextanhongpin/goql"
)

func TestQuery(t *testing.T) {
	v := make(url.Values)
	v.Set("name.eq", "john")
	v.Set("name.neq", "jane")
	v.Set("age.gt", "10")
	v.Set("age.lt", "100")
	v.Set("married.is", "true")
	v.Set("and", "(age.is:true,or(age.eq:13,age.eq.17))")
	v.Set("or", "(married.isnot:true)")

	v = goql.FilterValues(v, goql.OpAnd.String(), goql.OpOr.String())
	queries, err := goql.ParseQuery(v)
	if err != nil {
		t.Fatalf("failed to parse query: %s", err)
	}

	if exp, got := 5, len(queries); exp != got {
		t.Fatalf("expected %d, got %d: %v", exp, got, queries)
	}
	t.Logf("queries: %+v", queries)
}
