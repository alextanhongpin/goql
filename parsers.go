package goql

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const (
	bitSize16 = 16
	bitSize32 = 32
	bitSize64 = 64

	base10 = 10
)

// ParserFn handles conversion of the querystring `string` input to the
// designated types.
type ParserFn func(s string) (any, error)

// NewParsers returns a list of default parsers. This can be extended, and set
// back to the Decoder.
func NewParsers() map[string]ParserFn {
	return map[string]ParserFn{
		"time.Time":       ParseTime,
		"*time.Time":      ParseStringPointer[time.Time],
		"bool":            ParseBool,
		"*bool":           ParseStringPointer[bool],
		"float32":         ParseFloat32,
		"*float32":        ParseNumericPointer[float32],
		"float64":         ParseFloat64,
		"*float64":        ParseNumericPointer[float64],
		"int16":           ParseInt16,
		"*int16":          ParseNumericPointer[int16],
		"int32":           ParseInt32,
		"*int32":          ParseNumericPointer[int32],
		"int64":           ParseInt64,
		"*int64":          ParseNumericPointer[int64],
		"int":             ParseInt,
		"*int":            ParseNumericPointer[int],
		"string":          ParseString,
		"*string":         ParseStringPointer[string],
		"json.RawMessage": ParseJSON,
		"[]byte":          ParseByte,
		"":                ParseNop,
	}
}

// Map applies a function to the list of values.
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

// ParseTime parses string to time with the format RFC3339.
func ParseTime(in string) (any, error) {
	t, err := time.Parse(time.RFC3339, in)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrBadValue, err)
	}

	return t, nil
}

// ParseNop is a nop parser.
func ParseNop(in string) (any, error) {
	panic("goql: parser not implemented")
}

// ParseBool parses string to bool.
func ParseBool(in string) (any, error) {
	t, err := strconv.ParseBool(in)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrBadValue, err)
	}

	return t, nil
}

// ParseFloat32 parses string to float32.
func ParseFloat32(in string) (any, error) {
	return parseFloat(in, bitSize32)
}

// ParseFloat64 parses string to float64.
func ParseFloat64(in string) (any, error) {
	return parseFloat(in, bitSize64)
}

// ParseInt parses string to int.
func ParseInt(in string) (any, error) {
	n, err := strconv.Atoi(in)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrBadValue, err)
	}

	return n, nil
}

// ParseInt16 parses string to int16.
func ParseInt16(in string) (any, error) {
	return parseInt(in, bitSize16)
}

// ParseInt32 parses string to int32.
func ParseInt32(in string) (any, error) {
	return parseInt(in, bitSize32)
}

// ParseInt64 parses string to int64.
func ParseInt64(in string) (any, error) {
	return parseInt(in, bitSize64)
}

// ParseString returns the string as it is.
func ParseString(in string) (any, error) {
	return in, nil
}

// ParseByte parses string to byte slice.
func ParseByte(in string) (any, error) {
	return []byte(in), nil
}

// ParseJSON parses string to json.
func ParseJSON(in string) (any, error) {
	var m map[string]any
	if err := json.Unmarshal([]byte(in), &m); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrBadValue, err)
	}

	return json.RawMessage(in), nil
}

// ParseStringPointer parses json string to pointer type.
func ParseStringPointer[T string | bool | time.Time](in string) (any, error) {
	var b []byte
	switch in {
	// This can still be a string 'true' or 'false', not boolean.
	// Otherwise, we could have just called strconv.ParseBool
	case "null", "true", "false":
		b = []byte(in)
	default:
		b = strconv.AppendQuote(nil, in)
	}

	var t *T
	if err := json.Unmarshal(b, &t); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrBadValue, err)
	}

	return t, nil
}

// ParseNumericPointer parses json string to pointer type.
func ParseNumericPointer[T int | int16 | int32 | int64 | float32 | float64](in string) (any, error) {
	var t *T
	if err := json.Unmarshal([]byte(in), &t); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrBadValue, err)
	}

	return t, nil
}

func parseFloat(in string, bitSize int) (any, error) {
	f, err := strconv.ParseFloat(in, bitSize)
	if err != nil {
		return nil, fmt.Errorf("%w, %s", ErrBadValue, err)
	}

	return f, nil
}

func parseInt(in string, bitSize int) (any, error) {
	f, err := strconv.ParseInt(in, base10, bitSize)
	if err != nil {
		return nil, fmt.Errorf("%w, %s", ErrBadValue, err)
	}

	return f, nil
}
