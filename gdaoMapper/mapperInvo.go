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
	"reflect"
	"time"
)

type mapperInvoke[T any] mapperHandler

func (m *mapperInvoke[T]) Select(mapperId string, args ...any) (r *T, er error) {
	pb, _ := mapperparser.getParamBean(mapperId)
	if pb != nil {
		databean, err := (*mapperHandler)(m)._selectBean(mapperId, pb, args...)
		if err != nil || databean == nil {
			return nil, err
		}
		if isDBType(pb.outputType) {
			if r, er = toT[T](databean); r != nil {
				return
			}
		} else {
			if r, er = gdao.Scan[T](databean); er != nil || r == nil {
				r, er = toT[T](databean)
			}
		}
		return
	} else {
		if base.Logger.IsVaild {
			base.Logger.Errorf("Mapper Id not found [%s]", mapperId)
		}
		return nil, fmt.Errorf("Mapper Id not found [%s]", mapperId)
	}
	return
}

func (m *mapperInvoke[T]) SelectAny(mapperId string, parameter any) (r *T, er error) {
	pb, _ := mapperparser.getParamBean(mapperId)
	if pb != nil {
		databean, err := (*mapperHandler)(m).SelectAny(mapperId, parameter)
		if err != nil || databean == nil {
			return nil, err
		}
		if isDBType(pb.outputType) {
			if r, er = toT[T](databean); r != nil {
				return
			}
		} else {
			if r, er = gdao.Scan[T](databean); er != nil {
				r, er = toT[T](databean)
			}
		}
		return
	} else {
		if base.Logger.IsVaild {
			base.Logger.Errorf("Mapper Id not found [%s]", mapperId)
		}
		return nil, fmt.Errorf("Mapper Id not found [%s]", mapperId)
	}
	return
}

func (m *mapperInvoke[T]) Selects(mapperId string, args ...any) (r []*T, er error) {
	pb, _ := mapperparser.getParamBean(mapperId)
	if pb != nil {
		databeans, err := (*mapperHandler)(m)._selectsBean(mapperId, pb, args...)
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
		return
	} else {
		if base.Logger.IsVaild {
			base.Logger.Errorf("Mapper Id not found [%s]", mapperId)
		}
		return nil, fmt.Errorf("Mapper Id not found [%s]", mapperId)
	}
	return
}

func (m *mapperInvoke[T]) SelectsAny(mapperId string, parameter any) (r []*T, er error) {
	pb, _ := mapperparser.getParamBean(mapperId)
	if pb != nil {
		databeans, err := (*mapperHandler)(m).SelectsAny(mapperId, parameter)
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
		return
	} else {
		if base.Logger.IsVaild {
			base.Logger.Errorf("Mapper Id not found [%s]", mapperId)
		}
		return nil, fmt.Errorf("Mapper Id not found [%s]", mapperId)
	}

}

func toT[T any](databean *base.DataBean) (r *T, err error) {
	defer base.Recover(&err)
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
