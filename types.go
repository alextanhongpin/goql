package goql

import (
	"reflect"
)

// Type represents the type of the field.
type Type struct {
	Name  string
	Null  bool
	Array bool
}

// Valid returns true if the type is not empty.
func (t *Type) Valid() bool {
	return t.Name != ""
}

/*
TypeOf handles conversion for the following:

type Struct {
	field1 *Type
	field2 Type
	field3 []Type
	field4 []*Type
	field5 *[]Type
	field6 *[]*Type
}
*/
func TypeOf(t reflect.Type) Type {
	res := Type{
		Name: t.String(),
	}

	switch t.Kind() {
	case reflect.Pointer:
		t = t.Elem()
		res.Null = true
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
