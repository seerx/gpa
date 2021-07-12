package engine

import (
	"errors"

	"github.com/seerx/gpa/engine/sql/metas/schema"
)

func (e *Engine) HasTable(table *schema.Table) error {
	sql, args := e.dialect.SQLTableExists(table.Name)
	row := e.provider.Executor().QueryRow(sql, args...)
	var tbName string
	if err := row.Scan(&tbName); err != nil {
		return err
	}
	if tbName != table.Name {
		return errors.New("table is not exists")
	}
	return nil
}
