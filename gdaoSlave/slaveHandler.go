// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoSlave

import (
	"database/sql"
	"fmt"
	. "github.com/donnie4w/gdao/base"
	"github.com/donnie4w/gofer/hashmap"
	"github.com/donnie4w/gofer/util"
	"sync"
)

type slaveHandler struct {
	slavemap *hashmap.MapL[string, []DBhandle]
	mutex    *sync.Mutex
}

var defaultSlaveHandler *slaveHandler

func init() {
	defaultSlaveHandler = newSlaveHandler()
	Len = defaultSlaveHandler.len
	BindTableWithDBhandle = defaultSlaveHandler.bindTableWithDBhandle
	BindTable = defaultSlaveHandler.bindTable
	UnbindTable = defaultSlaveHandler.unbindTable

	BindMapper = defaultSlaveHandler.bindMapper
	BindMapperWithDBhandle = defaultSlaveHandler.bindMapperWithDBhandle
	UnbindMapper = defaultSlaveHandler.unbindMapper

	BindMapperId = defaultSlaveHandler.bindMapperId
	BindMapperIdWithDBhandle = defaultSlaveHandler.bindMapperIdWithDBhandle
	UnbindMapperId = defaultSlaveHandler.unbindMapperId

	Len = defaultSlaveHandler.len
	Get = defaultSlaveHandler.get
	GetMapper = defaultSlaveHandler.getMapper
}

func newSlaveHandler() *slaveHandler {
	return &slaveHandler{slavemap: hashmap.NewMapL[string, []DBhandle](), mutex: &sync.Mutex{}}
}

var err_no_mapperid = fmt.Errorf("mapper binding error: no valid mapping id could be found")

func (t *slaveHandler) bind(s string, db *sql.DB, dbtype DBType) {
	t.bindWithDBhandle(s, Newdbhandle(db, dbtype))
}

func (t *slaveHandler) bindWithDBhandle(s string, dbHandle DBhandle) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if dblist, ok := t.slavemap.Get(s); ok {
		t.slavemap.Put(s, append(dblist, dbHandle))
	} else {
		t.slavemap.Put(s, []DBhandle{dbHandle})
	}
}

func (t *slaveHandler) bindTable(db *sql.DB, dbtype DBType, tableNames ...string) {
	t.bindTableWithDBhandle(Newdbhandle(db, dbtype), tableNames...)
}

func (t *slaveHandler) bindTableWithDBhandle(dbHandle DBhandle, tableNames ...string) {
	for _, tableName := range tableNames {
		t.bindWithDBhandle(tableName, dbHandle)
	}
}

func (t *slaveHandler) unbindTable(tableNames ...string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	for _, tableName := range tableNames {
		t.slavemap.Del(tableName)
	}
}

func (t *slaveHandler) bindMapper(namespace string, db *sql.DB, dbtype DBType) error {
	return t.bindMapperWithDBhandle(namespace, Newdbhandle(db, dbtype))
}

func (t *slaveHandler) bindMapperWithDBhandle(namespace string, dbHandle DBhandle) error {
	ids := GetMapperIds(namespace)
	if len(ids) == 0 {
		return err_no_mapperid
	} else {
		for _, id := range ids {
			t.bindWithDBhandle(mapperId(namespace, id), dbHandle)
		}
	}
	return nil
}

func (t *slaveHandler) bindMapperId(namespace, id string, db *sql.DB, dbtype DBType) error {
	return t.bindMapperIdWithDBhandle(namespace, id, Newdbhandle(db, dbtype))
}

func (t *slaveHandler) bindMapperIdWithDBhandle(namespace, id string, dbHandle DBhandle) error {
	if !HasMapperId(namespace + "." + id) {
		return err_no_mapperid
	}
	t.bindWithDBhandle(mapperId(namespace, id), dbHandle)
	return nil
}

func (t *slaveHandler) unbindMapper(namespace string) error {
	ids := GetMapperIds(namespace)
	if len(ids) == 0 {
		return err_no_mapperid
	} else {
		for _, id := range ids {
			t.slavemap.Del(mapperId(namespace, id))
		}
	}
	return nil
}

func (t *slaveHandler) unbindMapperId(namespace, id string) error {
	if !HasMapperId(namespace + "." + id) {
		return err_no_mapperid
	}
	t.slavemap.Del(mapperId(namespace, id))
	return nil
}

func (t *slaveHandler) len() int64 {
	return t.slavemap.Len()
}

func (t *slaveHandler) getMapper(namespace, id string) DBhandle {
	if namespace != "" && id != "" {
		dblist, _ := t.slavemap.Get(mapperId(namespace, id))
		if length := len(dblist); length > 0 {
			if length == 1 {
				return dblist[0]
			} else if length > 1 {
				return dblist[util.RandUint(uint(t.len()))]
			}
		}
	}
	return nil
}

func (t *slaveHandler) get(classname, tableName string) DBhandle {
	var dblist []DBhandle

	if classname != "" {
		dblist, _ = t.slavemap.Get(classname)
	}

	if len(dblist) == 0 && tableName != "" {
		dblist, _ = t.slavemap.Get(tableName)
	}

	if length := len(dblist); length > 0 {
		if length == 1 {
			return dblist[0]
		} else if length > 1 {
			return dblist[util.RandUint(uint(t.len()))]
		}
	}
	return nil
}

var Newdbhandle func(db *sql.DB, dbtype DBType) DBhandle

func mapperId(namespace, id string) string {
	return MapperPre + namespace + "." + id
}
