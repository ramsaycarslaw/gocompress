// Copyright (C) 2020 Ramsay Carslaw

package whitespace

const (
	// TabStopDistance is how many spaces the compressor should llook for for tab replacement
	TabStopDistance = 4
)

// RemoveTrailingWhitespace consumes any whitespace at the end of a source string
func RemoveTrailingWhitespace(src string) string {
	i := len(src) - 1

	if !isWhitespace(rune(src[i])) {
		return src
	}

	for isWhitespace(rune(src[i])) {
		i--
		continue
	}
	return src[:i+1]
}

// ReplaceSpaces replaces every TabStopDistance spaces in a row with a tab
func ReplaceSpaces(src string) string {
	var spacecount int
	for i, v := range src {
		if spacecount == TabStopDistance {
			var before, after string
			before = src[:i-4]
			after = src[i:]
			src = before + "\t" + after
		}

		if v == ' ' {
			spacecount++
			continue
		} else {
			spacecount = 0
		}
	}
	return src
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r'
}
