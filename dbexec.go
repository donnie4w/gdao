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
)

func executeQueryBeans(tx *sql.Tx, db *sql.DB, sqlstr string, args ...any) (databases []*DataBean, err error) {
	if tx == nil && db == nil {
		return nil, errInit
	}
	defer util.Recover(&err)
	var stmt *sql.Stmt
	if tx != nil {
		stmt, err = tx.Prepare(sqlstr)
	} else {
		stmt, err = db.Prepare(sqlstr)
	}
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var rows *sql.Rows
	if rows, err = stmt.Query(args...); err == nil {
		defer rows.Close()
		databases = make([]*DataBean, 0)
		if types, er := rows.ColumnTypes(); er == nil {
			for rows.Next() {
				databean := NewDataBean()
				buff := make([]any, 0, len(types))
				for i, columntype := range types {
					columntype.DatabaseTypeName()
					fb := new(FieldBeen)
					fb.FieldName = columntype.Name()
					fb.FieldIndex = i
					buff = append(buff, &fb.FieldValue)
					databean.Put(columntype.Name(), i, fb)
				}
				if err = rows.Scan(buff...); err != nil {
					databean.SetError(err)
					return
				}
				databases = append(databases, databean)
			}
		} else {
			err = er
		}
	}
	return
}

func executeQueryBean(tx *sql.Tx, db *sql.DB, sqlstr string, args ...any) (dataBean *DataBean, err error) {
	if tx == nil && db == nil {
		return nil, errInit
	}
	defer util.Recover(&err)
	var stmt *sql.Stmt
	if tx != nil {
		stmt, err = tx.Prepare(sqlstr)
	} else {
		stmt, err = db.Prepare(sqlstr)
	}
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var rows *sql.Rows
	if rows, err = stmt.Query(args...); err == nil {
		defer rows.Close()
		if types, er := rows.ColumnTypes(); er == nil {
			dataBean = NewDataBean()
			if rows.Next() {
				buff := make([]any, 0, len(types))
				for i, columntype := range types {
					fb := new(FieldBeen)
					fb.FieldName = columntype.Name()
					fb.FieldIndex = i
					buff = append(buff, &fb.FieldValue)
					dataBean.Put(columntype.Name(), i, fb)
				}
				if err = rows.Scan(buff...); err != nil {
					dataBean.SetError(err)
					return
				}
			}
		} else {
			err = er
		}
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
