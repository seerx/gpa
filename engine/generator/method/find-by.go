package method

import (
	"strings"

	"github.com/seerx/gpa/engine/generator/defines"
	rdesc "github.com/seerx/gpa/engine/generator/repo-desc"
	"github.com/seerx/gpa/engine/generator/sqlgenerator"
	"github.com/seerx/gpa/rt/dbutil"
)

type findby struct {
	BaseMethod
}

func (g *findby) Test(fn *defines.Func) bool {
	if strings.Index(fn.Name, "Find") == 0 {
		if strings.Index(fn.Name, "By") > 0 {
			g.BaseMethod.Test(fn)
			fn.Template = defines.FIND
			return true
		}
	}
	return false
}

func (g *findby) Parse() (*rdesc.FuncDesc, error) {
	byIdx := strings.Index(g.fn.Name, "By")
	if byIdx < 4 {
		return nil, g.fn.CreateError("invalid name of FindXXXBy func")
	}

	fd := rdesc.NewFuncDesc(g.fn, 2, g.logger)
	// fd, err := desc.Explain(g.fn, 2, false, nil)
	if err := fd.Explain(); err != nil {
		return nil, g.fn.CreateError(err.Error())
	}

	if err := g.CheckFindReturns(fd); err != nil {
		return nil, g.fn.CreateError(err.Error())
	}

	if fd.BeanObj == nil {
		return nil, g.fn.CreateError("no bean struct found in funcion")
	}

	whereInName := g.fn.Name[byIdx+2:]
	sqlWhere, whereParams, err := parseWhereFromFuncName(whereInName)
	if err != nil {
		return nil, g.fn.CreateError(err.Error())
	}

	if _, err := g.prepareParams(fd, whereParams, false); err != nil {
		return nil, g.fn.CreateError(err.Error())
	}

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
		sql.Columns = append(sql.Columns, f.FieldName())

		isJSON := false
		isTime := false
		isBlob := false
		varAliasName := ""
		goType := ""

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
				goType = "[]byte"
			}
		}

		var timeProp *dbutil.TimePropDesc
		if !isBlob {
			if f.IsJSON {
				varAliasName = fd.NextVarName()
				fd.DBUtilPackage = g.fn.AddDBUtilPackage()
				isJSON = true
				goType = "interface{}"
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
				goType = fd.DBUtilPackage + ".NullTime"
			}
		}

		fd.Fields = append(fd.Fields, &rdesc.BeanField{
			Name:     f.VarName,
			VarAlias: varAliasName,
			JSON:     isJSON,
			Time:     isTime,
			Blob:     isBlob,
			TimeProp: timeProp,
			VarType:  goType,
			Ptr:      f.Field.Type.IsPtr,
		})
		// isJSON := false
		// varAliasName := ""
		// if f.IsJSON {
		// 	varAliasName = fd.NextVarName()
		// 	fd.DBUtilPackage = g.fn.AddDBUtilPackage()
		// 	isJSON = true
		// }
		// sql.Params = append(sql.Params, &desc.SQLParam{
		// 	VarName:  arg.Name,
		// 	VarAlias: varAliasName,
		// 	JSON:     isJSON,
		// 	Time:     isTime,
		// 	TimeProp: timeProp,
		// })
	}

	fd.SQL, fd.SQLWhereParams = g.sqlg.Query(&sql)

	return fd, nil
}

// func (g *findby) where(fd *desc.FuncDesc, whereParams []*desc.SQLParam) error {
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
// 			if fd.ParamBeanObj != nil {
// 				for _, f := range fd.ParamBeanObj.Type.Bean.Fields {
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
