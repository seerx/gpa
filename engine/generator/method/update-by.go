package method

import (
	"strings"

	"github.com/seerx/gpa/engine/generator/defines"
	rdesc "github.com/seerx/gpa/engine/generator/repo-desc"
	"github.com/seerx/gpa/engine/generator/sqlgenerator"
	"github.com/seerx/gpa/rt/dbutil"
	"github.com/seerx/logo/log"
)

type updateby struct {
	BaseMethod
}

func (g *updateby) Test(fn *defines.Func) bool {
	if strings.Index(fn.Name, "Update") == 0 {
		if strings.Index(fn.Name, "By") > 0 {
			g.BaseMethod.Test(fn)
			fn.Template = defines.UPDATE
			return true
		}
	}
	return false
}

// func (g *updateby) where(fd *desc.FuncDesc, whereParams []*desc.FieldVarPair) ([]*desc.SQLParam, map[string]bool, error) {
// 	var params []*desc.SQLParam
// 	var fieldMap = map[string]bool{}
// 	// 组织 where 参数
// 	for _, p := range whereParams {
// 		// name = utils.LowerFirstChar(name)
// 		fieldMap[p.FieldName] = true
// 		found := false
// 		for _, arg := range g.fn.Args {
// 			if arg.Type.IsContext() || arg.Type.IsStruct() || arg.Type.IsFunc() {
// 				// 以上三种类型不能作为 where 参数
// 				continue
// 			}

// 			if arg.Name == utils.LowerFirstChar(p.VarName) {
// 				// 找到输入参数
// 				params = append(params, &desc.SQLParam{
// 					VarName: arg.Name,
// 				})
// 				found = true
// 			}
// 		}
// 		if !found {
// 			if fd.Input.Bean != nil {
// 				for _, f := range fd.Input.Bean.Type.Bean.Fields {
// 					if f.GoVarName == p.VarName && !f.Type.IsJson() {
// 						// 找到输入参数
// 						if f.SQLType.IsTime() {
// 							fd.DBUtilPackage = g.fn.AddDBUtilPackage()
// 							var timeProp = &dbutil.TimePropDesc{
// 								TypeName: f.SQLType.Name,
// 								Nullable: f.Nullable,
// 							}
// 							if f.TimeZone != nil {
// 								timeProp.TimeZone = f.TimeZone.String()
// 							}
// 							params = append(params, &desc.SQLParam{
// 								VarName:  fd.Input.Bean.Name + "." + f.GoVarName,
// 								VarAlias: fd.NextVarName(),
// 								Time:     true,
// 								TimeProp: timeProp,
// 							})
// 						} else {
// 							params = append(params, &desc.SQLParam{
// 								VarName: fd.Input.Bean.Name + "." + f.GoVarName,
// 							})
// 						}
// 						found = true
// 					}
// 				}
// 			}
// 		}
// 		if !found {
// 			return nil, nil, g.fn.Error("no where param [%s] found in func args", p.VarName)
// 			// return
// 		}
// 	}
// 	return params, fieldMap, nil
// }

// func (g *updateby) where(fd *rdesc.FuncDesc, whereParams []*intf.SQLParam) (map[string]bool, error) {
// 	// var params []*desc.SQLParam
// 	var fieldMap = map[string]bool{}
// 	// 组织 where 参数
// 	for _, p := range whereParams {
// 		// name = utils.LowerFirstChar(name)
// 		if p.IsInOperator {
// 			fd.DBUtilPackage = g.fn.AddDBUtilPackage()
// 		}
// 		fieldMap[p.SQLParamFieldName] = true
// 		found := g.findParamInFuncArgs(p, true, func(argName, sqlParamName string) bool {
// 			return argName == sqlParamName || argName == names.LowerFirstChar(sqlParamName)
// 		})
// 		if !found {
// 			found = g.findParamInFuncArgs(p, true, func(argName, sqlParamName string) bool {
// 				return strings.EqualFold(argName, sqlParamName)
// 			})
// 			if found {
// 				g.logger.Warnf(g.fn.Format("variable %s as sql param %s", p.VarName, p.SQLParamName))
// 			}
// 		}

