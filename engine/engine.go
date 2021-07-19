package engine

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/seerx/gpa/engine/constants"
	"github.com/seerx/gpa/engine/generator/method"
	"github.com/seerx/gpa/engine/generator/sqlgenerator"
	"github.com/seerx/gpa/engine/sql/dialect"
	"github.com/seerx/gpa/engine/sql/dialect/intf"
	"github.com/seerx/gpa/engine/sql/metas/rflt"
	"github.com/seerx/gpa/logger"
	"github.com/seerx/gpa/rt"
)

const TagName = "gpa"

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

// func New

func New(dialectName constants.DIALECT, source string) (e *Engine, err error) {
	driver := dialectName.GetDRIVER()
	if driver == constants.DB_UNKNOWN {
		return nil, fmt.Errorf("unkown databse driver of dialect %s", dialectName)
	}
	log := logger.GetLogger()
	dial, err := dialect.OpenDialect(dialectName, source)
	if err != nil {
		return nil, err
	}

	sqlg, err := sqlgenerator.GetGenerator(dialectName)
	if err != nil {
		return nil, err
	}
	method.InitMethods(sqlg, log)

	propsParser := rflt.NewPropsParser(TagName, dial)
	db, err := sql.Open(string(driver), source)
	if err != nil {
		log.WithError(err).Error("connect database error")
		return nil, err
	}

	prvd := rt.NewProvider(context.Background(), driver, db, time.Local, log)

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
