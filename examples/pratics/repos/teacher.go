package repos

import "github.com/seerx/gpa/examples/pratics/models"

type Teacher interface {
	InsertTeacher(name string) (*models.Teacher, error)
	Insert1Teacher(tc *models.Teacher) (*models.Teacher, error)

	UpdateByID(tc *models.Teacher) error

	FindByID(id int64) (*models.Teacher, error)
	FindAll() ([]*models.Teacher, error)
	FindMap(kg func(*models.Teacher) int64) (map[int64]*models.Teacher, error)
	FindCallbck(cb func(*models.Teacher) error) error
}
