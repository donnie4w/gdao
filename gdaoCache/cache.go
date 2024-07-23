// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoCache

import (
	"github.com/donnie4w/gdao/base"
)

func BindClass[T base.TableBase[T]]() {
	gdaocache.Bind(base.Classname[T]())
}

func BindClassWithCacheHandle[T base.TableBase[T]](cacheHandle *CacheHandle) {
	gdaocache.BindWithCacheHandle(base.Classname[T](), cacheHandle)
}

func RemoveClass[T base.TableBase[T]]() {
	gdaocache.Remove(base.Classname[T]())
}

func BindMapper(mapperId string) {
	gdaocache.BindMapper(base.Pre + mapperId)
}
func BindMapperWithCacheHandle(mapperId string, cacheHandle *CacheHandle) {
	gdaocache.BindMapperWithCacheHandle(base.Pre+mapperId, cacheHandle)
}
func RemoveMapper(mapperId string) {
	gdaocache.RemoveMapper(base.Pre + mapperId)
}

func GetCache(domain, cacheId string, condition *Condition) any {
	return gdaocache.GetCache(domain, cacheId, condition)
}
func SetCache(domain string, cacheId string, condition *Condition, value any) bool {
	return gdaocache.SetCache(domain, cacheId, condition, value)
}

func GetDomain(cacheId string) string {
	return gdaocache.GetDomain(cacheId)
}

type cache interface {
	Bind(tablename string)

	BindWithCacheHandle(tablename string, cacheHandle *CacheHandle)

	Remove(tablename string)

	BindMapper(mapperId string)

	BindMapperWithCacheHandle(mapperId string, cacheHandle *CacheHandle)

	RemoveMapper(mapperId string)

	GetCache(domain, cacheId string, condition *Condition) any

	SetCache(domain string, cacheId string, condition *Condition, value any) bool

	GetDomain(cacheId string) string
}
