package goql

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const Separator = "."

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
	Op  Op
	Key string
	Val any
	Neg bool
}

type Decoder struct {
	rules map[string]Op
}

func NewDecoder(v any) *Decoder {
	return &Decoder{
		rules: NewRules(v),
	}
}

func (d *Decoder) Decode(values url.Values) ([]FieldSet, error) {
	return Decode(d.rules, values)
}

func Decode(rules map[string]Op, values url.Values) ([]FieldSet, error) {
	used := make(map[string]bool)

	var sets []FieldSet

	for field, rule := range rules {
		vs := values[field]

		for _, v := range vs {
			ops, val := split2(v, Separator)
			op, ok := ParseOp(ops)
			if !ok {
				return nil, fmt.Errorf("%w: %s", ErrUnknownOperator, ops)
			}

			if rule&op != op {
				return nil, fmt.Errorf("%w: %s", ErrUnknownOperator, ops)
			}

			var neg bool
			if op == OpNot {
				neg = true

				ops, val = split2(val, Separator)
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
			}

			if op == OpIs {
				if strings.HasPrefix(val, "not.") {
					neg = true
					val = strings.TrimPrefix(val, "not.")
				}
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

			switch rule {
			case RuleText:
				switch op {
				case OpIn:
					val, ok = Unquote(val, '(', ')')
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
						Neg: neg,
					})
				default:
					sets = append(sets, FieldSet{
						Typ: rule,
						Op:  op,
						Key: field,
						Val: val,
						Neg: neg,
					})

				}
			case RuleInt:
				switch op {
				case OpIn:
					val, ok = Unquote(val, '(', ')')
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
						Neg: neg,
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
						Neg: neg,
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
					Neg: neg,
				})
			}
		}
	}

	return sets, nil
}
