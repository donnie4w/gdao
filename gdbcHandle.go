// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdao

import (
	"database/sql"
	"github.com/donnie4w/gdao/base"
)

type gdbcHandle interface {
	ExecuteQueryBeans(sqlstr string, args ...any) *base.DataBeans
	ExecuteQueryBean(sqlstr string, args ...any) *base.DataBean
	ExecuteUpdate(sqlstr string, args ...any) (int64, error)
	ExecuteBatch(sqlstr string, args [][]any) (r []int64, err error)
	GetDBType() base.DBType
	GetDB() *sql.DB
	Close() error
}

type gdbcHandler struct {
	TX     *sql.Tx
	DB     *sql.DB
	DBType base.DBType
}

func newGdbcHandle(tx *sql.Tx, db *sql.DB, dbType base.DBType) gdbcHandle {
	return &gdbcHandler{TX: tx, DB: db, DBType: dbType}
}

func (g *gdbcHandler) GetDBType() base.DBType {
	return g.DBType
}
func (g *gdbcHandler) GetDB() *sql.DB {
	return g.DB
}

func (g *gdbcHandler) ExecuteQueryBeans(sqlstr string, args ...any) (r *base.DataBeans) {
	r = &base.DataBeans{}
	sqlstr = parseSql(g.DBType, sqlstr, args)
	if dbs, err := executeQueryBeans(g.TX, g.DB, sqlstr, args...); err == nil {
		r.Beans = dbs
	} else {
		r.SetError(err)
	}
	return
}

func (g *gdbcHandler) ExecuteQueryBean(sqlstr string, args ...any) (r *base.DataBean) {
	sqlstr = parseSql(g.DBType, sqlstr, args)
	if db, err := executeQueryBean(g.TX, g.DB, sqlstr, args...); err == nil {
		return db
	} else {
		r = &base.DataBean{}
		r.SetError(err)
	}
	return
}

func (g *gdbcHandler) ExecuteUpdate(sqlstr string, args ...any) (int64, error) {
	sqlstr = parseSql(g.DBType, sqlstr, args)
	return executeUpdate(g.TX, g.DB, sqlstr, args...)
}

func (g *gdbcHandler) ExecuteBatch(sqlstr string, args [][]any) ([]int64, error) {
	sqlstr = parseSql(g.DBType, sqlstr, args)
	return executeBatch(g.TX, g.DB, sqlstr, args)
}

func (g *gdbcHandler) Close() error {
	return g.DB.Close()
}
