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

	// SelectBeanDirect executes a query based on the specified XML mapping mapper ID and returns a single row of data as a DataBean.
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
	SelectBeanDirect(mapperId string, args ...any) *base.DataBean

	// SelectBeansDirect executes a query based on the specified XML mapping mapper ID and returns multiple rows of data as DataBeans.
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
	SelectBeansDirect(mapperId string, args ...any) *base.DataBeans

	// InsertDirect executes an insert operation based on the specified XML mapping mapper ID.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   args: Variable length argument list, which corresponds to placeholder arguments of mapperId.
	//
	// Returns:
	//   The number of rows affected by the insert operation as an int64, or an error if the operation fails.
	//
	// Description:
	//   This function executes an insert operation based on the specified XML mapping mapper ID.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "insertUser" is the ID of the CRUD operation within the namespace
	//   rowsAffected, err := gdaoMapper.InsertDirect("com.example.mappers.users.insertUser", 1,"hello world",100)
	//   if err != nil {
	//       log.Fatalf("Failed to insert user: %v", err)
	//   }
	InsertDirect(mapperId string, args ...any) (int64, error)

	// UpdateDirect executes an update operation based on the specified XML mapping mapper ID.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   args: Variable length argument list, which corresponds to placeholder arguments of mapperId.
	//
	// Returns:
	//   The number of rows affected by the update operation as an int64, or an error if the operation fails.
	//
	// Description:
	//   This function executes an update operation based on the specified XML mapping mapper ID.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "updateUserByEmail" is the ID of the CRUD operation within the namespace
	//   rowsAffected, err := gdaoMapper.UpdateDirect("com.example.mappers.users.updateUserByEmail","hello world" ,"donnie4w@gmail.com")
	//   if err != nil {
	//       log.Fatalf("Failed to update user: %v", err)
	//   }
	UpdateDirect(mapperId string, args ...any) (int64, error)

	// DeleteDirect executes a delete operation based on the specified XML mapping mapper ID.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   args: Variable length argument list, which corresponds to placeholder arguments of mapperId.
	//
	// Returns:
	//   The number of rows affected by the delete operation as an int64, or an error if the operation fails.
	//
	// Description:
	//   This function executes a delete operation based on the specified XML mapping mapper ID.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "deleteUserById" is the ID of the CRUD operation within the namespace
	//   rowsAffected, err := gdaoMapper.DeleteDirect("com.example.mappers.users.deleteUserById", 1)
	//   if err != nil {
	//       log.Fatalf("Failed to delete user: %v", err)
	//   }
	DeleteDirect(mapperId string, args ...any) (int64, error)

	// SelectBeans executes a query based on the specified XML mapping mapper ID and returns multiple rows of data as DataBeans.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   parameter: The parameter to pass to the query. Can be a basic data type, an entity class object, a map, or a slice, depending on the requirements of the query.
	//
	// Returns:
	//   A slice of pointers to DataBeans containing the data retrieved by the query.
	//
	// Description:
	//   This function executes a query based on the specified XML mapping mapper ID and returns multiple rows of data.
	//   The parameter can be a basic data type, an entity class object, a map, or a slice, depending on the requirements of the query.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "getUsersByLimit" is the ID of the CRUD operation within the namespace
	//   usersBeans, err := gdaoMapper.SelectBeans("com.example.mappers.users.getUsersByLimit", 10)
	//   if err != nil {
	//       log.Fatalf("Failed to select users: %v", err)
	//   }
	//
	SelectBeans(mapperId string, parameter any) *base.DataBeans

	// SelectBean executes a query based on the specified XML mapping mapper ID and returns a single row of data as a DataBean.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   parameter: The parameter to pass to the query. Can be a basic data type, an entity class object, a map, or a slice, depending on the requirements of the query.
	//
	// Returns:
	//   A pointer to a DataBean containing the data retrieved by the query.
	//
	// Description:
	//   This function executes a query based on the specified XML mapping mapper ID and returns a single row of data.
	//   The parameter can be a basic data type, an entity class object, a map, or a slice, depending on the requirements of the query.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "getUserById" is the ID of the CRUD operation within the namespace
	//   userBean, err := gdaoMapper.SelectBean("com.example.mappers.users.getUserById", 1)
	//   if err != nil {
	//       log.Fatalf("Failed to select user: %v", err)
	//   }
	SelectBean(mapperId string, parameter any) *base.DataBean

	// Insert executes an insert operation based on the specified XML mapping mapper ID.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   parameter: The parameter to pass to the insert operation. Can be a basic data type, an entity class object, a map, or a slice, depending on the requirements of the insert operation.
	//
	// Returns:
	//   The number of rows affected by the insert operation as an int64, or an error if the operation fails.
	//
	// Description:
	//   This function executes an insert operation based on the specified XML mapping mapper ID.
	//   The parameter can be a basic data type, an entity class object, a map, or a slice, depending on the requirements of the insert operation.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "insertUser" is the ID of the CRUD operation within the namespace
	//   // And "newUser" is the user object to insert
	//   rowsAffected, err := gdaoMapper.Insert("com.example.mappers.users.insertUser", newUser)
	//   if err != nil {
	//       log.Fatalf("Failed to insert user: %v", err)
	//   }
	Insert(mapperId string, parameter any) (int64, error)

	// Update executes an update operation based on the specified XML mapping mapper ID.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   parameter: The parameter to pass to the update operation. Can be a basic data type, an entity class object, a map, or a slice, depending on the requirements of the update operation.
	//
	// Returns:
	//   The number of rows affected by the update operation as an int64, or an error if the operation fails.
	//
	// Description:
	//   This function executes an update operation based on the specified XML mapping mapper ID.
	//   The parameter can be a basic data type, an entity class object, a map, or a slice, depending on the requirements of the update operation.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "updateUserByEmail" is the ID of the CRUD operation within the namespace
	//   // And "userToUpdate" is the user object with updated fields
	//   rowsAffected, err := gdaoMapper.Update("com.example.mappers.users.updateUserByEmail", userToUpdate)
	//   if err != nil {
	//       log.Fatalf("Failed to update user: %v", err)
	//   }
	Update(mapperId string, parameter any) (int64, error)

	// Delete executes a delete operation based on the specified XML mapping mapper ID.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   parameter: The parameter to pass to the delete operation. Can be a basic data type, an entity class object, a map, or a slice, depending on the requirements of the delete operation.
	//
	// Returns:
	//   The number of rows affected by the delete operation as an int64, or an error if the operation fails.
	//
	// Description:
	//   This function executes a delete operation based on the specified XML mapping mapper ID.
	//   The parameter can be a basic data type, an entity class object, a map, or a slice, depending on the requirements of the delete operation.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "deleteUserById" is the ID of the CRUD operation within the namespace
	//   rowsAffected, err := gdaoMapper.Delete("com.example.mappers.users:deleteUserById", 1)
	//   if err != nil {
	//       log.Fatalf("Failed to delete user: %v", err)
	//   }
	Delete(mapperId string, parameter any) (int64, error)
}

