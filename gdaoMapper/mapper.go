// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoMapper

import (
	"database/sql"
	"github.com/donnie4w/gdao/base"
)

// GdaoMapper is the interface for the gdaoMapper module, providing methods to manage transactions and database connections.
// This interface defines the basic operations required for CRUD functionalities.
type GdaoMapper interface {
	IsAutocommit() bool
	SetAutocommit(autocommit bool) (err error)
	UseTransaction(tx base.Transaction)
	Rollback() (err error)
	Commit() (err error)
	UseDBhandle(dbhandler base.DBhandle)
	UseDBhandleWithDB(db *sql.DB, dbType base.DBType)

	// SelectBean executes a query based on the specified XML mapping mapper ID and returns a single row of data as a DataBean.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//	 args: Variable length argument list, which corresponds to placeholder arguments of mapperId.
	//
	// Returns:
	//   A pointer to a DataBean containing the data retrieved by the query.
	//
	// Description:
	//   This function executes a query based on the specified XML mapping mapper ID and returns a single row of data.
	//   The arguments can be a basic data type, an entity class object, a map, or a slice, depending on the requirements of the query.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "getUserById" is the ID of the CRUD operation within the namespace
	//   userBean, err := gdaoMapper.SelectBeanDirect("com.example.mappers.users.getUserById", 1)
	//   if err != nil {
	//       log.Fatalf("Failed to select user: %v", err)
	//   }
	SelectBean(mapperId string, args ...any) *base.DataBean

	// SelectBeans executes a query based on the specified XML mapping mapper ID and returns multiple rows of data as DataBeans.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   args: Variable length argument list, which corresponds to placeholder arguments of mapperId.
	//
	// Returns:
	//   A slice of pointers to DataBeans containing the data retrieved by the query.
	//
	// Description:
	//   This function executes a query based on the specified XML mapping mapper ID and returns multiple rows of data.
	//   The arguments can be a basic data type, an entity class object, a map, or a slice, depending on the requirements of the query.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "getUsersByLimit" is the ID of the CRUD operation within the namespace
	//   usersBeans, err := gdaoMapper.SelectBeansDirect("com.example.mappers.users.getUsersByLimit", 10)
	//   if err != nil {
	//       log.Fatalf("Failed to select users: %v", err)
	//   }
	//
	SelectBeans(mapperId string, args ...any) *base.DataBeans

	// Insert executes an insertXml operation based on the specified XML mapping mapper ID.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   args: Variable length argument list, which corresponds to placeholder arguments of mapperId.
	//
	// Returns:
	//   The number of rows affected by the insertXml operation as an int64, or an error if the operation fails.
	//
	// Description:
	//   This function executes an insertXml operation based on the specified XML mapping mapper ID.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "insertUser" is the ID of the CRUD operation within the namespace
	//   rowsAffected, err := gdaoMapper.InsertDirect("com.example.mappers.users.insertUser", 1,"hello world",100)
	//   if err != nil {
	//       log.Fatalf("Failed to insertXml user: %v", err)
	//   }
	Insert(mapperId string, args ...any) (sql.Result, error)

	// Update executes an updateXml operation based on the specified XML mapping mapper ID.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   args: Variable length argument list, which corresponds to placeholder arguments of mapperId.
	//
	// Returns:
	//   The number of rows affected by the updateXml operation as an int64, or an error if the operation fails.
	//
	// Description:
	//   This function executes an updateXml operation based on the specified XML mapping mapper ID.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "updateUserByEmail" is the ID of the CRUD operation within the namespace
	//   rowsAffected, err := gdaoMapper.UpdateDirect("com.example.mappers.users.updateUserByEmail","hello world" ,"donnie4w@gmail.com")
	//   if err != nil {
	//       log.Fatalf("Failed to updateXml user: %v", err)
	//   }
	Update(mapperId string, args ...any) (sql.Result, error)

	// Delete executes a deleteXml operation based on the specified XML mapping mapper ID.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   args: Variable length argument list, which corresponds to placeholder arguments of mapperId.
	//
	// Returns:
	//   The number of rows affected by the deleteXml operation as an int64, or an error if the operation fails.
	//
	// Description:
	//   This function executes a deleteXml operation based on the specified XML mapping mapper ID.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "deleteUserById" is the ID of the CRUD operation within the namespace
	//   rowsAffected, err := gdaoMapper.DeleteDirect("com.example.mappers.users.deleteUserById", 1)
	//   if err != nil {
	//       log.Fatalf("Failed to deleteXml user: %v", err)
	//   }
	Delete(mapperId string, args ...any) (sql.Result, error)
}

