package method

import (
	"fmt"
	"strings"

	"github.com/seerx/gpa/engine/generator/defines"
	rdesc "github.com/seerx/gpa/engine/generator/repo-desc"
	"github.com/seerx/gpa/engine/generator/sqlgenerator"
	"github.com/seerx/gpa/engine/generator/xtype"
	"github.com/seerx/gpa/engine/sql/names"
	"github.com/seerx/gpa/logger"
	"github.com/seerx/gpa/rt/dbutil"
)

type BaseMethod struct {
	// dialect intf.Dialect
	sqlg   sqlgenerator.SQLGenerator
	fn     *defines.Func
	logger logger.GpaLogger
}

func (bg *BaseMethod) Test(fn *defines.Func) bool {
	bg.fn = fn
	return false
}

func (bg *BaseMethod) CheckFindReturns(fd *rdesc.FuncDesc) error {
	if fd.Result.Bean != nil {
		// 此时返回值应该是 (struct, error) 形式
		// 下面的判断可以忽略
		// if !fd.Result.List[0].Type.IsStruct() {
		// 	return bg.fn.Error("the first return must be bean struct")
		// }

		if fd.Result.List[0].IsMap {
			// 如果返回的是 map
			// 在参数中查找 MapKey 函数
			// func(*beanStruct) keyType
			foundKeyFunc := false
			for _, arg := range fd.Input.Args {
				if arg.Arg.IsFunc {
					if len(arg.Arg.Params) == 1 && arg.Arg.Params[0].Type.Equals(&fd.BeanObj.Type) {
						fd.Input.KeyGeneratorArgIsPtr = arg.Arg.Params[0].Type.IsPtr
						if len(arg.Arg.Results) != 1 {
							return bg.fn.CreateError("the key generate func shuld only has one result")
						}
						if arg.Arg.Results[0].Type.EqualsExactly(fd.Result.List[0].Key) {
							arg.IsMapKeyFunc = true
							fd.Input.KeyGenerator = arg.Arg
							// arg.Arg.Args
							foundKeyFunc = true
							// fd.BeanObj = fd.Result.List[0]
						} else {
							return bg.fn.CreateError("the key generate func return type is not match map's key type")
						}
					}
				}
			}
			if !foundKeyFunc {
				return bg.fn.CreateError("no key generate func found in args")
			}
		}
		fd.Result.FindOne = !fd.Result.List[0].IsSlice && !fd.Result.List[0].IsMap
	} else {
		// 此时返回值应该使用回调函数完成
		// 回调函数的形式为 func(*beanStruct) error
		foundReturnFunc := false
		for _, arg := range fd.Input.Args {
			if arg.Arg.IsFunc {
				if len(arg.Arg.Params) == 1 && arg.Arg.Params[0].Type.IsCustom() {
					// if arg.Arg.Args[0].Ptr != fd.BeanObj.Ptr {
					// 	if arg.Arg.Args[0].Ptr {
					// 		fd.Input.KeyGeneratorArgPrefix = "&"
					// 	} else {
					// 		fd.Input.KeyGeneratorArgPrefix = "*"
					// 	}
					// }
					fd.Input.CallbackArgIsPtr = arg.Arg.Params[0].Type.IsPtr
					if len(arg.Arg.Results) == 1 && arg.Arg.Results[0].Type.IsError() {
						fd.BeanObj = bg.fn.MakeObject(arg.Arg.Params[0])
						fd.BeanVarName = fd.NextVarName()
						fd.BeanTypeName = fd.BeanObj.Type.String()
						arg.IsReturnFunc = true
						fd.Input.Callback = arg.Arg
						foundReturnFunc = true
						break
					}
				}
			}
		}
		if !foundReturnFunc {
			return bg.fn.CreateError("no callback func to return results found in args")
		}
	}

	// 查找参数 bean
	for _, arg := range fd.Input.Args {
		if arg.Arg.Type.IsCustom() {
			if fd.ParamBeanObj == nil {
				fd.ParamBeanObj = bg.fn.MakeObject(arg.Arg)
			} else if arg.Arg.Type.Equals(&fd.BeanObj.Type) {
				fd.ParamBeanObj = bg.fn.MakeObject(arg.Arg)
				break
			}
		}
	}

	fd.SQLReturnVarName = fd.NextVarName()
	fd.SQLPackage = bg.fn.AddSQLPackage()
	return nil
}

