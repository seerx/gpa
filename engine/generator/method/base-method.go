package method

import (
	"github.com/seerx/gpa/engine/generator/defines"
	rdesc "github.com/seerx/gpa/engine/generator/repo-desc"
	"github.com/seerx/gpa/engine/sql/dialect/intf"
	"github.com/seerx/gpa/logger"
)

type BaseMethod struct {
	dialect intf.Dialect
	fn      *defines.Func
	logger  logger.GpaLogger
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
						if len(arg.Arg.Results) == 1 && arg.Arg.Results[0].Type.EqualsExactly(fd.Result.List[0].Key) {
							arg.IsMapKeyFunc = true
							fd.Input.KeyGenerator = arg.Arg
							// arg.Arg.Args
							foundKeyFunc = true
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
				if len(arg.Arg.Params) == 1 && arg.Arg.Params[0].Type.IsStruct() {
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
		if arg.Arg.Type.IsStruct() {
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
			if !fd.Result.List[1].Type.IsStruct() {
				return bg.fn.CreateError("the second return must be bean struct")
			}
		}
	} else {
		// 不返回 Affect 行数
		// 可能是 (struct, error) 或者 error
		if fd.Result.Bean != nil {
			// 此时返回值应该是 (struct, error) 形式
			if !fd.Result.List[0].Type.IsStruct() {
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
			if !fd.Result.List[1].Type.IsStruct() {
				return bg.fn.CreateError("the second return must be bean struct")
			}
		}
	} else {
		// 不返回 Affect 行数
		// 可能是 (struct, error) 或者 error
		if fd.Result.Bean != nil {
			// 此时返回值应该是 (struct, error) 形式
			if !fd.Result.List[0].Type.IsStruct() {
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
		if !fd.Result.List[1].Type.IsStruct() {
			return bg.fn.CreateError("the first return must be bean struct")
		}
	}
	// fd.Result.CountVar = fd.NextVarName()
	return nil
}
