package goql

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

var StructTag = "filter"

var (
	ErrMultipleOperators = errors.New("goql: multiple operators")
	ErrOperatorNotFound  = errors.New("goql: op not found")
	ErrUnknownOperator   = errors.New("goql: unknown op")
)

type FieldSet struct {
	typ RuleType
	op  Op
	key string
	val any
}

// name=eq:john&age=gt:13&married=eq:true
func Decode(values url.Values, v any) error {
	used := make(map[string]bool)

	var sets []FieldSet

	rules := Infer(v)
	for field, rule := range rules {
		vs := values[field]
		for _, v := range vs {
			opval := strings.SplitN(v, ":", 2)
			if len(opval) != 2 {
				return fmt.Errorf("%w: %s", ErrOperatorNotFound, opval)
			}

			ops, val := opval[0], opval[1]
			op, ok := rule.ops.Get(ops)
			if !ok {
				return fmt.Errorf("%w: %s", ErrUnknownOperator, ops)
			}
			usedKey := fmt.Sprintf("%s:%s", field, v)
			if used[usedKey] {
				return fmt.Errorf("%w: %s", ErrMultipleOperators, usedKey)
			}
			used[usedKey] = true

			switch rule.typ {
			case RuleTypeString:
				switch op {
				case OpIn:
					// Must be a list of strings.
					vals, err := splitString(val)
					if err != nil {
						return err
					}
					sets = append(sets, FieldSet{
						typ: rule.typ,
						op:  op,
						key: field,
						val: vals,
					})
				default:
					sets = append(sets, FieldSet{
						typ: rule.typ,
						op:  op,
						key: field,
						val: val,
					})

				}
			case RuleTypeInt:
				switch op {
				case OpIn:
					vals, err := ParseInts(strings.Split(val, string(QueryDelimiter)))
					if err != nil {
						return err
					}

					sets = append(sets, FieldSet{
						typ: rule.typ,
						op:  op,
						key: field,
						val: vals,
					})
				default:
					n, err := strconv.Atoi(val)
					if err != nil {
						return err
					}

					sets = append(sets, FieldSet{
						typ: rule.typ,
						op:  op,
						key: field,
						val: n,
					})
				}
			case RuleTypeBool:
				t, err := strconv.ParseBool(val)
				if err != nil {
					panic(err)
				}
				sets = append(sets, FieldSet{
					typ: rule.typ,
					op:  op,
					key: field,
					val: t,
				})
			}
		}
	}

	for _, set := range sets {
		fmt.Printf("%s %s %s %v\n", set.typ, set.op, set.key, set.val)
	}

	return nil
}

// Provide different infer, one for query, another for sql, another for mongo etc
func Infer(i any) map[string]Rule {
	rules := make(map[string]Rule)

	v := reflect.Indirect(reflect.ValueOf(i))
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag
		name := coalesce(tag.Get(StructTag), f.Name)

		// NULL Pointer
		switch v.Field(i).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			rules[name] = IntRule.Copy()
		case reflect.String:
			rules[name] = StringRule.Copy()
		case reflect.Float32, reflect.Float64:
			rules[name] = FloatRule.Copy()
		case reflect.Bool:
			rules[name] = BoolRule.Copy()
		case reflect.Struct:
			// Handle time,Time, json.RawMessage, []byte
		default:
		}
	}

	return rules
}

func coalesce[T comparable](ts ...T) (t T) {
	var z T
	for _, t = range ts {
		if t != z {
			return
		}
	}
	return z
}
