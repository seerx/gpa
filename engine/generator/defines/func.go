package defines

import (
	"errors"
	"fmt"
	"go/ast"
)

type Template string

const (
	INSERT Template = "insert"
	UPDATE Template = "update"
	DELETE Template = "delete"
	FIND   Template = "find"
	COUNT  Template = "count"
)

type Func struct {
	// repoIntf *RepoInterface
	*Object
	Template Template
}

func NewFuncWithObject(o *Object) *Func {
	return &Func{Object: o}
}

func NewFunc(repo *RepoInterface) *Func {
	return &Func{Object: NewEmptyObject(repo)}
}

func (f *Func) Parse(method *ast.Field, dialect string) error {
	f.Name = GetName(method.Names)               // 函数名称
	f.SQL = ParseSQL(method.Doc.Text(), dialect) // SQL 语句定义

	typ, ok := method.Type.(*ast.FuncType)
	if !ok {
		return fmt.Errorf("%s is not a valid method", f.Name)
	}

	// 遍历参数列表
	if typ.Params != nil {
		for _, mp := range typ.Params.List {
			param := NewEmptyObject(f.repo)
			// param.Object.Name = param.Name
			if err := param.Parse(mp, mp.Type, dialect, 0); err != nil {
				return err
			}
			param.Name = GetName(mp.Names)
			f.AddParam(param.Object)
		}
	}
	// 遍历返回值列表
	if typ.Results != nil {
		for _, p := range typ.Results.List {
			result := NewEmptyObject(f.repo) // {Name: getName(p.Names)}
			// result.Object.Name = result.Name
			if err := result.Parse(p, p.Type, dialect, 0); err != nil {
				return err
			}
			result.Name = GetName(p.Names)
			f.AddResult(result.Object)
		}
	}

	return nil
}

// func NewFunc(repoIntf *RepoInterface, obj *Object, sql string) *Func {
// 	return &Func{repoIntf: repoIntf, SQL: sql, Object: *obj}
// }

func (f *Func) CreateError(format string, v ...interface{}) error {
	return errors.New(f.Format(format, v...))
}

func (f *Func) Format(format string, v ...interface{}) string {
	return fmt.Sprintf("%s.%s %s\n%s", f.repo.Name, f.Name, f.repo.repoFile.Path, fmt.Sprintf(format, v...))
}

func (f *Func) AddSQLPackage() string {
	return f.repo.repoFile.AddSQLPackage()
}

func (f *Func) AddDBUtilPackage() string {
	return f.repo.repoFile.AddDBUtilPackage()
}

func (f *Func) AddRunTimePackage() string {
	return f.repo.repoFile.AddRuntimePackage()
}

func (f *Func) AddXTypePackage(fullPkgPath string) string {
	return f.repo.repoFile.AddXTypePackage(fullPkgPath)
}
