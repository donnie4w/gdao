// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package base

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type DBType int8

func AsString(src any) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	case reflect.Slice:
		return string(rv.Bytes())
	}
	return fmt.Sprintf("%v", src)
}

func scanBytes(buf []byte, rv reflect.Value) (b []byte, ok bool) {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.AppendInt(buf, rv.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.AppendUint(buf, rv.Uint(), 10), true
	case reflect.Float32:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 32), true
	case reflect.Float64:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 64), true
	case reflect.Bool:
		return strconv.AppendBool(buf, rv.Bool()), true
	case reflect.String:
		s := rv.String()
		return append(buf, s...), true
	case reflect.Slice:
		return rv.Bytes(), true
	}
	return
}

func AsBytes(v any) (b []byte) {
	switch (v).(type) {
	case []byte:
		return v.([]byte)
	case string:
		return []byte(v.(string))
	}
	return []byte{}
}

func AsBool(v any) bool {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Bool:
		return val.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return val.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return val.Float() != 0
	case reflect.String:
		return val.String() != ""
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan, reflect.Ptr, reflect.Interface, reflect.Func:
		return !val.IsNil()
	default:
		return false
	}
}

func AsByte(v any) byte {
	return byte(AsInt64(v))
}

func AsInt8(v any) int8 {
	return int8(AsInt64(v))
}

func AsUint8(v any) byte {
	return uint8(AsUint64(v))
}

func AsInt16(v any) int16 {
	return int16(AsInt64(v))
}

func AsUint16(v any) uint16 {
	return uint16(AsUint64(v))
}

func AsInt32(v any) int32 {
	return int32(AsInt64(v))
}

func AsUint32(v any) uint32 {
	return uint32(AsUint64(v))
}

func AsInt64(v any) int64 {
	switch v.(type) {
	case int64:
		return v.(int64)
	case int32:
		return int64(v.(int32))
	case int16:
		return int64(v.(int16))
	case int8:
		return int64(v.(int8))
	case uint64:
		return int64(v.(uint64))
	case uint32:
		return int64(v.(uint32))
	case uint16:
		return int64(v.(uint16))
	case uint8:
		return int64(v.(uint8))
	case uint:
		return int64(v.(uint))
	case int:
		return int64(v.(int))
	case float32:
		return int64(v.(float32))
	case float64:
		return int64(v.(float64))
	case []uint8:
		i, _ := strconv.ParseInt(string(v.([]uint8)), 10, 64)
		return i
	case string:
		i, _ := strconv.ParseInt(string(v.(string)), 10, 64)
		return i
	default:
		return 0
	}
}

func AsUint64(v any) uint64 {
	switch v.(type) {
	case int64:
		return uint64(v.(int64))
	case int32:
		return uint64(v.(int32))
	case int16:
		return uint64(v.(int16))
	case int8:
		return uint64(v.(int8))
	case uint64:
		return v.(uint64)
	case uint32:
		return uint64(v.(uint32))
	case uint16:
		return uint64(v.(uint16))
	case uint8:
		return uint64(v.(uint8))
	case uint:
		return uint64(v.(uint))
	case int:
		return uint64(v.(int))
	case float32:
		return uint64(v.(float32))
	case float64:
		return uint64(v.(float64))
	case []uint8:
		i, _ := strconv.ParseUint(string(v.([]uint8)), 10, 64)
		return i
	case string:
		i, _ := strconv.ParseUint(string(v.(string)), 10, 64)
		return i
	default:
		return 0
	}
}

func AsFloat32(v any) float32 {
	return float32(AsFloat64(v))
}

func AsFloat64(v any) float64 {
	switch (v).(type) {
	case int64:
		return float64(v.(int64))
	case int32:
		return float64(v.(int32))
	case int16:
		return float64(v.(int16))
	case int8:
		return float64(v.(int8))
	case uint64:
		return float64(v.(uint64))
	case uint32:
		return float64(v.(uint32))
	case uint16:
		return float64(v.(uint16))
	case uint8:
		return float64(v.(uint8))
	case uint:
		return float64(v.(uint))
	case int:
		return float64(v.(int))
	case float32:
		return float64(v.(float32))
	case float64:
		return float64(v.(float64))
	case []uint8:
		i, _ := strconv.ParseFloat(string(v.([]uint8)), 64)
		return i
	case string:
		i, _ := strconv.ParseFloat(string(v.(string)), 64)
		return i
	default:
		return 0
	}
}

