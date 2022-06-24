package goql

import (
	"reflect"
)

type Type struct {
	Name  string
	Array bool
	Null  bool
}

func (t *Type) Valid() bool {
	return t.Name != ""
}

func TypeOf(t reflect.Type) Type {
	res := Type{
		Name: t.String(),
	}

	switch t.Kind() {
	case reflect.Pointer:
		t = t.Elem()
		res.Null = true

	case reflect.Slice, reflect.Array:
		t = t.Elem()
		res.Array = true

		switch t.Kind() {
		case reflect.Pointer:
			t = t.Elem()
			res.Null = true
		}
	}

	return res
}
