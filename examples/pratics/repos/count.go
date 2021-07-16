package repos

import "github.com/seerx/gpa/examples/pratics/models"

type Count interface {
	CountByName(user *models.User) (int64, models.User, error)
	// sql:select count(*) from "user"  where id=:id and age>:age
	CountName(user *models.User, name string) (int64, error)
	// sql:select count(1) where id = :id and age>:age
	CountXName(user *models.User, name string) (int64, error)
}
