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
	"github.com/donnie4w/gdao/gdaoSlave"
)

type mapperHandler struct {
	transaction Transaction
	dBhandle    DBhandle
}

func newMapperHandler() *mapperHandler {
	if !mapperparser.hasMapper() {
		panic("The mapping file is not parsed.call JdaoMapper.build(.xml) first")
	}
	return &mapperHandler{}
}

func (t *mapperHandler) SetDBhandle(dbhandler DBhandle) {
	t.dBhandle = dbhandler
}

func (t *mapperHandler) SetDBhandleWithDB(db *sql.DB, dbType DBType) {
	t.dBhandle = gdao.NewDBHandler(db, dbType)
}

func (t *mapperHandler) IsAutocommit() bool {
	return t.transaction == nil
}

func (t *mapperHandler) SetAutocommit(autocommit bool) (err error) {
	if autocommit {
		if dbHandle := t.getDBhandle("", false); dbHandle != nil {
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

func (t *mapperHandler) getDBhandle(mapperId string, queryType bool) (dbhandle DBhandle) {
	if t.dBhandle != nil {
		return t.dBhandle
	}
	if mapperId != "" && queryType && gdaoSlave.Len() > 0 {
		dbhandle = gdaoSlave.Get("", "", mapperId)
	}
	if dbhandle == nil {
		dbhandle = gdao.GetDefaultDBHandle()
	}
	return
}

func (t *mapperHandler) SelectBean(mapperId string, args ...any) (r *DataBean, err error) {
	var pb *paramBean
	if pb, _, err = t.parseParameter(mapperId, nil); err != nil {
		return r, err
	}
	return t._selectBean(mapperId, pb, args...)
}

func (t *mapperHandler) SelectAny(mapperId string, parameter any) (r *DataBean, err error) {
	var pb *paramBean
	var args []any
	if pb, args, err = t.parseParameter(mapperId, parameter); err != nil {
		return r, err
	}
	return t._selectBean(mapperId, pb, args...)
}

func (t *mapperHandler) _selectBean(mapperId string, pb *paramBean, args ...any) (r *DataBean, err error) {
	if Logger.IsVaild {
		Logger.Debug("[Mapper Id] "+mapperId+" \nSELECTONE SQL["+pb.sql+"]ARGS", args)
	}
	cacheid := cacheId(mapperId)
	domain := gdaoCache.GetDomain(cacheid)
	isCache := domain != ""
	var condition *gdaoCache.Condition
	if isCache {
		condition = gdaoCache.NewCondition("one", pb.sql, args...)
		if result := gdaoCache.GetCache(domain, cacheid, condition); result != nil {
			if Logger.IsVaild {
				Logger.Debug("[GET CACHE]["+pb.sql+"]", args)
			}
			return result.(*DataBean), nil
		}
	}
	if r, err = t.getDBhandle(mapperId, true).ExecuteQueryBean(pb.sql, args...); err == nil {
		if isCache {
			gdaoCache.SetCache(domain, cacheid, condition, r)
			if Logger.IsVaild {
				Logger.Debug("[SET CACHE]["+pb.sql+"]", args)
			}
		}
	}
	return
}

func (t *mapperHandler) SelectsBean(mapperId string, args ...any) (r []*DataBean, err error) {
	var pb *paramBean
	if pb, _, err = t.parseParameter(mapperId, nil); err != nil {
		return r, err
	}
	return t._selectsBean(mapperId, pb, args...)
}

func (t *mapperHandler) SelectsAny(mapperId string, parameter any) (r []*DataBean, err error) {
	var pb *paramBean
	var args []any
	if pb, args, err = t.parseParameter(mapperId, parameter); err != nil {
		return r, err
	}
	return t._selectsBean(mapperId, pb, args...)
}

func (t *mapperHandler) _selectsBean(mapperId string, pb *paramBean, args ...any) (r []*DataBean, err error) {
	if Logger.IsVaild {
		Logger.Debug("[Mapper Id] "+mapperId+" \nSELECTLIST SQL["+pb.sql+"]ARGS", args)
	}
	cacheid := cacheId(mapperId)
	domain := gdaoCache.GetDomain(cacheid)
	isCache := domain != ""
	var condition *gdaoCache.Condition
	if isCache {
		condition = gdaoCache.NewCondition("list", pb.sql, args...)
		if result := gdaoCache.GetCache(domain, cacheid, condition); result != nil {
			if Logger.IsVaild {
				Logger.Debug("[GET CACHE]["+pb.sql+"]", args)
			}
			return result.([]*DataBean), nil
		}
	}
	if r, err = t.getDBhandle(mapperId, true).ExecuteQueryBeans(pb.sql, args...); err == nil {
		if isCache {
			gdaoCache.SetCache(domain, cacheid, condition, r)
			if Logger.IsVaild {
				Logger.Debug("[SET CACHE]["+pb.sql+"]", args)
			}
		}
	}
	return
}

func (t *mapperHandler) Insert(mapperId string, args ...any) (r int64, err error) {
	var pb *paramBean
	if pb, _, err = t.parseParameter(mapperId, nil); err != nil {
		return r, err
	}
	if Logger.IsVaild {
		Logger.Debug("[Mapper Id] "+mapperId+" \nINSERT SQL["+pb.sql+"]ARGS", args)
	}
	return t.getDBhandle(mapperId, false).ExecuteUpdate(pb.sql, args...)
}

func (t *mapperHandler) InsertAny(mapperId string, parameter any) (r int64, err error) {
	var pb *paramBean
	var args []any
	if pb, args, err = t.parseParameter(mapperId, parameter); err != nil {
		return r, err
	}
	if Logger.IsVaild {
		Logger.Debug("[Mapper Id] "+mapperId+" \nINSERT SQL["+pb.sql+"]ARGS", args)
	}
	return t.getDBhandle(mapperId, false).ExecuteUpdate(pb.sql, args...)
}

func (t *mapperHandler) Update(mapperId string, args ...any) (r int64, err error) {
	var pb *paramBean
	if pb, _, err = t.parseParameter(mapperId, nil); err != nil {
		return r, err
	}
	if Logger.IsVaild {
		Logger.Debug("[Mapper Id] "+mapperId+" \nUPDATE SQL["+pb.sql+"]ARGS", args)
	}
	return t.getDBhandle(mapperId, false).ExecuteUpdate(pb.sql, args...)
}

func (t *mapperHandler) UpdateAny(mapperId string, parameter any) (r int64, err error) {
	var pb *paramBean
	var args []any
	if pb, args, err = t.parseParameter(mapperId, parameter); err != nil {
		return r, err
	}
	if Logger.IsVaild {
		Logger.Debug("[Mapper Id] "+mapperId+" \nUPDATE SQL["+pb.sql+"]ARGS", args)
	}
	return t.getDBhandle(mapperId, false).ExecuteUpdate(pb.sql, args...)
}

func (t *mapperHandler) Delete(mapperId string, args ...any) (r int64, err error) {
	var pb *paramBean
	if pb, _, err = t.parseParameter(mapperId, nil); err != nil {
		return r, err
	}
	if Logger.IsVaild {
		Logger.Debug("[Mapper Id] "+mapperId+" \nDELETE SQL["+pb.sql+"]ARGS", args)
	}
	return t.getDBhandle(mapperId, false).ExecuteUpdate(pb.sql, args...)
}

func (t *mapperHandler) DeleteAny(mapperId string, parameter any) (r int64, err error) {
	var pb *paramBean
	var args []any
	if pb, args, err = t.parseParameter(mapperId, parameter); err != nil {
		return r, err
	}
	if Logger.IsVaild {
		Logger.Debug("[Mapper Id] "+mapperId+" \nDELETE SQL["+pb.sql+"]ARGS", args)
	}
	return t.getDBhandle(mapperId, false).ExecuteUpdate(pb.sql, args...)
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

func cacheId(mapperId string) string {
	return Pre + mapperId
}

var defaultMapperHandler *mapperHandler

func NewInstance() JdaoMapper {
	return newMapperHandler()
}
