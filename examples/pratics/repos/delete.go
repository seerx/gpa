package repos

import "github.com/seerx/gpa/examples/pratics/models"

type Delete interface {
	DeleteByName(user *models.User, gender bool) (int64, error)
	DeleteByAgeAndCret(user *models.User) (int64, *models.User, error)
	DeleteXByAge(age int, URL string) (int64, models.User, error)

	// sql:delete "user"  where id=:id and age in :ages
	DeleteName(user *models.User, name string, ages []int) (int64, error)
	// sql:delete user   where id = :id and age>:age
	DeleteXName(user *models.User, name string) (int64, error)
	DeleteAll() (int64, models.User, error)
	// sql:where id=:id
	Delete1Name(id uint64) (int64, models.User, error)
	// sql:where id=:id
	Delete2Name(user models.User, id uint64) (int64, error)
}
