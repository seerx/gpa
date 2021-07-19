package method

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/seerx/gpa/engine/generator/sqlgenerator"
	"github.com/seerx/gpa/engine/sql/names"
	"github.com/seerx/gpa/utils"
)

func parseWhereFromFuncName(whereInName string) (string, []*sqlgenerator.SQLParam, error) {
	var inPlaceIndex = time.Now().Unix()
	// fnName := g.fn.Name
	// whereInName := fnName[len("UpdateBy"):]
	if len(whereInName) <= 0 {
		return "", nil, errors.New("func like UpdateByXxx needs condition field(s)")
	}
	words := []string{}
	var appendWord = func(word string, end bool) {
		words = append(words, word)
		if word == "And" || word == "Or" || end {
			// if (word == "And" || word == "Or") {
			// 	if len(words) > 0  {
			// 	}
			// }
			andOrSkip := 0
			if !end {
				andOrSkip = 1
			}
			var m int
			foundConnect := 0
			for m = len(words) - 1 - andOrSkip; m >= 0; m-- {
				wd := words[m]
				if wd == "And" || wd == "Or" {
					foundConnect++
				} else if foundConnect > 0 {
					m++
					break
				}
			}
			for k := m + 2; k < len(words)-andOrSkip; k++ {
				words[m+1] = words[m+1] + words[k]
				words[k] = ""
			}
		}
		// if !end {
		// 	words = append(words, word)
		// }
	}

	word := ""
	// bytes.ToLower()
	for n := range whereInName {
		ch := whereInName[n]
		if utils.IsCapital(ch) {
			if word != "" {
				appendWord(word, false)
			}
			word = string(ch)
		} else {
			word = word + string(ch)
		}
	}
	// if word != "And" && word != "Or" {
	appendWord(word, true)
	foundConnect := false
	for n := len(words) - 1; n >= 0; n-- {
		if words[n] == "And" || words[n] == "Or" {
			foundConnect = true
		} else if words[n] != "" {
			if foundConnect {
				for m := n; m < len(words)-1; m++ {
					words[n] += words[m+1]
					words[m+1] = ""
				}
			}
			break
		}
	}
	// }

	var params []*sqlgenerator.SQLParam
	// nextIsField := false
	sql := ""
	for _, wd := range words {
		if wd == "" {
			continue
		}

		if wd == "And" || wd == "Or" {
			if sql != "" {
				// if nextIsField {
				// 	if n == len(words)-1 {

				// 	}
				// }
				sql += wd + " "
				// nextIsField = true
			}
		} else {
			oper := "="
			suffix := ""
			if len(wd) > 2 {
				suffix = wd[len(wd)-2:]
			}
			hasParameter := true
			switch suffix {
			case "IN":
				oper = "IN"
				wd = wd[:len(wd)-2]
			case "GE":
				oper = ">="
				wd = wd[:len(wd)-2]
			case "GT":
				oper = ">"
				wd = wd[:len(wd)-2]
			case "LE":
				oper = "<="
				wd = wd[:len(wd)-2]
			case "LT":
				oper = "<"
				wd = wd[:len(wd)-2]
			case "NE":
				oper = "<>"
				wd = wd[:len(wd)-2]
			case "NN":
				oper = " is not null "
				wd = wd[:len(wd)-2]
				hasParameter = false
			case "NU":
				oper = " is null "
				wd = wd[:len(wd)-2]
				hasParameter = false
			case "EQ":
				wd = wd[:len(wd)-2]
			}

			// if sql == "" {
			// 	if n > 0 {
			// 		for m := n - 1; m >= 0; m-- {
			// 			wd = words[m] + wd
			// 		}
			// 	}
			// }

			if hasParameter {
				fieldName := names.ToTableName(wd)
				// pos := -1
				holder := ""
				if oper == "IN" {
					inPlaceIndex++
					holder = fmt.Sprintf("@in-%d", inPlaceIndex)
					sql += fieldName + " IN(" + holder + ") "
					// pos = len(sql) - 2 //   ^
				} else {
					sql += fieldName + oper + "? "
				}

				params = append(params, &sqlgenerator.SQLParam{
					SQLParamFieldName:  fieldName,
					SQLParamName:       wd,
					IsInOperator:       oper == "IN",
					InParamPlaceHolder: holder,
				})
			} else {
				sql += names.ToTableName(wd) + oper
			}
			// nextIsField = false
		}
	}
	return sql, params, nil
}

