package goql

import (
	"encoding/csv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var QueryDelimiter = ','

var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

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

func split3(str, by string) (string, string, string) {
	a, b := split2(str, by)
	c, d := split2(b, by)

	return a, c, d
}

func lowerCommonInitialism(field string) string {
	if field == "" {
		return ""
	}

	if commonInitialisms[field] {
		return strings.ToLower(field)
	}

	r, i := utf8.DecodeRuneInString(field)
	return string(unicode.ToLower(r)) + field[i:]
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
