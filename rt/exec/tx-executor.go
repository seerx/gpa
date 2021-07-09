package exec

import (
	"context"
	"database/sql"

	"github.com/seerx/mro/log"
)

type TXExecutor struct {
	tx *sql.Tx
}

func NewTXExecutor(ctx context.Context, db *sql.DB) (*TXExecutor, error) {
	tx, err := db.BeginTx(ctx, nil)
	return &TXExecutor{tx: tx}, err
}

func (tx *TXExecutor) Rollback() error {
	return tx.tx.Rollback()
}

func (tx *TXExecutor) Commit() error {
	return tx.tx.Commit()
}

func (tx *TXExecutor) QueryRow(sql string, args ...interface{}) *sql.Row {
	if log.IsPrintSQL() {
		log.Info(sql, args)
	}
	return tx.tx.QueryRow(sql, args...)
}

func (tx *TXExecutor) QueryRows(sql string, args ...interface{}) (*sql.Rows, error) {
	if log.IsPrintSQL() {
		log.Info(sql, args)
	}
	return tx.tx.Query(sql, args...)
}

func (tx *TXExecutor) QueryContextRow(ctx context.Context, sql string, args ...interface{}) *sql.Row {
	if log.IsPrintSQL() {
		log.Info(sql, args)
	}
	return tx.tx.QueryRowContext(ctx, sql, args...)
}

func (tx *TXExecutor) QueryContextRows(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error) {
	if log.IsPrintSQL() {
		log.Info(sql, args)
	}
	return tx.tx.QueryContext(ctx, sql, args...)
}

func (tx *TXExecutor) Exec(sql string, args ...interface{}) (sql.Result, error) {
	if log.IsPrintSQL() {
		log.Info(sql, args)
	}
	return tx.tx.Exec(sql, args...)
}

func (tx *TXExecutor) ExecContext(ctx context.Context, sql string, args ...interface{}) (sql.Result, error) {
	if log.IsPrintSQL() {
		log.Info(sql, args)
	}
	return tx.tx.ExecContext(ctx, sql, args...)
}
