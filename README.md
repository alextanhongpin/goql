# goql [![Go Reference](https://pkg.go.dev/badge/github.com/alextanhongpin/goql.svg)](https://pkg.go.dev/github.com/alextanhongpin/goql)

Parse query string to Postgres SQL.

Requires `go 1.18+`.

## Features

- customizable field names, field ops, query string fields, as well as struct tags
- supports many operators used by Postgres
- handles conjunction (and/or)
- handles type conversion from query string to designated struct field's type
- handles limit/offset
- handles sorting
- handles parsing slices for array operators like `IN`, `NOT IN`, `LIKE`, `ILIKE`

## Installation

```bash
go get github.com/alextanhongpin/goql
```

## Basic example

```go
package main

import (
	"fmt"
	"net/url"

	"github.com/alextanhongpin/goql"
)

type Book struct {
	Author      string
	Title       string
	PublishYear int `q:"publish_year" sort:"true"`
}

func main() {
	// Register a new decoder for the type Book.
	// Information are extracted from the individual struct
	// fields, as well as the struct tag.
	dec := goql.NewDecoder[Book]()

	/*
		Query: Find books that are published by 'Robert
		Greene' which has the keyword 'law' and 'master'
		published after 2010. The results should be ordered
		descending nulls last, limited to 10.


		SELECT *
		FROM books
		WHERE author = 'Robert Greene'
		AND publish_year > 2010
		AND title ilike any(array['law', 'master'])
		ORDER BY publish_year desc nullslast
		LIMIT 10
	*/

	v := make(url.Values)
	v.Set("author.eq", "Robert Greene")
	v.Set("publish_year.gt", "2010")
	v.Add("title.ilike", "law")
	v.Add("title.ilike", "master")
	v.Add("sort_by", "publish_year.desc.nullslast")
	v.Add("limit", "10")

	f, err := dec.Decode(v)
	if err != nil {
		panic(err)
	}

	for _, and := range f.And {
		fmt.Println("and:", and)
	}

	fmt.Println("limit:", *f.Limit)
	for _, sort := range f.Sort {
		fmt.Println("sort:", sort)
	}
}
```

Output:

```
and: author eq "Robert Greene"
and: publish_year gt 2010
and: title ilike []interface {}{"law", "master"}
limit: 10
sort: publish_year desc nullslast
```

## FAQ

> What if I need to filter some fields from `url.Values`?

Filter it manually before passing to `NewDecoder(values)`.


> What if I need to add validation to the values?

Create a new type (aka `value object), and add a parser for that type which includes validation. Parsers are type-specific.

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/alextanhongpin/goql"
)

var ErrInvalidEmail = errors.New("bad email format")

type Email string

func (e Email) Validate() error {
	if !strings.Contains(string(e), "@") {
		return ErrInvalidEmail
	}

	return nil
}

type User struct {
	Email Email `q:",type:email"` // Register a new type "email"
}

func main() {
	dec := goql.NewDecoder[User]()
	dec.SetParser("email", parseEmail)

	v := make(url.Values)
	v.Set("email.eq", "bad email") // Register a parser for type "email"
	_, err := dec.Decode(v)
	fmt.Println(err)
	fmt.Println(errors.Is(err, ErrInvalidEmail))
}

func parseEmail(in string) (any, error) {
	email := Email(in)
	return email, email.Validate()
}
```


> Why is there no support for keyset/cursor pagination, only limit/offset?

Because you can construct the cursor pagination directly.
