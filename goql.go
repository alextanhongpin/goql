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
		if strings.HasPrefix(tag, "-") {
			continue
		}

		c := ParseTag(tag)
		if c.Name == "" {
			c.Name = f.Name
		}

		if c.SQLType != "" {
			columnByFieldName[c.Name] = c
			continue
		}

		sqlType, null, array := GetSQLType(f.Type)
		c.SQLType = sqlType
		c.IsNull = null
		c.IsArray = array

		columnByFieldName[c.Name] = c
	}

	return columnByFieldName
}
