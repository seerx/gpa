package exec

import (
	"context"
	"database/sql"
)

type SQLExecutor interface {
	QueryRow(sql string, args ...interface{}) *sql.Row
	QueryRows(sql string, args ...interface{}) (*sql.Rows, error)
	QueryContextRow(ctx context.Context, sql string, args ...interface{}) *sql.Row
	QueryContextRows(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error)

	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// type Queryer interface {
// 	QueryRow(sql string, args ...interface{}) *sql.Row
// 	QueryRows(sql string, args ...interface{}) (*sql.Rows, error)
// 	QueryContextRow(ctx context.Context, sql string, args ...interface{}) *sql.Row
// 	QueryContextRows(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error)
// }

// // Executer represents an interface to execute a SQL
// type Executer interface {
// 	Exec(query string, args ...interface{}) (sql.Result, error)
// 	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
// }

// // QueryExecuter combines the Queryer and Executer
// type QueryExecuter interface {
// 	Queryer
// 	Executer
// }
