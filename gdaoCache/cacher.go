// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoCache

import (
	"fmt"
	"github.com/donnie4w/gdao/base"
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

type storeMode uint8

const (
	STRONG storeMode = 1
	SOFT   storeMode = 2
)

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
	mm        *Map[string, *Map[uint64, *CacheBean]]
	domain    string
	expire    int64
	storemode storeMode
}

func NewCacheHandle() *CacheHandle {
	domain := string(util.Base58EncodeForInt64(uint64(util.RandId())))
	expire := int64(5 * 60 * 1000)
	return &CacheHandle{mm: NewMap[string, *Map[uint64, *CacheBean]](), domain: domain, expire: expire, storemode: SOFT}
}

// SetDomain set the domain of cacheHandle
func (c *CacheHandle) SetDomain(domain string) *CacheHandle {
	c.domain = domain
	return c
}

// SetExpire set the data validity period in milliseconds.The default is 300*1000
func (c *CacheHandle) SetExpire(expire int64) *CacheHandle {
	c.expire = expire
	return c
}

// SetStoreMode set the storage mode. The default is SOFT
func (c *CacheHandle) SetStoreMode(mode storeMode) *CacheHandle {
	c.storemode = mode
	return c
}

func NewCacheHandle2(expire int64, mode storeMode) *CacheHandle {
	domain := string(util.Base58EncodeForInt64(uint64(util.RandId())))
	return &CacheHandle{mm: NewMap[string, *Map[uint64, *CacheBean]](), domain: domain, expire: expire, storemode: mode}
}

var defaultCacheHandle = NewCacheHandle()

var err_no_mapperid = fmt.Errorf("mapper binding error: no valid mapping id could be found")

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

func (c *cacher) Unbind(tablename string) {
	c.rmap.Del(tablename)
}

func (c *cacher) BindMapper(namespace string) error {
	ids := base.GetMapperIds(namespace)
	if len(ids) == 0 {
		return err_no_mapperid
	} else {
		for _, id := range ids {
			c.BindMapperId(namespace, id)
		}
	}
	return nil
}

func (c *cacher) BindMapperWithCacheHandle(namespace string, cacheHandle *CacheHandle) error {
	ids := base.GetMapperIds(namespace)
	if len(ids) == 0 {
		return err_no_mapperid
	} else {
		for _, id := range ids {
			c.BindMapperIdWithCacheHandle(namespace, id, cacheHandle)
		}
	}
	return nil
}

func (c *cacher) UnbindMapper(namespace string) {
	ids := base.GetMapperIds(namespace)
	if len(ids) > 0 {
		for _, id := range ids {
			c.UnbindMapperId(namespace, id)
		}
	}
}

func (c *cacher) BindMapperId(namespace, id string) error {
	if !base.HasMapperId(namespace + "." + id) {
		return err_no_mapperid
	}
	c.rmap.Put(mapperId(namespace, id), defaultCacheHandle.domain)
	return nil
}

func (c *cacher) BindMapperIdWithCacheHandle(namespace, id string, cacheHandle *CacheHandle) error {
	if !base.HasMapperId(namespace + "." + id) {
		return err_no_mapperid
	}
	c.rmap.Put(mapperId(namespace, id), cacheHandle.domain)
	c.cacheMap.Put(cacheHandle.domain, cacheHandle)
	return nil
}

func (c *cacher) UnbindMapperId(namespace, id string) {
	c.rmap.Del(mapperId(namespace, id))
}

func (c *cacher) GetMapperCache(domain, namepace, id string, condition *Condition) any {
	return c.GetCache(domain, mapperId(namepace, id), condition)
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

func (c *cacher) SetMapperCache(domain string, namespace, id string, condition *Condition, value any) bool {
	return c.SetCache(domain, mapperId(namespace, id), condition, value)
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

func (c *cacher) Clear(domain string, cacheId string) bool {
	if cacheHandle, b := c.cacheMap.Get(domain); b {
		return cacheHandle.mm.Del(cacheId)
	}
	return false
}

func (c *cacher) ClearMapper(domain string, namespace, id string) bool {
	if cacheHandle, b := c.cacheMap.Get(domain); b {
		return cacheHandle.mm.Del(mapperId(namespace, id))
	}
	return false
}

var memorymonitor = newMemoryMonitor(0.9, 3*time.Second, 1)

func (c *cacher) ticker() {
	tk := time.NewTicker(10 * time.Second)
	for {
		func() {
			defer func() { recover() }()
			select {
			case <-tk.C:
				b := memorymonitor.CheckMemoryPressure()
				c.cacheMap.Range(func(domain string, cachehandle *CacheHandle) bool {
					cachehandle.mm.Range(func(cacheId string, cacheBeanMap *Map[uint64, *CacheBean]) bool {
						cacheBeanMap.Range(func(condition uint64, cb *CacheBean) bool {
							if (b && cachehandle.storemode == SOFT) || time.Now().UnixMilli()-cachehandle.expire-cb.timestamp > 0 {
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

func (c *cacher) GetDomain(classname, tablename string) string {
	if classname != "" {
		if domain, b := c.rmap.Get(classname); b {
			return domain
		}
	}
	if tablename != "" {
		if domain, b := c.rmap.Get(tablename); b {
			return domain
		}
	}
	return ""
}

func (c *cacher) GetMapperDomain(namespace, id string) string {
	if domain, b := c.rmap.Get(mapperId(namespace, id)); b {
		return domain
	}
	return ""
}

func mapperId(namespace, id string) string {
	return base.MapperPre + namespace + "." + id
}
