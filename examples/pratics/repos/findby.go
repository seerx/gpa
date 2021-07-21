package repos

import (
	"context"

	"github.com/seerx/gpa/examples/pratics/models"
)

type FindBy interface {
	FindAById(ctx context.Context, id uint64) (*models.User, error)
	FindById(id uint64) (models.User, error)
	FindXuByName(name string) ([]*models.User, error)
	FindMapByName(name string, kg func(*models.User) (uint64, error)) (map[uint64]*models.User, error)
	FindCbById(id uint64, fn func(*models.User) error) error
	// FindMapById(id uint64, fn func(*models.User) uint64) (map[uint64]*models.User, error)
}
