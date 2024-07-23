// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package base

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"
)

type Serialize[T any] interface {
	Encode(T) ([]byte, error)
	Decode([]byte) (T, error)
}

type Serializer struct {
}

func (t *Serializer) Encode(data map[string]interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)

	for key, value := range data {
		if err := binary.Write(buf, binary.LittleEndian, int8(len(key))); err != nil {
			return nil, err
		}
		if _, err := buf.WriteString(key); err != nil {
			return nil, err
		}

		switch v := value.(type) {
		case bool:
			if err := binary.Write(buf, binary.LittleEndian, byte(1)); err != nil {
				return nil, err
			}
			var b byte
			if v {
				b = 1
			} else {
				b = 0
			}
			if err := binary.Write(buf, binary.LittleEndian, b); err != nil {
				return nil, err
			}
		case string:
			if err := binary.Write(buf, binary.LittleEndian, byte(2)); err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.LittleEndian, int32(len(v))); err != nil {
				return nil, err
			}
			if _, err := buf.WriteString(v); err != nil {
				return nil, err
			}
		case float64:
			if err := binary.Write(buf, binary.LittleEndian, byte(3)); err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
				return nil, err
			}
		case float32:
			if err := binary.Write(buf, binary.LittleEndian, byte(4)); err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
				return nil, err
			}
		case int64:
			if err := binary.Write(buf, binary.LittleEndian, byte(5)); err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
				return nil, err
			}
		case int32:
			if err := binary.Write(buf, binary.LittleEndian, byte(6)); err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
				return nil, err
			}
		case int16:
			if err := binary.Write(buf, binary.LittleEndian, byte(7)); err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
				return nil, err
			}
		case int8:
			if err := binary.Write(buf, binary.LittleEndian, byte(8)); err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
				return nil, err
			}
		case uint:
			if err := binary.Write(buf, binary.LittleEndian, byte(9)); err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.LittleEndian, uint64(v)); err != nil {
				return nil, err
			}
		case uint64:
			if err := binary.Write(buf, binary.LittleEndian, byte(10)); err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
				return nil, err
			}
		case uint32:
			if err := binary.Write(buf, binary.LittleEndian, byte(11)); err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
				return nil, err
			}
		case uint16:
			if err := binary.Write(buf, binary.LittleEndian, byte(12)); err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
				return nil, err
			}
		case uint8:
			if err := binary.Write(buf, binary.LittleEndian, byte(13)); err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
				return nil, err
			}
		case []byte:
			if err := binary.Write(buf, binary.LittleEndian, byte(14)); err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.LittleEndian, int32(len(v))); err != nil {
				return nil, err
			}
			if _, err := buf.Write(v); err != nil {
				return nil, err
			}
		case time.Time:
			if err := binary.Write(buf, binary.LittleEndian, byte(15)); err != nil {
				return nil, err
			}
			timestamp := v.UnixNano()
			if err := binary.Write(buf, binary.LittleEndian, timestamp); err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("unsupported type")
		}
	}

	return buf.Bytes(), nil
}

func (t *Serializer) Decode(data []byte) (map[string]interface{}, error) {
	buf := bytes.NewReader(data)
	result := make(map[string]interface{})

	for buf.Len() > 0 {
		var keyLen int8
		if err := binary.Read(buf, binary.LittleEndian, &keyLen); err != nil {
			return nil, err
		}

		key := make([]byte, keyLen)
		if _, err := buf.Read(key); err != nil {
			return nil, err
		}

		var valueType byte
		if err := binary.Read(buf, binary.LittleEndian, &valueType); err != nil {
			return nil, err
		}

		switch valueType {
		case 1:
			var b byte
			if err := binary.Read(buf, binary.LittleEndian, &b); err != nil {
				return nil, err
			}
			result[string(key)] = b == 1
		case 2:
			var valueLen int32
			if err := binary.Read(buf, binary.LittleEndian, &valueLen); err != nil {
				return nil, err
			}
			value := make([]byte, valueLen)
			if _, err := buf.Read(value); err != nil {
				return nil, err
			}
			result[string(key)] = string(value)
		case 3:
			var value float64
			if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
				return nil, err
			}
			result[string(key)] = value
		case 4:
			var value float32
			if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
				return nil, err
			}
			result[string(key)] = value
		case 5:
			var value int64
			if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
				return nil, err
			}
			result[string(key)] = value
		case 6:
			var value int32
			if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
				return nil, err
			}
			result[string(key)] = value
		case 7:
			var value int16
			if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
				return nil, err
			}
			result[string(key)] = value
		case 8:
			var value int8
			if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
				return nil, err
			}
			result[string(key)] = value
		case 9:
			var value uint64
			if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
				return nil, err
			}
			result[string(key)] = uint(value)
		case 10:
			var value uint64
			if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
				return nil, err
			}
			result[string(key)] = value
		case 11:
			var value uint32
			if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
				return nil, err
			}
			result[string(key)] = value
		case 12:
			var value uint16
			if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
				return nil, err
			}
			result[string(key)] = value
		case 13:
			var value uint8
			if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
				return nil, err
			}
			result[string(key)] = value
		case 14:
			var valueLen int32
			if err := binary.Read(buf, binary.LittleEndian, &valueLen); err != nil {
				return nil, err
			}
			value := make([]byte, valueLen)
			if _, err := buf.Read(value); err != nil {
				return nil, err
			}
			result[string(key)] = value
		case 15:
			var timestamp int64
			if err := binary.Read(buf, binary.LittleEndian, &timestamp); err != nil {
				return nil, err
			}
			result[string(key)] = time.Unix(0, timestamp)
		default:
			return nil, errors.New("unsupported type")
		}
	}

	return result, nil
}
