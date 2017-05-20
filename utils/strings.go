package utils

import (
	"strings"
)

// Trim removes trailing blanks and replace all carriage returns by a space
func Trim(s string) string {
	s = strings.Replace(s, "\n", " ", -1)
	return strings.TrimSpace(s)
}

// Expand resplits a list of values by applying all the separator specified
// The result is still a one dimensional array, but may include more elements
func Expand(separators string, params []string) []string {
	if len(separators) == 0 {
		return params
	}

	result := make([]string, 0, len(params))

	for _, param := range params {
		for _, sep := range separators[1:] {
			param = strings.Replace(param, string(sep), string(separators[0]), -1)
		}

		for _, param := range strings.Split(param, string(separators[0])) {
			result = append(result, strings.TrimSpace(param))
		}
	}
	return result
}
