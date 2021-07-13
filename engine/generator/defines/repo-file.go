package defines

import (
	"fmt"
	"go/ast"
	"math/rand"
)

type RepoInterface struct {
	repoFile       *RepoFile
	Name           string
	Funcs          []*Func
	RunTimePackage string
}

func NewRepoInterface(name string, repoFile *RepoFile) *RepoInterface {
	return &RepoInterface{
		Name:     name,
		repoFile: repoFile,
	}
}

type RepoFile struct {
	mro  *Info     `json:"-"`
	File *ast.File `json:"-"`
	// Parsed  bool
	Name    string
	Path    string
	Package string
	Imports map[string]string
	Repos   []*RepoInterface

	SQLPackage     string // database/sql
	RunTimePackage string // github.com/seerx/gpa/rt
	DBUtilPackage  string // github.com/seerx/gpa/rt/dbutil
}

func NewRepoFile(mro *Info, file *ast.File) *RepoFile {
	return &RepoFile{
		mro:     mro,
		File:    file,
		Path:    file.Name.Name,
		Imports: map[string]string{},
	}
}

func (rf *RepoFile) AddSQLPackage() string {
	if rf.SQLPackage == "" {
		rf.SQLPackage = rf.addPackage("sql", "database/sql")
	}
	return rf.SQLPackage
}

// func (rf *RepoFile) AddContextPackage() string {
// 	if rf.ContextPackage == "" {
// 		rf.ContextPackage = rf.addPackage("context", "context")
// 	}
// 	return rf.ContextPackage
// }

func (rf *RepoFile) AddRuntimePackage() string {
	if rf.RunTimePackage == "" {
		rf.RunTimePackage = rf.addPackage("rt", "github.com/seerx/gpa/rt")
	}
	for _, intf := range rf.Repos {
		intf.RunTimePackage = rf.RunTimePackage
	}
	return rf.RunTimePackage
}

func (rf *RepoFile) AddDBUtilPackage() string {
	if rf.DBUtilPackage == "" {
		rf.DBUtilPackage = rf.addPackage("dbutil", "github.com/seerx/gpa/rt/dbutil")
	}
	return rf.DBUtilPackage
}

func (rf *RepoFile) addPackage(pkgNamePrefix, pkg string) string {
	var pkgName string
	for {
		pkgName = fmt.Sprintf("%s%d", pkgNamePrefix, rand.Intn(1000))
		if _, ok := rf.Imports[pkgName]; ok {
			continue
		}
		rf.Imports[pkgName] = pkg
		break
	}
	return pkgName
}

func (rf *RepoFile) FindImport(pkg string) string {
	return rf.Imports[pkg]
}
