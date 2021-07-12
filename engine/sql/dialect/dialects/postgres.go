package dialects

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/seerx/gpa/engine/sql/dialect/intf"
	"github.com/seerx/gpa/engine/sql/metas/schema"
	"github.com/seerx/gpa/engine/sql/types"
	"github.com/seerx/gpa/rt/exec"
)

type postgres struct {
	baseDialect
}

var (
	// DefaultPostgresSchema default postgres schema
	DefaultPostgresSchema = "public"
)

var postgresQuoter = intf.Quoter{
	Prefix:     '"',
	Suffix:     '"',
	IsReserved: intf.AlwaysReserve,
}

func RegisterPostgres(r intf.Regitser) {
	r("postgres", &postgres{}, &pqDriver{})
}

func (p *postgres) Init(uri *intf.URI) error {
	p.quoter = postgresQuoter
	return p.baseDialect.Init(p, uri)
}

func (p *postgres) getDatabaseSchema() string {
	if p.uri.Schema != "" {
		return p.uri.Schema
	}
	return DefaultPostgresSchema
}

func (p *postgres) TableNameWithSchema(tableName string) string {
	if p.getDatabaseSchema() != "" && !strings.Contains(tableName, ".") {
		return fmt.Sprintf("%s.%s", p.getDatabaseSchema(), tableName)
	}
	return tableName
}

// SQLType 转换为 SQL 数据类型
func (p *postgres) SQLType(col *schema.Column) string {
	var res string
	// var tag = col
	switch t := col.Type.Name; t {
	case types.TinyInt:
		res = types.SmallInt
		return res
	case types.Bit:
		res = types.Boolean
		return res
	case types.MediumInt, types.Int, types.Integer:
		if col.IsAutoIncrement {
			return types.Serial
		}
		return types.Integer
	case types.BigInt:
		if col.IsAutoIncrement {
			return types.BigSerial
		}
		return types.BigInt
	case types.Serial, types.BigSerial:
		col.IsAutoIncrement = true
		col.Nullable = false
		res = t
	case types.Binary, types.VarBinary:
		return types.Bytea
	case types.DateTime:
		res = types.TimeStamp
	case types.TimeStampz:
		return "timestamp with time zone"
	case types.Float:
		res = types.Real
	case types.TinyText, types.MediumText, types.LongText:
		res = types.Text
	case types.NChar:
		res = types.Char
	case types.NVarchar:
		res = types.Varchar
	case types.Uuid:
		return types.Uuid
	case types.Blob, types.TinyBlob, types.MediumBlob, types.LongBlob:
		return types.Bytea
	case types.Double:
		return "DOUBLE PRECISION"
	default:
		if col.IsAutoIncrement {
			return types.Serial
		}
		res = t
	}

	if strings.EqualFold(res, "bool") {
		// for bool, we don't need length information
		return res
	}
	hasLen1 := (col.Length > 0)
	hasLen2 := (col.Length2 > 0)

	if hasLen2 {
		res += "(" + strconv.Itoa(col.Length) + "," + strconv.Itoa(col.Length2) + ")"
	} else if hasLen1 {
		res += "(" + strconv.Itoa(col.Length) + ")"
	}
	return res
}

func (p *postgres) Quoter() intf.Quoter {
	return p.quoter
}

// AutoIncrStr 自增字段标志字符串
func (p *postgres) AutoIncrStr() string {
	return ""
}

func (p *postgres) GetTables(ex exec.SQLExecutor, ctx context.Context) ([]*schema.Table, error) {
	args := []interface{}{}
	s := "SELECT tablename FROM pg_tables"
	dbSchema := p.getDatabaseSchema()
	if dbSchema != "" {
		args = append(args, dbSchema)
		s = s + " WHERE schemaname = $1"
	}

	rows, err := ex.QueryContextRows(ctx, s, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []*schema.Table
	for rows.Next() {
		table := schema.NewEmptyTable()
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		table.Name = name
		tables = append(tables, table)
	}
	return tables, nil
}

// SQLTableExists 生成判断表是否存在的 SQL
func (p *postgres) SQLTableExists(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "select tablename from pg_tables where tablename=$1", args
}

// SQLCreateTable 生成创建表结构的 SQL
func (p *postgres) SQLCreateTable(table *schema.Table, tableName string) ([]string, error) {
	var sql string
	sql = "CREATE TABLE IF NOT EXISTS "

	if tableName == "" {
		tableName = table.Name
	}
	quoter := p.Quoter()
	sql += quoter.Quote(tableName)
	sql += " ("
	if len(table.Columns) > 0 {
		// 查找全部主键
		var pkList []string
		for _, col := range table.Columns {
			// col := schema.GetColumn(n)
			if col.IsPrimaryKey {
				pkList = append(pkList, col.FieldName())
			}
		}
		// inlinePK := len(pkList) == 1
		var cols = []string{}
		for _, col := range table.Columns {
			// col := schema.GetColumn(n)
			s, err := p.SQLColumn(col, len(pkList) == 1)
			if err != nil {
				return nil, err
			}
			cols = append(cols, s)
		}
		sql += strings.Join(cols, ",")
		if len(pkList) > 1 {
			sql += "PRIMARY KEY ( "
			sql += quoter.Join(pkList, ",")
			sql += " ) "
		}
		sql = strings.TrimSpace(sql)
		// sql = sql[:len(sql)-1]
	}
	sql += ")"

	return []string{sql}, nil
}
