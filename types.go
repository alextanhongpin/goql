package goql

import (
	"database/sql"
	"encoding/json"
	"net"
	"reflect"
	"time"
)

var (
	timeType           = reflect.TypeOf((*time.Time)(nil)).Elem()
	ipType             = reflect.TypeOf((*net.IP)(nil)).Elem()
	ipNetType          = reflect.TypeOf((*net.IPNet)(nil)).Elem()
	scannerType        = reflect.TypeOf((*sql.Scanner)(nil)).Elem()
	nullBoolType       = reflect.TypeOf((*sql.NullBool)(nil)).Elem()
	nullFloatType      = reflect.TypeOf((*sql.NullFloat64)(nil)).Elem()
	nullIntType        = reflect.TypeOf((*sql.NullInt64)(nil)).Elem()
	nullStringType     = reflect.TypeOf((*sql.NullString)(nil)).Elem()
	nullTimeType       = reflect.TypeOf((*sql.NullTime)(nil)).Elem()
	jsonRawMessageType = reflect.TypeOf((*json.RawMessage)(nil)).Elem()
)

func sqlNullType(t reflect.Type) bool {
	switch t {
	case
		nullBoolType,
		nullFloatType,
		nullIntType,
		nullStringType,
		nullTimeType:
		return true
	default:
		return false
	}
}

func sqlType(t reflect.Type) string {
	switch t {
	case timeType, nullTimeType:
		return pgTypeTimestampTz
	case ipType:
		return pgTypeInet
	case ipNetType:
		return pgTypeCidr
	case nullBoolType:
		return pgTypeBoolean
	case nullFloatType:
		return pgTypeDoublePrecision
	case nullIntType:
		return pgTypeBigint
	case nullStringType:
		return pgTypeText
	case jsonRawMessageType:
		return pgTypeJSONB
	}

	switch t.Kind() {
	case reflect.Int8, reflect.Uint8, reflect.Int16:
		return pgTypeSmallint
	case reflect.Uint16, reflect.Int32:
		return pgTypeInteger
	case reflect.Uint32, reflect.Int64, reflect.Int:
		return pgTypeBigint
	case reflect.Uint, reflect.Uint64:
		// Unsigned bigint is not supported - use bigint.
		return pgTypeBigint
	case reflect.Float32:
		return pgTypeReal
	case reflect.Float64:
		return pgTypeDoublePrecision
	case reflect.Bool:
		return pgTypeBoolean
	case reflect.String:
		return pgTypeText
	case reflect.Map, reflect.Struct:
		return pgTypeJSONB
	case reflect.Array, reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return pgTypeBytea
		}
		return pgTypeJSONB
	default:
		return t.Kind().String()
	}
}

func GetSQLType(t reflect.Type) (typ string, null, array bool) {
	/*
		Handles only the following
		field *Struct
		field Struct
		field []*Struct
	*/
	typ = sqlType(t)
	null = sqlNullType(t)

	switch typ {
	case "ptr":
		null = true

		// Get the value of the pointer.
		t = t.Elem()

		typ = sqlType(t)
	case "jsonb":
		t = t.Elem()
		switch t.Kind() {
		case reflect.Uint8:
		default:
			array = true
			typ = sqlType(t)
		}
	}

	return
}
