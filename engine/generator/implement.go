package generator

import (
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/seerx/gpa/engine/generator/defines"
)

type inp struct {
	Instance string
	Name     string
}

// createImplementFile 生成实现文件
func (g *Generator) createImplementFile() error {
	var err error
	dest := g.Info.CreateImplementFilePath("implement.go")
	repos := []*inp{}

	if err := g.Info.TraverseRepos(func(intf *defines.RepoInterface, rf *defines.RepoFile) error {
		repos = append(repos, &inp{
			Name:     intf.Name,
			Instance: strings.ToLower(intf.Name[:1]) + intf.Name[1:],
		})
		return nil
	}); err != nil {
		return err
	}

	// tmpl, err := template.New("implement").Parse(implemntsgo)
	tmpl, err := template.ParseFS(templates, "resources/implement.tpl")
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}

	defer file.Close()

	// buf := bytes.NewBuffer([]byte{})
	return tmpl.Execute(file, map[string]interface{}{
		"dialect":          g.Info.Dialect,
		"reposPackage":     g.Info.Package,
		"reposPackageName": g.Info.PackageName,
		"Time":             time.Now().Format("2006-01-02 15:04:05"),
		"packageName":      g.Info.Dialect,
		"Repos":            repos,
	})
}
