package defines

import (
	"errors"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/seerx/gpa/engine/constants"
	"github.com/seerx/gpa/engine/generator/xtype"
	"github.com/seerx/gpa/logger"
	"github.com/seerx/gpa/utils"
)

type Info struct {
	Parsed             bool
	Package            string
	PackageName        string
	OldRepositoriesMap map[string]bool // 从已存在的 interface.go 文件中读取的已经提供的 Repository 列表
	Dir                string
	Dialect            string
	FSet               *token.FileSet `json:"-"`
	Files              []*RepoFile
	logger             logger.GpaLogger
	xtypeParser        *xtype.XTypeParser
}

func NewInfo(pkg string, dialect constants.DIALECT, logger logger.GpaLogger) *Info {
	return &Info{
		FSet:        token.NewFileSet(),
		Package:     pkg,
		Dialect:     string(dialect),
		logger:      logger,
		xtypeParser: xtype.NewXTypeParser(string(dialect), logger),
	}
}

// func (m *Info) ParseXType(obj *Object, dir string) (*xtype.XType, error) {
// 	obj.repo.repoFile.
// 	if obj.Type.IsStruct() {
// 		return
// 	}
// 	m.xtypeParser.Parse()
// }

func (m *Info) IsProvidesChanged() bool {
	for k := range m.OldRepositoriesMap {
		m.OldRepositoriesMap[k] = false
	}
	// changed := false
	if err := m.TraverseRepos(func(intf *RepoInterface, rf *RepoFile) error {
		if _, ok := m.OldRepositoriesMap[intf.Name]; ok {
			m.OldRepositoriesMap[intf.Name] = true
		} else {
			// changed = true
			return errors.New("")
		}
		return nil
	}); err != nil {
		return true
	}
	for _, v := range m.OldRepositoriesMap {
		if !v {
			// 某些 repo 已经删除
			return true
		}
	}
	return false
}

func (m *Info) TraverseRepoFiles(fn func(rf *RepoFile) error) error {
	for _, rf := range m.Files {
		if err := fn(rf); err != nil {
			return err
		}
	}
	return nil
}

func (m *Info) TraverseFuncs(fn func(f *Func, intf *RepoInterface, rf *RepoFile) error) error {
	for _, rf := range m.Files {
		for _, ri := range rf.Repos {
			for _, f := range ri.Funcs {
				if err := fn(f, ri, rf); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (m *Info) TraverseRepos(fn func(intf *RepoInterface, rf *RepoFile) error) error {
	for _, rf := range m.Files {
		for _, ri := range rf.Repos {
			if err := fn(ri, rf); err != nil {
				return err
			}
		}
	}
	return nil
}

// oldname: CheckDialectDir
func (m *Info) CreateRepositoryDirIfNotExists() error {
	dialectDir := filepath.Join(m.Dir, m.Dialect)
	return utils.MakeDirsIfNotExists(dialectDir)
}

// oldname:  CreateImplementFilePath
func (m *Info) GenerateImplementFilePath(fileName string) string {
	return filepath.Join(m.Dir, m.Dialect, fileName)
}

// oldname:  CreateInterfaceFilePath
func (m *Info) GenerateInterfaceFilePath(fileName string) string {
	return filepath.Join(m.Dir, fileName)
}

func (m *Info) FindRepoFiles() (err error) {

	fs, err := ioutil.ReadDir(m.Dir)
	if err != nil {
		return err
	}

	// 查找 mro 定义文件
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), ".go") {
			fileName := f.Name()
			file := filepath.Join(m.Dir, f.Name())
			// fmt.Println(file)
			f, err := parser.ParseFile(m.FSet, file, nil, parser.ParseComments)
			if err != nil {
				return err
			}

			if len(f.Comments) > 0 {
				head := f.Comments[0].Text()
				lines := strings.Split(head, "\n")
				ignore := false
				for _, line := range lines {
					if line == "+gpa-ignore" {
						ignore = true
						continue
					}
					if strings.Index(line, "+gpa-provides:") == 0 {
						// 提供的 repository 列表
						contents := line[len("+gpa-provides:"):]
						ary := strings.Split(contents, ",")
						for _, item := range ary {
							if item == "" {
								continue
							}
							if m.OldRepositoriesMap == nil {
								m.OldRepositoriesMap = map[string]bool{}
							}
							m.OldRepositoriesMap[item] = false
						}
					}
				}
				if ignore {
					// mro 生成的文件，忽略
					continue
				}
			}

			if m.PackageName == "" {
				m.PackageName = f.Name.Name
			}
			repo := NewRepoFile(m, f)
			repo.Path = file
			repo.Name = fileName //  fileName[:len(fileName)-3]
			m.Files = append(m.Files, repo)
		}
	}
	return
}

func (m *Info) CreateImplementFilePath(fileName string) string {
	return filepath.Join(m.Dir, m.Dialect, fileName)
}

func (m *Info) CreateInterfaceFilePath(fileName string) string {
	return filepath.Join(m.Dir, fileName)
}