var (
	IsAutocommit      func() bool
	SetAutocommit     func(autocommit bool) (err error)
	UseTransaction    func(tx base.Transaction)
	Rollback          func() (err error)
	Commit            func() (err error)
	UseDBhandle       func(dbhandler base.DBhandle)
	UseDBhandleWithDB func(db *sql.DB, dbType base.DBType)

	// SelectBeanDirect executes a query based on the specified XML mapping mapper ID and returns a single row of data as a DataBean.
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
	SelectBeanDirect func(mapperId string, args ...any) *base.DataBean

	SelectBean func(mapperId string, parameter any) *base.DataBean

	// SelectBeansDirect executes a query based on the specified XML mapping mapper ID and returns multiple rows of data as DataBeans.
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
	SelectBeansDirect func(mapperId string, args ...any) *base.DataBeans

	SelectBeans func(mapperId string, parameter any) *base.DataBeans

	// InsertDirect executes an insert operation based on the specified XML mapping mapper ID.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   args: Variable length argument list, which corresponds to placeholder arguments of mapperId.
	//
	// Returns:
	//   The number of rows affected by the insert operation as an int64, or an error if the operation fails.
	//
	// Description:
	//   This function executes an insert operation based on the specified XML mapping mapper ID.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "insertUser" is the ID of the CRUD operation within the namespace
	//   // And "newUser" is the user object to insert
	//   rowsAffected, err := gdaoMapper.InsertDirect("com.example.mappers.users.insertUser", newUser)
	//   if err != nil {
	//       log.Fatalf("Failed to insert user: %v", err)
	//   }
	InsertDirect func(mapperId string, args ...any) (r int64, err error)

	// UpdateDirect executes an update operation based on the specified XML mapping mapper ID.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   args: Variable length argument list, which corresponds to placeholder arguments of mapperId.
	//
	// Returns:
	//   The number of rows affected by the update operation as an int64, or an error if the operation fails.
	//
	// Description:
	//   This function executes an update operation based on the specified XML mapping mapper ID.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "updateUser" is the ID of the CRUD operation within the namespace
	//   rowsAffected, err := gdaoMapper.UpdateDirect("com.example.mappers.users.updateUser", 10)
	//   if err != nil {
	//       log.Fatalf("Failed to insert user: %v", err)
	//   }
	UpdateDirect func(mapperId string, args ...any) (r int64, err error)

	// DeleteDirect executes an delete operation based on the specified XML mapping mapper ID.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   args: Variable length argument list, which corresponds to placeholder arguments of mapperId.
	//
	// Returns:
	//   The number of rows affected by the delete operation as an int64, or an error if the operation fails.
	//
	// Description:
	//   This function executes an delete operation based on the specified XML mapping mapper ID.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "deleteUser" is the ID of the CRUD operation within the namespace
	//   rowsAffected, err := gdaoMapper.DeleteDirect("com.example.mappers.users.deleteUser", 10)
	//   if err != nil {
	//       log.Fatalf("Failed to insert user: %v", err)
	//   }
	DeleteDirect func(mapperId string, args ...any) (r int64, err error)

	// Insert executes an insert operation based on the specified XML mapping mapper ID.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   parameter: The parameter to pass to the query. Can be a basic data type, an entity class object, a map, or a slice.
	//
	// Returns:
	//   The number of rows affected by the insert operation as an int64, or an error if the operation fails.
	//
	// Description:
	//   This function executes an insert operation based on the specified XML mapping mapper ID.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "insertUser" is the ID of the CRUD operation within the namespace
	//   // And "newUser" is the user object to insert
	//   rowsAffected, err := gdaoMapper.Insert("com.example.mappers.users.insertUser", newUser)
	//   if err != nil {
	//       log.Fatalf("Failed to insert user: %v", err)
	//   }
	Insert func(mapperId string, parameter any) (r int64, err error)

	// Update executes an update operation based on the specified XML mapping mapper ID.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   parameter: The parameter to pass to the query. Can be a basic data type, an entity class object, a map, or a slice.
	//
	// Returns:
	//   The number of rows affected by the update operation as an int64, or an error if the operation fails.
	//
	// Description:
	//   This function executes an update operation based on the specified XML mapping mapper ID.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "updateUser" is the ID of the CRUD operation within the namespace
	//   rowsAffected, err := gdaoMapper.Update("com.example.mappers.users.updateUser", []any{"hello world",10})
	//   if err != nil {
	//       log.Fatalf("Failed to insert user: %v", err)
	//   }
	Update func(mapperId string, parameter any) (r int64, err error)

	// Delete executes an delete operation based on the specified XML mapping mapper ID.
	//
	// Parameters:
	//   mapperId: The ID of the CRUD operation within the XML mapping namespace.
	//   parameter: The parameter to pass to the query. Can be a basic data type, an entity class object, a map, or a slice.
	//
	// Returns:
	//   The number of rows affected by the delete operation as an int64, or an error if the operation fails.
	//
	// Description:
	//   This function executes an delete operation based on the specified XML mapping mapper ID.
	//
	// Example:
	//   // Assuming "com.example.mappers.users" is the namespace in the XML mapping files
	//   // And "deleteUser" is the ID of the CRUD operation within the namespace
	//   rowsAffected, err := gdaoMapper.Delete("com.example.mappers.users.deleteUser", 10)
	//   if err != nil {
	//       log.Fatalf("Failed to insert user: %v", err)
	//   }
	Delete func(mapperId string, parameter any) (r int64, err error)
)

