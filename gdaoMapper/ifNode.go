// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoMapper

import (
	"github.com/donnie4w/gdao/gdaoStruct"
	"github.com/donnie4w/gdao/util"
)

type ifNode struct {
	content string
	test    string
	list    []sqlNode
}

func newIfNode(test, content string) *ifNode {
	return &ifNode{content: content, test: test}
}

func (t *ifNode) addSqlNode(node sqlNode) {
	t.list = append(t.list, node)
}

func (t *ifNode) apply(context *gdaoStruct.ParamContext) *ackContext {
	if t.test != "" {
		if b, err := util.EvaluateExpression(t.test, context); err != nil {
			panic("test[" + t.test + "]Express failed:" + err.Error())
		} else if !b {
			return nil
		}
	}
	if ac, err := newAckContext2(t.content, context); err == nil {
		for _, node := range t.list {
			if _ac := node.apply(context); _ac != nil {
				ac.Append(_ac)
			}
		}
		return ac
	} else {
		panic(err)
	}
}
