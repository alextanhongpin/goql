package goql_test

import (
	"database/sql"
	"encoding/json"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/alextanhongpin/goql"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
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

	uuid  uuid.UUID
	uuidp *uuid.UUID
}

func TestTypes(t *testing.T) {
	null := true
	array := true

	var typ types
	tests := []struct {
		typ any
		exp goql.Type
	}{
		{typ.s, goql.Type{"string", false, false}},
		{typ.sp, goql.Type{"*string", null, false}},
		{typ.sn, goql.Type{"sql.NullString", false, false}},
		{typ.i, goql.Type{"int", false, false}},
		{typ.ip, goql.Type{"*int", null, false}},
		{typ.in, goql.Type{"sql.NullInt64", false, false}},
		{typ.f, goql.Type{"float64", false, false}},
		{typ.fp, goql.Type{"*float64", null, false}},
		{typ.fn, goql.Type{"sql.NullFloat64", false, false}},
		{typ.b, goql.Type{"bool", false, false}},
		{typ.bp, goql.Type{"*bool", null, false}},
		{typ.bn, goql.Type{"sql.NullBool", false, false}},
		{typ.t, goql.Type{"time.Time", false, false}},
		{typ.tp, goql.Type{"*time.Time", null, false}},
		{typ.tn, goql.Type{"sql.NullTime", false, false}},
		{typ.nip, goql.Type{"net.IP", false, array}},
		{typ.nipp, goql.Type{"*net.IP", null, array}},
		{typ.nipn, goql.Type{"net.IPNet", false, false}},
		{typ.nipnp, goql.Type{"*net.IPNet", null, false}},
		{typ.jr, goql.Type{"json.RawMessage", false, array}},
		{typ.js, goql.Type{"[]string", false, array}},
		{typ.ji, goql.Type{"[]int", false, array}},
		{typ.jf, goql.Type{"[]float64", false, array}},
		{typ.jb, goql.Type{"[]bool", false, array}},
		{typ.jt, goql.Type{"[]time.Time", false, array}},
		{typ.uuid, goql.Type{"uuid.UUID", false, array}},
		{typ.uuidp, goql.Type{"*uuid.UUID", null, array}},
	}

	for _, tt := range tests {
		t.Run(tt.exp.Name, func(t *testing.T) {
			got := goql.TypeOf(reflect.TypeOf(tt.typ))

			if diff := cmp.Diff(tt.exp, got); diff != "" {
				t.Fatalf("exp+, got-: %s", diff)
			}
		})
	}
}