func AsTime(v any) (r time.Time, err error) {
	switch (v).(type) {
	case []uint8:
		return convertToDate(string(v.([]uint8)))
	case string:
		return convertToDate(v.(string))
	case time.Time:
		return v.(time.Time), nil
	case int64, uint64, int, uint, int32, uint32:
		return time.Unix(AsInt64(v), 0), nil
	}
	err = fmt.Errorf("unable to parse the value %v", v)
	return
}

func strconvErr(err error) error {
	if ne, ok := err.(*strconv.NumError); ok {
		return ne.Err
	}
	return err
}

func IsBytes(fieldValueVal reflect.Value) bool {
	return fieldValueVal.Kind() == reflect.Slice && fieldValueVal.Type().Elem().Kind() == reflect.Uint8
}

func ScanValue(desc reflect.Value, src any) {
	if desc.IsValid() && desc.CanSet() {
		fieldValueVal := reflect.ValueOf(src)
		fieldType := desc.Type()
		switch fieldType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			desc.SetInt(AsInt64(src))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			desc.SetUint(AsUint64(src))
		case reflect.Float64, reflect.Float32:
			desc.SetFloat(AsFloat64(src))
		case reflect.Bool:
			if bv, err := driver.Bool.ConvertValue(src); err == nil {
				desc.SetBool(bv.(bool))
			}
		case reflect.String:
			desc.SetString(AsString(src))
		case reflect.Slice:
			if fieldType.Elem().Kind() == reflect.Uint8 {
				switch src.(type) {
				case []uint8:
					desc.Set(reflect.ValueOf((src).([]uint8)))
				case string:
					desc.Set(reflect.ValueOf([]byte((src).(string))))
				default:
					if bs, ok := scanBytes(nil, fieldValueVal); ok {
						desc.Set(reflect.ValueOf(bs))
					}
				}
			}
		case reflect.Struct:
			if fieldType == reflect.TypeOf(time.Time{}) {
				if v, err := AsTime(src); err == nil {
					desc.Set(reflect.ValueOf(v))
				}
			}
		case reflect.Ptr:
			if desc.IsNil() {
				desc.Set(reflect.New(desc.Type().Elem()))
			}
			ScanValue(desc.Elem(), src)
		default:
			if fieldValueVal.Type().ConvertibleTo(desc.Type()) {
				desc.Set(fieldValueVal.Convert(desc.Type()))
			}
		}
	} else {
		if logger.IsVaild {
			logger.Warn("field '", desc, "' not found or cannot be set in type ", src)
		}
	}
}

func GetValue(fieldType reflect.Type, src any) any {
	switch fieldType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return AsInt64(src)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return AsUint64(src)
	case reflect.Float64, reflect.Float32:
		return AsFloat64(src)
	case reflect.Bool:
		if bv, err := driver.Bool.ConvertValue(src); err == nil {
			return bv.(bool)
		}
		return false
	case reflect.String:
		return AsString(src)
	case reflect.Slice:
		fieldValueVal := reflect.ValueOf(src)
		if fieldType.Elem().Kind() == reflect.Uint8 {
			switch src.(type) {
			case []uint8:
				return src.([]uint8)
			case string:
				return []byte((src).(string))
			default:
				if bs, ok := scanBytes(nil, fieldValueVal); ok {
					return bs
				}
			}
		}
	case reflect.Struct:
		if fieldType == reflect.TypeOf(time.Time{}) {
			if v, err := AsTime(src); err == nil {
				return v
			}
		}
	case reflect.Ptr:
		elemType := fieldType.Elem()
		if val := GetValue(elemType, src); val != nil {
			return &val
		}
	default:
		return src
	}
	return nil
}
