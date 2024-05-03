package utils

import "regexp"

func Sanitize_Discord_Text(input string) string {
	// Regular expression to find custom Discord emojis which are in format <:name:id>
	re := regexp.MustCompile(`<:([^:]+):\d+>`)
	// Replace all instances of the custom emoji with just its name wrapped in colons
	return re.ReplaceAllString(input, ":$1:")
}
