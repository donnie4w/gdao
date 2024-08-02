// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoMapper

import (
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	"github.com/donnie4w/gdao/gdaoCache"
	"reflect"
	"time"
)

type mapperInvoke[T any] mapperHandler

func (m *mapperInvoke[T]) Select(mapperId string, args ...any) (r *T, er error) {
	if pb, ok := mapperparser.getParamBean(mapperId); !ok {
		return nil, fmt.Errorf("Mapper Id not found [%s]", mapperId)
	} else {
		if base.Logger.IsVaild {
			base.Logger.Debug("[Mapper Id] "+mapperId+" \nSelect SQL["+pb.sql+"]ARGS", args)
		}
		return _select[T]((*mapperHandler)(m), pb, args...)
	}
}

func (m *mapperInvoke[T]) SelectAny(mapperId string, parameter any) (r *T, err error) {
	var pb *paramBean
	var args []any
	mh := (*mapperHandler)(m)
	if pb, args, err = mh.parseParameter(mapperId, parameter); err != nil {
		return r, err
	}
	if base.Logger.IsVaild {
		base.Logger.Debug("[Mapper Id] "+mapperId+" \nSelectAny SQL["+pb.sql+"]ARGS", args)
	}
	return _select[T](mh, pb, args...)
}

func _select[T any](mh *mapperHandler, pb *paramBean, args ...any) (r *T, err error) {
	domain := gdaoCache.GetMapperDomain(pb.namespace, pb.id)
	isCache := domain != ""
	var condition *gdaoCache.Condition
	if isCache {
		condition = gdaoCache.NewCondition("*"+base.Classname[T](), pb.sql, args...)
		if result := gdaoCache.GetMapperCache(domain, pb.namespace, pb.id, condition); result != nil {
			if base.Logger.IsVaild {
				base.Logger.Debug("[GET CACHE]["+pb.sql+"]", args)
			}
			return result.(*T), nil
		}
	}
	var databean *base.DataBean
	if databean, err = mh.getDBhandle(pb.namespace, pb.id, true).ExecuteQueryBean(pb.sql, args...); err == nil {
		if isDBType(pb.outputType) {
			r, err = toT[T](databean)
		}
		if r == nil {
			if r, err = gdao.Scan[T](databean); err != nil || r == nil {
				r, err = toT[T](databean)
			}
		}
	}

	if isCache && r != nil && err == nil {
		gdaoCache.SetMapperCache(domain, pb.namespace, pb.id, condition, r)
		if base.Logger.IsVaild {
			base.Logger.Debug("[SET CACHE]["+pb.sql+"]", args)
		}
	}
	return
}

func (m *mapperInvoke[T]) Selects(mapperId string, args ...any) (r []*T, er error) {
	if pb, ok := mapperparser.getParamBean(mapperId); !ok {
		return nil, fmt.Errorf("Mapper Id not found [%s]", mapperId)
	} else {
		if base.Logger.IsVaild {
			base.Logger.Debug("[Mapper Id] "+mapperId+" \nSelects SQL["+pb.sql+"]ARGS", args)
		}
		return selects[T]((*mapperHandler)(m), pb, args...)
	}
}

func (m *mapperInvoke[T]) SelectsAny(mapperId string, parameter any) (r []*T, err error) {
	var pb *paramBean
	var args []any
	mh := (*mapperHandler)(m)
	if pb, args, err = mh.parseParameter(mapperId, parameter); err != nil {
		return r, err
	}
	if base.Logger.IsVaild {
		base.Logger.Debug("[Mapper Id] "+mapperId+" \nSelectsAny SQL["+pb.sql+"]ARGS", args)
	}
	return selects[T](mh, pb, args...)
}

func selects[T any](mh *mapperHandler, pb *paramBean, args ...any) (r []*T, err error) {
	domain := gdaoCache.GetMapperDomain(pb.namespace, pb.id)
	isCache := domain != ""
	var condition *gdaoCache.Condition
	if isCache {
		condition = gdaoCache.NewCondition("[]*"+base.Classname[T](), pb.sql, args...)
		if result := gdaoCache.GetMapperCache(domain, pb.namespace, pb.id, condition); result != nil {
			if base.Logger.IsVaild {
				base.Logger.Debug("[GET CACHE]["+pb.sql+"]", args)
			}
			return result.([]*T), nil
		}
	}
	var databeans []*base.DataBean
	if databeans, err = mh.getDBhandle(pb.namespace, pb.id, true).ExecuteQueryBeans(pb.sql, args...); err == nil {
		if err != nil || databeans == nil {
			return nil, err
		}
		r = make([]*T, 0)
		if isDBType(pb.outputType) {
			for _, databean := range databeans {
				if v, _ := toT[T](databean); v != nil {
					r = append(r, v)
				}
			}
		} else {
			ok := true
			for _, databean := range databeans {
				if v, err := gdao.Scan[T](databean); err == nil && v != nil {
					r = append(r, v)
				} else {
					ok = false
					break
				}
			}
			if !ok {
				if len(r) > 0 {
					r = make([]*T, 0)
				}
				for _, databean := range databeans {
					if v, _ := toT[T](databean); v != nil {
						r = append(r, v)
					}
				}
			}
		}
		if err == nil && len(r) > 0 && isCache {
			gdaoCache.SetMapperCache(domain, pb.namespace, pb.id, condition, r)
			if base.Logger.IsVaild {
				base.Logger.Debug("[SET CACHE]["+pb.sql+"]", args)
			}
		}
	}
	return
}

func toT[T any](databean *base.DataBean) (r *T, err error) {
	defer base.Recover(&err)
	if databean == nil {
		return
	}
	var t T
	field := databean.FieldMapIndex[0]
	if field == nil {
		return nil, nil
	}
	var v any = nil
	val := reflect.ValueOf(&t).Elem()
	switch val.Kind() {
	case reflect.Int64, reflect.Int32, reflect.Int, reflect.Int8, reflect.Int16:
		v = field.ValueInt64()
	case reflect.Uint64, reflect.Uint32, reflect.Uint, reflect.Uint8, reflect.Uint16:
		v = field.ValueUint64()
	case reflect.Float32, reflect.Float64:
		v = field.ValueFloat64()
	case reflect.Bool:
		v = field.ValueBool()
	case reflect.String:
		v = field.ValueString()
	case reflect.Slice:
		if val.Type().Elem().Kind() == reflect.Uint8 {
			v = field.ValueBytes()
		}
	}

	if v == nil && val.Type() == reflect.TypeOf(time.Time{}) {
		v = field.ValueTime()
	}

	if v != nil {
		value := reflect.ValueOf(v).Convert(reflect.TypeOf(t))
		result := value.Interface().(T)
		return &result, nil
	} else {
		err = fmt.Errorf("value:%v ,type %v cannot be converted to type %v", field.Value(), reflect.TypeOf(field.Value()), reflect.TypeOf(t))
	}
	return
}
