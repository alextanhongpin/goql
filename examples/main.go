package main

import (
	"net/url"
	"time"

	"github.com/alextanhongpin/goql"
)

type User struct {
	Name     string    `filter:"name" sort:"true"` // query=, sql=?
	Age      int       `filter:"age" sort:"true"`
	Married  bool      `filter:"married" sort:"true"`
	Hobbies  []string  `filter:"hobbies" sort:"true"`
	Birthday time.Time `filter:"birthday" sort:"true"`
}

func main() {
	v, err := url.ParseQuery(`name=eq:john&age=gt:13&married=eq:true`)
	if err != nil {
		panic(err)
	}
	goql.Parser(v, &User{})
}
