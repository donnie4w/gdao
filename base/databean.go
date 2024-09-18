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
	"github.com/donnie4w/gofer/pool/buffer"
	"reflect"
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

// ScanAndFree copies the data from the DataBeans into the provided variable 'v'.
// This method is typically used to transfer data out of the DataBean and into
// another data structure or variable.
//
// After calling this method, the data in the DataBeans will be recycled and
// should not be used anymore. It is the caller's responsibility to ensure that
// the DataBeans is not accessed after calling Scan, as the internal data may
// have been cleared or reused for subsequent operations.
func (d *DataBeans) ScanAndFree(v any) error {
	return d.scan(v, true)
}

// Scan copies the data from the DataBeans into the provided variable 'v'.
// This method is typically used to transfer data out of the DataBean and into
// another data structure or variable.
func (d *DataBeans) Scan(v any) error {
	return d.scan(v, false)
}

func (d *DataBeans) scan(v any, free bool) error {
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
			if err := bean.scan(elem.Interface(), free); err != nil {
				return err
			}
			sliceValue.Set(reflect.Append(sliceValue, elem))
		} else {
			elem := reflect.New(elemType).Elem()
			if err := bean.scan(elem.Addr().Interface(), free); err != nil {
				return err
			}
			sliceValue.Set(reflect.Append(sliceValue, elem))
		}
	}
	return nil
}

var dataBeanPool = buffer.NewPool[DataBean](func() *DataBean {
	databean := new(DataBean)
	return databean
}, func(d *DataBean) {
	d.reset()
})

type DataBean struct {
	fieldMapName  map[string]*FieldBeen
	fieldMapIndex []*FieldBeen
	err           error
}

func NewDataBean(length int) *DataBean {
	r := dataBeanPool.Get()
	if r.fieldMapName == nil {
		r.fieldMapName = make(map[string]*FieldBeen, length)
	}
	return r
}

func (g *DataBean) free() {
	if len(g.fieldMapIndex) > 0 {
		g.fieldMapIndex = g.fieldMapIndex[:0]
	}
	dataBeanPool.Put(&g)
}

func (g *DataBean) reset() {
	defer util.Recover(nil)
	for _, f := range g.fieldMapIndex {
		f.free()
	}
	for k := range g.fieldMapName {
		delete(g.fieldMapName, k)
	}
	g.err = nil
}

func (g *DataBean) FirstField() *FieldBeen {
	if len(g.fieldMapIndex) > 0 {
		return g.fieldMapIndex[0]
	}
	return nil
}

func (g *DataBean) Map() map[string]*FieldBeen {
	return g.fieldMapName
}

func (g *DataBean) SetError(err error) {
	g.err = err
}

func (g *DataBean) GetError() error {
	return g.err
}

func (g *DataBean) Put(name string, fb *FieldBeen) {
	g.fieldMapName[name] = fb
	g.fieldMapIndex = append(g.fieldMapIndex, fb)
}

func (g *DataBean) FieldByName(name string) (_r *FieldBeen) {
	if g.Len() > 0 {
		defer util.Recover(nil)
		_r, _ = g.fieldMapName[name]
	}
	return
}

func (g *DataBean) FieldByIndex(index int) (_r *FieldBeen) {
	if index > 0 && index <= g.Len() {
		return g.fieldMapIndex[index-1]
	}
	return
}

func (g *DataBean) ValueByName(name string) (_val any) {
	if g.Len() > 0 {
		defer util.Recover(nil)
		if f, ok := g.fieldMapName[name]; ok {
			return f.Value()
		}
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
	sb := []string{}
	for _, fm := range g.fieldMapIndex {
		sb = append(sb, fmt.Sprint(fm.Value()))
	}
	return strings.Join(sb, ",")
}

// ScanAndFree copies the data from the DataBean into the provided variable 'v'.
// This method is typically used to transfer data out of the DataBean and into
// another data structure or variable.
//
// After calling this method, the data in the DataBean will be recycled and
// should not be used anymore. It is the caller's responsibility to ensure that
// the DataBean is not accessed after calling Scan, as the internal data may
// have been cleared or reused for subsequent operations.
func (g *DataBean) ScanAndFree(v any) (err error) {
	return g.scan(v, true)
}

// Scan copies the data from the DataBean into the provided variable 'v'.
// This method is typically used to transfer data out of the DataBean and into
// another data structure or variable.
func (g *DataBean) Scan(v any) (err error) {
	return g.scan(v, false)
}

func (g *DataBean) scan(v any, free bool) (err error) {
	//defer util.Recover(&err)
	if g != nil {
		if g.err != nil {
			return g.err
		}

		//if typ, ok := v.(reflect.Type); ok {
		//	elem := reflect.New(typ).Elem()
		//	v = elem.Interface()
		//}

		if scanner, ok := v.(Scanner); ok {
			scanner.ToGdao()
			for name, fieldBean := range g.fieldMapName {
				scanner.Scan(name, fieldBean.Value())
			}
			if free {
				g.free()
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
				//fieldName := strings.ToLower(util.DecodeFieldname(field.Name))
				if value := g.ValueByName(field.Name); value != nil {
					ScanValue(val.Field(i), value)
				} else {
					hasScan = false
					break
				}
			}
		}

		if !hasScan || num >= g.Len() {
			for fieldName, fieldBean := range g.fieldMapName {
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
		if free {
			g.free()
		}
		return
	}
	return fmt.Errorf("DataBean is nil")
}

func (g *DataBean) ToInt64() (r int64) {
	if g != nil && g.err == nil && g.Len() > 0 {
		r = g.FirstField().ValueInt64()
	}
	return
}

func (g *DataBean) ToUint64() (r uint64) {
	if g != nil && g.err == nil && g.Len() > 0 {
		r = g.FirstField().ValueUint64()
	}
	return
}

func (g *DataBean) ToFloat64() (r float64) {
	if g != nil && g.err == nil && g.Len() > 0 {
		r = g.FirstField().ValueFloat64()
	}
	return
}

func (g *DataBean) ToString() (r string) {
	if g != nil && g.err == nil && g.Len() > 0 {
		r = g.FirstField().ValueString()
	}
	return
}

func (g *DataBean) ToBytes() (r []byte) {
	if g != nil && g.err == nil && g.Len() > 0 {
		r = g.FirstField().ValueBytes()
	}
	return
}

func (g *DataBean) ToBool() (r bool) {
	if g != nil && g.err == nil && g.Len() > 0 {
		r = g.FirstField().ValueBool()
	}
	return
}

func (g *DataBean) ToTime() (r time.Time) {
	if g != nil && g.err == nil && g.Len() > 0 {
		r = g.FirstField().ValueTime()
	}
	return
}

func (g *DataBean) Len() int {
	return len(g.fieldMapIndex)
}
