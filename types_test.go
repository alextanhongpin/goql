package goql_test

import (
	"database/sql"
	"encoding/json"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/alextanhongpin/goql"
)

type types struct {
	// String
	s  string
	sp *string
	sn sql.NullString

	// Int
	i  int
	ip *int
	in sql.NullInt64

	// Float
	f  float64
	fp *float64
	fn sql.NullFloat64

	// Bool
	b  bool
	bp *bool
	bn sql.NullBool

	// Time
	t  time.Time
	tp *time.Time
	tn sql.NullTime

	// JSON
	jr json.RawMessage
	js []string
	ji []int
	jf []float64
	jb []bool
	jt []time.Time

	// Net
	nip   net.IP
	nipp  *net.IP
	nipn  net.IPNet
	nipnp *net.IPNet
}

func TestTypes(t *testing.T) {
	null := true
	array := true

	var typ types
	tests := []struct {
		exp         string
		val         any
		null, array bool
	}{
		{"text", typ.s, false, false},
		{"text", typ.sp, null, false},
		{"text", typ.sn, null, false},
		{"bigint", typ.i, false, false},
		{"bigint", typ.ip, null, false},
		{"bigint", typ.in, null, false},
		{"double precision", typ.f, false, false},
		{"double precision", typ.fp, null, false},
		{"double precision", typ.fn, null, false},
		{"boolean", typ.b, false, false},
		{"boolean", typ.bp, null, false},
		{"boolean", typ.bn, null, false},
		{"timestamptz", typ.t, false, false},
		{"timestamptz", typ.tp, null, false},
		{"timestamptz", typ.tn, null, false},
		{"jsonb", typ.jr, false, false},
		{"text", typ.js, false, array},
		{"bigint", typ.ji, false, array},
		{"double precision", typ.jf, false, array},
		{"boolean", typ.jb, false, array},
		{"timestamptz", typ.jt, false, array},
		{"inet", typ.nip, false, false},
		{"inet", typ.nipp, null, false},
		{"cidr", typ.nipn, false, false},
		{"cidr", typ.nipnp, null, false},
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
