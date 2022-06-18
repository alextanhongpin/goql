package main

import (
	"fmt"
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
	v, err := url.ParseQuery(`name=eq:john&age=gt:13&married=eq:true&name=in:football,basketball,tennis&age=in:10,20,100`)
	if err != nil {
		panic(err)
	}

	dec := goql.NewDecoder(&User{})
	sets, err := dec.Decode(v)
	if err != nil {
		panic(err)
	}

	for _, set := range sets {
		fmt.Printf("%s %s %s %v\n", set.Typ, set.Op, set.Key, set.Val)
	}
}