// 		// for _, arg := range g.fn.Params {
// 		// 	if arg.Type.IsContext() || arg.Type.IsStruct() || arg.IsFunc {
// 		// 		// 以上三种类型不能作为 where 参数
// 		// 		continue
// 		// 	}

// 		// 	if arg.Name == p.SQLParamName || arg.Name == names.LowerFirstChar(p.SQLParamName) {
// 		// 		// 找到输入参数
// 		// 		p.VarName = arg.Name
// 		// 		found = true
// 		// 	}
// 		// 	if !found {
// 		// 		if strings.EqualFold(arg.Name, p.SQLParamName) {
// 		// 			p.VarName = arg.Name
// 		// 			found = true
// 		// 		}
// 		// 	}
// 		// }
// 		if !found && fd.Input.Bean != nil {
// 			bean, err := fd.Input.Bean.GetBeanType()
// 			if err != nil {
// 				return nil, g.fn.CreateError("get ben type error: %s", err.Error())
// 			}
// 			found = g.findParamInBeanFieldsAndFill(fd, bean, p, func(varName, sqlParamName string) bool {
// 				return varName == sqlParamName || varName == names.UpperFirstChar(sqlParamName)
// 			})
// 			if !found {
// 				found = g.findParamInBeanFieldsAndFill(fd, bean, p, func(varName, sqlParamName string) bool {
// 					return strings.EqualFold(varName, sqlParamName)
// 				})
// 				if found {
// 					g.logger.Warnf(g.fn.Format("variable %s as sql param %s", p.VarName, p.SQLParamName))
// 				}
// 			}
// 			// for _, f := range bean.Fields {
// 			// 	if !f.IsJSON {
// 			// 		if f.VarName == p.SQLParamName || f.VarName == names.UpperFirstChar(p.SQLParamName) {
// 			// 			// 找到输入参数
// 			// 			p.VarName = fd.Input.Bean.Name + "." + f.VarName

// 			// 			if f.Type.IsTime() {
// 			// 				fd.DBUtilPackage = g.fn.AddDBUtilPackage()
// 			// 				var timeProp = &dbutil.TimePropDesc{
// 			// 					TypeName: f.Type.Name,
// 			// 					Nullable: f.Nullable,
// 			// 				}
// 			// 				if f.TimeZone != nil {
// 			// 					timeProp.TimeZone = f.TimeZone.String()
// 			// 				}
// 			// 				p.VarName = fd.Input.Bean.Name + "." + f.VarName
// 			// 				p.VarAlias = fd.NextVarName()
// 			// 				p.Time = true
// 			// 				p.TimeProp = timeProp
// 			// 			}
// 			// 			found = true
// 			// 		}
// 			// 	}
// 			// }
// 		}
// 		if !found {
// 			return nil, g.fn.CreateError("no where param [%s] found in func args", p.VarName)
// 			// return
// 		}
// 	}
// 	return fieldMap, nil
// }

