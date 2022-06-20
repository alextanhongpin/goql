package goql

import (
	"encoding/csv"
	"strings"
)

var QueryDelimiter = ','

func joinStrings(ss []string) (string, error) {
	var sb strings.Builder
	w := csv.NewWriter(&sb)
	w.Comma = QueryDelimiter
	w.UseCRLF = false

	if err := w.Write(ss); err != nil {
		return "", err
	}

	w.Flush()

	return strings.TrimSpace(sb.String()), nil
}

func splitString(s string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(s))
	r.Comma = QueryDelimiter
	return r.Read()
}

func split2(str, by string) (string, string) {
	paths := strings.SplitN(str, by, 2)
	switch len(paths) {
	case 1:
		return paths[0], ""
	case 2:
		return paths[0], paths[1]
	default:
		return "", ""
	}
}

func Unquote(str string, l, r rune) (string, bool) {
	if len(str) < 2 {
		return str, false
	}

	if rune(str[0]) == l && rune(str[len(str)-1]) == r {
		return str[1 : len(str)-1], true
	}

	return str, false
}
