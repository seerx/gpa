package method

import (
	"fmt"
	"strings"

	"github.com/seerx/gpa/engine/generator/defines"
	rdesc "github.com/seerx/gpa/engine/generator/repo-desc"
	"github.com/seerx/gpa/engine/generator/xtype"
	"github.com/seerx/gpa/engine/sql/dialect/intf"
	"github.com/seerx/gpa/engine/sql/names"
	"github.com/seerx/gpa/rt/dbutil"
	"github.com/seerx/gpa/utils"
)

type find struct {
	BaseMethod
}

func (g *find) Test(fn *defines.Func) bool {
	// 如果函数名称中有 By 字符串，则优先为 FindXBy 操作
	if strings.Index(fn.Name, "Find") == 0 {
		// fn.SQL 的形式
		// 1 select {fields} from {table} where id=:arg3
		//		类似原生 SQL 语句，arg1 arg2 arg3 从函数的参数中获取，不需要提供表相关 struct 类型数据
		// 2 from {table} where id=:arg3
		// 		需要提供表相关的 struct 结构，表名由 sql 语句中的 table 定义，根据参数自行生成 要查询的列列表
		// 3 where id=:arg3
		//		需要提供表相关的 struct 结构，根据参数自行生成 要查询的列列表
		// 4 空，没有 where 条件
		g.BaseMethod.Test(fn)
		fn.Template = defines.FIND
		return true
		// }
	}

	return false
}

func (g *find) Parse() (*rdesc.FuncDesc, error) {
	xsql := g.fn.SQL
	sql, err := parseFindSQL(xsql)
	if err != nil {
		return nil, g.fn.CreateError(err.Error())
	}

	fd := rdesc.NewFuncDesc(g.fn, 2, g.logger)
	// fd, err := desc.Explain(g.fn, 2, false, nil)
	if err := fd.Explain(); err != nil {
		return nil, g.fn.CreateError(err.Error())
	}
	if fd.BeanObj == nil {
		return nil, g.fn.CreateError("no bean struct found in funcion")
	}
	// if (len(sql.SelectFields) == 0 || sql.TableName == "") && fd.BeanObj == nil {
	// 	return nil, g.fn.Error("no table name in sql and no bean struct in funcion")
	// }

	if err := g.CheckFindReturns(fd); err != nil {
		return nil, g.fn.CreateError(err.Error())
	}

	if len(sql.WhereParams) > 0 {
		if _, err := g.prepareParams(fd, sql.WhereParams, false); err != nil {
			return nil, g.fn.CreateError(err.Error())
		}
	}

	// var bean *beans.Bean
	bean, err := fd.BeanObj.GetBeanType()
	if err != nil {
		return nil, g.fn.CreateError(err.Error())
	}
	if sql.TableName == "" {
		sql.TableName = bean.TableName
	}
	if len(sql.SelectFields) == 0 {
		// SQL 语句中没有 select 字段
		for _, f := range bean.Fields {
			if f.Ignore {
				continue
			}
			sql.Columns = append(sql.Columns, f.FieldName())

			isJSON := false
			isTime := false
			isBlob := false
			var timeProp *dbutil.TimePropDesc
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
				TimeProp: timeProp,
				VarType:  goType,
				Ptr:      f.Field.Type.IsPtr,
			})
		}
	} else {
		// sql 语句中 select 字段
		for _, field := range sql.SelectFields {
			tableField := names.ToTableName(field) //  fieldMapper.Obj2Table(field)
			var f *xtype.Field = nil
			// found := false
			for _, fd := range bean.Fields {
				if fd.FieldName() == tableField {
					f = fd
					break
				}
			}
			if f == nil {
				return nil, g.fn.CreateError("field [%s] not found in bean struct", field)
			}

			isJSON := false
			isTime := false
			isBlob := false
			var timeProp *dbutil.TimePropDesc
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
		}
	}

	fd.SQL, fd.SQLWhereParams = g.dialect.CreateQuerySQL(sql)
	return fd, nil
}

// func (g *find) fillParamProp(fd *desc.FuncDesc, params []*desc.SQLParam) error {
// 	// var params []*desc.SQLParam
// 	// 组织 where 参数
// 	for _, p := range params {
// 		// name = utils.LowerFirstChar(name)
// 		if p.IsInOperator {
// 			fd.DBUtilPackage = g.fn.AddDBUtilPackage()
// 		}
// 		found := false
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

func parseFindSQL(sql string) (*intf.SQL, error) {
	terms, err := splitSQL(sql)
	if err != nil {
		return nil, err
	}
	selectIndex := -1
	fromIndex := -1
	whereIndex := -1
	for n, t := range terms {
		term := strings.ToLower(t)
		switch term {
		case "where":
			if whereIndex == -1 {
				whereIndex = n
			}
		case "from":
			if fromIndex == -1 {
				fromIndex = n
			}
		case "select":
			if selectIndex == -1 {
				selectIndex = n
			}
		}
	}

	var s intf.SQL
	if whereIndex-fromIndex >= 2 && fromIndex >= 0 {
		s.TableName = strings.Join(terms[fromIndex+1:whereIndex-1], " ")
	}

	// 解析 select 子句
	if selectIndex == 0 {
		selectEnd := len(terms)
		if fromIndex > 0 {
			selectEnd = fromIndex
		} else if whereIndex > 0 {
			selectEnd = whereIndex
		}

		starSelect := false
		if selectEnd-selectIndex == 2 {
			if terms[selectIndex+1] == "*" {
				// select * from 或者 select * where
				starSelect = true
			}
		}
		if !starSelect {
			selectFields := strings.Join(terms[selectIndex+1:selectEnd], " ")
			ary := strings.Split(selectFields, ",")
			for _, f := range ary {
				f = strings.TrimSpace(f)
				items := strings.Split(f, " ")
				fieldName := items[len(items)-1]

				if !utils.IsValidSQLFieldName(fieldName) {
					return nil, fmt.Errorf("field name [%s] is invalid", fieldName)
				}
				s.SelectFields = append(s.SelectFields, fieldName)
			}
			s.Columns = []string{selectFields}
		}

	}

	// 解析 where 子句
	whereSQL, whereParams, err := ParseWhere(terms, whereIndex)
	if err != nil {
		return nil, err
	}
	s.Where, s.WhereParams = whereSQL, whereParams
	return &s, nil
}
