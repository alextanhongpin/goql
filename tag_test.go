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
		exp  goql.Tag
	}{
		{
			name: "empty",
			tag:  "",
			exp:  goql.Tag{},
		},
		{
			name: "field only",
			tag:  "name",
			exp: goql.Tag{
				Name: "name",
				Tag:  "name",
			},
		},
		{
			name: "field",
			tag:  "age,null",
			exp: goql.Tag{
				Name: "age",
				Type: goql.Type{
					Null: true,
				},
				Tag: "age,null",
			},
		},
		{
			name: "field type",
			tag:  "id,type:uuid",
			exp: goql.Tag{
				Name: "id",
				Type: goql.Type{
					Name: "uuid",
				},
				Ops: goql.OpsComparable,
				Tag: "id,type:uuid",
			},
		},
		{
			name: "field type array",
			tag:  "id,type:[]uuid",
			exp: goql.Tag{
				Name: "id",
				Type: goql.Type{
					Name:  "uuid",
					Array: true,
				},
				Ops: goql.OpsComparable | goql.OpsRange,
				Tag: "id,type:[]uuid",
			},
		},
		{
			name: "field type null alternative",
			tag:  "birthday,type:*date",
			exp: goql.Tag{
				Name: "birthday",
				Type: goql.Type{
					Name: "*date",
					Null: true,
				},
				Ops: goql.OpsComparable | goql.OpsNull,
				Tag: "birthday,type:*date",
			},
		},
		{
			name: "field ops",
			tag:  "name,ops:eq",
			exp: goql.Tag{
				Name: "name",
				Tag:  "name,ops:eq",
				Ops:  goql.OpEq,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := goql.ParseTag(tt.tag)
			if err != nil {
				t.Fatalf("error parsing tag: %v", err)
			}

			if diff := cmp.Diff(tt.exp, *got); diff != "" {
				t.Fatalf("+exp, -got: %s", diff)
			}
		})
	}
}
