package defines

import (
	"go/ast"
	"strings"

	"github.com/seerx/gpa/engine/objs"
)

func GetName(names []*ast.Ident) string {
	for _, name := range names {
		return name.Name
	}
	return ""
}

func ParseSQL(comment, dialect string) string {
	lines := strings.Split(comment, "\n")
	sql := ""
	dialect += ":"
	for _, line := range lines {
		if strings.Index(line, dialect) == 0 {
			sql = line[len(dialect):]
			break
		}
		if strings.Index(line, "sql:") == 0 {
			sql = line[len("sql:"):]
		}
	}
	return sql
}

func ParseSelectorExprType(se *ast.SelectorExpr, ptr bool) *objs.Type {
	// name := ""
	x, ok := se.X.(*ast.Ident)
	pkg := ""
	if ok {
		pkg = x.Name
	}
	if ptr {
		return objs.NewPtrType(pkg, se.Sel.Name)
	}
	return objs.NewType(pkg, se.Sel.Name)
}
