package goql_test

import (
	"fmt"
	"testing"

	"github.com/alextanhongpin/goql"
)

func TestUnquote(t *testing.T) {
	tests := []struct {
		exp   string
		str   string
		valid bool
	}{
		{exp: "", str: "", valid: false},
		{exp: "(", str: "(", valid: false},
		{exp: "((", str: "((", valid: false},
		{exp: "", str: "()", valid: true},
		{exp: "hello world", str: "(hello world)", valid: true},
		{exp: "hello world ", str: "(hello world )", valid: true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("exp %s for %s", tt.exp, tt.str), func(t *testing.T) {
			got, valid := goql.Unquote(tt.str, '(', ')')
			if valid != tt.valid {
				t.Errorf("expected %t, got %t", tt.valid, valid)
			}
			if got != tt.exp {
				t.Errorf("expected %s, got %s", tt.exp, got)
			}
		})
	}
}
