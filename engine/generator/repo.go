package generator

import (
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/seerx/gpa/engine/generator/defines"
	rdesc "github.com/seerx/gpa/engine/generator/repo-desc"
	"github.com/seerx/gpa/engine/sql/names"
)

type FuncInfo struct {
	rdesc.FuncDesc
	Repo     *repo
	Name     string
	Template string
}

type repo struct {
	Instance       string
	Name           string
	RunTimePackage string
	Funcs          []*FuncInfo
}

func (g *Generator) createRepoFile(rf *defines.RepoFile) error {
	dest := g.Info.CreateImplementFilePath(rf.Name)
	repos := []*repo{}
	for _, r := range rf.Repos {
		rp := &repo{
			Name:           r.Name,
			RunTimePackage: r.RunTimePackage,
			Instance:       names.LowerFirstChar(r.Name), //  strings.ToLower(r.Name[:1]) + r.Name[1:],
		}
		repos = append(repos, rp)
		for _, fn := range r.Funcs {
			m := g.getMethod(fn)
			// g := gen.GetGenerator(fn)
			if m != nil {
				desc, err := m.Parse()
				if err != nil {
					g.logger.Error(err)
					continue
				}

				fnInfo := &FuncInfo{
					FuncDesc: *desc,
					Name:     fn.Name,
					Repo:     rp,
					Template: string(fn.Template),
				}
				rp.Funcs = append(rp.Funcs, fnInfo)
				// f.SQL = sql
				// fmt.Println("\t", f.Name, sql)
			} else {
				g.logger.Warnf("func name [%s] is not support", fn.Name)
			}

		}
	}

	funcs := template.FuncMap{"join": strings.Join}
	tmpl, err := template.ParseFS(templates, "resources/repo.tpl")
	if err != nil {
		return err
	}
	tmpl, err = tmpl.Funcs(funcs).ParseFS(templates,
		"resources/repo-body.tpl",
		"resources/func-insert.tpl",
		"resources/func-update.tpl",
		"resources/func-delete.tpl",
		"resources/func-find.tpl",
		"resources/func-find-block-read-row.tpl",
		"resources/func-find-block-read-rows.tpl",
		"resources/func-find-block-read-rows-callback.tpl",
		"resources/func-count.tpl")
	if err != nil {
		return err
	}

	// tmpl.ParseFS()  .Funcs(funcs)

	file, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	// buf := bytes.NewBuffer([]byte{})
	return tmpl.Execute(file, map[string]interface{}{
		"Space":   "",
		"package": g.Info.Dialect,
		"imports": rf.Imports,
		"Repos":   repos,
		"Time":    time.Now().Format("2006-01-02 15:04:05"),
	})
}

/**   $. ???????????? .??????
// ????????? $.Instance ???????????? pipeline.Instance
// .Name ???????????? .Funcs ??????????????? .Name
{{ range .Funcs }}
func ({{$.Instance}} *{{$.Name}}) {{.Name}}(user *models.User, name string) (*models.User, error) {
	user.Name = name

	_, err := {{$.Instance}}.query.Exec("insert into user ")
	if err != nil {
		return nil, err
	}
	return user, nil
}
{{ end }}
**/
