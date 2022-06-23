package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/alextanhongpin/goql"
)

type Hobby string

type User struct {
	// TODO: Handle sortable
	Name      string     `sql:"name" sort:"true"` // query=, sql=?
	Age       int        `sql:"age" sort:"true"`
	Married   *bool      `sql:"married" sort:"true"`
	Hobbies   []Hobby    `sql:"hobbies" sort:"true"`
	Birthday  *time.Time `sql:"birthday" sort:"true"`
	MarriedAt time.Time  `sql:"marriedAt" sort:"true"`
}

func main() {
	marriedAt := time.Now().UTC().Format(time.RFC3339)
	v, err := url.ParseQuery(`name=eq:john&age=gt:13&married=is:true&name=in:{football,basketball,tennis}&age=in:{10,20,100}&name=not.in:{alice,bob}&birthday=is:null&hobbies=eq:{1,2,3}&marriedAt=gt:` + marriedAt)
	if err != nil {
		panic(err)
	}
	fmt.Println("query:", v)

	dec := goql.NewDecoder(&User{})
	sets, err := dec.Decode(v)
	if err != nil {
		panic(err)
	}

	fmt.Println("sets:", sets)
	for _, set := range sets {
		fmt.Printf("%+v\n", set)
	}
}
