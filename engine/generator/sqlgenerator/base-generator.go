package sqlgenerator

import (
	"fmt"
	"strings"
)

type baseGen struct {
}

func (bd *baseGen) QuoteExpr(sql string) string {
	return strings.ReplaceAll(sql, "\"", "\\\"")
}

func (bd *baseGen) Insert(sql *SQL) (string, []*SQLParam) {
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		sql.TableName,
		strings.Join(sql.Columns, ","),
		strings.Join(sql.ParamPlaceHolder, ",")), sql.Params
}

func (bd *baseGen) Update(sql *SQL) (string, []*SQLParam, []*SQLParam) {
	// quoter := bd.Dialect.Quoter()
	if sql.Where != "" {
		sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
			sql.TableName,
			strings.Join(sql.Columns, ","),
			sql.Where)
		return bd.QuoteExpr(sqlStr), sql.Params, sql.WhereParams
	}
	sqlStr := fmt.Sprintf("UPDATE %s SET %s",
		sql.TableName,
		strings.Join(sql.Columns, ","))
	return bd.QuoteExpr(sqlStr), sql.Params, nil
}
func (bd *baseGen) Delete(sql *SQL) (string, []*SQLParam) {
	// quoter := bd.Dialect.Quoter()
	if sql.Where != "" {
		sqlStr := fmt.Sprintf("DELETE FROM %s WHERE %s",
			sql.TableName,
			sql.Where)
		return bd.QuoteExpr(sqlStr), sql.WhereParams
	}
	sqlStr := fmt.Sprintf("DELETE FROM %s",
		sql.TableName)
	return bd.QuoteExpr(sqlStr), sql.Params
}

func (bd *baseGen) Query(sql *SQL) (string, []*SQLParam) {
	if sql.Where != "" {
		sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
			strings.Join(sql.Columns, ","),
			sql.TableName,
			sql.Where)
		return bd.QuoteExpr(sqlStr), sql.WhereParams
	}
	sqlStr := fmt.Sprintf("SELECT %s FROM %s",
		strings.Join(sql.Columns, ","),
		sql.TableName)
	return bd.QuoteExpr(sqlStr), sql.WhereParams
}
