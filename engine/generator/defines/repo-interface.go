package defines

import (
	"go/ast"

	"github.com/seerx/gpa/logger"
)

type RepoInterface struct {
	repoFile       *RepoFile
	Name           string
	Funcs          []*Func
	RunTimePackage string
	logger         logger.GpaLogger
}

func NewRepoInterface(name string, logger logger.GpaLogger) *RepoInterface {
	return &RepoInterface{
		Name:   name,
		logger: logger,
	}
}

func (rf *RepoInterface) AddFunction(fn *Func) {
	rf.Funcs = append(rf.Funcs, fn)
}

func (rf *RepoInterface) Parse(body *ast.InterfaceType, dialect string) error {
	for _, method := range body.Methods.List {
		fn := NewFunc(rf)
		if err := fn.Parse(method, dialect); err != nil {
			rf.logger.Errorf(err, "%s.%s is invalid", rf.Name, fn.Name)
			return err
		}

		rf.AddFunction(fn)
		// obj, er := parseObject(rf, repo, method)
		// if er != nil {
		// 	err = er
		// 	log.Errorf("%s.%s 定义不合法: %v", repo.Name, obj.Name, err)
		// 	return false
		// }
		// fmt.Println(method.Doc.Text())
		// repo.Funcs = append(repo.Funcs,
		// 	NewFunc(repo, obj, parseSQL(method.Doc.Text(), dialect)))
		// fmt.Println(field.Comment.Text())
	}
	return nil
}
