// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package base

import (
	"database/sql"
	"github.com/donnie4w/gofer/util"
)

type TableBase interface {
	TableName() string
}

type DBhandle interface {
	GetTransaction() (r Transaction, err error)

	ExecuteQueryBean(sql string, args ...any) *DataBean

	ExecuteQueryBeans(sql string, args ...any) *DataBeans

	ExecuteUpdate(sql string, args ...any) (sql.Result, error)

	ExecuteBatch(sql string, args [][]any) (r []sql.Result, err error)

	GetDBType() DBType

	GetDB() *sql.DB

	Close() error
}

type Scanner interface {
	Scan(fieldname string, value any)

	// ToGdao
	// : when don't create an object by calling a New method of the standardized entity class,
	// but by using some other method such as the new keyword, then should call the ToGdao function,
	// which initializes the relevant data for database operations
	ToGdao()
}

var MapperPre = string(util.Base58EncodeForInt64(uint64(util.RandId())))

type In struct {
	Value any
}

func NewIn(v any) In {
	return In{Value: v}
}

func (t In) GetValue() any {
	return t.Value
}

func (t In) SetValue(v any) {
}

type Out struct {
	Value any
}

func NewOut(v any) Out {
	return Out{Value: v}
}

func (t Out) GetValue() *sql.Out {
	return &sql.Out{Dest: t.Value}
}

type InOut struct {
	Value any
}

func NewInOut(v any) InOut {
	return InOut{Value: v}
}

func (t InOut) GetValue() *sql.Out {
	return &sql.Out{Dest: t.Value, In: true}
}

var (
	GetMapperIds      func(string) []string
	HasMapperId       func(string) bool
	GetMapperDBhandle func(string, string, bool) DBhandle
)
