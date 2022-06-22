package goql

import (
	"reflect"
	"strings"
)

type Column struct {
	IsNull  bool
	IsArray bool
	SQLType string
	Name    string
	Format  string
	Tag     string
}

// Provide different infer, one for query, another for sql, another for mongo etc
// TODO: https://github.com/go-pg/pg/blob/782c9d35ba243106ba6445fc753c3ac6a14c3324/orm/table.go
func StructToColumns(unk any, key string) map[string]Column {
	columnByFieldName := make(map[string]Column)

	v := reflect.Indirect(reflect.ValueOf(unk))
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get(key)
		if tag == "-" {
			continue
		}

		c := Column{Tag: tag}

		if tag == "" {
			tag = f.Name
		}

		values := strings.Split(tag, ",")
		switch len(values) {
		case 1:
			c.Name = values[0]
		case 2, 3, 4:
			c.Name = values[0]
			isNull, valid := IsSQLNull(values[1])
			if !valid {
				panic("goql: second argument must be null or notnull")
			}

			c.IsNull = isNull

			for _, val := range values[2:] {
				k, v := split2(val, ":")
				switch k {
				case "type":
					v, c.IsArray = ParseSQLArray(v)
				case "format":
					c.Format = v
				default:
					panic("goql: unexpected tag")
				}
			}
		default:
			panic("goql: invalid tag")
		}

		if c.SQLType != "" {
			columnByFieldName[c.Name] = c
			continue
		}

		sqlType, null, array := getSQLType(f.Type)
		c.SQLType = sqlType
		c.IsNull = null
		c.IsArray = array
		columnByFieldName[c.Name] = c
	}

	return columnByFieldName
}

func IsSQLNull(s string) (null bool, valid bool) {
	switch s {
	case "null", "notnull":
		null = s == "null"
		valid = true
	}

	return
}

func ParseSQLArray(s string) (base string, array bool) {
	base = s
	array = strings.HasPrefix(s, "[]")
	if array {
		base = strings.TrimPrefix(s, "[]")
	}

	return
}
