// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoBuilder

import (
	"database/sql"
	"fmt"
	"github.com/donnie4w/gdao/base"
	"log"
	"reflect"
	"strings"
	"time"
)

type TableBean struct {
	TableName string
	Fieldlist []*FieldBean
	Fieldmap  map[string]*FieldBean
}

func (t *TableBean) ContainTime() bool {
	for _, field := range t.Fieldlist {
		if field.FieldType == reflect.TypeOf(time.Time{}) {
			return true
		}
	}
	return false
}

func newTableBean() *TableBean {
	return &TableBean{Fieldlist: make([]*FieldBean, 0), Fieldmap: make(map[string]*FieldBean)}
}

type FieldBean struct {
	FieldName     string
	FieldIndex    int
	FieldType     reflect.Type
	FieldTypeName string
}

func (f *FieldBean) String() string {
	return fmt.Sprint(f.FieldIndex+1, ".ColumnName:", f.FieldName, ",Type:", f.FieldTypeName)
}

func executeForTableInfo(sql string, db *sql.DB) (tb *TableBean, err error) {
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tb = newTableBean()
	types, _ := rows.ColumnTypes()
	for i, columntype := range types {
		fb := &FieldBean{}
		fb.FieldTypeName = columntype.DatabaseTypeName()
		fb.FieldType = columntype.ScanType()
		fb.FieldName = columntype.Name()
		fb.FieldIndex = i
		tb.Fieldmap[columntype.Name()] = fb
		tb.Fieldlist = append(tb.Fieldlist, fb)
	}
	return
}

func GetTableBean(tablename string, db *sql.DB) (tb *TableBean, err error) {
	sql := fmt.Sprint("select * from ", tablename, " where 0=1")
	tb, err = executeForTableInfo(sql, db)
	if tb != nil {
		tb.TableName = tablename
	}
	return
}

