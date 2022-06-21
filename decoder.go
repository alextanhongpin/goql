package goql

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
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
	Typ Op
	Key string
	Val any

	Op Op

	//Type string
	Not bool // sql NOT
	//Format string // Additional information for parsing data types, e.g. date:20060102 for formatting text date based on yyyymmdd
	//Tag    string
	//Column string
	//Value  any
}

type Decoder struct {
	rules   map[string]Op
	columns map[string]column
	tag     string
}

func NewDecoder(v any) *Decoder {
	rules := make(map[string]Op)
	columns := structToColumns(v, StructTag)
	for name, col := range columns {
		ops, ok := opsByPgType[col.sqlType]
		if !ok {
			ops = RuleWhere
		}
		if col.null {
			ops |= RuleNull
		}

		rules[name] = ops
	}

	return &Decoder{
		columns: columns,
		rules:   rules,
		tag:     StructTag,
		// parsers: func(string) any, error
	}
}

func (d *Decoder) SetStructTag(tag string) {
	if tag == "" {
		panic("tag cannot be empty")
	}

	d.tag = tag
}

func (d *Decoder) SetFieldOps(opsByField map[string]Op) {
	d.rules = opsByField
}

func (d *Decoder) Decode(values url.Values) ([]FieldSet, error) {
	return Decode(d.rules, values)
}

func Decode(rules map[string]Op, values url.Values) ([]FieldSet, error) {
	used := make(map[string]bool)

	var sets []FieldSet

	// TODO: Handle unknown url.values?

	for field, rule := range rules {
		vs := values[field]

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

				if RuleNegate&op != op {
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
			switch rule {
			case RuleText:
				switch op {
				case OpIn:
					val, ok = Unquote(val, '{', '}')
					if !ok {
						return nil, fmt.Errorf("%w: missing parantheses: %s", ErrInvalidArray, val)
					}

					// Must be a list of strings.
					vals, err := splitString(val)
					if err != nil {
						return nil, err
					}
					sets = append(sets, FieldSet{
						Typ: rule,
						Op:  op,
						Key: field,
						Val: vals,
						Not: not,
					})
				default:
					sets = append(sets, FieldSet{
						Typ: rule,
						Op:  op,
						Key: field,
						Val: val,
						Not: not,
					})

				}
			case RuleInt:
				switch op {
				case OpIn:
					val, ok = Unquote(val, '{', '}')
					if !ok {
						return nil, fmt.Errorf("%w: missing parantheses: %s", ErrInvalidArray, val)
					}

					vals, err := ParseInts(strings.Split(val, string(QueryDelimiter)))
					if err != nil {
						return nil, err
					}

					sets = append(sets, FieldSet{
						Typ: rule,
						Op:  op,
						Key: field,
						Val: vals,
						Not: not,
					})
				default:
					n, err := strconv.Atoi(val)
					if err != nil {
						return nil, fmt.Errorf("%w: %s", ErrInvalidInt, val)
					}

					sets = append(sets, FieldSet{
						Typ: rule,
						Op:  op,
						Key: field,
						Val: n,
						Not: not,
					})
				}
			case RuleBool:
				t, err := strconv.ParseBool(val)
				if err != nil {
					return nil, ErrInvalidBool
				}
				sets = append(sets, FieldSet{
					Typ: rule,
					Op:  op,
					Key: field,
					Val: t,
					Not: not,
				})
			}
		}
	}

	return sets, nil
}
