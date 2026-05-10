package search

import (
	"strings"
	"unicode"
)

func BuildFTSQuery(input string) string {
	fields := strings.Fields(boundaryNormalize(input))
	parts := make([]string, 0, len(fields))
	for _, field := range fields {
		token := cleanToken(field)
		if token == "" {
			continue
		}
		parts = append(parts, `"`+strings.ReplaceAll(token, `"`, `""`)+`"`)
	}
	return strings.Join(parts, " AND ")
}

func boundaryNormalize(value string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.' || r == '/' || r == ':' {
			return r
		}
		return ' '
	}, value)
}

func cleanToken(value string) string {
	value = strings.TrimSpace(value)
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.' || r == '/' || r == ':' {
			return r
		}
		return -1
	}, value)
}
