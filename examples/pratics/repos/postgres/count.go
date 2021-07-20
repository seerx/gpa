//+mro-ignore
// DO NOT EDIT THIS FILE
// Generated by mro at 2021-07-20 16:36:10
package postgres

import (
	models "github.com/seerx/gpa/examples/pratics/models"
	rt81 "github.com/seerx/gpa/rt"
)

type Count struct {
	p *rt81.Provider
}

func (count *Count) CountByName(user *models.User) (int64, models.User, error) {
	var err error
	sql := "SELECT count(0) FROM user WHERE name=? "
	var1 := count.p.Executor().QueryRow(sql, user.Name)
	var var2 int64
	err = var1.Scan(&var2)
	if err != nil {
		return 0, models.User{}, err
	}
	return var2, models.User{}, nil
}

func (count *Count) CountName(user *models.User, name string) (int64, error) {
	var err error
	sql := "SELECT count(*) FROM user WHERE id=? and age>?"
	var1 := count.p.Executor().QueryRow(sql, user.ID, user.Age)
	var var2 int64
	err = var1.Scan(&var2)
	if err != nil {
		return 0, err
	}
	return var2, nil
}

func (count *Count) CountXName(user *models.User, name string) (int64, error) {
	var err error
	sql := "SELECT count(1) FROM user WHERE id = ? and age>?"
	var1 := count.p.Executor().QueryRow(sql, user.ID, user.Age)
	var var2 int64
	err = var1.Scan(&var2)
	if err != nil {
		return 0, err
	}
	return var2, nil
}
