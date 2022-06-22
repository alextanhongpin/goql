package goql

import (
	"strconv"
	"time"
)

type parserFn func(s string, format ...string) (any, error)

var defaultParsers = map[string]parserFn{
	pgTypeTimestamp:       parseTime,
	pgTypeTimestampTz:     parseTime,
	pgTypeDate:            parseTime,
	pgTypeTime:            parseTime,
	pgTypeTimeTz:          parseTime,
	pgTypeInterval:        parseNop,
	pgTypeInet:            parseNop,
	pgTypeCidr:            parseNop,
	pgTypeMacaddr:         parseNop,
	pgTypeBoolean:         parseBool,
	pgTypeReal:            parseFloat32,
	pgTypeDoublePrecision: parseFloat64,
	pgTypeSmallint:        parseInt16,
	pgTypeInteger:         parseInt32,
	pgTypeBigint:          parseInt64,
	pgTypeSmallserial:     parseNop,
	pgTypeSerial:          parseNop,
	pgTypeBigserial:       parseNop,
	pgTypeVarchar:         parseText,
	pgTypeChar:            parseText,
	pgTypeText:            parseText,
	pgTypeJSON:            parseText,
	pgTypeJSONB:           parseText,
	pgTypeBytea:           parseText,
}

func MapAny[T any](in []string, format []string, fn func(string, ...string) (T, error)) ([]T, error) {
	res := make([]T, len(in))
	for i, s := range in {
		var err error
		res[i], err = fn(s, format...)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func parseTime(in string, format ...string) (any, error) {
	if len(format) == 1 {
		return time.Parse(format[0], in)
	} else {
		return time.Parse(time.RFC3339, in)
	}
}

func parseNop(in string, format ...string) (any, error) {
	panic("goql: parser not implemented")
}

func parseBool(in string, format ...string) (any, error) {
	if len(format) > 0 {
		panic("goql: invalid args")
	}
	return strconv.ParseBool(in)
}

func parseFloat32(in string, format ...string) (any, error) {
	if len(format) > 0 {
		panic("goql: invalid args")
	}
	return strconv.ParseFloat(in, 32)
}

func parseFloat64(in string, format ...string) (any, error) {
	if len(format) > 0 {
		panic("goql: invalid args")
	}
	return strconv.ParseFloat(in, 64)
}

func parseInt16(in string, format ...string) (any, error) {
	if len(format) > 0 {
		panic("goql: invalid args")
	}
	return strconv.ParseInt(in, 10, 16)
}

func parseInt32(in string, format ...string) (any, error) {
	if len(format) > 0 {
		panic("goql: invalid args")
	}
	return strconv.ParseInt(in, 10, 32)
}

func parseInt64(in string, format ...string) (any, error) {
	if len(format) > 0 {
		panic("goql: invalid args")
	}
	return strconv.ParseInt(in, 10, 64)
}

func parseText(in string, format ...string) (any, error) {
	if len(format) > 0 {
		panic("goql: invalid args")
	}
	return in, nil
}
