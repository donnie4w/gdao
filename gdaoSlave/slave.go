// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoSlave

import (
	"database/sql"
	"github.com/donnie4w/gdao/base"
	"github.com/donnie4w/gdao/util"
)

var (
	// BindTable binds one or more table names to use the specified SQL database connection and database type for database qurey operation.
	//
	// Parameters:
	//   db: The SQL database connection to use for the table operations.
	//   dbtype: The type of the database (e.g., gdao.MYSQL, gdao.POSTGRESQL).
	//   tableNames: One or more table names to bind to the specified SQL database connection.
	//
	// Description:
	//   This function sets up one or more table names to use the provided SQL database connection and database type for database query operation.
	//
	// Example:
	//   // Assuming "users" and "orders" are table names corresponding to the standardized entity classes
	//   // And "db" is an instance of sql.DB configured for the database connection
	//   // And "gdao.MYSQL" represents the type of the database
	//   gdaoSlave.BindTable(db, gdao.MYSQL, "users", "orders")
	BindTable func(db *sql.DB, dbtype base.DBType, tableNames ...string)

	// BindTableWithDBhandle binds one or more table names to use the specified DB handle for database qurey operation.
	//
	// Parameters:
	//   dbhandle: The DB handle that encapsulates the CRUD operations and database connection management.
	//   tableNames: One or more table names to bind to the specified DB handle.
	//
	// Description:
	//   This function sets up one or more table names to use the provided DB handle for database query operation.
	//
	// Example:
	//   // Assuming "users" and "orders" are table names corresponding to the standardized entity classes
	//   // And "dbhandle" is an instance of DBhandle configured for the database connection and CRUD operations
	//   gdaoSlave.BindTableWithDBhandle(dbhandle, "users", "orders")
	BindTableWithDBhandle func(dbhandle base.DBhandle, tableNames ...string)

	UnbindTable func(tableNames ...string)

	// BindMapper binds the specified XML mapping namespace to use the given SQL database connection and database type for all methods within the namespace.
	//
	// Parameters:
	//   namespace: The namespace in the XML mapping files that corresponds to the CRUD operations to bind.
	//   db: The SQL database connection to use for the operations.
	//   dbtype: The type of the database (e.g., gdao.MYSQL, gdao.POSTGRESQL).
	//
	// Returns:
	//   An error if the binding fails, nil otherwise.
	//
	// Description:
	//   This function sets up the specified XML mapping namespace to use the provided SQL database connection and database type for database qurey operation.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "db" is an instance of sql.DB configured for the database connection
	//   // And "gdao.MYSQL" represents the type of the database
	//   err := gdaoSlave.BindMapper("com.example.mappers.users", db, gdao.MYSQL)
	//   if err != nil {
	//       log.Fatalf("Failed to bind the 'com.example.mappers.users' namespace with the specified data source and database type: %v", err)
	//   }
	BindMapper func(namespace string, db *sql.DB, dbtype base.DBType) error

	// BindMapperWithDBhandle binds the specified XML mapping namespace to use the given DB handle for all methods within the namespace.
	//
	// Parameters:
	//   namespace: The namespace in the XML mapping files that corresponds to the CRUD operations to bind.
	//   dbhandle: The DB handle that encapsulates the CRUD operations and database connection management.
	//
	// Returns:
	//   An error if the binding fails, nil otherwise.
	//
	// Description:
	//   This function sets up the specified XML mapping namespace to use the provided DB handle for database qurey operation.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "dbhandle" is an instance of DBhandle configured for the database connection and CRUD operations
	//   err := gdaoSlave.BindMapperWithDBhandle("com.example.mappers.users", dbhandle)
	//   if err != nil {
	//       log.Fatalf("Failed to bind the 'com.example.mappers.users' namespace with the specified DB handle: %v", err)
	//   }
	BindMapperWithDBhandle func(namespace string, dbhandle base.DBhandle) error

	UnbindMapper func(namespace string) error

	// BindMapperId binds a specific CRUD operation within the specified XML mapping namespace to use the given SQL database connection and database type.
	//
	// Parameters:
	//   namespace: The namespace in the XML mapping files that corresponds to the CRUD operations to bind.
	//   id: The ID of the CRUD operation within the namespace.
	//   db: The SQL database connection to use for the operation.
	//   dbtype: The type of the database (e.g., gdao.MYSQL, gdao.POSTGRESQL).
	//
	// Returns:
	//   An error if the binding fails, nil otherwise.
	//
	// Description:
	//   This function sets up a specific CRUD operation within the specified XML mapping namespace to use the provided SQL database connection and database type for database qurey operation.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "getUserById" is the ID of the query operation within the namespace
	//   // And "db" is an instance of sql.DB configured for the database connection
	//   // And "gdao.MYSQL" represents the type of the database
	//   err := gdaoSlave.BindMapperId("com.example.mappers.users", "getUserById", db, gdao.MYSQL)
	//   if err != nil {
	//       log.Fatalf("Failed to bind the 'getUserById' qurey operation with the specified data source and database type: %v", err)
	//   }
	BindMapperId func(namespace, id string, db *sql.DB, dbtype base.DBType) error

	// BindMapperIdWithDBhandle binds a specific CRUD operation within the specified XML mapping namespace to use the given DB handle.
	//
	// Parameters:
	//   namespace: The namespace in the XML mapping files that corresponds to the CRUD operations to bind.
	//   id: The ID of the CRUD operation within the namespace.
	//   dbhandle: The DB handle that encapsulates the CRUD operations and database connection management.
	//
	// Returns:
	//   An error if the binding fails, nil otherwise.
	//
	// Description:
	//   This function sets up a specific CRUD operation within the specified XML mapping namespace to use the provided DB handle for database qurey operation.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "getUserById" is the ID of the qurey operation within the namespace
	//   // And "dbhandle" is an instance of DBhandle configured for the database connection and CRUD operations
	//   err := gdaoSlave.BindMapperIdWithDBhandle("com.example.mappers.users", "getUserById", dbhandle)
	//   if err != nil {
	//       log.Fatalf("Failed to bind the 'getUserById' query operation with the specified DB handle: %v", err)
	//   }
	BindMapperIdWithDBhandle func(namespace, id string, dbhandle base.DBhandle) error

	UnbindMapperId func(namespace, id string) error

	Len       func() int64
	Get       func(classname, tableName string) base.DBhandle
	GetMapper func(namespace, id string) base.DBhandle
)

