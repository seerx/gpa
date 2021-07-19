package intf

import (
	"strings"
	"time"

	"github.com/seerx/gpa/engine/constants"
)

// type DriverParser interface {
// 	Parse(string, string) (*URI, error)
// 	// GetDialect() Dialect
// }

type Driver interface {
	Parse(constants.DIALECT, string) (*URI, error)
	GetDialect() Dialect
}

// func (d *Driver) GetDialect() Dialect {
// 	return d.Dialect
// }

// URI represents an uri to visit database
type URI struct {
	DRIVER  constants.DRIVER
	Proto   string
	Host    string
	Port    string
	DBName  string
	User    string
	Passwd  string
	Charset string
	Laddr   string
	Raddr   string
	Timeout time.Duration
	Schema  string
}

// SetSchema set schema
func (uri *URI) SetSchema(schema string) {
	// hack me
	if uri.DRIVER == constants.DB_POSTGRES {
		uri.Schema = strings.TrimSpace(schema)
	}
}

// SupportColumnVarchar2Text 是否支持把现有列的数据类型从 varchar 扩展为 text
func (uri *URI) SupportColumnVarchar2Text() bool {
	return uri.DRIVER == constants.DB_MYSQL || uri.DRIVER == constants.DB_POSTGRES
}

// SupportColumnVarcharIncLength 是否支持把现有列的数据类型从 varchar 增加长度
func (uri *URI) SupportColumnVarcharIncLength() bool {
	return uri.DRIVER == constants.DB_MYSQL || uri.DRIVER == constants.DB_POSTGRES
}
