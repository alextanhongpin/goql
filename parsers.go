package goql

import (
	"encoding/json"
	"strconv"
	"time"
)

func newByte(in string) (b []byte) {
	switch in {
	case
		"null",
		"true",
		"false":
		b = []byte(in)
	default:
		b = strconv.AppendQuote(nil, in)
	}

	return
}

type ParserFn func(s string) (any, error)

func NewParsers() map[string]ParserFn {
	return map[string]ParserFn{
		"time.Time":  ParseTime,
		"*time.Time": ParsePointer[time.Time],
		"bool":       ParseBool,
		"*bool":      ParsePointer[bool],
		"float32":    ParseFloat32,
		"*float32":   ParsePointer[float32],
		"float64":    ParseFloat64,
		"*float64":   ParsePointer[float64],
		"int16":      ParseInt16,
		"*int16":     ParsePointer[int16],
		"int32":      ParseInt32,
		"*int32":     ParsePointer[int32],
		"int64":      ParseInt64,
		"*int64":     ParsePointer[int64],
		"int":        ParseInt64,
		"*int":       ParsePointer[int],
		"string":     ParseString,
		"*string":    ParsePointer[string],
		"":           ParseNop,
	}
}

func Map[T, K any](in []T, fn func(T) (K, error)) ([]K, error) {
	res := make([]K, len(in))
	for i, s := range in {
		var err error
		res[i], err = fn(s)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func ParseTime(in string) (any, error) {
	return time.Parse(time.RFC3339, in)
}

func ParseNop(in string) (any, error) {
	panic("goql: parser not implemented")
}

func ParseBool(in string) (any, error) {
	return strconv.ParseBool(in)
}

func ParseFloat32(in string) (any, error) {
	return strconv.ParseFloat(in, 32)
}

func ParseFloat64(in string) (any, error) {
	return strconv.ParseFloat(in, 64)
}

func ParseInt16(in string) (any, error) {
	return strconv.ParseInt(in, 10, 16)
}

func ParseInt32(in string) (any, error) {
	return strconv.ParseInt(in, 10, 32)
}

func ParseInt64(in string) (any, error) {
	return strconv.ParseInt(in, 10, 64)
}

func ParseString(in string) (any, error) {
	return in, nil
}

func ParsePointer[T any](in string) (any, error) {
	b := newByte(in)

	var t *T
	if err := json.Unmarshal(b, &t); err != nil {
		return nil, err
	}

	return t, nil
}
