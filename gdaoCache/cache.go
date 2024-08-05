// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoCache

import (
	"github.com/donnie4w/gdao/base"
	"github.com/donnie4w/gdao/util"
)

func BindClass[T base.TableBase[T]]() {
	gdaocache.Bind(util.Classname[T]())
}

func BindClassWithCacheHandle[T base.TableBase[T]](cacheHandle *CacheHandle) {
	gdaocache.BindWithCacheHandle(util.Classname[T](), cacheHandle)
}

func UnbindClass[T base.TableBase[T]]() {
	gdaocache.Unbind(util.Classname[T]())
}

// BindTableNames binds one or more table names to enable the gdao caching mechanism for data operations on these tables.
// Parameters:
//
//	tableNames: A variadic list of strings representing the names of the tables to bind.
//
// The function configures the caching system to recognize and cache data operations for the specified tables.
func BindTableNames(tableNames ...string) {
	for _, tablename := range tableNames {
		gdaocache.Bind(tablename)
	}
}

// BindTableNamesWithCacheHandle binds one or more table names with a CacheHandle to enable the gdao caching mechanism for data operations on these tables.
// This function is useful when you want to enable caching for specific tables in your application with custom cache settings.
// Parameters:
//
//	cacheHandle: A pointer to a CacheHandle object that defines the caching behavior such as expiration time and eviction policies.
//	tableNames: A variadic list of strings representing the names of the tables to bind.
//
// The function configures the caching system to recognize and cache data operations for the specified tables with the provided cache settings.
func BindTableNamesWithCacheHandle(cacheHandle *CacheHandle, tableNames ...string) {
	for _, tablename := range tableNames {
		gdaocache.BindWithCacheHandle(tablename, cacheHandle)
	}
}

func UnbindTableNames(tableNames ...string) {
	for _, tablename := range tableNames {
		gdaocache.Unbind(tablename)
	}
}

// BindMapper binds the specified XML mapping namespace to use the gdao caching mechanism for all query operations.
//
// Parameters:
//
//	namespace: The namespace in the XML mapping files that corresponds to the CRUD operations to bind.
//
// Returns:
//
//	An error if the binding fails, nil otherwise.
//
// Description:
//
//	This function sets up the specified XML mapping namespace to use a caching mechanism for all query operations.
//
// Example:
//
//	// Assuming "com.example.mappers.users" is the namespace in the XML mapping files
//	gdaoCache.BindMapper("com.example.mappers.users")
func BindMapper(namespace string) error {
	return gdaocache.BindMapper(namespace)
}

// BindMapperWithCacheHandle binds the specified XML mapping namespace to use the given cache handle for all query operations.
//
// Parameters:
//
//	namespace: The namespace in the XML mapping files that corresponds to the CRUD operations to bind.
//	cacheHandle: The cache handle that manages the caching mechanism, including cache timeouts and eviction policies.
//
// Returns:
//
//	An error if the binding fails, nil otherwise.
//
// Description:
//
//	This function sets up the specified XML mapping namespace to use the provided cache handle for all query operations.
//
// Example:
//
//	// Assuming "com.example.mappers.users" is the namespace in the XML mapping files
//	err := gdaoCache.BindMapperWithCacheHandle("com.example.mappers.users", gdaoCache.NewCacheHandle().SetExpire(10000))
//	if err != nil {
//	    log.Fatalf("Failed to bind the 'com.example.mappers.users' namespace with the specified cache handle: %v", err)
//	}
func BindMapperWithCacheHandle(namespace string, cacheHandle *CacheHandle) error {
	return gdaocache.BindMapperWithCacheHandle(namespace, cacheHandle)
}

func UnbindMapper(namespace string) {
	gdaocache.UnbindMapper(namespace)
}

