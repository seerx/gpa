package defines

import (
	"errors"
	"fmt"
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
	repoIntf *RepoInterface
	Object
	Template Template
	SQL      string
}

func NewFunc(repoIntf *RepoInterface, obj *Object, sql string) *Func {
	return &Func{repoIntf: repoIntf, SQL: sql, Object: *obj}
}

func (f *Func) CreateError(format string, v ...interface{}) error {
	return errors.New(f.Format(format, v...))
}

func (f *Func) Format(format string, v ...interface{}) string {
	return fmt.Sprintf("%s.%s %s\n%s", f.repoIntf.Name, f.Name, f.repo.repoFile.Path, fmt.Sprintf(format, v...))
}
