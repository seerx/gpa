package constants

type DRIVER string

const (
	DB_POSTGRES DRIVER = "postgres"
	DB_SQLITE   DRIVER = "sqlite3"
	DB_MYSQL    DRIVER = "mysql"
	DB_MSSQL    DRIVER = "mssql"
	DB_ORACLE   DRIVER = "oracle"
	DB_UNKNOWN  DRIVER = "unknown"
)
