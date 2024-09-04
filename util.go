// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdao

import (
	"database/sql"
	"fmt"
	"github.com/donnie4w/gdao/base"
	"reflect"
	"strconv"
)

var errInit = fmt.Errorf("the gdao DataSource was not initialized(Hint: gdao.Init(db, dbtype))")

const (
	_ base.DBType = iota
	MYSQL
	POSTGRESQL
	MARIADB
	SQLITE
	ORACLE
	SQLSERVER
	DB2
	SYBASE
	DERBY
	FIREBIRD
	INGRES
	GREENPLUM
	TERADATA
	NETEZZA
	VERTICA
	TIDB
	OCEANBASE
	OPENGAUSS
	HSQLDB
	ENTERPRISEDB
	SAPHANA
	COCKROACHDB
	INFORMIX
)

func SetLogger(on bool) {
	base.Logger.SetLogger(on)
}

func parseSql(dbtype base.DBType, sqlstr string, args ...any) string {
	if len(args) > 0 {
		switch dbtype {
		case POSTGRESQL, GREENPLUM, OPENGAUSS:
			s := ""
			k := 1
			for _, c := range sqlstr {
				if c == '?' {
					s = s + "$" + strconv.Itoa(k)
					k++
				} else {
					s = s + string(c)
				}
			}
			return s
		case ORACLE:
			s := ""
			k := 1
			for _, c := range sqlstr {
				if c == '?' {
					s = s + ":v" + strconv.Itoa(k)
					k++
				} else {
					s = s + string(c)
				}
			}
			for i, arg := range args {
				if vs, ok := arg.([]any); ok {
					for j, v := range vs {
						vs[j] = sql.Named("v"+strconv.Itoa(j+1), v)
					}
				} else {
					args[i] = sql.Named("v"+strconv.Itoa(i+1), arg)
				}
			}
			return s
		}
	}
	return sqlstr
}

func toArray(arg any) (r []any) {
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
