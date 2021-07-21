package method

import (
	"github.com/seerx/gpa/engine/generator/defines"
	rdesc "github.com/seerx/gpa/engine/generator/repo-desc"
	"github.com/seerx/gpa/engine/generator/sqlgenerator"
	"github.com/seerx/gpa/logger"
)

type Method interface {
	Test(fn *defines.Func) bool
	Parse() (*rdesc.FuncDesc, error)
}

// var methods []Method

func CreateMethods(sqlg sqlgenerator.SQLGenerator, logger logger.GpaLogger) []Method {
	return []Method{
		&insert{BaseMethod: BaseMethod{sqlg: sqlg, logger: logger}},
		&updateby{BaseMethod: BaseMethod{sqlg: sqlg, logger: logger}}, // updateby 排在 update 之前，优先考虑 updateby 操作
		&update{BaseMethod: BaseMethod{sqlg: sqlg, logger: logger}},
		&deleteby{BaseMethod: BaseMethod{sqlg: sqlg, logger: logger}},
		&delete{BaseMethod: BaseMethod{sqlg: sqlg, logger: logger}},
		&findby{BaseMethod: BaseMethod{sqlg: sqlg, logger: logger}},
		&find{BaseMethod: BaseMethod{sqlg: sqlg, logger: logger}},
		&countby{BaseMethod: BaseMethod{sqlg: sqlg, logger: logger}},
		&count{BaseMethod: BaseMethod{sqlg: sqlg, logger: logger}},
	}
}

// func GetMethod(fn *defines.Func) Method {
// 	for _, g := range methods {
// 		if g.Test(fn) {
// 			return g
// 		}
// 	}
// 	return nil
// }
