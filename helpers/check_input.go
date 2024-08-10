package helpers

import (
	"regexp"
	"strings"
)

func RemoveAlphaNum(text string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")

	processedString := reg.ReplaceAllString(text, "")

	return processedString
}

func RemoveSpaces(text string) string {
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", " ")

	return text
}

func Remove(text, char string) string {
	text = strings.ReplaceAll(text, char, "")

	return text
}

func GetValidName(input string) string {
	return strings.ToLower(Remove(RemoveSpaces(input), " "))
}

func GetKind(text string) string {
	trimmed := strings.TrimPrefix(text, "/")

	parts := strings.Split(trimmed, "/")

	return parts[0]
}
