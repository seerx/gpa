package rflt

import (
	"go/ast"

	"github.com/seerx/gpa/engine/objs"
	"github.com/seerx/gpa/engine/sql/metas/schema"
	"github.com/seerx/gpa/engine/sql/names"
	"github.com/seerx/mro/utils"
)

func Parse(model interface{}, pp *PropsParser) (*schema.Table, error) {
	mValue, mType := utils.ReflectValueAndType(model)
	tableName := names.ToTableName(mValue)
	table := schema.NewTable(mType, tableName)
	for n := 0; n < mType.NumField(); n++ {
		field := mType.Field(n)
		// field.Type
		if !field.Anonymous && ast.IsExported(field.Name) {
			fieldObj := objs.NewObjectFromStructField(&field)
			col := &schema.Column{
				Field: *fieldObj,
			}
			fieldVal := mValue.Field(n)
			if err := pp.Parse(col, field, fieldVal); err != nil {
				return nil, err
			}
			if col.Ignore {
				continue
			}
			table.AddColumn(col)
			for name, typ := range col.Indexes {
				table.AddIndex(name, typ, col)
			}
		}
	}
	return table, nil
}
