// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoSlave

import (
	"database/sql"
	"github.com/donnie4w/gdao/base"
)

var (
	BindWithDBhandle       func(tablename string, dbhandle base.DBhandle)
	Bind                   func(tablename string, db *sql.DB, dbtype base.DBType)
	BindMapperWithDBhandle func(mapperId string, dbhandle base.DBhandle)
	BindMapper             func(mapperId string, db *sql.DB, dbtype base.DBType)
	Remove                 func(tablename string) bool
	RemoveMapper           func(mapperId string) bool

	Len func() int64
	Get func(classname, tableName, mapperId string) base.DBhandle
)

func BindClass[T any](db *sql.DB, dbtype base.DBType) {
	classname := base.Classname[T]()
	Bind(classname, db, dbtype)
}

func BindClassWithDBhandle[T any](dbhandle base.DBhandle) {
	classname := base.Classname[T]()
	BindWithDBhandle(classname, dbhandle)
}

func RemoveClass[T any]() bool {
	classname := base.Classname[T]()
	return Remove(classname)
}
