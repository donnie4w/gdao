// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdao

import (
	"database/sql"
	"errors"
	. "github.com/donnie4w/gdao/base"
	"github.com/donnie4w/gdao/util"
	"github.com/donnie4w/gofer/hashmap"
	"sync/atomic"
)

var sqlWare = hashmap.NewLimitMap[string, *int64](1 << 19)
var stmtExec = &stmtexec{stmtMap: hashmap.NewMapL[string, *sql.Stmt]()}
var errorStmt = errors.New("")

type stmtexec struct {
	stmtMap *hashmap.MapL[string, *sql.Stmt]
	lock    int64
}

func (se *stmtexec) Exec(db *sql.DB, sqlStr string, args ...any) (rs sql.Result, err error) {
	if stmt, err := se.Prepare(sqlStr, db); err == nil {
		return stmt.Exec(args...)
	} else {
		return db.Exec(sqlStr, args...)
	}
}

func (se *stmtexec) Qurey(db *sql.DB, sqlStr string, args ...any) (rs *sql.Rows, err error) {
	if stmt, err := se.Prepare(sqlStr, db); err == nil {
		return stmt.Query(args...)
	} else {
		return db.Query(sqlStr, args...)
	}
}

func (se *stmtexec) clear() {
	if atomic.CompareAndSwapInt64(&se.lock, 0, stmtLimit) {
		defer atomic.StoreInt64(&se.lock, 0)
		if se.stmtMap.Len() >= stmtLimit {
			se.stmtMap.Range(func(k string, v *sql.Stmt) bool {
				v.Close()
				se.stmtMap.Del(k)
				return true
			})
		}
	}
}

func (se *stmtexec) Prepare(sqlStr string, db *sql.DB) (stmt *sql.Stmt, err error) {
	if se.len() >= stmtLimit {
		se.clear()
		return stmt, errorStmt
	}
	if v, ok := se.stmtMap.Get(sqlStr); ok {
		return v, nil
	}
	if stmt, err = db.Prepare(sqlStr); err == nil {
		if p, ok := se.stmtMap.Swap(sqlStr, stmt); ok && p != nil {
			p.Close()
		}
	}
	return
}

func (se *stmtexec) len() int64 {
	if se.lock > 0 {
		return se.lock
	}
	return se.stmtMap.Len()
}

func (se *stmtexec) nostmt(tx *sql.Tx, sqlstr string) (b bool) {
	if stmtLimit == 0 || tx != nil || se.len() >= stmtLimit {
		return true
	}
	if v, ok := sqlWare.Get(sqlstr); ok {
		b = atomic.AddInt64(v, 1) < 16
	} else {
		sqlWare.Put(sqlstr, new(int64))
		b = true
	}
	return
}

func (se *stmtexec) executeQueryBeans(tx *sql.Tx, db *sql.DB, sqlstr string, args ...any) (databases []*DataBean, err error) {
	if se.nostmt(tx, sqlstr) {
		return executeQueryBeans(tx, db, sqlstr, args...)
	}
	if tx == nil && db == nil {
		return nil, errInit
	}
	var rows *sql.Rows
	if tx != nil {
		return executeQueryBeans(tx, db, sqlstr, args...)
	} else {
		rows, err = se.Qurey(db, sqlstr, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	databases = make([]*DataBean, 0)
	if names, er := rows.Columns(); er == nil {
		for rows.Next() {
			databean := NewDataBean(len(names))
			buff := newAnys(len(names))
			for _, name := range names {
				fb := NewFieldBeen()
				*buff = append(*buff, &fb.FieldValue)
				databean.Put(name, fb)
			}
			if err = rows.Scan(*buff...); err != nil {
				databean.SetError(err)
				return
			}
			bufpool.Put(&buff)
			databases = append(databases, databean)
		}
	} else {
		err = er
	}
	return
}

func (se *stmtexec) executeQueryBean(tx *sql.Tx, db *sql.DB, sqlstr string, args ...any) (dataBean *DataBean, err error) {
	if se.nostmt(tx, sqlstr) {
		return executeQueryBean(tx, db, sqlstr, args...)
	}
	if tx == nil && db == nil {
		return nil, errInit
	}
	var rows *sql.Rows
	if tx != nil {
		return executeQueryBean(tx, db, sqlstr, args...)
	} else {
		rows, err = se.Qurey(db, sqlstr, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if names, er := rows.Columns(); er == nil {
		dataBean = NewDataBean(len(names))
		if rows.Next() {
			buff := newAnys(len(names))
			for _, name := range names {
				fb := NewFieldBeen()
				*buff = append(*buff, &fb.FieldValue)
				dataBean.Put(name, fb)
			}
			if err = rows.Scan(*buff...); err != nil {
				dataBean.SetError(err)
				return
			}
			bufpool.Put(&buff)
		}
	} else {
		err = er
	}
	return
}

func (se *stmtexec) executeUpdate(tx *sql.Tx, db *sql.DB, sqlstr string, args ...any) (r int64, err error) {
	if se.nostmt(tx, sqlstr) {
		return executeUpdate(tx, db, sqlstr, args...)
	}
	if tx == nil && db == nil {
		return 0, errInit
	}
	defer util.Recover(&err)
	var rs sql.Result
	if tx != nil {
		return executeUpdate(tx, db, sqlstr, args...)
	} else {
		rs, err = se.Exec(db, sqlstr, args...)
	}
	if err == nil {
		return rs.RowsAffected()
	}
	return
}
