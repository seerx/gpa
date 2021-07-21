package rt

import (
	"context"
	"database/sql"
	"time"

	"github.com/seerx/gpa/engine/constants"
	"github.com/seerx/gpa/logger"
	"github.com/seerx/gpa/rt/exec"
	"github.com/seerx/logo/log"
)

type Provider struct {
	ctx    context.Context
	driver constants.DRIVER
	// dialect        intf.Dialect
	db             *sql.DB
	exe            exec.SQLExecutor
	tx             *exec.TXExecutor
	transactioning bool
	timezone       *time.Location
	logger         logger.GpaLogger
}

type Option struct {
	Conext   context.Context
	Timezone *time.Location
	Logger   logger.GpaLogger
}

func NewProvider(driver constants.DRIVER, db *sql.DB, opt *Option) *Provider {
	if opt == nil {
		opt = &Option{
			Conext:   context.Background(),
			Timezone: time.Local,
			Logger:   logger.GetLogger(),
		}
	} else {
		if opt.Conext == nil {
			opt.Conext = context.Background()
		}
		if opt.Timezone == nil {
			opt.Timezone = time.Local
		}
		if opt.Logger == nil {
			opt.Logger = logger.GetLogger()
		}
	}
	return &Provider{
		ctx:      opt.Conext,
		driver:   driver,
		db:       db,
		timezone: opt.Timezone,
		logger:   opt.Logger,
		exe:      exec.NewExecutor(db, opt.Logger),
	}
}

func (p *Provider) Executor() exec.SQLExecutor {
	if p.transactioning {
		return p.tx
	}
	return p.exe
}

func (p *Provider) GetTimezone() *time.Location {
	if p.timezone == nil {
		p.timezone = time.Local
	}
	return p.timezone
}

func (p *Provider) GetTimeStampzFormat() string {
	// if p.dialect.URI().DRIVER == constants.DB_MSSQL {
	// 	return "2006-01-02T15:04:05.9999999Z07:00"
	// }
	if p.driver == constants.DB_MSSQL {
		return "2006-01-02T15:04:05.9999999Z07:00"
	}
	return ""
}

func (p *Provider) Transaction(fn func() error) (err error) {
	p.tx, err = exec.NewTXExecutor(p.ctx, p.db, p.logger)
	if err != nil {
		return err
	}
	p.transactioning = true
	defer func() {
		if err != nil {
			if er := p.tx.Rollback(); er != nil {
				log.Error(er)
			}
		} else {
			if er := p.tx.Commit(); er != nil {
				log.Error(er)
			}
		}
		p.transactioning = false
		p.tx = nil
	}()
	err = fn()
	return
}
