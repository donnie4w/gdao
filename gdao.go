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

type GStruct[P any, T any] interface {
	Scanner
	// UseCache use default gdaoCache or not use cache
	UseCache(use bool)

	// UseTransaction use specified transaction
	UseTransaction(transaction Transaction)

	// UseDBHandle use specified DBhandle
	UseDBHandle(db DBhandle) *Table[T]

	// UseCommentLine set annotations for sql
	UseCommentLine(commentline string)

	MustMaster(must bool)

	// Where adds a WHERE clause to the query with one or more conditions.
	//
	// Parameters:
	//
	//	wheres: Variable length argument list of *Where[T] objects representing the conditions to add to the WHERE clause.
	//
	// Returns:
	//
	//	A pointer to the Table[T] instance to allow method chaining.
	//
	// Description:
	//
	//	This function allows you to specify one or more conditions that will be added to the WHERE clause of the SQL query.
	//	Each *Where[T] object represents a condition that must be satisfied by the rows returned by the query.
	//	Multiple conditions can be combined to form complex queries.
	//
	// Example:
	//
	//	// Assuming "hs" is an instance of a Table struct that represents a table named "hstest"
	//	// And "Rowname" and "Id" are columns in the "hstest" table
	//	hs := dao.NewHstest()
	//	hs = hs.Where(hs.Rowname.RLIKE(1)).GroupBy(hs.Id).Having(hs.Id.Count().LT(2)).Limit(2)
	//	hslist, _ := hs.Selects()
	Where(wheres ...*Where[T]) *Table[T]
	// OrderBy sql: order by
	OrderBy(sorts ...*Sort[T]) *Table[T]
	// GroupBy sql: group by
	GroupBy(columns ...Column[T]) *Table[T]
	// Having sql: having
	Having(havings ...*Having[T]) *Table[T]
	Limit2(offset, limit int64)
	Limit(limit int64)
	// Selects sql:select from table and Return data slice
	Selects(columns ...Column[T]) (_r []P, err error)
	// Select sql:select from table and Return first data
	Select(columns ...Column[T]) (_r P, err error)
	// Update sql: update
	Update() (sql.Result, error)
	// Insert sql: insert
	Insert() (sql.Result, error)
	// Delete sql: delete
	Delete() (sql.Result, error)
	// AddBatch sql: add data to batch sql
	AddBatch()
	// ExecBatch sql:database batch operation
	ExecBatch() ([]sql.Result, error)
	//Copy object data
	Copy(h P) P
	// Encode Serialized object
	Encode() ([]byte, error)
	// Decode deserialization
	Decode(bs []byte) (err error)
	String() string
	// TableName return table name
	TableName() string
}
