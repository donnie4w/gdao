gdao ��һ��golang��orm�⣬
gdao���Խ����ݿ��ӳ��������Ӧ��***.go�ļ�(����.go)��
֮��������ݿⵥ���ֱ�Ӳ�����Ӧ��go�ļ����ɡ�ͬʱ֧��ԭ��sql��䡣

����hstest���hstest.go�ļ�
	gdao.CreateDaoFile("hstest", "dao", "d:/gdao/src/example/dao")
 

��ѯ������select id,age,createtime,name from hstest where id between 1 and 10 and age in(30, 31)
	hstest := dao.NewHstest()
 	hstest.Where(hstest.Id.Between(1, 10), hstest.Age.IN(30, 31))
	hstests, _:= hstest.Query(hstest.Id, hstest.Age, hstest.Createtime, hstest.Name)
	for _, u := range hstests {
		fmt.Println(">>>>", u.GetId(), u.GetAge(), u.GetCreatetime(), u.GetName())
	}

���²�����update hstest set name="wu",age=34 where id=2
	hstest := dao.NewHstest()
	hstest.SetName("wu")
	hstest.SetAge(34)
	hstest.Where(hstest.Id.EQ(2))
	hstest.Update()

�������: insert into hstest(id,name,age)values(1,"wu",30,time.Now())
	hstest := dao.NewHstest()
	hstest.SetId(1)
	hstest.SetName("wu")
	hstest.SetAge(30)
	hstest.SetCreatetime(time.Now())
	hstest.Insert()

ɾ��������delete from hstest where id=1
	hstest := dao.NewHstest()
	hstest.Where(hstest.Id.EQ(1))
	hstest.Delete()
