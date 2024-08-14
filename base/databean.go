// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package base

import (
	"fmt"
	"github.com/donnie4w/gdao/util"
	"reflect"
	"sort"
	"strings"
	"time"
)

type DataBeans struct {
	Beans []*DataBean
	err   error
}

func (d *DataBeans) Len() int {
	return len(d.Beans)
}

func (d *DataBeans) SetError(err error) {
	d.err = err
}

func (d *DataBeans) GetError() error {
	return d.err
}

func (d *DataBeans) Scan(v any) error {
	if d == nil {
		return nil
	}
	if d.err != nil {
		return d.err
	}
	if len(d.Beans) == 0 {
		return nil
	}
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("DataBeans scan expected a pointer,but got a %s", val.Kind())
	}
	if val.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("DataBeans scan expected a slice,but got a %s", val.Elem().Kind())
	}
	sliceValue := val.Elem()
	elemType := sliceValue.Type().Elem()
	for _, bean := range d.Beans {
		if elemType.Kind() == reflect.Ptr {
			elem := reflect.New(elemType.Elem())
			if err := bean.Scan(elem.Interface()); err != nil {
				return err
			}
			sliceValue.Set(reflect.Append(sliceValue, elem))
		} else {
			elem := reflect.New(elemType).Elem()
			if err := bean.Scan(elem.Addr().Interface()); err != nil {
				return err
			}
			sliceValue.Set(reflect.Append(sliceValue, elem))
		}
	}
	return nil
}

type DataBean struct {
	FieldMapName  map[string]*FieldBeen
	FieldMapIndex map[int]*FieldBeen
	err           error
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

func (g *DataBean) SetError(err error) {
	g.err = err
}

func (g *DataBean) GetError() error {
	return g.err
}

func (g *DataBean) Put(name string, index int, fb *FieldBeen) {
	g.FieldMapName[name] = fb
	g.FieldMapIndex[index] = fb
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

func (g *DataBean) Scan(v any) (err error) {
	//defer util.Recover(&err)
	if g != nil {
		if g.err != nil {
			return g.err
		}
		if scanner, ok := v.(Scanner); ok {
			scanner.ToGdao()
			for _, fieldBean := range g.FieldMapName {
				scanner.Scan(fieldBean.FieldName, fieldBean.Value())
			}
			return nil
			//val := reflect.ValueOf(scanner).Elem()
			//return val.Addr().Interface().(*T), nil
		}
		valptr := reflect.ValueOf(v)

		if valptr.Kind() != reflect.Ptr {
			return fmt.Errorf("DataBean scan expected a pointer,but got a %s", valptr.Type())
		}

		val := valptr.Elem()

		if val.Kind() == reflect.Ptr {
			return fmt.Errorf("DataBean scan expected a pointer,but got a %s", valptr.Type())
		}

		typ := val.Type()
		if typ.Kind() != reflect.Struct {
			return fmt.Errorf("reflect: NumField of non-struct type " + typ.String())
		}
		hasScan := true
		num := typ.NumField()
		if num < g.Len() {
			for i := 0; i < num; i++ {
				field := typ.Field(i)
				fieldName := strings.ToLower(util.DecodeFieldname(field.Name))
				if value := g.ValueByName(fieldName); value != nil {
					ScanValue(val.Field(i), value)
				} else {
					hasScan = false
					break
				}
			}
		}

		if !hasScan || num >= g.Len() {
			for _, fieldBean := range g.FieldMapName {
				fieldName := fieldBean.FieldName
				field := val.FieldByNameFunc(func(s string) bool {
					return strings.EqualFold(s, fieldName)
				})
				if field.IsValid() && field.CanSet() {
					ScanValue(field, fieldBean.Value())
				} else {
					columnName := util.ToUpperFirstLetter(fieldName)
					col := val.MethodByName("Set" + columnName)
					if !col.IsValid() {
						col = valptr.MethodByName("Set" + columnName)
					}
					if col.IsValid() {
						methodType := col.Type()
						numIn := methodType.NumIn()
						if numIn == 1 {
							expectedType := methodType.In(0)
							if val := GetValue(expectedType, fieldBean.Value()); val != nil {
								args := []reflect.Value{reflect.ValueOf(val)}
								col.Call(args)
							}
						}
					} else {
						if Logger.IsVaild {
							Logger.Warn("Failed to assign value to the field [", fieldName, "]")
						}
					}
				}
			}
		}
		return
	}
	return fmt.Errorf("DataBean is nil")
}

func (g *DataBean) ToInt64() (r int64) {
	if g != nil && g.err == nil && g.Len() > 0 {
		r = g.FieldMapIndex[0].ValueInt64()
	}
	return
}

func (g *DataBean) ToUint64() (r uint64) {
	if g != nil && g.err == nil && g.Len() > 0 {
		r = g.FieldMapIndex[0].ValueUint64()
	}
	return
}

func (g *DataBean) ToFloat64() (r float64) {
	if g != nil && g.err == nil && g.Len() > 0 {
		r = g.FieldMapIndex[0].ValueFloat64()
	}
	return
}

func (g *DataBean) ToString() (r string) {
	if g != nil && g.err == nil && g.Len() > 0 {
		r = g.FieldMapIndex[0].ValueString()
	}
	return
}

func (g *DataBean) ToBytes() (r []byte) {
	if g != nil && g.err == nil && g.Len() > 0 {
		r = g.FieldMapIndex[0].ValueBytes()
	}
	return
}

func (g *DataBean) ToBool() (r bool) {
	if g != nil && g.err == nil && g.Len() > 0 {
		r = g.FieldMapIndex[0].ValueBool()
	}
	return
}

func (g *DataBean) ToTime() (r time.Time) {
	if g != nil && g.err == nil && g.Len() > 0 {
		r = g.FieldMapIndex[0].ValueTime()
	}
	return
}

func (g *DataBean) Len() int {
	return len(g.FieldMapName)
}
