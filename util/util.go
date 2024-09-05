// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package util

import (
	"fmt"
	"reflect"
	"strings"
)

func ToUpperFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

//func EncodeFieldname(k string) string {
//	if Iskey(k) {
//		return k + "_"
//	}
//	return k
//}

//func Iskey(name string) bool {
//	switch name {
//	case "break", "default", "func", "interface", "select", "case", "defer", "go", "map", "struct", "chan", "else", "goto", "package", "switch", "const", "fallthrough", "if", "range", "type", "continue", "for", "import", "return", "var":
//		return true
//	default:
//		return false
//	}
//}

//func DecodeFieldname(name string) string {
//	if name[len(name)-1:] == "_" {
//		if n := name[:len(name)-1]; Iskey(n) {
//			return n
//		}
//	}
//	return name
//}

func Recover(errp *error) {
	if r := recover(); r != nil {
		if errp != nil {
			*errp = fmt.Errorf("panic recovering: %v", r)
		}
	}
}

func Classname[T any]() string {
	var t T
	tType := reflect.TypeOf(t)
	if tType.Kind() == reflect.Ptr {
		return tType.Elem().String()
	} else {
		return tType.String()
	}
}

func ToArray(arg any) (r []any) {
	value := reflect.ValueOf(arg)
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		r = make([]any, value.Len(), value.Len())
		for i := 0; i < value.Len(); i++ {
			r[i] = value.Index(i).Interface()
		}
	case reflect.Map:
		r = make([]any, 0)
		for _, key := range value.MapKeys() {
			if _, ok := key.Interface().(string); ok {
				r = append(r, value.MapIndex(key).Interface())
			}
		}
	}
	return
}
