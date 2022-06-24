package goql

import (
	"encoding/json"
	"strconv"
	"time"
)

type ParserFn func(s string) (any, error)

func NewParsers() map[string]ParserFn {
	return map[string]ParserFn{
		"time.Time":  ParseTime,
		"*time.Time": ParseStringPointer[time.Time],
		"bool":       ParseBool,
		"*bool":      ParseStringPointer[bool],
		"float32":    ParseFloat32,
		"*float32":   ParseNumericPointer[float32],
		"float64":    ParseFloat64,
		"*float64":   ParseNumericPointer[float64],
		"int16":      ParseInt16,
		"*int16":     ParseNumericPointer[int16],
		"int32":      ParseInt32,
		"*int32":     ParseNumericPointer[int32],
		"int64":      ParseInt64,
		"*int64":     ParseNumericPointer[int64],
		"int":        ParseInt64,
		"*int":       ParseNumericPointer[int],
		"string":     ParseString,
		"*string":    ParseStringPointer[string],
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

func ParseStringPointer[T string | bool | time.Time](in string) (any, error) {
	var b []byte
	switch in {
	case
		"null",
		"true",
		"false":
		b = []byte(in)
	default:
		b = strconv.AppendQuote(nil, in)
	}

	var t *T
	if err := json.Unmarshal(b, &t); err != nil {
		return nil, err
	}

	return t, nil
}

func ParseNumericPointer[T int | int16 | int32 | int64 | float32 | float64](in string) (any, error) {
	var t *T
	if err := json.Unmarshal([]byte(in), &t); err != nil {
		return nil, err
	}

	return t, nil
}
