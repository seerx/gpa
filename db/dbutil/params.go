package dbutil

import "strings"

func TakePlaceHolder(sql, placeHolder string, count int) string {
	var ps []string
	for n := 0; n < count; n++ {
		ps = append(ps, "?")
	}
	return strings.ReplaceAll(sql, placeHolder, strings.Join(ps, ","))
}