// // FieldVarPair 版本
// func parseWhereFromFuncName(whereInName string, fieldMapper names.Mapper) (string, []*desc.FieldVarPair, error) {
// 	// fnName := g.fn.Name
// 	// whereInName := fnName[len("UpdateBy"):]
// 	if len(whereInName) <= 0 {
// 		return "", nil, errors.New("func like UpdateByXxx needs condition field(s)")
// 	}
// 	words := []string{}
// 	var appendWord = func(word string, end bool) {
// 		words = append(words, word)
// 		if word == "And" || word == "Or" || end {
// 			// if (word == "And" || word == "Or") {
// 			// 	if len(words) > 0  {
// 			// 	}
// 			// }
// 			andOrSkip := 0
// 			if !end {
// 				andOrSkip = 1
// 			}
// 			var m int
// 			foundConnect := 0
// 			for m = len(words) - 1 - andOrSkip; m >= 0; m-- {
// 				wd := words[m]
// 				if wd == "And" || wd == "Or" {
// 					foundConnect++
// 				} else if foundConnect > 0 {
// 					m++
// 					break
// 				}
// 			}
// 			for k := m + 2; k < len(words)-andOrSkip; k++ {
// 				words[m+1] = words[m+1] + words[k]
// 				words[k] = ""
// 			}
// 		}
// 		// if !end {
// 		// 	words = append(words, word)
// 		// }
// 	}

// 	word := ""
// 	// bytes.ToLower()
// 	for n := range whereInName {
// 		ch := whereInName[n]
// 		if utils.IsCapital(ch) {
// 			if word != "" {
// 				appendWord(word, false)
// 			}
// 			word = string(ch)
// 		} else {
// 			word = word + string(ch)
// 		}
// 	}
// 	// if word != "And" && word != "Or" {
// 	appendWord(word, true)
// 	foundConnect := false
// 	for n := len(words) - 1; n >= 0; n-- {
// 		if words[n] == "And" || words[n] == "Or" {
// 			foundConnect = true
// 		} else if words[n] != "" {
// 			if foundConnect {
// 				for m := n; m < len(words)-1; m++ {
// 					words[n] += words[m+1]
// 					words[m+1] = ""
// 				}
// 			}
// 			break
// 		}
// 	}
// 	// }

// 	var params []*desc.FieldVarPair
// 	// nextIsField := false
// 	sql := ""
// 	for _, wd := range words {
// 		if wd == "" {
// 			continue
// 		}

// 		if wd == "And" || wd == "Or" {
// 			if sql != "" {
// 				// if nextIsField {
// 				// 	if n == len(words)-1 {

// 				// 	}
// 				// }
// 				sql += wd + " "
// 				// nextIsField = true
// 			}
// 		} else {
// 			oper := "="
// 			suffix := ""
// 			if len(wd) > 2 {
// 				suffix = wd[len(wd)-2:]
// 			}
// 			hasParameter := true
// 			switch suffix {
// 			case "IN":
// 				oper = "IN"
// 				wd = wd[:len(wd)-2]
// 			case "GE":
// 				oper = ">="
// 				wd = wd[:len(wd)-2]
// 			case "GT":
// 				oper = ">"
// 				wd = wd[:len(wd)-2]
// 			case "LE":
// 				oper = "<="
// 				wd = wd[:len(wd)-2]
// 			case "LT":
// 				oper = "<"
// 				wd = wd[:len(wd)-2]
// 			case "NE":
// 				oper = "<>"
// 				wd = wd[:len(wd)-2]
// 			case "NN":
// 				oper = " is not null "
// 				wd = wd[:len(wd)-2]
// 				hasParameter = false
// 			case "NU":
// 				oper = " is null "
// 				wd = wd[:len(wd)-2]
// 				hasParameter = false
// 			case "EQ":
// 				wd = wd[:len(wd)-2]
// 			}

