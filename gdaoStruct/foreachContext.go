// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoStruct

type ForeachContext struct {
	*ParamContext
	itemMap  map[string]any
	indexMap map[string]int
}

func NewForeachContext(pc *ParamContext, arg any) (r *ParamContext) {
	r = NewParamContext3(pc)
	fc := new(ForeachContext)
	fc.ParamContext = NewParamContext(arg)
	r.ForeachContext = fc
	return
}

func NewForeachContext2(pc *ParamContext) (r *ParamContext) {
	r = NewParamContext3(pc)
	fc := new(ForeachContext)
	fc.ParamContext = NewParamContext3(pc)
	r.ForeachContext = fc
	return
}

func (r *ForeachContext) GetItem(itemName string) any {
	if r.itemMap == nil {
		return nil
	}
	return r.itemMap[itemName]
}

func (r *ForeachContext) SetItem(itemName string, itemValue any) {
	if r.itemMap == nil {
		r.itemMap = make(map[string]any)
	}
	r.itemMap[itemName] = itemValue
}

func (r *ForeachContext) GetIndex(indexName string) int {
	if r.indexMap == nil {
		return 0
	}
	return r.indexMap[indexName]
}

func (r *ForeachContext) SetIndex(indexName string, indexValue int) {
	if r.indexMap == nil {
		r.indexMap = make(map[string]int)
	}
	r.indexMap[indexName] = indexValue
}