var (
	IsAutocommit      func() bool
	SetAutocommit     func(autocommit bool) (err error)
	UseTransaction    func(tx base.Transaction)
	Rollback          func() (err error)
	Commit            func() (err error)
	UseDBhandle       func(dbhandler base.DBhandle)
	UseDBhandleWithDB func(db *sql.DB, dbType base.DBType)

	// SelectBean executes a query based on the specified XML mapping mapper ID and returns a single row of data as a DataBean.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//	 args: Variable length argument list, which corresponds to placeholder arguments of mapperId.
	//
	// Returns:
	//   A pointer to a DataBean containing the data retrieved by the query.
	//
	// Description:
	//   This function executes a query based on the specified XML mapping mapper ID and returns a single row of data.
	//   The arguments can be a basic data type, an entity class object, a map, or a slice, depending on the requirements of the query.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "getUserById" is the ID of the CRUD operation within the namespace
	//   userBean, err := gdaoMapper.SelectBeanDirect("com.example.mappers.users.getUserById", 1)
	//   if err != nil {
	//       log.Fatalf("Failed to select user: %v", err)
	//   }
	SelectBean func(mapperId string, args ...any) *base.DataBean

	// SelectBeans executes a query based on the specified XML mapping mapper ID and returns multiple rows of data as DataBeans.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   args: Variable length argument list, which corresponds to placeholder arguments of mapperId.
	//
	// Returns:
	//   A slice of pointers to DataBeans containing the data retrieved by the query, or an error if the query fails.
	//
	// Description:
	//   This function executes a query based on the specified XML mapping mapper ID and returns multiple rows of data.
	//   The arguments can be a basic data type, an entity class object, a map, or a slice, depending on the requirements of the query.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "getUsersByLimit" is the ID of the CRUD operation within the namespace
	//   usersBeans, err := gdaoMapper.SelectBeansDirect("com.example.mappers.users.getUsersByLimit", 10)
	//   if err != nil {
	//       log.Fatalf("Failed to select users: %v", err)
	//   }
	//
	SelectBeans func(mapperId string, args ...any) *base.DataBeans

	// Insert executes an insertXml operation based on the specified XML mapping mapper ID.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   args: Variable length argument list, which corresponds to placeholder arguments of mapperId.
	//
	// Returns:
	//   The number of rows affected by the insertXml operation as an int64, or an error if the operation fails.
	//
	// Description:
	//   This function executes an insertXml operation based on the specified XML mapping mapper ID.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "insertUser" is the ID of the CRUD operation within the namespace
	//   // And "newUser" is the user object to insertXml
	//   rowsAffected, err := gdaoMapper.InsertDirect("com.example.mappers.users.insertUser", newUser)
	//   if err != nil {
	//       log.Fatalf("Failed to insertXml user: %v", err)
	//   }
	Insert func(mapperId string, args ...any) (r sql.Result, err error)

	// Update executes an updateXml operation based on the specified XML mapping mapper ID.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   args: Variable length argument list, which corresponds to placeholder arguments of mapperId.
	//
	// Returns:
	//   The number of rows affected by the updateXml operation as an int64, or an error if the operation fails.
	//
	// Description:
	//   This function executes an updateXml operation based on the specified XML mapping mapper ID.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "updateUser" is the ID of the CRUD operation within the namespace
	//   rowsAffected, err := gdaoMapper.UpdateDirect("com.example.mappers.users.updateUser", 10)
	//   if err != nil {
	//       log.Fatalf("Failed to insertXml user: %v", err)
	//   }
	Update func(mapperId string, args ...any) (r sql.Result, err error)

	// Delete executes an deleteXml operation based on the specified XML mapping mapper ID.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   args: Variable length argument list, which corresponds to placeholder arguments of mapperId.
	//
	// Returns:
	//   The number of rows affected by the deleteXml operation as an int64, or an error if the operation fails.
	//
	// Description:
	//   This function executes an deleteXml operation based on the specified XML mapping mapper ID.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "deleteUser" is the ID of the CRUD operation within the namespace
	//   rowsAffected, err := gdaoMapper.DeleteDirect("com.example.mappers.users.deleteUser", 10)
	//   if err != nil {
	//       log.Fatalf("Failed to insertXml user: %v", err)
	//   }
	Delete func(mapperId string, args ...any) (r sql.Result, err error)
)

