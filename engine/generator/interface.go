package generator

import (
	"os"
	"text/template"
	"time"

	"github.com/seerx/gpa/engine/generator/defines"
)

// createInterfaceFile 生成接口文件
func (g *Generator) createInterfaceFile() error {
	dest := g.Info.CreateInterfaceFilePath("interface.go")
	repos := []string{}

	if err := g.Info.TraverseRepos(func(intf *defines.RepoInterface, rf *defines.RepoFile) error {
		repos = append(repos, intf.Name)
		return nil
	}); err != nil {
		return err
	}

	tmpl, err := template.ParseFS(templates, "resources/interface.tpl") // .Parse(interfacego)
	// tmpl, err := template.New("interface").Parse(interfacego)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, map[string]interface{}{
		"Space":       "",
		"Time":        time.Now().Format("2006-01-02 15:04:05"),
		"PackageName": g.Info.PackageName,
		"Repos":       repos,
	})
}
