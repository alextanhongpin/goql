package goql

import (
	"reflect"
	"strings"
)

type column struct {
	null    bool
	array   bool
	sqlType string
	name    string
	parse   string
	tag     string
}

// Provide different infer, one for query, another for sql, another for mongo etc
// TODO: https://github.com/go-pg/pg/blob/782c9d35ba243106ba6445fc753c3ac6a14c3324/orm/table.go
func structToColumns(unk any, tagKey string) map[string]column {
	columnByFieldName := make(map[string]column)

	v := reflect.Indirect(reflect.ValueOf(unk))
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get(tagKey)

		c := column{tag: tag}
		if tag == "" {
			tag = f.Name
		}

		tagValues := strings.Split(tag, ",")
		switch len(tagValues) {
		case 1:
			c.name = tagValues[0]
		case 2, 3, 4:
			c.name = tagValues[0]
			if tagValues[1] != "null" {
				panic("expected null as second argument")
			}
			c.null = tagValues[1] == "null"
			for _, val := range tagValues[2:] {
				k, v := split2(val, ":")
				switch k {
				case "type":
					if strings.HasPrefix(v, "[]") {
						v = strings.TrimPrefix(v, "[]")
						c.array = true
					}

					c.sqlType = v
				case "parse":
					c.parse = v
				default:
					panic("unexpected tag")
				}
			}
		default:
			panic("invalid tag")
		}

		if c.sqlType != "" {
			columnByFieldName[c.name] = c
			continue
		}

		ft := f.Type

		/*
			Handles only the following
			field *Struct
			field Struct
			field []*Struct
			field *[]*Struct
			field *[]Struct
		*/

		switch ft.Kind() {
		case reflect.Pointer:
			c.null = true

			// Get the value of the pointer.
			ft = ft.Elem()
		}

		switch ft.Kind() {
		case reflect.Pointer:
			c.null = true

			// Get the value of the pointer.
			ft = ft.Elem()
		case reflect.Slice:
			c.array = true

			// Get the item of the slice.
			ft = ft.Elem()

			// Check the item if it is a pointer.
			switch ft.Kind() {
			case reflect.Pointer:
				c.null = true

				ft = ft.Elem()
			}
		}

		c.sqlType = sqlType(ft)
		columnByFieldName[c.name] = c
	}

	return columnByFieldName
}