func buildstruct(dbtype, dbname, tableName, tableAlias string, packageName string, tableBean *TableBean) string {
	datetime := time.Now().Format(time.DateTime)
	ua := ToUpperFirstLetter
	structName := ua(tableAlias)

	timePackage := func() string {
		b := false
		for _, bean := range tableBean.Fieldlist {
			if bean.FieldType == reflect.TypeOf(sql.NullTime{}) {
				b = true
				break
			}
		}
		if b {
			return `"time"`
		} else {
			return ""
		}
	}()

	r := `// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao
//
// datetime :` + datetime + `
// gdao version ` + base.VERSION + `
// dbtype:` + dbtype + ` ,database:` + dbname + ` ,tablename:` + tableName + `

package ` + packageName + `

import (
	"fmt"
	"github.com/donnie4w/gdao/base"
	"github.com/donnie4w/gdao"
	` + timePackage + `
)`

	for _, bean := range tableBean.Fieldlist {
		log.Println(bean)
		rtype := goPtrType(bean.FieldType)
		structField := strings.ToLower(tableAlias) + "_" + ua(bean.FieldName)
		s := `
type ` + structField + `[T any] struct {
	base.Field[T]
	fieldName  string
	fieldValue ` + rtype + `
}

func (t *` + structField + `[T]) Name() string {
	return t.fieldName
}

func (t *` + structField + `[T]) Value() any {
	return t.fieldValue
}
`
		r = r + s
	}

	r = r + `
type ` + structName + ` struct {
	gdao.Table[` + structName + `]
`
	for _, bean := range tableBean.Fieldlist {
		s := `
	` + ua(bean.FieldName) + `		*` + strings.ToLower(tableAlias) + "_" + ua(bean.FieldName) + `[` + structName + `]`
		r = r + s
	}
	r = r + `
}
`
	mustptr := func(t reflect.Type, s string) string {
		if mustPtr(t) {
			return s
		} else {
			return ""
		}
	}

	for _, bean := range tableBean.Fieldlist {
		fieldName := ua(bean.FieldName)
		rtype := goType(bean.FieldType)
		s := `
func (u *` + structName + `) Get` + fieldName + `() (_r ` + rtype + `){
`
		s1 := `	if u.` + fieldName + `.fieldValue != nil {
		_r = ` + mustptr(bean.FieldType, "*") + `u.` + fieldName + `.fieldValue
	}`
		if mustPtr(bean.FieldType) {
			s = s + s1
		} else {
			s = s + `	_r = ` + mustptr(bean.FieldType, "*") + `u.` + fieldName + `.fieldValue`
		}
		s = s + `
	return
}

func (u *` + structName + `) Set` + fieldName + `(arg ` + rtype + `) *` + structName + `{
	u.Put0(u.` + fieldName + `.fieldName, arg)
	u.` + fieldName + `.fieldValue = ` + mustptr(bean.FieldType, "&") + `arg
	return u
}
`
		r = r + s
	}

	r = r + `

func (u *` + structName + `) Scan(fieldname string, value any) {
	switch fieldname {`

	for _, bean := range tableBean.Fieldlist {
		fieldName := ua(bean.FieldName)
		rtype := goType(bean.FieldType)
		var s string
		if bean.FieldType == reflect.TypeOf(sql.NullTime{}) {

			s = `
	case "` + bean.FieldName + `":
		if t, err := base.AsTime(value); err == nil {
			u.Set` + fieldName + `(t)
		}`
		} else if bean.FieldType.Kind() == reflect.Slice && bean.FieldType.Elem().Kind() == reflect.Uint8 {

			s = `
	case "` + bean.FieldName + `":
		u.Set` + fieldName + `(base.AsBytes(value))`
		} else {

			s = `
	case "` + bean.FieldName + `":
		u.Set` + fieldName + `(base.As` + ua(rtype) + `(value))`
		}
		r = r + s
	}
	r = r + `
	}
}
`
	columns := ""
	fields := ""
	fieldsString := ""
	columnsStr := ""
	for i, bean := range tableBean.Fieldlist {
		columns = columns + "t." + ua(bean.FieldName)
		fieldsString = fieldsString + "\"" + ua(bean.FieldName) + ":\"" + ",t.Get" + ua(bean.FieldName) + "()"
		fields = fields + ua(bean.FieldName) + ":" + encodeFieldname(bean.FieldName)
		columnsStr = columnsStr + encodeFieldname(bean.FieldName)
		if i < len(tableBean.Fieldlist)-1 {
			columns = columns + ","
			fieldsString = fieldsString + `, ",",`
			fields = fields + ","
			columnsStr = columnsStr + ","
		}
	}

	//	selectfunc := `
	//
	//func (t *` + structName + `) Selects(columns ...base.Column[` + structName + `]) (_r []*` + structName + `, err error) {
	//	if columns == nil {
	//		columns = []base.Column[` + structName + `]{` + columns + `}
	//	}
	//	databeans, err := t.ExecuteQueryBeans(columns...)
	//	if err != nil || len(databeans) == 0 {
	//		return nil, err
	//	}
	//	_r = make([]*` + structName + `, 0)
	//	for _, beans := range databeans {
	//		__` + structName + ` := New` + structName + `()
	//		for name, bean := range beans.Map() {
	//			__` + structName + `.Scan(name, bean.Value())
	//		}
	//		_r = append(_r, __` + structName + `)
	//	}
	//	return
	//}
	//
	//func (t *` + structName + `) Select(columns ...base.Column[` + structName + `]) (_r *` + structName + `, err error) {
	//	if columns == nil {
	//		columns = []base.Column[` + structName + `]{` + columns + `}
	//	}
	//	databean, err := t.ExecuteQueryBean(columns...)
	//	if err != nil || databean == nil {
	//		return nil, err
	//	}
	//	_r = New` + structName + `()
	//	for name, bean := range databean.Map() {
	//		_r.Scan(name, bean.Value())
	//	}
	//	return
	//}
	//`
	//	r = r + selectfunc

	r = r + `
func (t *` + structName + `) ToGdao() {
	_t := New` + structName + `()
	*t = *_t
}
`

	copy := `
func (t *` + structName + `) Copy(h *` + structName + `) *` + structName + `{`
	r = r + copy
	for _, bean := range tableBean.Fieldlist {
		r = r + `
	t.Set` + ua(bean.FieldName) + `(h.Get` + ua(bean.FieldName) + `())`
	}
	r = r + `
	return t
}
`

	stringBody := `
func (t *` + structName + `) String() string {
	return fmt.Sprint(` + fieldsString + `)
}
`
	r = r + stringBody

	newfunc := `
func New` + structName + `(tablename ...string) (_r *` + structName + `) {
`
	for _, bean := range tableBean.Fieldlist {
		structField := strings.ToLower(tableAlias) + "_" + ua(bean.FieldName)
		varfield := encodeFieldname(bean.FieldName)
		s := `
	` + varfield + ` := &` + structField + `[` + structName + `]{fieldName: "` + bean.FieldName + `"}
	` + varfield + `.Field.FieldName = "` + bean.FieldName + `"
`
		newfunc = newfunc + s
	}

	newfunc = newfunc + `
	_r = &` + structName + `{` + fields + `}
	s := "` + tableName + `"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.Init(s, []base.Column[` + structName + `]{` + columnsStr + `})
	return
}
`
	r = r + newfunc

	serialstr := `
func (t *` + structName + `) Encode() ([]byte, error) {
	m := make(map[string]any, 0)`
	r += serialstr
	for _, bean := range tableBean.Fieldlist {
		r += `
	m["` + bean.FieldName + `"] = t.Get` + ua(bean.FieldName) + `()`
	}

	serialstr = `
	return t.Table.Encode(m)
}

func (t *` + structName + `) Decode(bs []byte) (err error) {
	var m map[string]any
	if m, err = t.Table.Decode(bs); err == nil {
		if !t.IsInit() {
			t.ToGdao()
		}
		for name, bean := range m {
			t.Scan(name, bean)
		}
	}
	return
}

`
	r = r + serialstr
	return r
}
