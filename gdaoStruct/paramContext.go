// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoStruct

import "reflect"

type ParamContext struct {
	paramMap   map[string]any
	paramArray []any

	ForeachContext *ForeachContext
}

func NewParamContext(arg any) *ParamContext {
	p := &ParamContext{}
	p.ToParamContext(arg)
	return p
}

func NewParamContext2(arg ...any) *ParamContext {
	p := &ParamContext{}
	p.paramArray = arg
	return p
}

func NewParamContext3(pc *ParamContext) *ParamContext {
	return &ParamContext{paramMap: pc.paramMap, paramArray: pc.paramArray}
}

func (p *ParamContext) Put(key string, value any) {
	if p.paramMap == nil {
		p.paramMap = make(map[string]any)
	}
	p.paramMap[key] = value
}

func (p *ParamContext) Get(key string) any {
	if p.paramMap == nil {
		return nil
	}
	if r, ok := p.paramMap[key]; ok {
		return r
	} else {
		return nil
	}
}

func (p *ParamContext) GetMap() map[string]any {
	if p.ForeachContext != nil {
		return p.ForeachContext.GetMap()
	}
	return p.paramMap
}

func (p *ParamContext) GetArray() []any {
	if p.ForeachContext != nil {
		return p.ForeachContext.GetArray()
	}
	return p.paramArray
}

func (p *ParamContext) SetMap(paramMap map[string]any) {
	p.paramMap = paramMap
}

func (pc *ParamContext) ToParamContext(arg any) {
	switch v := arg.(type) {
	case map[string]any:
		pc.paramMap = v
		return
	case []any:
		pc.paramArray = v
	default:
		value := reflect.ValueOf(arg)
		switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:
			pc.paramArray = []any{arg}
		case reflect.Slice, reflect.Array:
			tempArray := make([]any, value.Len())
			for i := 0; i < value.Len(); i++ {
				tempArray[i] = value.Index(i).Interface()
			}
			pc.paramArray = tempArray
		case reflect.Map:
			pc.paramMap = make(map[string]any)
			for _, key := range value.MapKeys() {
				strKey, ok := key.Interface().(string)
				if ok {
					pc.paramMap[strKey] = value.MapIndex(key).Interface()
				}
			}
		}
	}

	if pc.paramArray == nil && pc.paramMap == nil {
		pc.paramMap = ToMap(arg)
	}
}
