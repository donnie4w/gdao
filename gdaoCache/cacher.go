// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoCache

import (
	"fmt"
	. "github.com/donnie4w/gofer/hashmap"
	"github.com/donnie4w/gofer/util"
	"hash/fnv"
	"strconv"
	"time"
)

type SqlKV struct {
	sql  string
	args []any
}

func newSqlKV(sql string, args ...any) *SqlKV {
	return &SqlKV{sql: sql, args: args}
}

func (s *SqlKV) hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(s.sql))
	for _, arg := range s.args {
		switch v := arg.(type) {
		case int:
			h.Write([]byte(strconv.Itoa(v)))
		case int8:
			h.Write([]byte(strconv.FormatInt(int64(v), 10)))
		case int16:
			h.Write([]byte(strconv.FormatInt(int64(v), 10)))
		case int32:
			h.Write([]byte(strconv.FormatInt(int64(v), 10)))
		case int64:
			h.Write([]byte(strconv.FormatInt(v, 10)))
		case uint:
			h.Write([]byte(strconv.FormatUint(uint64(v), 10)))
		case uint8:
			h.Write([]byte(strconv.FormatUint(uint64(v), 10)))
		case uint16:
			h.Write([]byte(strconv.FormatUint(uint64(v), 10)))
		case uint32:
			h.Write([]byte(strconv.FormatUint(uint64(v), 10)))
		case uint64:
			h.Write([]byte(strconv.FormatUint(v, 10)))
		case float32:
			h.Write([]byte(strconv.FormatFloat(float64(v), 'f', -1, 32)))
		case float64:
			h.Write([]byte(strconv.FormatFloat(v, 'f', -1, 64)))
		case bool:
			h.Write([]byte(strconv.FormatBool(v)))
		case string:
			h.Write([]byte(v))
		case []uint8:
			h.Write(v)
		case time.Time:
			timestamp := v.UnixNano()
			h.Write([]byte(strconv.FormatInt(timestamp, 10)))
		default:
			h.Write([]byte(fmt.Sprintf("%v", arg)))
		}
	}
	return h.Sum64()
}

type CacheBean struct {
	timestamp int64
	value     any
}

type Condition struct {
	sqlKV *SqlKV
	node  string
}

func NewCondition(node, sql string, args ...any) *Condition {
	return &Condition{sqlKV: newSqlKV(sql, args...), node: node}
}

func (c *Condition) hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(c.node))
	h.Write([]byte(strconv.FormatInt(int64(c.sqlKV.hash()), 10)))
	return h.Sum64()
}

type CacheHandle struct {
	mm     *Map[string, *Map[uint64, *CacheBean]]
	domain string
	expire int64
}

func newCacheHandle() *CacheHandle {
	domain := string(util.Base58EncodeForInt64(uint64(util.RandId())))
	expire := int64(5 * 60 * 1000)
	return &CacheHandle{mm: NewMap[string, *Map[uint64, *CacheBean]](), domain: domain, expire: expire}
}

func NewCacheHandleWithDomain(domain string) *CacheHandle {
	expire := int64(5 * 60 * 1000)
	return &CacheHandle{mm: NewMap[string, *Map[uint64, *CacheBean]](), domain: domain, expire: expire}
}

func NewCacheHandleWithExpire(expire int64) *CacheHandle {
	domain := string(util.Base58EncodeForInt64(uint64(util.RandId())))
	return &CacheHandle{mm: NewMap[string, *Map[uint64, *CacheBean]](), domain: domain, expire: expire}
}

func NewCacheHandle(domain string, expire int64) *CacheHandle {
	return &CacheHandle{mm: NewMap[string, *Map[uint64, *CacheBean]](), domain: domain, expire: expire}
}

var defaultCacheHandle = newCacheHandle()

type cacher struct {
	cacheMap *Map[string, *CacheHandle]
	rmap     *Map[string, string]
}

func newcache() cache {
	c := &cacher{cacheMap: NewMap[string, *CacheHandle](), rmap: NewMap[string, string]()}
	go c.ticker()
	return c
}

var gdaocache = newcache()

func (c *cacher) Bind(tablename string) {
	c.rmap.Put(tablename, defaultCacheHandle.domain)
}

func (c *cacher) BindWithCacheHandle(tablename string, cacheHandle *CacheHandle) {
	c.rmap.Put(tablename, cacheHandle.domain)
	c.cacheMap.Put(cacheHandle.domain, cacheHandle)
}

func (c *cacher) Remove(tablename string) {
	c.rmap.Del(tablename)
}

func (c *cacher) BindMapper(mapperId string) {
	c.rmap.Put(mapperId, defaultCacheHandle.domain)
}

func (c *cacher) BindMapperWithCacheHandle(mapperId string, cacheHandle *CacheHandle) {
	c.rmap.Put(mapperId, cacheHandle.domain)
	c.cacheMap.Put(cacheHandle.domain, cacheHandle)
}

func (c *cacher) RemoveMapper(mapperId string) {
	c.rmap.Del(mapperId)
}

func (c *cacher) GetCache(domain, cacheId string, condition *Condition) any {
	if cacheId == "" || condition == nil {
		return nil
	}
	if domain == "" {
		domain = defaultCacheHandle.domain
	}
	if cacheHandle, b := c.cacheMap.Get(domain); b {
		if cacheBeanMap, b := cacheHandle.mm.Get(cacheId); b {
			hashcode := condition.hash()
			if cacheBean, b := cacheBeanMap.Get(hashcode); b {
				if time.Now().UnixMilli()-cacheBean.timestamp-cacheHandle.expire <= 0 {
					return cacheBean.value
				} else {
					cacheBeanMap.Del(hashcode)
				}
			}
		}
	}
	return nil
}

func (c *cacher) SetCache(domain string, cacheId string, condition *Condition, value any) bool {
	if cacheId == "" || condition == nil || value == nil {
		return false
	}
	if domain == "" {
		domain = defaultCacheHandle.domain
	}
	var cacheHandle *CacheHandle
	if cacheHandle, _ = c.cacheMap.Get(domain); cacheHandle == nil {
		cacheHandle = defaultCacheHandle
		c.cacheMap.Put(domain, cacheHandle)
	}
	var cacheBeanMap *Map[uint64, *CacheBean]
	if cacheBeanMap, _ = cacheHandle.mm.Get(cacheId); cacheBeanMap == nil {
		cacheBeanMap = NewMap[uint64, *CacheBean]()
		cacheHandle.mm.Put(cacheId, cacheBeanMap)
	}
	cacheBeanMap.Put(condition.hash(), &CacheBean{time.Now().UnixMilli(), value})
	return true
}

func (c *cacher) ticker() {
	tk := time.NewTicker(10 * time.Second)
	for {
		func() {
			defer func() { recover() }()
			select {
			case <-tk.C:
				c.cacheMap.Range(func(domain string, cachehandle *CacheHandle) bool {
					cachehandle.mm.Range(func(cacheId string, cacheBeanMap *Map[uint64, *CacheBean]) bool {
						cacheBeanMap.Range(func(condition uint64, cb *CacheBean) bool {
							if time.Now().UnixMilli()-cachehandle.expire-cb.timestamp > 0 {
								cacheBeanMap.Del(condition)
							}
							return true
						})
						return true
					})
					return true
				})
			}
		}()
	}
}

func (c *cacher) GetDomain(cacheId string) string {
	if domain, b := c.rmap.Get(cacheId); b {
		return domain
	}
	return ""
}
