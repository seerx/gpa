package rdesc

import (
	"fmt"

	"github.com/seerx/gpa/engine/generator/defines"
	"github.com/seerx/gpa/engine/sql/dialect/intf"
	"github.com/seerx/gpa/logger"
	"github.com/seerx/gpa/rt/dbutil"
)

type VarSetPair struct {
	VarName   string
	ValueName string
}

type BeanField struct {
	Name     string
	VarAlias string
	VarType  string
	JSON     bool
	Time     bool
	Blob     bool
	SQLType  string
	Ptr      bool
	TimeProp *dbutil.TimePropDesc
}

type FuncDesc struct {
	varNameIndex int
	// SQLReturnVarName string

	FuncName string
	// JSONPackage     string
	SQLPackage     string
	RumTimePackage string
	DBUtilPackage  string
	// TimePackage     string
	// ContextPackage  string

	Result *Result
	Input  *Input

	// 表关联结构信息
	BeanObj           *defines.Object // 函数中关联表的 struct 类型参数
	ParamBeanObj      *defines.Object // 函数中关联查询参数的 struct 类型参数，用在 Find 函数中
	BeanTypeName      string          // 函数中关联表的 struct 类型名称
	BeanVarNeedCreate bool            // 是否需要创建 BeanType 的变量,当函数的参数中没有对应的 struct 参数时，需要创建
	BeanVarName       string          // 在函数体中使用的 关联表的 struct 类型变量的名称

	BeanFieldSetValues []VarSetPair // 需要给 VarBeanName 设置的变量

	SQLReturnVarName string // 接收执行 SQL 时返回值的变量名称

	// 与select 语句中出现的 column 一一对应
	Fields []*BeanField
	// SQL 相关
	SQL            string
	SQLParams      []*intf.SQLParam
	SQLWhereParams []*intf.SQLParam

	fn         *defines.Func
	logger     logger.GpaLogger
	maxResults int
	countFunc  bool
}

func NewFuncDesc(fn *defines.Func, maxResults int, logger logger.GpaLogger) *FuncDesc {
	return newFuncDesc(fn, maxResults, false, logger)
}

func NewCountFuncDesc(fn *defines.Func, maxResults int, logger logger.GpaLogger) *FuncDesc {
	return newFuncDesc(fn, maxResults, true, logger)
}

func newFuncDesc(fn *defines.Func, maxResults int, countFunc bool, logger logger.GpaLogger) *FuncDesc {
	return &FuncDesc{
		fn:         fn,
		logger:     logger,
		maxResults: maxResults,
		countFunc:  countFunc,
	}
}

func (fd *FuncDesc) NextVarName() string {
	fd.varNameIndex++
	return fmt.Sprintf("var%d", fd.varNameIndex)
}

type argWrap struct {
	Arg      *defines.Object
	SetValue bool
	InWhere  bool
}

func (fd *FuncDesc) Explain() (err error) {
	fd.Result, err = explainResult(fd.fn, fd, fd.maxResults, fd.countFunc)
	if err != nil {
		return
	}
	fd.Input, fd.BeanObj, err = explainInput(fd.fn, fd.Result, fd.logger)
	if err != nil {
		return
	}
	if fd.Result.AffectVar != "" {
		fd.SQLPackage = fd.fn.AddSQLPackage()
	}
	fd.FuncName = fd.fn.Name

	if fd.BeanObj != nil {
		fd.BeanTypeName = fd.BeanObj.Type.String()
		if fd.Input.Bean == nil {
			// 输入参数中，没有 bean struct
			fd.BeanVarName = fd.NextVarName() // fmt.Sprintf("%s%d", utils.LowerFirstChar(fd.BeanObj.Type.Bean.Name), time.Now().Unix())
			// fd.BeanVarNeedCreate = useSetValues
		} else {
			fd.BeanVarName = fd.Input.Bean.Name
		}
	}

	return
}

func (fd *FuncDesc) ExplainSetBeanFieldsValueWithArgs(whereParams map[string]bool) (err error) {
	if fd.BeanObj != nil {
		// fd.BeanTypeName = fd.BeanObj.Type.String()
		// if fd.Input.Bean == nil {
		// 	// 输入参数中，没有 bean struct
		// 	fd.BeanVarName = fd.NextVarName() // fmt.Sprintf("%s%d", utils.LowerFirstChar(fd.BeanObj.Type.Bean.Name), time.Now().Unix())
		// 	fd.BeanVarNeedCreate = useSetValues
		// } else {
		// 	fd.BeanVarName = fd.Input.Bean.Name
		// }

		// if useSetValues {
		fd.BeanVarNeedCreate = fd.Input.Bean == nil
		// insert 和 update 中涉及的函数的参数赋值给表关联的 struct 的 field 的列表
		inputArgs := map[string]*argWrap{}
		for _, inputArg := range fd.fn.Params {
			inputArgs[inputArg.Name] = &argWrap{
				Arg:      defines.NewObject(fd.fn.GetRepoInterface(), inputArg),
				SetValue: false,
			}
		}

		bean, err := fd.BeanObj.GetBeanType()
		if err != nil {
			return err
		}
		for _, f := range bean.Fields {
			if f.Ignore {
				continue
			}
			names := f.GetArgNames()
			var arg *argWrap
			var ok bool
			for _, name := range names {
				arg, ok = inputArgs[name] // f.ArgName
				if ok {
					break
				}
			}
			// arg, ok := inputArgs[f.ArgName]
			if whereParams != nil {
				if _, exsists := whereParams[f.FieldName()]; exsists {
					// 在 sql 语句 where 条件中已经有该字段，则该字段不需要赋值
					if ok {
						arg.InWhere = true
					}
					continue
				}
			}

			if fd.Input.Bean == nil {
				// 输入参数中，没有与 beanObject 一致的对象
				if !ok {
					// log.Warn(g.fn.Format("input arg [%s] may be ignored"))
					continue
				}
				valName := arg.Arg.Name

				if f.Field.Type.IsPtr != arg.Arg.Type.IsPtr {
					if f.Field.Type.IsPtr {
						valName = "&" + valName
					} else {
						valName = "*" + valName
					}
				}
				fd.BeanFieldSetValues = append(fd.BeanFieldSetValues, VarSetPair{
					VarName:   f.VarName,
					ValueName: valName,
				})
				arg.SetValue = true
			} else {
				// 输入参数中有与 beanObject 一致的对象
				if ok {
					valName := arg.Arg.Name
					if f.Field.Type.IsPtr != arg.Arg.Type.IsPtr {
						if f.Field.Type.IsPtr {
							valName = "&" + valName
						} else {
							valName = "*" + valName
						}
					}
					// 输入参数中找到对应的字段
					fd.BeanFieldSetValues = append(fd.BeanFieldSetValues, VarSetPair{
						VarName:   f.VarName,
						ValueName: valName,
					})
					arg.SetValue = true
					// sqlParams = append(sqlParams, arg.Name)
				}
			}
		}

		for _, arg := range inputArgs {
			if !arg.SetValue && !arg.InWhere &&
				(fd.Input.Bean != nil && fd.Input.Bean.Name != arg.Arg.Name) &&
				!arg.Arg.Type.IsContext() {
				fd.logger.Warn(fd.fn.Format("input arg [%s] is ignored", arg.Arg.Name))
			}
		}
		// }
	}

	return
}
