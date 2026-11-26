package utils

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Adds the given character to the left of the text,
// returning a string of the given size
func ShiftLeftFillingBlanks(text string, character string, length int) string {
	if len(text) > length {
		return text[:4]
	} else {

		return shiftFillingBlanks(text, character, length, true)
	}
}

// Adds the given character to the right of the text,
// returning a string of the given size
func ShiftRightFillingBlanks(text string, character string, length int) string {
	return shiftFillingBlanks(text, character, length, false)
}

func shiftFillingBlanks(text string, character string, length int, isLeft bool) string {
	if length < 1 {
		return text
	}

	shiftedText := text

	for {

		if len(shiftedText) < length {
			if isLeft {
				shiftedText = character + shiftedText
			} else {
				shiftedText = shiftedText + character
			}
		} else {

			break
		}
	}

	return shiftedText
}

// Reverses a string content
func Reverse(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}

func ReplaceTemplateWithMapping(template string, templateMap map[string]string) string {

	result := template
	for key, value := range templateMap {
		result = strings.Replace(result, key, value, -1)
	}
	return result
}

func ToSnakeCase(txt string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(txt, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
