// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoStruct

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func isBasicType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool, reflect.String:
		return true
	case reflect.Ptr:
		return isBasicType(t.Elem())
	case reflect.Slice:
		// Check if it is a []byte
		return t.Elem().Kind() == reflect.Uint8
	default:
		if t == reflect.TypeOf(time.Time{}) {
			return true
		}
	}
	return false
}

func ToMap(arg any) map[string]any {
	switch v := arg.(type) {
	case map[string]any:
		return v
	}
	result := make(map[string]any)
	val := reflect.ValueOf(arg)
	valStruct := val
	typ := reflect.TypeOf(arg)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		valStruct = val.Elem()
	}
	if typ.Kind() != reflect.Struct {
		panic(fmt.Sprintf("argument must be a struct or a pointer to struct: %s", typ))
	}
	if _, ok := arg.(TableClass); !ok {
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			fieldValue := valStruct.Field(i)
			if field.PkgPath != "" {
				continue
			}
			if isBasicType(field.Type) {
				if !fieldValue.IsZero() {
					result[strings.ToLower(field.Name)] = fieldValue.Interface()
				}
			}
		}
	}
	for i := 0; i < reflect.TypeOf(arg).NumMethod(); i++ {
		method := reflect.TypeOf(arg).Method(i)
		if strings.HasPrefix(method.Name, "Get") && method.Type.NumIn() == 1 && method.Type.NumOut() == 1 {
			name := strings.ToLower(method.Name[3:])
			if _, exists := result[name]; !exists {
				out := val.Method(i).Call(nil)
				if len(out) > 0 && !out[0].IsZero() {
					result[name] = out[0].Interface()
				}
			}
		}
	}
	return result
}