// BindMapperId binds a specific CRUD operation within the specified XML mapping namespace to use the gdao cache mechanism for query operations.
//
// Parameters:
//
//	namespace: The namespace in the XML mapping files that corresponds to the CRUD operations to bind.
//	id: The ID of the CRUD operation within the namespace.
//
// Returns:
//
//	An error if the binding fails, nil otherwise.
//
// Description:
//
//	This function sets up a specific CRUD operation within the specified XML mapping namespace to use the gdao cache mechanism for query operations.
//
// Example:
//
//	// Assuming "com.example.mappers.users" is the namespace in the XML mapping files
//	// And "getUserById" is the ID of the CRUD operation within the namespace
//	err := gdaoCache.BindMapperId("com.example.mappers.users", "getUserById")
//	if err != nil {
//	    log.Fatalf("Failed to bind the 'getUserById' query operation with the gdao cache mechanism: %v", err)
//	}
func BindMapperId(namespace, id string) error {
	return gdaocache.BindMapperId(namespace, id)
}

// BindMapperIdWithCacheHandle binds a specific CRUD operation within the specified XML mapping namespace to use the given cache handle for query operations.
//
// Parameters:
//
//	namespace: The namespace in the XML mapping files that corresponds to the CRUD operations to bind.
//	id: The ID of the CRUD operation within the namespace.
//	cacheHandle: The cache handle that manages the caching mechanism, including cache timeouts and eviction policies.
//
// Returns:
//
//	An error if the binding fails, nil otherwise.
//
// Description:
//
//	This function sets up a specific CRUD operation within the specified XML mapping namespace to use the provided cache handle for query operations.
//
// Example:
//
//	// Assuming "com.example.mappers.users" is the namespace in the XML mapping files
//	// And "getUserById" is the ID of the CRUD operation within the namespace
//	err := gdaoCache.BindMapperIdWithCacheHandle("com.example.mappers.users", "getUserById", gdaoCache.NewCacheHandle().SetExpire(10000))
//	if err != nil {
//	    log.Fatalf("Failed to bind the 'getUserById' qurey operation with the specified cache handle: %v", err)
//	}
func BindMapperIdWithCacheHandle(namespace, id string, cacheHandle *CacheHandle) error {
	return gdaocache.BindMapperIdWithCacheHandle(namespace, id, cacheHandle)
}

func UnbindMapperId(namespace, id string) {
	gdaocache.UnbindMapperId(namespace, id)
}

func GetCache(domain, cacheId string, condition *Condition) any {
	return gdaocache.GetCache(domain, cacheId, condition)
}

func GetMapperCache(domain, namepace, id string, condition *Condition) any {
	return gdaocache.GetMapperCache(domain, namepace, id, condition)
}

func SetMapperCache(domain string, namespace, id string, condition *Condition, value any) bool {
	return gdaocache.SetMapperCache(domain, namespace, id, condition, value)
}

func SetCache(domain string, cacheId string, condition *Condition, value any) bool {
	return gdaocache.SetCache(domain, cacheId, condition, value)
}

func GetMapperDomain(namespace, id string) string {
	return gdaocache.GetMapperDomain(namespace, id)
}

func GetDomain(classname, tablename string) string {
	return gdaocache.GetDomain(classname, tablename)
}

type cache interface {
	Bind(tablename string)

	BindWithCacheHandle(tablename string, cacheHandle *CacheHandle)

	Unbind(tablename string)

	BindMapper(namespace string) error

	BindMapperWithCacheHandle(namespace string, cacheHandle *CacheHandle) error

	UnbindMapper(namespace string)

	BindMapperId(namespace, id string) error

	BindMapperIdWithCacheHandle(namespace, id string, cacheHandle *CacheHandle) error

	UnbindMapperId(namespace, id string)

	GetMapperCache(domain, namepace, id string, condition *Condition) any

	GetCache(domain, cacheId string, condition *Condition) any

	SetMapperCache(domain string, namespace, id string, condition *Condition, value any) bool

	SetCache(domain string, cacheId string, condition *Condition, value any) bool

	GetMapperDomain(namespace, id string) string

	GetDomain(classname, tablename string) string
}
