package method

import (
	"fmt"
	"strings"

	"github.com/seerx/gpa/engine/generator/defines"
	rdesc "github.com/seerx/gpa/engine/generator/repo-desc"
	"github.com/seerx/gpa/engine/generator/xtype"
	"github.com/seerx/gpa/engine/sql/dialect/intf"
	"github.com/seerx/gpa/rt/dbutil"
)

type update struct {
	BaseMethod
}

func (g *update) Test(fn *defines.Func) bool {
	// 如果函数名称中有 By 字符串，则优先为 UpdateXBy 操作
	if strings.Index(fn.Name, "Update") == 0 {
		// if fn.SQL != "" {
		// Update 操作必须提供 SQL 注释
		// fn.SQL 的形式
		// 1 update table set field=:arg1,field1=:arg2 where id=:arg3
		//		类似原生 SQL 语句，arg1 arg2 arg3 从函数的参数中获取，不需要提供表相关 struct 类型数据
		// 2 set field=:arg1,field1=:arg2 where id=:arg3
		//      需要提供表相关的 struct 结构，arg1 arg2 arg3 从函数的参数中获取
		// 3 where id=:arg3
		//		需要提供表相关的 struct 结构，根据参数自行生成 set field=? 语句
		// 4 空，没有 where 条件
		g.BaseMethod.Test(fn)
		fn.Template = defines.UPDATE
		return true
		// }
	}

	return false
}

// func lastIndex(s string, substrs ...string) (string, int) {
// 	for _, ss := range substrs {
// 		if p := strings.LastIndex(s, ss); p >= 0 {
// 			return ss, p
// 		}
// 	}
// 	return "", -1
// }

func parseSQL(sql string) (*intf.SQL, error) {
	terms, err := splitSQL(sql)
	if err != nil {
		return nil, err
	}
	updateIndex := -1
	setIndex := -1
	whereIndex := -1
	for n, t := range terms {
		term := strings.ToLower(t)
		switch term {
		case "where":
			if whereIndex == -1 {
				whereIndex = n
			}
		case "update":
			if updateIndex == -1 {
				updateIndex = n
			}
		case "set":
			if setIndex == -1 {
				setIndex = n
			}
		}
	}

	var s intf.SQL
	if setIndex-updateIndex == 2 {
		s.TableName = terms[updateIndex+1]
	}

	if setIndex >= 0 {
		setEnd := len(terms)
		if whereIndex > 0 {
			setEnd = whereIndex
		}

		setSQL := ""

		for n := setIndex + 1; n < setEnd; n++ {
			col := terms[n]

			ps, err := FindParams(col)
			if err != nil {
				return nil, err
			}
			var termParams []*intf.SQLParam
			for m := len(ps) - 1; m >= 0; m-- {
				var fieldName string
				for k := n; k >= 0; k-- {
					term := terms[k]
					eqPos := strings.LastIndex(term, "=")
					if eqPos > 0 {
						commaPos := strings.LastIndex(term, ",")
						if commaPos > 0 {
							fieldName = term[commaPos+1 : eqPos]
						} else {
							fieldName = term[:eqPos]
						}
						break
					} else if eqPos == 0 {
						// 从上一个 term 开始找
						if k > 0 {
							term := terms[k-1]
							if term == "=" {
								if k > 1 {
									term = terms[k-2]
								} else {
									return nil, fmt.Errorf("cann't find field name of params %s", ps[m].Name)
								}
							}
							commaPos := strings.LastIndex(term, ",")
							if commaPos > 0 {
								fieldName = term[commaPos+1 : eqPos]
							} else {
								fieldName = term[:eqPos]
							}
							break
						} else {
							return nil, fmt.Errorf("cann't find field name of params %s", ps[m].Name)
						}
					}
				}
				if fieldName == "" {
					return nil, fmt.Errorf("cann't find field name of params %s", ps[m].Name)
				}

				col = ReplaceParam(col, ps[m], "?")
				termParams = append(termParams, &intf.SQLParam{
					SQLParamName:      ps[m].Name,
					SQLParamFieldName: fieldName,
				})
				// s.Params = append([]*desc.SQLParam{{}}, s.Params...)
			}
			setSQL += col
			for n := len(termParams) - 1; n >= 0; n-- {
				s.Params = append(s.Params, termParams[n])
			}
		}
		s.Columns = append(s.Columns, setSQL)
	}

	whereSQL, whereParams, err := ParseWhere(terms, whereIndex)
	if err != nil {
		return nil, err
	}
	s.Where, s.WhereParams = whereSQL, whereParams
	return &s, nil
}

