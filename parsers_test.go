package goql_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/alextanhongpin/goql"
	"github.com/google/go-cmp/cmp"
)

func TestParser(t *testing.T) {
	i := 100

	b, err := json.Marshal(map[string]any{
		"name":    "john",
		"age":     17,
		"married": false,
	})
	if err != nil {
		t.FailNow()
	}

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
		{string(b), "json.RawMessage", json.RawMessage(b)},
	}

	parsers := goql.NewParsers()

	for _, tt := range tests {
		name := fmt.Sprintf("parsing %s: %s", tt.typ, tt.inp)
		t.Run(name, func(t *testing.T) {
			parser, ok := parsers[tt.typ]
			if !ok {
				t.FailNow()
			}

			res, err := parser(tt.inp)
			if err != nil {
				t.Fatalf("parser error: %s", err)
			}
			if diff := cmp.Diff(tt.exp, res); diff != "" {
				t.Fatalf("expected+, got-: %s", diff)
			}
		})
	}

}
