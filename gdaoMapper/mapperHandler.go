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
	"github.com/donnie4w/gdao"
	. "github.com/donnie4w/gdao/base"
	"github.com/donnie4w/gdao/gdaoCache"
)

type mapperHandler struct {
	transaction Transaction
	dBhandle    DBhandle
}

func newMapperHandler() *mapperHandler {
	if !mapperparser.hasMapper() {
		panic("The mapping file is not parsed; call gdaoMapper.build() first")
	}
	return &mapperHandler{}
}

func newMapperHandlerWithMapperparser() *mapperHandler {
	return &mapperHandler{}
}

func (t *mapperHandler) UseDBhandle(dbhandler DBhandle) {
	t.dBhandle = dbhandler
}

func (t *mapperHandler) UseDBhandleWithDB(db *sql.DB, dbType DBType) {
	t.dBhandle = gdao.NewDBHandle(db, dbType)
}

func (t *mapperHandler) IsAutocommit() bool {
	return t.transaction == nil
}

func (t *mapperHandler) SetAutocommit(autocommit bool) (err error) {
	if autocommit {
		if dbHandle := t.getDBhandle("", "", false); dbHandle != nil {
			t.transaction, err = dbHandle.GetTransaction()
		} else {
			fmt.Errorf("no data source was found")
		}
	} else {
		t.transaction = nil
	}
	return
}

func (t *mapperHandler) UseTransaction(tx Transaction) {
	t.transaction = tx
}

func (t *mapperHandler) Rollback() (err error) {
	if t.transaction != nil {
		err = t.transaction.Rollback()
		t.transaction = nil
	}
	return
}

func (t *mapperHandler) Commit() (err error) {
	if t.transaction != nil {
		err = t.transaction.Commit()
		t.transaction = nil
	}
	return
}

func (t *mapperHandler) getDBhandle(namespace, id string, queryType bool) (dbhandle DBhandle) {
	if t.dBhandle != nil {
		return t.dBhandle
	}
	if dbhandle = GetMapperDBhandle(namespace, id, queryType); dbhandle == nil {
		dbhandle = gdao.GetDefaultDBHandle()
	}
	return
}

func (t *mapperHandler) SelectBeanDirect(mapperId string, args ...any) (r *DataBean) {
	var pb *paramBean
	var err error
	if pb, _, err = t.parseParameter(mapperId, nil); err != nil {
		r = &DataBean{}
		r.SetError(err)
		return
	}
	if Logger.IsVaild {
		Logger.Debug("[Mapper Id] "+mapperId+" \nSelectBeanDirect SQL["+pb.sql+"]ARGS", args)
	}
	return t._selectBean(mapperId, pb, args...)
}

func (t *mapperHandler) SelectBean(mapperId string, parameter any) (r *DataBean) {
	var pb *paramBean
	var args []any
	var err error
	if pb, args, err = t.parseParameter(mapperId, parameter); err != nil {
		r = &DataBean{}
		r.SetError(err)
		return
	}
	if Logger.IsVaild {
		Logger.Debug("[Mapper Id] "+mapperId+" \nSelectBean SQL["+pb.sql+"]ARGS", args)
	}
	return t._selectBean(mapperId, pb, args...)
}

func (t *mapperHandler) _selectBean(mapperId string, pb *paramBean, args ...any) (r *DataBean) {
	domain := gdaoCache.GetMapperDomain(pb.namespace, pb.id)
	isCache := domain != ""
	var condition *gdaoCache.Condition
	if isCache {
		condition = gdaoCache.NewCondition("*DataBean", pb.sql, args...)
		if result := gdaoCache.GetMapperCache(domain, pb.namespace, pb.id, condition); result != nil {
			if Logger.IsVaild {
				Logger.Debug("[GET CACHE]["+pb.sql+"]", args)
			}
			return result.(*DataBean)
		}
	}
	if r = t.getDBhandle(pb.namespace, pb.id, true).ExecuteQueryBean(pb.sql, args...); r.GetError() == nil {
		if isCache {
			gdaoCache.SetMapperCache(domain, pb.namespace, pb.id, condition, r)
			if Logger.IsVaild {
				Logger.Debug("[SET CACHE]["+pb.sql+"]", args)
			}
		}
	}
	return
}

func (t *mapperHandler) SelectBeansDirect(mapperId string, args ...any) *DataBeans {
	var pb *paramBean
	var err error
	if pb, _, err = t.parseParameter(mapperId, nil); err != nil {
		r := &DataBeans{}
		r.SetError(err)
		return r
	}
	if Logger.IsVaild {
		Logger.Debug("[Mapper Id] "+mapperId+" \nSelectsBean SQL["+pb.sql+"]ARGS", args)
	}
	return t._selectBeans(mapperId, pb, args...)
}

func (t *mapperHandler) SelectBeans(mapperId string, parameter any) *DataBeans {
	var pb *paramBean
	var args []any
	var err error
	if pb, args, err = t.parseParameter(mapperId, parameter); err != nil {
		r := &DataBeans{}
		r.SetError(err)
		return r
	}
	if Logger.IsVaild {
		Logger.Debug("[Mapper Id] "+mapperId+" \nSelectBeans SQL["+pb.sql+"]ARGS", args)
	}
	return t._selectBeans(mapperId, pb, args...)
}