// func (g *update) fillParamProp(fd *rdesc.FuncDesc, params []*intf.SQLParam, forSetParam bool) error {
// 	// var params []*desc.SQLParam
// 	// 组织 where 参数
// 	for _, p := range params {
// 		if p.IsInOperator {
// 			fd.DBUtilPackage = g.fn.AddDBUtilPackage()
// 		}
// 		found := g.findParamInFuncArgs(p, !forSetParam, func(argName, sqlParamName string) bool {
// 			return argName == sqlParamName || argName == names.LowerFirstChar(sqlParamName)
// 		})
// 		if !found {
// 			found = g.findParamInFuncArgs(p, !forSetParam, func(argName, sqlParamName string) bool {
// 				return strings.EqualFold(argName, sqlParamName)
// 			})
// 			if found {
// 				g.logger.Warnf(g.fn.Format("variable %s as sql param %s", p.VarName, p.SQLParamName))
// 			}
// 		}

// 		// for _, arg := range g.fn.Params {
// 		// 	if forSetParam {
// 		// 		if arg.Type.IsContext() || arg.IsFunc {
// 		// 			// 以上三种类型不能作为 set 参数
// 		// 			continue
// 		// 		}
// 		// 	} else {
// 		// 		if arg.Type.IsContext() || arg.Type.IsStruct() || arg.IsFunc {
// 		// 			// 以上三种类型不能作为 where 参数
// 		// 			continue
// 		// 		}
// 		// 	}

// 		// 	if arg.Name == p.SQLParamName {
// 		// 		// 找到输入参数
// 		// 		p.VarName = arg.Name
// 		// 		found = true
// 		// 	}
// 		// }
// 		if !found && fd.Input.Bean != nil {
// 			bean, err := fd.Input.Bean.GetBeanType()
// 			if err != nil {
// 				return g.fn.CreateError("get ben type error: %s", err.Error())
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
// 			// 	if f.GoVarName == utils.UpperFirstChar(p.SQLParamName) && !f.Type.IsJson() {
// 			// 		// 找到输入参数
// 			// 		p.VarName = fd.Input.Bean.Name + "." + f.GoVarName
// 			// 		if f.SQLType.IsTime() {
// 			// 			fd.DBUtilPackage = g.fn.AddDBUtilPackage()
// 			// 			var timeProp = &dbutil.TimePropDesc{
// 			// 				TypeName: f.SQLType.Name,
// 			// 				Nullable: f.Nullable,
// 			// 			}
// 			// 			if f.TimeZone != nil {
// 			// 				timeProp.TimeZone = f.TimeZone.String()
// 			// 			}
// 			// 			p.VarAlias = fd.NextVarName()
// 			// 			p.Time = true
// 			// 			p.TimeProp = timeProp
// 			// 		}
// 			// 		found = true
// 			// 		break
// 			// 	}
// 			// }
// 		}
// 		if !found {
// 			return g.fn.CreateError("no where param [%s] found in func args", p.SQLParamName)
// 			// return
// 		}
// 	}
// 	return nil
// }

