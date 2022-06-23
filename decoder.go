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
	Not      bool
	IsNull   bool
	IsArray  bool
	Format   string
	Tag      string
	Op       string
}

type Decoder struct {
	ops     map[string]Op
	columns map[string]Column
	parsers map[string]parserFn
	tag     string
}

func NewDecoder(v any) *Decoder {
	parsers := NewParsers()
	opsByField := make(map[string]Op)
	columns := StructToColumns(v, StructTag)

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

	return &Decoder{
		columns: columns,
		ops:     opsByField,
		tag:     StructTag,
		parsers: parsers,
	}
}

func (d *Decoder) SetStructTag(tag string) {
	if tag == "" {
		panic("tag cannot be empty")
	}

	d.tag = tag
}

func (d *Decoder) SetFieldOps(opsByField map[string]Op) {
	d.ops = opsByField
}

func (d *Decoder) SetParsers(parsers map[string]parserFn) {
	d.parsers = parsers
}

func (d *Decoder) Decode(values url.Values) ([]FieldSet, error) {
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

			var op Op
			var ok bool
			var not bool

			// ops can be chained, e.g. is.not:true, not.in:{1,2,3}
			subops := strings.Split(ops, OpSeparator)
			switch len(subops) {
			case 1:
				op, ok = ParseOp(subops[0])
				if !ok {
					return nil, fmt.Errorf("%w: %s", ErrUnknownOperator, v)
				}
			case 2:
				// Allow chaining of rule with not or is, e.g.
				// is.not:true, is:not:unknown
				// not.eq:john, not.in:{1,2,3}
				op1, ok := ParseOp(subops[0])
				if !ok {
					return nil, fmt.Errorf("%w: %s", ErrUnknownOperator, v)
				}

				op2, ok := ParseOp(subops[1])
				if !ok {
					return nil, fmt.Errorf("%w: %s", ErrUnknownOperator, v)
				}

				if !IsOpChainable(op1, op2) {
					return nil, fmt.Errorf("%w: %s", ErrBadOperator, v)
				}

				not = true

				switch op1 {
				case OpIs:
					op = op1
					// OpIs must have value: true, false, unknown or null.
					if !sqlIs(val) {
						return nil, fmt.Errorf("%w: %s", ErrInvalidIs, v)
					}
				case OpNot:
					op = op2
				default:
					return nil, fmt.Errorf("%w: %s", ErrBadOperator, v)
				}
			default:
				return nil, fmt.Errorf("%w: %s", ErrBadOperator, v)
			}

			if !rule.Has(op) {
				return nil, fmt.Errorf("%w: %s", ErrBadOperator, v)
			}

			cacheKey := fmt.Sprintf("%s:%s:%t", field, op, not)
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
				Not:      not,
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
			case op.Is(OpIn), col.IsArray:
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
