// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoBuilder

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
)

// Build creates a source code string for a standardized gdao entity class.
//
// Parameters:
// - tableName: The name of the database table to query for structure information.
// - dbType: The type of the database, e.g., "mysql", "postgresql", "tidb", "oceanbase", "opengauss".
// - dbName: The name of the database to connect to.
// - packageName: The name of the Go package where the generated entity class will reside.
// - db: An open database connection.
//
// Returns:
// - err: An error if the gdao builder fails, nil otherwise.
//
// Example usage:
// db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/my_database?charset=utf8mb4")
//
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// defer db.Close()
// sourceCode := gdaoBuilder.Build("employees", "mysql", "my_database", "dao", db)
// fmt.Println(sourceCode)
func Build(tableName, dbType, dbName, packageName string, db *sql.DB) (err error) {
	return BuildWithAlias(tableName, tableName, dbType, dbName, packageName, db)
}

// BuildWithAlias creates a source code string for a standardized gdao entity class.
//
// Parameters:
// - tableName: The name of the database table to query for structure information.
// - tableAlias: An alias for the table used when generating the entity class.
// - dbType: The type of the database, e.g., "mysql", "postgresql", "tidb", "oceanbase", "opengauss".
// - dbName: The name of the database to connect to.
// - packageName: The name of the Go package where the generated entity class will reside.
// - db: An open database connection.
//
// Returns:
// - err: An error if the gdao builder fails, nil otherwise.
//
// Example usage:
// db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/my_database?charset=utf8mb4")
//
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// defer db.Close()
// sourceCode := gdaoBuilder.BuildWithAlias("employees", "", "mysql", "my_database", "dao", db)
// fmt.Println(sourceCode)
func BuildWithAlias(tableName, tableAlias, dbType, dbName, packageName string, db *sql.DB) (err error) {
	return BuildDirWithAlias("", tableName, tableAlias, dbType, dbName, packageName, db)
}

// BuildDir creates a source code string for a standardized gdao entity class.
//
// Parameters:
// - dir: Path for storing the generated file.
// - tableName: The name of the database table to query for structure information.
// - dbType: The type of the database, e.g., "mysql", "postgresql", "tidb", "oceanbase", "opengauss".
// - dbName: The name of the database to connect to.
// - packageName: The name of the Go package where the generated entity class will reside.
// - db: An open database connection.
//
// Returns:
// - err: An error if the gdao builder fails, nil otherwise.
//
// Example usage:
// db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/my_database?charset=utf8mb4")
//
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// defer db.Close()
// sourceCode := gdaoBuilder.BuildDir("/usr/local/gdao", "employees", "mysql", "my_database", "dao", db)
// fmt.Println(sourceCode)
func BuildDir(dir, tableName, dbType, dbName, packageName string, db *sql.DB) (err error) {
	return BuildDirWithAlias(dir, tableName, tableName, dbType, dbName, packageName, db)
}

// BuildDirWithAlias creates a source code string for a standardized gdao entity class.
//
// Parameters:
// - dir: Path for storing the generated file.
// - tableName: The name of the database table to query for structure information.
// - tableAlias: An alias for the table used when generating the entity class.
// - dbType: The type of the database, e.g., "mysql", "postgresql", "tidb", "oceanbase", "opengauss".
// - dbName: The name of the database to connect to.
// - packageName: The name of the Go package where the generated entity class will reside.
// - db: An open database connection.
//
// Returns:
// - err: An error if the gdao builder fails, nil otherwise.
//
// Example usage:
// db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/my_database?charset=utf8mb4")
//
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// defer db.Close()
// sourceCode := gdaoBuilder.BuildDirWithAlias("/usr/local/gdao", "employees", "", "mysql", "my_database", "dao", db)
// fmt.Println(sourceCode)
func BuildDirWithAlias(dir, tableName, tableAlias, dbType, dbName, packageName string, db *sql.DB) (err error) {
	var tb *TableBean
	if tb, err = GetTableBean(tableName, db); err == nil {
		if tableAlias == "" {
			tableAlias = tableName
		}
		if structstr := buildstruct(dbType, dbName, tableName, tableAlias, packageName, tb); structstr != "" {
			fileName := filepath.Join(packageName, tableAlias) + ".go"
			if dir != "" {
				fileName = filepath.Join(dir, fileName)
			}
			if err = os.MkdirAll(filepath.Dir(fileName), os.ModePerm); err == nil {
				var f *os.File
				if f, err = os.Create(fileName); err == nil {
					defer f.Close()
					if _, err = f.WriteString(structstr); err == nil {
						log.Println("[successfully created gdao struct]", "[table:", tableName, "]["+fileName+"]")
					}
				}
			}
		}
	}
	if err != nil {
		log.Println("[failed to created gdao struct]", aslog(tableName, tableAlias))
	}
	return
}
