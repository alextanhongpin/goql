package goql

import (
	"reflect"
	"regexp"
	"strings"
)

var tagRe = regexp.MustCompile(`(?P<name>^[\w-]*)(,(?P<null1>null))?(,type:(?P<array>\[\])?(?P<null2>\*)?(?P<type>\w+))?$`)

type Tag struct {
	Type Type
	Name string
	Tag  string
}

func match(re *regexp.Regexp, str string) map[string]string {
	if str == "" {
		return nil
	}

	match := re.FindStringSubmatch(str)
	m := make(map[string]string)
	for i, name := range tagRe.SubexpNames() {
		if i != 0 && name != "" {
			m[name] = match[i]
		}
	}

	return m
}

func ParseTag(tag string) Tag {
	m := match(tagRe, tag)

	return Tag{
		Name: m["name"],
		Type: Type{
			Name:  m["null2"] + m["type"],
			Null:  m["null1"] == "null" || m["null2"] != "",
			Array: m["array"] != "",
		},
		Tag: tag,
	}
}

func ParseStruct(unk any, key string) map[string]Tag {
	tagByField := make(map[string]Tag)

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
			c.Name = strings.ToLower(f.Name)
		}

		if c.Type.Name != "" {
			tagByField[c.Name] = c
			continue
		}

		c.Type = TypeOf(f.Type)
		tagByField[c.Name] = c
	}

	return tagByField
}
