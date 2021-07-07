package dialects

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/seerx/gpa/engine/sql/types"
)

type postgres struct {
	baseDialect
}

func init() {
	RegisterDialect("postgres", &postgres{})
	RegisterDriver("postgres", &pqDriver{})
}

func (p *postgres) Init(uri *URI) error {
	// p.quoter = postgresQuoter
	return p.baseDialect.Init(p, uri)
}

func (p *postgres) DataTypeOf(val reflect.Value) *types.SQLType {
	return types.Type2SQLType(val.Type())
}

// ----------------- 以下为 Driver ----------------------

type values map[string]string

func (vs values) Set(k, v string) {
	vs[k] = v
}

func (vs values) Get(k string) (v string) {
	return vs[k]
}

type pqDriver struct {
}

func parseURL(connstr string) (string, error) {
	u, err := url.Parse(connstr)
	if err != nil {
		return "", err
	}

	if u.Scheme != "postgresql" && u.Scheme != "postgres" {
		return "", fmt.Errorf("invalid connection protocol: %s", u.Scheme)
	}

	escaper := strings.NewReplacer(` `, `\ `, `'`, `\'`, `\`, `\\`)

	if u.Path != "" {
		return escaper.Replace(u.Path[1:]), nil
	}

	return "", nil
}

func parseOpts(name string, o values) error {
	if len(name) == 0 {
		return fmt.Errorf("invalid options: %s", name)
	}

	name = strings.TrimSpace(name)

	ps := strings.Split(name, " ")
	for _, p := range ps {
		kv := strings.Split(p, "=")
		if len(kv) < 2 {
			return fmt.Errorf("invalid option: %q", p)
		}
		o.Set(kv[0], kv[1])
	}

	return nil
}

func (p *pqDriver) Parse(driverName, dataSourceName string) (*URI, error) {
	db := &URI{DBType: types.POSTGRES}
	var err error

	if strings.HasPrefix(dataSourceName, "postgresql://") || strings.HasPrefix(dataSourceName, "postgres://") {
		db.DBName, err = parseURL(dataSourceName)
		if err != nil {
			return nil, err
		}
	} else {
		o := make(values)
		err = parseOpts(dataSourceName, o)
		if err != nil {
			return nil, err
		}

		db.DBName = o.Get("dbname")
	}

	if db.DBName == "" {
		return nil, errors.New("dbname is empty")
	}

	return db, nil
}
