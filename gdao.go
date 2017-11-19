package gdao

/**
  donnie4w@gmail.com
*/

import (
	"database/sql"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type DbType int32

const (
	_VER = "1.0.2"
)
const (
	_ DbType = iota
	MYSQL
	PostgreSQL
)

var adapter DbType = MYSQL

var _MYSQL_TABLE_SCHMAINFO_SQL string = "SHOW COLUMNS FROM "
var _PostgreSQL_TABLE_SCHMAINFO_SQL string = "SELECT column_name,data_type FROM information_schema.COLUMNS WHERE TABLE_NAME="

//var _SQLITE_TABLE_SCHMAINFO_SQL string = ""
//var _ORACLE_TABLE_SCHMAINFO_SQL string = ""
//var _SQLSERVER_TABLE_SCHMAINFO_SQL string = ""
//var _DB2_TABLE_SCHMAINFO_SQL string = ""

var adapterMap = map[DbType]string{MYSQL: _MYSQL_TABLE_SCHMAINFO_SQL, PostgreSQL: _PostgreSQL_TABLE_SCHMAINFO_SQL}

var db *sql.DB

var dbMap map[string]*sql.DB = make(map[string]*sql.DB)

func SetAdapterType(t DbType) {
	adapter = t
}

func SetDB(_db *sql.DB) {
	db = _db
}

func SetDBSrouceByTableName(tableName string, db *sql.DB) {
	dbMap[tableName] = db
}

type Column interface {
	Name() string
	Value() interface{}
}

type FieldBeen struct {
	fieldName  string
	fieldIndex int
	fieldValue interface{}
}

func (f *FieldBeen) Value() interface{} {
	return GetValue(&f.fieldValue)
}

func (f *FieldBeen) ValueString() string {
	return GetValue(&f.fieldValue).(string)
}

func (f *FieldBeen) ValueInt64() int64 {
	if f.fieldValue == nil {
		return 0
	}
	switch (f.fieldValue).(type) {
	case int64:
		return int64((f.fieldValue).(int64))
	case int32:
		return int64((f.fieldValue).(int32))
	case int16:
		return int64((f.fieldValue).(int16))
	case int8:
		return int64((f.fieldValue).(int8))
	case uint64:
		return int64((f.fieldValue).(uint64))
	case uint32:
		return int64((f.fieldValue).(uint32))
	case uint16:
		return int64((f.fieldValue).(uint16))
	case uint8:
		return int64((f.fieldValue).(uint8))
	case int:
		return int64((f.fieldValue).(int))
	case float32:
		return int64((f.fieldValue).(float32))
	case float64:
		return int64((f.fieldValue).(float64))
	case []uint8:
		i, _ := strconv.ParseInt(string((f.fieldValue).([]uint8)), 0, 0)
		return int64(i)
	default:
		return int64((f.fieldValue).(int64))
	}
}

func (f *FieldBeen) ValueInt32() int32 {
	return int32(f.ValueInt64())
}

func (f *FieldBeen) ValueInt16() int16 {
	return int16(f.ValueInt64())
}

func (f *FieldBeen) ValueFloat64() float64 {
	if f.fieldValue == nil {
		return 0
	}
	switch (f.fieldValue).(type) {
	case int64:
		return float64((f.fieldValue).(int64))
	case int32:
		return float64((f.fieldValue).(int32))
	case int16:
		return float64((f.fieldValue).(int16))
	case int8:
		return float64((f.fieldValue).(int8))
	case uint64:
		return float64((f.fieldValue).(uint64))
	case uint32:
		return float64((f.fieldValue).(uint32))
	case uint16:
		return float64((f.fieldValue).(uint16))
	case uint8:
		return float64((f.fieldValue).(uint8))
	case int:
		return float64((f.fieldValue).(int))
	case float32:
		return float64((f.fieldValue).(float32))
	case float64:
		return float64((f.fieldValue).(float64))
	case []uint8:
		i, _ := strconv.ParseInt(string((f.fieldValue).([]uint8)), 0, 0)
		return float64(i)
	default:
		return float64((f.fieldValue).(int64))
	}
}

func (f *FieldBeen) ValueFloat32() float32 {
	return float32(f.ValueFloat64())
}

func (f *FieldBeen) Name() string {
	return f.fieldName
}

func (f *FieldBeen) Index() int {
	return f.fieldIndex
}

type GoBeen struct {
	FieldBeens    []*FieldBeen
	FieldMapName  map[string]*FieldBeen
	FieldMapIndex map[int]*FieldBeen
}

func (g *GoBeen) MapName(name string) *FieldBeen {
	v, ok := g.FieldMapName[name]
	if ok {
		return v
	} else {
		return nil
	}
}

func (g *GoBeen) MapIndex(index int) *FieldBeen {
	v, ok := g.FieldMapIndex[index]
	if ok {
		return v
	} else {
		return nil
	}
}

type Table struct {
	islog       bool
	commentline string
	TableName   string
	querySql    string
	whereSql    string
	args        []interface{}
	groupBySql  string
	havingSql   string
	orderBySql  string
	limitSql    string
	sql         string
	ModifyMap   map[string]interface{}
	modifySql   string
	DB          *sql.DB
	Tx          *TX
}

type TX struct {
	tx   *sql.Tx
	isBg bool
}

func (t *Table) SetTx(tx *TX) {
	t.Tx = tx
}

func (x *TX) Begin(dbsource ...*sql.DB) {
	if dbsource != nil && len(dbsource) == 1 {
		x.tx, _ = dbsource[0].Begin()
		x.isBg = true
	} else {
		x.tx, _ = db.Begin()
		x.isBg = true
	}
}

func (x *TX) Commit() {
	if x.tx != nil {
		x.tx.Commit()
		x.isBg = false
	}
}

func (x *TX) RollBack() {
	if x.tx != nil {
		x.tx.Rollback()
		x.isBg = false
	}
}

func GetTX() *TX {
	return &TX{isBg: false}
}

type Field struct {
	FieldName string
}

type Sort struct {
	OrderByArg string
}

type SetOperation struct {
	fieldName  string
	FieldValue interface{}
}

type Where struct {
	WhereSql string
	Value    interface{}
	Values   []interface{}
}

type Having struct {
	havingSql string
	Value     interface{}
	Values    []interface{}
}

func (w *Where) And(wheres ...*Where) *Where {
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

func (w *Where) Or(wheres ...*Where) *Where {
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

func (t *Table) Where(wheres ...*Where) {
	whereSqls := make([]string, len(wheres))
	t.args = make([]interface{}, 0, len(wheres))
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
}

func (t *Table) IsLog(islog bool) {
	t.islog = islog
}

func (t *Table) logger(v ...interface{}) {
	if t.islog {
		log.Println("[gdao log]", v)
	}
}

func (t *Table) Selects(columns ...Column) (rows *sql.Rows, err error) {
	t.completeSql4Columns(columns...)
	t.completeSql4Query()
	return t.executeQuery_()
}

func (t *Table) executeQuery_() (rows *sql.Rows, err error) {
	rows, err = t.getDB().Query(t.sql, t.args...)
	if err != nil {
		return nil, err
	}
	return
}

func (t *Table) Query(columns ...Column) ([][]interface{}, error) {
	t.completeSql4Columns(columns...)
	t.completeSql4Query()
	return t.executeQuery()
}

func (t *Table) QueryBeen(columns ...Column) ([]*GoBeen, error) {
	t.completeSql4Columns(columns...)
	t.completeSql4Query()
	return t.executeQueryBeen()
}

func (t *Table) completeSql4Columns(columns ...Column) {
	querycolumns := make([]string, len(columns))
	for i, c := range columns {
		name := c.Name()
		querycolumns[i] = name
	}
	s := strings.Join(querycolumns, ",")
	t.querySql = t.commentline + " select " + s + " from " + t.TableName
}

func (t *Table) completeSql4Query() {
	t.sql = t.querySql
	switch t.sql != "" {
	case t.whereSql != "":
		t.sql = t.sql + t.whereSql
		fallthrough
	case t.groupBySql != "":
		t.sql = t.sql + t.groupBySql
		fallthrough
	case t.havingSql != "":
		t.sql = t.sql + t.havingSql
		fallthrough
	case t.orderBySql != "":
		t.sql = t.sql + t.orderBySql
		fallthrough
	case t.limitSql != "":
		t.sql = t.sql + t.limitSql
	}
	t.logger(t.sql, t.args)
}

func (t *Table) completeSql4Update() {
	t.sql = t.modifySql
	switch t.sql != "" {
	case t.whereSql != "":
		t.sql = t.sql + t.whereSql
		fallthrough
	case t.groupBySql != "":
		t.sql = t.sql + t.groupBySql
		fallthrough
	case t.havingSql != "":
		t.sql = t.sql + t.havingSql
	}
	t.logger(t.sql, t.args)
}

func (t *Table) QuerySingle(columns ...Column) ([]interface{}, error) {
	t.completeSql4Columns(columns...)
	t.completeSql4Query()
	return t.executeQuerySingle()
}

func ToUpperFirstLetter(arg string) string {
	return strings.ToUpper(string(arg[0])) + Substr(arg, 1, len(arg)-1)
}

func Substr(str string, start, length int) string {
	if start < 0 || length < 0 {
		return str
	}
	rs := []rune(str)
	end := start + length
	return string(rs[start:end])
}

func exception(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func (t *Table) executeQuery() ([][]interface{}, error) {
	rows, err := t.getDB().Query(t.sql, t.args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	cols, _ := rows.Columns()
	ts := make([][]interface{}, 0, 1)
	for rows.Next() {
		buff := make([]interface{}, len(cols))
		data := make([]interface{}, len(cols))
		for i, _ := range buff {
			buff[i] = &data[i]
			//buff[i] = data[i]
		}
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		ts = append(ts, data)
	}
	if len(ts) == 0 {
		return nil, nil
	}
	return ts, nil
}

func (t *Table) executeQueryBeen() ([]*GoBeen, error) {
	return executeQuery_(t.getDB(), t.sql, t.args...)
}

func Query(db *sql.DB, sql string, args ...interface{}) ([]*GoBeen, error) {
	return executeQuery_(db, sql, args...)
}

func ExecuteQuery(sql string, args ...interface{}) ([]*GoBeen, error) {
	return executeQuery_(db, sql, args...)
}

func executeQuery_(dbsource *sql.DB, sql string, args ...interface{}) ([]*GoBeen, error) {
	rows, err := dbsource.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//exception(err)
	cols, _ := rows.Columns()
	gb := make([]*GoBeen, 0)
	for rows.Next() {
		gobeen := new(GoBeen)
		gobeen.FieldMapName = make(map[string]*FieldBeen, 0)
		gobeen.FieldMapIndex = make(map[int]*FieldBeen, 0)
		buff := make([]interface{}, 0, len(cols))
		i := 1
		for _, c := range cols {
			fb := new(FieldBeen)
			fb.fieldName = c
			fb.fieldIndex = i
			gobeen.FieldBeens = append(gobeen.FieldBeens, fb)
			buff = append(buff, &fb.fieldValue)
			gobeen.FieldMapName[c] = fb
			gobeen.FieldMapIndex[i] = fb
			i++
		}
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		gb = append(gb, gobeen)
	}
	if len(gb) == 0 {
		return nil, nil
	}
	return gb, nil
}

func ExecuteUpdate(sql string, args ...interface{}) (int64, error) {
	return executeUpdate_(nil, db, sql, args...)
}

func ExecuteUpdateTx(x *TX, sql string, args ...interface{}) (int64, error) {
	return executeUpdate_(x, db, sql, args...)
}

func executeUpdate_(x *TX, dbsource *sql.DB, sqlstr string, args ...interface{}) (int64, error) {
	var rs sql.Result
	var err error
	if x != nil && x.tx != nil && x.isBg {
		rs, err = x.tx.Exec(sqlstr, args...)
	} else {
		rs, err = dbsource.Exec(sqlstr, args...)
	}
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}

func Execute(dbsource *sql.DB, sql string, args ...interface{}) (int64, error) {
	rs, err := dbsource.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}

func Exec(sql string, args ...interface{}) (int64, error) {
	rs, err := db.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}

func (t *Table) getDB() *sql.DB {
	var dbsource *sql.DB
	if t.DB != nil {
		dbsource = t.DB
	} else {
		db_, ok := dbMap[t.TableName]
		if !ok {
			dbsource = db
		} else {
			dbsource = db_
		}
	}
	return dbsource
}

func (t *Table) executeQuerySingle() ([]interface{}, error) {
	rows, err := t.getDB().Query(t.sql, t.args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		cols, _ := rows.Columns()
		buff := make([]interface{}, len(cols))
		data := make([]interface{}, len(cols))
		for i, _ := range buff {
			buff[i] = &data[i]
		}
		row_err := rows.Scan(buff...)
		if row_err != nil {
			return nil, row_err
		}
		return data, nil
	} else {
		return nil, nil
	}

}

func GetValue(data *interface{}) interface{} {
	switch (*data).(type) {
	case int64:
		return (*data).(int64)
	case int32:
		return (*data).(int32)
	case int16:
		return (*data).(int16)
	case int8:
		return (*data).(int8)
	case uint64:
		return (*data).(uint64)
	case uint32:
		return (*data).(uint32)
	case uint16:
		return (*data).(uint16)
	case uint8:
		return (*data).(uint8)
	case int:
		return (*data).(int)
	case float32:
		return (*data).(float32)
	case float64:
		return (*data).(float64)
	case []uint8:
		return string((*data).([]uint8))
	default:
		return (*data)
	}
}

func getTypeString(data *interface{}) string {
	switch (*data).(type) {
	case int64:
		return "int64"
	case int32:
		return "int32"
	case int16:
		return "int16"
	case int8:
		return "int8"
	case uint64:
		return "uint64"
	case uint32:
		return "uint32"
	case uint16:
		return "uint16"
	case uint8:
		return "uint8"
	case int:
		return "int"
	case float32:
		return "float32"
	case float64:
		return "float64"
	case []uint8:
		return "string"
	default:
		return "string"
	}
}

func getValues(data []interface{}) []interface{} {
	retdata := make([]interface{}, 0)
	for _, d := range data {
		switch d.(type) {
		case int64:
			retdata = append(retdata, d)
		case []uint8:
			retdata = append(retdata, string(d.([]uint8)))
		default:
			retdata = append(retdata, d)
		}
	}
	return retdata
}

func (t *Table) GroupBy(columns ...Column) {
	ss := make([]string, 0, len(columns))
	for _, v := range columns {
		ss = append(ss, v.Name())
	}
	t.groupBySql = " group by " + strings.Join(ss, ",")
}

func (t *Table) Having(havings ...*Having) {
	ss := make([]string, 0, len(havings))
	for _, w := range havings {
		ss = append(ss, w.havingSql)
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
}

func (t *Table) OrderBy(sorts ...*Sort) {
	ss := make([]string, 0, len(sorts))
	for _, v := range sorts {
		ss = append(ss, v.OrderByArg)
	}
	t.orderBySql = " order by " + strings.Join(ss, ",")
}

func (t *Table) Limit(from, nums int32) {
	if nums != 0 {
		t.limitSql = " limit ?,? "
		t.args = append(t.args, from, nums)
	} else {
		t.limitSql = " ? "
		t.args = append(t.args, from)
	}
}

func (t *Table) Update() (int64, error) {
	modifystr := make([]string, 0)
	args := make([]interface{}, 0)
	for k, v := range t.ModifyMap {
		modifystr = append(modifystr, k+"=?")
		args = append(args, v)
	}
	t.modifySql = " update " + t.TableName + " set " + strings.Join(modifystr, ",")
	for _, v := range t.args {
		args = append(args, v)
	}
	t.args = args
	t.completeSql4Update()
	return executeUpdate_(t.Tx, t.getDB(), t.sql, t.args...)
}

func (t *Table) Insert() (int64, error) {
	insertField := make([]string, 0)
	insert_ := make([]string, 0)
	args := make([]interface{}, 0)
	for k, v := range t.ModifyMap {
		insertField = append(insertField, k)
		insert_ = append(insert_, "?")
		args = append(args, v)
	}
	t.sql = " insert  into " + t.TableName + "(" + strings.Join(insertField, ",") + " )values(" + strings.Join(insert_, ",") + ")"
	for _, v := range t.args {
		args = append(args, v)
	}
	t.args = args
	t.logger(t.sql, t.args)
	return executeUpdate_(t.Tx, t.getDB(), t.sql, t.args...)
}

func (t *Table) Delete() (int64, error) {
	t.modifySql = " delete from " + t.TableName
	t.completeSql4Update()
	return executeUpdate_(t.Tx, t.getDB(), t.sql, t.args...)
}

func (f *Field) EQ(arg interface{}) *Where {
	return &Where{f.FieldName + "=?", arg, nil}
}

func (f *Field) NEQ(arg interface{}) *Where {
	return &Where{f.FieldName + "<>?", arg, nil}
}

func (f *Field) LT(arg interface{}) *Where {
	return &Where{f.FieldName + "<?", arg, nil}
}

func (f *Field) LE(arg interface{}) *Where {
	return &Where{f.FieldName + "<=?", arg, nil}
}

func (f *Field) GT(arg interface{}) *Where {
	return &Where{f.FieldName + ">?", arg, nil}
}

func (f *Field) GE(arg interface{}) *Where {
	return &Where{f.FieldName + ">=?", arg, nil}
}

func (f *Field) LIKE(arg interface{}) *Where {
	return &Where{f.FieldName + " like %?%", arg, nil}
}

func (f *Field) RLIKE(arg interface{}) *Where {
	return &Where{f.FieldName + " like ?%", arg, nil}
}

func (f *Field) LLIKE(arg interface{}) *Where {
	return &Where{f.FieldName + " like %?", arg, nil}
}

func (f *Field) Between(from, to interface{}) *Where {
	return &Where{f.FieldName + " between ? and ?", nil, []interface{}{from, to}}
}

func (f *Field) IN(args ...interface{}) *Where {
	s := make([]string, len(args))
	for i := 0; i < len(args); i++ {
		s[i] = "?"
	}
	return &Where{f.FieldName + " in (" + strings.Join(s, ",") + ")", nil, args}
}

func (f *Field) NOTIN(args ...interface{}) *Where {
	s := make([]string, len(args))
	for i := 0; i < len(args); i++ {
		s[i] = "?"
	}
	return &Where{f.FieldName + " not in (" + strings.Join(s, ",") + ")", nil, args}
}

func (f *Field) Asc() *Sort {
	return &Sort{f.FieldName + " asc "}
}

func (f *Field) Desc() *Sort {
	return &Sort{f.FieldName + " desc "}
}

func (f *Field) Count(aliasName ...string) *SetOperation {
	if len(aliasName) == 1 {
		return &SetOperation{fieldName: " count(" + f.FieldName + ") as " + aliasName[0] + " "}
	} else {
		return &SetOperation{fieldName: " count(" + f.FieldName + ") "}
	}
}

func (f *Field) Distinct(aliasName ...string) *SetOperation {
	if len(aliasName) == 1 {
		return &SetOperation{fieldName: " distinct " + f.FieldName + " as " + aliasName[0] + " "}
	} else {
		return &SetOperation{fieldName: " distinct " + f.FieldName + " "}
	}
}

func (f *Field) Sum(aliasName ...string) *SetOperation {
	if len(aliasName) == 1 {
		return &SetOperation{fieldName: " sum(" + f.FieldName + ") as " + aliasName[0] + " "}
	} else {
		return &SetOperation{fieldName: " sum(" + f.FieldName + ") "}
	}
}

func (f *Field) Avg(aliasName ...string) *SetOperation {
	if len(aliasName) == 1 {
		return &SetOperation{fieldName: " avg(" + f.FieldName + ") as " + aliasName[0] + " "}
	} else {
		return &SetOperation{fieldName: " avg(" + f.FieldName + ") "}
	}
}

func (f *Field) Max(aliasName ...string) *SetOperation {
	if len(aliasName) == 1 {
		return &SetOperation{fieldName: " max(" + f.FieldName + ") as " + aliasName[0] + " "}
	} else {
		return &SetOperation{fieldName: " max(" + f.FieldName + ") "}
	}
}

func (f *Field) Min(aliasName ...string) *SetOperation {
	if len(aliasName) == 1 {
		return &SetOperation{fieldName: " min(" + f.FieldName + ") as " + aliasName[0] + " "}
	} else {
		return &SetOperation{fieldName: " min(" + f.FieldName + ") "}
	}
}

func (f *Field) Operation(qurey4SetOperation string) *SetOperation {
	return &SetOperation{fieldName: " " + qurey4SetOperation + " "}
}

func (s *SetOperation) EQ(arg interface{}) *Having {
	return &Having{s.fieldName + "=?", arg, nil}
}

func (s *SetOperation) NEQ(arg interface{}) *Having {
	return &Having{s.fieldName + "<>?", arg, nil}
}

func (s *SetOperation) LT(arg interface{}) *Having {
	return &Having{s.fieldName + "<?", arg, nil}
}

func (s *SetOperation) LE(arg interface{}) *Having {
	return &Having{s.fieldName + "<=?", arg, nil}
}

func (s *SetOperation) GT(arg interface{}) *Having {
	return &Having{s.fieldName + ">?", arg, nil}
}

func (s *SetOperation) GE(arg interface{}) *Having {
	return &Having{s.fieldName + ">=?", arg, nil}
}

func (s *SetOperation) Between(from, to interface{}) *Having {
	return &Having{s.fieldName + " between ? and ?", nil, []interface{}{from, to}}
}

func (s *SetOperation) Name() string {
	return s.fieldName
}

func (s *SetOperation) Value() interface{} {
	return GetValue(&s.FieldValue)
}

func (t *Table) SetCommentLine(commentline string) {
	t.commentline = commentline
}

func getTableColumnInfo(tablName string) *map[string][2]string {
	rows, _ := db.Query(getAdapterSqlStr() + tablName)
	defer rows.Close()
	mapname := make(map[string][2]string)
	for rows.Next() {
		var column_name string
		var data_type string
		rows.Scan(&column_name, &data_type, nil, nil, nil, nil)
		mapname[column_name] = getTypeStrs(data_type)
	}
	return &mapname
}

func getAdapterSqlStr() string {
	switch adapter {
	case MYSQL:
		return _MYSQL_TABLE_SCHMAINFO_SQL
	case PostgreSQL:
		return _PostgreSQL_TABLE_SCHMAINFO_SQL
	default:
		return _MYSQL_TABLE_SCHMAINFO_SQL
	}
}

func CreateDaoFile(tableName, packageName, destPath string) error {
	str := createFile(tableName, getTableColumnInfo(tableName), packageName)
	fileName := destPath + "/" + tableName + ".go"
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		return err
	}
	f.WriteString(str)
	log.Println("[create dao file] [" + fileName + "]")
	return nil
}

func createFile(table string, columnMap *map[string][2]string, packageName string) string {
	tableName := ToUpperFirstLetter(table)
	fileContent := "package " + packageName + "\n\n"
	fileContent = fileContent + "/**\n"
	fileContent = fileContent + "tablename:" + table + "\n"
	t := time.Now()
	fileContent = fileContent + "datetime :" + t.Format("2006-01-02 15:04:05") + "\n"
	fileContent = fileContent + "*/\n"
	fileContent = fileContent + "import (\n"
	fileContent = fileContent + "\t\"github.com/donnie4w/gdao\"\n"
	fileContent = fileContent + "\t\"reflect\"\n"
	fileContent = fileContent + ")\n\n"

	for field, data_type := range *columnMap {
		f := ToUpperFirstLetter(field)
		fileContent = fileContent + "type " + table + "_" + f + " struct {\n"
		fileContent = fileContent + "\tgdao.Field\n"
		fileContent = fileContent + "\tfieldName  string\n"
		fileContent = fileContent + "\tFieldValue *" + data_type[0] + "\n"
		fileContent = fileContent + "}\n\n"

		fileContent = fileContent + "func (c *" + table + "_" + f + ") Name() string {\n"
		fileContent = fileContent + "\treturn c.fieldName\n"
		fileContent = fileContent + "}\n\n"

		fileContent = fileContent + "func (c *" + table + "_" + f + ") Value() interface{} {\n"
		fileContent = fileContent + "\treturn c.FieldValue\n"
		fileContent = fileContent + "}\n\n"
	}
	fileContent = fileContent + "type " + tableName + " struct {\n"
	fileContent = fileContent + "\tgdao.Table\n"
	for field, _ := range *columnMap {
		f := ToUpperFirstLetter(field)
		fileContent = fileContent + "\t" + f + " *" + table + "_" + f + "\n"
	}
	fileContent = fileContent + "}\n\n"

	for field, data_type := range *columnMap {
		f := ToUpperFirstLetter(field)
		fileContent = fileContent + "func (u *" + tableName + ") Get" + f + "() " + data_type[0] + " {\n"
		fileContent = fileContent + "\treturn *u." + f + ".FieldValue\n}\n\n"
		fileContent = fileContent + "func (u *" + tableName + ") Set" + f + "(arg " + data_type[1] + ") {\n"
		fileContent = fileContent + "\tu.Table.ModifyMap[u." + f + ".fieldName] = arg\n"
		fileContent = fileContent + "\tv := " + data_type[0] + "(arg)\n"
		fileContent = fileContent + "\tu." + f + ".FieldValue = &v\n"
		fileContent = fileContent + "}\n\n"
	}
	fileContent = fileContent + "func (t *" + tableName + ") Query(columns ...gdao.Column) ([]" + tableName + ",error) {\n"
	fileContent = fileContent + "\tif columns == nil {\n"
	fs := make([]string, 0)
	for field, _ := range *columnMap {
		f := ToUpperFirstLetter(field)
		fs = append(fs, "t."+f)
	}
	fileContent = fileContent + "\t\tcolumns = []gdao.Column{ " + strings.Join(fs, ",") + "}\n"
	fileContent = fileContent + "\t}\n"
	fileContent = fileContent + "\trs,err := t.Table.Query(columns...)\n"
	fileContent = fileContent + "\tif rs == nil || err != nil {\n"
	fileContent = fileContent + "\t\treturn nil, err\n"
	fileContent = fileContent + "\t}\n"
	fileContent = fileContent + "\tts := make([]" + tableName + ", 0, len(rs))\n"
	fileContent = fileContent + "\tc := make(chan int16,len(rs))\n"
	fileContent = fileContent + "\tfor _, rows := range rs {\n"
	fileContent = fileContent + "\t\tt := New" + tableName + "()\n"
	fileContent = fileContent + "\t\tgo copy" + tableName + "(c, rows, t, columns)\n"
	fileContent = fileContent + "\t\t<-c\n"
	fileContent = fileContent + "\t\tts = append(ts, *t)\n"
	fileContent = fileContent + "\t}\n"
	fileContent = fileContent + "\treturn ts,nil\n"
	fileContent = fileContent + "}\n\n"

	copy := `func copy` + tableName + `(channle chan int16, rows []interface{}, t *` + tableName + `, columns []gdao.Column) {
	defer func() { channle <- 1 }()
	for j, core := range rows {
		if core == nil {
			continue
		}
		field := columns[j].Name()
		setfield := "Set" + gdao.ToUpperFirstLetter(field)
		reflect.ValueOf(t).MethodByName(setfield).Call([]reflect.Value{reflect.ValueOf(gdao.GetValue(&core))})
	}
}`
	fileContent = fileContent + copy + "\n\n"
	fileContent = fileContent + "func (t *" + tableName + ") QuerySingle(columns ...gdao.Column) (*" + tableName + ",error) {\n"
	fileContent = fileContent + "\tif columns == nil {\n"
	fileContent = fileContent + "\t\tcolumns = []gdao.Column{ " + strings.Join(fs, ",") + "}\n"
	fileContent = fileContent + "\t}\n"
	fileContent = fileContent + "\trs,err := t.Table.QuerySingle(columns...)\n"
	fileContent = fileContent + "\tif rs == nil || err != nil {\n"
	fileContent = fileContent + "\t\treturn nil, err\n"
	fileContent = fileContent + "\t}\n"
	fileContent = fileContent + "\trt := New" + tableName + "()\n"
	fileContent = fileContent + "\tfor j, core := range rs {\n"
	fileContent = fileContent + "\t\tif core == nil {\n"
	fileContent = fileContent + "\t\t\tcontinue\n"
	fileContent = fileContent + "\t\t}\n"
	fileContent = fileContent + "\t\tfield := columns[j].Name()\n"
	fileContent = fileContent + "\t\tsetfield := \"Set\" + gdao.ToUpperFirstLetter(field)\n"
	fileContent = fileContent + "\t\treflect.ValueOf(rt).MethodByName(setfield).Call([]reflect.Value{reflect.ValueOf(gdao.GetValue(&core))})\n"
	fileContent = fileContent + "\t}\n"
	fileContent = fileContent + "\treturn rt,nil\n"
	fileContent = fileContent + "}\n\n"

	fileContent = fileContent + "func (t *" + tableName + ") Select(columns ...gdao.Column) (*" + tableName + ",error) {\n"
	fileContent = fileContent + "\tif columns == nil {\n"
	fileContent = fileContent + "\t\tcolumns = []gdao.Column{ " + strings.Join(fs, ",") + "}\n"
	fileContent = fileContent + "\t}\n"
	fileContent = fileContent + "\trows,err := t.Table.Selects(columns...)\n"
	fileContent = fileContent + "\tdefer rows.Close()\n"
	fileContent = fileContent + "\tif err != nil || rows==nil {\n"
	fileContent = fileContent + "\t\treturn nil, err\n"
	fileContent = fileContent + "\t}\n"
	fileContent = fileContent + "\tbuff := make([]interface{}, len(columns))\n"
	fileContent = fileContent + "\tif rows.Next() {\n"
	fileContent = fileContent + "\t\tn := New" + tableName + "()\n"
	fileContent = fileContent + "\t\tcp" + tableName + "(buff, n, columns)\n"
	fileContent = fileContent + "\t\trow_err := rows.Scan(buff...)\n"
	fileContent = fileContent + "\t\tif row_err != nil {\n"
	fileContent = fileContent + "\t\t\treturn nil, row_err\n"
	fileContent = fileContent + "\t\t}\n"
	fileContent = fileContent + "\t\treturn n, nil\n"
	fileContent = fileContent + "\t}\n"
	fileContent = fileContent + "\treturn nil, nil\n"
	fileContent = fileContent + "}\n\n"

	fileContent = fileContent + "func (t *" + tableName + ") Selects(columns ...gdao.Column) ([]*" + tableName + ",error) {\n"
	fileContent = fileContent + "\tif columns == nil {\n"
	fileContent = fileContent + "\t\tcolumns = []gdao.Column{ " + strings.Join(fs, ",") + "}\n"
	fileContent = fileContent + "\t}\n"
	fileContent = fileContent + "\trows,err := t.Table.Selects(columns...)\n"
	fileContent = fileContent + "\tdefer rows.Close()\n"
	fileContent = fileContent + "\tif err != nil || rows==nil {\n"
	fileContent = fileContent + "\t\treturn nil, err\n"
	fileContent = fileContent + "\t}\n"
	fileContent = fileContent + "\tns := make([]*" + tableName + ", 0)\n"
	fileContent = fileContent + "\tbuff := make([]interface{}, len(columns))\n"
	fileContent = fileContent + "\tfor rows.Next() {\n"
	fileContent = fileContent + "\t\tn := New" + tableName + "()\n"
	fileContent = fileContent + "\t\tcp" + tableName + "(buff, n, columns)\n"
	fileContent = fileContent + "\t\trow_err := rows.Scan(buff...)\n"
	fileContent = fileContent + "\t\tif row_err != nil {\n"
	fileContent = fileContent + "\t\t\treturn nil, row_err\n"
	fileContent = fileContent + "\t\t}\n"
	fileContent = fileContent + "\t\tns = append(ns, n)\n"
	fileContent = fileContent + "\t}\n"
	fileContent = fileContent + "\treturn ns, nil\n"
	fileContent = fileContent + "}\n\n"

	fileContent = fileContent + "func  cp" + tableName + "(buff []interface{}, t *" + tableName + ", columns []gdao.Column) {\n"
	fileContent = fileContent + "\tfor i, column := range columns {\n"
	fileContent = fileContent + "\t\tfield := column.Name()\n"
	fileContent = fileContent + "\t\tswitch field {\n"
	for field, _ := range *columnMap {
		fileContent = fileContent + "\t\tcase \"" + field + "\":\n"
		fileContent = fileContent + "\t\t\tbuff[i] = &t." + ToUpperFirstLetter(field) + ".FieldValue\n"
	}
	fileContent = fileContent + "\t\t}\n"
	fileContent = fileContent + "\t}\n"
	fileContent = fileContent + "}\n\n"

	fileContent = fileContent + "func New" + tableName + "(tableName ...string) *" + tableName + " {\n"
	for field, _ := range *columnMap {
		f := ToUpperFirstLetter(field)
		fileContent = fileContent + "\t" + checkReserveKey(field) + " := &" + table + "_" + f + "{fieldName: \"" + field + "\"}\n"
		fileContent = fileContent + "\t" + checkReserveKey(field) + ".Field.FieldName = \"" + field + "\"\n"
	}
	fileContent = fileContent + "\ttable := &" + tableName + "{"
	ss := make([]string, 0)
	for field, _ := range *columnMap {
		f := ToUpperFirstLetter(field)
		ss = append(ss, f+":"+checkReserveKey(field))
	}
	fileContent = fileContent + strings.Join(ss, ",") + "}\n"
	fileContent = fileContent + "\ttable.Table.ModifyMap = make(map[string]interface{})\n"
	fileContent = fileContent + "\tif len(tableName) == 1 {\n"
	fileContent = fileContent + "\t\ttable.Table.TableName = tableName[0]\n"
	fileContent = fileContent + "\t} else {\n"
	fileContent = fileContent + "\t\ttable.Table.TableName = \"" + table + "\"\n"
	fileContent = fileContent + "\t}\n"
	fileContent = fileContent + "\treturn table\n"
	fileContent = fileContent + "}\n"
	return fileContent
}

func checkReserveKey(k string) string {
	b, err := regexp.MatchString("break|default|func|interface|select|case|defer|go|map|struct|chan|else|goto|package|switch|const|fallthrough|if|range|type|continue|for|import|return|var", k)
	exception(err)
	if b {
		return k + "_"
	}
	return k
}

func getTypeStrs(t string) [2]string {
	switch adapter {
	case MYSQL:
		return typeOfMysql2go(t)
	case PostgreSQL:
		return typePostgreSQL2go(t)
	default:
		return typeOfMysql2go(t)
	}
}

func typeOfMysql2go(t string) [2]string {
	t = strings.Replace(t, regexp.MustCompile("\\(\\d*?\\)").FindString(t), "", -1)
	switch strings.ToUpper(t) {
	case "CHAR":
		return [...]string{"string", "string"}
	case "VARCHAR":
		return [...]string{"string", "string"}
	case "TINYTEXT":
		return [...]string{"string", "string"}
	case "TEXT":
		return [...]string{"string", "string"}
	case "MEDIUMTEXT":
		return [...]string{"string", "string"}
	case "LONGTEXT":
		return [...]string{"string", "string"}
	case "BIT":
		return [...]string{"int64", "int64"}
	case "TINYINT":
		return [...]string{"int16", "int64"}
	case "BOOL":
		return [...]string{"int16", "int64"}
	case "BOOLEAN":
		return [...]string{"int16", "int64"}
	case "SMALLINT":
		return [...]string{"int32", "int64"}
	case "MEDIUMINT":
		return [...]string{"int32", "int64"}
	case "INT":
		return [...]string{"int32", "int64"}
	case "INTEGER":
		return [...]string{"int64", "int64"}
	case "BIGINT":
		return [...]string{"int64", "int64"}
	case "FLOAT":
		return [...]string{"float32", "float64"}
	case "DOUBLE":
		return [...]string{"float64", "float64"}
	case "DECIMAL":
		return [...]string{"float64", "float64"}
	case "DATE":
		return [...]string{"string", "string"}
	case "DATETIME":
		return [...]string{"string", "string"}
	case "TIMESTAMP":
		return [...]string{"string", "string"}
	case "TIME":
		return [...]string{"string", "string"}
	case "YEAR":
		return [...]string{"string", "string"}
	default:
		return [...]string{"string", "string"}
	}
}

func typePostgreSQL2go(t string) [2]string {
	switch strings.ToUpper(t) {
	case "SMALLINT":
		return [...]string{"int16", "int64"}
	case "INTEGER":
		return [...]string{"int32", "int64"}
	case "BIGINT":
		return [...]string{"int64", "int64"}
	case "NUMERIC":
		return [...]string{"float64", "float64"}
	case "REAL":
		return [...]string{"float32", "float64"}
	case "DOUBLE":
		return [...]string{"float64", "float64"}
	case "SERIAL":
		return [...]string{"int64", "int64"}
	case "BIGSERIAL":
		return [...]string{"int64", "int64"}
	case "VARCHAR":
		return [...]string{"string", "string"}
	case "CHAR":
		return [...]string{"string", "string"}
	case "TEXT":
		return [...]string{"string", "string"}
	case "TIMESTAMP":
		return [...]string{"string", "string"}
	case "DATE":
		return [...]string{"string", "string"}
	case "TIME":
		return [...]string{"string", "string"}
	default:
		return [...]string{"string", "string"}
	}
}
