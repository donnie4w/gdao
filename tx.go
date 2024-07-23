// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdao

import (
	"database/sql"
	. "github.com/donnie4w/gdao/base"
)

type tx struct {
	tx      *sql.Tx
	dbtype  DBType
	gdbc    gdbcHandle
	isclose bool
}

func newTX(db DBhandle) (x *tx, err error) {
	x = new(tx)
	if x.tx, err = db.GetDB().Begin(); err == nil {
		x.dbtype = db.GetDBType()
		x.gdbc = newGdbcHandle(x.tx, db.GetDB(), db.GetDBType())
	}
	return x, err
}

func (x *tx) IsClose() bool {
	return x.isclose
}

func (x *tx) Commit() (err error) {
	return x.tx.Commit()
}

func (x *tx) Rollback() (err error) {
	return x.tx.Rollback()
}

func (x *tx) Close() (err error) {
	return
}

func (x *tx) GetDBType() DBType {
	return x.dbtype
}

func (x *tx) GetTransaction() (Transaction, error) {
	return x, nil
}

func (x *tx) GetDB() *sql.DB {
	return x.gdbc.GetDB()
}

func (x *tx) ExecuteUpdate(sqlstr string, args ...any) (int64, error) {
	return x.gdbc.ExecuteUpdate(sqlstr, args...)
}

func (x *tx) ExecuteBatch(sqlstr string, args [][]any) (r []int64, err error) {
	return x.gdbc.ExecuteBatch(sqlstr, args)
}

func (x *tx) ExecuteQueryBeans(sqlstr string, args ...any) ([]*DataBean, error) {
	return x.gdbc.ExecuteQueryBeans(sqlstr, args...)
}

func (x *tx) ExecuteQueryBean(sqlstr string, args ...any) (*DataBean, error) {
	return x.gdbc.ExecuteQueryBean(sqlstr, args...)
}
