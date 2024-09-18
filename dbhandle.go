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

// ExecuteQuery executes an SQL query and returns the parsed results.
// T is a generic type that represents the specific struct or type the query result should be converted to.
// sql is the SQL query statement to execute.
// args is an optional list of parameters to substitute placeholders in the SQL query.
// The function returns a pointer to a value of type *T, which is typically a pointer to a struct that holds the query results.
// If there's an error, it returns nil and the specific error information; otherwise, it returns a filled result object and nil.
func ExecuteQuery[T any](sql string, args ...any) (r *T, err error) {
	if databean := defaultDBhandle.ExecuteQueryBean(sql, args...); databean.GetError() == nil && databean.Len() > 0 {
		r = new(T)
		err = databean.ScanAndFree(r)
	} else {
		err = databean.GetError()
	}
	return
}

// ExecuteQueryList executes an SQL query and returns a list of parsed results.
// T is a generic type that represents the specific struct or type the query results should be converted to.
// sql is the SQL query statement to execute.
// args is an optional list of parameters to substitute placeholders in the SQL query.
// The function returns a slice of pointers to values of type T, where each element represents one row of the query results.
// If there's an error, it returns nil and the specific error information; otherwise, it returns a slice of filled result objects and nil.
func ExecuteQueryList[T any](sql string, args ...any) (r []*T, err error) {
	if databeans := defaultDBhandle.ExecuteQueryBeans(sql, args...); databeans.GetError() == nil && databeans.Len() > 0 {
		r = make([]*T, 0)
		for _, databean := range databeans.Beans {
			t := new(T)
			if err = databean.ScanAndFree(t); err == nil {
				r = append(r, t)
			} else {
				break
			}
		}
	} else {
		return nil, databeans.GetError()
	}
	return
}

// ExecuteQueryBean executes an SQL query and returns a single DataBean object.
// sql is the SQL query statement to execute.
// args is an optional list of parameters to substitute placeholders in the SQL query.
// The function returns a pointer to a DataBean object, which typically holds the data retrieved from a single row in the query results.
// If there's an error, it returns nil and the specific error information; otherwise, it returns a filled DataBean object and nil.
func ExecuteQueryBean(sql string, args ...any) *base.DataBean {
	if defaultDBhandle == nil {
		r := &base.DataBean{}
		r.SetError(errInit)
		return r
	}
	return defaultDBhandle.ExecuteQueryBean(sql, args...)
}

// ExecuteQueryBeans executes an SQL query and returns a list of DataBean objects.
// sql is the SQL query statement to execute.
// args is an optional list of parameters to substitute placeholders in the SQL query.
// The function returns a slice of pointers to DataBean objects, where each element represents one row of the query results.
// If there's an error, it returns nil and the specific error information; otherwise, it returns a slice of filled DataBean objects and nil.
func ExecuteQueryBeans(sql string, args ...any) *base.DataBeans {
	if defaultDBhandle == nil {
		r := &base.DataBeans{}
		r.SetError(errInit)
		return r
	}
	return defaultDBhandle.ExecuteQueryBeans(sql, args...)
}

// ExecuteUpdate executes an SQL update, insert, or delete statement.
// sql is the SQL statement to execute.
// args is an optional list of parameters to substitute placeholders in the SQL statement.
// The function returns the number of rows affected by the SQL statement and any error encountered.
// If there's an error, it returns -1 and the specific error information; otherwise, it returns the number of affected rows and nil.
func ExecuteUpdate(sql string, args ...any) (sql.Result, error) {
	if defaultDBhandle == nil {
		return nil, errInit
	}
	return defaultDBhandle.ExecuteUpdate(sql, args...)
}

// ExecuteBatch executes a batch of SQL statements.
// sql is the SQL statement to execute for each batch item.
// args is a slice of slices, where each inner slice contains the arguments for a single SQL statement.
// The function returns a slice of int64 values representing the number of rows affected by each SQL statement and any error encountered.
// If there's an error, it returns nil and the specific error information; otherwise, it returns the slice of affected rows and nil.
func ExecuteBatch(sql string, args [][]any) ([]sql.Result, error) {
	if defaultDBhandle == nil {
		return nil, errInit
	}
	return defaultDBhandle.ExecuteBatch(sql, args)
}
