package goql_test

import (
	"testing"

	"github.com/alextanhongpin/goql"
	"github.com/google/go-cmp/cmp"
)

func TestTag(t *testing.T) {
	tests := []struct {
		name string
		tag  string
		exp  goql.Column
	}{
		{
			name: "empty",
			tag:  "",
			exp:  goql.Column{},
		},
		{
			name: "field only",
			tag:  "name",
			exp: goql.Column{
				Name: "name",
				Tag:  "name",
			},
		},
		{
			name: "field null",
			tag:  "age,null",
			exp: goql.Column{
				Name:   "age",
				IsNull: true,
				Tag:    "age,null",
			},
		},
		{
			name: "field notnull",
			tag:  "age,notnull",
			exp: goql.Column{
				Name: "age",
				Tag:  "age,notnull",
			},
		},
		{
			name: "field notnull",
			tag:  "age,notnull",
			exp: goql.Column{
				Name: "age",
				Tag:  "age,notnull",
			},
		},
		{
			name: "field type",
			tag:  "id,type:uuid",
			exp: goql.Column{
				Name:    "id",
				SQLType: "uuid",
				Tag:     "id,type:uuid",
			},
		},
		{
			name: "field type array",
			tag:  "id,type:uuid[]",
			exp: goql.Column{
				Name:    "id",
				SQLType: "uuid",
				IsArray: true,
				Tag:     "id,type:uuid[]",
			},
		},
		{
			name: "field type format",
			tag:  "birthday,type:date,format:20060102",
			exp: goql.Column{
				Format:  "20060102",
				Name:    "birthday",
				SQLType: "date",
				Tag:     "birthday,type:date,format:20060102",
			},
		},
		{
			name: "field type format array",
			tag:  "birthday,type:date[],format:20060102",
			exp: goql.Column{
				Format:  "20060102",
				Name:    "birthday",
				IsArray: true,
				SQLType: "date",
				Tag:     "birthday,type:date[],format:20060102",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := goql.ParseTag(tt.tag)
			if diff := cmp.Diff(tt.exp, got); diff != "" {
				t.Fatalf("+exp, -got: %s", diff)
			}
		})
	}
}
