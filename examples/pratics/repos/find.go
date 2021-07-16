package repos

import (
	"context"

	"github.com/seerx/gpa/examples/pratics/models"
)

type Find interface {
	//sql: select id as id, name name, sum(age) as Age where id=:id group by name
	FindA(ctx context.Context, id uint64) (*models.User, error)
	FindById(id uint64) (models.User, error)
	// sql: select 	* from user where name like :nm
	FindUsers(nm string) ([]*models.User, error)

	// sql: select * from user where id in :ids
	FindUsers1(ids []uint64) ([]*models.User, error)
	FindMapByName(name string, kg func(*models.User) uint64) (map[uint64]*models.User, error)
	FindCbById(id uint64, fn func(*models.User) error) error
	FindMapById(id uint64, fn func(*models.User) uint64) (map[uint64]*models.User, error)
}
