package gentype

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	commonInitialisms = map[string]struct{}{
		"ACL":   {},
		"API":   {},
		"ASCII": {},
		"CPU":   {},
		"CSS":   {},
		"DNS":   {},
		"DSN":   {},
		"DTO":   {},
		"EOF":   {},
		"GUID":  {},
		"HTML":  {},
		"HTTP":  {},
		"HTTPS": {},
		"ID":    {},
		"IP":    {},
		"JSON":  {},
		"JWT":   {},
		"LHS":   {},
		"MIME":  {},
		"OCR":   {},
		"QPS":   {},
		"RAM":   {},
		"RHS":   {},
		"RPC":   {},
		"SLA":   {},
		"SMTP":  {},
		"SQL":   {},
		"SSH":   {},
		"SSL":   {},
		"TCP":   {},
		"TLS":   {},
		"TTL":   {},
		"UDP":   {},
		"UI":    {},
		"UID":   {},
		"UUID":  {},
		"URI":   {},
		"URL":   {},
		"UTF8":  {},
		"VM":    {},
		"XML":   {},
		"XMPP":  {},
		"XSRF":  {},
		"XSS":   {},
	}

	rxSeparators = regexp.MustCompile(`[^a-zA-Z0-9]+`)
)

func ToCamel(name string) string {
	return wordsToCamel(nameToWords(name))
}

func ToLowerCamel(name string) string {
	return wordsToLowerCamel(nameToWords(name))
}

//nolint:cyclop
func nameToWords(s string) []string {
	// Normalize separators to space
	s = rxSeparators.ReplaceAllString(s, " ")

	// Split by spaces
	parts := strings.Fields(s)

	var result []string

	for _, part := range parts {
		var buf strings.Builder

		for i, r := range part {
			if i == 0 {
				buf.WriteRune(r)

				continue
			}

			prev := rune(part[i-1])

			// Insert a space before:
			// 1. upper followed by lower (APIConfig -> API Config)
			// 2. letter/digit boundaries
			if (unicode.IsLower(prev) && unicode.IsUpper(r)) ||
				(unicode.IsLetter(prev) && unicode.IsDigit(r)) ||
				(unicode.IsDigit(prev) && unicode.IsLetter(r)) ||
				(unicode.IsUpper(prev) && unicode.IsUpper(r) && i+1 < len(part) && unicode.IsLower(rune(part[i+1]))) {
				buf.WriteRune(' ')
			}

			buf.WriteRune(r)
		}

		words := strings.Fields(buf.String())

		for _, w := range words {
			result = append(result, strings.ToLower(w))
		}
	}

	return result
}

func wordsToCamel(words []string) string {
	res := ""
	start := 0

	for i := start; i < len(words); i++ {
		word := words[i]

		initialism := strings.ToUpper(word)
		if _, ok := commonInitialisms[initialism]; ok {
			res += initialism

			continue
		}

		runes := []rune(word)
		runes[0] = unicode.ToTitle(runes[0])

		res += string(runes)
	}

	return res
}

func wordsToLowerCamel(words []string) string {
	if len(words) == 0 {
		return ""
	}

	res := strings.ToLower(words[0])

	if len(words) > 1 {
		res += wordsToCamel(words[1:])
	}

	return res
}
