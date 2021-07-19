package sqlgenerator

import (
	"fmt"

	"github.com/seerx/gpa/engine/constants"
)

type pgsql struct {
	baseGen
}

func init() {
	register(constants.POSTGRES, &pgsql{})
}

func (p *pgsql) Insert(sql *SQL) (string, []*SQLParam) {
	s, ps := p.baseGen.Insert(sql)
	if sql.ReturnAutoincrPrimaryKey != "" {
		s += fmt.Sprintf(" RETURNING %s", sql.ReturnAutoincrPrimaryKey)
	}
	return s, ps
}