func (bg *BaseMethod) CheckUpdateReturns(fd *rdesc.FuncDesc) error {
	if fd.Result.AffectVar != "" {
		// 如果返回 Affect 行数,其返回值的第一个参数必须是 int64
		// 可能是 (int64, error) 或者 (int64, struct, error)
		if !fd.Result.List[0].Type.IsInt64() {
			return bg.fn.CreateError("the first return must be int64")
		}
		if fd.Result.Bean != nil {
			// 此时返回值应该是 (int64, struct, error) 形式
			if !fd.Result.List[1].Type.IsCustom() {
				return bg.fn.CreateError("the second return must be bean struct")
			}
		}
	} else {
		// 不返回 Affect 行数
		// 可能是 (struct, error) 或者 error
		if fd.Result.Bean != nil {
			// 此时返回值应该是 (struct, error) 形式
			if !fd.Result.List[0].Type.IsCustom() {
				return bg.fn.CreateError("the first return must be bean struct")
			}
		}
	}
	return nil
}

func (bg *BaseMethod) CheckDeleteReturns(fd *rdesc.FuncDesc) error {
	if fd.Result.AffectVar != "" {
		// 如果返回 Affect 行数,其返回值的第一个参数必须是 int64
		// 可能是 (int64, error) 或者 (int64, struct, error)
		if !fd.Result.List[0].Type.IsInt64() {
			return bg.fn.CreateError("the first return must be int64")
		}
		if fd.Result.Bean != nil {
			// 此时返回值应该是 (int64, struct, error) 形式
			if !fd.Result.List[1].Type.IsCustom() {
				return bg.fn.CreateError("the second return must be bean struct")
			}
		}
	} else {
		// 不返回 Affect 行数
		// 可能是 (struct, error) 或者 error
		if fd.Result.Bean != nil {
			// 此时返回值应该是 (struct, error) 形式
			if !fd.Result.List[0].Type.IsCustom() {
				return bg.fn.CreateError("the first return must be bean struct")
			}
		}
	}
	return nil
}

func (bg *BaseMethod) CheckCountReturns(fd *rdesc.FuncDesc) error {
	// 第一个返回值必须是 int64
	if !fd.Result.List[0].Type.IsInt64() {
		return bg.fn.CreateError("the first return must be int64")
	}
	// 如果返回 bean struct ，则为第二个返回值
	if fd.Result.Bean != nil {
		// 此时返回值应该是 (struct, error) 形式
		if !fd.Result.List[1].Type.IsCustom() {
			return bg.fn.CreateError("the first return must be bean struct")
		}
	}
	// fd.Result.CountVar = fd.NextVarName()
	return nil
}

func (bg *BaseMethod) findParamInFuncArgs(p *sqlgenerator.SQLParam, asWhere bool, matchFn func(argName, sqlParamName string) bool) (bool, error) {
	for _, arg := range bg.fn.Params {
		if asWhere {
			if arg.Type.IsContext() || arg.Type.IsCustom() || arg.IsFunc {
				// 以上三种类型不能作为 where 参数
				continue
			}
		} else {
			if arg.Type.IsContext() || arg.IsFunc {
				// 以上两种类型不能作为 set 参数
				continue
			}
		}

		if matchFn(arg.Name, p.SQLParamName) {
			p.VarName = arg.Name
			if p.IsInOperator {
				// IN 操作符
				if !arg.IsSlice {
					return false, fmt.Errorf("arg [%s] shuld be array", p.VarName)
				}
			}

			return true, nil
		}

		// if arg.Name == p.SQLParamName || arg.Name == names.LowerFirstChar(p.SQLParamName) {
		// 	// 找到输入参数
		// 	p.VarName = arg.Name
		// 	found = true
		// }
		// if !found {
		// 	if strings.EqualFold(arg.Name, p.SQLParamName) {
		// 		p.VarName = arg.Name
		// 		found = true
		// 	}
		// }
	}
	return false, nil
}

func (bg *BaseMethod) findParamInBeanFieldsAndFill(fd *rdesc.FuncDesc, bean *xtype.XType, p *sqlgenerator.SQLParam, matchFn func(varName, sqlParamName string) bool) (bool, error) {
	for _, f := range bean.Fields {
		if !f.IsJSON {
			if matchFn(f.VarName, p.SQLParamName) {
				p.VarName = fd.Input.Bean.Name + "." + f.VarName
				if p.IsInOperator {
					// IN 操作符
					if !f.Field.IsSlice {
						return false, fmt.Errorf("arg [%s] shuld be array", p.VarName)
					}
				}
				// if f.VarName == p.SQLParamName || f.VarName == names.UpperFirstChar(p.SQLParamName) {
				// 找到输入参数

				if f.Type.IsTime() {
					fd.DBUtilPackage = bg.fn.AddDBUtilPackage()
					var timeProp = &dbutil.TimePropDesc{
						TypeName: f.Type.Name,
						Nullable: f.Nullable,
					}
					if f.TimeZone != nil {
						timeProp.TimeZone = f.TimeZone.String()
					}
					p.VarName = fd.Input.Bean.Name + "." + f.VarName
					p.VarAlias = fd.NextVarName()
					p.Time = true
					p.TimeProp = timeProp
				}
				return true, nil
			}
		}
	}
	return false, nil
}

