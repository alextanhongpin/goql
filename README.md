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

## Operators

Basic datatypes (int, float, string, bool, time):

| op       | querystring                                   | sql                                     |
|----------|-----------------------------------------------|-----------------------------------------|
| eq       | `name.eq=john appleseed`                        | `name = 'john appleseed'`                 |
| neq      | `name.neq=john appleseed`                       | `name <> 'john appleseed'`                |
| lt       | `age.lt=17`                                     | `age < 17`                                |
| lte      | `age.lte=17`                                    | `age <= 17`                               |
| gt       | `age.gt=17`                                     | `age > 17`                                |
| gte      | `age.gte=17`                                    | `age >= 17`                               |
| like     | `title.like=programming%`                       | `title like 'programming%'`               |
| ilike    | `title.ilike=programming%`                      | `title ilike 'programming%'`              |
| notlike  | `title.notlike=programming%`                    | `title not like 'programming%'`           |
| notilike | `title.notilike=programming%`                   | `title not ilike 'programming%'`          |
| in       | `hobbies.in=programming&hobbies.in=music`       | `hobbies in ('programming', 'music')`     |
| notin    | `hobbies.notin=programming&hobbies.notin=music` | `hobbies not in ('programming', 'music')` |
| is       | `married.is=true`                               | `married is true`                         |
| isnot    | `married_at.isnot=null`                         | `married_at is not null`                  |


Some operators such as `IN`, `LIKE`, `ILIKE` and their negation `NOT` supports multiple values:

| op       | querystring                                       | sql                                                  |
|----------|---------------------------------------------------|------------------------------------------------------|
| eq       | `hobbies.eq=swimming&hobbies.eq=dancing`            | `hobbies = array['swimming', 'dancing']`               |
| neq      | `hobbies.neq=swimming&hobbies.neq=dancing`          | `hobbies <> array['swimming', 'dancing']`              |
| lt       | `scores.lt=50&scores.lt=100`                        | `scores < array[10, 100]`                              |
| lte      | `scores.lte=50&scores.lte=100`                      | `scores <= array[10, 100]`                             |
| gt       | `scores.gt=50&scores.gt=100`                        | `scores >= array[10, 100]`                             |
| gte      | `scores.gte=50&scores.gte=100`                      | `scores >= array[10, 100]`                             |


If the target type is an `array` [^1], then multiple values are accepted too:

| op       | querystring                                         | sql                                                    |
|----------|-----------------------------------------------------|--------------------------------------------------------|
| like     | `title.like=programming%&title.like=music%`         | `title like any(array['programming%', 'music%'])`      |
| ilike    | `title.ilike=programming%&title.ilike=music%`       | `title ilike any(array['programming%', 'music%'])`     |
| notlike  | `title.notlike=programming%&title.notlike=music%`   | `title not like all(array['programming%', 'music%'])`  |
| notilike | `title.notilike=programming%&title.notilike=music%` | `title not ilike all(array['programming%', 'music%'])` |
| in       | `hobbies.in=programming&hobbies.in=music`           | `hobbies in ('programming', 'music')`                  |
| notin    | `hobbies.notin=programming&hobbies.notin=music`     | `hobbies not in ('programming', 'music')`              |

## And/Or


`AND`/`OR` needs to be wrapped in open/close brackets:

| op  | querystring                                                         | sql                                                                           |
|-----|---------------------------------------------------------------------|-------------------------------------------------------------------------------|
| and | `and=(age.gt:13,age.lt:30,or.(name.ilike:alice%,name.notilike:bob%))` | `AND (age > 13 && age < 30 OR (name ilike 'alice%' and name not ilike 'bob%'))` |
| or  | `or=(height.isnot:null,height.gte:170)`                               | `OR (height is not null AND height >= 170)`                                     |
| or  | `or=(height.isnot:null)&or=(height.gte:170)`                          | `OR height is not null OR height >= 170`                                        |

## Limit/Offset


The default naming for the limit/offset in querystring is `limit` and `offset`. The query string name can changed by calling:

```go
dec.SetQueryLimitName("_limit")
dec.SetQueryOffsetName("_offset")
```

The default min/max for the limit is `1` and `20` respectively. To change the limit:


