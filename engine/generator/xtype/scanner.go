package xtype

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"path/filepath"

	"github.com/seerx/gpa/engine/sql/names"
	"github.com/seerx/gpa/logger"
)

func scan(dir string, tagName string, logger logger.GpaLogger) (*poolObj, error) {
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
				switch t := n.(type) {
				case *ast.TypeSpec: // type 定义行
					param = nil
					if _, ok := t.Type.(*ast.StructType); ok {
						param = &XType{
							File:      path.Join(dir, f.Name.Name),
							Name:      t.Name.Name,
							TableName: names.ToTableName(t.Name.Name),
						}
						params.AddXType(param)
						// params[param.Name] = param
					}
				case *ast.StructType: // struct 定义结构体
					if param != nil {
						if err := param.ParseFields(t, tagName); err != nil {
							logger.Error(err, "parse struct fields")
							return false
						}
					}
				case *ast.FuncDecl:
					fn, err := parseFunc(t)
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

	return params, nil
}

func parseFunc(decl *ast.FuncDecl) (*Func, error) {
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

	for _, p := range decl.Type.Params.List {
		fmt.Println(p)
	}

	return fn, nil
}
