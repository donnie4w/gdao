// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoMapper

import (
	"fmt"
	. "github.com/donnie4w/gdao/gdaoStruct"
	"github.com/donnie4w/gdao/util"
	"reflect"
	"regexp"
	"strings"
)

type paramBean struct {
	namespace      string
	id             string
	sqltype        sqlType
	sql            string
	parameterNames []string
	inputType      string
	outputType     string
	sqlNode        sqlNode
}

func newParamBean2(namespace, id, sql, inputType, outputType string, sqltype sqlType) *paramBean {
	p := &paramBean{namespace: namespace, id: id, sqltype: sqltype, sql: sql, inputType: inputType, outputType: outputType}
	return p
}

func newParamBean(namespace, id, sqltype, sql, inputType, outputType string) (r *paramBean) {
	r = &paramBean{namespace: namespace, id: id, inputType: inputType, outputType: outputType}
	switch sqltype {
	case "select":
		r.sqltype = _SELECT
	case "insert":
		r.sqltype = _INSERT
	case "update":
		r.sqltype = _UPDATE
	case "delete":
		r.sqltype = _DELETE
	}
	r.sql, r.parameterNames = parseSql(sql)
	return
}

func parseSql(sqlstr string) (sql string, parameterNames []string) {
	pattern := `#\{(\w+)\}`
	re := regexp.MustCompile(pattern)
	sql = re.ReplaceAllStringFunc(sqlstr, func(match string) string {
		parameterNames = append(parameterNames, strings.Trim(match[2:len(match)-1], " "))
		return "?"
	})
	return strings.TrimSpace(sql), parameterNames
}

func (p *paramBean) hasSqlNode() bool {
	return p.sqlNode != nil
}

func (p *paramBean) err_num_no_match(expectvalue, gotvalue int) error {
	return fmt.Errorf("the parameter number does not match the configuration:Expected %d but got %d. [namespace:%s][mapper id:%s]", expectvalue, gotvalue, p.namespace, p.id)
}

func (p *paramBean) err_invalid_parameter(getType string) error {
	return fmt.Errorf("invalid parameter type: Expected %s but get %s [namespace:%s][mapper id:%s]", p.inputType, getType, p.namespace, p.id)
}

func (p *paramBean) err_params_no_found(param string) error {
	return fmt.Errorf("parameter not found:  %s  [namespace:%s][mapper id:%s]", param, p.namespace, p.id)
}

func (p *paramBean) setParameter(parameter any) (args []any, err error) {
	defer util.Recover(&err)
	if p.inputType == "" || parameter == nil || len(p.parameterNames) == 0 {
		return args, nil
	}
	if isDBType(p.inputType) {
		if len(p.parameterNames) == 1 {
			return []any{parameter}, nil
		} else {
			return nil, p.err_num_no_match(len(p.parameterNames), 1)
		}
	}
	if isSlice(p.inputType) {
		typ := reflect.TypeOf(parameter)
		if typ.Kind() == reflect.Slice {
			return p.toAnySlice(parameter, typ.Elem().String())
		} else {
			fmt.Println("The variable is not a slice.")
		}
	}
	if p.inputType == "map" {
		if m, ok := parameter.(map[string]any); ok {
			args = make([]any, 0)
			for _, name := range p.parameterNames {
				if v, ok := m[name]; ok {
					args = append(args, v)
				} else {
					return nil, p.err_params_no_found(name)
				}
			}
			if len(args) == len(p.parameterNames) {
				return args, nil
			} else {
				return nil, p.err_num_no_match(len(p.parameterNames), len(args))
			}
		} else {
			return nil, p.err_invalid_parameter(reflect.TypeOf(parameter).String())
		}
	}

	val := reflect.ValueOf(parameter)
	typ := reflect.TypeOf(parameter)
	if val.Kind() == reflect.Struct {
		ptr := reflect.New(val.Type())
		ptr.Elem().Set(val)
		val = ptr
		typ = val.Type()
	}
	kind := typ.Kind()
	switch kind {
	case reflect.Struct, reflect.Ptr:
		args = make([]any, 0)
		for _, name := range p.parameterNames {
			method := val.MethodByName("Get" + upperFirstLetter(name))
			if method.IsValid() {
				//in := []reflect.Value{}
				//if method.Type().NumIn() == 1 {
				//	in = []reflect.Value{val}
				//}
				if vals := method.Call(nil); len(vals) > 0 {
					args = append(args, vals[0].Interface())
				}
			} else {
				return nil, p.err_params_no_found(name)
			}
		}
		if len(args) == len(p.parameterNames) {
			return args, nil
		} else {
			return nil, p.err_num_no_match(len(p.parameterNames), len(args))
		}
	}
	return nil, p.err_invalid_parameter(reflect.TypeOf(parameter).String())
}

func (p *paramBean) toAnySlice(parameter any, typename string) (r []any, err error) {
	length := len(p.parameterNames)
	switch typename {
	case "string":
		if slice, ok := parameter.([]string); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "int64":
		if slice, ok := parameter.([]int64); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "int":
		if slice, ok := parameter.([]int); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "int32":
		if slice, ok := parameter.([]int32); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "int16":
		if slice, ok := parameter.([]int16); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "int8":
		if slice, ok := parameter.([]int8); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "uint64":
		if slice, ok := parameter.([]uint64); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "uint32":
		if slice, ok := parameter.([]uint32); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "uint16":
		if slice, ok := parameter.([]uint16); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "uint8":
		if slice, ok := parameter.([]uint8); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "float64":
		if slice, ok := parameter.([]float64); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "float32":
		if slice, ok := parameter.([]float32); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "bool":
		if slice, ok := parameter.([]bool); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	default:
		return nil, p.err_invalid_parameter(reflect.TypeOf(parameter).String())
	}
	return
}

func isDBType(putType string) bool {
	if putType[0:1] == "*" {
		putType = putType[1:]
	}
	switch putType {
	case "int64", "int32", "int", "int16", "int8":
		return true
	case "float64", "float32":
		return true
	case "uint64", "uint32", "uint", "uint16", "uint8":
		return true
	case "string", "byte", "rune", "[]byte", "bool", "time", "Time", "time.Time":
		return true
	}
	return false
}

func isSlice(putType string) bool {
	if putType[:2] == "[]" {
		return true
	}
	return false
}

func upperFirstLetter(s string) string {
	return strings.ToUpper(string(s[0])) + s[1:]
}

func (p *paramBean) parseSqlNode(parameter any) (r *paramBean, params []any) {
	ac := p.sqlNode.apply(NewParamContext(parameter))
	if ac.params != nil {
		params = ac.params
	}
	return newParamBean2(p.namespace, p.id, ac.GetSql(), p.inputType, p.outputType, p.sqltype), params
}

func (p *paramBean) parseSqlNode2(args ...any) (r *paramBean, params []any) {
	ac := p.sqlNode.apply(NewParamContext2(args...))
	if ac.params != nil {
		params = ac.params
	}
	return newParamBean2(p.namespace, p.id, ac.GetSql(), p.inputType, p.outputType, p.sqltype), params
}