// Select executes a query based on the specified XML mapping mapper ID and returns a single row of data as an instance of the generic type T.
//
// Parameters:
//
//	T: A generic type parameter representing the type of the data to be returned.
//	mapperId: The ID of the CRUD operation within the XML mapping namespace.
//	args: Variable length argument list, which corresponds to placeholder arguments of mapperId.
//
// Returns:
//
//	A pointer to an instance of type T containing the data retrieved by the query, or an error if the query fails.
//
// Description:
//
//	This function executes a query based on the specified XML mapping mapper ID and returns a single row of data.
//
// Example:
//
//	// Assuming "com.example.mappers.users" is the namespace in the XML mapping files
//	// And "getUserById" is the ID of the CRUD operation within the namespace
//	userResult, err := gdaoMapper.SelectDirect[dao.User]("com.example.mappers.users.getUserById", 1)
//	if err != nil {
//	    log.Fatalf("Failed to select user: %v", err)
//	}
func Select[T any](mapperId string, args ...any) (*T, error) {
	if len(args) == 1 {
		return selectAny[T](mapperId, args[0])
	}
	return (*mapperInvoke[T])(defaultMapperHandler).SelectDirect(mapperId, args...)
}

// selectAny executes a query based on the specified XML mapping mapper ID and returns a single row of data as an instance of the generic type T.
//
// Parameters:
//
//	T: A generic type parameter representing the type of the data to be returned.
//	mapperId: The ID of the CRUD operation within the XML mapping namespace.
//	parameter: The parameter to pass to the query. Can be a basic data type, an entity class object, a map, or a slice.
//
// Returns:
//
//	A pointer to an instance of type T containing the data retrieved by the query, or an error if the query fails.
//
// Description:
//
//	This function executes a query based on the specified XML mapping mapper ID and returns a single row of data.
//	The parameter can be a basic data type, an entity class object, a map, or a slice, depending on the requirements of the query.
//	It is typically used when you need to retrieve a single row of data from the database.
//
// Example:
//
//	// Assuming "com.example.mappers.users" is the namespace in the XML mapping files
//	// And "getUserById" is the ID of the CRUD operation within the namespace
//
//	userResult, err := gdaoMapper.Select[dao.User]("com.example.mappers.users.getUserById", 1)
//	if err != nil {
//	    log.Fatalf("Failed to select user: %v", err)
//	}
func selectAny[T any](mapperId string, parameter any) (*T, error) {
	return (*mapperInvoke[T])(defaultMapperHandler).Select(mapperId, parameter)
}

// Selects executes a query based on the specified XML mapping mapper ID and returns multiple rows of data as instances of the generic type T.
//
// Parameters:
//
//		T: A generic type parameter representing the type of the data to be returned.
//		mapperId: The ID of the CRUD operation within the XML mapping namespace.
//	 args: Variable length argument list, which corresponds to placeholder arguments of mapperId.
//
// Returns:
//
//	A slice of pointers to instances of type T containing the data retrieved by the query, or an error if the query fails.
//
// Description:
//
//	This function executes a query based on the specified XML mapping mapper ID and returns multiple rows of data.
//
// Example:
//
//	// Assuming "com.example.mappers.users" is the namespace in the XML mapping files
//	// And "getUsersByLimit" is the ID of the CRUD operation within the namespace
//	usersResult, err := gdaoMapper.SelectsDirect[dao.User]("com.example.mappers.users.getUsersByLimit", 10)
//	if err != nil {
//	    log.Fatalf("Failed to select users: %v", err)
//	}
func Selects[T any](mapperId string, args ...any) ([]*T, error) {
	if len(args) == 1 {
		return selectsAny[T](mapperId, args[0])
	}
	return (*mapperInvoke[T])(defaultMapperHandler).SelectsDirect(mapperId, args...)
}

// selectsAny executes a query based on the specified XML mapping mapper ID and returns multiple rows of data as instances of the generic type T.
//
// Parameters:
//
//	T: A generic type parameter representing the type of the data to be returned.
//	mapperId: The ID of the CRUD operation within the XML mapping namespace.
//	parameter: The parameter to pass to the query. Can be a basic data type, an entity class object, a map, or a slice.
//
// Returns:
//
//	A slice of pointers to instances of type T containing the data retrieved by the query, or an error if the query fails.
//
// Description:
//
//	This function executes a query based on the specified XML mapping mapper ID and returns multiple rows of data.
//	The parameter can be a basic data type, an entity class object, a map, or a slice, depending on the requirements of the query.
//
// Example:
//
//	// Assuming "com.example.mappers.users" is the namespace in the XML mapping files
//	// And "getUsersByLimit" is the ID of the CRUD operation within the namespace
//	usersResult, err := gdaoMapper.Selects[dao.User]("com.example.mappers.users.getUsersByLimit", 10)
//	if err != nil {
//	    log.Fatalf("Failed to select users: %v", err)
//	}
func selectsAny[T any](mapperId string, parameter any) ([]*T, error) {
	return (*mapperInvoke[T])(defaultMapperHandler).Selects(mapperId, parameter)
}
