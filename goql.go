package goql

import (
	"reflect"
)

const StructTag = "sql"

// Provide different infer, one for query, another for sql, another for mongo etc
func NewRules(i any) map[string]Op {
	rules := make(map[string]Op)

	v := reflect.Indirect(reflect.ValueOf(i))
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag
		name := coalesce(tag.Get(StructTag), f.Name)

		//if v.Field(i).Kind() == reflect.Pointer {
		//v.Field(i).Elem().Interface()
		//}

		// NULL Pointer
		switch v.Field(i).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			rules[name] = RuleInt
		case reflect.String:
			rules[name] = RuleText
		case reflect.Float32, reflect.Float64:
			rules[name] = RuleFloat
		case reflect.Bool:
			rules[name] = RuleBool
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
