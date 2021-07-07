package rflt

import (
	"reflect"
	"strings"

	"github.com/seerx/gpa/engine/sql/metas/schema"
)

func ParseTagProps(col *schema.Column, field reflect.StructField, fieldVal reflect.Value) error {
	return nil
}

func SplitTag(tag string) (tags []string) {
	tag = strings.TrimSpace(tag)
	var hasQuote = false
	var lastIdx = 0
	for i, t := range tag {
		if t == '\'' {
			hasQuote = !hasQuote
		} else if t == ' ' {
			if lastIdx < i && !hasQuote {
				tags = append(tags, strings.TrimSpace(tag[lastIdx:i]))
				lastIdx = i + 1
			}
		}
	}
	if lastIdx < len(tag) {
		tags = append(tags, strings.TrimSpace(tag[lastIdx:]))
	}
	return
}
