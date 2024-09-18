// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package sqlBuilder

import (
	"database/sql"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	"github.com/donnie4w/gdao/util"
	"strings"
)

type sqlBuilder struct {
	sql        strings.Builder
	parameters []any
	dbhandle   base.DBhandle
	tx         base.Transaction
}

func NewSqlBuilder() SqlBuilder {
	return &sqlBuilder{}
}

func (b *sqlBuilder) UseDBhandle(dbhandle base.DBhandle) {
	b.dbhandle = dbhandle
}

func (b *sqlBuilder) UseTransaction(transaction base.Transaction) {
	b.tx = transaction
}

func (b *sqlBuilder) Append(text string, params ...any) SqlBuilder {
	return b.append(text, params...)
}

func (b *sqlBuilder) append(text string, params ...any) *sqlBuilder {
	b.sql.WriteString(" ")
	b.sql.WriteString(text)
	b.sql.WriteString(" ")
	b.addParameters(params...)
	return b
}

func (b *sqlBuilder) AppendIf(expression string, context any, text string, params ...any) SqlBuilder {
	if evaluate(expression, context) {
		b.sql.WriteString(" ")
		b.sql.WriteString(text)
		b.addParameters(params...)
	}
	return b
}

func (b *sqlBuilder) AppendChoose(context any, chooseBuilderConsumer func(ChooseBuilder)) SqlBuilder {
	chooseBuilder := newChooseBuilder(b, context)
	chooseBuilderConsumer(chooseBuilder)
	return b
}

func (b *sqlBuilder) AppendForeach(collectionName string, context any, item, separator, open, close string, foreachConsumer func(ForeachBuilder)) SqlBuilder {
	var collectionObj any
	if collectionName != "" && collectionName != "list" && collectionName != "array" {
		switch m := context.(type) {
		case map[string]any:
			if v, ok := m[collectionName]; ok {
				collectionObj = v
			}
		}
	} else {
		collectionObj = context
	}
	if collectionObj == nil {
		if base.Logger.IsVaild {
			base.Logger.Warn("AppendForeach unable to find collection data")
		}
		return b
	}

	collection := util.ToArray(collectionObj)
	if len(collection) > 0 {
		if open != "" {
			b.sql.WriteString(open)
		}

		foreachBuilder := newForeachBuilder(b, separator)
		foreachConsumer(foreachBuilder)

		for i, currentItem := range collection {
			replacedBody := strings.Replace(foreachBuilder.body, "#{"+item+"}", "?", -1)
			b.sql.WriteString(replacedBody)
			b.addParameter(currentItem)
			if separator != "" && i < len(collection)-1 {
				b.sql.WriteString(separator)
			}
		}

		if close != "" {
			b.sql.WriteString(close)
		}
	}

	return b
}

func (b *sqlBuilder) AppendTrim(prefix, suffix, prefixOverrides, suffixOverrides string, contentBuilder func(SqlBuilder)) SqlBuilder {
	tempBuilder := &sqlBuilder{}
	contentBuilder(tempBuilder)
	tempSql := strings.TrimSpace(tempBuilder.sql.String())

	if prefixOverrides != "" {
		prefixes := strings.Split(prefixOverrides, "|")
		for _, override := range prefixes {
			if strings.HasPrefix(tempSql, override) {
				tempSql = tempSql[len(override):]
				break
			}
		}
	}

	if suffixOverrides != "" {
		suffixes := strings.Split(suffixOverrides, "|")
		for _, override := range suffixes {
			if strings.HasSuffix(tempSql, override) {
				tempSql = tempSql[:len(tempSql)-len(override)]
			}
		}
	}

	if len(tempSql) > 0 {
		if prefix != "" {
			tempSql = prefix + tempSql
		}
		if suffix != "" {
			tempSql += suffix
		}
		b.sql.WriteString(tempSql + " ")
	}
	b.addParameters(tempBuilder.parameters...)
	return b
}

func (b *sqlBuilder) AppendSet(contentBuilder func(SqlBuilder)) SqlBuilder {
	tempBuilder := &sqlBuilder{}
	contentBuilder(tempBuilder)
	tempSql := strings.TrimRight(tempBuilder.sql.String(), ", ")
	if len(tempSql) > 0 {
		b.sql.WriteString("SET ")
		b.sql.WriteString(tempSql)
		b.sql.WriteString(" ")
	}
	b.addParameters(tempBuilder.parameters...)
	return b
}

func (b *sqlBuilder) GetSql() string {
	return b.sql.String()
}

func (b *sqlBuilder) GetParameters() []any {
	return b.parameters
}

func (b *sqlBuilder) addParameter(param any) {
	b.parameters = append(b.parameters, param)
}

func (b *sqlBuilder) addParameters(params ...any) {
	if len(params) > 0 {
		b.parameters = append(b.parameters, params...)
	}
}

func (b *sqlBuilder) SelectOne() *base.DataBean {
	if base.Logger.IsVaild {
		base.Logger.Debug("[SqlBuilder SQL]", b.GetSql(), "[ARGS]", b.GetParameters())
	}
	if b.dbhandle != nil {
		return b.dbhandle.ExecuteQueryBean(b.GetSql(), b.GetParameters()...)
	}
	return gdao.ExecuteQueryBean(b.GetSql(), b.GetParameters()...)
}

func (b *sqlBuilder) SelectList() *base.DataBeans {
	if base.Logger.IsVaild {
		base.Logger.Debug("[SqlBuilder SQL]", b.GetSql(), "[ARGS]", b.GetParameters())
	}
	return b.getDBHandle().ExecuteQueryBeans(b.GetSql(), b.GetParameters()...)
}

func (b *sqlBuilder) Exec() (sql.Result, error) {
	if base.Logger.IsVaild {
		base.Logger.Debug("[SqlBuilder SQL]", b.GetSql(), "[ARGS]", b.GetParameters())
	}
	return b.getDBHandle().ExecuteUpdate(b.GetSql(), b.GetParameters()...)
}

func (b *sqlBuilder) getDBHandle() (r base.DBhandle) {
	if r = b.tx; r != nil {
		return
	}
	if r = b.dbhandle; r != nil {
		return
	}
	if r = gdao.GetDefaultDBHandle(); r != nil {
		return
	}
	panic("no datasource handle found")
}

type chooseBuilder struct {
	parentBuilder *sqlBuilder
	context       any
	conditionMet  bool
}

func newChooseBuilder(parentBuilder *sqlBuilder, context any) *chooseBuilder {
	return &chooseBuilder{parentBuilder: parentBuilder, context: context}
}

func (cb *chooseBuilder) When(expression, sql string, params ...any) ChooseBuilder {
	if !cb.conditionMet && evaluate(expression, cb.context) {
		cb.parentBuilder.append(sql).addParameters(params...)
		cb.conditionMet = true
	}
	return cb
}

func (cb *chooseBuilder) Otherwise(sql string, params ...any) ChooseBuilder {
	if !cb.conditionMet {
		cb.parentBuilder.append(sql).addParameters(params...)
	}
	return cb
}

type foreachBuilder struct {
	parentBuilder *sqlBuilder
	separator     string
	body          string
}

func newForeachBuilder(parentBuilder *sqlBuilder, separator string) *foreachBuilder {
	return &foreachBuilder{parentBuilder: parentBuilder, separator: separator}
}

func (fb *foreachBuilder) Body(body string) ForeachBuilder {
	fb.body = body
	return fb
}

func evaluate(expression string, context any) bool {
	if b, err := util.EvaluateWithAny(expression, context); err == nil {
		return b
	} else {
		panic(err)
	}
}
