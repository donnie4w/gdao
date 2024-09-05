// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoMapper

import (
	. "github.com/donnie4w/gdao/gdaoStruct"
	"github.com/donnie4w/gdao/util"
)

type ChooseNode struct {
	content string
	test    string
	nodes   []sqlNode
}

func newChooseNode(test, content string) *ChooseNode {
	return &ChooseNode{
		content: content,
		test:    test,
	}
}

func (cn *ChooseNode) addSqlNode(node sqlNode) {
	cn.nodes = append(cn.nodes, node)
}

func (cn *ChooseNode) apply(context *ParamContext) *ackContext {
	if cn.test != "" {
		if b, err := util.EvaluateExpression(cn.test, context); err != nil {
			panic(err)
		} else if !b {
			return nil
		}
	}
	ac, err := newAckContext2(cn.content, context)
	if err != nil {
		panic(err)
	}
	for _, sqlNode := range cn.nodes {
		if sqlNode != nil {
			a := sqlNode.apply(context)
			if a != nil {
				ac.Append(a)
				break
			}
		}
	}
	return ac
}

type whenNode struct {
	content string
	test    string
	nodes   []sqlNode
}

func newWhenNode(test, content string) *whenNode {
	return &whenNode{
		content: content,
		test:    test,
	}
}

func (wn *whenNode) addSqlNode(node sqlNode) {
	wn.nodes = append(wn.nodes, node)
}

func (wn *whenNode) apply(context *ParamContext) *ackContext {
	if wn.test != "" {
		if b, err := util.EvaluateExpression(wn.test, context); err != nil {
			panic(err)
		} else if !b {
			return nil
		}
	}
	ac, err := newAckContext2(wn.content, context)
	if err != nil {
		panic(err)
	}
	for _, sqlNode := range wn.nodes {
		if sqlNode != nil {
			a := sqlNode.apply(context)
			if a != nil {
				ac.Append(a)
			}
		}
	}
	return ac
}

type otherWiseNode struct {
	content string
	test    string
	nodes   []sqlNode
}

func newOtherWiseNode(test, content string) *otherWiseNode {
	return &otherWiseNode{
		content: content,
		test:    test,
	}
}

func (on *otherWiseNode) addSqlNode(node sqlNode) {
	on.nodes = append(on.nodes, node)
}

func (on *otherWiseNode) apply(context *ParamContext) *ackContext {
	if on.test != "" {
		if b, err := util.EvaluateExpression(on.test, context); err != nil {
			panic(err)
		} else if !b {
			return nil
		}
	}
	ac, err := newAckContext2(on.content, context)
	if err != nil {
		panic(err)
	}
	for _, sqlNode := range on.nodes {
		if sqlNode != nil {
			a := sqlNode.apply(context)
			if a != nil {
				ac.Append(a)
			}
		}
	}
	return ac
}
