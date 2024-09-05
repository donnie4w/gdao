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

type foreach struct {
	Text       string
	Collection string
	Item       string
	Index      string
	Open       string
	Close      string
	Separator  string
	nodes      []sqlNode
}

func newForeach(text, collection, item, index, open, close, separator string) *foreach {
	return &foreach{
		Text:       text,
		Collection: collection,
		Item:       item,
		Index:      index,
		Open:       open,
		Close:      close,
		Separator:  separator,
	}
}

func (f *foreach) addSqlNode(node sqlNode) {
	f.nodes = append(f.nodes, node)
}

func (f *foreach) apply(context *ParamContext) *ackContext {
	var collectionContext *ParamContext
	if f.Collection != "" {
		if f.Collection == "list" || f.Collection == "array" {
			collectionContext = NewForeachContext2(context)
		} else {
			value := context.Get(f.Collection)
			if value != nil {
				collectionContext = NewForeachContext(context, value)
			} else {
				collectionContext = NewForeachContext2(context)
			}
		}
	} else {
		collectionContext = NewForeachContext2(context)
	}

	sqlBuilder := strings.Builder{}
	var params []any

	if array := collectionContext.GetArray(); array != nil {
		for i := 0; i < len(array); i++ {
			if f.Index != "" {
				collectionContext.ForeachContext.SetIndex(f.Index, i)
			}
			if f.Item != "" {
				collectionContext.ForeachContext.SetItem(f.Item, array[i])
			}
			ackContext, err := newAckContext2(f.Text, collectionContext)
			if err != nil {
				panic(err)
			}
			for _, node := range f.nodes {
				ackContext.Append(node.apply(collectionContext))
			}
			if ackContext.sqlbuilder != nil && ackContext.sqlbuilder.Len() > 0 {
				if f.Separator != "" && sqlBuilder.Len() > 0 && i > 0 {
					sqlBuilder.WriteString(f.Separator)
				}
				sqlBuilder.WriteString(ackContext.GetSql())
				if ackContext.GetParams() != nil {
					params = append(params, ackContext.GetParams()...)
				}
			}
		}
	} else {
		i := 0
		for _, value := range collectionContext.GetMap() {
			if f.Index != "" {
				collectionContext.ForeachContext.SetIndex(f.Index, i)
			}
			if f.Item != "" {
				collectionContext.ForeachContext.SetItem(f.Item, value)
			}
			ackContext, err := newAckContext2(f.Text, collectionContext)
			if err != nil {
				panic(err)
			}
			for _, node := range f.nodes {
				ackContext.Append(node.apply(collectionContext))
			}
			if ackContext.sqlbuilder != nil && ackContext.sqlbuilder.Len() > 0 {
				if f.Separator != "" && sqlBuilder.Len() > 0 && i > 0 {
					sqlBuilder.WriteString(f.Separator)
				}
				sqlBuilder.WriteString(ackContext.GetSql())
				if ackContext.GetParams() != nil {
					params = append(params, ackContext.GetParams()...)
				}
			}
			i++
		}
	}

	result := &ackContext{params: params}
	result.sqlbuilder = &strings.Builder{}
	if f.Open != "" {
		result.sqlbuilder.WriteString(f.Open)
		result.sqlbuilder.WriteString(sqlBuilder.String())
	}
	if f.Close != "" {
		result.sqlbuilder.WriteString(f.Close)
	}
	return result
}
