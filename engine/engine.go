package engine

import (
	"time"

	"github.com/seerx/gpa/engine/sql/dialect/intf"
	"github.com/seerx/gpa/engine/sql/metas/rflt"
	"github.com/seerx/gpa/rt"
)

const tagName = "gpa"

type Engine struct {
	// db          *sql.DB
	provider    rt.Provider
	dialect     intf.Dialect
	propsParser *rflt.PropsParser
	TZLocation  *time.Location // The timezone of the application
	DatabaseTZ  *time.Location // The timezone of the database
}
