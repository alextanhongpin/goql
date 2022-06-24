package goql

import (
	"regexp"
)

var tagRe = regexp.MustCompile(`(?P<name>^[\w-]*)(,(?P<null1>null|notnull))?(,type:(?P<array>\[\])?(?P<type>\w+)(?P<null2>\?)?)?(,format:(?P<format>[\w-]+))?$`)

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
		IsNull:  m["null1"] == "null" || m["null2"] != "",
		IsArray: m["array"] != "",
		SQLType: m["type"],
		Name:    m["name"],
		Format:  m["format"],
		Tag:     tag,
	}
}
