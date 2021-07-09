package dialect

import (
	"fmt"

	"github.com/seerx/gpa/engine/sql/dialect/dialects"
	"github.com/seerx/gpa/engine/sql/dialect/intf"
	"github.com/seerx/gpa/engine/sql/types"
)

var (
	drivers    = map[string]intf.Driver{}
	dialectMap = map[string]intf.Dialect{}
)

func init() {
	dialects.RegisterPostgres(register)
}

func register(name string, d intf.Dialect, drv intf.Driver) {
	dialectMap[name] = d

	if drv == nil {
		panic("core: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("core: Register called twice for driver " + name)
	}
	drivers[name] = drv
}

// func RegisterDialect(name string, dialect Dialect) {
// 	dialectMap[name] = dialect
// }

func GetDialect(name types.DBType) intf.Dialect {
	return dialectMap[string(name)]
}

// func RegisterDriver(driverName string, driver Driver) {
// 	if driver == nil {
// 		panic("core: Register driver is nil")
// 	}
// 	if _, dup := drivers[driverName]; dup {
// 		panic("core: Register called twice for driver " + driverName)
// 	}
// 	drivers[driverName] = driver
// }

func GetDriver(driverName string) intf.Driver {
	return drivers[driverName]
}

// func RegisteredDriverCount() int {
// 	return len(drivers)
// }

// OpenDialect opens a dialect via driver name and connection string
func OpenDialect(driverName, connstr string) (intf.Dialect, error) {
	driver := GetDriver(driverName)
	if driver == nil {
		return nil, fmt.Errorf("unsupported driver name: %v", driverName)
	}

	uri, err := driver.Parse(driverName, connstr)
	if err != nil {
		return nil, err
	}

	dialect := GetDialect(uri.DBType)
	if dialect == nil {
		return nil, fmt.Errorf("unsupported dialect type: %v", uri.DBType)
	}

	dialect.Init(uri)

	return dialect, nil
}
