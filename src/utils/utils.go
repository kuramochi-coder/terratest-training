package utils

// TrimQuotes trims the quotes from a string.
func TrimQuotes(s string) string {
	// Set the default length.
	defaultLength := 2

	// Trim the quotes from the string.
	if len(s) >= defaultLength {
		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
			return s[1 : len(s)-1]
		}
	}

	return s
}
