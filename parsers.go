package goql

import (
	"strconv"
	"time"
)

type ParserFn func(s string, format ...string) (any, error)

func NewParsers() map[string]ParserFn {
	return map[string]ParserFn{
		pgTypeTimestamp:       ParseTime,
		pgTypeTimestampTz:     ParseTime,
		pgTypeDate:            ParseTime,
		pgTypeTime:            ParseTime,
		pgTypeTimeTz:          ParseTime,
		pgTypeInterval:        ParseNop,
		pgTypeInet:            ParseNop,
		pgTypeCidr:            ParseNop,
		pgTypeMacaddr:         ParseNop,
		pgTypeBoolean:         ParseBool,
		pgTypeReal:            ParseFloat32,
		pgTypeDoublePrecision: ParseFloat64,
		pgTypeSmallint:        ParseInt16,
		pgTypeInteger:         ParseInt32,
		pgTypeBigint:          ParseInt64,
		pgTypeSmallserial:     ParseNop,
		pgTypeSerial:          ParseNop,
		pgTypeBigserial:       ParseNop,
		pgTypeVarchar:         ParseText,
		pgTypeChar:            ParseText,
		pgTypeText:            ParseText,
		pgTypeJSON:            ParseText,
		pgTypeJSONB:           ParseText,
		pgTypeBytea:           ParseText,
	}
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

func ParseTime(in string, format ...string) (any, error) {
	if len(format) == 1 {
		return time.Parse(format[0], in)
	} else {
		return time.Parse(time.RFC3339, in)
	}
}

func ParseNop(in string, format ...string) (any, error) {
	panic("goql: parser not implemented")
}

func ParseBool(in string, format ...string) (any, error) {
	if len(format) > 0 {
		panic("goql: invalid args")
	}
	return strconv.ParseBool(in)
}

func ParseFloat32(in string, format ...string) (any, error) {
	if len(format) > 0 {
		panic("goql: invalid args")
	}
	return strconv.ParseFloat(in, 32)
}

func ParseFloat64(in string, format ...string) (any, error) {
	if len(format) > 0 {
		panic("goql: invalid args")
	}
	return strconv.ParseFloat(in, 64)
}

func ParseInt16(in string, format ...string) (any, error) {
	if len(format) > 0 {
		panic("goql: invalid args")
	}
	return strconv.ParseInt(in, 10, 16)
}

func ParseInt32(in string, format ...string) (any, error) {
	if len(format) > 0 {
		panic("goql: invalid args")
	}
	return strconv.ParseInt(in, 10, 32)
}

func ParseInt64(in string, format ...string) (any, error) {
	if len(format) > 0 {
		panic("goql: invalid args")
	}
	return strconv.ParseInt(in, 10, 64)
}

func ParseText(in string, format ...string) (any, error) {
	if len(format) > 0 {
		panic("goql: invalid args")
	}
	return in, nil
}
