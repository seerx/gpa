package exec

import (
	"context"
	"database/sql"

	"github.com/seerx/gpa/logger"
)

type Executor struct {
	db     *sql.DB
	logger logger.GpaLogger
}

func NewExecutor(db *sql.DB, logger logger.GpaLogger) *Executor {
	return &Executor{db: db, logger: logger}
}

func (e *Executor) QueryRow(sql string, args ...interface{}) *sql.Row {
	if e.logger.IsLogSQL() {
		e.logger.Info(sql, args)
	}
	return e.db.QueryRow(sql, args...)
}

func (e *Executor) QueryRows(sql string, args ...interface{}) (*sql.Rows, error) {
	if e.logger.IsLogSQL() {
		e.logger.Info(sql, args)
	}
	return e.db.Query(sql, args...)
}

func (e *Executor) QueryContextRow(ctx context.Context, sql string, args ...interface{}) *sql.Row {
	if e.logger.IsLogSQL() {
		e.logger.Info(sql, args)
	}
	return e.db.QueryRowContext(ctx, sql, args...)
}

func (e *Executor) QueryContextRows(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error) {
	if e.logger.IsLogSQL() {
		e.logger.Info(sql, args)
	}
	return e.db.QueryContext(ctx, sql, args...)
}

func (e *Executor) Exec(sql string, args ...interface{}) (sql.Result, error) {
	if e.logger.IsLogSQL() {
		e.logger.Info(sql, args)
	}
	return e.db.Exec(sql, args...)
}

func (e *Executor) ExecContext(ctx context.Context, sql string, args ...interface{}) (sql.Result, error) {
	if e.logger.IsLogSQL() {
		e.logger.Info(sql, args)
	}
	return e.db.ExecContext(ctx, sql, args...)
}
