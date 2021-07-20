package defines

import (
	"fmt"
	"go/ast"
	"go/build"
	"math/rand"
	"strings"

	"github.com/seerx/gpa/logger"
)

type RepoFile struct {
	info *Info     `json:"-"`
	File *ast.File `json:"-"`
	// Parsed  bool
	Name              string
	Path              string
	Package           string
	Imports           map[string]string // key: 简写包名称  value: 包全称
	ImportsReverseMap map[string]string // key: 包全称      value: 简写包名称
	Repos             []*RepoInterface

	SQLPackage     string // database/sql
	RunTimePackage string // github.com/seerx/gpa/rt
	DBUtilPackage  string // github.com/seerx/gpa/rt/dbutil
	logger         logger.GpaLogger
}

func NewRepoFile(info *Info, file *ast.File) *RepoFile {
	return &RepoFile{
		info:              info,
		File:              file,
		Path:              file.Name.Name,
		Imports:           map[string]string{},
		ImportsReverseMap: map[string]string{},
		logger:            info.logger,
	}
}

func (rf *RepoFile) AddRepo(repo *RepoInterface) {
	repo.repoFile = rf
	rf.Repos = append(rf.Repos, repo)
}

func (rf *RepoFile) AddSQLPackage() string {
	if rf.SQLPackage == "" {
		rf.SQLPackage = rf.addPackage("sql", "database/sql")
	}
	return rf.SQLPackage
}

// func (rf *RepoFile) AddContextPackage() string {
// 	if rf.ContextPackage == "" {
// 		rf.ContextPackage = rf.addPackage("context", "context")
// 	}
// 	return rf.ContextPackage
// }

func (rf *RepoFile) AddRuntimePackage() string {
	if rf.RunTimePackage == "" {
		rf.RunTimePackage = rf.addPackage("rt", "github.com/seerx/gpa/rt")
	}
	for _, intf := range rf.Repos {
		intf.RunTimePackage = rf.RunTimePackage
	}
	return rf.RunTimePackage
}

func (rf *RepoFile) AddDBUtilPackage() string {
	if rf.DBUtilPackage == "" {
		rf.DBUtilPackage = rf.addPackage("dbutil", "github.com/seerx/gpa/rt/dbutil")
	}
	return rf.DBUtilPackage
}

func (rf *RepoFile) AddXTypePackage(pkgPath string) string {
	pkgName := pkgPath
	lastIdx := strings.LastIndex(pkgPath, "/")
	if lastIdx > 0 {
		pkgName = pkgPath[lastIdx+1:]
	}
	pkgName = rf.addPackage(pkgName, pkgPath)
	return pkgName
}

func (rf *RepoFile) addPackage(pkgNamePrefix, pkg string) string {
	pkgName, ok := rf.ImportsReverseMap[pkg]
	if ok {
		return pkgName
	}
	// var pkgName string
	for {
		pkgName = fmt.Sprintf("%s%d", pkgNamePrefix, rand.Intn(1000))
		if _, ok := rf.Imports[pkgName]; ok {
			continue
		}
		rf.Imports[pkgName] = pkg
		rf.ImportsReverseMap[pkg] = pkgName
		break
	}
	return pkgName
}

func (rf *RepoFile) FindImport(pkg string) string {
	return rf.Imports[pkg]
}

