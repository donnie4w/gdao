// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package base

import (
	"fmt"
	"strings"
)

type Field[T any] struct {
	FieldName string
}

func (f *Field[T]) EQ(arg any) *Where[T] {
	return &Where[T]{f.FieldName + "=?", arg, nil}
}

func (f *Field[T]) NEQ(arg any) *Where[T] {
	return &Where[T]{f.FieldName + "<>?", arg, nil}
}

func (f *Field[T]) LT(arg any) *Where[T] {
	return &Where[T]{f.FieldName + "<?", arg, nil}
}

func (f *Field[T]) LE(arg any) *Where[T] {
	return &Where[T]{f.FieldName + "<=?", arg, nil}
}

func (f *Field[T]) GT(arg any) *Where[T] {
	return &Where[T]{f.FieldName + ">?", arg, nil}
}

func (f *Field[T]) GE(arg any) *Where[T] {
	return &Where[T]{f.FieldName + ">=?", arg, nil}
}

func (f *Field[T]) LIKE(arg any) *Where[T] {
	return &Where[T]{f.FieldName + " like ?", fmt.Sprint("%", arg, "%"), nil}
}

func (f *Field[T]) RLIKE(arg any) *Where[T] {
	return &Where[T]{f.FieldName + " like ?", fmt.Sprint("%", arg), nil}
}

func (f *Field[T]) LLIKE(arg any) *Where[T] {
	return &Where[T]{f.FieldName + " like ?", fmt.Sprint(arg, "%"), nil}
}

func (f *Field[T]) Between(from, to any) *Where[T] {
	return &Where[T]{f.FieldName + " between ? and ?", nil, []any{from, to}}
}

func (f *Field[T]) IN(args ...any) *Where[T] {
	s := make([]string, len(args))
	for i := 0; i < len(args); i++ {
		s[i] = "?"
	}
	return &Where[T]{f.FieldName + " in (" + strings.Join(s, ",") + ")", nil, args}
}

func (f *Field[T]) NOTIN(args ...any) *Where[T] {
	s := make([]string, len(args))
	for i := 0; i < len(args); i++ {
		s[i] = "?"
	}
	return &Where[T]{f.FieldName + " not in (" + strings.Join(s, ",") + ")", nil, args}
}

func (f *Field[T]) Asc() *Sort[T] {
	return &Sort[T]{f.FieldName + " asc "}
}

func (f *Field[T]) Desc() *Sort[T] {
	return &Sort[T]{f.FieldName + " desc "}
}

func (f *Field[T]) Count() *Func[T] {
	return &Func[T]{FieldName: " count(" + f.FieldName + ") "}
}

func (f *Field[T]) Distinct() *Func[T] {
	return &Func[T]{FieldName: " distinct " + f.FieldName + " "}
}

func (f *Field[T]) Sum() *Func[T] {
	return &Func[T]{FieldName: " sum(" + f.FieldName + ") "}
}

func (f *Field[T]) Avg() *Func[T] {
	return &Func[T]{FieldName: " avg(" + f.FieldName + ") "}
}

func (f *Field[T]) Max() *Func[T] {
	return &Func[T]{FieldName: " max(" + f.FieldName + ") "}
}

func (f *Field[T]) Min() *Func[T] {
	return &Func[T]{FieldName: " min(" + f.FieldName + ") "}
}

func (f *Field[T]) Operation(qurey4SetOperation string) *Func[T] {
	return &Func[T]{FieldName: " " + qurey4SetOperation + " "}
}
