package parse

// func ParseRepos(pkg, dialect string, logger logger.GpaLogger) (*defines.Info, error) {
// 	// 解析路径信息
// 	pkgInfo, err := build.Import(pkg, "", build.FindOnly)
// 	if err != nil {
// 		return nil, err
// 	}
// 	info := defines.NewInfo(pkg, dialect, logger)
// 	info.FSet = token.NewFileSet()
// 	info.Dir = pkgInfo.Dir
// 	// info := &defines.Info{
// 	// 	FSet:    token.NewFileSet(),
// 	// 	Package: pkg,
// 	// 	Dialect: dialect,
// 	// 	Dir:     pkgInfo.Dir,
// 	// }
// 	logger.Infof("parsing %s with dialect %s ...", info.Dir, dialect)

// 	// 查找全部用户定义的 repostory 接口文件
// 	if err := info.FindRepoFiles(); err != nil {
// 		logger.Error(err, "search repo define files error")
// 		return nil, err
// 	}
// 	logger.Infof("found %d repo files", len(info.Files))
// 	logger.Infof("got repos package: %s, named: %s", info.Package, info.PackageName)
// 	for _, r := range info.Files {
// 		log.Infof("parsing repo %s", r.Path)
// 		if err := r.Parse(info.Dialect); err != nil {
// 			logger.Errorf(err, "parse repo file %s error", r.Path)
// 			return nil, err
// 		}

// 		_ = r.AddRuntimePackage()
// 	}

// 	return info, nil
// }
