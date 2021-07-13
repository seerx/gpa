package engine

import (
	"context"
	"database/sql"
	"time"

	"github.com/seerx/gpa/engine/sql/dialect"
	"github.com/seerx/gpa/engine/sql/dialect/intf"
	"github.com/seerx/gpa/engine/sql/metas/rflt"
	"github.com/seerx/gpa/logger"
	"github.com/seerx/gpa/rt"
	"github.com/seerx/logo/log"
)

const tagName = "gpa"

type Engine struct {
	db          *sql.DB
	provider    rt.Provider
	dialect     intf.Dialect
	propsParser *rflt.PropsParser
	TZLocation  *time.Location // The timezone of the application
	DatabaseTZ  *time.Location // The timezone of the database
	logger      logger.GpaLogger
}

func (e *Engine) GetProvider() *rt.Provider {
	return &e.provider
}

func New(driver, source string) (e *Engine, err error) {
	dial, err := dialect.OpenDialect(driver, source)
	if err != nil {
		return nil, err
	}

	propsParser := rflt.NewPropsParser(tagName, dial)
	db, err := sql.Open(driver, source)
	if err != nil {
		log.WithError(err).Error("connect database error")
		return nil, err
	}
	log := logger.GetLogger()
	prvd := rt.NewProvider(context.Background(), dial, db, time.Local, log)

	e = &Engine{
		db:          db,
		provider:    *prvd,
		dialect:     dial,
		propsParser: propsParser,
		DatabaseTZ:  time.Local,
		logger:      log,
	}
	return
}
