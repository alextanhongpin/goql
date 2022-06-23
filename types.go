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
	sqlNullTimeType    = reflect.TypeOf((*sql.NullTime)(nil)).Elem()
	ipType             = reflect.TypeOf((*net.IP)(nil)).Elem()
	ipNetType          = reflect.TypeOf((*net.IPNet)(nil)).Elem()
	scannerType        = reflect.TypeOf((*sql.Scanner)(nil)).Elem()
	nullBoolType       = reflect.TypeOf((*sql.NullBool)(nil)).Elem()
	nullFloatType      = reflect.TypeOf((*sql.NullFloat64)(nil)).Elem()
	nullIntType        = reflect.TypeOf((*sql.NullInt64)(nil)).Elem()
	nullStringType     = reflect.TypeOf((*sql.NullString)(nil)).Elem()
	jsonRawMessageType = reflect.TypeOf((*json.RawMessage)(nil)).Elem()
)

func sqlType(t reflect.Type) string {
	switch t {
	case timeType, sqlNullTimeType:
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

func GetSQLType(t reflect.Type) (pt string, null bool, array bool) {
	/*
		Handles only the following
		field *Struct
		field Struct
		field []*Struct
	*/

	// Unwraps all pointer first, could be slice or base type.
	switch t.Kind() {
	case reflect.Pointer:
		null = true

		// Get the value of the pointer.
		t = t.Elem()
	case reflect.Slice:
		array = true

		// Get the item of the slice.
		t = t.Elem()

		// Check the item if it is a pointer.
		switch t.Kind() {
		case reflect.Pointer:
			null = true

			t = t.Elem()
		}
	}

	pt = sqlType(t)

	return
}
