package defines

import "github.com/seerx/gpa/engine/objs"

type Object struct {
	repo *RepoInterface
	*objs.Object
}

func NewObject(obj *objs.Object, repo *RepoInterface) *Object {
	return &Object{
		repo:   repo,
		Object: obj,
	}
}
