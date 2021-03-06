package types

const (
	Bit       = "BIT"
	TinyInt   = "TINYINT"
	SmallInt  = "SMALLINT"
	MediumInt = "MEDIUMINT"
	Int       = "INT"
	Integer   = "INTEGER"
	BigInt    = "BIGINT"

	Enum = "ENUM"
	Set  = "SET"

	Char             = "CHAR"
	Varchar          = "VARCHAR"
	NChar            = "NCHAR"
	NVarchar         = "NVARCHAR"
	TinyText         = "TINYTEXT"
	Text             = "TEXT"
	NText            = "NTEXT"
	Clob             = "CLOB"
	MediumText       = "MEDIUMTEXT"
	LongText         = "LONGTEXT"
	Uuid             = "UUID"
	UniqueIdentifier = "UNIQUEIDENTIFIER"
	SysName          = "SYSNAME"

	Date          = "DATE"
	DateTime      = "DATETIME"
	SmallDateTime = "SMALLDATETIME"
	Time          = "TIME"
	TimeStamp     = "TIMESTAMP"
	TimeStampz    = "TIMESTAMPZ"
	Year          = "YEAR"

	Decimal    = "DECIMAL"
	Numeric    = "NUMERIC"
	Money      = "MONEY"
	SmallMoney = "SMALLMONEY"

	Real   = "REAL"
	Float  = "FLOAT"
	Double = "DOUBLE"

	Binary     = "BINARY"
	VarBinary  = "VARBINARY"
	TinyBlob   = "TINYBLOB"
	Blob       = "BLOB"
	MediumBlob = "MEDIUMBLOB"
	LongBlob   = "LONGBLOB"
	Bytea      = "BYTEA"

	Bool    = "BOOL"
	Boolean = "BOOLEAN"

	Serial    = "SERIAL"
	BigSerial = "BIGSERIAL"

	Json  = "JSON"
	Jsonb = "JSONB"

	XML   = "XML"
	Array = "ARRAY"
)

var SqlTypes = map[string]int{
	Bit:       NUMERIC_TYPE,
	TinyInt:   NUMERIC_TYPE,
	SmallInt:  NUMERIC_TYPE,
	MediumInt: NUMERIC_TYPE,
	Int:       NUMERIC_TYPE,
	Integer:   NUMERIC_TYPE,
	BigInt:    NUMERIC_TYPE,

	Enum:  TEXT_TYPE,
	Set:   TEXT_TYPE,
	Json:  TEXT_TYPE,
	Jsonb: TEXT_TYPE,

	XML: TEXT_TYPE,

	Char:       TEXT_TYPE,
	NChar:      TEXT_TYPE,
	Varchar:    TEXT_TYPE,
	NVarchar:   TEXT_TYPE,
	TinyText:   TEXT_TYPE,
	Text:       TEXT_TYPE,
	NText:      TEXT_TYPE,
	MediumText: TEXT_TYPE,
	LongText:   TEXT_TYPE,
	Uuid:       TEXT_TYPE,
	Clob:       TEXT_TYPE,
	SysName:    TEXT_TYPE,

	Date:          TIME_TYPE,
	DateTime:      TIME_TYPE,
	Time:          TIME_TYPE,
	TimeStamp:     TIME_TYPE,
	TimeStampz:    TIME_TYPE,
	SmallDateTime: TIME_TYPE,
	Year:          TIME_TYPE,

	Decimal:    NUMERIC_TYPE,
	Numeric:    NUMERIC_TYPE,
	Real:       NUMERIC_TYPE,
	Float:      NUMERIC_TYPE,
	Double:     NUMERIC_TYPE,
	Money:      NUMERIC_TYPE,
	SmallMoney: NUMERIC_TYPE,

	Binary:    BLOB_TYPE,
	VarBinary: BLOB_TYPE,

	TinyBlob:         BLOB_TYPE,
	Blob:             BLOB_TYPE,
	MediumBlob:       BLOB_TYPE,
	LongBlob:         BLOB_TYPE,
	Bytea:            BLOB_TYPE,
	UniqueIdentifier: BLOB_TYPE,

	Bool: NUMERIC_TYPE,

	Serial:    NUMERIC_TYPE,
	BigSerial: NUMERIC_TYPE,

	Array: ARRAY_TYPE,
}
