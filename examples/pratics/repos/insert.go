//mro=github.com/seerx/mro
//接口敌营

package repos

// ***
// 函数参数中的参数名称(sql语句中参数名称)与 bean struct 成员对应关系
// ID    ID
// id    Id
// name  Name
// 函数参数中的参数名称 非必要首字母使用小写
// 如果 bean struct 中，使用全大写定义的 field ,函数参数中的参数名称 也是用相同的名称定义
// 全大写定义参见: utils.LintGonicMapper 中的内容，也可以使用首字符大写定义
// ***

type UserRepo interface {
	// Insert(ctx context.Context, user *models.User, name string) (*models.User, error)
	// InsertA(user *models.User) error
	// InsertB(user *models.User, DD string, URL string) (models.User, error)
	// InsertC(user *models.User) (*models.User, error)
	// InsertUser(name string) (*models.User, error)
}
