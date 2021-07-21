package generator

import (
	"go/build"
	"go/token"

	"github.com/seerx/gpa/engine/constants"
	"github.com/seerx/gpa/engine/generator/defines"
	"github.com/seerx/gpa/engine/generator/method"
	"github.com/seerx/gpa/engine/generator/sqlgenerator"
	"github.com/seerx/gpa/logger"
	"github.com/seerx/logo/log"
)

type Generator struct {
	Info    *defines.Info
	methods []method.Method
	dialect constants.DIALECT
	logger  logger.GpaLogger
}

func New(dialectName constants.DIALECT, logger logger.GpaLogger) (*Generator, error) {
	sqlg, err := sqlgenerator.GetGenerator(dialectName)
	if err != nil {
		return nil, err
	}
	var g = Generator{
		dialect: dialectName,
		logger:  logger,
	}
	g.methods = method.CreateMethods(sqlg, logger)
	return &g, nil
}

func (g *Generator) getMethod(fn *defines.Func) method.Method {
	for _, g := range g.methods {
		if g.Test(fn) {
			return g
		}
	}
	return nil
}

func (g *Generator) Generate(fullPackagePath string, forceRegenerate bool) error {
	if err := g.parse(fullPackagePath); err != nil {
		return err
	}
	// 检查源码路径是否存在，不存在则创建
	if err := g.Info.CreateRepositoryDirIfNotExists(); err != nil {
		return err
	}

	// 依据接口定义创建 repository
	if err := g.Info.TraverseRepoFiles(func(rf *defines.RepoFile) error {
		if err := g.createRepoFile(rf); err != nil {
			panic(err)
		}
		return nil
	}); err != nil {
		return err
	}

	if forceRegenerate || g.Info.IsProvidesChanged() {
		// 重新生成 interface.go
		if err := g.createInterfaceFile(); err != nil {
			return err
		}
		// 生成实现 implement.go
		if err := g.createImplementFile(); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) parse(pkg string) error {
	// 解析路径信息
	pkgInfo, err := build.Import(pkg, "", build.FindOnly)
	if err != nil {
		return err
	}
	info := defines.NewInfo(pkg, g.dialect, g.logger)
	info.FSet = token.NewFileSet()
	info.Dir = pkgInfo.Dir
	g.logger.Infof("parsing %s with dialect %s ...", info.Dir, string(g.dialect))

	// 查找全部用户定义的 repostory 接口文件
	if err := info.FindRepoFiles(); err != nil {
		g.logger.Error(err, "search repo define files error")
		return err
	}
	g.logger.Infof("found %d repo files", len(info.Files))
	g.logger.Infof("got repos package: %s, named: %s", info.Package, info.PackageName)
	for _, r := range info.Files {
		log.Infof("parsing repo %s", r.Path)
		if err := r.Parse(info.Dialect); err != nil {
			g.logger.Errorf(err, "parse repo file %s error", r.Path)
			return err
		}

		_ = r.AddRuntimePackage()
	}
	g.Info = info
	return nil
}
