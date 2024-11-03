package distext

import (
	"strings"
	"unicode/utf8"
)

func Truncate(text, truncatedText string, size int) string {
	if utf8.RuneCountInString(text) > size {
		slice := []rune(text)
		strarr := []string{string(slice[:size]), truncatedText}
		text = strings.Join(strarr, " ")
	}
	return text
}
