//+gpa-ignore
// DO NOT EDIT THIS FILE
//+gpa-provides:Count,Delete,UserRepo,Update,
// Generated by gpa at 2021-07-16 16:43:02
package repos

import (
	"fmt"
	"github.com/seerx/gpa/rt"
)

type Repository interface {
	Count() Count
	Delete() Delete
	UserRepo() UserRepo
	Update() Update
	
}

type Maker func(p *rt.Provider) Repository

var makerMap = map[string]interface{}{}
var defaultMaker interface{}

func Register(dialect string, maker interface{}) {
	if len(makerMap) == 0 {
		defaultMaker = maker
	}
	makerMap[dialect] = maker
}

func GetRepository(p *rt.Provider, dialect ...string) Repository {
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
