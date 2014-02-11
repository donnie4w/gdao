package main

import (
	"database/sql"
	"example/dao"
	"fmt"
	"github.com/donnie4w/gdao"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func initdb() {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/dbtest")
	if err != nil {
		fmt.Println("any error on open database ", err.Error())
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	//	defer db.Close()
	gdao.SetDB(db)
}

func getDB() *sql.DB {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/dbtest")
	//defer db.Close()
	if err != nil {
		fmt.Println("any error on open database ", err.Error())
		return nil
	}
	return db
}

func init() {
	fmt.Println("init()")
	initdb()
}

//创建表
func createTable() {
	sql := "CREATE TABLE `hstest` ( `id` int(10) NOT NULL DEFAULT '-1' COMMENT 'id',   `name` varchar(20) NOT NULL COMMENT '名字',   `age` int(10) NOT NULL DEFAULT '-1' COMMENT '年龄',   `createtime` datetime NOT NULL DEFAULT '1900-01-01 00:00:00' COMMENT '创建时间'  ) ENGINE=InnoDB DEFAULT CHARSET=utf8"
	i, err := gdao.ExecuteUpdate(sql)
	fmt.Println(">>>>", i)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func goBeenTest1() {
	hstest := dao.NewHstest()
	//设置打印日志，默认false
	hstest.IsLog(true)
	// sql注释部分
	hstest.SetCommentLine("/*master*/")
	//查询自动非表字段名时，使用QueryBeen，返回([]*GoBeen, error)
	gb, err := hstest.QueryBeen(hstest.Id.Count(), hstest.Age)
	if err != nil {
		fmt.Println(err.Error())
	}
	//取出[]*GoBeen中的值
	for _, g := range gb {
		//GoBeen中属性FieldBeen包含了查询的字段名，字段索引，字段值
		for _, v := range g.FieldBeens {
			fmt.Println(v.Name(), v.Index(), v.Value())
		}
		//可以直接使用字段名来获取值
		fmt.Println(g.MapName("age").Value())
		//可以直接使用字段索引来获取值
		fmt.Println(g.MapIndex(1).Value())
	}
}

func goBeenTest2() {
	i := 1
	for i > 0 {
		fmt.Println(">>>>>>>>", i)
		//支持原生sql查询，返回[]*GoBeen，取值方式与上同
		gbs, err := gdao.ExecuteQuery("select * from hstest limit ?", 1)
		if err != nil {
			fmt.Println(err.Error())
		}
		for _, g := range gbs {
			for _, v := range g.FieldBeens {
				fmt.Println(v.Name(), v.Index(), v.Value())
			}
			fmt.Println(g.MapName("id").Index(), g.MapName("id").Name(), g.MapName("id").Value())
			fmt.Println(g.MapIndex(2).Index(), g.MapIndex(2).Name(), g.MapIndex(2).Value())
			i--
		}
	}
}

func queryTest() {
	//使用时先初始化
	hstest := dao.NewHstest()
	hstest.IsLog(true)
	hstest.Where(hstest.Id.Between(1, 10), hstest.Age.IN(30, 31, 32, 33, 34))
	hstest.GroupBy(hstest.Id)
	hstest.Having(hstest.Id.Count().GT(0))
	hstest.OrderBy(hstest.Id.Asc())
	hstest.Limit(0, 3)
	//Query方法查询将返回[]Hstest
	hstests, err := hstest.Query(hstest.Id, hstest.Age, hstest.Createtime, hstest.Name)
	if err != nil {
		fmt.Println(err.Error())
	}
	//取值
	for _, u := range hstests {
		//使用字段的Get方法 Get***()
		fmt.Println(">>>>", u.GetId(), u.GetAge(), u.GetCreatetime(), u.GetName())
	}
}

func queryTest2() {
	hstest := dao.NewHstest()
	hstest.IsLog(true)
	hstest.DB = getDB()
	hstest.SetCommentLine("/*master*/")
	hstest.Where(hstest.Id.EQ(3), hstest.Age.IN(30, 31, 32, 33, 34))
	hstest.GroupBy(hstest.Id)
	hstest.Having(hstest.Id.Count().GT(0))
	hstest.OrderBy(hstest.Id.Asc())
	hstest.Limit(0, 1)
	//QuerySingle方法返回一行数据
	hstest, _ = hstest.QuerySingle(hstest.Id, hstest.Age, hstest.Createtime, hstest.Name)
	fmt.Println(">>>>", hstest.GetId(), hstest.GetAge(), hstest.GetCreatetime(), hstest.GetName())
}

func updateTest() {
	fmt.Println("updateTest() ")
	hstest := dao.NewHstest()
	hstest.IsLog(true)
	hstest.SetName("wuxiaodong3")
	hstest.SetAge(34)
	hstest.Where(hstest.Id.EQ(2))
	//更新
	hstest.Update()
}

func inserteTest() {
	fmt.Println("inserteTest() ")
	hstest := dao.NewHstest()
	hstest.IsLog(true)
	hstest.SetId(1)
	hstest.SetName("wuxiaodong")
	hstest.SetAge(30)
	hstest.SetCreatetime(time.Now())
	//插入数据
	i, err := hstest.Insert()
	fmt.Println(i)
	if err != nil {
		fmt.Println(err.Error())
	}

}

func deleteTest() {
	fmt.Println("deleteTest() ")
	hstest := dao.NewHstest()
	hstest.IsLog(true)
	hstest.Where(hstest.Id.EQ(1))
	//删除数据
	hstest.Delete()
}

func createDaoTest() {
	//创建数据表hstest相应的hstest.go文件，包名为dao,路径为d:/gdao/src/dao
	err := gdao.CreateDaoFile("hstest", "dao", "d:/gdao/src/dao")
	if err != nil {
		fmt.Println(err.Error())
	}
}

//事务测试
func txTest() {
	tx := gdao.GetTX()
	tx.Begin()
	hstest := dao.NewHstest()
	hstest.SetTx(tx)
	hstest.SetId(13)
	hstest.SetName("wuxiaodong")
	hstest.SetAge(30)
	hstest.SetCreatetime(time.Now())
	//插入数据
	hstest.Insert()
	c1 := make(chan int32, 2)

	go func(c *chan int32) {
		defer func() {
			*c <- 1
		}()
		hstest2 := dao.NewHstest()
		hstest2.SetTx(tx)
		hstest2.Where(hstest2.Id.EQ(1))
		//删除数据
		hstest2.Delete()
	}(&c1)

	go func(c *chan int32) {
		defer func() {
			*c <- 2
		}()
		sqlstr := "insert into hstest(id,age,name)values(?,?,?)"
		gdao.ExecuteUpdateTx(tx, sqlstr, 20, 41, "wuxiaodong")
	}(&c1)
	//go inserteTest()
	//go inserteTest()
	//go inserteTest()

	fmt.Println(">>>>>", <-c1)
	fmt.Println(">>>>>", <-c1)
	//tx.RollBack()
	tx.Commit()
}

func main() {
	fmt.Println("main()")
	//createTable()
	//goBeenTest1()
	goBeenTest2()
	//queryTest()
	//queryTest2()
	//updateTest()
	//inserteTest()
	//deleteTest()
	//createDaoTest()
	//txTest()
}
