package engine

import (
	"time"

	"github.com/seerx/gpa/logger"
)

func (e *Engine) SetLogSQL(log bool) *Engine {
	e.logger.SetLogSQL(log)
	return e
}

func (e *Engine) SetMaxIdleConns(n int) *Engine {
	e.db.SetMaxIdleConns(n)
	return e
}

func (e *Engine) SetMaxOpenConns(n int) *Engine {
	e.db.SetMaxOpenConns(n)
	return e
}

func (e *Engine) SetConnMaxIdleTime(d time.Duration) *Engine {
	e.db.SetConnMaxIdleTime(d)
	return e
}

func (e *Engine) SetConnMaxLifetime(d time.Duration) *Engine {
	e.db.SetConnMaxLifetime(d)
	return e
}

func (e *Engine) SetLogger(log logger.GpaLogger) {
	e.logger = log
}
