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
	ErrOperatorNotFound = errors.New("goql: op not found")
	ErrUnknownOperator  = errors.New("goql: unknown op")
	ErrInvalidNot       = errors.New("goql: invalid not")
	ErrInvalidIs        = errors.New("goql: invalid is")
	ErrInvalidArray     = errors.New("goql: invalid array")
	ErrInvalidBool      = errors.New("goql: invalid bool")
	ErrInvalidFloat     = errors.New("goql: invalid float")
	ErrInvalidInt       = errors.New("goql: invalid int")
)

type FieldSet struct {
	Name    string
	Value   any
	SQLType string
	Not     bool
	Format  string
	Tag     string
	Op      string
}

type Decoder struct {
	ops     map[string]Op
	columns map[string]Column
	parsers map[string]parserFn
	tag     string
}

func NewDecoder(v any) *Decoder {
	opsByField := make(map[string]Op)
	columns := StructToColumns(v, StructTag)
	parsers := make(map[string]parserFn)
	for k, v := range defaultParsers {
		parsers[k] = v
	}

	for name, col := range columns {
		// TODO: Allow registering custom types.
		if !IsPgType(col.SQLType) {
			panic("goql: not a sql type")
		}

		// By default, all datatypes are comparable.
		ops := OpsComparable

		if col.IsNull {
			ops |= OpsNull
		}

		if col.IsArray {
			ops |= OpsRange
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

func (d *Decoder) Decode(values url.Values) ([]FieldSet, error) {
	return Decode(d.ops, d.columns, d.parsers, values)
}

func Decode(ops map[string]Op, columns map[string]Column, parsers map[string]parserFn, values url.Values) ([]FieldSet, error) {
	used := make(map[string]bool)

	var sets []FieldSet

	// TODO: Handle unknown url.values?

	for field, rule := range ops {
		vs := values[field]
		col, ok := columns[field]
		if !ok {
			return nil, errors.New("goql: column not found")
		}

		for _, v := range vs {
			ops, val := split2(v, ValSeparator)

			var op Op
			var ok bool
			var not bool

			subops := strings.Split(ops, OpSeparator)
			switch len(subops) {
			case 1:
				op, ok = ParseOp(subops[0])
				if !ok {
					return nil, fmt.Errorf("%w: ops(%s) field(%s)", ErrUnknownOperator, ops, field)
				}
			case 2:
				// Allow chaining of rule with not or is, e.g.
				// is.not:true, is:not:unknown
				// not.eq:john, not.in:{1,2,3}
				op1, ok := ParseOp(subops[0])
				if !ok {
					return nil, fmt.Errorf("%w: ops(%s) field(%s)", ErrUnknownOperator, ops, field)
				}
				op2, ok := ParseOp(subops[1])
				if !ok {
					return nil, fmt.Errorf("%w: ops(%s) field(%s)", ErrUnknownOperator, ops, field)
				}

				if op1.Is(OpIs) && !op2.Is(OpNot) {
					return nil, fmt.Errorf("%w: %s %s", ErrUnknownOperator, ops, field)
				}
				switch op1 {
				case OpIs:
					not = true
					op = op1
				case OpNot:
					not = true
					op = op2
				default:
					return nil, fmt.Errorf("%w: %s %s", ErrUnknownOperator, ops, field)
				}
			default:
				return nil, fmt.Errorf("%w: %s %s", ErrUnknownOperator, ops, field)
			}

			if !rule.Has(op) {
				return nil, fmt.Errorf("%w: ops(%s) field(%s)", ErrUnknownOperator, ops, field)
			}

			switch {
			case op.Is(OpNot):
				not = true

				ops, val = split2(val, ":")
				op, ok = ParseOp(ops)
				if !ok {
					return nil, fmt.Errorf("%w: %s", ErrUnknownOperator, ops)
				}

				if OpsNot.Has(op) {
					return nil, fmt.Errorf("%w: %s", ErrInvalidNot, v)
				}

				if rule&op != op {
					return nil, fmt.Errorf("%w: %s", ErrUnknownOperator, ops)
				}
			case op.Is(OpIs):
				switch val {
				case
					"0", "f", "F", "false", "FALSE", "False",
					"1", "t", "T", "true", "TRUE", "True",
					"null", "NULL", "Null",
					"unk", "UNK", "Unk",
					"unknown", "UNKNOWN", "Unknown":
				default:
					return nil, fmt.Errorf("%w: %s", ErrInvalidIs, v)
				}
			}

			usedKey := fmt.Sprintf("%s:%s", field, v)
			if used[usedKey] {
				return nil, fmt.Errorf("%w: %s", ErrMultipleOperator, usedKey)
			}
			used[usedKey] = true

			// TODO: Handle parsing for all types based on parser.
			parser, ok := parsers[col.SQLType]
			if !ok {
				panic(fmt.Errorf("goql: parser not found: %+v", col))
			}

			switch {
			case op.Is(OpIn), col.IsArray:
				val, ok = Unquote(val, '{', '}')
				if !ok {
					return nil, fmt.Errorf("%w: missing parantheses: %s", ErrInvalidArray, val)
				}

				// Must be a list of strings.
				vals, err := splitString(val)
				if err != nil {
					return nil, err
				}
				var format []string
				if col.Format != "" {
					format = append(format, col.Format)
				}

				res, err := MapAny(vals, format, parser)
				if err != nil {
					return nil, err
				}

				sets = append(sets, FieldSet{
					Name:    field,
					Value:   res,
					SQLType: col.SQLType,
					Not:     not,
					Format:  col.Format,
					Tag:     col.Tag,
					Op:      op.String(),
				})
			default:
				sets = append(sets, FieldSet{
					Name:    field,
					Value:   val,
					SQLType: col.SQLType,
					Not:     not,
					Format:  col.Format,
					Tag:     col.Tag,
					Op:      op.String(),
				})
			}
		}
	}

	return sets, nil
}
