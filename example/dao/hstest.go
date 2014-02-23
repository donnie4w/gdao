package dao

/**
tablename:hstest
datetime :2014-02-20 12:35:24
*/
import (
	"github.com/donnie4w/gdao"
	"reflect"
)

type hstest_Id struct {
	gdao.Field
	fieldName  string
	FieldValue *int32
}

func (c *hstest_Id) Name() string {
	return c.fieldName
}

func (c *hstest_Id) Value() interface{} {
	return c.FieldValue
}

type hstest_Name struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *hstest_Name) Name() string {
	return c.fieldName
}

func (c *hstest_Name) Value() interface{} {
	return c.FieldValue
}

type hstest_Age struct {
	gdao.Field
	fieldName  string
	FieldValue *int16
}

func (c *hstest_Age) Name() string {
	return c.fieldName
}

func (c *hstest_Age) Value() interface{} {
	return c.FieldValue
}

type hstest_Createtime struct {
	gdao.Field
	fieldName  string
	FieldValue *string
}

func (c *hstest_Createtime) Name() string {
	return c.fieldName
}

func (c *hstest_Createtime) Value() interface{} {
	return c.FieldValue
}

type hstest_Money struct {
	gdao.Field
	fieldName  string
	FieldValue *float32
}

func (c *hstest_Money) Name() string {
	return c.fieldName
}

func (c *hstest_Money) Value() interface{} {
	return c.FieldValue
}

type Hstest struct {
	gdao.Table
	Id *hstest_Id
	Name *hstest_Name
	Age *hstest_Age
	Createtime *hstest_Createtime
	Money *hstest_Money
}

func (u *Hstest) GetId() int32 {
	return *u.Id.FieldValue
}

func (u *Hstest) SetId(arg int64) {
	u.Table.ModifyMap[u.Id.fieldName] = arg
	v := int32(arg)
	u.Id.FieldValue = &v
}

func (u *Hstest) GetName() string {
	return *u.Name.FieldValue
}

func (u *Hstest) SetName(arg string) {
	u.Table.ModifyMap[u.Name.fieldName] = arg
	v := string(arg)
	u.Name.FieldValue = &v
}

func (u *Hstest) GetAge() int16 {
	return *u.Age.FieldValue
}

func (u *Hstest) SetAge(arg int64) {
	u.Table.ModifyMap[u.Age.fieldName] = arg
	v := int16(arg)
	u.Age.FieldValue = &v
}

func (u *Hstest) GetCreatetime() string {
	return *u.Createtime.FieldValue
}

func (u *Hstest) SetCreatetime(arg string) {
	u.Table.ModifyMap[u.Createtime.fieldName] = arg
	v := string(arg)
	u.Createtime.FieldValue = &v
}

func (u *Hstest) GetMoney() float32 {
	return *u.Money.FieldValue
}

func (u *Hstest) SetMoney(arg float64) {
	u.Table.ModifyMap[u.Money.fieldName] = arg
	v := float32(arg)
	u.Money.FieldValue = &v
}

func (t *Hstest) Query(columns ...gdao.Column) ([]Hstest,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Id,t.Name,t.Age,t.Createtime,t.Money}
	}
	rs,err := t.Table.Query(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	ts := make([]Hstest, 0, len(rs))
	for _, rows := range rs {
		t := NewHstest()
		c := make(chan int16)
		go copyHstest(c, rows, t, columns)
		<-c
		ts = append(ts, *t)
	}
	return ts,nil
}

func copyHstest(channle chan int16, rows []interface{}, t *Hstest, columns []gdao.Column) {
	defer func() { channle <- 1 }()
	for j, core := range rows {
		if core == nil {
			continue
		}
		field := columns[j].Name()
		setfield := "Set" + gdao.ToUpperFirstLetter(field)
		reflect.ValueOf(t).MethodByName(setfield).Call([]reflect.Value{reflect.ValueOf(gdao.GetValue(&core))})
	}
}

func (t *Hstest) QuerySingle(columns ...gdao.Column) (*Hstest,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Id,t.Name,t.Age,t.Createtime,t.Money}
	}
	rs,err := t.Table.QuerySingle(columns...)
	if rs == nil || err != nil {
		return nil, err
	}
	rt := NewHstest()
	for j, core := range rs {
		if core == nil {
			continue
		}
		field := columns[j].Name()
		setfield := "Set" + gdao.ToUpperFirstLetter(field)
		reflect.ValueOf(rt).MethodByName(setfield).Call([]reflect.Value{reflect.ValueOf(gdao.GetValue(&core))})
	}
	return rt,nil
}

func NewHstest(tableName ...string) *Hstest {
	id := &hstest_Id{fieldName: "id"}
	id.Field.FieldName = "id"
	name := &hstest_Name{fieldName: "name"}
	name.Field.FieldName = "name"
	age := &hstest_Age{fieldName: "age"}
	age.Field.FieldName = "age"
	createtime := &hstest_Createtime{fieldName: "createtime"}
	createtime.Field.FieldName = "createtime"
	money := &hstest_Money{fieldName: "money"}
	money.Field.FieldName = "money"
	table := &Hstest{Id:id,Name:name,Age:age,Createtime:createtime,Money:money}
	table.Table.ModifyMap = make(map[string]interface{})
	if len(tableName) == 1 {
		table.Table.TableName = tableName[0]
	} else {
		table.Table.TableName = "hstest"
	}
	return table
}
