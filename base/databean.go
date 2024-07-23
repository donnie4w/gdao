// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package base

import (
	"fmt"
	"sort"
)

type DataBean struct {
	FieldMapName  map[string]*FieldBeen
	FieldMapIndex map[int]*FieldBeen
	length        int
}

func NewDataBean() *DataBean {
	databean := new(DataBean)
	databean.FieldMapName = make(map[string]*FieldBeen, 0)
	databean.FieldMapIndex = make(map[int]*FieldBeen, 0)
	return databean
}

func (g *DataBean) Map() map[string]*FieldBeen {
	return g.FieldMapName
}

func (g *DataBean) Put(name string, index int, fb *FieldBeen) {
	g.FieldMapName[name] = fb
	g.FieldMapIndex[index] = fb
	g.length++
}

func (g *DataBean) FieldByName(name string) (_r *FieldBeen) {
	_r, _ = g.FieldMapName[name]
	return
}

func (g *DataBean) FieldByIndex(index int) (_r *FieldBeen) {
	_r, _ = g.FieldMapIndex[index]
	return
}

func (g *DataBean) ValueByName(name string) (_val any) {
	if v := g.FieldByName(name); v != nil {
		return v.Value()
	}
	return
}

func (g *DataBean) ValueByIndex(index int) (_val any) {
	if v := g.FieldByIndex(index); v != nil {
		return v.Value()
	}
	return
}

func (g *DataBean) String() (r string) {
	fs := make([]*FieldBeen, 0)
	for _, fm := range g.FieldMapName {
		fs = append(fs, fm)
	}
	sort.Slice(fs, func(i, j int) bool { return fs[i].FieldIndex < fs[j].FieldIndex })
	for _, m := range fs {
		r += fmt.Sprint(m.String())
	}
	return
}

func (g *DataBean) Len() int {
	return g.length
}
