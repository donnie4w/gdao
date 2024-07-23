// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoSlave

import (
	"database/sql"
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
	BindWithDBhandle = defaultSlaveHandler.addWithDBhandle
	Bind = defaultSlaveHandler.add
	BindMapper = defaultSlaveHandler.addMapper
	BindMapperWithDBhandle = defaultSlaveHandler.addMapperWithDBhandle
	Remove = defaultSlaveHandler.remove
	RemoveMapper = defaultSlaveHandler.removeMapper
	Len = defaultSlaveHandler.len
	Get = defaultSlaveHandler.get
}

func newSlaveHandler() *slaveHandler {
	return &slaveHandler{slavemap: hashmap.NewMapL[string, []DBhandle](), mutex: &sync.Mutex{}}
}

func (t *slaveHandler) add(tableName string, db *sql.DB, dbtype DBType) {
	t.addWithDBhandle(tableName, Newdbhandle(db, dbtype))
}

func (t *slaveHandler) addWithDBhandle(tableName string, dbHandle DBhandle) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if dblist, ok := t.slavemap.Get(tableName); ok {
		dblist = append(dblist, dbHandle)
		t.slavemap.Put(tableName, dblist)
	} else {
		t.slavemap.Put(tableName, []DBhandle{dbHandle})
	}
}

func (t *slaveHandler) addMapper(mapperId string, db *sql.DB, dbtype DBType) {
	t.add(Pre+mapperId, db, dbtype)
}

func (t *slaveHandler) addMapperWithDBhandle(mapperId string, dbHandle DBhandle) {
	t.addWithDBhandle(Pre+mapperId, dbHandle)
}

func (t *slaveHandler) remove(tableName string) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.slavemap.Del(tableName)
}

func (t *slaveHandler) removeMapper(mapperId string) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.slavemap.Del(Pre + mapperId)
}

func (t *slaveHandler) len() int64 {
	return t.slavemap.Len()
}

func (t *slaveHandler) get(classname, tableName, mapperId string) DBhandle {
	var dblist []DBhandle

	if classname != "" {
		dblist, _ = t.slavemap.Get(classname)
	}

	if len(dblist) == 0 && tableName != "" {
		dblist, _ = t.slavemap.Get(tableName)
	}

	if len(dblist) == 0 && mapperId != "" {
		dblist, _ = t.slavemap.Get(Pre + mapperId)
	}

	if length := len(dblist); length > 0 {
		if length == 1 {
			return dblist[0]
		} else if length > 1 {
			return dblist[util.Rand(int(t.len()))]
		}
	}
	return nil
}

var Newdbhandle func(db *sql.DB, dbtype DBType) DBhandle
