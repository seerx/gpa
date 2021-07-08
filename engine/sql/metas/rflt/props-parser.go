package rflt

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/seerx/gpa/db/dbutil"
	"github.com/seerx/gpa/engine/sql/dialects"
	"github.com/seerx/gpa/engine/sql/metas/schema"
	"github.com/seerx/gpa/engine/sql/types"
)

type PropsParser struct {
	tagName string
	dialect dialects.Dialect
}

func NewPropsParser(tagName string, dialect dialects.Dialect) *PropsParser {
	return &PropsParser{tagName: tagName, dialect: dialect}
}

func (p *PropsParser) Parse(col *schema.Column, field reflect.StructField, fieldVal reflect.Value) error {
	tag, found := field.Tag.Lookup(p.tagName)
	if found {
		tag = strings.TrimSpace(tag)
	}
	col.Tag, col.Nullable = tag, true
	context := &Context{col: col, indexNames: map[string]int{}}
	fieldType := field.Type

	if tag != "" {
		// tag 不是空
		tags := SplitTag(tag)
		for _, item := range tags {
			ITEM := strings.ToUpper(item)
			context.tagName = ITEM
			// 查看 tag 是否带参数
			pStart := strings.Index(ITEM, "(")
			if pStart == 0 {
				return errors.New("( could not be the first character")
			}
			if pStart > -1 {
				if !strings.HasSuffix(ITEM, ")") {
					return fmt.Errorf("field %s tag %s cannot match ) character", field.Name, item)
				}
				context.tagName = ITEM[:pStart]
				params := strings.Split(item[pStart+1:len(ITEM)-1], ",")
				context.params = params
			}
			// 查找 tag 对应的 handler
			handler, ok := tagHandlers[context.tagName]
			if ok {
				// 找到 handler
				if err := handler(context); err != nil {
					return err
				}
				if col.Ignore {
					// 忽略该字段
					return nil
				}
			}
			if col.Type.Name == "" {
				// tag 中没有定义数据类型
				// 使用 golang 数据类型确定 sql 数据类型
				col.Type = *p.dialect.ToSQLType(fieldType)
			}

			// 特定处理
			if col.Type.Name == types.Serial || col.Type.Name == types.BigSerial {
				col.IsAutoIncrement = true
				col.Nullable = false
			}
			// 如果数据长度是 0， 使用默认长度
			if col.Length == 0 {
				col.Length = col.Type.Length
			}
			if col.Length2 == 0 {
				col.Length2 = col.Type.Length2
			}
			if col.IsUnique {
				// 把当前字段的名称作为唯一索引
				context.indexNames[col.FieldName()] = types.UniqueType
			} else if col.IsIndex {
				// 把当前字段的名称作为索引
				context.indexNames[col.FieldName()] = types.IndexType
			}
		}
	} else if fieldVal.CanSet() {
		// 没有 tag 时，依据 golang 定义生成 sql 数据类型
		var sqlType *types.SQLType
		if fieldVal.CanAddr() {
			if _, ok := fieldVal.Addr().Interface().(dbutil.BlobReadWriter); ok {
				sqlType = &types.SQLType{Name: types.Blob}
			}
		}
		if _, ok := fieldVal.Interface().(dbutil.BlobReadWriter); ok {
			sqlType = &types.SQLType{Name: types.Blob}
		} else {
			sqlType = p.dialect.ToSQLType(fieldType)
		}
		col.Type = *sqlType
		col.Length = sqlType.Length
		col.Length2 = sqlType.Length2
	} else {
		// 不可用字段
		return nil
	}

	// col.FieldName() = names.ToTableFieldName(field.Name)
	if col.IsAutoIncrement {
		// 自增长类型不能为 null
		col.Nullable = false
	}
	col.Indexes = context.indexNames
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
