// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package util

import (
	"errors"
	"fmt"
	"github.com/donnie4w/gdao/gdaoStruct"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func ResolveVariableValue(variable string, paramContext *gdaoStruct.ParamContext) (value any, err error) {
	if strings.Contains(variable, ".") {
		return getValueByPath(variable, paramContext)
	} else if strings.Contains(variable, "[") {
		startIndex := strings.Index(variable, "[")
		endIndex := strings.Index(variable, "]")
		if startIndex > 0 && endIndex > startIndex {
			objectName := variable[:startIndex]
			indexStr := variable[startIndex+1 : endIndex]
			index := 0
			if isNumeric(indexStr) {
				if index, err = strconv.Atoi(indexStr); err != nil {
					return nil, err
				}
			} else if foreachContext := paramContext.ForeachContext; foreachContext != nil {
				index = foreachContext.GetIndex(indexStr)
			}
			value = getObjectFromContext(objectName, paramContext)
			if value != nil && reflect.TypeOf(value).Kind() == reflect.Slice {
				sliceValue := reflect.ValueOf(value)
				if index >= 0 && index < sliceValue.Len() {
					value = sliceValue.Index(index).Interface()
				} else {
					value = nil
				}
			}
		}
	} else {
		return getValueByPath(variable, paramContext)
	}
	return
}

func getValueByPath(path string, obj any) (any, error) {
	if obj == nil || path == "" {
		return nil, nil
	}
	parts := strings.Split(path, ".")
	currentObj := obj
	var part string

	for _, part = range parts {
		if currentObj == nil {
			break
		}
		if v, err := getValueByFiled(part, currentObj); v == nil {
			if strings.ToLower(part) != part {
				if currentObj, err = getValueByFiled(strings.ToLower(part), currentObj); err != nil {
					return nil, err
				}
			}
			currentObj = v
		} else {
			currentObj = v
		}
	}

	if currentObj == nil {
		return nil, errors.New("Can't find object or value [" + part + "] in " + path)
	}

	return currentObj, nil
}

func getValueByFiled(filed string, currentObj any) (any, error) {
	if currentObj == nil {
		return nil, nil
	}

	switch v := currentObj.(type) {
	case gdaoStruct.ParamContext:
		currentObj = getObjectFromContext(filed, &v)
	case *gdaoStruct.ParamContext:
		currentObj = getObjectFromContext(filed, v)
	case map[string]any:
		currentObj = v[filed]
	default:
		fieldValue, err := getFieldValue(currentObj, filed)
		if err != nil {
			return nil, err
		}
		currentObj = fieldValue
	}
	return currentObj, nil
}

func getFieldValue(obj any, fieldName string) (any, error) {
	if obj == nil || fieldName == "" {
		return nil, errors.New("object or fieldName cannot be nil or empty")
	}

	if _, ok := obj.(gdaoStruct.TableClass); ok {
		value := reflect.ValueOf(obj)
		if value.Kind() != reflect.Ptr {
			ptrToObj := reflect.New(reflect.TypeOf(obj))
			ptrToObj.Elem().Set(value)
			value = ptrToObj
		}
		methodName := "Get" + strings.ToUpper(fieldName[:1]) + fieldName[1:]
		method := value.MethodByName(methodName)
		if method.IsValid() && method.Type().NumIn() == 0 {
			results := method.Call(nil)
			if len(results) > 0 {
				return results[0].Interface(), nil
			}
		} else {
			return nil, fmt.Errorf("method %s does not exist or is not valid", methodName)
		}
	}

	value := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
		typ = typ.Elem()
	}

	if value.Kind() != reflect.Struct {
		return nil, errors.New("object or fieldName must be struct")
	}

	field := value.FieldByName(fieldName)
	if field.IsValid() {
		return field.Interface(), nil
	}
	field = value.FieldByNameFunc(func(s string) bool {
		return strings.ToLower(fieldName) == strings.ToLower(s)
	})
	if field.IsValid() {
		return field.Interface(), nil
	}

	return nil, errors.New("field or getter method not found")
}

func getObjectFromContext(key string, context *gdaoStruct.ParamContext) any {
	var value any
	if foreachCtx := context.ForeachContext; foreachCtx != nil {
		value = foreachCtx.GetItem(key)
		if value == nil {
			value = foreachCtx.GetIndex(key)
		}
	}
	if value == nil {
		value = context.Get(key)
	}
	return value
}

func extractVariables(expression string, paramContext *gdaoStruct.ParamContext) string {
	var currentVariable strings.Builder
	insideQuotes := false
	var result strings.Builder
	for i := 0; i < len(expression); i++ {
		c := rune(expression[i])
		if c == '\'' || c == '"' {
			insideQuotes = !insideQuotes
		}
		if insideQuotes {
			result.WriteRune(c)
			continue
		}
		if unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' || c == '.' || c == '[' || c == ']' {
			currentVariable.WriteRune(c)
		} else {
			if currentVariable.Len() > 0 {
				s := currentVariable.String()
				if isVariable(s) {
					value, _ := ResolveVariableValue(s, paramContext)
					if value != nil {
						switch v := value.(type) {
						case string:
							result.WriteString(fmt.Sprintf("'%s'", v))
						default:
							result.WriteString(fmt.Sprintf("%v", v))
						}
					} else {
						result.WriteString("nil")
					}
				} else {
					result.WriteString(s)
				}
				currentVariable.Reset()
			}
			result.WriteRune(c)
		}
	}

	if currentVariable.Len() > 0 {
		s := currentVariable.String()
		if isVariable(s) {
			value, _ := ResolveVariableValue(s, paramContext)
			if value != nil {
				switch v := value.(type) {
				case string:
					result.WriteString(fmt.Sprintf("'%s'", v))
				default:
					result.WriteString(fmt.Sprintf("%v", v))
				}
			} else {
				result.WriteString("nil")
			}
		} else {
			result.WriteString(s)
		}
	}

	return result.String()
}

func isVariable(s string) bool {
	if s == "nil" || s == "true" || s == "false" {
		return false
	}
	matched, _ := regexp.MatchString("^\\d+$", s)
	return !matched
}
