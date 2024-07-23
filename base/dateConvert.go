// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package base

import (
	"github.com/donnie4w/gofer/buffer"
	"time"
	"unicode"
)

func convertToDate(isoDateString string) (time.Time, error) {
	buf := buffer.NewBufferByPool()
	defer buf.Free()
	isdecimals := true
	hasZ := false
	for i, ch := range isoDateString {
		switch i + 1 {
		case 4:
			buf.WriteString("2006")
		case 5:
			buf.WriteString(string(ch))
		case 7:
			buf.WriteString("01")
		case 8:
			buf.WriteString(string(ch))
		case 10:
			buf.WriteString("02")
		case 11:
			if ch == 'T' {
				buf.WriteString("T")
			} else {
				buf.WriteString(string(ch))
			}
		case 13:
			buf.WriteString("15")
		case 14:
			buf.WriteString(string(ch))
		case 16:
			buf.WriteString("04")
		case 17:
			buf.WriteString(string(ch))
		case 19:
			buf.WriteString("05")
		case 20:
			buf.WriteString(string(ch))
		default:
			if i > 19 {
				if unicode.IsDigit(ch) && isdecimals {
					buf.WriteString("9")
				} else if ch == '+' {
					isdecimals = false
					if !hasZ {
						buf.WriteByte('Z')
						hasZ = true
					}
					if i+4 <= len(isoDateString) && isoDateString[i+3:i+4] == ":" {
						buf.WriteString("07:00")
					} else {
						buf.WriteString("0700")
					}
				} else if ch == '-' {
					isdecimals = false
					if !hasZ {
						buf.WriteByte('Z')
						hasZ = true
					}
					if i+4 <= len(isoDateString) && isoDateString[i+3:i+4] == ":" {
						buf.WriteString("-07:00")
					} else {
						buf.WriteString("-0700")
					}
				} else if ch == 'Z' {
					isdecimals = false
					if !hasZ {
						buf.WriteByte('Z')
						hasZ = true
					}
				}
			}
		}
	}
	return time.Parse(buf.String(), isoDateString)
}
