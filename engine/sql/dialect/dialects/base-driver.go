package dialects

import (
	"fmt"

	"github.com/seerx/gpa/engine/constants"
	"github.com/seerx/gpa/engine/sql/dialect/intf"
)

type baseDriver struct {
	// Dialect       intf.Dialect
	fnMakeDialect func() intf.Dialect
}

func (b *baseDriver) GetDialect() intf.Dialect {
	return b.fnMakeDialect()
}

func (b *baseDriver) uri(dialect constants.DIALECT) (*intf.URI, error) {
	dbt := dialect.GetDRIVER()
	if dbt == constants.DB_UNKNOWN {
		return nil, fmt.Errorf("no db type surport od dialcet %s", dialect)
	}
	return &intf.URI{DRIVER: dbt}, nil
}
