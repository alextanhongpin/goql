package goql

import (
	"regexp"
)

var tagRe = regexp.MustCompile(`(?P<name>^[\w-]+)(,(?P<null>null|notnull))?(,type:(?P<type>\w+)(?P<array>\[\])?)?(,format:(?P<format>[\w-]+))?$`)

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

func ParseTag(tag string) Column {
	m := match(tagRe, tag)

	return Column{
		IsNull:  m["null"] == "null",
		IsArray: m["array"] != "",
		SQLType: m["type"],
		Name:    m["name"],
		Format:  m["format"],
		Tag:     tag,
	}
}
