gdao 是一个golang的orm库，
gdao可以将数据库表映射生成响应的***.go文件(表名.go)，
之后操作数据库单表就直接操作响应的go文件即可。同时支持原生sql语句。

生成hstest表的hstest.go文件
	gdao.CreateDaoFile("hstest", "dao", "d:/gdao/src/example/dao")
 

查询操作：select id,age,createtime,name from hstest where id between 1 and 10 and age in(30, 31)
	hstest := dao.NewHstest()
 	hstest.Where(hstest.Id.Between(1, 10), hstest.Age.IN(30, 31))
	hstests, _:= hstest.Query(hstest.Id, hstest.Age, hstest.Createtime, hstest.Name)
	for _, u := range hstests {
		fmt.Println(">>>>", u.GetId(), u.GetAge(), u.GetCreatetime(), u.GetName())
	}

更新操作：update hstest set name="wu",age=34 where id=2
	hstest := dao.NewHstest()
	hstest.SetName("wu")
	hstest.SetAge(34)
	hstest.Where(hstest.Id.EQ(2))
	hstest.Update()

插入操作: insert into hstest(id,name,age)values(1,"wu",30,time.Now())
	hstest := dao.NewHstest()
	hstest.SetId(1)
	hstest.SetName("wu")
	hstest.SetAge(30)
	hstest.SetCreatetime(time.Now())
	hstest.Insert()

删除操作：delete from hstest where id=1
	hstest := dao.NewHstest()
	hstest.Where(hstest.Id.EQ(1))
	hstest.Delete()