// 			// if sql == "" {
// 			// 	if n > 0 {
// 			// 		for m := n - 1; m >= 0; m-- {
// 			// 			wd = words[m] + wd
// 			// 		}
// 			// 	}
// 			// }

// 			if hasParameter {
// 				fieldName := fieldMapper.Obj2Table(wd)
// 				pos := -1
// 				if oper == "IN" {
// 					sql += fieldName + " IN() "
// 					pos = len(sql) - 2 //   ^
// 				} else {
// 					sql += fieldName + oper + "? "
// 				}
// 				params = append(params, &desc.FieldVarPair{
// 					FieldName:           fieldName,
// 					VarName:             wd,
// 					IsInOperator:        oper == "IN",
// 					ParamInsertPosition: pos,
// 				})
// 			} else {
// 				sql += fieldMapper.Obj2Table(wd) + oper
// 			}
// 			// nextIsField = false
// 		}
// 	}
// 	return sql, params, nil
// }

func splitSQL(sql string) ([]string, error) {
	var quoteChar rune = 0
	var terms []string
	// var term string
	var lastSpaceIndex = -1
	for n, ch := range sql {
		if ch == '\'' || ch == '"' {
			if quoteChar == 0 {
				// 进入引号内容
				quoteChar = ch
			} else if quoteChar == ch {
				// 退出引号内容
				quoteChar = 0
			} // 否则作为普通字符处理
		}
		if quoteChar == 0 {
			// 在引号中，所有字符作为普通字符，包括空格
			if ch == ' ' {
				if lastSpaceIndex+1 != n {
					terms = append(terms, strings.TrimSpace(sql[lastSpaceIndex+1:n]))
				}
				lastSpaceIndex = n
			}
		}
	}
	if quoteChar != 0 {
		return nil, errors.New("sql error")
	}
	if lastSpaceIndex < len(sql)-1 {
		terms = append(terms, strings.TrimSpace(sql[lastSpaceIndex+1:]))
	}
	return terms, nil
}

type Param struct {
	Start int
	End   int
	Name  string
}

func ReplaceParam(sql string, p *Param, newStr string) string {
	return sql[:p.Start] + newStr + sql[p.End:]
}

