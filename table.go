// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdao

import (
	. "github.com/donnie4w/gdao/base"
	"github.com/donnie4w/gdao/gdaoCache"
	"strings"
)

type Table[T any] struct {
	TableBase[T]
	commentline string
	tableName   string
	querySql    string
	whereSql    string
	args        []any
	groupSql    string
	havingSql   string
	orderSql    string
	limitSql    string
	sql         string
	modifymap   map[string]any
	modifySql   string
	dbhandler   DBhandle
	transaction Transaction
	batchArgs   [][]any
	mustMaster  bool
	isCache     int8
	classname   string
}

func (t *Table[T]) Init(s string) {
	t.tableName = s
	t.args = make([]any, 0)
	t.modifymap = make(map[string]any, 0)
}

func (t *Table[T]) IsInit() bool {
	return t.tableName != "" && t.args != nil && t.modifymap != nil
}

func (t *Table[T]) Put0(k string, v any) {
	t.modifymap[k] = v
}

func (t *Table[T]) UseCache(use bool) {
	if use {
		t.isCache = 1
	} else {
		t.isCache = 2
	}
}

func (t Table[T]) ClassName() {
}

func (t *Table[T]) Where(wheres ...*Where[T]) *Table[T] {
	whereSqls := make([]string, len(wheres))
	for i, w := range wheres {
		whereSqls[i] = w.WhereSql
		t.whereSql = " " + w.WhereSql + " "
		if w.Value != nil {
			t.args = append(t.args, w.Value)
		}
		if w.Values != nil {
			for _, v := range w.Values {
				t.args = append(t.args, v)
			}
		}
	}
	s := strings.Join(whereSqls, " and ")
	t.whereSql = " where " + s
	return t
}

func (t *Table[T]) UseTransaction(transaction Transaction) {
	t.transaction = transaction
}

func (t *Table[T]) MustMaster(must bool) {
	t.mustMaster = must
}

func (t *Table[T]) UseDBHandler(db DBhandle) *Table[T] {
	t.dbhandler = db
	return t
}

func (t *Table[T]) ExecuteQueryBeans(columns ...Column[T]) ([]*DataBean, error) {
	t.completeSql4Columns(columns...)
	t.completeSql4Query()

	if Logger.IsVaild {
		Logger.Debug("[SELETE LIST]["+t.sql+"]", t.args)
	}
	if t.classname == "" {
		t.classname = Classname[T]()
	}
	domain := gdaoCache.GetDomain(t.classname)
	iscache := (t.isCache == 1 || domain != "") && t.isCache != 2
	var condition *gdaoCache.Condition
	if iscache {
		condition = gdaoCache.NewCondition("list", t.sql, t.args...)
		if result := gdaoCache.GetCache(domain, t.classname, condition); result != nil {
			if Logger.IsVaild {
				Logger.Debug("[GET CACHE]["+t.sql+"]", t.args)
			}
			return result.([]*DataBean), nil
		}
	}

	if g := t.getDB(true); g != nil {
		if v, err := g.ExecuteQueryBeans(t.sql, t.args...); err == nil {
			if iscache {
				gdaoCache.SetCache(domain, t.classname, condition, v)
				if Logger.IsVaild {
					Logger.Debug("[SET CACHE]["+t.sql+"]", t.args)
				}
			}
			return v, nil
		} else {
			return nil, err
		}
	} else {
		return nil, errInit
	}
}

func (t *Table[T]) ExecuteQueryBean(columns ...Column[T]) (*DataBean, error) {
	t.completeSql4Columns(columns...)
	t.completeSql4Query()

	if Logger.IsVaild {
		Logger.Debug("[SELETE ONE]["+t.sql+"]", t.args)
	}
	if t.classname == "" {
		t.classname = Classname[T]()
	}
	domain := gdaoCache.GetDomain(t.classname)
	iscache := (t.isCache == 1 || domain != "") && t.isCache != 2
	var condition *gdaoCache.Condition
	if iscache {
		condition = gdaoCache.NewCondition("one", t.sql, t.args...)
		if result := gdaoCache.GetCache(domain, t.classname, condition); result != nil {
			if Logger.IsVaild {
				Logger.Debug("[GET CACHE]["+t.sql+"]", t.args)
			}
			return result.(*DataBean), nil
		}
	}

	if g := t.getDB(true); g != nil {
		if v, err := g.ExecuteQueryBean(t.sql, t.args...); err == nil {
			if iscache {
				gdaoCache.SetCache(domain, t.classname, condition, v)
				if Logger.IsVaild {
					Logger.Debug("[SET CACHE]["+t.sql+"]", t.args)
				}
			}
			return v, nil
		} else {
			return nil, err
		}
	} else {
		return nil, errInit
	}
}