func (g *updateby) Parse() (*rdesc.FuncDesc, error) {
	byIdx := strings.Index(g.fn.Name, "By")
	if byIdx < 6 {
		return nil, g.fn.CreateError("invalid name of UpdateXyzBy func")
	}

	whereInName := g.fn.Name[byIdx+2:]
	sqlWhere, whereParams, err := parseWhereFromFuncName(whereInName)
	if err != nil {
		err = g.fn.CreateError(err.Error())
		log.Error(err, "pasre where from name error")
		return nil, err
	}

	fd := rdesc.NewFuncDesc(g.fn, 3, g.logger)
	if err := fd.Explain(); err != nil {
		err = g.fn.CreateError(err.Error())
		log.Error(err, "expain error")
		return nil, err
	}
	// wpmap := map[string]bool{}
	// for _, wp := range whereParams {
	// 	wpmap[wp.SQLParamFieldName] = true
	// }
	if err := fd.ExplainSetBeanFieldsValueWithArgs(nil); err != nil {
		err = g.fn.CreateError(err.Error())
		log.Error(err, "expain set values error")
		return nil, err
	}

	if fd.BeanObj == nil {
		err := g.fn.CreateError("no struct bean found in funcion")
		log.Error(err, "")
		return nil, err
	}
	if err := g.CheckUpdateReturns(fd); err != nil {
		return nil, g.fn.CreateError(err.Error())
	}

	paramFieldMap, err := g.prepareParams(fd, whereParams, false)
	if err != nil {
		return nil, g.fn.CreateError(err.Error())
	}

	// sql := strings.Builder{}
	// sqlParams := []*desc.SQLParam{}
	// columns := []string{}
	bean, err := fd.BeanObj.GetBeanType()
	if err != nil {
		return nil, g.fn.CreateError(err.Error())
	}
	var sql = sqlgenerator.SQL{
		TableName:   bean.TableName,
		Where:       sqlWhere,
		WhereParams: whereParams,
	}
	for _, f := range bean.Fields {
		if f.Ignore {
			continue
		}
		if f.IsPrimaryKey {
			// 主键不进行 update 操作
			continue
		}
		if _, ok := paramFieldMap[f.FieldName()]; ok {
			// 在 where 参数中出现且有参数值的字段不进行 update 操作
			continue
		}

		arg := g.fn.FindParam(f.GetArgNames())
		// arg, ok := g.fn.ArgMap[f.ArgName]
		if fd.Input.Bean == nil {
			// 输入参数中，没有与 beanObject 一致的对象
			if arg == nil {
				// log.Warn(g.fn.Format("input arg [%s] may be ignored"))
				continue
			}

			// 输入参数中找到对应的字段
			sql.Params = append(sql.Params, &sqlgenerator.SQLParam{
				VarName: arg.Name,
			})
		} else {
			// 输入参数中有与 beanObject 一致的对象
			// 使用 输入参数中有与 beanObject 一致的对象 的数据作为参数

			varAliasName := ""
			isBlob := false
			isJSON := false
			isTime := false

			if f.Field.Type.IsStruct() {
				obj := g.fn.MakeObject(&f.Field)
				fb, err := obj.GetBeanType()
				if err != nil {
					return nil, err
				}
				varAliasName = fd.NextVarName()
				isBlob = fb.IsBlobReadWriter()
				if isBlob {
					fd.DBUtilPackage = g.fn.AddDBUtilPackage()
				}
			}

			var timeProp *dbutil.TimePropDesc
			if !isBlob {
				if f.IsJSON {
					varAliasName = fd.NextVarName()
					fd.DBUtilPackage = g.fn.AddDBUtilPackage()
					isJSON = true
					// }
				} else if f.Type.IsTime() {
					varAliasName = fd.NextVarName()
					fd.DBUtilPackage = g.fn.AddDBUtilPackage()
					isTime = true
					timeProp = &dbutil.TimePropDesc{
						TypeName: f.Type.Name,
						Nullable: f.Nullable,
					}
					if f.TimeZone != nil {
						timeProp.TimeZone = f.TimeZone.String()
					}
				}
			}

			sql.Params = append(sql.Params, &sqlgenerator.SQLParam{
				VarName:  fd.Input.Bean.Name + "." + f.VarName,
				VarAlias: varAliasName,
				JSON:     isJSON,
				Blob:     isBlob,
				Time:     isTime,
				TimeProp: timeProp,
			})
			// }
		}
		sql.Columns = append(sql.Columns, f.FieldName()+"=?")
	}

	if len(sql.Columns) <= 0 {
		return nil, g.fn.CreateError("no field to be update in this func")
	}

	// _, err = sql.WriteString("UPDATE " + bean.TableName + " SET ")
	// if err != nil {
	// 	return nil, err
	// }
	// _, err = sql.WriteString(strings.Join(columns, ",") + " WHERE " + sqlWhere)
	// if err != nil {
	// 	return nil, err
	// }

	fd.SQL, fd.SQLParams, fd.SQLWhereParams = g.sqlg.Update(&sql) // sql.CreateUpdate() //   append(sqlParams, whereParams...)

	return fd, nil
}