func (rf *RepoFile) FindPackagePath(pkg string) (string, error) {
	impPkg := rf.FindImport(pkg)
	if impPkg == "" {
		return "", fmt.Errorf("package %s is not found", pkg)
	}
	pkgInfo, err := build.Import(impPkg, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return pkgInfo.Dir, nil
}

func (rf *RepoFile) Parse(dialect string) error {
	file := rf.File
	var err error
	var repo *RepoInterface
	ast.Inspect(file, func(node ast.Node) bool {
		// ast.Field
		switch t := node.(type) {
		// case *ast.CommentGroup:
		// 	fmt.Println(t.Text())
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
			rf.Imports[name] = path
			rf.ImportsReverseMap[path] = name
			// fmt.Println(t.Path.Value, t.Name)
		case *ast.TypeSpec:
			// repo 定义  type RepoXxx
			repo = nil
			if _, ok := t.Type.(*ast.InterfaceType); ok {
				repo = NewRepoInterface(t.Name.Name, rf.logger)
				rf.AddRepo(repo)
			}
		case *ast.InterfaceType:
			// repo 定义的实体内容(接口方法列表)
			if repo != nil {
				if err := repo.Parse(t, dialect); err != nil {
					return false
				}
			}
		}
		// if node != nil {
		// 	fmt.Println("node", node.Pos())
		// }
		return true
	})
	// })
	return err
}

// func parseObject(repo *RepoFile, repoIntf *RepoInterface, method *ast.Field) (*Object, error) {
// 	obj := NewObject(repoIntf) //  &metas.Object{}

// 	typ, ok := method.Type.(*ast.FuncType)
// 	if !ok {
// 		return nil, errors.New("not a valid method")
// 	}
// 	for _, p := range typ.Params.List {
// 		// arg := &metas.Object{Name: getName(p.Names)}
// 		// objs.NewObject()
// 		// objs.NewObject()
// 		// arg := metas.NewObject(repoIntf)
// 		// arg.Name = getName(p.Names)
// 		var err error
// 		arg, err = parseParam(repo, repoIntf, arg, &pp{expr: p.Type, field: p}, 0)
// 		if err != nil {
// 			return arg, err
// 		}
// 		obj.AddArg(arg)
// 		// obj.Args = append(obj.Args, arg)
// 		// fmt.Println("param", p.Names, p.Type)
// 	}

// 	// for _, name := range method.Names {
// 	obj.Name = getName(method.Names)
// 	// }
// 	if typ.Results != nil {
// 		for _, p := range typ.Results.List {
// 			// result := &metas.Object{Name: getName(p.Names)}
// 			result := metas.NewObject(repoIntf) // {Name: getName(p.Names)}
// 			result.Name = getName(p.Names)
// 			_, err := parseParam(repo, repoIntf, result, &pp{expr: p.Type, field: p}, 0)
// 			// result, err := parseParam(p)
// 			if err != nil {
// 				return result, err
// 			}
// 			// fd := p.(*ast.Field)
// 			// fd := p.Type.(*ast.StarExpr)
// 			// fmt.Println("result", p.Names, p.Type)
// 			// se := fd.X.(*ast.SelectorExpr)
// 			// fmt.Println("result", p.Names, p.Type, se.Sel.Name)
// 			obj.Results = append(obj.Results, result)
// 		}
// 	}

// 	return obj, nil
// }

// type pp struct {
// 	field *ast.Field
// 	expr  ast.Expr
// }

// func parseParam(repo *RepoFile, repoIntf *RepoInterface, arg *Object, p *pp, level int) (*Object, error) {
// 	// arg := &Object{Name: getName(p.Names)}
// 	var err error
// 	switch pt := p.expr.(type) {
// 	// case *ast.Im:
// 	// 	fmt.Println(pt.Path, pt.Name.Name)
// 	case *ast.Ident:
// 		// arg.Type = pkg.TypesInfo.ObjectOf(pt).Type()
// 		// arg.Type = metas.Type{Name: pt.Name}
// 		arg.Type = *objs.NewPrimitiveType(pt.Name)
// 		// arg.Type = *metas.NewPrimitiveType(arg, pt.Name)
// 		// fmt.Println(pt.Name)
// 	case *ast.FuncType:
// 		if level > 1 {
// 			return nil, errors.New("不支持多层嵌套函数类型")
// 		}
// 		arg, err = parseObject(repo, repoIntf, p.field)
// 		if err != nil {
// 			return nil, err
// 		}
// 		// arg.Type = metas.Type{Name: "func"} // pkg.TypesInfo.ObjectOf(p.field.Names[0]).Type()
// 		arg.Type = *objs.NewFuncType()
// 	case *ast.SelectorExpr:
// 		arg.Type = *parseSelectorExprType(pt, false)
// 	case *ast.StarExpr:
// 		arg.Type = *parseSelectorExprType(pt.X.(*ast.SelectorExpr), true)
// 	case *ast.ArrayType:
// 		if level > 1 {
// 			return nil, errors.New("不支持多层嵌套数据类型")
// 		}
// 		if _, err := parseParam(repo, repoIntf, arg, &pp{expr: pt.Elt}, level+1); err != nil {
// 			return nil, err
// 		}
// 		// arg.IsSlice()
// 		arg.Slice = true
// 	case *ast.SliceExpr:
// 		if level > 1 {
// 			return nil, errors.New("不支持多层嵌套数据类型")
// 		}
// 		if _, err := parseParam(repo, repoIntf, arg, &pp{expr: pt.X}, level+1); err != nil {
// 			return nil, err
// 		}
// 		arg.Slice = true
// 	case *ast.MapType:
// 		if level > 1 {
// 			return nil, errors.New("不支持多层嵌套数据类型")
// 		}
// 		// pt.Key
// 		if _, err := parseParam(repo, repoIntf, arg, &pp{expr: pt.Value}, level+1); err != nil {
// 			return nil, err
// 		}
// 		o := metas.NewObject(nil)
// 		// o := &metas.Object{}
// 		if _, err := parseParam(repo, repoIntf, o, &pp{expr: pt.Key}, level+1); err != nil {
// 			return nil, err
// 		}
// 		arg.KeyType = &o.Type
// 		arg.Map = true
// 	default:
// 		err = errors.New("不支持的数据类型")
// 	}
// 	// fmt.Println(arg.Type.String())
// 	return arg, err
// }
