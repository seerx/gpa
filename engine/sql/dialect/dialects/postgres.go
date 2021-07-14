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

func (p *postgres) GetColumns(ex exec.SQLExecutor, ctx context.Context, tableName string) ([]string, map[string]*schema.Column, error) {
	args := []interface{}{tableName}
	s := `SELECT column_name, column_default, is_nullable, data_type, character_maximum_length,
    CASE WHEN p.contype = 'p' THEN true ELSE false END AS primarykey,
    CASE WHEN p.contype = 'u' THEN true ELSE false END AS uniquekey
FROM pg_attribute f
    JOIN pg_class c ON c.oid = f.attrelid JOIN pg_type t ON t.oid = f.atttypid
    LEFT JOIN pg_attrdef d ON d.adrelid = c.oid AND d.adnum = f.attnum
    LEFT JOIN pg_namespace n ON n.oid = c.relnamespace
    LEFT JOIN pg_constraint p ON p.conrelid = c.oid AND f.attnum = ANY (p.conkey)
    LEFT JOIN pg_class AS g ON p.confrelid = g.oid
    LEFT JOIN INFORMATION_SCHEMA.COLUMNS s ON s.column_name=f.attname AND c.relname=s.table_name
WHERE n.nspname= s.table_schema AND c.relkind = 'r'::char AND c.relname = $1%s AND f.attnum > 0 ORDER BY f.attnum;`

	sName := p.getDatabaseSchema()
	if sName != "" {
		s = fmt.Sprintf(s, " AND s.table_schema = $2")
		args = append(args, sName)
	} else {
		s = fmt.Sprintf(s, "")
	}

	rows, err := ex.QueryContextRows(ctx, s, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	cols := make(map[string]*schema.Column)
	colSeq := make([]string, 0)

	for rows.Next() {
		col := new(schema.Column)
		col.Indexes = make(map[string]int)

		var colName, isNullable, dataType string
		var maxLenStr, colDefault *string
		var isPK, isUnique bool
		err = rows.Scan(&colName, &colDefault, &isNullable, &dataType, &maxLenStr, &isPK, &isUnique)
		if err != nil {
			return nil, nil, err
		}

		var maxLen int
		if maxLenStr != nil {
			maxLen, err = strconv.Atoi(*maxLenStr)
			if err != nil {
				return nil, nil, err
			}
		}

		col.Field.Name = strings.Trim(colName, `" `)

		if colDefault != nil {
			var theDefault = *colDefault
			// cockroach has type with the default value with :::
			// and postgres with ::, we should remove them before store them
			idx := strings.Index(theDefault, ":::")
			if idx == -1 {
				idx = strings.Index(theDefault, "::")
			}
			if idx > -1 {
				theDefault = theDefault[:idx]
			}

			if strings.HasSuffix(theDefault, "+00:00'") {
				theDefault = theDefault[:len(theDefault)-7] + "'"
			}

			col.Default = theDefault
			// col.DefaultIsEmpty = false
			if strings.HasPrefix(col.Default, "nextval(") {
				col.IsAutoIncrement = true
				col.Default = ""
				// col.DefaultIsEmpty = true
			}
		} // else {
		// col.DefaultIsEmpty = true
		//}

		if isPK {
			col.IsPrimaryKey = true
		}

		col.Nullable = (isNullable == "YES")

		switch strings.ToLower(dataType) {
		case "character varying", "string":
			col.Type = &types.SQLType{Name: types.Varchar, Length: 0, Length2: 0}
		case "character":
			col.Type = &types.SQLType{Name: types.Char, Length: 0, Length2: 0}
		case "timestamp without time zone":
			col.Type = &types.SQLType{Name: types.DateTime, Length: 0, Length2: 0}
		case "timestamp with time zone":
			col.Type = &types.SQLType{Name: types.TimeStampz, Length: 0, Length2: 0}
		case "double precision":
			col.Type = &types.SQLType{Name: types.Double, Length: 0, Length2: 0}
		case "boolean":
			col.Type = &types.SQLType{Name: types.Bool, Length: 0, Length2: 0}
		case "time without time zone":
			col.Type = &types.SQLType{Name: types.Time, Length: 0, Length2: 0}
		case "bytes":
			col.Type = &types.SQLType{Name: types.Binary, Length: 0, Length2: 0}
		case "oid":
			col.Type = &types.SQLType{Name: types.BigInt, Length: 0, Length2: 0}
		case "array":
			col.Type = &types.SQLType{Name: types.Array, Length: 0, Length2: 0}
		default:
			startIdx := strings.Index(strings.ToLower(dataType), "string(")
			if startIdx != -1 && strings.HasSuffix(dataType, ")") {
				length := dataType[startIdx+8 : len(dataType)-1]
				l, _ := strconv.Atoi(length)
				col.Type = &types.SQLType{Name: "STRING", Length: l, Length2: 0}
			} else {
				col.Type = &types.SQLType{Name: strings.ToUpper(dataType), Length: 0, Length2: 0}
			}
		}
		if _, ok := types.SqlTypes[col.Type.Name]; !ok {
			return nil, nil, fmt.Errorf("unknown colType: %s - %s", dataType, col.Type.Name)
		}

		col.Length = maxLen

		if col.Default != "" {
			if col.Type.IsText() {
				if strings.HasSuffix(col.Default, "::character varying") {
					// col.Default = strings.TrimRight(col.Default, "::character varying")
					col.Default = strings.TrimSuffix(col.Default, "::character varying")
				} else if !strings.HasPrefix(col.Default, "'") {
					col.Default = "'" + col.Default + "'"
				}
			} else if col.Type.IsTime() {
				if strings.HasSuffix(col.Default, "::timestamp without time zone") {
					// col.Default = strings.TrimRight(col.Default, "::timestamp without time zone")
					col.Default = strings.TrimSuffix(col.Default, "::timestamp without time zone")
				}
			}
		}
		cols[col.FieldName()] = col
		colSeq = append(colSeq, col.FieldName())
	}

	return colSeq, cols, nil
}

func (p *postgres) GetIndexes(ex exec.SQLExecutor, ctx context.Context, tableName string) (map[string]*schema.Index, error) {
	args := []interface{}{tableName}
	s := "SELECT indexname, indexdef FROM pg_indexes WHERE tablename=$1"
	if len(p.getDatabaseSchema()) != 0 {
		args = append(args, p.getDatabaseSchema())
		s = s + " AND schemaname=$2"
	}

	rows, err := ex.QueryContextRows(ctx, s, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indexes := make(map[string]*schema.Index)
	for rows.Next() {
		var indexType int
		var indexName, indexdef string
		var colNames []string
		err = rows.Scan(&indexName, &indexdef)
		if err != nil {
			return nil, err
		}

		if indexName == "primary" {
			continue
		}
		indexName = strings.Trim(indexName, `" `)
		if strings.HasSuffix(indexName, "_pkey") {
			continue
		}
		if strings.HasPrefix(indexdef, "CREATE UNIQUE INDEX") {
			indexType = types.UniqueType
		} else {
			indexType = types.IndexType
		}
		colNames = getIndexColName(indexdef)
		var isRegular bool
		if strings.HasPrefix(indexName, "IDX_"+tableName) || strings.HasPrefix(indexName, "UQE_"+tableName) {
			newIdxName := indexName[5+len(tableName):]
			isRegular = true
			if newIdxName != "" {
				indexName = newIdxName
			}
		}

		index := &schema.Index{Name: indexName, Type: indexType, Cols: make([]string, 0)}
		for _, colName := range colNames {
			index.Cols = append(index.Cols, strings.TrimSpace(strings.Replace(colName, `"`, "", -1)))
		}
		index.Regular = isRegular
		indexes[index.Name] = index
	}
	return indexes, nil
}

func getIndexColName(indexdef string) []string {
	var colNames []string

	cs := strings.Split(indexdef, "(")
	for _, v := range strings.Split(strings.Split(cs[1], ")")[0], ",") {
		colNames = append(colNames, strings.Split(strings.TrimLeft(v, " "), " ")[0])
	}

	return colNames
}
