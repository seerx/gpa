//+mro-ignore
// DO NOT EDIT THIS FILE
// Generated by mro at 2021-07-16 09:23:20
package postgres

import (
	"github.com/seerx/gpa/rt"
	repos "github.com/seerx/gpa/examples/pratics/repos"
)

type repository struct {
	p *rt.Provider
	
	userRepo *UserRepo
	update *Update
}

func maker(p *rt.Provider) *repository {
	return &repository{p: p}
}

func init() {
	repos.Register("postgres", maker)
}


func (r *repository) UserRepo() *UserRepo {
	if r.userRepo == nil {
		r.userRepo = &UserRepo{p: r.p}
	}
	return r.userRepo
}

func (r *repository) Update() *Update {
	if r.update == nil {
		r.update = &Update{p: r.p}
	}
	return r.update
}
