package method

import (
	"strings"

	"github.com/seerx/gpa/engine/generator/defines"
	rdesc "github.com/seerx/gpa/engine/generator/repo-desc"
	"github.com/seerx/gpa/engine/generator/xtype"
	"github.com/seerx/gpa/engine/sql/dialect/intf"
)

type delete struct {
	BaseMethod
}

func (g *delete) Test(fn *defines.Func) bool {
	// 如果函数名称中有 By 字符串，则优先为 DeleteXBy 操作
	if strings.Index(fn.Name, "Delete") == 0 {
		// fn.SQL 的形式
		// 1 delete table where id=:arg3
		//      delete 后面不要加 from
		//		类似原生 SQL 语句，arg1 arg2 arg3 从函数的参数中获取，不需要提供表相关 struct 类型数据
		// 2 where id=:arg3
		//		需要提供表相关的 struct 结构，根据参数自行生成 set field=? 语句
		// 3 空，没有 where 条件
		g.BaseMethod.Test(fn)
		fn.Template = defines.DELETE
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

func parseDeleteSQL(sql string) (*intf.SQL, error) {
	terms, err := splitSQL(sql)
	if err != nil {
		return nil, err
	}
	deleteIndex := -1
	// setIndex := -1
	whereIndex := -1
	for n, t := range terms {
		term := strings.ToLower(t)
		switch term {
		case "where":
			if whereIndex == -1 {
				whereIndex = n
			}
		case "delete":
			if deleteIndex == -1 {
				deleteIndex = n
			}
		}
	}

	var s intf.SQL
	if whereIndex-deleteIndex == 2 {
		s.TableName = terms[deleteIndex+1]
	}

	// whereTerms := []string{}
	// var whereParams []*desc.SQLParam
	// if whereIndex >= 0 {

	// 	for n := whereIndex + 1; n < len(terms); n++ {
	// 		ps, err := FindParams(terms[n])
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		col := terms[n]
	// 		var termParams []*desc.SQLParam
	// 		for m := len(ps) - 1; m >= 0; m-- {
	// 			var fieldName string
	// 			for k := n; k >= 0; k-- {
	// 				term := terms[k]

	// 				_, keyPos := lastIndex(term, "=", "<>", "<", "<=", ">", ">=")
	// 				// eqPos := strings.LastIndex(term, "=")
	// 				if keyPos > 0 {
	// 					_, openParenPos := lastIndex(term, "(")
	// 					if openParenPos > 0 {
	// 						fieldName = term[openParenPos+1 : keyPos]
	// 					} else {
	// 						fieldName = term[:keyPos]
	// 					}
	// 					break
	// 				} else if keyPos == 0 {
	// 					// 从上一个 term 开始找
	// 					if k > 0 {
	// 						term := terms[k-1]
	// 						if ss, _ := lastIndex(term, "=", "<>", "<", "<=", ">", ">="); ss == term {
	// 							if k > 1 {
	// 								term = terms[k-2]
	// 							} else {
	// 								return nil, fmt.Errorf("cann't find field name of params %s", ps[m].Name)
	// 							}
	// 						}

	// 						_, openParenPos := lastIndex(term, "(")
	// 						if openParenPos > 0 {
	// 							fieldName = term[openParenPos+1 : keyPos]
	// 						} else {
	// 							fieldName = term[:keyPos]
	// 						}
	// 						break
	// 					} else {
	// 						return nil, fmt.Errorf("cann't find field name of params %s", ps[m].Name)
	// 					}
	// 				} else {
	// 					if k > 1 {
	// 						kw := strings.ToLower(terms[k-1])
	// 						if kw != "in" && kw != "like" {
	// 							return nil, fmt.Errorf("invalid key word of params %s", ps[m].Name)
	// 						}
	// 						fieldName = terms[k-2]
	// 					} else {
	// 						return nil, fmt.Errorf("cann't find field name of params %s", ps[m].Name)
	// 					}
	// 				}
	// 			}
	// 			if fieldName == "" {
	// 				return nil, fmt.Errorf("cann't find field name of params %s", ps[m].Name)
	// 			}

	// 			col = ReplaceParam(col, ps[m], "?")
	// 			termParams = append(termParams, &desc.SQLParam{
	// 				SQLParamName:      ps[m].Name,
	// 				SQLParamFieldName: fieldName,
	// 			})

	// 		}
	// 		for n := len(termParams) - 1; n >= 0; n-- {
	// 			whereParams = append(whereParams, termParams[n])
	// 		}
	// 		whereTerms = append(whereTerms, col)
	// 	}
	// }

	whereSQL, whereParams, err := ParseWhere(terms, whereIndex)
	if err != nil {
		return nil, err
	}
	s.Where, s.WhereParams = whereSQL, whereParams
	return &s, nil
}

// func (g *delete) fillParamProp(fd *desc.FuncDesc, params []*desc.SQLParam) error {
// 	// var params []*desc.SQLParam
// 	// 组织 where 参数
// 	for _, p := range params {
// 		// name = utils.LowerFirstChar(name)
// 		found := false
// 		if p.IsInOperator {
// 			fd.DBUtilPackage = g.fn.AddDBUtilPackage()
// 		}
// 		for _, arg := range g.fn.Args {
// 			if arg.Type.IsContext() || arg.Type.IsStruct() || arg.Type.IsFunc() {
// 				// 以上三种类型不能作为 where 参数
// 				continue
// 			}

// 			if arg.Name == p.SQLParamName {
// 				// 找到输入参数
// 				p.VarName = arg.Name
// 				found = true
// 			}
// 		}
// 		if !found {
// 			if fd.Input.Bean != nil {
// 				for _, f := range fd.Input.Bean.Type.Bean.Fields {
// 					if f.GoVarName == utils.UpperFirstChar(p.SQLParamName) && !f.Type.IsJson() {
// 						// 找到输入参数
// 						p.VarName = fd.Input.Bean.Name + "." + f.GoVarName
// 						if f.SQLType.IsTime() {
// 							fd.DBUtilPackage = g.fn.AddDBUtilPackage()
// 							var timeProp = &dbutil.TimePropDesc{
// 								TypeName: f.SQLType.Name,
// 								Nullable: f.Nullable,
// 							}
// 							if f.TimeZone != nil {
// 								timeProp.TimeZone = f.TimeZone.String()
// 							}
// 							p.VarAlias = fd.NextVarName()
// 							p.Time = true
// 							p.TimeProp = timeProp
// 						}
// 						found = true
// 						break
// 					}
// 				}
// 			}
// 		}
// 		if !found {
// 			return fmt.Errorf("no where param [%s] found in func args", p.SQLParamName)
// 			// return
// 		}
// 	}
// 	return nil
// }

func (g *delete) Parse() (*rdesc.FuncDesc, error) {
	xsql := g.fn.SQL
	sql, err := parseDeleteSQL(xsql)
	if err != nil {
		return nil, g.fn.CreateError(err.Error())
	}
	// hasSetFieldsInSQL := sql.Columns != nil
	hasWhereInSQL := sql.WhereParams != nil

	// var whereParamsMap map[string]bool
	// if hasWhereInSQL {
	// 	whereParamsMap = map[string]bool{}
	// 	for _, p := range whereParams {
	// 		whereParamsMap[p.SQLParamFieldName] = true
	// 	}
	// }
	fd := rdesc.NewFuncDesc(g.fn, 3, g.logger)
	if err := fd.Explain(); err != nil {
		return nil, g.fn.CreateError(err.Error())
	}
	// fd, err := desc.Explain(g.fn, 3, false, nil)
	// // rst, err := desc.ParseResult(g.fn, 2)
	// if err != nil {
	// 	return nil, g.fn.Error(err.Error())
	// }
	if sql.TableName == "" && fd.BeanObj == nil {
		return nil, g.fn.CreateError("no table name in sql and no bean struct in funcion")
	}
	var bean *xtype.XType
	if sql.TableName == "" {
		// bean = fd.BeanObj.Type.Bean
		bean, err = fd.BeanObj.GetBeanType()
		if err != nil {
			return nil, g.fn.CreateError(err.Error())
		}
		sql.TableName = bean.TableName
	}
	// 检查返回值
	if err := g.CheckDeleteReturns(fd); err != nil {
		return nil, g.fn.CreateError(err.Error())
	}
	// 完善条件信息
	if hasWhereInSQL {
		if _, err := g.prepareParams(fd, sql.WhereParams, false); err != nil {
			return nil, g.fn.CreateError(err.Error())
		}
		// if err := g.fillParamProp(&fd, sql.WhereParams); err != nil {
		// 	return nil, g.fn.Error(err.Error())
		// }
	}

	fd.SQL, fd.SQLWhereParams = g.dialect.CreateDeleteSQL(sql)
	// 解析 xsql
	return fd, nil
}
