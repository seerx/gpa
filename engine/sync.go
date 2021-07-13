package engine

import (
	"context"
	"fmt"
	"strings"

	"github.com/seerx/gpa/engine/sql/metas/rflt"
	"github.com/seerx/gpa/engine/sql/metas/schema"
	"github.com/seerx/gpa/engine/sql/types"
	"github.com/seerx/logo/log"
)

func (e *Engine) Sync(beans ...interface{}) error {
	dbTables, err := e.dialect.GetTables(e.provider.Executor(), context.Background())
	if err != nil {
		e.logger.Error(err, "get tables from database")
		return err
	}

	for _, bean := range beans {
		table, err := rflt.Parse(bean, e.propsParser)
		if err != nil {
			e.logger.Error(err, "parse table bean")
			return err
		}
		tbName := table.Name
		tbNameWithSchema := e.dialect.TableNameWithSchema(tbName)

		var dbTable *schema.Table
		// 从 tables 中查找 table 是否存在
		for _, tb := range dbTables {
			if strings.EqualFold(e.dialect.TableNameWithSchema(tb.Name), tbNameWithSchema) {
				dbTable = tb
				break
			}
		}
		if dbTable == nil {
			// 创建表
			if err := e.CreateTable(table); err != nil {
				e.logger.Errorf(err, "create table [%s]", tbName)
				return err
			}
			break
		}
		// 从数据库加载表结构
		if err := e.loadTableInfo(dbTable); err != nil {
			e.logger.Errorf(err, "load table [%s] info from database", tbName)
			return err
		}

		// 检查列变化
		if err := e.checkColumns(tbNameWithSchema, dbTable, table); err != nil {
			return err
		}
		// 检查索引变化
		if err := e.checkIndexes(tbNameWithSchema, dbTable, table); err != nil {
			return err
		}

		// 检查所有从结构体中移除的字段，但在数据库中仍然存在
		for _, colName := range dbTable.ColumnNames {
			if table.GetColumn(colName) == nil {
				log.Warnf("Table %s has column %s but struct has not related field", tbNameWithSchema, colName)
			}
		}
	}

	return nil
}

func (e *Engine) checkIndexes(tbNameWithSchema string, oriTable, table *schema.Table) error {
	var existsIndexInDB = make(map[string]bool)
	var indexesTobeCreate = make(map[string]*schema.Index)
	// 遍历结构体定义的索引
	for name, index := range table.Indexes {
		var oriIndex *schema.Index
		// 查找表中的索引
		for name2, index2 := range oriTable.Indexes {
			if index.Equal(index2) {
				// 找到索引
				oriIndex = index2
				existsIndexInDB[name2] = true
				break
			}
		}

		// 索引存在
		// if oriIndex != nil {
		// 	if oriIndex.Type != index.Type {
		// 		sql := s.dialect.SQLDropIndex(tbNameWithSchema, oriIndex)
		// 		if _, err := s.Exec(sql); err != nil {
		// 			return err
		// 		}
		// 		oriIndex = nil
		// 	}
		// }
		// 索引不存在
		if oriIndex == nil {
			indexesTobeCreate[name] = index
		}
	}

	// 遍历数据库中的索引
	for name2, index2 := range oriTable.Indexes {
		if _, ok := existsIndexInDB[name2]; !ok {
			// 该索引在结构体中不存在，删除索引
			sql := e.dialect.SQLDropIndex(tbNameWithSchema, index2)
			if _, err := e.provider.Executor().Exec(sql); err != nil {
				e.logger.Error(err, sql)
				return err
			}
		}
	}

	for _, index := range indexesTobeCreate {
		// s.SetTable(table)
		sql := e.dialect.SQLCreateIndex(tbNameWithSchema, index)
		if _, err := e.provider.Executor().Exec(sql); err != nil {
			e.logger.Error(err, sql)
			return err
		}
		// if index.Type == types.UniqueType {
		// 	s.SetTable(table)
		// 	sql := s.dialect.SQLCreateIndex(tbNameWithSchema, index)
		// 	if _, err := s.Exec(sql); err != nil {
		// 		return err
		// 	}
		// } else if index.Type == types.IndexType {
		// 	s.SetTable(table)
		// 	sql := s.dialect.SQLCreateIndex(tbNameWithSchema, index)
		// 	if _, err := s.Exec(sql); err != nil {
		// 		return err
		// 	}
		// }
	}

	return nil
}

func (e *Engine) addColumn(tbNameWithSchema string, col *schema.Column) error {
	// col := s.refTable.GetColumn(colName)
	sql := e.dialect.SQLAddColumn(tbNameWithSchema, col)
	_, err := e.provider.Executor().Exec(sql)
	if err != nil {
		e.logger.Errorf(err, "exec sql: %s", sql)
	}
	return err
}

