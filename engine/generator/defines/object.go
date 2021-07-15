package defines

import (
	"go/ast"

	"github.com/seerx/gpa/engine/generator/xtype"
	"github.com/seerx/gpa/engine/objs"
)

type Object struct {
	*objs.Object
	repo     *RepoInterface
	SQL      string
	beanType *xtype.XType
}

func NewObject(repo *RepoInterface) *Object {
	return &Object{
		repo:   repo,
		Object: objs.NewEmptyObject(),
	}
}

func (o *Object) GetBeanType() *xtype.XType {
	if o.beanType == nil {
		if o.Type.IsStruct() {
			dir, err := o.repo.repoFile.FindPackagePath(o.Type.Package)
			if err != nil {
				o.repo.logger.Error(err, "finding package dir (%s)", o.Type.String())
			} else {
				o.beanType, err = o.repo.repoFile.info.xtypeParser.Parse(o.Type.Name, dir)
				if err != nil {
					o.repo.logger.Error(err, "parsing struct (%s)", o.Type.String())
				}
			}
		}
	}
	return o.beanType
}

func (o *Object) Parse(field *ast.Field, expr ast.Expr, dialect string, level int) error {
	obj := objs.NewEmptyObject()
	obj.Parse(field, expr, dialect,
		func(obj *objs.Object) error {
			o.SQL = ParseSQL(field.Doc.Text(), dialect)
			return nil
		},
		level)
	o.Object = obj
	return nil
}

// func (o *Object) Parse1(field *ast.Field, expr ast.Expr, dialect string, level int) error {
// 	var err error
// 	switch pt := expr.(type) {
// 	// case *ast.Im:
// 	// 	fmt.Println(pt.Path, pt.Name.Name)
// 	case *ast.Ident:
// 		// 普通类型
// 		o.Type = *objs.NewPrimitiveType(pt.Name)
// 	case *ast.FuncType:
// 		// 函数类型
// 		if level > 1 {
// 			return errors.New("不支持多层嵌套函数类型")
// 		}

// 		fn := NewFuncWithObject(o)
// 		if err := fn.Parse(field, dialect); err != nil {
// 			return err
// 		}
// 		o.Type = *objs.NewFuncType()

// 		// arg, err = parseObject(repo, repoIntf, p.field)
// 		// if err != nil {
// 		// 	return nil, err
// 		// }
// 		// arg.Type = metas.Type{Name: "func"} // pkg.TypesInfo.ObjectOf(p.field.Names[0]).Type()
// 		// arg.Type = *objs.NewFuncType()
// 	case *ast.SelectorExpr:
// 		o.Type = *ParseSelectorExprType(pt, false)
// 	case *ast.StarExpr:
// 		o.Type = *ParseSelectorExprType(pt.X.(*ast.SelectorExpr), true)
// 	case *ast.ArrayType:
// 		if level > 1 {
// 			return errors.New("不支持多层嵌套数据类型")
// 		}
// 		obj := NewObject(o.repo)
// 		if err := obj.Parse(nil, pt.Elt, dialect, level+1); err != nil {
// 			return err
// 		}
// 		// if _, err := parseParam(repo, repoIntf, arg, &pp{expr: pt.Elt}, level+1); err != nil {
// 		// 	return err
// 		// }
// 		// arg.IsSlice()
// 		o.Object = obj.Object
// 		o.IsSlice = true
// 		// arg.Slice = true
// 	case *ast.SliceExpr:
// 		if level > 1 {
// 			return errors.New("不支持多层嵌套数据类型")
// 		}
// 		// if _, err := parseParam(repo, repoIntf, arg, &pp{expr: pt.X}, level+1); err != nil {
// 		// 	return nil, err
// 		// }
// 		obj := NewObject(o.repo)
// 		if err := obj.Parse(nil, pt.X, dialect, level+1); err != nil {
// 			return err
// 		}
// 		// arg.Slice = true
// 		o.Object = obj.Object
// 		o.IsSlice = true
// 	case *ast.MapType:
// 		if level > 1 {
// 			return errors.New("不支持多层嵌套数据类型")
// 		}
// 		o.IsMap = true
// 		obj := NewObject(o.repo)
// 		if err := obj.Parse(nil, pt.Value, dialect, level+1); err != nil {
// 			return err
// 		}
// 		o.Object = obj.Object
// 		// pt.Key
// 		// if _, err := parseParam(repo, repoIntf, arg, &pp{expr: pt.Value}, level+1); err != nil {
// 		// 	return nil, err
// 		// }
// 		// o := metas.NewObject(nil)
// 		// // o := &metas.Object{}
// 		// if _, err := parseParam(repo, repoIntf, o, &pp{expr: pt.Key}, level+1); err != nil {
// 		// 	return nil, err
// 		// }

// 		obj = NewObject(o.repo)
// 		if err := obj.Parse(nil, pt.Key, dialect, level+1); err != nil {
// 			return err
// 		}
// 		o.Key = &obj.Type

// 		// arg.KeyType = &o.Type
// 		// arg.Map = true
// 	default:
// 		err = errors.New("不支持的数据类型")
// 	}
// 	// fmt.Println(arg.Type.String())
// 	return err
// }
