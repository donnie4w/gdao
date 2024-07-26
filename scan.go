// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdao

import (
	"fmt"
	. "github.com/donnie4w/gdao/base"
	"reflect"
	"strings"
)

type scaner[T any] DataBean

func (s *scaner[T]) Scan() (t *T, err error) {
	defer Recover(&err)
	if s == nil {
		return nil, nil
	}
	dataBean := (*DataBean)(s)
	t = new(T)
	if scanner, ok := any(t).(Scanner); ok {
		scanner.ToGdao()
		for _, fieldBean := range dataBean.FieldMapName {
			scanner.Scan(fieldBean.FieldName, fieldBean.Value())
		}
		return any(scanner).(*T), nil
		//val := reflect.ValueOf(scanner).Elem()
		//return val.Addr().Interface().(*T), nil
	}
	valptr := reflect.ValueOf(t)
	val := valptr.Elem()
	if val.Kind() == reflect.Ptr {
		return nil, fmt.Errorf("Scan failed, the generic parameter cannot be of pointer type: %v", valptr.Type())
	}

	typ := val.Type()
	hasScan := true
	num := typ.NumField()
	if num < dataBean.Len() {
		for i := 0; i < num; i++ {
			field := typ.Field(i)
			fieldName := strings.ToLower(decodeFieldname(field.Name))
			if value := dataBean.ValueByName(fieldName); value != nil {
				ScanValue(val.Field(i), value)
			} else {
				hasScan = false
				break
			}
		}
	}

	if !hasScan || num >= dataBean.Len() {
		for _, fieldBean := range dataBean.FieldMapName {
			fieldName := fieldBean.FieldName
			field := val.FieldByNameFunc(func(s string) bool {
				return strings.EqualFold(s, fieldName)
			})
			if field.IsValid() && field.CanSet() {
				ScanValue(field, fieldBean.Value())
			} else {
				columnName := upperFirstLetter(fieldName)
				col := val.MethodByName("Set" + columnName)
				if !col.IsValid() {
					col = valptr.MethodByName("Set" + columnName)
				}
				if col.IsValid() {
					methodType := col.Type()
					numIn := methodType.NumIn()
					if numIn == 1 {
						expectedType := methodType.In(0)
						if val := GetValue(expectedType, fieldBean.Value()); val != nil {
							args := []reflect.Value{reflect.ValueOf(val)}
							col.Call(args)
						}
					}
				} else {
					if Logger.IsVaild {
						Logger.Warn("Failed to assign value to the field [", fieldName, "]")
					}
				}
			}
		}
	}
	return
}
