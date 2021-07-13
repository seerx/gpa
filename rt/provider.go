package rt

import (
	"context"
	"database/sql"
	"time"

	"github.com/seerx/gpa/engine/sql/dialect/intf"
	"github.com/seerx/gpa/engine/sql/types"
	"github.com/seerx/gpa/logger"
	"github.com/seerx/gpa/rt/exec"
	"github.com/seerx/logo/log"
)

type Provider struct {
	ctx            context.Context
	dialect        intf.Dialect
	db             *sql.DB
	exe            exec.SQLExecutor
	tx             *exec.TXExecutor
	transactioning bool
	timezone       *time.Location
	logger         logger.GpaLogger
}

func NewProvider(ctx context.Context,
	dialect intf.Dialect,
	db *sql.DB,
	timezone *time.Location,
	logger logger.GpaLogger) *Provider {
	return &Provider{
		ctx:      ctx,
		dialect:  dialect,
		db:       db,
		timezone: timezone,
		exe:      exec.NewExecutor(db, logger),
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
	if p.dialect.URI().DBType == types.MSSQL {
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
