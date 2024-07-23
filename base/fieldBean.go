// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package base

import (
	"fmt"
	"time"
)

type FieldBeen struct {
	FieldName  string
	FieldIndex int
	FieldValue *any
}

func (f *FieldBeen) Value() (r any) {
	if f.FieldValue != nil {
		r = *f.FieldValue
	}
	return
}

func (f *FieldBeen) ValueTime() (t time.Time) {
	if f.FieldValue != nil {
		t, _ = AsTime(*f.FieldValue)
	}
	return
}

func (f *FieldBeen) ValueString() string {
	if f.FieldValue == nil {
		return ""
	}
	v := *f.FieldValue
	switch v.(type) {
	case []uint8:
		return string(v.([]uint8))
	case string:
		return v.(string)
	}
	return fmt.Sprint(v)
}

func (f *FieldBeen) ValueInt64() int64 {
	if f.FieldValue == nil {
		return 0
	}
	return AsInt64(*f.FieldValue)
}

func (f *FieldBeen) ValueInt32() int32 {
	return int32(f.ValueInt64())
}

func (f *FieldBeen) ValueInt16() int16 {
	return int16(f.ValueInt64())
}

func (f *FieldBeen) ValueUint64() uint64 {
	if f.FieldValue == nil {
		return 0
	}
	return AsUint64(*f.FieldValue)
}

func (f *FieldBeen) ValueBytes() []byte {
	if f.FieldValue == nil {
		return nil
	}
	return AsBytes(*f.FieldValue)
}

func (f *FieldBeen) ValueUint32() uint32 {
	return uint32(f.ValueUint64())
}

func (f *FieldBeen) ValueUint16() uint16 {
	return uint16(f.ValueUint64())
}

func (f *FieldBeen) ValueFloat64() float64 {
	if f.FieldValue == nil {
		return 0
	}
	return AsFloat64(*f.FieldValue)
}

func (f *FieldBeen) ValueFloat32() float32 {
	return float32(f.ValueFloat64())
}

func (f *FieldBeen) ValueBool() bool {
	if f.FieldValue == nil {
		return false
	}
	return AsBool(*f.FieldValue)
}

func (f *FieldBeen) Name() string {
	return f.FieldName
}

func (f *FieldBeen) Index() int {
	return f.FieldIndex
}

func (f *FieldBeen) String() string {
	return fmt.Sprint("[", f.FieldName, ":", f.Value(), "]")
}
