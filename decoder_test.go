package goql_test

import (
	"errors"
	"net/url"
	"testing"

	"github.com/alextanhongpin/goql"
	"github.com/google/uuid"
)

func TestDecoderCustomStructTag(t *testing.T) {
	type User struct {
		Name    string `sql:"name"`
		Age     int    `sql:"age"`
		Married bool
	}

	dec, err := goql.NewDecoder[User]()
	if err != nil {
		t.Fatalf("error constructing new decoder: %v", err)
	}

	if err := dec.SetStructTag("sql"); err != nil {
		t.Fatalf("error setting struct tag: %v", err)
	}

	v, err := url.ParseQuery(`name=eq:hello&age=eq:10&married=is:true`)
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
	type Book struct {
		ID uuid.UUID `sql:"id,type:uuid"`
	}

	dec, err := goql.NewDecoder[Book]()
	if err != nil {
		t.Fatalf("error constructing new decoder: %v", err)
	}

	dec.SetStructTag("sql")
	dec.SetParsers(map[string]goql.ParserFn{
		// Register type
		"uuid": parseUUID,
	})

	id := uuid.New()
	v, err := url.ParseQuery(`id=eq:` + id.String())
	if err != nil {
		t.FailNow()
	}

	fieldSets, err := dec.Decode(v)
	if err != nil {
		t.Fatal(err)
	}

	if exp, got := id, fieldSets[0].Value; exp != got {
		t.Fatalf("expected %v, got %v", exp, got)
	}

	t.Log(fieldSets)
}

func parseUUID(in string) (any, error) {
	return uuid.Parse(in)
}

func TestDecoderTagOps(t *testing.T) {
	type User struct {
		Name string `q:"name,ops:eq"`
	}

	dec, err := goql.NewDecoder[User]()
	if err != nil {
		t.Fatalf("error constructing new decoder: %v", err)
	}

	t.Run("valid ops", func(t *testing.T) {
		v, err := url.ParseQuery(`name=eq:hello`)
		if err != nil {
			t.FailNow()
		}

		fieldSets, err := dec.Decode(v)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(fieldSets)
	})

	t.Run("invalid ops", func(t *testing.T) {
		v, err := url.ParseQuery(`name=neq:hello`)
		if err != nil {
			t.FailNow()
		}

		_, err = dec.Decode(v)
		if err == nil {
			t.FailNow()
		}

		if exp, got := true, errors.Is(err, goql.ErrUnknownOperator); exp != got {
			t.Fatalf("expected %v, got %v", exp, err)
		}
	})
}
