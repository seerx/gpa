//+mro-ignore
// DO NOT EDIT THIS FILE
//+mro-provides:{{- range $i, $v := .Repos -}}
{{$v}},
{{- end -}}
{{.Space}}
// Generated by mro at {{.Time}}
package {{.PackageName}}

import (
	"fmt"
	"github.com/seerx/mro/db"
)

type Repository interface {
	{{ range $i, $v := .Repos -}}
	{{$v}}() {{$v}}
	{{ end }}
}

type Maker func(p *db.Provider) Repository

var makerMap = map[string]interface{}{}
var defaultMaker interface{}

func Register(dialect string, maker interface{}) {
	if len(makerMap) == 0 {
		defaultMaker = maker
	}
	makerMap[dialect] = maker
}

func GetRepository(p *db.Provider, dialect ...string) Repository {
	mk := defaultMaker
	if len(dialect) > 0 {
		var ok bool
		mk, ok = makerMap[dialect[0]]
		if !ok {
			panic(fmt.Sprintf("no provider maker named [%s] registered", dialect[0]))
		}
	}
	if mk == nil {
		panic("no provider maker registered")
	}
	maker := mk.(Maker)
	return maker(p)
}