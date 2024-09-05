// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoMapper

import (
	. "github.com/donnie4w/gdao/gdaoStruct"
)

type setNode struct {
	nodes []sqlNode
}

func newSetNode() *setNode {
	return &setNode{}
}

func (t *setNode) addSqlNode(node sqlNode) {
	t.nodes = append(t.nodes, node)
}

func (t *setNode) apply(context *ParamContext) *ackContext {
	ac, err := newAckContext2("", context)
	if err != nil {
		panic(err)
	}
	for _, sqlNode := range t.nodes {
		if sqlNode != nil {
			a := sqlNode.apply(context)
			if a != nil {
				ac.Append(a)
			}
		}
	}
	if ac.sqlbuilder != nil {
		s := ac.GetSql()
		ac.sqlbuilder.Reset()
		ac.sqlbuilder.WriteString(" SET ")
		ac.sqlbuilder.WriteString(s)
	}
	return ac
}
