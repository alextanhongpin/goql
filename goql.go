package goql

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

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
func Parser(values url.Values, v any) error {
	var sets []FieldSet

	rules := Infer(v)
	for field, rule := range rules {
		vs := values[field]
		if len(vs) > 1 {
			return fmt.Errorf("%w: %s", ErrMultipleOperators, vs[0])
		}

		v := vs[0]
		opval := strings.SplitN(v, ":", 2)
		if len(opval) != 2 {
			return fmt.Errorf("%w: %s", ErrOperatorNotFound, opval)
		}

		ops, val := opval[0], opval[1]
		op, ok := rule.ops.Get(ops)
		if !ok {
			return fmt.Errorf("%w: %s", ErrUnknownOperator, ops)
		}

		switch rule.typ {
		case RuleTypeString:
			sets = append(sets, FieldSet{
				typ: rule.typ,
				op:  op,
				key: field,
				val: val,
			})
		case RuleTypeInt:
			n, err := strconv.Atoi(val)
			if err != nil {
				panic(err)
			}
			sets = append(sets, FieldSet{
				typ: rule.typ,
				op:  op,
				key: field,
				val: n,
			})
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
		name := coalesce(tag.Get("filter"), f.Name)

		// NULL Pointer
		switch v.Field(i).Kind() {
		case reflect.Int:
			rules[name] = IntRule.Copy()
		case reflect.String:
			rules[name] = StringRule.Copy()
		case reflect.Bool:
			rules[name] = BoolRule.Copy()
		case reflect.Struct:
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
