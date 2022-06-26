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
