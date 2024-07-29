// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdao

import (
	. "github.com/donnie4w/gdao/base"
	"github.com/donnie4w/gofer/hashmap"
)

type container interface {
	putTables(dbhandle DBhandle, tables ...string)

	delTables(tables ...string)

	putMapper(namespace string, dbhandle DBhandle)

	putMapperId(namespace, id string, dbhandle DBhandle)

	delMapper(namespace string)

	delMapperId(namespace, id string)

	get(name string) (DBhandle, bool)

	getMapper(namespace, id string) (DBhandle, bool)

	len() int64
}

type containerHandle struct {
	handleMap *hashmap.MapL[string, DBhandle]
}

func newContainer() container {
	return &containerHandle{handleMap: hashmap.NewMapL[string, DBhandle]()}
}

func (c *containerHandle) putMapper(namespace string, dbhandle DBhandle) {
	ids := GetMapperIds(namespace)
	if len(ids) > 0 {
		for _, id := range ids {
			c.putMapperId(namespace, id, dbhandle)
		}
	}
}

func (c *containerHandle) putMapperId(namespace, id string, dbhandle DBhandle) {
	c.handleMap.Put(mapperId(namespace, id), dbhandle)
}

func (c *containerHandle) delMapper(namespace string) {
	ids := GetMapperIds(namespace)
	if len(ids) > 0 {
		for _, id := range ids {
			c.delMapperId(namespace, id)
		}
	}
}

func (c *containerHandle) delMapperId(namespace, id string) {
	c.handleMap.Del(mapperId(namespace, id))
}

func (c *containerHandle) putTables(dbhandle DBhandle, tables ...string) {
	for _, table := range tables {
		c.handleMap.Put(table, dbhandle)
	}
}

func (c *containerHandle) delTables(tables ...string) {
	for _, table := range tables {
		c.handleMap.Del(table)
	}
}

func (c *containerHandle) get(name string) (DBhandle, bool) {
	return c.handleMap.Get(name)
}

func (c *containerHandle) getMapper(namespace, id string) (DBhandle, bool) {
	return c.handleMap.Get(mapperId(namespace, id))
}

func (c *containerHandle) len() int64 {
	return c.handleMap.Len()
}

func mapperId(namespace, id string) string {
	return MapperPre + namespace + "." + id
}
