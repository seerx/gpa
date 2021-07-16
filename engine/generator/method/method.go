package method

import (
	"github.com/seerx/gpa/engine/generator/defines"
	rdesc "github.com/seerx/gpa/engine/generator/repo-desc"
	"github.com/seerx/gpa/engine/sql/dialect/intf"
	"github.com/seerx/gpa/logger"
)

type Method interface {
	Test(fn *defines.Func) bool
	Parse() (*rdesc.FuncDesc, error)
}

var methods []Method

func InitMethods(d intf.Dialect, logger logger.GpaLogger) {
	methods = []Method{
		&insert{BaseMethod: BaseMethod{dialect: d, logger: logger}},
		&updateby{BaseMethod: BaseMethod{dialect: d, logger: logger}}, // updateby 排在 update 之前，优先考虑 updateby 操作
		&update{BaseMethod: BaseMethod{dialect: d, logger: logger}},
		&deleteby{BaseMethod: BaseMethod{dialect: d, logger: logger}},
		&delete{BaseMethod: BaseMethod{dialect: d, logger: logger}},
		// &findby{BaseGenerator: BaseGenerator{dialect: d, logger: logger}},
		// &find{BaseGenerator: BaseGenerator{dialect: d, logger: logger}},
		&countby{BaseMethod: BaseMethod{dialect: d, logger: logger}},
		&count{BaseMethod: BaseMethod{dialect: d, logger: logger}},
	}
}

func GetMethod(fn *defines.Func) Method {
	for _, g := range methods {
		if g.Test(fn) {
			return g
		}
	}
	return nil
}
