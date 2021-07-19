package sqlgenerator

import (
	"fmt"

	"github.com/seerx/gpa/engine/constants"
)

var generators = map[constants.DIALECT]SQLGenerator{}

func register(dialect constants.DIALECT, gen SQLGenerator) {
	generators[dialect] = gen
}

func GetGenerator(dialect constants.DIALECT) (SQLGenerator, error) {
	gen := generators[dialect]
	if gen == nil {
		return nil, fmt.Errorf("unsupported dialect type: %v", dialect)
	}
	return gen, nil
}