func (t *mapperHandler) _selectBeans(mapperId string, pb *paramBean, args ...any) (r *DataBeans) {
	domain := gdaoCache.GetMapperDomain(pb.namespace, pb.id)
	isCache := domain != ""
	var condition *gdaoCache.Condition
	if isCache {
		condition = gdaoCache.NewCondition("[]*DataBean", pb.sql, args...)
		if result := gdaoCache.GetMapperCache(domain, pb.namespace, pb.id, condition); result != nil {
			if Logger.IsVaild {
				Logger.Debug("[GET CACHE]["+pb.sql+"]", args)
			}
			return result.(*DataBeans)
		}
	}
	if r = t.getDBhandle(pb.namespace, pb.id, true).ExecuteQueryBeans(pb.sql, args...); r.GetError() == nil && r.Len() > 0 {
		if isCache {
			gdaoCache.SetMapperCache(domain, pb.namespace, pb.id, condition, r)
			if Logger.IsVaild {
				Logger.Debug("[SET CACHE]["+pb.sql+"]", args)
			}
		}
	}
	return
}

func (t *mapperHandler) InsertDirect(mapperId string, args ...any) (r int64, err error) {
	var pb *paramBean
	if pb, _, err = t.parseParameter(mapperId, nil); err != nil {
		return r, err
	}
	if Logger.IsVaild {
		Logger.Debug("[Mapper Id] "+mapperId+" \nInsertDirect SQL["+pb.sql+"]ARGS", args)
	}
	return t.getDBhandle(pb.namespace, pb.id, false).ExecuteUpdate(pb.sql, args...)
}

func (t *mapperHandler) Insert(mapperId string, parameter any) (r int64, err error) {
	var pb *paramBean
	var args []any
	if pb, args, err = t.parseParameter(mapperId, parameter); err != nil {
		return r, err
	}
	if Logger.IsVaild {
		Logger.Debug("[Mapper Id] "+mapperId+" \nInsert SQL["+pb.sql+"]ARGS", args)
	}
	return t.getDBhandle(pb.namespace, pb.id, false).ExecuteUpdate(pb.sql, args...)
}

func (t *mapperHandler) UpdateDirect(mapperId string, args ...any) (r int64, err error) {
	var pb *paramBean
	if pb, _, err = t.parseParameter(mapperId, nil); err != nil {
		return r, err
	}
	if Logger.IsVaild {
		Logger.Debug("[Mapper Id] "+mapperId+" \nUpdateDirect SQL["+pb.sql+"]ARGS", args)
	}
	return t.getDBhandle(pb.namespace, pb.id, false).ExecuteUpdate(pb.sql, args...)
}

func (t *mapperHandler) Update(mapperId string, parameter any) (r int64, err error) {
	var pb *paramBean
	var args []any
	if pb, args, err = t.parseParameter(mapperId, parameter); err != nil {
		return r, err
	}
	if Logger.IsVaild {
		Logger.Debug("[Mapper Id] "+mapperId+" \nUpdate SQL["+pb.sql+"]ARGS", args)
	}
	return t.getDBhandle(pb.namespace, pb.id, false).ExecuteUpdate(pb.sql, args...)
}

func (t *mapperHandler) DeleteDirect(mapperId string, args ...any) (r int64, err error) {
	var pb *paramBean
	if pb, _, err = t.parseParameter(mapperId, nil); err != nil {
		return r, err
	}
	if Logger.IsVaild {
		Logger.Debug("[Mapper Id] "+mapperId+" \nDeleteDirect SQL["+pb.sql+"]ARGS", args)
	}
	return t.getDBhandle(pb.namespace, pb.id, false).ExecuteUpdate(pb.sql, args...)
}

func (t *mapperHandler) Delete(mapperId string, parameter any) (r int64, err error) {
	var pb *paramBean
	var args []any
	if pb, args, err = t.parseParameter(mapperId, parameter); err != nil {
		return r, err
	}
	if Logger.IsVaild {
		Logger.Debug("[Mapper Id] "+mapperId+" \nDelete SQL["+pb.sql+"]ARGS", args)
	}
	return t.getDBhandle(pb.namespace, pb.id, false).ExecuteUpdate(pb.sql, args...)
}

func (t *mapperHandler) parseParameter(mapperId string, parameter any) (pb *paramBean, args []any, err error) {
	var ok bool
	if pb, ok = mapperparser.getParamBean(mapperId); !ok {
		return nil, nil, fmt.Errorf("Mapper Id not found [%s]", mapperId)
	}
	if parameter != nil {
		args, err = pb.setParameter(parameter)
	}
	return
}

var defaultMapperHandler *mapperHandler

// NewInstance create a GdaoMapper Object
func NewInstance() GdaoMapper {
	return newMapperHandler()
}

func init() {
	defaultMapperHandler = newMapperHandlerWithMapperparser()
	IsAutocommit = defaultMapperHandler.IsAutocommit
	SetAutocommit = defaultMapperHandler.SetAutocommit
	UseTransaction = defaultMapperHandler.UseTransaction
	Rollback = defaultMapperHandler.Rollback
	Commit = defaultMapperHandler.Commit
	UseDBhandle = defaultMapperHandler.UseDBhandle
	UseDBhandleWithDB = defaultMapperHandler.UseDBhandleWithDB

	SelectBean = defaultMapperHandler.SelectBean
	SelectBeans = defaultMapperHandler.SelectBeans

	SelectBeanDirect = defaultMapperHandler.SelectBeanDirect
	SelectBeansDirect = defaultMapperHandler.SelectBeansDirect
	InsertDirect = defaultMapperHandler.InsertDirect
	UpdateDirect = defaultMapperHandler.UpdateDirect
	DeleteDirect = defaultMapperHandler.DeleteDirect

	Insert = defaultMapperHandler.Insert
	Update = defaultMapperHandler.Update
	Delete = defaultMapperHandler.Delete
}
