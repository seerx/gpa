package constants

type DIALECT string

const (
	POSTGRES DIALECT = "postgres"
	SQLITE   DIALECT = "sqlite3"
	MYSQL    DIALECT = "mysql"
	MSSQL    DIALECT = "mssql"
	ORACLE   DIALECT = "oracle"
)

var dialdbmap = map[DIALECT]DRIVER{
	POSTGRES: DB_POSTGRES,
	SQLITE:   DB_SQLITE,
	MYSQL:    DB_MYSQL,
	MSSQL:    DB_MSSQL,
	ORACLE:   DB_ORACLE,
}

func (d DIALECT) GetDRIVER() DRIVER {
	dbt, ok := dialdbmap[d]
	if ok {
		return dbt
	}
	return DB_UNKNOWN
}
