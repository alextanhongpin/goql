package goql_test

import (
	"database/sql"
	"errors"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/alextanhongpin/goql"
	"github.com/google/uuid"
)

func TestDecoderCustomStructTag(t *testing.T) {
	type User struct {
		Name    string `sql:"name"`
		Age     *int   `sql:"age"`
		Married bool
	}

	v, err := url.ParseQuery(`name.eq=hello&age.eq=10&married.is=true`)
	if err != nil {
		t.FailNow()
	}

	dec := goql.NewDecoder[User]()
	dec.SetFilterTag("sql")

	f, err := dec.Decode(v)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(f.And)
}

func TestDecoderCustomParser(t *testing.T) {
	type Book struct {
		ID uuid.UUID `q:"id,type:uuid"`
	}

	id := uuid.New()

	v := make(url.Values)
	v.Set("id.eq", id.String())

	dec := goql.NewDecoder[Book]()
	err := dec.Validate()
	if !errors.Is(err, goql.ErrUnknownParser) {
		t.FailNow()
	}
	t.Logf("validateError: %s", err)

	dec.SetParser("uuid", parseUUID)
	err = dec.Validate()
	if err != nil {
		t.Fatal(err)
	}

	filter, err := dec.Decode(v)
	if err != nil {
		t.Fatal(err)
	}

	if exp, got := id, filter.And[0].Value; exp != got {
		t.Fatalf("expected %v, got %v", exp, got)
	}

	t.Log(filter)
}

func parseUUID(in string) (any, error) {
	return uuid.Parse(in)
}

func TestDecoderSetOps(t *testing.T) {
	type User struct {
		Name string `q:"name,ops:eq"`
	}

	dec := goql.NewDecoder[User]()
	dec.SetOps("name", goql.OpNeq)

	t.Run("valid ops", func(t *testing.T) {
		v := make(url.Values)
		v.Set("name.eq", "hello")

		_, err := dec.Decode(v)
		if exp, got := goql.ErrUnknownOperator, err; !errors.Is(err, exp) {
			t.Fatalf("expected %v, got %v", exp, got)
		}
	})

	t.Run("invalid ops", func(t *testing.T) {
		v := make(url.Values)
		v.Set("name.neq", "hello")

		f, err := dec.Decode(v)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if exp, got := "hello", f.And[0].Value; !reflect.DeepEqual(exp, got) {
			t.Fatalf("expected %v, got %v", exp, got)
		}
	})
}

func TestDecoderTagOps(t *testing.T) {
	type User struct {
		Name string `q:"name,ops:eq"`
	}

	t.Run("valid ops", func(t *testing.T) {
		v := make(url.Values)
		v.Set("name.eq", "hello world")

		f, err := goql.NewDecoder[User]().Decode(v)
		if err != nil {
			t.Fatal(err)
		}

		if exp, got := "hello world", f.And[0].Value; !reflect.DeepEqual(exp, got) {
			t.Fatalf("expected %v, got %v", exp, got)
		}
	})

	t.Run("invalid ops", func(t *testing.T) {
		v := make(url.Values)
		v.Set("name.neq", "hello")

		_, err := goql.NewDecoder[User]().Decode(v)
		if err == nil {
			t.FailNow()
		}

		if exp, got := true, errors.Is(err, goql.ErrUnknownOperator); exp != got {
			t.Fatalf("expected %v, got %v", exp, err)
		}
	})
}

func TestDecoderNullTime(t *testing.T) {
	type User struct {
		MarriedAt sql.NullTime `q:"marriedAt,ops:is,gt"`
	}

	t.Run("null time", func(t *testing.T) {
		u := make(url.Values)
		u.Set("marriedAt.is", "null")

		dec := goql.NewDecoder[User]()
		dec.SetParser("sql.NullTime", parseSQLNullTime)
		if err := dec.Validate(); err != nil {
			t.Fatal(err)
		}

		f, err := dec.Decode(u)
		if err != nil {
			t.Fatal(err)
		}

		var exp sql.NullTime
		if got := f.And[0].Value; exp != got {
			t.Fatalf("expected %v, got %v", exp, got)
		}
	})

	t.Run("nonnull time", func(t *testing.T) {
		now := time.Now()

		u := make(url.Values)
		u.Set("marriedAt.gt", now.Format(time.RFC3339))

		dec := goql.NewDecoder[User]()
		dec.SetParser("sql.NullTime", parseSQLNullTime)
		if err := dec.Validate(); err != nil {
			t.Fatal(err)
		}

		f, err := dec.Decode(u)
		if err != nil {
			t.Fatal(err)
		}

		now, err = time.Parse(time.RFC3339, now.Format(time.RFC3339))
		if err != nil {
			t.Fatal(err)
		}

		exp := sql.NullTime{
			Valid: true,
			Time:  now,
		}
		if got := f.And[0].Value; exp != got {
			t.Fatalf("expected %v, got %v", exp, got)
		}
	})
}

