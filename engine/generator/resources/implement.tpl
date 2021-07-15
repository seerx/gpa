//+mro-ignore
// DO NOT EDIT THIS FILE
// Generated by mro at {{.Time}}
package {{.packageName}}

import (
	"github.com/seerx/mro/db"
	{{.reposPackageName}} "{{.reposPackage}}"
)

type repository struct {
	p *db.Provider
	{{ range $i, $v := .Repos }}
	{{$v.Instance}} *{{$v.Name -}}
	{{ end }}
}

func maker(p *db.Provider) *repository {
	return &repository{p: p}
}

func init() {
	{{.reposPackageName}}.Register("{{.dialect}}", maker)
}

{{ range $i, $v := .Repos }}
func (r *repository) {{$v.Name}}() *{{$v.Name}} {
	if r.{{$v.Instance}} == nil {
		r.{{$v.Instance}} = &{{$v.Name}}{p: r.p}
	}
	return r.{{$v.Instance}}
}
{{ end }}