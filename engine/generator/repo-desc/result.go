package rdesc

import (
	"fmt"

	"github.com/seerx/gpa/engine/generator/defines"
)

type Result struct {
	Count     int
	Bean      *defines.Object // 与表关联的 struct 参数
	AffectVar string          // 返回影响行数的变量名称，适用于 update 和 delete 等操作
	// ReturnAffect bool // 是否返回影响函数，适用于 update 和 delete 等操作
	List           []*defines.Object
	FindOne        bool
	ReturnTypeName string
	CountVar       string // 返回结果数量的变量名称 适用于 select count(),即 Count 操作
	// Third  *metas.Object
}

// explainResult 解析函数返回值
func explainResult(fn *defines.Func, fd *FuncDesc, maxResults int, countFunc bool) (*Result, error) {
	rst := &Result{
		Count: len(fn.Results),
	}
	if rst.Count < 1 {
		return nil, fn.CreateError("need 1 return at leat")
	}
	if rst.Count > maxResults {
		return nil, fn.CreateError("%d returns at most", maxResults)
	}
	// 最后的返回值必须是 error
	lastResult := fn.Results[rst.Count-1]
	// 只有一个返回值，返回值必须是 error
	if !lastResult.Type.IsError() {
		return nil, fn.CreateError("the only return value must be error")
	}

	for n, r := range fn.Results {
		obj := defines.NewObject(fn.GetRepoInterface(), r)
		// obj.Object = r
		rst.List = append(rst.List, obj)
		if n >= rst.Count-1 {
			break
		}
		if r.Type.IsCustom() {
			rst.Bean = obj
			rst.ReturnTypeName = r.Type.StringExt()
			if r.IsMap {
				rst.ReturnTypeName = fmt.Sprintf("map[%s]%s", r.Key.StringExt(), r.Type.StringExt())
			} else if r.IsSlice {
				rst.ReturnTypeName = fmt.Sprintf("[]%s", r.Type.StringExt())
			}
		}

		if r.Type.IsInt64() {
			// sql 操作返回的变量名称
			fd.SQLReturnVarName = fd.NextVarName()
			if countFunc {
				// select count(0) from ...
				rst.CountVar = fd.NextVarName()
			} else {
				// 修改和删除操作，返回影响行数
				rst.AffectVar = fd.NextVarName()
			}
		}
	}
	return rst, nil
}
