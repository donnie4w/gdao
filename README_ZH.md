## GDAO —— Go 持久层框架  [[English](https://github.com/donnie4w/gdao/blob/master/README.md)]

### 简介

gdao是一种创新的持久层解决方案。主要目的在于 减少编程量，提高生产力，提高性能，支持多数据源整合操作，支持数据读写分离，制定持久层编程规范。 灵活运用gdao，可以在持久层设计上，减少30%甚至50%以上的编程量，同时形成持久层的统一编程规范，减少持久层错误，同时易于维护和扩展。
gdao对于go语言，相当于[hibernate](https://github.com/hibernate/hibernate-orm)+[mybatis](https://github.com/mybatis/mybatis-3) 对于java语言，gdao框架融合了Hibernate的抽象性和MyBatis的灵活性，并解决了它们各自在ORM框架上长久以来使用上的痛点。关于hibernate与mybatis在痛点问题，可以参考[jdao使用文档](https://tlnet.top/jdaodoc)
* gdao设计结构简洁且严谨，所有接口与函数均能见名知意。
* 即使从未接触过gdao，看到gdao持久层代码，也能马上知道它的代码表达的意思和相关的数据行为。
* 你可以在几分钟内，掌握它的用法，这是它的简洁性与设计规范性带来的优势。

### [官网](https://tlnet.top/gdao)

### [使用文档](https://tlnet.top/gdaodoc)

### [Demo](https://github.com/donnie4w/gdaodemo)

### 主要特点

1. **生成代码**：运行gdao代码生成工具，创建数据库表的标准化实体类。类似thrift/protobuf。
2. **高效序列化**：表的标准化实体类实现了高效的序列化与反序列化。性能更高，数据体积更小。
3. **支持数据读写分离**：gdao支持绑定多数据源，并支持数据源绑定表，类，映射接口等属性。并支持数据读写分离
4. **支持数据缓存**：gdao支持数据缓存，并支持对缓存数据时效，与数据回收等特性进行细致控制
5. **广泛兼容性**：gdao理论上支持所有实现 go 数据库驱动接口的数据库
6. **高级特性**：支持事务，存储过程，批处理等数据库操作
7. **支持SQL与程序分离**：类似 mybatis，gdao支持xml文件写sql映射调用。这是少有的orm支持SQL与程序分离的功能，但是该功能非常强大。


#### GDAO 创新的orm解决方案

Gdao是[Jdao](https://github.com/donnie4w/jdao)的功能等价框架，设计模式均来自[Jdao](https://github.com/donnie4w/jdao)。

1. **标准化映射实体类，处理单表CRUD操作**：90%以上的数据库单表操作，可用通过实体类操作完成。这些对单表的增删改查操作，一般不涉及复杂的SQL优化，由实体类封装生成，可以减少错误率，更易于维护。
   利用缓存，读写分离等机制，在优化持久层上，更为高效和方便
   标准化实体类的数据操作格式并非简单的对象函数拼接，而是更类似SQL操作的对象化，使得操作上更加易于理解。
2. **复杂SQL的执行**：在实践中发现，复杂SQL，特别是多表关联的SQL，通常需要优化，这需要对表结构，表索引性质等数据库属性有所了解。
   而将复杂SQL使用对象进行拼接，通常会增加理解上的难度。甚至，开发者都不知道对象拼接后的最终执行SQL是什么，这无疑增加了风险和维护难度。
   因此Gdao在复杂SQL问题上，建议调用Gdao的CURD接口执行，Gdao提供了灵活的数据转换和高效的对象映射实现，可以避免过渡使用反射等耗时的操作。
3. **Sql映射文件**： 对于复杂的sql操作，Gdao提供了相应的crud接口。同时也支持通过xml配置sql进行接口映射调用，这点与java的mybatis Orm框架相似，区别在于mybatis需要映射所有SQL操作，
   而Gdao虽然提供了完整的sql映射接口，但是建议只映射复杂SQL，或操作部分标准实体类无法完成的CURD操作。
   Gdao的SQL配置文件参考mybatis配置文件格式，实现自己新的解析器，使得配置参数在类型的容忍度上更高，更灵活。（可以参考文档）

------

![](https://tlnet.top/statics/tlnet/29197.jpg)

------
### 核心组件

#### 1. gdao

主要核心入口，提供以下功能：
- 设置数据源
- SQL CRUD 函数

#### 2. gdaoCache

缓存入口，支持以下功能：
- 绑定或移除包、类等属性，开启或移除它们的查询缓存

#### 3. gdaoSlave

读写分离操作入口，支持以下功能：
- 绑定或移除包、类等属性，开启或移除它们的读写分离操作

#### 4. gdaoMapper

执行Sql可以直接调用Gdao的CRUD接口。也可以用xml文件映射，通过gdaoMapper调用Mapper Id执行

------

## 快速入门

### 1. 安装

```bash
# gdao 导入
go get github.com/donnie4w/gdao
```

### 2. 配置数据源

```go
gdao.Init(mysqlDB,gdao.MYSQL)
// dataSource 为数据源
// gdao.MYSQL 为数据库类型
```

### 3. 生成表实体类

使用 Gdao 代码生成工具生成数据库表的标准化实体类。

### 4. 实体类操作

```go
//  数据源设置
gdao.init(mysqlDB,gdao.MYSQL)

// 读取
hs := dao.NewHstest()
hs.Where(hs.Id.EQ(10))
h, _ := hs.Select(hs.Id, hs.Value, hs.Rowname)
logger.Debug(h)
//[DEBUG][SELETE ONE][ select id,value,rowname from hstest where id=?][10] 

// 更新
hs := dao.NewHstest()
hs.SetRowname("hello10")
hs.Where(hs.Id.EQ(10))
hs.Update()
//[DEBUG][UPDATE][update hstest set rowname=? where id=?][hello10 10]

// 删除
hs := dao.NewHstest()
hs.Where(hs.Id.EQ(10))
t.delete()
//[DEBUG][UPDATE][delete from hstest where id=?][10]

//新增
hs := dao.NewHstest()
hs.SetValue("hello123")
hs.SetLevel(12345)
hs.SetBody([]byte("hello"))
hs.SetRowname("hello1234")
hs.SetUpdatetime(time.Now())
hs.SetFloa(123456)
hs.SetAge(123)
hs.Insert()
//[DEBUG][INSERT][insert  into hstest(floa,age,value,level,body,rowname,updatetime )values(?,?,?,?,?,?,?)][123456 123 hello123 12345 [104 101 108 108 111] hello1234 2024-07-17 19:36:44]
```

### 5. gdao api

###### CRUD操作

```go
//查询，返回单条
bean, _ := gdao.ExecuteQueryBean("select id,value,rowname from hstest where id=?", 10)
logger.Debug(bean)

//insert
int  i = gdao.ExecuteUpdate("insert into hstest2(rowname,value) values(?,?)", "helloWorld", "123456789");

//update
int  i = gdao.ExecuteUpdate("update hstest set value=? where id=1", "hello");

//delete
int  i = gdao.ExecuteUpdate("delete from hstest where id = ?", 1);
```

### 6. gdaoCache

###### 配置缓存

```go
//绑定Hstest  启用缓存, 缓存时效为 300秒
gdaoCache.BindClass[dao.Hstest]()

//Hstest 第一次查询，并根据条件设置数据缓存
hs := dao.NewHstest()
hs.Where((hs.Id.Between(0, 2)).Or(hs.Id.Between(10, 15)))
hs.Limit(3)
hs.Selects()

//相同条件，数据直接由缓存获取
hs = dao.NewHstest()
hs.Where((hs.Id.Between(0, 2)).Or(hs.Id.Between(10, 15)))
hs.Limit(3)
hs.Selects()

```

##### 执行结果

```text
[DEBUG][SELETE LIST][ select id,age,rowname,value,updatetime,body,floa,level from hstest where id between ? and ? or (id between ? and ?) LIMIT ? ][0 2 10 15 3]
[DEBUG][SET CACHE][ select id,age,rowname,value,updatetime,body,floa,level from hstest where id between ? and ? or (id between ? and ?) LIMIT ? ][0 2 10 15 3]
[DEBUG][SELETE LIST][ select id,age,rowname,value,updatetime,body,floa,level from hstest where id between ? and ? or (id between ? and ?) LIMIT ? ][0 2 10 15 3]
[DEBUG][GET CACHE][ select id,age,rowname,value,updatetime,body,floa,level from hstest where id between ? and ? or (id between ? and ?) LIMIT ? ][0 2 10 15 3]
```

### 7. gdaoSlave

###### 读写分离

```go
设置备库数据源：mysql
mysql, _ := getDataSource("mysql.json")
gdaoSlave.BindClass[dao.Hstest](mysql, gdao.MYSQL)
//这里主数据库为sqlite，备数据库为mysql，Hstest读取数据源为mysql
hs := dao.NewHstest()
hs.Where(hs.Id.Between(0, 5))
hs.OrderBy(hs.Id.Desc())
hs.Limit(3)
hs.Selects()
```

### 8. gdaoMapper

###### 使用 XML 映射 SQL

```xml
<!-- MyBatis 风格的 XML 配置文件 -->
<mapper namespace="user">
   <select id="selectHstest1" parameterType="int64" resultType="hstest1">
      SELECT * FROM hstest1  order by id desc limit #{limit}
   </select>
</mapper>
```

```go
//数据源设置
if db, err := getDataSource("sqlite.json"); err == nil {
   gdao.Init(db, gdao.SQLITE)  
   gdao.SetLogger(true)
}

//读取解析xml配置
hs1, _ := gdaoMapper.Select[dao.Hstest1]("user.selectHstest1", 1)
fmt.Println(hs1)

id, _ := gdaoMapper.Select[int64]("user.selectHstest1", 1)
fmt.Println(*id)
```

##### 执行结果

```text
[DEBUG][Mapper Id] user.selectHstest1 
SELECTONE SQL[SELECT * FROM hstest1  order by id desc limit ?]ARGS[1]
Id:52,Rowname:rowname>>>123456789,Value:[104 101 108 108 111 32 103 100 97 111],Goto:[49 50 51 52 53]
[DEBUG][Mapper Id] user.selectHstest1 
SELECTONE SQL[SELECT * FROM hstest1  order by id desc limit ?]ARGS[1]
52
```
