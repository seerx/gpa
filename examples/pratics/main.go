package main

import (
	"fmt"

	"github.com/seerx/gpa/engine/sql/dialects"
	"github.com/seerx/gpa/engine/sql/metas/rflt"
	"github.com/seerx/gpa/examples/pratics/models"
)

func main() {
	dial, err := dialects.OpenDialect("postgres", "host=192.168.0.10 port=5432 user=checkin dbname=checkin password=hcdj&*@HDSBddns776^&^&DW sslmode=disable connect_timeout=10")
	if err != nil {
		panic(err)
	}

	table, err := rflt.Parse(&models.User{}, rflt.NewPropsParser("gpa", dial))
	if err != nil {
		panic(err)
	}

	fmt.Println(table.Name)
}
