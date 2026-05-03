package utils

import "strings"

// Pluralize returns a basic English plural form for generator naming use-cases.
func Pluralize(word string) string {
	trimmed := strings.TrimSpace(word)
	if trimmed == "" {
		return trimmed
	}

	lower := strings.ToLower(trimmed)

	if strings.HasSuffix(lower, "ch") || strings.HasSuffix(lower, "sh") || strings.HasSuffix(lower, "s") || strings.HasSuffix(lower, "x") || strings.HasSuffix(lower, "z") {
		return lower + "es"
	}

	if strings.HasSuffix(lower, "y") && len(lower) > 1 {
		prev := lower[len(lower)-2]
		if !strings.ContainsRune("aeiou", rune(prev)) {
			return lower[:len(lower)-1] + "ies"
		}
	}

	if strings.HasSuffix(lower, "fe") {
		return lower[:len(lower)-2] + "ves"
	}

	if strings.HasSuffix(lower, "f") {
		return lower[:len(lower)-1] + "ves"
	}

	return lower + "s"
}