func (t *Table[T]) completeSql4Columns(columns ...Column[T]) {
	querycolumns := make([]string, len(columns))
	for i, c := range columns {
		name := c.Name()
		querycolumns[i] = name
	}
	s := strings.Join(querycolumns, ",")
	t.querySql = t.commentline + " select " + s + " from " + t.tableName
}

func (t *Table[T]) completeSql4Query() {
	t.sql = t.querySql
	if t.sql != "" {
		if t.whereSql != "" {
			t.sql = t.sql + t.whereSql
		}
		if t.groupSql != "" {
			t.sql = t.sql + t.groupSql
		}
		if t.havingSql != "" {
			t.sql = t.sql + t.havingSql
		}
		if t.orderSql != "" {
			t.sql = t.sql + t.orderSql
		}
		if t.limitSql != "" {
			t.sql = t.sql + t.limitSql
		}
	}
}

func (t *Table[T]) completeSql4Update() {
	t.sql = t.modifySql
	if t.sql != "" {
		if t.whereSql != "" {
			t.sql = t.sql + t.whereSql
		}
		if t.groupSql != "" {
			t.sql = t.sql + t.groupSql
		}
		if t.havingSql != "" {
			t.sql = t.sql + t.havingSql
		}
	}
}

func (t *Table[T]) getDB(queryType bool) (r DBhandle) {
	if t.transaction != nil {
		return t.transaction
	}
	if t.dbhandler != nil {
		return t.dbhandler
	}
	return getDBhandle(Classname[T](), t.tableName, queryType && !t.mustMaster)
}

func (t *Table[T]) GroupBy(columns ...Column[T]) *Table[T] {
	ss := make([]string, 0, len(columns))
	for _, v := range columns {
		ss = append(ss, v.Name())
	}
	t.groupSql = " group by " + strings.Join(ss, ",")
	return t
}

func (t *Table[T]) Having(havings ...*Having[T]) *Table[T] {
	ss := make([]string, 0, len(havings))
	for _, w := range havings {
		ss = append(ss, w.HavingSql)
		if w.Value != nil {
			t.args = append(t.args, w.Value)
		}
		if w.Values != nil {
			for _, v := range w.Values {
				t.args = append(t.args, v)
			}
		}
	}
	t.havingSql = " having " + strings.Join(ss, ",")
	return t
}

func (t *Table[T]) OrderBy(sorts ...*Sort[T]) *Table[T] {
	ss := make([]string, 0, len(sorts))
	for _, v := range sorts {
		ss = append(ss, v.OrderByArg)
	}
	t.orderSql = " order by " + strings.Join(ss, ",")
	return t
}

func (t *Table[T]) Limit(limit int64) {
	if limit > 0 {
		t.limitAdapt(limit)
	}
}

func (t *Table[T]) Limit2(offset, limit int64) {
	if limit != 0 {
		t.limit2Adapt(offset, limit)
	}
}

func (t *Table[T]) limitAdapt(limit int64) {
	switch t.getDB(true).GetDBType() {
	case SQLSERVER:
		t.limitSql = " OFFSET 0 ROWS FETCH NEXT ? ROWS ONLY "
	case ORACLE:
		t.limitSql = " FETCH FIRST ? ROWS ONLY "
	case NETEZZA, GREENPLUM, POSTGRESQL:
		t.limitSql = " LIMIT ? OFFSET 0 "
	case DB2:
		t.limitSql = " FETCH FIRST ? ROWS ONLY "
	case TERADATA, FIREBIRD, SYBASE:
		t.limitSql = ""
	case DERBY:
		t.limitSql = " FETCH FIRST ? ROWS ONLY "
	case INGRES, VERTICA, MYSQL, MARIADB, SQLITE:
		t.limitSql = " LIMIT ? "
	default:
		t.limitSql = ""
	}
	if t.limitSql != "" {
		t.args = append(t.args, limit)
	}
}

