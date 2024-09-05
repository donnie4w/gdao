// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoMapper

import (
	. "github.com/donnie4w/gdao/gdaoStruct"
	"strings"
)

type whereNode struct {
	content string
	nodes   []sqlNode
}

func newWhereNode(content string) *whereNode {
	return &whereNode{content: content}
}

func (wn *whereNode) addSqlNode(node sqlNode) {
	wn.nodes = append(wn.nodes, node)
}

func (wn *whereNode) apply(context *ParamContext) *ackContext {
	var contentdb strings.Builder
	if wn.content != "" {
		contentdb.WriteString(wn.content)
	}

	ackContext, err := newAckContext2(contentdb.String(), context)
	if err != nil {
		panic(err)
	}

	for _, sqlNode := range wn.nodes {
		a := sqlNode.apply(context)
		if a != nil {
			ackContext.Append(a)
		}
	}

	sqlbuilder := ackContext.sqlbuilder
	if sqlbuilder == nil || sqlbuilder.Len() == 0 {
		return nil
	}

	sqlstr := strings.TrimSpace(sqlbuilder.String())
	if len(sqlstr) > 3 {
		lowerSql := strings.ToLower(sqlstr[:4])
		if strings.HasPrefix(lowerSql, "and ") {
			sqlstr = sqlstr[4:]
		} else if strings.HasPrefix(lowerSql, "or ") {
			sqlstr = sqlstr[3:]
		}
	}

	hasPrefix := false
	if len(sqlstr) > 5 {
		lowerSql := strings.ToLower(sqlstr[:5])
		if strings.HasPrefix(lowerSql, "where ") {
			hasPrefix = true
		}
	}

	if !hasPrefix {
		sqlstr = "WHERE " + sqlstr
	}

	ackContext.sqlbuilder.Reset()
	ackContext.sqlbuilder.WriteString(sqlstr)

	return ackContext
}