// BindClass binds the specified entity class to use the given SQL database connection and database type for database qurey operation.
//
// Parameters:
//
//	T: A generic type parameter representing the standardized entity class.
//	db: The SQL database connection to use for the operations.
//	dbtype: The type of the database (e.g., gdao.MYSQL, gdao.POSTGRESQL).
//
// Description:
//
//	This function sets up the specified entity class to use the provided SQL database connection and database type for database qurey operation.
//
// Example:
//
//	// Assuming User is the entity class
//	// And "db" is an instance of sql.DB configured for the database connection
//	// And "gdao.MYSQL" represents the type of the database
//	gdaoSlave.BindClass[User](db, gdao.MYSQL)
func BindClass[T base.TableBase[T]](db *sql.DB, dbtype base.DBType) {
	BindTable(db, dbtype, util.Classname[T]())
}

// BindClassWithDBhandle binds the specified entity class to use the given DB handle for database qurey operation.
//
// Parameters:
//
//	T: A generic type parameter representing the standardized entity class.
//	dbhandle: The DB handle that encapsulates the CRUD operations and database connection management.
//
// Description:
//
//	This function sets up the specified entity class to use the provided DB handle for database qurey operation.
//
// Example:
//
//	// Assuming User is the entity class
//	// And "dbhandle" is an instance of DBhandle configured for the database connection and CRUD operations
//	gdaoSlave.BindClassWithDBhandle[User](dbhandle)
func BindClassWithDBhandle[T base.TableBase[T]](dbhandle base.DBhandle) {
	BindTableWithDBhandle(dbhandle, util.Classname[T]())
}

func UnbindClass[T base.TableBase[T]]() {
	UnbindTable(util.Classname[T]())
}
