package intf

import (
	"strings"
	"time"

	"github.com/seerx/gpa/engine/sql/types"
)

type Driver interface {
	Parse(string, string) (*URI, error)
}

// URI represents an uri to visit database
type URI struct {
	DBType  types.DBType
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
	if uri.DBType == types.POSTGRES {
		uri.Schema = strings.TrimSpace(schema)
	}
}

// SupportColumnVarchar2Text 是否支持把现有列的数据类型从 varchar 扩展为 text
func (uri *URI) SupportColumnVarchar2Text() bool {
	return uri.DBType == types.MYSQL || uri.DBType == types.POSTGRES
}

// SupportColumnVarcharIncLength 是否支持把现有列的数据类型从 varchar 增加长度
func (uri *URI) SupportColumnVarcharIncLength() bool {
	return uri.DBType == types.MYSQL || uri.DBType == types.POSTGRES
}