```go
dec.SetLimitRange(5, 100)
```

| op           | querystring        | sql                 |
|--------------|--------------------|---------------------|
| limit        | `limit=20`           | `LIMIT 20`            |
| offset       | `offset=42`          | `OFFSET 42`           |
| limit/offset | `limit=20&offset=42` | `LIMIT 20, OFFSET 42` |

## Sort

To enable sorting for a field, set the struct tag `sort:"true"`. E.g.

```go
type User struct {
	Age int `sort:"true"`
}
```

The `sort` struct name can be changed:

```go
dec.SetSortTag("sortable")
```

And the example above will then be:

```go
type User struct {
	Age int `sortable:"true"`
}
```

| op   | querystring             | sql                                            |
|------|-------------------------|------------------------------------------------|
| sort | `sort=age`                | `ORDER BY AGE ASC NULLSLAST`                     |
|      | `sort=age.asc`            | `ORDER BY age ASC NULLSLAST`                     |
|      | `sort=age.desc`           | `ORDER BY age DESC NULLSFIRST`                   |
|      | `sort=age.asc.nullsfirst` | `ORDER BY age ASC NULLSFIRST`                    |
|      | `sort=age.desc.nullslast` | `ORDER BY age DESC NULLSLAST`                    |
|      | `sort=id.desc&sort=age`   | `ORDER BY id DESC NULLSFIRST, age ASC NULLSLAST` |

## Tags


```go
type User struct {
	Age int `q:"age"`
}
```

The default tag name is `q`. To change it:


```go
dec.SetFilterTag("filter")
```

And the example above will be:

```go
type User struct {
	Age int `filter:"age"`
}
```

| field                           | tag                | description                                                                                                                                    |
|---------------------------------|--------------------|------------------------------------------------------------------------------------------------------------------------------------------------|
| name string `q:"name"`          | name               | modifies the query string name                                                                                                                 |
| MarriedAt time.Time `q:","`     | name               | if no name is specified, it defaults to lower common initialism of the name                                                                    |
| ID uuid.UUID `q:",type:uuid"`   | type:<your-type>   | specifies the type that would be used by the `parser`                                                                                          |
| IDs []string `q:",type:[]uuid"` | type:[]<your-type> | specifies an `array` type. `array` types have special operators                                                                                |
| ID *string `q:",type:*uuid"`    | type:*<your-type>  | specifies a `null` type. `null` types have special operators                                                                                   |
| ID string `q:",null"`           | null               | another approach of specifying `null` types                                                                                                    |
| ID string `q:",ops:eq,neq"`     | ops                | specifies the list of supported ops. In this example, only `id.eq=v` and `id.neq=v` is valid. This can be further overwritten by `dec.SetOps`. |


## Ops

To customize `ops` for a specific field, either set the struct tag `ops:<comma-separate-list-of-ops>`, or set it through the method `SetOps`:

```go
// To allow only `eq` and `neq` for the field `name`:
dec.SetOps("name", goql.OpEq | goql.OpNeq)
```

## Parsers

Query string parameters are string (or list of string). Parsers are responsible for parsing the string to the desired types that are either

1) inferred through the struct field's type through reflection
2) set at the struct tag through `type:"yourtype"`


Custom parsers can be registered as follow:

```go
	type Book struct {
		ID uuid.UUID `q:"id,type:uuid"` // Register a new type `uuid`.
	}

	id := uuid.New()

	v := make(url.Values)
	v.Set("id.eq", id.String())
	v.Add("id.in", id.String()) // Automatically handles conversion for a list of values.
	v.Add("id.in", id.String())

	dec := goql.NewDecoder[Book]()
	dec.SetParser("uuid", parseUUID)

	f, err := dec.Decode(v)
	if err != nil {
		t.Fatal(err)
	}
```

where `parseUUID` fulfills the `goql.ParseFn` method definition:

```go
func parseUUID(in string) (any, error) {
	return uuid.Parse(in)
}
```


## FAQ

> What if I need to filter some fields from `url.Values`?

Filter it manually before passing to `Decode(v url.Values)`.


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


## Reference

[^1]: [List of array operators](https://www.postgresql.org/docs/9.1/functions-array.html)
