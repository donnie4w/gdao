// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package sqlBuilder

import (
	"github.com/donnie4w/gdao/base"
)

// SqlBuilder is an interface used for building dynamic SQL queries.
type SqlBuilder interface {

	// UseDBhandle sets the database handle to be used by the current SQL builder.
	// This allows the builder to make appropriate adjustments based on different database handles.
	// The parameter dbhandle is an object that implements the base.DBhandle interface.
	UseDBhandle(dbhandle base.DBhandle)

	// UseTransaction initiates a database transaction and takes an object that implements the base.Transaction interface.
	// This function is primarily used to execute a series of database operations within a transactional context.
	// If all operations complete successfully, the entire transaction will be committed, making the changes persistent in the database.
	// If any operation fails, the entire transaction will be rolled back, undoing all changes made within that transaction.
	// This function is typically used in scenarios where ACID (Atomicity, Consistency, Isolation, Durability) properties need to be ensured.
	// Parameters:
	//   tx: An object that implements the base.Transaction interface, providing methods to start, commit, and rollback a transaction.
	UseTransaction(transaction base.Transaction)

	// Append appends a piece of text to the current SQL statement.
	// The parameter text is the string to append.
	// The parameter params is a variadic list of values that may be needed for subsequent parameters.
	// Returns the SqlBuilder instance itself, supporting method chaining.
	Append(text string, params ...any) SqlBuilder

	// AppendIf conditionally appends a piece of text to the current SQL statement.
	// The parameter expression is a boolean expression string used to determine whether to append the text.
	// The parameter context is an object used to evaluate the expression.
	// The parameter text is the string to append if the condition is met.
	// The parameter params is a variadic list of values that may be needed for subsequent parameters.
	// Returns the SqlBuilder instance itself, supporting method chaining.
	AppendIf(expression string, context any, text string, params ...any) SqlBuilder

	// AppendChoose conditionally appends different parts to the current SQL statement based on choices.
	// The parameter context is an object used to evaluate the choice logic.
	// The parameter chooseBuilderConsumer is a function that defines the content of different choice branches.
	// Returns the SqlBuilder instance itself, supporting method chaining.
	AppendChoose(context any, chooseBuilderConsumer func(ChooseBuilder)) SqlBuilder

	// AppendForeach appends a loop structure to the current SQL statement.
	// The parameter collectionName is the name of the variable used for iterating over a collection.
	// The parameter context is an object used to evaluate the loop logic.
	// The parameter item is the name of the variable representing each element within the loop.
	// The parameter separator is the text between loop iterations.
	// The parameter open is the opening text of the loop body.
	// The parameter close is the closing text of the loop body.
	// The parameter foreachConsumer is a function that defines the content of the loop body.
	// Returns the SqlBuilder instance itself, supporting method chaining.
	AppendForeach(collectionName string, context any, item, separator, open, close string, foreachConsumer func(ForeachBuilder)) SqlBuilder

	// AppendTrim appends a trimming structure to the current SQL statement to remove extra prefixes or suffixes.
	// The parameter prefix is the trimming prefix text.
	// The parameter suffix is the trimming suffix text.
	// The parameter prefixOverrides is the text used to override the prefix.
	// The parameter suffixOverrides is the text used to override the suffix.
	// The parameter contentBuilder is a function that defines the trimming content.
	// Returns the SqlBuilder instance itself, supporting method chaining.
	AppendTrim(prefix, suffix, prefixOverrides, suffixOverrides string, contentBuilder func(SqlBuilder)) SqlBuilder

	// AppendSet appends a SET clause to the current SQL statement.
	// The parameter contentBuilder is a function that defines the content of the SET clause.
	// Returns the SqlBuilder instance itself, supporting method chaining.
	AppendSet(contentBuilder func(SqlBuilder)) SqlBuilder

	// GetSql retrieves the final constructed SQL statement.
	// Returns a string representing the complete SQL statement.
	GetSql() string

	// GetParameters retrieves all collected parameters during the construction process.
	// Returns a slice of any type containing all parameters.
	GetParameters() []any

	// SelectOne executes a SQL query and returns the first record.
	// Returns a *base.DataBean object representing the first record of the query result.
	SelectOne() *base.DataBean

	// SelectList executes a SQL query and returns all records.
	// Returns a *base.DataBeans object representing all records of the query result.
	SelectList() *base.DataBeans

	// Exec executes a SQL statement and returns the number of affected rows.
	// Returns an int64 value representing the number of affected rows.
	// If there is an error during execution, it also returns the corresponding error information.
	Exec() (int64, error)
}

type ForeachBuilder interface {
	Body(body string) ForeachBuilder
}

type ChooseBuilder interface {
	When(expression, sql string, params ...any) ChooseBuilder
	Otherwise(sql string, params ...any) ChooseBuilder
}
