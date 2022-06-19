package main

import (
	"fmt"
	"net/url"

	"github.com/alextanhongpin/goql"
)

type User struct {
	Name    string   `sql:"name" sort:"true"` // query=, sql=?
	Age     int      `sql:"age" sort:"true"`
	Married bool     `sql:"married" sort:"true"`
	Hobbies []string `sql:"hobbies" sort:"true"`
	// Avoid using reflect, instead, let user specify the type in the tag...
	//Birthday time.Time `sql:"birthday" sort:"true"`
	Birthday string `sql:"birthday" sort:"true"`
}

func main() {
	v, err := url.ParseQuery(`name=eq.john&age=gt.13&married=is.true&name=in.(football,basketball,tennis)&age=in.(10,20,100)&name=not.in.(alice,bob)&birthday=is.null&hobbies=is.not.null`)
	if err != nil {
		panic(err)
	}
	fmt.Println(v)

	dec := goql.NewDecoder(&User{})
	sets, err := dec.Decode(v)
	if err != nil {
		panic(err)
	}

	fmt.Println("sets:", sets)
	for _, set := range sets {
		fmt.Printf("%s %s %s %v neg: %t\n", set.Typ, set.Op, set.Key, set.Val, set.Neg)
	}
}
