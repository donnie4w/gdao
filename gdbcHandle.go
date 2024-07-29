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
	ExecuteQueryBeans(sqlstr string, args ...any) ([]*base.DataBean, error)
	ExecuteQueryBean(sqlstr string, args ...any) (*base.DataBean, error)
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

func (g *gdbcHandler) ExecuteQueryBeans(sqlstr string, args ...any) ([]*base.DataBean, error) {
	sqlstr = parseSql(g.DBType, sqlstr, args)
	return executeQueryBeans(g.TX, g.DB, sqlstr, args...)
}

func (g *gdbcHandler) ExecuteQueryBean(sqlstr string, args ...any) (*base.DataBean, error) {
	sqlstr = parseSql(g.DBType, sqlstr, args)
	return executeQueryBean(g.TX, g.DB, sqlstr, args...)
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

func ExecuteQuery[T any](sql string, args ...any) (r *T, err error) {
	if databean, err := defaultDBhandle.ExecuteQueryBean(sql, args...); err == nil {
		return Scan[T](databean)
	} else {
		return nil, err
	}
}

func ExecuteQueryList[T any](sql string, args ...any) (r []*T, err error) {
	var databeans []*base.DataBean
	if databeans, err = defaultDBhandle.ExecuteQueryBeans(sql, args...); err == nil && len(databeans) > 0 {
		r = make([]*T, 0)
		for _, databean := range databeans {
			var t *T
			if t, err = Scan[T](databean); err == nil {
				r = append(r, t)
			}
		}
	}
	return
}

func ExecuteQueryBean(sql string, args ...any) (*base.DataBean, error) {
	if defaultDBhandle == nil {
		return nil, errInit
	}
	return defaultDBhandle.ExecuteQueryBean(sql, args...)
}

func ExecuteQueryBeans(sql string, args ...any) ([]*base.DataBean, error) {
	if defaultDBhandle == nil {
		return nil, errInit
	}
	return defaultDBhandle.ExecuteQueryBeans(sql, args...)
}

func ExecuteUpdate(sql string, args ...any) (int64, error) {
	if defaultDBhandle == nil {
		return 0, errInit
	}
	return defaultDBhandle.ExecuteUpdate(sql, args...)
}

func ExecuteBatch(sql string, args [][]any) ([]int64, error) {
	if defaultDBhandle == nil {
		return nil, errInit
	}
	return defaultDBhandle.ExecuteBatch(sql, args)
}
