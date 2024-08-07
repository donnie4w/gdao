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

const VERSION = "1.1.0"

type TableBase[T any] interface {
	ClassName()
}

type DBhandle interface {
	GetTransaction() (r Transaction, err error)

	ExecuteQueryBeans(sql string, args ...any) ([]*DataBean, error)

	ExecuteQueryBean(sql string, args ...any) (*DataBean, error)

	ExecuteUpdate(sql string, args ...any) (int64, error)

	ExecuteBatch(sql string, args [][]any) (r []int64, err error)

	GetDBType() DBType

	GetDB() *sql.DB
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
