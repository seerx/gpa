package main

import (
	_ "github.com/lib/pq"
	"github.com/seerx/gpa/engine"
	"github.com/seerx/gpa/engine/generator/parse"
	"github.com/seerx/gpa/examples/pratics/models"
	"github.com/seerx/gpa/logger"
	"github.com/seerx/logo/log"
)

func main() {
	log.Info("starting ...")
	// log.SetLogFileLine(true)
	log.SetLogErrorCallStacks(true)
	eg, err := engine.New("postgres", "host=192.168.0.10 port=5432 user=checkin dbname=checkin password=hcdj&*@HDSBddns776^&^&DW sslmode=disable connect_timeout=10")
	if err != nil {
		log.WithError(err).Error("connect error")
		return
	}

	if err := eg.Sync(&models.User{}, &models.Student{}); err != nil {
		log.WithError(err).Error("sync tables error")
	}

	doParse()
	log.Info("exiting ...")
}

func doParse() {
	parse.Parse("github.com/seerx/gpa/examples/pratics/repos", "postgres", logger.GetLogger())
}
