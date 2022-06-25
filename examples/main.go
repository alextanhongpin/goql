package main

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/alextanhongpin/goql"
)

// type alias needs to be registered manually
type Hobby string

type User struct {
	// TODO: Handle sortable
	Name      string `sort:"true"` // query=, q=?
	Age       int    `sort:"true"`
	Married   *bool
	Hobbies   []Hobby `q:"hobbies,type:[]string"`
	Birthday  *time.Time
	MarriedAt time.Time
	Height    *int
}

func main() {
	marriedAt := time.Now().UTC().Format(time.RFC3339)
	v := url.Values{
		"name":      []string{"eq:john", "in:{football,basketball,tennis}", `notin:{alice,bob,"charles, junior"}`},
		"age":       []string{"gt:13", "in:{10,20,100}"},
		"married":   []string{"is:true"},
		"birthday":  []string{"is:null"},
		"hobbies":   []string{"eq:{1,2,3,}"},
		"marriedAt": []string{"gt:" + marriedAt},
		"height":    []string{"eq:10"},
		"and":       []string{"(age.lt:10,age.gt:13,or.(name.eq:john,name.neq:jessie))", "(or.(height.isnot:null,height.lt:100))"},
		"or":        []string{`(name.eq:"alice,ms",name.neq:bob)`},
		"sort_by":   []string{"name.asc", "age.desc"},
	}
	fmt.Println(v.Encode())

	dec, err := goql.NewDecoder[User]()
	if err != nil {
		panic(err)
	}

	filter, err := dec.Decode(v)
	if err != nil {
		panic(err)
	}

	var debug func(sets []goql.FieldSet, depth int)
	debug = func(sets []goql.FieldSet, depth int) {
		for i, set := range sets {
			tab := strings.Repeat("\t", depth)
			fmt.Printf("%s%d. %s %s %#v\n", tab, i+1, set.Name, set.Op, set.Value)

			debug(set.And, depth+1)
			debug(set.Or, depth+1)
		}
	}

	debug(filter.And, 0)
	debug(filter.Or, 0)
	for _, sort := range filter.Sort {
		fmt.Printf("sort: %+v\n", sort)
	}
}