func parseSQLNullTime(in string) (any, error) {
	t, err := goql.ParseStringPointer[time.Time](in)
	if err != nil {
		return nil, err
	}

	tp, ok := t.(*time.Time)
	if !ok || tp == nil {
		return sql.NullTime{}, nil
	}

	return sql.NullTime{Time: *tp, Valid: true}, nil
}

func TestDecodeLimit(t *testing.T) {
	type User struct {
		ID string
	}

	dec := goql.NewDecoder[User]().
		SetLimitRange(5, 25).         // Default is 1 to 20.
		SetQueryLimitName("_limit").  // Default is "limit".
		SetQueryOffsetName("_offset") // Default is "offset".

	t.Run("not set", func(t *testing.T) {
		u := make(url.Values)

		f, err := dec.Decode(u)
		if err != nil {
			t.Fatal(err)
		}

		if exp, got := (*int)(nil), f.Limit; exp != got {
			t.Fatalf("expected %v, got %v", exp, got)
		}

		if exp, got := (*int)(nil), f.Offset; exp != got {
			t.Fatalf("expected %v, got %v", exp, got)
		}
	})

	t.Run("above limit", func(t *testing.T) {
		u := make(url.Values)
		u.Set("_limit", "50")
		u.Set("_offset", "30")

		f, err := dec.Decode(u)
		if err != nil {
			t.Fatal(err)
		}

		if exp, got := 25, *f.Limit; exp != got {
			t.Fatalf("expected %v, got %v", exp, got)
		}

		if exp, got := 30, *f.Offset; exp != got {
			t.Fatalf("expected %v, got %v", exp, got)
		}
	})

	t.Run("below limit", func(t *testing.T) {
		u := make(url.Values)
		u.Set("_limit", "-10")
		u.Set("_offset", "-20")

		f, err := dec.Decode(u)
		if err != nil {
			t.Fatal(err)
		}

		if exp, got := 5, *f.Limit; exp != got {
			t.Fatalf("expected %v, got %v", exp, got)
		}

		if exp, got := 0, *f.Offset; exp != got {
			t.Fatalf("expected %v, got %v", exp, got)
		}
	})
}

func TestSort(t *testing.T) {
	type User struct {
		ID   int    `sortable:"true"`
		Name string `sortable:"true"`
	}

	dec := goql.NewDecoder[User]().
		SetSortTag("sortable").      // Default is "sort".
		SetQuerySortName("_sort_by") // Default is "sort_by".

	t.Run("not set", func(t *testing.T) {
		u := make(url.Values)

		f, err := dec.Decode(u)
		if err != nil {
			t.Fatal(err)
		}

		if exp, got := 0, len(f.Sort); exp != got {
			t.Fatalf("expected %v, got %v", exp, got)
		}
	})

	t.Run("sort multiple", func(t *testing.T) {
		u := make(url.Values)
		u.Add("_sort_by", "id.desc.nullslast")
		u.Add("_sort_by", "name.asc.nullsfirst")

		f, err := dec.Decode(u)
		if err != nil {
			t.Fatal(err)
		}

		byID := goql.Order{
			Field:     "id",
			Direction: "desc",
			Option:    "nullslast",
		}
		byName := goql.Order{
			Field:     "name",
			Direction: "asc",
			Option:    "nullsfirst",
		}

		if exp, got := byID, f.Sort[0]; exp != got {
			t.Fatalf("expected %v, got %v", exp, got)
		}

		if exp, got := byName, f.Sort[1]; exp != got {
			t.Fatalf("expected %v, got %v", exp, got)
		}
	})
}

var ErrInvalidEmail = errors.New("bad email format")

type Email string

func (e Email) Validate() error {
	if !strings.Contains(string(e), "@") {
		return ErrInvalidEmail
	}

	return nil
}

type User struct {
	Email Email `q:",type:email"` // Register custom type.
}

func TestParserValidator(t *testing.T) {
	dec := goql.NewDecoder[User]()
	dec.SetParser("email", parseEmail)

	v := make(url.Values)
	v.Set("email.eq", "bad email")
	_, err := dec.Decode(v)
	if err == nil {
		t.FailNow()
	}

	if !errors.Is(err, ErrInvalidEmail) {
		t.Fatalf("expected %v, got %v", ErrInvalidEmail, err)
	}

	t.Logf("validator: %v", err)
}

func parseEmail(in string) (any, error) {
	email := Email(in)
	return email, email.Validate()
}
