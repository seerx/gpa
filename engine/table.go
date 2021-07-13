package engine

import (
	"errors"

	"github.com/seerx/gpa/engine/sql/metas/rflt"
	"github.com/seerx/gpa/engine/sql/metas/schema"
	"github.com/seerx/gpa/engine/sql/types"
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

func (e *Engine) CreateTable(table *schema.Table) error {
	sqls, err := e.dialect.SQLCreateTable(table, "")
	if err != nil {
		return err
	}

	for _, sql := range sqls {
		if _, err := e.provider.Executor().Exec(sql); err != nil {
			return err
		}
	}

	return nil
}

func (e *Engine) dropTable(table *schema.Table) error {
	sql, err := e.dialect.SQLDropTable(table.Name)
	if err != nil {
		return err
	}
	_, err = e.provider.Executor().Exec(sql)
	return err
}

func (e *Engine) DropTable(models ...interface{}) error {
	for _, model := range models {
		table, err := rflt.Parse(model, e.propsParser)
		if err != nil {
			return err
		}
		if err := e.dropTable(table); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) CreateIndex(table *schema.Table) error {
	for _, index := range table.Indexes {
		if index.Type == types.IndexType {
			sql := e.dialect.SQLCreateIndex(table.Name, index)
			if _, err := e.provider.Executor().Exec(sql); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Engine) CreateUnique(table *schema.Table) error {
	for _, index := range table.Indexes {
		if index.Type == types.UniqueType {
			sql := e.dialect.SQLCreateIndex(table.Name, index)
			if _, err := e.provider.Executor().Exec(sql); err != nil {
				return err
			}
		}
	}
	return nil
}

// DropIndexes drop indexes
func (e *Engine) DropIndexes(table *schema.Table) error {
	for _, index := range table.Indexes {
		sql := e.dialect.SQLDropIndex(table.Name, index)
		if _, err := e.provider.Executor().Exec(sql); err != nil {
			return err
		}
	}
	return nil
}
