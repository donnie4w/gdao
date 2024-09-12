// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package sqlBuilder

import (
	"fmt"
	"testing"
)

func Test_AppendIf(t *testing.T) {
	context := map[string]any{
		"username": "John Doe",
		"age":      30,
		"emails":   []any{"john@example.com", "doe@example.com"},
	}
	builder := NewSqlBuilder()
	builder.Append("SELECT * FROM users where").
		AppendIf("age>0", context, "age=?", context["age"]).
		Append("ORDER BY id ASC")
	fmt.Println(builder.GetSql())
	fmt.Println(builder.GetParameters())
}

func Test_AppendTrim(t *testing.T) {
	context := map[string]any{
		"username": "John Doe",
		"age":      30,
		"emails":   []any{"john@example.com", "doe@example.com"},
	}
	builder := NewSqlBuilder()
	builder.Append("SELECT * FROM users").
		AppendTrim("WHERE ", "", "AND", "", func(trimBuilder SqlBuilder) {
			trimBuilder.AppendIf("username != nil", context, "AND username = ?", context["username"]).
				AppendIf("age > 18", context, "AND age = ?", context["age"])
		}).
		Append("ORDER BY id ASC")
	fmt.Println(builder.GetSql())
	fmt.Println(builder.GetParameters())
}

func Test_AppendChoose(t *testing.T) {
	context := map[string]any{
		"username": "John Doe",
		"age":      30,
	}
	builder := NewSqlBuilder()
	builder.Append("SELECT * FROM users where 1=1").
		AppendChoose(context, func(chooseBuilder ChooseBuilder) {
			chooseBuilder.When("age > 38", "AND age > ?", context["age"])
			chooseBuilder.When("age > 19", "AND age >= ?", context["age"])
			chooseBuilder.When("username != nil", "AND username = ?", context["username"])
			chooseBuilder.When("username != nil || age>0 ", "AND username = ?", context["username"])
			chooseBuilder.Otherwise("AND email IS NULL")
		})
	fmt.Println(builder.GetSql())
	fmt.Println(builder.GetParameters())
}

func Test_AppendForeach(t *testing.T) {
	context := map[string]any{
		"username": "John Doe",
		"age":      30,
		"emails":   []any{"john@example.com", "doe@example.com"},
	}
	builder := NewSqlBuilder()
	builder.Append("SELECT * FROM users").
		Append("where email in").
		AppendForeach("emails", context, "email", ",", "(", ")", func(foreach ForeachBuilder) {
			foreach.Body("?")
		}).
		Append("ORDER BY id ASC")
	fmt.Println(builder.GetSql())
	fmt.Println(builder.GetParameters())
}

func Test_AppendSet(t *testing.T) {
	context := map[string]any{
		"username": "John Doe",
		"age":      30,
	}
	builder := NewSqlBuilder()
	builder.Append("UPDATE users").
		AppendSet(func(setBuilder SqlBuilder) {
			setBuilder.AppendIf("username != nil", context, "username = ?,", context["username"]).
				AppendIf("age == 30", context, "age = ?,", context["age"])
		}).
		Append("WHERE id = ?", 10)
	fmt.Println(builder.GetSql())
	fmt.Println(builder.GetParameters())
}
