package goql_test

import (
	"net/url"
	"testing"

	"github.com/alextanhongpin/goql"
)

type User struct {
	Name    string `sql:"name"`
	Age     int    `sql:"age"`
	Married bool
}

func TestDecoder(t *testing.T) {
	dec := goql.NewDecoder[User]()
	dec.SetStructTag("sql")
	dec.SetFieldOps(map[string]goql.Op{
		"name":    goql.OpEq,
		"age":     goql.OpsComparable,
		"Married": goql.OpsNull,
	})

	v, err := url.ParseQuery(`name=eq:hello&age=eq:10&Married=is:true`)
	if err != nil {
		t.FailNow()
	}

	fieldSets, err := dec.Decode(v)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fieldSets)
}
