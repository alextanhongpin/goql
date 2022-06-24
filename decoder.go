package goql

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	StructTag    = "sql"
	ValSeparator = ":"
	OpSeparator  = "."
)

var (
	ErrMultipleOperator = errors.New("goql: multiple op")
	ErrUnknownOperator  = errors.New("goql: unknown op")
	ErrBadOperator      = errors.New("goql: bad op")
	ErrInvalidNot       = errors.New("goql: invalid not")
	ErrInvalidIs        = errors.New("goql: 'IS' must be followed by {true, false, null, unknown}")
	ErrInvalidArray     = errors.New("goql: invalid array")
	ErrUnknownField     = errors.New("goql: unknown field")
	ErrUnknownParser    = errors.New("goql: unknown parser")
)

type FieldSet struct {
	Name     string
	Value    any
	RawValue string
	SQLType  string
	IsNull   bool
	IsArray  bool
	Format   string
	Tag      string
	Op       string
}

type Decoder[T any] struct {
	ops     map[string]Op
	columns map[string]Column
	parsers map[string]parserFn
	tag     string
}

func NewDecoder[T any]() *Decoder[T] {
	var t T

	parsers := NewParsers()
	opsByField := make(map[string]Op)
	columns := StructToColumns(t, StructTag)

	for name, col := range columns {
		// By default, all datatypes are comparable.
		ops := OpsComparable

		if col.IsNull {
			ops |= OpsNull
		}

		if col.IsArray {
			ops |= OpsRange
		}

		if IsPgText(col.SQLType) {
			ops |= OpsText
		}

		opsByField[name] = ops
	}

	return &Decoder[T]{
		columns: columns,
		ops:     opsByField,
		tag:     StructTag,
		parsers: parsers,
	}
}

func (d *Decoder[T]) SetStructTag(tag string) {
	if tag == "" {
		panic("tag cannot be empty")
	}

	var t T
	d.tag = tag
	d.columns = StructToColumns(t, tag)
}

func (d *Decoder[T]) SetFieldOps(opsByField map[string]Op) {
	d.ops = opsByField
}

func (d *Decoder[T]) SetParsers(parsers map[string]parserFn) {
	d.parsers = parsers
}

func (d *Decoder[T]) Decode(values url.Values) ([]FieldSet, error) {
	return Decode(d.ops, d.columns, d.parsers, values)
}

func Decode(ops map[string]Op, columns map[string]Column, parsers map[string]parserFn, values url.Values) ([]FieldSet, error) {
	cache := make(map[string]bool)

	var sets []FieldSet

	for field, rule := range ops {
		col, ok := columns[field]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownField, field)
		}

		for _, v := range values[field] {
			ops, val := split2(v, ValSeparator)

			op, ok := ParseOp(ops)
			if !ok {
				return nil, fmt.Errorf("%w: %s", ErrUnknownOperator, v)
			}

			if !rule.Has(op) {
				return nil, fmt.Errorf("%w: %s", ErrBadOperator, v)
			}

			if OpsIs.Has(op) && !sqlIs(val) {
				// OpIs/OpIsNot must have value: true, false, unknown or null.
				return nil, fmt.Errorf("%w: %s", ErrInvalidIs, v)
			}

			cacheKey := fmt.Sprintf("%s:%s", field, op)
			if cache[cacheKey] {
				return nil, fmt.Errorf("%w: %q.%q", ErrMultipleOperator, field, strings.ToLower(op.String()))
			}
			cache[cacheKey] = true

			parser, ok := parsers[col.SQLType]
			if !ok {
				return nil, fmt.Errorf("%w: %s", ErrUnknownParser, col.SQLType)
			}

			fs := FieldSet{
				Name:     field,
				SQLType:  col.SQLType,
				IsNull:   col.IsNull,
				IsArray:  col.IsArray,
				Format:   col.Format,
				Tag:      col.Tag,
				Op:       op.String(),
				RawValue: val,
			}

			var format []string
			if col.Format != "" {
				format = append(format, col.Format)
			}

			switch {
			case OpsIn.Has(op), col.IsArray:
				val, ok = Unquote(val, '{', '}')
				if !ok {
					return nil, fmt.Errorf("%w: missing parantheses: %s", ErrInvalidArray, val)
				}

				vals, err := splitString(val)
				if err != nil {
					return nil, err
				}

				res, err := MapAny(vals, format, parser)
				if err != nil {
					return nil, err
				}

				fs.Value = res
			default:
				if col.IsNull && (val == "null" || val == "") {
					fs.Value = nil
				} else {
					res, err := parser(val, format...)
					if err != nil {
						return nil, err
					}
					fs.Value = res
				}
			}

			sets = append(sets, fs)
		}
	}

	return sets, nil
}
