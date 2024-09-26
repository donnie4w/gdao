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
	goutil "github.com/donnie4w/gofer/util"
	"sync"
	"sync/atomic"
)

var sqlWare = hashmap.NewLimitHashMap[uint64, *int64](1 << 19)
var stmtExec = &stmtexec{stmtMap: hashmap.NewMap[*sql.DB, *hashmap.MapL[uint64, *sql.Stmt]](), mux: &sync.Mutex{}}
var errorStmt = errors.New("")

type stmtexec struct {
	stmtMap *hashmap.Map[*sql.DB, *hashmap.MapL[uint64, *sql.Stmt]]
	mux     *sync.Mutex
	lock    int64
}

func (se *stmtexec) Exec(db *sql.DB, sqlStr string, args ...any) (rs sql.Result, err error) {
	if stmt, e := se.Prepare(sqlStr, db); e == nil {
		return stmt.Exec(args...)
	} else {
		return db.Exec(sqlStr, args...)
	}
}

func (se *stmtexec) Qurey(db *sql.DB, sqlStr string, args ...any) (rs *sql.Rows, err error) {
	if stmt, e := se.Prepare(sqlStr, db); e == nil {
		return stmt.Query(args...)
	} else {
		return db.Query(sqlStr, args...)
	}
}

func (se *stmtexec) clear(db *sql.DB) {
	if atomic.CompareAndSwapInt64(&se.lock, 0, stmtLimit) {
		defer atomic.StoreInt64(&se.lock, 0)
		if sm, _ := se.stmtMap.Get(db); sm != nil {
			if sm.Len() >= stmtLimit {
				sm.Range(func(k uint64, v *sql.Stmt) bool {
					sm.Del(k)
					v.Close()
					return true
				})
			}
		}
	}
}

func (se *stmtexec) newmap(db *sql.DB) (r *hashmap.MapL[uint64, *sql.Stmt]) {
	se.mux.Lock()
	defer se.mux.Unlock()
	if !se.stmtMap.Has(db) {
		r = hashmap.NewMapL[uint64, *sql.Stmt]()
		se.stmtMap.Put(db, r)
	}
	return r
}

func (se *stmtexec) Prepare(sqlStr string, db *sql.DB) (stmt *sql.Stmt, err error) {
	if se.len(db) >= stmtLimit {
		se.clear(db)
		return stmt, errorStmt
	}
	var hm *hashmap.MapL[uint64, *sql.Stmt]
	var sqlhs uint64
	if hm, _ = se.stmtMap.Get(db); hm != nil {
		sqlhs := goutil.Hash64([]byte(sqlStr))
		if a, b := hm.Get(sqlhs); b {
			return a, nil
		}
	} else {
		hm = se.newmap(db)
	}
	if stmt, err = db.Prepare(sqlStr); err == nil {
		if sqlhs == 0 {
			sqlhs = goutil.Hash64([]byte(sqlStr))
		}
		if p, ok := hm.Put(sqlhs, stmt); ok && p != nil {
			p.Close()
		}
	}
	return
}

func (se *stmtexec) len(db *sql.DB) int64 {
	if se.lock > 0 {
		return se.lock
	}
	if v, _ := se.stmtMap.Get(db); v != nil {
		return v.Len()
	}
	return 0
}

func (se *stmtexec) nostmt(tx *sql.Tx, db *sql.DB, sqlstr string) (b bool) {
	if stmtLimit == 0 || tx != nil || se.len(db) >= stmtLimit {
		return true
	}
	sqlhs := goutil.Hash64([]byte(sqlstr))
	if v, ok := sqlWare.Get(sqlhs); ok {
		b = atomic.AddInt64(v, 1) < 16
	} else {
		sqlWare.Put(sqlhs, new(int64))
		b = true
	}
	return
}

func (se *stmtexec) executeQueryBeans(tx *sql.Tx, db *sql.DB, sqlstr string, args ...any) (databases []*DataBean, err error) {
	if se.nostmt(tx, db, sqlstr) {
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
	if se.nostmt(tx, db, sqlstr) {
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

func (se *stmtexec) executeUpdate(tx *sql.Tx, db *sql.DB, sqlstr string, args ...any) (rs sql.Result, err error) {
	if se.nostmt(tx, db, sqlstr) {
		return executeUpdate(tx, db, sqlstr, args...)
	}
	if tx == nil && db == nil {
		return nil, errInit
	}
	defer util.Recover(&err)
	if tx != nil {
		return executeUpdate(tx, db, sqlstr, args...)
	} else {
		return se.Exec(db, sqlstr, args...)
	}
}
