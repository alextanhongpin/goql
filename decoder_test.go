package goql_test

import (
	"net/url"
	"testing"

	"github.com/alextanhongpin/goql"
	"github.com/google/uuid"
)

type User struct {
	Name    string `sql:"name"`
	Age     int    `sql:"age"`
	Married bool
}

type Book struct {
	ID uuid.UUID `sql:"id,type:uuid"`
}

func TestDecoder(t *testing.T) {
	dec := goql.NewDecoder[User]()
	dec.SetStructTag("sql")
	dec.SetFieldOps(map[string]goql.Op{
		"name":    goql.OpEq, // Only allow equality comparison.
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

func TestDecoderCustomParser(t *testing.T) {
	dec := goql.NewDecoder[Book]()
	dec.SetStructTag("sql")
	dec.SetParsers(map[string]goql.ParserFn{
		// Register type
		"uuid": parseUUID,
	})

	v, err := url.ParseQuery(`id=eq:` + uuid.Nil.String())
	if err != nil {
		t.FailNow()
	}

	fieldSets, err := dec.Decode(v)
	if err != nil {
		t.Fatal(err)
	}

	if exp, got := uuid.Nil, fieldSets[0].Value; exp != got {
		t.Fatalf("expected %v, got %v", exp, got)
	}

	t.Log(fieldSets)
}

func parseUUID(in string, format ...string) (any, error) {
	return uuid.Parse(in)
}
