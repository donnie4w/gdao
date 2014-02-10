package dao

/**
tablename:hstest
datetime :2014-02-10 09:49:08
*/
import (
	"github.com/gdao"
	"reflect"
)

type hstest_Id struct {
	gdao.Field
	fieldName  string
	FieldValue interface{}
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
	FieldValue interface{}
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
	FieldValue interface{}
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
	FieldValue interface{}
}

func (c *hstest_Createtime) Name() string {
	return c.fieldName
}

func (c *hstest_Createtime) Value() interface{} {
	return c.FieldValue
}

type Hstest struct {
	gdao.Table
	Id *hstest_Id
	Name *hstest_Name
	Age *hstest_Age
	Createtime *hstest_Createtime
}

func (u *Hstest) GetId() interface{} {
	return u.Id.FieldValue
}

func (u *Hstest) SetId(arg interface{}) {
	u.Table.ModifyMap[u.Id.fieldName] = arg
	u.Id.FieldValue = arg
}

func (u *Hstest) GetName() interface{} {
	return u.Name.FieldValue
}

func (u *Hstest) SetName(arg interface{}) {
	u.Table.ModifyMap[u.Name.fieldName] = arg
	u.Name.FieldValue = arg
}

func (u *Hstest) GetAge() interface{} {
	return u.Age.FieldValue
}

func (u *Hstest) SetAge(arg interface{}) {
	u.Table.ModifyMap[u.Age.fieldName] = arg
	u.Age.FieldValue = arg
}

func (u *Hstest) GetCreatetime() interface{} {
	return u.Createtime.FieldValue
}

func (u *Hstest) SetCreatetime(arg interface{}) {
	u.Table.ModifyMap[u.Createtime.fieldName] = arg
	u.Createtime.FieldValue = arg
}

func (t *Hstest) Query(columns ...gdao.Column) ([]Hstest,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Id,t.Name,t.Age,t.Createtime}
	}
	rs,err := t.Table.Query(columns...)
	if err != nil {
		return nil, err
	}
	ts := make([]Hstest, 0, len(rs))
	for _, rows := range rs {
		t := NewHstest()
		for j, core := range rows {
			if core == nil {
				continue
			}
			field := columns[j].Name()
			setfield := "Set" + gdao.ToUpperFirstLetter(field)
			reflect.ValueOf(t).MethodByName(setfield).Call([]reflect.Value{reflect.ValueOf(gdao.GetValue(&core))})
		}
		ts = append(ts, *t)
	}
	return ts,nil
}

func (t *Hstest) QuerySingle(columns ...gdao.Column) (*Hstest,error) {
	if columns == nil {
		columns = []gdao.Column{ t.Id,t.Name,t.Age,t.Createtime}
	}
	rs,err := t.Table.QuerySingle(columns...)
	if err != nil {
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
	table := &Hstest{Id:id,Name:name,Age:age,Createtime:createtime}
	table.Table.ModifyMap = make(map[string]interface{})
	if len(tableName) == 1 {
		table.Table.TableName = tableName[0]
	} else {
		table.Table.TableName = "hstest"
	}
	return table
}