func (bg *BaseMethod) prepareParams(fd *rdesc.FuncDesc, params []*sqlgenerator.SQLParam, forSet bool) (map[string]bool, error) {
	// var params []*desc.SQLParam
	var fieldMap = map[string]bool{}
	// 组织 where 参数
	for _, p := range params {
		// name = utils.LowerFirstChar(name)
		if p.IsInOperator {
			fd.DBUtilPackage = bg.fn.AddDBUtilPackage()
		}
		// var err error
		fieldMap[p.SQLParamFieldName] = true
		found, err := bg.findParamInFuncArgs(p, !forSet, func(argName, sqlParamName string) bool {
			return argName == sqlParamName || argName == names.LowerFirstChar(sqlParamName)
		})
		if err != nil {
			bg.logger.Error(err, bg.fn.Format("find param [%s] in func args error", p.VarName))
		}
		if !found {
			found, err = bg.findParamInFuncArgs(p, !forSet, func(argName, sqlParamName string) bool {
				return strings.EqualFold(argName, sqlParamName)
			})
			if err != nil {
				bg.logger.Error(err, bg.fn.Format("find param [%s] in func args error", p.VarName))
			}
			if found {
				bg.logger.Warnf(bg.fn.Format("variable %s as sql param %s", p.VarName, p.SQLParamName))
			}
		}

		// for _, arg := range g.fn.Params {
		// 	if arg.Type.IsContext() || arg.Type.IsStruct() || arg.IsFunc {
		// 		// 以上三种类型不能作为 where 参数
		// 		continue
		// 	}

		// 	if arg.Name == p.SQLParamName || arg.Name == names.LowerFirstChar(p.SQLParamName) {
		// 		// 找到输入参数
		// 		p.VarName = arg.Name
		// 		found = true
		// 	}
		// 	if !found {
		// 		if strings.EqualFold(arg.Name, p.SQLParamName) {
		// 			p.VarName = arg.Name
		// 			found = true
		// 		}
		// 	}
		// }
		if !found && fd.Input.Bean != nil {
			bean, err := fd.Input.Bean.GetBeanType()
			if err != nil {
				return nil, bg.fn.CreateError("get ben type error: %s", err.Error())
			}
			found, err = bg.findParamInBeanFieldsAndFill(fd, bean, p, func(varName, sqlParamName string) bool {
				return varName == sqlParamName || varName == names.UpperFirstChar(sqlParamName)
			})
			if err != nil {
				bg.logger.Error(err, bg.fn.Format("find param [%s] in bean's field error", p.VarName))
			}
			if !found {
				found, err = bg.findParamInBeanFieldsAndFill(fd, bean, p, func(varName, sqlParamName string) bool {
					return strings.EqualFold(varName, sqlParamName)
				})
				if err != nil {
					bg.logger.Error(err, bg.fn.Format("find param [%s] in bean's field error", p.VarName))
				}
				if found {
					bg.logger.Warnf(bg.fn.Format("variable %s as sql param %s", p.VarName, p.SQLParamName))
				}
			}
			// for _, f := range bean.Fields {
			// 	if !f.IsJSON {
			// 		if f.VarName == p.SQLParamName || f.VarName == names.UpperFirstChar(p.SQLParamName) {
			// 			// 找到输入参数
			// 			p.VarName = fd.Input.Bean.Name + "." + f.VarName

			// 			if f.Type.IsTime() {
			// 				fd.DBUtilPackage = g.fn.AddDBUtilPackage()
			// 				var timeProp = &dbutil.TimePropDesc{
			// 					TypeName: f.Type.Name,
			// 					Nullable: f.Nullable,
			// 				}
			// 				if f.TimeZone != nil {
			// 					timeProp.TimeZone = f.TimeZone.String()
			// 				}
			// 				p.VarName = fd.Input.Bean.Name + "." + f.VarName
			// 				p.VarAlias = fd.NextVarName()
			// 				p.Time = true
			// 				p.TimeProp = timeProp
			// 			}
			// 			found = true
			// 		}
			// 	}
			// }
		}
		if !found {
			return nil, bg.fn.CreateError("no where param [%s] found in func args", p.VarName)
			// return
		}
	}
	return fieldMap, nil
}
