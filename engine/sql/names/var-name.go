package names

import "strings"

func LowerFirstChar(str string) string {
	if str == "" {
		return str
	}
	// for key := range LintGonicMapper {
	// 	if strings.Index(str, key) == 0 {
	// 		return str
	// 	}
	// }
	return strings.ToLower(str[:1]) + str[1:]
}

func UpperFirstChar(str string) string {
	if str == "" {
		return str
	}
	// for key := range LintGonicMapper {
	// 	if strings.Index(str, strings.ToLower(key)) == 0 {
	// 		return key + str[len(key):]
	// 	}
	// }
	return strings.ToUpper(str[:1]) + str[1:]
}
