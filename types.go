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

func sqlType(typ reflect.Type) string {
	switch typ {
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

	switch typ.Kind() {
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
		if typ.Elem().Kind() == reflect.Uint8 {
			return pgTypeBytea
		}
		return pgTypeJSONB
	default:
		return typ.Kind().String()
	}
}
