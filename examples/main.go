package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/alextanhongpin/goql"
)

// type alias needs to be registered manually
type Hobby string

type User struct {
	// TODO: Handle sortable
	Name      string     `q:"name" sort:"true"` // query=, q=?
	Age       int        `q:"age" sort:"true"`
	Married   *bool      `q:"married" sort:"true"`
	Hobbies   []Hobby    `q:"hobbies,type:[]string" sort:"true"`
	Birthday  *time.Time `q:"birthday" sort:"true"`
	MarriedAt time.Time  `q:"marriedAt" sort:"true"`
}

func main() {
	marriedAt := time.Now().UTC().Format(time.RFC3339)
	v, err := url.ParseQuery(`name=eq:john&age=gt:13&married=is:true&name=in:{football,basketball,tennis}&age=in:{10,20,100}&name=notin:{alice,bob}&birthday=is:null&hobbies=eq:{1,2,3}&marriedAt=gt:` + marriedAt)
	if err != nil {
		panic(err)
	}
	fmt.Println("query:", v)

	dec := goql.NewDecoder[User]()
	sets, err := dec.Decode(v)
	if err != nil {
		panic(err)
	}

	fmt.Println("sets:", sets)
	for i, set := range sets {
		fmt.Printf("%d. %+v\n\n", i+1, set)
	}
}
