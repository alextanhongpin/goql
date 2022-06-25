package goql_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/alextanhongpin/goql"
)

func TestParser(t *testing.T) {
	i := 100
	tests := []struct {
		inp string
		typ string
		exp any
	}{
		{"100", "int", int(i)},
		{"null", "*int", (*int)(nil)},
		{"100", "*int", &i},
		{"100", "int64", int64(i)},
		{"true", "bool", true},
		{"null", "*bool", (*bool)(nil)},
	}

	parsers := goql.NewParsers()

	for _, tt := range tests {
		name := fmt.Sprintf("%s, %s => %v", tt.inp, tt.typ, tt.exp)
		t.Run(name, func(t *testing.T) {
			parser, ok := parsers[tt.typ]
			if !ok {
				t.FailNow()
			}

			res, err := parser(tt.inp)
			if err != nil {
				t.Fatalf("parser error: %s", err)
			}
			if exp, got := tt.exp, res; !reflect.DeepEqual(exp, got) {
				t.Fatalf("expected %v, got %v", exp, got)
			}
		})
	}

}
