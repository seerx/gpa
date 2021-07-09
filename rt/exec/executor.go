package exec

import (
	"context"
	"database/sql"

	"github.com/seerx/mro/log"
)

type Executor struct {
	db *sql.DB
}

func NewExecutor(db *sql.DB) *Executor {
	return &Executor{db: db}
}

func (e *Executor) QueryRow(sql string, args ...interface{}) *sql.Row {
	if log.IsPrintSQL() {
		log.Info(sql, args)
	}
	return e.db.QueryRow(sql, args...)
}

func (e *Executor) QueryRows(sql string, args ...interface{}) (*sql.Rows, error) {
	if log.IsPrintSQL() {
		log.Info(sql, args)
	}
	return e.db.Query(sql, args...)
}

func (e *Executor) QueryContextRow(ctx context.Context, sql string, args ...interface{}) *sql.Row {
	if log.IsPrintSQL() {
		log.Info(sql, args)
	}
	return e.db.QueryRowContext(ctx, sql, args...)
}

func (e *Executor) QueryContextRows(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error) {
	if log.IsPrintSQL() {
		log.Info(sql, args)
	}
	return e.db.QueryContext(ctx, sql, args...)
}

func (e *Executor) Exec(sql string, args ...interface{}) (sql.Result, error) {
	if log.IsPrintSQL() {
		log.Info(sql, args)
	}
	return e.db.Exec(sql, args...)
}

func (e *Executor) ExecContext(ctx context.Context, sql string, args ...interface{}) (sql.Result, error) {
	if log.IsPrintSQL() {
		log.Info(sql, args)
	}
	return e.db.ExecContext(ctx, sql, args...)
}
