package dialect

import (
	"fmt"

	"github.com/seerx/gpa/engine/constants"
	"github.com/seerx/gpa/engine/sql/dialect/dialects"
	"github.com/seerx/gpa/engine/sql/dialect/intf"
)

var (
	drivers = map[constants.DIALECT]intf.Driver{}
	// dialectMap = map[constants.DIALECT]intf.Dialect{}
)

func init() {
	dialects.RegisterPostgres(register)
}

func register(dialect constants.DIALECT, drv intf.Driver) {
	// dialectMap[name] = d

	if drv == nil {
		panic("core: Register driver is nil")
	}
	// drv.Dialect = d
	if _, dup := drivers[dialect]; dup {
		panic("core: Register called twice for driver " + dialect)
	}
	drivers[dialect] = drv
}

// func RegisterDialect(name string, dialect Dialect) {
// 	dialectMap[name] = dialect
// }

// func GetDialect(name types.DRIVER) intf.Dialect {
// 	return dialectMap[string(name)]
// }

// func RegisterDriver(driverName string, driver Driver) {
// 	if driver == nil {
// 		panic("core: Register driver is nil")
// 	}
// 	if _, dup := drivers[driverName]; dup {
// 		panic("core: Register called twice for driver " + driverName)
// 	}
// 	drivers[driverName] = driver
// }

// func GetDriver(dialect constants.DIALECT) intf.Driver {
// 	return drivers[dialect]
// }

// func RegisteredDriverCount() int {
// 	return len(drivers)
// }

// OpenDialect opens a dialect via driver name and connection string
func OpenDialect(dialect constants.DIALECT, connstr string) (intf.Dialect, error) {
	driver := drivers[dialect] // GetDriver(driverName)
	if driver == nil {
		return nil, fmt.Errorf("unsupported driver name: %v", dialect)
	}

	uri, err := driver.Parse(dialect, connstr)
	if err != nil {
		return nil, err
	}

	dial := driver.GetDialect() //  GetDialect(uri.DRIVER)
	if dial == nil {
		return nil, fmt.Errorf("unsupported dialect type: %v", uri.DRIVER)
	}

	dial.Init(uri)

	return dial, nil
}
