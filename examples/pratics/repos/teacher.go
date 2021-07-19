package repos

import "github.com/seerx/gpa/examples/pratics/models"

type Teacher interface {
	InsertTeacher(name string) (*models.Teacher, error)
}
