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
