package method

import (
	"strings"

	"github.com/seerx/gpa/engine/generator/defines"
	rdesc "github.com/seerx/gpa/engine/generator/repo-desc"
	"github.com/seerx/gpa/engine/sql/dialect/intf"
)

type countby struct {
	BaseMethod
}

func (g *countby) Test(fn *defines.Func) bool {
	if strings.Index(fn.Name, "Count") == 0 {
		if strings.Index(fn.Name, "By") > 0 {
			g.BaseMethod.Test(fn)
			fn.Template = defines.COUNT
			return true
		}
	}
	return false
}

// func (g *countby) where(fd *desc.FuncDesc, whereParams []*desc.SQLParam) error {
// 	// var params []*desc.SQLParam
// 	// var fieldMap = map[string]bool{}
// 	// 组织 where 参数
// 	for _, p := range whereParams {
// 		// name = utils.LowerFirstChar(name)
// 		// fieldMap[p.FieldName] = true
// 		if p.IsInOperator {
// 			fd.DBUtilPackage = g.fn.AddDBUtilPackage()
// 		}
// 		found := false
// 		for _, arg := range g.fn.Args {
// 			if arg.Type.IsContext() || arg.Type.IsStruct() || arg.Type.IsFunc() {
// 				// 以上三种类型不能作为 where 参数
// 				continue
// 			}

// 			if arg.Name == utils.LowerFirstChar(p.SQLParamName) {
// 				// 找到输入参数
// 				p.VarName = arg.Name
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
// 						//  else {
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

func (g *countby) Parse() (*rdesc.FuncDesc, error) {
	byIdx := strings.Index(g.fn.Name, "By")
	if byIdx < 5 {
		return nil, g.fn.CreateError("invalid name of CountXXXBy func")
	}

	fd := rdesc.NewCountFuncDesc(g.fn, 3, g.logger)
	if err := fd.Explain(); err != nil {
		return nil, g.fn.CreateError(err.Error())
	}

	// fd, err := desc.ExplainCount(g.fn, 3, false, nil)
	// // rst, err := desc.ParseResult(g.fn, 2)
	// if err != nil {
	// 	return nil, g.fn.Error(err.Error())
	// }
	if fd.BeanObj == nil {
		return nil, g.fn.CreateError("no struct bean found in funcion")
	}
	if err := g.CheckCountReturns(fd); err != nil {
		return nil, g.fn.CreateError(err.Error())
	}

	whereInName := g.fn.Name[byIdx+2:]
	sqlWhere, whereParams, err := parseWhereFromFuncName(whereInName)
	if err != nil {
		return nil, g.fn.CreateError(err.Error())
	}

	// whereParams, _,
	if _, err := g.prepareParams(fd, whereParams, false); err != nil {
		return nil, g.fn.CreateError(err.Error())
	}
	// if err := g.where(&fd, whereParams); err != nil {
	// 	return nil, g.fn.Error(err.Error())
	// }

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
		Columns:     []string{"count(0)"},
	}

	fd.SQL, fd.SQLWhereParams = g.dialect.CreateQuerySQL(&sql) // sql.CreateDelete() //   append(sqlParams, whereParams...)

	return fd, nil
}
