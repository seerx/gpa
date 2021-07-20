package main

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/seerx/gpa/engine"
	"github.com/seerx/gpa/engine/generator"
	"github.com/seerx/gpa/engine/generator/defines"
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
	info, err := parse.ParseRepos(engine.TagName, "github.com/seerx/gpa/examples/pratics/repos", "postgres", logger.GetLogger())
	if err != nil {
		panic(err)
	}
	// 检查源码路径是否存在
	if err := info.CreateRepositoryDirIfNotExists(); err != nil {
		panic(err)
	}

	info.TraverseRepoFiles(func(rf *defines.RepoFile) error {
		if err := generator.CreateRepoFile(info, rf, logger.GetLogger()); err != nil {
			panic(err)
		}
		return nil
	})

	if info.IsProvidesChanged() {
		// 重新生成 interface.go
		if err := generator.CreateInterfaceFile(info); err != nil {
			panic(err)
		}
		// 生成实现 implement.go
		if err := generator.CreateImplementFile(info); err != nil {
			panic(err)
		}
	}

	for _, p := range info.Files[0].Repos[0].Funcs[0].Params {
		fmt.Println(p.Name)
	}

	fmt.Println(info.Dir)

	// x := xtype.NewXTypeParser("gpa", "postgres", logger.GetLogger())
	// xt, err := x.Parse("User", info.Dir+"/../models")
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(xt.TableName)
}