func (g *update) Parse() (*rdesc.FuncDesc, error) {
	xsql := g.fn.SQL
	sql, err := parseSQL(xsql)
	hasSetFieldsInSQL := sql.Columns != nil
	hasWhereInSQL := sql.WhereParams != nil
	if err != nil {
		return nil, g.fn.CreateError(err.Error())
	}
	var whereParamsMap map[string]bool
	if hasWhereInSQL {
		whereParamsMap = map[string]bool{}
		for _, p := range sql.WhereParams {
			whereParamsMap[p.SQLParamFieldName] = true
		}
	}

	fd := rdesc.NewFuncDesc(g.fn, 3, g.logger)
	if err := fd.Explain(); err != nil {
		return nil, g.fn.CreateError("explain error: %s", err.Error())
	}
	if err := fd.ExplainSetBeanFieldsValueWithArgs(whereParamsMap); err != nil {
		return nil, g.fn.CreateError("explain set fields error: %s", err.Error())
	}

	// fd, err := desc.Explain(g.fn, 3, true, whereParamsMap)
	// // rst, err := desc.ParseResult(g.fn, 2)
	// if err != nil {
	// 	return nil, g.fn.Error(err.Error())
	// }
	if sql.TableName == "" && fd.BeanObj == nil {
		return nil, g.fn.CreateError("no table name in sql and no bean struct in funcion")
	}
	var bean *xtype.XType
	if !hasSetFieldsInSQL {
		if fd.BeanObj == nil {
			return nil, g.fn.CreateError("no table name in sql and no bean struct in funcion")
		}
		bean, err = fd.BeanObj.GetBeanType()
		if err != nil {
			return nil, g.fn.CreateError("get bean type error: %s", err.Error())
		}
		// bean = fd.BeanObj.Type.Bean
	}

	if sql.TableName == "" {
		sql.TableName = bean.TableName
	}
	// 检查返回值
	if err := g.CheckUpdateReturns(fd); err != nil {
		return nil, g.fn.CreateError(err.Error())
	}
	// 完善条件信息
	if hasWhereInSQL {
		if _, err := g.prepareParams(fd, sql.WhereParams, false); err != nil {
			return nil, g.fn.CreateError(err.Error())
		}
	}

	// 完善 update 的字段信息
	if hasSetFieldsInSQL {
		if _, err := g.prepareParams(fd, sql.Params, true); err != nil {
			return nil, g.fn.CreateError(err.Error())
		}
	} else {
		// SQL 语句中没有 set 片段，生成 set 片段
		for _, f := range bean.Fields {
			if f.Ignore {
				continue
			}
			if f.IsPrimaryKey {
				// 主键不进行 update 操作
				continue
			}

			varAliasName := ""
			isJSON := false
			isTime := false
			isBlob := false

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
			arg := g.fn.FindParam(f.GetArgNames())
			// arg, ok := g.fn.ParamsMap[f.ArgName]
			if fd.Input.Bean == nil {
				// 输入参数中，没有与 beanObject 一致的对象
				if arg == nil {
					// log.Warn(g.fn.Format("input arg [%s] may be ignored"))
					continue
				}

				// 输入参数中找到对应的字段
				sql.Params = append(sql.Params, &intf.SQLParam{
					VarName:  arg.Name,
					VarAlias: varAliasName,
					JSON:     isJSON,
					Time:     isTime,
					Blob:     isBlob,
					TimeProp: timeProp,
				})
			} else {
				// 输入参数中有与 beanObject 一致的对象
				// 使用 输入参数中有与 beanObject 一致的对象 的数据作为参数
				sql.Params = append(sql.Params, &intf.SQLParam{
					VarName:  fd.Input.Bean.Name + "." + f.VarName,
					VarAlias: varAliasName,
					JSON:     isJSON,
					Time:     isTime,
					Blob:     isBlob,
					TimeProp: timeProp,
				})
				// }
			}
			sql.Columns = append(sql.Columns, f.FieldName()+"=?")
		}
		if len(sql.Columns) <= 0 {
			return nil, g.fn.CreateError("no field to be update in this func")
		}
	}

	fd.SQL, fd.SQLParams, fd.SQLWhereParams = g.dialect.CreateUpdateSQL(sql) // sql.CreateUpdate()
	// fd.SQL = strconv.Quote(fd.SQL)
	// fd.SQL = fd.SQL[1 : len(fd.SQL)-1]
	// 解析 xsql
	return fd, nil
}
