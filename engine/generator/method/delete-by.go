package method

import (
	"strings"

	"github.com/seerx/gpa/engine/generator/defines"
	rdesc "github.com/seerx/gpa/engine/generator/repo-desc"
	"github.com/seerx/gpa/engine/sql/dialect/intf"
)

type deleteby struct {
	BaseMethod
}

func (g *deleteby) Test(fn *defines.Func) bool {
	if strings.Index(fn.Name, "Delete") == 0 {
		if strings.Index(fn.Name, "By") > 0 {
			g.BaseMethod.Test(fn)
			fn.Template = defines.DELETE
			return true
		}
	}
	return false
}

// func (g *deleteby) where(fd *rdesc.FuncDesc, whereParams []*intf.SQLParam) error {
// 	// var params []*desc.SQLParam
// 	// var fieldMap = map[string]bool{}
// 	// 组织 where 参数
// 	for _, p := range whereParams {
// 		// name = utils.LowerFirstChar(name)
// 		// fieldMap[p.SQLParamFieldName] = true
// 		if p.IsInOperator {
// 			fd.DBUtilPackage = g.fn.AddDBUtilPackage()
// 		}
// 		found := false
// 		for _, arg := range g.fn.Params {
// 			if arg.Type.IsContext() || arg.Type.IsStruct() || arg.IsFunc() {
// 				// 以上三种类型不能作为 where 参数
// 				continue
// 			}

// 			if arg.Name == utils.LowerFirstChar(p.SQLParamName) {
// 				// 找到输入参数
// 				p.VarName = arg.Name
// 				if p.IsInOperator {
// 					// arg.IsSlice
// 					if !arg.Slice {
// 						return fmt.Errorf("arg [%s] shuld be array", p.VarName)
// 					}
// 				}
// 				// params = append(params, &desc.SQLParam{
// 				// 	VarName: arg.Name,
// 				// })
// 				found = true
// 			}
// 		}
// 		if !found {
// 			if fd.Input.Bean != nil {
// 				for _, f := range fd.Input.Bean.Type.Bean.Fields {
// 					if f.GoVarName == p.SQLParamName && !f.Type.IsJson() {
// 						// 找到输入参数
// 						p.VarName = fd.Input.Bean.Name + "." + f.GoVarName
// 						if p.IsInOperator {
// 							// 判断是否是数组类型，如果不是数组则报出类型错误
// 							if !f.Type.IsSlice {
// 								return fmt.Errorf("arg [%s] shuld be array", p.VarName)
// 							}
// 							// if !f.Type {
// 							// 	return fmt.Errorf("arg [%s] shuld be array", p.VarName)
// 							// }
// 						}
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
// 							// params = append(params, &desc.SQLParam{
// 							// 	VarName:  fd.Input.Bean.Name + "." + f.GoVarName,
// 							// 	VarAlias: fd.NextVarName(),
// 							// 	Time:     true,
// 							// 	TimeProp: timeProp,
// 							// })
// 						}
// 						// else {
// 						// 	params = append(params, &desc.SQLParam{
// 						// 		VarName: fd.Input.Bean.Name + "." + f.GoVarName,
// 						// 	})
// 						// }
// 						found = true
// 					}
// 				}
// 			}
// 		}
// 		if !found {
// 			return fmt.Errorf("no where param [%s] found in func args", p.VarName)
// 			// return
// 		}
// 	}
// 	return nil
// }

func (g *deleteby) Parse() (*rdesc.FuncDesc, error) {
	byIdx := strings.Index(g.fn.Name, "By")
	if byIdx < 6 {
		return nil, g.fn.CreateError("invalid name of DeleteXXXBy func")
	}

	fd := rdesc.NewFuncDesc(g.fn, 3, g.logger)
	if err := fd.Explain(); err != nil {
		return nil, g.fn.CreateError(err.Error())
	}
	// fd, err := desc.Explain(g.fn, 3, false, nil)
	// // rst, err := desc.ParseResult(g.fn, 2)
	// if err != nil {
	// 	return nil, g.fn.CreateError(err.Error())
	// }
	if fd.BeanObj == nil {
		return nil, g.fn.CreateError("no struct bean found in funcion")
	}
	if err := g.CheckDeleteReturns(fd); err != nil {
		return nil, g.fn.CreateError(err.Error())
	}

	whereInName := g.fn.Name[byIdx+2:]
	sqlWhere, whereParams, err := parseWhereFromFuncName(whereInName)
	if err != nil {
		return nil, g.fn.CreateError(err.Error())
	}

	if _, err = g.prepareParams(fd, whereParams, false); err != nil {
		return nil, g.fn.CreateError(err.Error())
	}

	// sql := strings.Builder{}
	// sqlParams := []*desc.SQLParam{}
	// columns := []string{}
	bean, err := fd.BeanObj.GetBeanType()
	if err != nil {
		return nil, g.fn.CreateError(err.Error())
	}
	var sql = intf.SQL{
		TableName:   bean.TableName,
		Where:       sqlWhere,
		WhereParams: whereParams,
	}

	fd.SQL, fd.SQLWhereParams = g.dialect.CreateDeleteSQL(&sql) // sql.CreateDelete() //   append(sqlParams, whereParams...)

	return fd, nil
}
