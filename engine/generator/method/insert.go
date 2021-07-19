package method

import (
	"strings"

	"github.com/seerx/gpa/engine/generator/defines"
	rdesc "github.com/seerx/gpa/engine/generator/repo-desc"
	"github.com/seerx/gpa/engine/generator/sqlgenerator"
	"github.com/seerx/gpa/rt/dbutil"
)

type insert struct {
	BaseMethod
}

func (g *insert) Test(fn *defines.Func) bool {
	if strings.Index(fn.Name, "Insert") == 0 {
		g.BaseMethod.Test(fn)
		fn.Template = defines.INSERT
		return true
	}
	return false
}

func (g *insert) Parse() (*rdesc.FuncDesc, error) {
	fd := rdesc.NewFuncDesc(g.fn, 2, g.logger)
	if err := fd.Explain(); err != nil {
		err = g.fn.CreateError(err.Error())
		g.logger.Error(err, "explain error")
		return nil, err
	}
	if err := fd.CheckAutoincrPrimaryKey(); err != nil {
		return nil, g.fn.CreateError("check primary key error: %s", err.Error())
	}
	if err := fd.ExplainSetBeanFieldsValueWithArgs(nil); err != nil {
		g.fn.CreateError(err.Error())
		g.logger.Error(err, "explain set bean fields error")
		return nil, err
	}
	// fd, err := desc.Explain(g.fn, 2, true, nil)
	// if err != nil {
	// 	return nil, g.fn.Error(err.Error())
	// }
	if fd.BeanObj == nil {
		err := g.fn.CreateError("no struct bean found in funcion")
		g.logger.Error(err, "explain set bean fields error")
		return nil, err
	}

	// inputArgs := map[string]*metas.Object{}
	// for _, inputArg := range g.fn.Args {
	// 	inputArgs[inputArg.Name] = inputArg
	// }
	bean, err := fd.BeanObj.GetBeanType()
	if err != nil {
		return nil, err
	}
	// bean := fd.BeanObj.Type.Bean
	// sql := strings.Builder{}
	// sqlParams := []*desc.SQLParam{}

	var sql = sqlgenerator.SQL{
		TableName:                bean.TableName,
		ReturnAutoincrPrimaryKey: fd.AutoinrPrimaryKeyField,
	}
	// sql.TableName = bean.TableName
	// _, err = sql.WriteString("INSERT INTO " + bean.TableName + " (")
	// if err != nil {
	// 	return nil, err
	// }

	// columns := []string{}
	// params := ""
	for _, f := range bean.Fields {
		if f.Ignore {
			continue
		}
		if f.Field.Name == fd.AutoinrPrimaryKeyField {
			// 自增主键，忽略
			continue
		}
		arg := g.fn.FindParam(f.GetArgNames())
		if fd.Input.Bean == nil { // 输入参数中没有与 beanObject 一致的对象
			if arg == nil { // 在函数的参数中也没有找到对应 f 名称的参数
				continue // 插入时忽略该字段
			}
		}

		isJSON := false
		isTime := false
		isBlob := false
		var timeProp *dbutil.TimePropDesc
		varAliasName := ""

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

		// arg, ok := g.fn.ArgMap[f.ArgName]
		if fd.Input.Bean == nil {
			// 输入参数中，没有与 beanObject 一致的对象
			// if arg == nil {
			// 	// g.logger.Warnf(g.fn.Format("input param [%s] has bean ignored", f.Field.Name))
			// 	continue
			// }

			// 输入参数中找到对应的字段
			sql.Params = append(sql.Params, &sqlgenerator.SQLParam{
				VarName:  arg.Name,
				VarAlias: varAliasName,
				JSON:     isJSON,
				Time:     isTime,
				Blob:     isBlob,
				TimeProp: timeProp,
			})
			// sqlParams = append(sqlParams, &desc.SQLParam{
			// 	VarName:  arg.Name,
			// 	VarAlias: varAliasName,
			// 	JSON:     isJSON,
			// 	Time:     isTime,
			// 	TimeProp: timeProp,
			// })
		} else {
			// 输入参数中有与 beanObject 一致的对象
			// 使用 输入参数中有与 beanObject 一致的对象 的数据作为参数
			sql.Params = append(sql.Params, &sqlgenerator.SQLParam{
				VarName:  fd.Input.Bean.Name + "." + f.VarName,
				VarAlias: varAliasName,
				JSON:     isJSON,
				Time:     isTime,
				Blob:     isBlob,
				TimeProp: timeProp,
			})
			// sqlParams = append(sqlParams, &desc.SQLParam{
			// 	VarName:  fd.Input.Bean.Name + "." + f.GoVarName,
			// 	VarAlias: varAliasName,
			// 	JSON:     isJSON,
			// 	Time:     isTime,
			// 	TimeProp: timeProp,
			// })
		}
		sql.Columns = append(sql.Columns, f.FieldName())
		sql.ParamPlaceHolder = append(sql.ParamPlaceHolder, "?")
		// columns = append(columns, f.FieldName)
		// params += ",?"
	}

	// _, err = sql.WriteString(strings.Join(columns, ",") + ") VALUES (" + params[1:] + ")")
	// if err != nil {
	// 	return nil, err
	// }

	fd.SQL, fd.SQLParams = g.sqlg.Insert(&sql) // sql.CreateInsert()
	// fd.SQLParams = s.Params // sqlParams

	return fd, nil
}
