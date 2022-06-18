package goql

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

var (
	ErrMultipleOperators = errors.New("goql: multiple operators")
	ErrOperatorNotFound  = errors.New("goql: op not found")
	ErrUnknownOperator   = errors.New("goql: unknown op")
)

type FieldSet struct {
	Typ RuleType
	Op  Op
	Key string
	Val any
}

type Decoder struct {
	rules map[string]Rule
}

func NewDecoder(v any) *Decoder {
	return &Decoder{
		rules: NewRules(v),
	}
}

func (d *Decoder) Decode(values url.Values) ([]FieldSet, error) {
	used := make(map[string]bool)

	var sets []FieldSet

	for field, rule := range d.rules {
		vs := values[field]
		for _, v := range vs {
			opval := strings.SplitN(v, ":", 2)
			if len(opval) != 2 {
				return nil, fmt.Errorf("%w: %s", ErrOperatorNotFound, opval)
			}

			ops, val := opval[0], opval[1]
			op, ok := rule.ops.Get(ops)
			if !ok {
				return nil, fmt.Errorf("%w: %s", ErrUnknownOperator, ops)
			}
			usedKey := fmt.Sprintf("%s:%s", field, v)
			if used[usedKey] {
				return nil, fmt.Errorf("%w: %s", ErrMultipleOperators, usedKey)
			}
			used[usedKey] = true

			switch rule.typ {
			case RuleTypeString:
				switch op {
				case OpIn:
					// Must be a list of strings.
					vals, err := splitString(val)
					if err != nil {
						return nil, err
					}
					sets = append(sets, FieldSet{
						Typ: rule.typ,
						Op:  op,
						Key: field,
						Val: vals,
					})
				default:
					sets = append(sets, FieldSet{
						Typ: rule.typ,
						Op:  op,
						Key: field,
						Val: val,
					})

				}
			case RuleTypeInt:
				switch op {
				case OpIn:
					vals, err := ParseInts(strings.Split(val, string(QueryDelimiter)))
					if err != nil {
						return nil, err
					}

					sets = append(sets, FieldSet{
						Typ: rule.typ,
						Op:  op,
						Key: field,
						Val: vals,
					})
				default:
					n, err := strconv.Atoi(val)
					if err != nil {
						return nil, err
					}

					sets = append(sets, FieldSet{
						Typ: rule.typ,
						Op:  op,
						Key: field,
						Val: n,
					})
				}
			case RuleTypeBool:
				t, err := strconv.ParseBool(val)
				if err != nil {
					return nil, err
				}
				sets = append(sets, FieldSet{
					Typ: rule.typ,
					Op:  op,
					Key: field,
					Val: t,
				})
			}
		}
	}

	return sets, nil
}
