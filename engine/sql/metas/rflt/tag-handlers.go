package rflt

import (
	"errors"
	"strconv"
	"strings"

	"github.com/seerx/gpa/engine/sql/metas/schema"
	"github.com/seerx/gpa/engine/sql/types"
)

const (
	// nameOfTag        = "gpa" // tag 名称
	TagIgnore        = "-"
	TagPrimaryKey    = "pk"
	TagAutoIncrement = "autoincr"
	TagNullAble      = "allow-null"
	TagNotNull       = "not-null"
	TagIndex         = "index"
	TagUnique        = "unique"
	TagDefault       = "default"
	TagUTC           = "utc"
	TagLocal         = "local"
)

type Context struct {
	tagName    string
	params     []string
	col        *schema.Column
	indexNames map[string]int
}

type handler func(ctx *Context) error

var tagHandlers = map[string]handler{
	strings.ToUpper(TagIgnore):        func(ctx *Context) error { ctx.col.Ignore = true; return nil },
	strings.ToUpper(TagPrimaryKey):    func(ctx *Context) error { ctx.col.IsPrimaryKey = true; return nil },
	strings.ToUpper(TagAutoIncrement): func(ctx *Context) error { ctx.col.IsAutoIncrement = true; return nil },
	strings.ToUpper(TagNullAble):      nullAbleHandler,
	strings.ToUpper(TagNotNull):       func(ctx *Context) error { ctx.col.Nullable = true; return nil },
	strings.ToUpper(TagIndex):         indexHandler,
	strings.ToUpper(TagUnique):        uniqueHandler,
	strings.ToUpper(TagDefault):       defaultValueHandler,
}

func init() {
	for k := range types.SqlTypes {
		tagHandlers[k] = sqlTypeTagHandler
	}
}

func nullAbleHandler(ctx *Context) error {
	if !ctx.col.Nullable {
		ctx.col.Nullable = true
	}
	return nil
}

func indexHandler(ctx *Context) error {
	if len(ctx.params) > 0 {
		// ctx.tag.IndexName = ctx.params[0]
		ctx.indexNames[ctx.params[0]] = types.IndexType
	} else {
		// ctx.tag.IndexName = ctx.tag.Name
		// ctx.indexNames[ctx.params[0]]
		ctx.col.IsIndex = true
	}
	return nil
}

func uniqueHandler(ctx *Context) error {
	if len(ctx.params) > 0 {
		// ctx.tag.UniqueName = ctx.params[0]
		ctx.indexNames[ctx.params[0]] = types.UniqueType
	} else {
		ctx.col.IsUnique = true
		// ctx.tag.UniqueName = ctx.tag.Name
	}
	return nil
}

func defaultValueHandler(ctx *Context) error {
	if len(ctx.params) > 0 {
		ctx.col.Default = ctx.params[0]
		return nil
	}
	return errors.New("default must with value and wraped by ( and ), like default(10)")
}

// sqlTypeTagHandler describes SQL Type tag handler
func sqlTypeTagHandler(ctx *Context) error {
	ctx.col.Type = types.SQLType{Name: ctx.tagName}
	if strings.EqualFold(ctx.tagName, "JSON") {
		ctx.col.IsJSON = true
	}
	if len(ctx.params) > 0 {
		if ctx.tagName == types.Enum {
			// ctx.col.EnumOptions = make(map[string]int)
			// for k, v := range ctx.params {
			// 	v = strings.TrimSpace(v)
			// 	v = strings.Trim(v, "'")
			// 	ctx.col.EnumOptions[v] = k
			// }
		} else if ctx.tagName == types.Set {
			// ctx.col.SetOptions = make(map[string]int)
			// for k, v := range ctx.params {
			// 	v = strings.TrimSpace(v)
			// 	v = strings.Trim(v, "'")
			// 	ctx.col.SetOptions[v] = k
			// }
		} else {
			var err error
			if len(ctx.params) == 2 {
				ctx.col.Length, err = strconv.Atoi(ctx.params[0])
				if err != nil {
					return err
				}
				ctx.col.Length2, err = strconv.Atoi(ctx.params[1])
				if err != nil {
					return err
				}
			} else if len(ctx.params) == 1 {
				ctx.col.Length, err = strconv.Atoi(ctx.params[0])
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
