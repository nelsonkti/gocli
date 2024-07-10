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

// ToSnakeCase 将驼峰命名法的字符串转换为下划线命名法
func ToSnakeCase(s string) string {
	var sb strings.Builder

	for i, r := range s {
		if unicode.IsUpper(r) {
			// 如果不是第一个字符，则在前面添加下划线
			if i > 0 {
				sb.WriteRune('_')
			}
			// 将大写字母转换为小写字母
			sb.WriteRune(unicode.ToLower(r))
		} else {
			sb.WriteRune(r)
		}
	}

	return sb.String()
}
