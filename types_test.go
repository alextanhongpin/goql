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
		{"string", typ.s, false, false},
		{"*string", typ.sp, null, false},
		{"sql.NullString", typ.sn, false, false},
		{"int", typ.i, false, false},
		{"*int", typ.ip, null, false},
		{"sql.NullInt64", typ.in, false, false},
		{"float64", typ.f, false, false},
		{"*float64", typ.fp, null, false},
		{"sql.NullFloat64", typ.fn, false, false},
		{"bool", typ.b, false, false},
		{"*bool", typ.bp, null, false},
		{"sql.NullBool", typ.bn, false, false},
		{"time.Time", typ.t, false, false},
		{"*time.Time", typ.tp, null, false},
		{"sql.NullTime", typ.tn, false, false},
		{"json.RawMessage", typ.jr, false, array},
		{"[]string", typ.js, false, array},
		{"[]int", typ.ji, false, array},
		{"[]float64", typ.jf, false, array},
		{"[]bool", typ.jb, false, array},
		{"[]time.Time", typ.jt, false, array},
		{"net.IP", typ.nip, false, array},
		{"*net.IP", typ.nipp, null, false},
		{"net.IPNet", typ.nipn, false, false},
		{"*net.IPNet", typ.nipnp, null, false},
	}

	for _, tt := range tests {
		t.Run(tt.exp, func(t *testing.T) {
			typ := goql.TypeOf(reflect.TypeOf(tt.val))
			if exp, got := tt.exp, typ.Name; exp != got {
				t.Fatalf("expected %v, got %v", exp, got)
			}
			if exp, got := tt.null, typ.Null; exp != got {
				t.Fatalf("expected %v, got %v", exp, got)
			}
			if exp, got := tt.array, typ.Array; exp != got {
				t.Fatalf("expected %v, got %v", exp, got)
			}
		})
	}
}
