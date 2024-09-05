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

type sqlNode interface {
	apply(context *ParamContext) *ackContext
	addSqlNode(node sqlNode)
}

type crudNode struct {
	content string
	nodes   []sqlNode
}

func newCrudNode(context string) *crudNode {
	s := &crudNode{content: context}
	return s
}

func (s *crudNode) addSqlNode(node sqlNode) {
	s.nodes = append(s.nodes, node)
}

func (s *crudNode) apply(context *ParamContext) *ackContext {
	ac, err := newAckContext2(s.content, context)
	if err != nil {
		return nil
	}
	for _, node := range s.nodes {
		ac.Append(node.apply(context))
	}
	return ac
}
