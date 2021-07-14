package xtype

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"

	"github.com/seerx/gpa/engine/objs"
	"github.com/seerx/gpa/engine/sql/metas/rflt"
	"github.com/seerx/gpa/engine/sql/metas/schema"
	"github.com/seerx/gpa/engine/sql/names"
	"github.com/seerx/gpa/engine/sql/types"
	"github.com/seerx/gpa/logger"
)

type Field struct {
	schema.Column

	VarName  string   // 对应的 golang 中的变量名称
	argNames []string // 对应的参数名称列表，首字母小写和原来的名称 即 VarName
}

type XType struct {
	File      string   // 所在文件名称
	Name      string   // 名称
	TableName string   // 对应的数据库表名
	Fields    []*Field // 结构体对应的成员列表
	Funcs     []*Func  // 结构体对应的函数列表
}

type Func struct {
	objs.Object
	RecvName  string
	RecvType  string
	RecvIsPtr bool
}

// var paramsPool = map[string]map[string]*ParamType{}

type poolObj struct {
	objs      map[string]*XType
	tempFuncs map[string][]*Func
}

func (p *poolObj) AddXType(xt *XType) {
	fns, ok := p.tempFuncs[xt.Name]
	if ok {
		xt.Funcs = fns
		delete(p.tempFuncs, xt.Name)
	}
	p.objs[xt.Name] = xt
}

func (p *poolObj) AddFunc(fn *Func) {
	xt, ok := p.objs[fn.RecvType]
	if ok {
		xt.Funcs = append(xt.Funcs, fn)
		return
	}

	p.tempFuncs[fn.RecvType] = append(p.tempFuncs[fn.RecvType], fn)
}

type XTypeParser struct {
	pool    map[string]*poolObj
	tagName string
	logger  logger.GpaLogger
}

func NewXTypeParser(tagName string, logger logger.GpaLogger) *XTypeParser {
	return &XTypeParser{
		tagName: tagName,
		logger:  logger,
		pool:    map[string]*poolObj{},
	}
}

func (x *XTypeParser) Parse(name, dir string) (*XType, error) {
	var err error
	params, ok := x.pool[dir]
	if !ok {
		params, err = scan(dir, x.tagName, x.logger)
		if err != nil {
			return nil, err
		}
		x.pool[dir] = params
	}

	param, ok := params.objs[name]
	if ok {
		return param, nil
	}
	return nil, fmt.Errorf("no struct %s is defined in %s", name, dir)
}

func (f *Field) GetArgNames() []string {
	if f.argNames == nil {
		f.argNames = []string{
			names.LowerFirstChar(f.VarName),
			f.VarName,
		}
	}
	return f.argNames
}

func (p *XType) ParseFields(st *ast.StructType, tagName string) error {
	for _, fd := range st.Fields.List {
		if !ast.IsExported(fd.Names[0].Name) {
			continue
		}
		field := &Field{
			VarName: fd.Names[0].Name,
		}
		// 类型
		typ := objs.NewTypeByExpr(fd.Type)
		// 解析 tag
		tag := ""
		if fd.Tag != nil {
			tag = strings.Trim(fd.Tag.Value, "`")
			ta := reflect.StructTag(tag)
			tag, _ = ta.Lookup(tagName)
		}

		col := &schema.Column{
			Tag:      tag,
			Nullable: true,
			Field:    *objs.NewObject(field.VarName, *typ),
		}
		if tag != "" {
			context := rflt.NewContext(col)
			if err := rflt.ParseTag(col, context, field.VarName, tag, func() *types.SQLType {
				return col.Field.GetSQLTypeByType()
			}); err != nil {
				return err
			}
		}

		if col.Type == nil {
			col.Type = col.Field.GetSQLTypeByType()
		}
		if col.Type == nil {
			return fmt.Errorf("no sql types for field %s.%s", p.Name, field.VarName)
		}
		field.Column = *col
		p.Fields = append(p.Fields, field)
	}
	return nil
}