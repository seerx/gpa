package repos

import "github.com/seerx/gpa/examples/pratics/models"

type Update interface {
	UpdateXuByIdAndAge(user *models.User) error
	UpdateByName(user *models.User, gender bool) (int64, error)
	UpdateByAgeAndCret(user *models.User) (int64, *models.User, error)
	UpdateXByAge(age int, URL string) (int64, models.User, error)
	// sql:update "user" set  "name"=:name,url=:URL,cret=:cret  where id=:id and age>:age
	UpdateName(user *models.User, name string) (int64, error)
	// sql:update user   where id=:id and age>:age
	UpdateXName(user *models.User, name string) (int64, error)
	UpdateYName(name string) (int64, models.User, error)
	// sql:where id=:id
	Update1Name(name string, id uint64) (int64, models.User, error)
	// sql:where name=:name or id in :ids
	Update2Name(user models.User, ids []uint64) (int64, error)
}