func (e *Engine) checkColumns(tbNameWithSchema string, oriTable, table *schema.Table) error {
	// check columns
	for _, col := range table.Columns {
		var dbCol *schema.Column
		for _, col2 := range oriTable.Columns {
			if strings.EqualFold(col.FieldName(), col2.FieldName()) {
				dbCol = col2
				break
			}
		}

		// 表中不存在此列，添加列
		if dbCol == nil {
			// s.SetTable(table)
			if err := e.addColumn(tbNameWithSchema, col); err != nil {
				return err
			}
			continue
		}

		// 数据类型检查
		var err error
		expectedType := e.dialect.SQLType(col)
		curType := e.dialect.SQLType(dbCol)
		if expectedType != curType {
			// 数据类型不一致
			if expectedType == types.Text &&
				strings.HasPrefix(curType, types.Varchar) {
				// 从 varchar 改为 text
				// currently only support mysql & postgres
				if e.dialect.URI().SupportColumnVarchar2Text() {
					e.logger.Infof("Table %s column %s change type from %s to %s\n",
						tbNameWithSchema, col.FieldName(), curType, expectedType)
					_, err = e.provider.Executor().Exec(e.dialect.SQLModifyColumn(tbNameWithSchema, col))
				} else {
					e.logger.Warnf("Table %s column %s db type is %s, struct type is %s\n",
						tbNameWithSchema, col.FieldName(), curType, expectedType)
				}
			} else if strings.HasPrefix(curType, types.Varchar) && strings.HasPrefix(expectedType, types.Varchar) {
				// varchar 数据长度不一致，在数据库支持的情况下，只允许增加长度
				if e.dialect.URI().SupportColumnVarcharIncLength() &&
					dbCol.Length < col.Length {
					e.logger.Infof("table %s column %s change type from varchar(%d) to varchar(%d)\n",
						tbNameWithSchema, col.FieldName(), dbCol.Length, col.Length)
					_, err = e.provider.Executor().Exec(e.dialect.SQLModifyColumn(tbNameWithSchema, col))
				} else {
					e.logger.Warnf("table %s column %s db type is varchar(%d), struct type is varchar(%d)\n",
						tbNameWithSchema, col.FieldName(), dbCol.Length, col.Length)
				}
			} else {
				// 不支持的类型转换
				if !(strings.HasPrefix(curType, expectedType) && curType[len(expectedType)] == '(') {
					e.logger.Error(fmt.Errorf("table %s column %s db type is %s, struct type is %s",
						tbNameWithSchema, col.FieldName(), curType, expectedType))
				}
			}
		} else if expectedType == types.Varchar {
			// 都是 varchar 类型
			if e.dialect.URI().SupportColumnVarcharIncLength() &&
				dbCol.Length < col.Length {
				e.logger.Infof("table %s column %s change type from varchar(%d) to varchar(%d)\n",
					tbNameWithSchema, col.FieldName(), dbCol.Length, col.Length)
				_, err = e.provider.Executor().Exec(e.dialect.SQLModifyColumn(tbNameWithSchema, col))
			} else {
				e.logger.Warnf("table %s column %s db type is varchar(%d), struct type is varchar(%d)\n",
					tbNameWithSchema, col.FieldName(), dbCol.Length, col.Length)
			}
		}
		// 默认值发生变化
		if col.Default != dbCol.Default {
			switch {
			case col.IsAutoIncrement: // For autoincrement column, don't check default
			case (col.Type.Name == types.Bool || col.Type.Name == types.Boolean) &&
				((strings.EqualFold(col.Default, "true") && dbCol.Default == "1") ||
					(strings.EqualFold(col.Default, "false") && dbCol.Default == "0")):
			default:
				e.logger.Warnf("Table %s Column %s db default is %s, struct default is %s",
					tbNameWithSchema, col.FieldName(), dbCol.Default, col.Default)
			}
		}
		if col.Nullable != dbCol.Nullable {
			e.logger.Warnf("Table %s Column %s db nullable is %v, struct nullable is %v",
				tbNameWithSchema, col.FieldName(), dbCol.Nullable, col.Nullable)
		}

		if err != nil {
			e.logger.Errorf(err, "check table %s columns", tbNameWithSchema)
			return err
		}
	}
	return nil
}

func (e *Engine) loadTableInfo(table *schema.Table) error {
	colSeq, cols, err := e.dialect.GetColumns(e.provider.Executor(), context.Background(), table.Name)
	if err != nil {
		e.logger.Errorf(err, "get table %s columns", table.Name)
		return err
	}
	for _, name := range colSeq {
		table.AddColumn(cols[name])
	}
	indexes, err := e.dialect.GetIndexes(e.provider.Executor(), context.Background(), table.Name)
	if err != nil {
		e.logger.Errorf(err, "get table %s indexes", table.Name)
		return err
	}
	table.Indexes = indexes

	var seq int
	for _, index := range indexes {
		for _, name := range index.Cols {
			parts := strings.Split(strings.TrimSpace(name), " ")
			if len(parts) > 1 {
				if parts[1] == "DESC" {
					seq = 1
				} else if parts[1] == "ASC" {
					seq = 0
				}
			}
			var colName = strings.Trim(parts[0], `"`)
			if col := table.GetColumn(colName); col != nil {
				col.Indexes[index.Name] = index.Type
			} else {
				err := fmt.Errorf("unknown col %s seq %d, in index %v of table %v, columns %v", name, seq, index.Name, table.Name, table.ColumnNames)
				e.logger.Errorf(err, "table %s column %s not found", table.Name, colName)
				return err
			}
		}
	}
	return nil
}
