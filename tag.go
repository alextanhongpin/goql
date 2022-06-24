package goql

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var tagRe = regexp.MustCompile(`(?P<name>^[\w-]*)(,(?P<null1>null))?(,type:(?P<array>\[\])?(?P<null2>\*)?(?P<type>\w+))?(,ops:(?P<ops>(\w+(,\w+)*)+))?$`)

type Tag struct {
	Type Type
	Name string
	Tag  string
	Ops  Op
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

func ParseTag(tag string) (*Tag, error) {
	m := match(tagRe, tag)

	var ops Op
	for _, raw := range strings.Split(m["ops"], ",") {
		if raw == "" {
			continue
		}

		op, ok := ParseOp(raw)
		if !ok {
			return nil, fmt.Errorf("%w: %q", ErrUnknownOperator, op)
		}

		ops |= op
	}

	t := Type{
		Name:  m["null2"] + m["type"],
		Null:  m["null1"] == "null" || m["null2"] != "",
		Array: m["array"] != "",
	}

	if ops == 0 {
		ops = NewOps(t)
	}

	return &Tag{
		Name: m["name"],
		Type: t,
		Tag:  tag,
		Ops:  ops,
	}, nil
}

func NewOps(t Type) Op {
	if t.Name == "" {
		return 0
	}

	// All types are comparable.
	ops := OpsComparable

	// Null type have special operators.
	if t.Null {
		ops |= OpsNull
	}

	// Array type have special operators.
	if t.Array {
		ops |= OpsRange
	}

	switch t.Name {
	// String types have special operators.
	case "string":
		ops |= OpsText
	// Bool types have special operators.
	case "bool":
		ops |= OpsNull
	}

	return ops
}

func ParseStruct(unk any, key string) (map[string]*Tag, error) {
	tagByField := make(map[string]*Tag)

	v := reflect.Indirect(reflect.ValueOf(unk))
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get(key)
		if strings.HasPrefix(tag, "-") {
			continue
		}

		c, err := ParseTag(tag)
		if err != nil {
			return nil, err
		}

		if c.Name == "" {
			c.Name = strings.ToLower(f.Name)
		}

		if c.Type.Name != "" {
			tagByField[c.Name] = c
			continue
		}

		c.Type = TypeOf(f.Type)
		if c.Ops == 0 {
			c.Ops = NewOps(c.Type)
		}
		tagByField[c.Name] = c
	}

	return tagByField, nil
}
