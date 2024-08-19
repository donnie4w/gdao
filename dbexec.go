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
	"github.com/donnie4w/gdao/util"
	"github.com/donnie4w/gofer/pool/buffer"
)

var bufpool = buffer.NewPool[[]any](func() *[]any {
	return nil
}, func(a *[]any) {
	*a = (*a)[:0]
})

func newAnys(length int) (r *[]any) {
	if r = bufpool.Get(); r == nil {
		var v []any
		return &v
	} else {
		return r
	}
}

func executeQueryBeans(tx *sql.Tx, db *sql.DB, sqlstr string, args ...any) (databases []*DataBean, err error) {
	if tx == nil && db == nil {
		return nil, errInit
	}
	//defer util.Recover(&err)
	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(sqlstr, args...)
	} else {
		rows, err = db.Query(sqlstr, args...)
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

func executeQueryBean(tx *sql.Tx, db *sql.DB, sqlstr string, args ...any) (dataBean *DataBean, err error) {
	if tx == nil && db == nil {
		return nil, errInit
	}
	//defer util.Recover(&err)
	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(sqlstr, args...)
	} else {
		rows, err = db.Query(sqlstr, args...)
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

func executeUpdate(tx *sql.Tx, db *sql.DB, sqlstr string, args ...any) (r int64, err error) {
	if tx == nil && db == nil {
		return 0, errInit
	}
	defer util.Recover(&err)
	var stmtIns *sql.Stmt
	if tx != nil {
		stmtIns, err = tx.Prepare(sqlstr)
	} else {
		stmtIns, err = db.Prepare(sqlstr)
	}
	if err != nil {
		return
	}
	defer stmtIns.Close()
	var rs sql.Result
	if rs, err = stmtIns.Exec(args...); err == nil {
		return rs.RowsAffected()
	}
	return
}

func executeBatch(tx *sql.Tx, db *sql.DB, sqlstr string, args [][]any) (r []int64, err error) {
	if tx == nil && db == nil {
		return nil, errInit
	}
	defer util.Recover(&err)
	var stmtIns *sql.Stmt
	if tx != nil {
		stmtIns, err = tx.Prepare(sqlstr)
	} else {
		stmtIns, err = db.Prepare(sqlstr)
	}
	if err != nil {
		return nil, err
	}
	defer stmtIns.Close()
	r = make([]int64, 0)
	for _, record := range args {
		if rs, er := stmtIns.Exec(record...); er == nil {
			i, _ := rs.RowsAffected()
			r = append(r, i)
		} else {
			err = er
			break
		}
	}
	return
}