func (t *Table[T]) limit2Adapt(offset, limit int64) {
	switch t.getDB(true).GetDBType() {
	case POSTGRESQL, GREENPLUM:
		t.limitSql = " OFFSET ? LIMIT ? "
		t.args = append(t.args, offset, limit)
	case ORACLE, SQLSERVER:
		t.limitSql = " OFFSET ? ROWS FETCH NEXT ? ROWS ONLY "
		t.args = append(t.args, offset, limit)
	case SQLITE, NETEZZA, INGRES, VERTICA:
		t.limitSql = " LIMIT ? OFFSET ? "
		t.args = append(t.args, limit, offset)
	case DB2, DERBY:
		t.limitSql = " FETCH FIRST ? ROWS ONLY OFFSET ? ROWS "
		t.args = append(t.args, limit, offset)
	case SYBASE, TERADATA, FIREBIRD:
		t.limitSql = ""
	case MYSQL, MARIADB:
		t.limitSql = " LIMIT ?,? "
		t.args = append(t.args, offset, limit)
	}

}

func (t *Table[T]) Update() (int64, error) {
	modifystr := make([]string, 0)
	args := make([]any, 0)
	for k, v := range t.modifymap {
		modifystr = append(modifystr, k+"=?")
		args = append(args, v)
	}
	t.modifySql = "update " + t.tableName + " set " + strings.Join(modifystr, ",")
	for _, v := range t.args {
		args = append(args, v)
	}
	t.args = args
	t.completeSql4Update()

	if Logger.IsVaild {
		Logger.Debug("[UPDATE]["+t.sql+"]", t.args)
	}

	if g := t.getDB(false); g != nil {
		return g.ExecuteUpdate(t.sql, t.args...)
	} else {
		return 0, errInit
	}
}

func (t *Table[T]) Insert() (int64, error) {
	insertField := make([]string, 0)
	insert_ := make([]string, 0)
	args := make([]any, 0)
	for k, v := range t.modifymap {
		insertField = append(insertField, k)
		insert_ = append(insert_, "?")
		args = append(args, v)
	}
	t.sql = "insert  into " + t.tableName + "(" + strings.Join(insertField, ",") + " )values(" + strings.Join(insert_, ",") + ")"
	for _, v := range t.args {
		args = append(args, v)
	}
	t.args = args

	if Logger.IsVaild {
		Logger.Debug("[INSERT]["+t.sql+"]", t.args)
	}

	if g := t.getDB(false); g != nil {
		return g.ExecuteUpdate(t.sql, t.args...)
	} else {
		return 0, errInit
	}
}

func (t *Table[T]) AddBatch() {
	insertField := make([]string, 0)
	insert_ := make([]string, 0)
	if len(t.batchArgs) == 0 {
		t.batchArgs = make([][]any, 0)
	}
	args := make([]any, 0)
	for k, v := range t.modifymap {
		insertField = append(insertField, k)
		insert_ = append(insert_, "?")
		args = append(args, v)
	}
	t.sql = " insert  into " + t.tableName + "(" + strings.Join(insertField, ",") + " )values(" + strings.Join(insert_, ",") + ")"
	for _, v := range t.args {
		args = append(args, v)
	}
	t.batchArgs = append(t.batchArgs, args)
}

func (t *Table[T]) ExecBatch() ([]int64, error) {

	if Logger.IsVaild {
		Logger.Debug("[BATCH]["+t.sql+"]", t.batchArgs)
	}

	if g := t.getDB(false); g != nil {
		return g.ExecuteBatch(t.sql, t.batchArgs)
	} else {
		return nil, errInit
	}
}

func (t *Table[T]) Delete() (int64, error) {
	t.modifySql = " delete from " + t.tableName
	t.completeSql4Update()

	if Logger.IsVaild {
		Logger.Debug("[DELETE]["+t.sql+"]", t.args)
	}

	if g := t.getDB(false); g != nil {
		return g.ExecuteUpdate(t.sql, t.args...)
	} else {
		return 0, errInit
	}
}

var serialize Serialize[map[string]any] = &Serializer{}

func (t *Table[T]) Encode(m map[string]any) ([]byte, error) {
	return serialize.Encode(m)
}

func (t *Table[T]) Decode(data []byte) (map[string]any, error) {
	return serialize.Decode(data)
}

func (t *Table[T]) CommentLine(commentline string) {
	t.commentline = `/*` + commentline + `*/`
}
