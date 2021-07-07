package dialects

import (
	"fmt"

	"github.com/seerx/gpa/engine/sql/types"
)

var (
	drivers    = map[string]Driver{}
	dialectMap = map[string]Dialect{}
)

func RegisterDialect(name string, dialect Dialect) {
	dialectMap[name] = dialect
}

func GetDialect(name types.DBType) Dialect {
	return dialectMap[string(name)]
}

func RegisterDriver(driverName string, driver Driver) {
	if driver == nil {
		panic("core: Register driver is nil")
	}
	if _, dup := drivers[driverName]; dup {
		panic("core: Register called twice for driver " + driverName)
	}
	drivers[driverName] = driver
}

func GetDriver(driverName string) Driver {
	return drivers[driverName]
}

// func RegisteredDriverCount() int {
// 	return len(drivers)
// }

// OpenDialect opens a dialect via driver name and connection string
func OpenDialect(driverName, connstr string) (Dialect, error) {
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
