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
	"github.com/donnie4w/gdao/gdaoSlave"
)

type dbHandler struct {
	gdbc gdbcHandle
}

func init() {
	gdaoSlave.Newdbhandle = newdbhandle
}

func newdbhandle(db *sql.DB, dbtype DBType) DBhandle {
	return &dbHandler{gdbc: newGdbcHandle(nil, db, dbtype)}
}

func (h *dbHandler) GetDB() *sql.DB {
	return h.gdbc.GetDB()
}

func (h *dbHandler) Close() error {
	return h.gdbc.Close()
}

func (h *dbHandler) GetDBType() DBType {
	return h.gdbc.GetDBType()
}

func (h *dbHandler) GetTransaction() (r Transaction, err error) {
	return NewTransactionWithDBhandle(h)
}

func (h *dbHandler) ExecuteQueryBean(sqlstr string, args ...any) *DataBean {
	return h.gdbc.ExecuteQueryBean(sqlstr, args...)
}

func (h *dbHandler) ExecuteQueryBeans(sqlstr string, args ...any) (r *DataBeans) {
	return h.gdbc.ExecuteQueryBeans(sqlstr, args...)
}

func (h *dbHandler) ExecuteUpdate(sqlstr string, args ...any) (sql.Result, error) {
	return h.gdbc.ExecuteUpdate(sqlstr, args...)
}

func (h *dbHandler) ExecuteBatch(sqlstr string, args [][]any) (r []sql.Result, err error) {
	return h.gdbc.ExecuteBatch(sqlstr, args)
}