func FindParams(term string) ([]*Param, error) {
	var quoteChar rune = 0
	paramsStart := -1
	var ps []*Param
	for n, ch := range term {
		if ch == '\'' || ch == '"' {
			if paramsStart != -1 {
				// 语法错误
				return nil, fmt.Errorf("sql syntax [%d]", n)
			}
			if quoteChar == 0 {
				// 进入引号内容
				quoteChar = ch
			} else if quoteChar == ch {
				// 退出引号内容
				quoteChar = 0
			} // 否则作为普通字符处理
		}
		if quoteChar == 0 {
			// 没有在引号内部
			if ch == ':' {
				if paramsStart != -1 {
					// 语法错误
					return nil, fmt.Errorf("sql syntax [%d]", n)
				}
				paramsStart = n
			} else {
				if paramsStart != -1 {
					validParamChar := (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
					if (!validParamChar) && n > paramsStart+1 {
						validParamChar = ch >= '0' && ch <= '9'
					}
					if !validParamChar {
						if n == paramsStart+1 {
							return nil, fmt.Errorf("sql syntax [%d]", n)
						}
						if ch != ',' {
							return nil, fmt.Errorf("sql syntax [%d]", n)
						}
						ps = append(ps, &Param{
							Start: paramsStart,
							End:   n,
							Name:  term[paramsStart+1 : n],
						})
						paramsStart = -1
					}
				}
			}
		}
	}
	if paramsStart != -1 {
		ps = append(ps, &Param{
			Start: paramsStart,
			End:   len(term),
			Name:  term[paramsStart+1:],
		})
	}
	return ps, nil
}

func ParseWhere(sqlTerms []string, whereIndex int) (string, []*sqlgenerator.SQLParam, error) {
	var whereTerms []string
	var inPlaceIndex = time.Now().Unix()
	var whereParams []*sqlgenerator.SQLParam
	if whereIndex >= 0 {
		for n := whereIndex + 1; n < len(sqlTerms); n++ {
			ps, err := FindParams(sqlTerms[n])
			if err != nil {
				return "", nil, err
			}
			col := sqlTerms[n]
			var termParams []*sqlgenerator.SQLParam
			for m := len(ps) - 1; m >= 0; m-- {
				var fieldName string
				k := n
				// for k := n; k > whereIndex; k-- {
				// if k == whereIndex {
				// 	break
				// }
				isIn := false
				term := sqlTerms[k]

				_, keyPos := lastIndex(term, "=", "<>", "<", "<=", ">", ">=")
				// eqPos := strings.LastIndex(term, "=")
				if keyPos > 0 {
					_, openParenPos := lastIndex(term, "(")
					if openParenPos > 0 {
						fieldName = term[openParenPos+1 : keyPos]
					} else {
						fieldName = term[:keyPos]
					}
					// break
				} else if keyPos == 0 {
					// 该 term 以 '"=", "<>", "<", "<=", ">", ">="' 开始， 从上一个 term 开始找
					if k > whereIndex+1 {
						term := sqlTerms[k-1]

						_, openParenPos := lastIndex(term, "(")
						if openParenPos > 0 {
							fieldName = term[openParenPos+1 : keyPos]
						} else {
							fieldName = term[:keyPos]
						}
						// break
					} else {
						return "", nil, fmt.Errorf("cann't find field name of params %s", ps[m].Name)
					}
				} else {
					// 该 term 中没有找到任何关键字
					if k == whereIndex+1 {
						term = sqlTerms[k-1]
						_, kp := lastIndex(term, "=", "<>", "<", "<=", ">", ">=")
						if kp > 0 {
							_, openParenPos := lastIndex(term, "(")
							if openParenPos > 0 {
								fieldName = term[openParenPos+1 : kp]
							} else {
								fieldName = term[:kp]
							}
							// break
						} else {
							return "", nil, fmt.Errorf("cann't find field name of params %s", ps[m].Name)
						}
					} else if k > whereIndex+1 {
						term = sqlTerms[k-1]
						kw := strings.ToLower(term)
						var kp int
						var wd string
						if wd, kp = lastIndex(kw, "=", "<>", "<", "<=", ">", ">=", "like", "in"); wd == kw {
							// 再向前找一个
							term = sqlTerms[k-2]
							kp = len(term)
							isIn = wd == "in"
						}

						_, openParenPos := lastIndex(term, "(")
						if openParenPos > 0 {
							fieldName = term[openParenPos+1 : kp]
						} else {
							fieldName = term[:kp]
						}
						// break
					} else {
						return "", nil, fmt.Errorf("cann't find field name of params %s", ps[m].Name)
					}
				}
				// }
				if fieldName == "" {
					return "", nil, fmt.Errorf("cann't find field name of params %s", ps[m].Name)
				}
				placeHolder := ""
				if isIn {
					inPlaceIndex++
					placeHolder = fmt.Sprintf("@in-%d", inPlaceIndex)
					col = ReplaceParam(col, ps[m], "("+placeHolder+")")
				} else {
					col = ReplaceParam(col, ps[m], "?")
				}
				termParams = append(termParams, &sqlgenerator.SQLParam{
					SQLParamName:       ps[m].Name,
					SQLParamFieldName:  fieldName,
					IsInOperator:       isIn,
					InParamPlaceHolder: placeHolder,
				})
			}
			for n := len(termParams) - 1; n >= 0; n-- {
				whereParams = append(whereParams, termParams[n])
			}
			whereTerms = append(whereTerms, col)
		}
	}

	return strings.Join(whereTerms, " "), whereParams, nil
}

func lastIndex(s string, substrs ...string) (string, int) {
	for _, ss := range substrs {
		if p := strings.LastIndex(s, ss); p >= 0 {
			return ss, p
		}
	}
	return "", -1
}
