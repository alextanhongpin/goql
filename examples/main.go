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
	Name      string `sort:"true"`
	Age       int    `sort:"true"`
	Married   *bool
	Hobbies   []Hobby `q:"hobbies,type:[]string"`
	Birthday  *time.Time
	MarriedAt time.Time
	Height    *int
}

func main() {
	marriedAt := time.Now().UTC().Format(time.RFC3339)
	v := make(url.Values)
	v.Set("name.eq", "john")
	v.Add("name.in", "alpha")
	v.Add("name.in", "beta")
	v.Add("name.in", "gamma")
	v.Add("name.notin", "alice")
	v.Add("name.notin", "bob")
	v.Add("name.notin", "charles, junior")
	v.Set("age.gt", "13")
	v.Add("age.in", "13")
	v.Add("age.in", "17")
	v.Set("married.is", "true")
	v.Set("birthday.is", "null")
	v.Add("hobbies.eq", "football")
	v.Add("hobbies.eq", "music")
	v.Add("hobbies.eq", "drawing")
	v.Set("marriedAt.gt", marriedAt)
	v.Set("height.eq", "10")
	v.Add("and", "age.lt:10")
	v.Add("and", "age.gt:13")
	v.Add("and", "or.(name.eq:john,name.neq:jessie)")
	v.Add("or", "and.(height.isnot:null,height.lt:100)")
	v.Add("or", `and.(name.eq:"alice,ms",name.neq:bob)`)
	v.Add("sort_by", "name.asc")
	v.Add("sort_by", "age.desc")

	fmt.Println(v.Encode())

	dec := goql.NewDecoder[User]()
	filter, err := dec.Decode(v)
	if err != nil {
		panic(err)
	}

	var debug func(sets []goql.FieldSet, msg string, depth int)
	debug = func(sets []goql.FieldSet, msg string, depth int) {
		for i, set := range sets {
			tab := strings.Repeat("\t", depth)
			fmt.Printf("%s[%s]:%d. %s %s %#v\n", tab, msg, i+1, set.Name, set.Op, set.Value)

			debug(set.And, "AND", depth+1)
			debug(set.Or, "OR", depth+1)
		}
	}

	debug(filter.And, "AND", 0)
	debug(filter.Or, "OR", 0)
	for _, sort := range filter.Sort {
		fmt.Printf("sort: %+v\n", sort)
	}
}
