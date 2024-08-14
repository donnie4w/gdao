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

func EncodeFieldname(k string) string {
	if Iskey(k) {
		return k + "_"
	}
	return k
}

func Iskey(name string) bool {
	switch name {
	case "break", "default", "func", "interface", "select", "case", "defer", "go", "map", "struct", "chan", "else", "goto", "package", "switch", "const", "fallthrough", "if", "range", "type", "continue", "for", "import", "return", "var":
		return true
	default:
		return false
	}
}

func DecodeFieldname(name string) string {
	if name[len(name)-1:] == "_" {
		if n := name[:len(name)-1]; Iskey(n) {
			return n
		}
	}
	return name
}

func Recover(errp *error) {
	if r := recover(); r != nil {
		*errp = fmt.Errorf("panic recovering: %v", r)
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
