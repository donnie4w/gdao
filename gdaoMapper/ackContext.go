// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoMapper

import (
	"fmt"
	. "github.com/donnie4w/gdao/gdaoStruct"
	"github.com/donnie4w/gdao/util"
	"regexp"
	"strings"
)

type ackContext struct {
	sqlbuilder *strings.Builder
	params     []any
}

func newAckContext(sqlbuilder *strings.Builder, params []any) *ackContext {
	return &ackContext{
		sqlbuilder: sqlbuilder,
		params:     params,
	}
}

func newAckContext2(text string, paramContext *ParamContext) (r *ackContext, err error) {
	r = &ackContext{}
	if text == "" {
		return
	}
	if list := r.parseSqltext(text); len(list) > 0 {
		params := []any{}
		for _, name := range list {
			v, _ := util.ResolveVariableValue(name, paramContext)
			if v != nil {
				params = append(params, v)
			} else {
				return nil, fmt.Errorf("parse text [%s] failed! The parameter value for [%s] could not be found", text, name)
			}
		}
		if len(params) == 0 && paramContext.GetArray() != nil && len(paramContext.GetArray()) >= len(list) {
			for i := 0; i < len(list); i++ {
				params = append(params, paramContext.GetArray()[i])
			}
		}
		r.params = params
	}
	return
}

func (ack *ackContext) parseSqltext(sql string) []string {
	pattern := regexp.MustCompile(`#\{([^}]+)\}`)
	matches := pattern.FindAllStringSubmatch(sql, -1)
	modifiedSql := pattern.ReplaceAllString(sql, "?")

	var list []string
	for _, match := range matches {
		list = append(list, match[1])
	}
	ack.sqlbuilder = &strings.Builder{}
	ack.sqlbuilder.WriteString(modifiedSql)
	return list
}

func (ack *ackContext) Append(ackContext *ackContext) {
	if ackContext == nil {
		return
	}

	if ack.sqlbuilder != nil && ack.sqlbuilder.Len() > 0 {
		if ackContext.sqlbuilder != nil && ackContext.sqlbuilder.Len() > 0 {
			ack.sqlbuilder.WriteString(" ")
			ack.sqlbuilder.WriteString(ackContext.sqlbuilder.String())
		}
	} else {
		ack.sqlbuilder = ackContext.sqlbuilder
	}

	if ack.params != nil {
		if ackContext.params != nil {
			ack.params = append(ack.params, ackContext.params...)
		}
	} else {
		ack.params = ackContext.params
	}
}

func (ack *ackContext) GetSqlbuilder() *strings.Builder {
	return ack.sqlbuilder
}

func (ack *ackContext) GetSql() string {
	if ack.sqlbuilder == nil {
		return ""
	}
	return ack.sqlbuilder.String()
}

func (ack *ackContext) GetParams() []any {
	return ack.params
}
