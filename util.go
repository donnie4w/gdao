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
)

func SetLogger(on bool) {
	base.Logger.SetLogger(on)
}

func iskey(name string) bool {
	switch name {
	case "break", "default", "func", "interface", "select", "case", "defer", "go", "map", "struct", "chan", "else", "goto", "package", "switch", "const", "fallthrough", "if", "range", "type", "continue", "for", "import", "return", "var":
		return true
	default:
		return false
	}
}

func encodeFieldname(name string) string {
	if iskey(name) {
		return name + "_"
	}
	return name
}

func decodeFieldname(name string) string {
	if name[len(name)-1:] == "_" {
		if n := name[:len(name)-1]; iskey(n) {
			return n
		}
	}
	return name
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
