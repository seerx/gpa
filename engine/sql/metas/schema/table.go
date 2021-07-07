package schema

import (
	"reflect"
)

type Table struct {
	// dialect     dialects.Dialect
	// Model       interface{}
	Type          reflect.Type
	Name          string
	Columns       []*Column
	ColumnNames   []string
	columnMap     map[string]*Column
	Indexes       map[string]*Index
	PrimaryKeys   []string
	AutoIncrement string
	// Created       map[string]bool
	// Updated string
	// Deleted     string
}

func NewEmptyTable() *Table {
	return NewTable(nil, "")
}

func NewTable(modelType reflect.Type, tableName string) *Table {
	return &Table{
		Type:      modelType,
		Name:      tableName,
		columnMap: map[string]*Column{},
		Indexes:   map[string]*Index{},
	}
}

// func (t *Table) GetName() string             { return t.Name }
// func (t *Table) GetColumn(index int) *Column { return t.Columns[index] }
// func (t *Table) GetColumnCount() int         { return len(t.Columns) }

func (t *Table) AddIndex(name string, typ int, col *Column) {
	fieldName := col.FieldName()
	if idx, ok := t.Indexes[name]; ok {
		idx.AddColumn(fieldName)
		col.Indexes[idx.Name] = typ
	} else {
		idx := NewIndex(name, typ)
		idx.AddColumn(fieldName)
		t.Indexes[name] = idx
		col.Indexes[idx.Name] = typ
	}
}

func (t *Table) AddColumn(col *Column) {
	fieldName := col.FieldName()
	t.Columns = append(t.Columns, col)
	t.ColumnNames = append(t.ColumnNames, fieldName)
	t.columnMap[fieldName] = col

	// colName := strings.ToLower(col.FieldName)

	if col.IsPrimaryKey {
		t.PrimaryKeys = append(t.PrimaryKeys, fieldName)
	}
	if col.IsAutoIncrement {
		t.AutoIncrement = fieldName
	}
	// if col.Created {
	// 	t.Created[col.FieldName] = true
	// }
	// if col.Updated {
	// 	t.Updated = col.FieldName
	// }
	// if col.IsDeleted {
	// 	table.Deleted = col.Name
	// }
}

func (t *Table) GetColumn(name string) *Column {
	return t.columnMap[name]
}

func (table *Table) GetPKColumns() []*Column {
	columns := make([]*Column, len(table.PrimaryKeys))
	for i, name := range table.PrimaryKeys {
		columns[i] = table.GetColumn(name)
	}
	return columns
}

// func (table *Table) UpdatedColumn() *Column {
// 	return table.GetColumn(table.Updated)
// }

// func (table *Table) DeletedColumn() *Column {
// 	return table.GetColumn(table.Deleted)
// }