// SelectDirect executes a query based on the specified XML mapping mapper ID and returns a single row of data as an instance of the generic type T.
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
func SelectDirect[T any](mapperId string, args ...any) (*T, error) {
	return (*mapperInvoke[T])(defaultMapperHandler).SelectDirect(mapperId, args...)
}

// Select executes a query based on the specified XML mapping mapper ID and returns a single row of data as an instance of the generic type T.
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
func Select[T any](mapperId string, parameter any) (*T, error) {
	return (*mapperInvoke[T])(defaultMapperHandler).Select(mapperId, parameter)
}

// SelectsDirect executes a query based on the specified XML mapping mapper ID and returns multiple rows of data as instances of the generic type T.
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
func SelectsDirect[T any](mapperId string, args ...any) ([]*T, error) {
	return (*mapperInvoke[T])(defaultMapperHandler).SelectsDirect(mapperId, args...)
}

// Selects executes a query based on the specified XML mapping mapper ID and returns multiple rows of data as instances of the generic type T.
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
func Selects[T any](mapperId string, parameter any) ([]*T, error) {
	return (*mapperInvoke[T])(defaultMapperHandler).Selects(mapperId, parameter)
}

//func SelectWithGdaoMapper[T any](gdaomapper GdaoMapper, mapperId string, args ...any) (*T, error) {
//	if v, ok := gdaomapper.(*mapperHandler); ok {
//		return (*mapperInvoke[T])(v).Select(mapperId, args...)
//	}
//	return nil, fmt.Errorf("gdaomapper is not a MapperHandler pointer")
//}

//
//func SelectAnyWithGdaoMapper[T any](gdaomapper GdaoMapper, mapperId string, parameter any) (*T, error) {
//	if v, ok := gdaomapper.(*mapperHandler); ok {
//		return (*mapperInvoke[T])(v).SelectAny(mapperId, parameter)
//	}
//	return nil, fmt.Errorf("gdaomapper is not a MapperHandler pointer")
//}
//
//func SelectsWithGdaoMapper[T any](gdaomapper GdaoMapper, mapperId string, args ...any) ([]*T, error) {
//	if v, ok := gdaomapper.(*mapperHandler); ok {
//		return (*mapperInvoke[T])(v).Selects(mapperId, args...)
//	}
//	return nil, fmt.Errorf("gdaomapper is not a MapperHandler pointer")
//}
//
//func SelectsAnyWithGdaoMapper[T any](gdaomapper GdaoMapper, mapperId string, parameter any) ([]*T, error) {
//	if v, ok := gdaomapper.(*mapperHandler); ok {
//		return (*mapperInvoke[T])(v).SelectsAny(mapperId, parameter)
//	}
//	return nil, fmt.Errorf("gdaomapper is not a MapperHandler pointer")
//}
