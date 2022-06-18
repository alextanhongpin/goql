package goql

import (
	"encoding/csv"
	"fmt"
	"net/url"
	"strings"
)

func main() {
	var s []string
	s = append(s, `hi, "John",this`)
	s = append(s, `hi, Jessie`)
	fmt.Println(len(strings.Split(strings.Join(s, ","), ",")))

	j, err := joinQuote(s, ',')
	if err != nil {
		panic(err)
	}
	fmt.Println("joined:", j)
	x, err := splitQuote(j, ',')
	if err != nil {
		panic(err)
	}
	fmt.Println("splitted:", x, len(x))

	v := make(url.Values)
	for _, xx := range x {
		v.Add("name", xx)
	}
	v.Add("test", strings.TrimSpace(j))

	fmt.Println(v.Encode())

	q, err := url.ParseQuery(v.Encode())
	if err != nil {
		panic(err)
	}
	fmt.Println(q)
}

func joinQuote(ss []string, delimiter rune) (string, error) {
	var sb strings.Builder
	w := csv.NewWriter(&sb)
	w.Comma = delimiter
	w.UseCRLF = false

	if err := w.Write(ss); err != nil {
		return "", err
	}

	w.Flush()

	return sb.String(), nil
}

func splitQuote(s string, delimiter rune) ([]string, error) {
	r := csv.NewReader(strings.NewReader(s))
	r.Comma = delimiter
	return r.Read()
}
