package xtype

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"path/filepath"
	"strings"

	"github.com/seerx/gpa/engine/objs"
	"github.com/seerx/gpa/engine/sql/names"
	"github.com/seerx/gpa/logger"
)

func (x *XTypeParser) scan(dialect string, thisPkg, dir string, tagName string, logger logger.GpaLogger) (*poolObj, error) {
	var err error
	dir, err = filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	params := &poolObj{
		objs:      map[string]*XType{},
		tempFuncs: map[string][]*Func{},
	}

	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			// f.Name.Name
			var param *XType
			ast.Inspect(f, func(n ast.Node) bool {
				filepath := path.Join(dir, f.Name.Name)
				switch t := n.(type) {
				case *ast.ImportSpec:
					// import 的包
					path := strings.Trim(t.Path.Value, `"`)
					var name string
					if t.Name == nil {
						idx := strings.LastIndex(path, "/")
						if idx < 0 {
							name = path
						} else {
							name = path[idx+1:]
						}
					} else {
						name = t.Name.Name
					}
					x.addImport(filepath, name, path)
				case *ast.TypeSpec: // type 定义行
					param = nil
					if _, ok := t.Type.(*ast.StructType); ok {
						param = &XType{
							File:      filepath,
							Package:   thisPkg,
							Name:      t.Name.Name,
							TableName: names.ToTableName(t.Name.Name),
						}
						params.AddXType(param)
						// params[param.Name] = param
					}
				case *ast.StructType: // struct 定义结构体
					if param != nil {
						param.tempStructType = t
						// if err := param.ParseFields(t, tagName); err != nil {
						// 	logger.Error(err, "parse struct fields")
						// 	return false
						// }
					}
				case *ast.FuncDecl:
					fn, err := parseFunc(t, dialect)
					if err != nil {
						logger.Error(err, "parse function")
						return false
					}
					if fn != nil {
						params.AddFunc(fn)
					}
				}
				return true
			})
		}
	}

	for _, p := range params.objs {
		if p.tempStructType != nil {
			if err := p.ParseFields(p.tempStructType, tagName, params.objs, x); err != nil {
				logger.Error(err, "parse struct fields")
			}
		}
	}

	return params, nil
}

func parseFunc(decl *ast.FuncDecl, dialect string) (*Func, error) {
	if decl.Recv == nil {
		// 不是结构体的函数，不解析
		return nil, nil
	}
	recv := (*ast.FieldList)(decl.Recv)
	recvObj := recv.List[0]
	recvName := ""
	if recvObj.Names != nil {
		recvName = recvObj.Names[0].Name
	}
	recvTypeName := ""
	ptr := false
	recvType := (ast.Expr)(recvObj.Type)
	temp, ok := recvType.(*ast.StarExpr)
	if ok {
		x := (ast.Expr)(temp.X)
		xi, _ := x.(*ast.Ident)
		recvTypeName = xi.Name
		ptr = true
	} else {
		xi, _ := recvType.(*ast.Ident)
		recvTypeName = xi.Name
	}
	// if recvObj.Type
	fn := &Func{
		RecvName:  recvName,
		RecvType:  recvTypeName,
		RecvIsPtr: ptr,
	}

	obj := objs.NewEmptyObject()
	obj.Name = decl.Name.Name
	obj.ParseFunc(decl.Type.Params, decl.Type.Results, dialect, nil)
	fn.Object = *obj

	return fn, nil
}
