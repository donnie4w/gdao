// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdao

import (
	"database/sql"
	"fmt"
	"github.com/donnie4w/gdao/gdaoSlave"

	. "github.com/donnie4w/gdao/base"
	"github.com/donnie4w/gofer/hashmap"
)

const VERSION = "1.1.0"

var errInit = fmt.Errorf("the gdao DataSource was not initialized(Hint: gdao.Init(db, dbtype))")

type GStruct[P any, T any] interface {
	Scanner
	TableBase[T]
	UseCache(use bool)
	UseTransaction(transaction Transaction)
	UseDBHandler(db DBhandle) *Table[T]
	UseCommentLine(commentline string)
	MustMaster(must bool)
	Where(wheres ...*Where[T]) *Table[T]
	OrderBy(sorts ...*Sort[T]) *Table[T]
	GroupBy(columns ...Column[T]) *Table[T]
	Having(havings ...*Having[T]) *Table[T]
	Limit2(offset, limit int64)
	Limit(limit int64)
	Selects(columns ...Column[T]) (_r []P, err error)
	Select(columns ...Column[T]) (_r P, err error)
	Update() (int64, error)
	Insert() (int64, error)
	Delete() (int64, error)
	AddBatch()
	ExecBatch() ([]int64, error)
	Copy(h P) P
	Encode() ([]byte, error)
	Decode(bs []byte) (err error)
	String() string
	TABLENAME() string
}

func NewDBHandler(db *sql.DB, dbtype DBType) DBhandle {
	return newdbhandle(db, dbtype)
}

func GetDefaultDBHandle() DBhandle {
	return defaultDBhandle
}

func ExecuteQuery[T any](sql string, args ...any) (r *T, err error) {
	if databean, err := defaultDBhandle.ExecuteQueryBean(sql, args...); err == nil {
		return Scan[T](databean)
	} else {
		return nil, err
	}
}

func ExecuteQueryList[T any](sql string, args ...any) (r []*T, err error) {
	var databeans []*DataBean
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

func ExecuteQueryBean(sql string, args ...any) (*DataBean, error) {
	if defaultDBhandle == nil {
		return nil, errInit
	}
	return defaultDBhandle.ExecuteQueryBean(sql, args...)
}

func ExecuteQueryBeans(sql string, args ...any) ([]*DataBean, error) {
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

func getDBhandle(classname, tableName string, queryType bool) (r DBhandle) {
	if gdaoSlave.Len() > 0 && queryType {
		if r = gdaoSlave.Get(classname, tableName, ""); r != nil {
			return
		}
	}
	if handleMap.Len() > 0 {
		if h, ok := handleMap.Get(classname); ok {
			return h
		}
		if h, ok := handleMap.Get(tableName); ok {
			return h
		}
	}
	return defaultDBhandle
}

var defaultDBhandle DBhandle

var handleMap = hashmap.MapL[string, DBhandle]{}

func Init(db *sql.DB, dbtype DBType) {
	defaultDBhandle = newdbhandle(db, dbtype)
}

func SetDataSource(tableName string, db *sql.DB, dbtype DBType) {
	handleMap.Put(tableName, newdbhandle(db, dbtype))
}

func SetDataSourceWithClass[T TableBase[T]](db *sql.DB, dbtype DBType) {
	handleMap.Put(Classname[T](), newdbhandle(db, dbtype))
}

func RemoveDataSource(tableName string) {
	handleMap.Del(tableName)
}

func RemoveDataSourceWithClass[T TableBase[T]]() {
	handleMap.Del(Classname[T]())
}

func NewTransaction() (r Transaction, err error) {
	return newTX(defaultDBhandle)
}

func NewTransactionWithDBhandle(db DBhandle) (r Transaction, err error) {
	return newTX(db)
}

type Scanner interface {
	Scan(fieldname string, value any)
	ToGdao()
}

func SetLogger(on bool) {
	Logger.SetLogger(on)
}
