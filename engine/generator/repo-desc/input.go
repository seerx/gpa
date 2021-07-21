package rdesc

import (
	"fmt"
	"strings"

	"github.com/seerx/gpa/engine/generator/defines"
	"github.com/seerx/gpa/engine/generator/xtype"
	"github.com/seerx/gpa/engine/objs"
	"github.com/seerx/gpa/logger"
)

type ArgPair struct {
	Name         string
	Type         string
	Arg          *objs.Object
	IsReturnFunc bool
	IsMapKeyFunc bool
}

type Input struct {
	Args           []*ArgPair      // 输入参数列表
	Bean           *defines.Object // 与表关联的 struct 参数
	ContextArgName string          // context.Context 参数名称

	Callback             *objs.Object // 返回数据的回调函数
	CallbackArgIsPtr     bool         // Callback 函数的参数是否指针
	KeyGenerator         *objs.Object // 生成 map 键值的回调函数
	KeyType              string       // 主键类型
	KeyGeneratorArgIsPtr bool         // KeyGenerator 函数的参数是否指针
}

// 输入参数中出现与表关联的 struct 时，只识别第一个参数作为 table 的 bean
func explainInput(fn *defines.Func, rst *Result, log logger.GpaLogger) (input *Input, beanObj *defines.Object, err error) {
	input = &Input{}
	beanObj = rst.Bean
	// 查找所有结构体类型的参数
	structArgs := []*objs.Object{}
	for _, arg := range fn.Params {
		// 组织参数列表
		if arg.IsFunc {
			// 参数是函数
			fnArgs := []string{}
			for _, a := range arg.Params {
				paramType := a.Type.StringExt()
				if a.Name == "" {
					if a.IsMap {
						paramType = fmt.Sprintf("map[%s]%s", a.Key.StringExt(), a.Type.StringExt())
					}
					if a.IsSlice {
						paramType = fmt.Sprintf("[]%s", a.Type.StringExt())
					}
				}
				if a.Name == "" {
					fnArgs = append(fnArgs, paramType)
				} else {
					fnArgs = append(fnArgs, a.Name+" "+paramType)
				}
			}
			fnRes := []string{}
			for _, a := range arg.Results {
				fnRes = append(fnRes, a.Type.StringExt())
			}
			resExpr := strings.Join(fnRes, ", ")
			if len(fnRes) > 0 {
				resExpr = fmt.Sprintf("(%s)", resExpr)
			}
			input.Args = append(input.Args, &ArgPair{
				Name: arg.Name,
				Type: fmt.Sprintf("func(%s) %s", strings.Join(fnArgs, ", "), resExpr),
				Arg:  arg,
			})
			continue
		}

		// argType := arg.Type.String()
		argType := arg.Type.StringExt()
		if arg.IsMap {
			argType = fmt.Sprintf("map[%s]%s", arg.Key.StringExt(), arg.Type.StringExt())
		}
		if arg.IsSlice {
			argType = fmt.Sprintf("[]%s", arg.Type.StringExt())
		}
		input.Args = append(input.Args, &ArgPair{
			Name: arg.Name,
			Type: argType,
			Arg:  arg,
		})

		if arg.Type.IsContext() {
			// 找到 context.Context 参数
			input.ContextArgName = arg.Name
			continue
		}
		if beanObj != nil {
			if beanObj.Type.Equals(&arg.Type) {
				if input.Bean == nil {
					// 找到 bean 的输入参数
					input.Bean = defines.NewObject(fn.GetRepoInterface(), arg) //  arg
				} else {
					log.Warnf("input arg [%s] is ignored", arg.Name)
				}
			}
		} else {
			if arg.Type.IsCustom() {
				structArgs = append(structArgs, arg)
			}
		}
	}

	if beanObj == nil {
		if len(structArgs) > 1 {
			// 结构体类型参数多于 1 个, 确定哪一个是
			// 确定方案：在函数的输入参数中找到与当前参数的 struct.field（不是 ignore）类型相同的参数，则认为当前 参数为 bean
		root:
			for n, arg := range structArgs {
				if arg.Type.IsCustom() {
					obj := defines.NewObject(fn.GetRepoInterface(), arg)
					var bean *xtype.XType
					bean, err = obj.GetBeanType()
					if err != nil {
						return
					}
					if bean == nil {
						// err = fmt.Errorf("no bean find in ")
						return
					}
					for _, f := range bean.Fields {
						for m, item := range fn.Params {
							if m != n {
								// f.Column.Object.Type.Equals(item.Type)
								// f.Type.Equals(item.Type)
								if !f.Ignore && f.Column.Field.Type.Equals(&item.Type) {
									// 在参数列表中找到与 arg 结构中的 f 字段 相同类型的参数
									// 此时认定 arg 为 beanObj
									beanObj = obj
									input.Bean = obj
									break root
								}
							}
						}
					}
				}
			}
		} else if len(structArgs) == 1 {
			// 只有一个结构体类型的参数
			input.Bean = defines.NewObject(fn.GetRepoInterface(), structArgs[0])
			beanObj = input.Bean
		}
	}

	// if beanObj == nil {
	// 	// 没有找到 bean 对象，函数定义有问题
	// 	err = fn.Error("no struct bean found in funcion")
	// 	return // nil, nil, "", g.fn.Error("no struct bean found in funcion")
	// }
	return
}
