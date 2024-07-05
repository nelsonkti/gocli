package helper

import (
	"strings"
	"unicode"
)

// SnakeToCamel 将下划线格式的字符串转换为驼峰格式的字符串
func SnakeToCamel(s string) string {
	s = strings.ToLower(s)
	var result strings.Builder
	upperNext := false

	for i, r := range s {
		if r == '_' {
			upperNext = true
		} else {
			if upperNext || i == 0 {
				result.WriteRune(unicode.ToUpper(r))
				upperNext = false
			} else {
				result.WriteRune(r)
			}
		}
	}

	return result.String()
}

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
