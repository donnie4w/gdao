// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package base

import "strings"

type Column[T any] interface {
	Name() string
}

type Col string

func (c Col) Name() string {
	return string(c)
}

type Sort[T any] struct {
	OrderByArg string
}

type Where[T any] struct {
	WhereSql string
	Value    any
	Values   []any
}

type Having[T any] struct {
	HavingSql string
	Value     any
	Values    []any
}

func (w *Where[T]) And(wheres ...*Where[T]) *Where[T] {
	whereSqls := make([]string, 0, len(wheres))
	for _, v := range wheres {
		whereSqls = append(whereSqls, v.WhereSql)
		if v.Value != nil {
			w.Values = append(w.Values, v.Value)
		}
		if v.Values != nil {
			for _, vv := range v.Values {
				w.Values = append(w.Values, vv)
			}
		}
	}
	w.WhereSql = w.WhereSql + " and (" + strings.Join(whereSqls, " or ") + ")"
	return w
}

func (w *Where[T]) Or(wheres ...*Where[T]) *Where[T] {
	whereSqls := make([]string, 0, len(wheres))
	for _, v := range wheres {
		whereSqls = append(whereSqls, v.WhereSql)
		if v.Value != nil {
			w.Values = append(w.Values, v.Value)
		}
		if v.Values != nil {
			for _, vv := range v.Values {
				w.Values = append(w.Values, vv)
			}
		}
	}
	w.WhereSql = w.WhereSql + " or (" + strings.Join(whereSqls, " and ") + ")"
	return w
}

type Func[T any] struct {
	FieldName  string
	FieldValue any
}

func (s *Func[T]) AS(alias Column[T]) *Func[T] {
	s.FieldName = s.FieldName + " as " + alias.Name()
	return s
}

func (s *Func[T]) EQ(arg any) *Having[T] {
	return &Having[T]{s.FieldName + "=?", arg, nil}
}

func (s *Func[T]) NEQ(arg any) *Having[T] {
	return &Having[T]{s.FieldName + "<>?", arg, nil}
}

func (s *Func[T]) LT(arg any) *Having[T] {
	return &Having[T]{s.FieldName + "<?", arg, nil}
}

func (s *Func[T]) LE(arg any) *Having[T] {
	return &Having[T]{s.FieldName + "<=?", arg, nil}
}

func (s *Func[T]) GT(arg any) *Having[T] {
	return &Having[T]{s.FieldName + ">?", arg, nil}
}

func (s *Func[T]) GE(arg any) *Having[T] {
	return &Having[T]{s.FieldName + ">=?", arg, nil}
}

func (s *Func[T]) Between(from, to any) *Having[T] {
	return &Having[T]{s.FieldName + " between ? and ?", nil, []any{from, to}}
}

func (s *Func[T]) Name() string {
	return s.FieldName
}

func (s *Func[T]) Value() any {
	return s.FieldValue
}
