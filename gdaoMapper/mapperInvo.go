// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoMapper

import (
	"fmt"
	"github.com/donnie4w/gdao/base"
	"github.com/donnie4w/gdao/gdaoCache"
	"github.com/donnie4w/gdao/util"
	"reflect"
	"time"
)

type mapperInvoke[T any] mapperHandler

func (m *mapperInvoke[T]) SelectDirect(mapperId string, args ...any) (r *T, er error) {
	if pb, ok := mapperparser.getParamBean(mapperId); !ok {
		return nil, fmt.Errorf("Mapper Id not found [%s]", mapperId)
	} else {
		if base.Logger.IsVaild {
			base.Logger.Debug("[Mapper Id] "+mapperId+" \nSelectDirect SQL["+pb.sql+"]ARGS", args)
		}
		return _select[T]((*mapperHandler)(m), pb, args...)
	}
}

func (m *mapperInvoke[T]) Select(mapperId string, parameter any) (r *T, err error) {
	var pb *paramBean
	var args []any
	mh := (*mapperHandler)(m)
	if pb, args, err = mh.parseParameter(mapperId, parameter); err != nil {
		return r, err
	}
	if base.Logger.IsVaild {
		base.Logger.Debug("[Mapper Id] "+mapperId+" \nSelect SQL["+pb.sql+"]ARGS", args)
	}
	return _select[T](mh, pb, args...)
}

func _select[T any](mh *mapperHandler, pb *paramBean, args ...any) (r *T, err error) {
	domain := gdaoCache.GetMapperDomain(pb.namespace, pb.id)
	isCache := domain != ""
	var condition *gdaoCache.Condition
	if isCache {
		condition = gdaoCache.NewCondition("*"+util.Classname[T](), pb.sql, args...)
		if result := gdaoCache.GetMapperCache(domain, pb.namespace, pb.id, condition); result != nil {
			if base.Logger.IsVaild {
				base.Logger.Debug("[GET CACHE]["+pb.sql+"]", args)
			}
			return result.(*T), nil
		}
	}
	var databean *base.DataBean
	if databean = mh.getDBhandle(pb.namespace, pb.id, true).ExecuteQueryBean(pb.sql, args...); databean.GetError() == nil {
		if isDBType(pb.outputType) {
			r, err = toT[T](databean)
		}
		if r == nil {
			if err = databean.GetError(); err == nil {
				r = new(T)
				if err = databean.Scan(r); err != nil {
					r, err = toT[T](databean)
				}
			}
		}
	} else {
		err = databean.GetError()
	}

	if isCache && r != nil && err == nil {
		gdaoCache.SetMapperCache(domain, pb.namespace, pb.id, condition, r)
		if base.Logger.IsVaild {
			base.Logger.Debug("[SET CACHE]["+pb.sql+"]", args)
		}
	}
	return
}

func (m *mapperInvoke[T]) SelectsDirect(mapperId string, args ...any) (r []*T, er error) {
	if pb, ok := mapperparser.getParamBean(mapperId); !ok {
		return nil, fmt.Errorf("Mapper Id not found [%s]", mapperId)
	} else {
		if base.Logger.IsVaild {
			base.Logger.Debug("[Mapper Id] "+mapperId+" \nSelectsDirect SQL["+pb.sql+"]ARGS", args)
		}
		return selects[T]((*mapperHandler)(m), pb, args...)
	}
}

func (m *mapperInvoke[T]) Selects(mapperId string, parameter any) (r []*T, err error) {
	var pb *paramBean
	var args []any
	mh := (*mapperHandler)(m)
	if pb, args, err = mh.parseParameter(mapperId, parameter); err != nil {
		return r, err
	}
	if base.Logger.IsVaild {
		base.Logger.Debug("[Mapper Id] "+mapperId+" \nSelects SQL["+pb.sql+"]ARGS", args)
	}
	return selects[T](mh, pb, args...)
}

func selects[T any](mh *mapperHandler, pb *paramBean, args ...any) (r []*T, err error) {
	domain := gdaoCache.GetMapperDomain(pb.namespace, pb.id)
	isCache := domain != ""
	var condition *gdaoCache.Condition
	if isCache {
		condition = gdaoCache.NewCondition("[]*"+util.Classname[T](), pb.sql, args...)
		if result := gdaoCache.GetMapperCache(domain, pb.namespace, pb.id, condition); result != nil {
			if base.Logger.IsVaild {
				base.Logger.Debug("[GET CACHE]["+pb.sql+"]", args)
			}
			return result.([]*T), nil
		}
	}
	if databeans := mh.getDBhandle(pb.namespace, pb.id, true).ExecuteQueryBeans(pb.sql, args...); databeans.GetError() == nil && databeans.Len() > 0 {
		r = make([]*T, 0)
		if isDBType(pb.outputType) {
			for _, databean := range databeans.Beans {
				if v, _ := toT[T](databean); v != nil {
					r = append(r, v)
				}
			}
		} else {
			ok := true
			for _, databean := range databeans.Beans {
				v := new(T)
				if err := databean.Scan(v); err == nil {
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
				for _, databean := range databeans.Beans {
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
	} else {
		err = databeans.GetError()
	}
	return
}

func toT[T any](databean *base.DataBean) (r *T, err error) {
	defer util.Recover(&err)
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
