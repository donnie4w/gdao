// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoMapper

import (
	"database/sql"
	"fmt"
	"github.com/donnie4w/gdao/base"
)

type JdaoMapper interface {
	IsAutocommit() bool
	SetAutocommit(autocommit bool) (err error)
	UseTransaction(tx base.Transaction)
	Rollback() (err error)
	Commit() (err error)
	SetDBhandle(dbhandler base.DBhandle)
	SetDBhandleWithDB(db *sql.DB, dbType base.DBType)

	SelectBean(mapperId string, args ...any) (*base.DataBean, error)
	SelectsBean(mapperId string, args ...any) ([]*base.DataBean, error)
	Insert(mapperId string, args ...any) (int64, error)
	Update(mapperId string, args ...any) (int64, error)
	Delete(mapperId string, args ...any) (int64, error)

	SelectsAny(mapperId string, parameter any) (r []*base.DataBean, err error)
	SelectAny(mapperId string, parameter any) (r *base.DataBean, err error)
	InsertAny(mapperId string, parameter any) (int64, error)
	UpdateAny(mapperId string, parameter any) (int64, error)
	DeleteAny(mapperId string, parameter any) (int64, error)
}

var (
	IsAutocommit      = defaultMapperHandler.IsAutocommit
	SetAutocommit     = defaultMapperHandler.SetAutocommit
	UseTransaction    = defaultMapperHandler.UseTransaction
	Rollback          = defaultMapperHandler.Rollback
	Commit            = defaultMapperHandler.Commit
	SetDBhandle       = defaultMapperHandler.SetDBhandle
	SetDBhandleWithDB = defaultMapperHandler.SetDBhandleWithDB

	SelectBean  = defaultMapperHandler.SelectBean
	SelectsBean = defaultMapperHandler.SelectsBean
	Insert      = defaultMapperHandler.Insert
	Update      = defaultMapperHandler.Update
	Delete      = defaultMapperHandler.Delete

	InsertAny = defaultMapperHandler.InsertAny
	UpdateAny = defaultMapperHandler.UpdateAny
	DeleteAny = defaultMapperHandler.DeleteAny
)

func Select[T any](mapperId string, args ...any) (*T, error) {
	return (*mapperInvoke[T])(defaultMapperHandler).Select(mapperId, args...)
}

func SelectAny[T any](mapperId string, parameter any) (*T, error) {
	return (*mapperInvoke[T])(defaultMapperHandler).SelectAny(mapperId, parameter)
}

func Selects[T any](mapperId string, args ...any) ([]*T, error) {
	return (*mapperInvoke[T])(defaultMapperHandler).Selects(mapperId, args...)
}

func SelectsAny[T any](mapperId string, parameter any) ([]*T, error) {
	return (*mapperInvoke[T])(defaultMapperHandler).SelectsAny(mapperId, parameter)
}

func SelectWithJdaoMapper[T any](jdaomapper JdaoMapper, mapperId string, args ...any) (*T, error) {
	if v, ok := jdaomapper.(*mapperHandler); ok {
		return (*mapperInvoke[T])(v).Select(mapperId, args...)
	}
	return nil, fmt.Errorf("jdaomapper is not a MapperHandler pointer")
}

func SelectAnyWithJdaoMapper[T any](jdaomapper JdaoMapper, mapperId string, parameter any) (*T, error) {
	if v, ok := jdaomapper.(*mapperHandler); ok {
		return (*mapperInvoke[T])(v).SelectAny(mapperId, parameter)
	}
	return nil, fmt.Errorf("jdaomapper is not a MapperHandler pointer")
}

func SelectsWithJdaoMapper[T any](jdaomapper JdaoMapper, mapperId string, args ...any) ([]*T, error) {
	if v, ok := jdaomapper.(*mapperHandler); ok {
		return (*mapperInvoke[T])(v).Selects(mapperId, args...)
	}
	return nil, fmt.Errorf("jdaomapper is not a MapperHandler pointer")
}

func SelectsAnyWithJdaoMapper[T any](jdaomapper JdaoMapper, mapperId string, parameter any) ([]*T, error) {
	if v, ok := jdaomapper.(*mapperHandler); ok {
		return (*mapperInvoke[T])(v).SelectsAny(mapperId, parameter)
	}
	return nil, fmt.Errorf("jdaomapper is not a MapperHandler pointer")
}
