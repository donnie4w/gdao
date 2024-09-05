## GDAO —— Go 持久层框架  [[English](https://github.com/donnie4w/gdao/blob/master/README.md)]

### 简介

gdao是一个全面的go持久层解决方案。主要目的在于 减少编程量，提高生产力，提高性能，支持多数据源整合操作，支持数据读写分离，制定持久层编程规范。 灵活运用gdao，可以在持久层设计上，减少30%甚至50%以上的编程量，同时形成持久层的统一编程规范，减少持久层错误，同时易于维护和扩展。
gdao对于go语言，相当于[hibernate](https://github.com/hibernate/hibernate-orm)+[mybatis](https://github.com/mybatis/mybatis-3) 对于java语言，gdao框架融合了Hibernate的抽象性和MyBatis的灵活性，并解决了它们各自在ORM框架上长久以来使用上的痛点。(框架痛点问题可参考[jdao使用文档](https://tlnet.top/jdaodoc))

gdao完整地在go语言中实现myBatis的核心功能，实现SQL与程序分离，实现强大的动态SQL功能，是目前唯一个用go语言完整实现SQL与程序分离的orm。


### [使用文档](https://tlnet.top/gdaodoc)

### [示例程序](https://github.com/donnie4w/gdaodemo)

------


### 主要特点
1. 完全表对象映射: 支持表操作映射为标准结构体的对象操作。
2. 完全SQL对象映射: 支持完全的SQL操作映射为对象操作。
3. SQL与程序分离: 支持SQL与程序分离，与[mybatis](https://github.com/mybatis/mybatis-3)相同。


### 主要功能

1. 生成代码：运行gdao代码生成工具，创建数据库表的标准化实体类。类似thrift/protobuf。
2. 高效序列化：表的标准化实体类实现了高效的序列化与反序列化。性能更高，数据体积更小。
3. 读写分离：gdao支持绑定多个备库数据源，并对绑定的表或SQL实现数据读写分离
4. 数据缓存：gdao支持数据缓存，并支持对缓存数据时效，与数据回收等特性进行细致控制
5. 广泛兼容性：gdao理论上支持所有实现 go 数据库驱动接口的数据库
6. 高级特性：支持事务，存储过程，批处理等数据库操作
7. 动态SQL：gdao实现丰富的动态SQL构建功能。支持动态SQL标签映射构建与原生SQL动态构建等多种模式。
8. myBatis特性：Gdao的映射模块是myBatis在go上的实现。也是所知的go orm中，唯一支持SQL与程序分离的orm。


### GDAO 持久层解决方案

**说明：** 完全的表映射与完全的SQL映射均可能导致一些问题。可能表现为对象过渡封装，大量重复相似行为的SQL等问题。gdao采用表映射的抽象性与SQL映射的半自动化管理的灵活性，解决持久层的问题。

1. **完全表对象映射，处理单表CRUD操作**: 将所有操作都用表对象映射操作是`hibernate`的经典做法，它导致了一些过渡封装的问题。
   Gdao支持自动生成表的映射实体类，专门处理单表的增删改查操作。它的底层动态SQL构建，有效地解决了相似的表操作行为带来的大量相似SQL的问题。
2. **完全SQL映射对象，处理复杂SQL操作**: 在实践中发现，复杂SQL，特别是多表关联的SQL，通常需要优化，这需要对表结构，表索引性质等数据库属性有所了解。
   而将复杂SQL字符串使用go对象进行拼接，通常会增加理解上的难度。甚至，最后开发者自己都不知道对象拼接后的最终执行SQL是什么，这无疑增加了风险和维护难度。
   Gdao实现了`mybatis`的核心功能SQL映射特性，可以在XML配置执行SQL，并映射为go对象。
3. **`SqlBuilder`，原生SQL动态构建**：Gdao实现了原生SQL的动态构建。对于极为复杂的SQL，甚至是function或存储过程等数据库高级特性，Gdao支持使用`SqlBuilder`进行动态构建。
   它是动态标签的程序性化实现，基于go程序实现动态SQL，可以构建出任何形式的SQL。

------

![](https://tlnet.top/statics/tlnet/29504.jpg)

------

# 快速入门

### 1. 安装

```bash
go get github.com/donnie4w/gdao
```

### 2. 配置数据源

```go
gdao.Init(mysqldb,gdao.MYSQL) //gdao初始化数据源  mysqldb 为sql.DB对象 ， gdao.MYSQL 为数据库类型
```

### 3. 标准化实体类操作

```go
// 读取
hs := dao.NewHstest()
hs.Where(hs.Id.EQ(10))
h, err := hs.Select(hs.Id, hs.Value, hs.Rowname)
//[DEBUG][SELETE ONE][ select id,value,rowname from hstest where id=?][10] 
```

```go
// 更新
hs := dao.NewHstest()
hs.SetRowname("hello10")
hs.Where(hs.Id.EQ(10))
hs.Update()
//[DEBUG][UPDATE][update hstest set rowname=? where id=?][hello10 10]
```

```go
// 删除
hs := dao.NewHstest()
hs.Where(hs.Id.EQ(10))
hs.delete()
//[DEBUG][UPDATE][delete from hstest where id=?][10]
```

```go
//新增
hs := dao.NewHstest()
hs.SetValue("hello123")
hs.SetLevel(12345)
hs.SetBody([]byte("hello"))
hs.SetRowname("hello1234")
hs.SetUpdatetime(time.Now())
hs.Insert()
//[DEBUG][INSERT][insert  into hstest(floa,age,value,level,body,rowname,updatetime )values(?,?,?,?,?)][hello123 12345 hello hello1234 2024-07-17 19:36:44]
```

### 4. 原生SQL操作

```go
//查询，返回单条 gdao是原生SQL操作入口
gdao.ExecuteQueryBean("select id,value,rowname from hstest where id=?", 10)

//insert
gdao.ExecuteUpdate("insert into hstest2(rowname,value) values(?,?)", "helloWorld", "123456789");

//update
gdao.ExecuteUpdate("update hstest set value=? where id=1", "hello");

//delete
gdao.ExecuteUpdate("delete from hstest where id = ?", 1);
```

### 5. 配置缓存

```go
//绑定Hstest  启用缓存, 缓存默认时效为300秒, gdaoCache为缓存操作入口
gdaoCache.BindClass[dao.Hstest]()
```

### 6. 读写分离

```go
mysqldb := getDataSource("mysql.json") // 获取备库数据源：mysqldb
gdaoSlave.BindClass[dao.Hstest](mysqldb, gdao.MYSQL) //这里主数据库为sqlite，备数据库为mysql，Hstest读取数据源为mysql， gdaoSlave为读写分离操作入口
```

### 7.  SQL映射

```xml
<!-- MyBatis 风格的 XML 配置文件 -->
<mapper namespace="user">
   <select id="selectHstest1" parameterType="int64" resultType="hstest1">
      SELECT * FROM hstest1  order by id desc limit #{limit}
   </select>
</mapper>
```

```go
//读取解析xml配置
hs1, _ := gdaoMapper.Select[dao.Hstest1]("user.selectHstest1", 1)
fmt.Println(hs1)
```

### 8.  SqlBuilder

```go
func Test_sqlBuilder_if(t *testing.T) {
   context := map[string]any{"id": 12}
   
   builder := &SqlBuilder{}
   builder.Append("SELECT * FROM hstest where 1=1").
           AppendIf("id>0", context, "and id=?", context["id"])  //动态SQL，当id>0时，动态添加 and id=?
   
   bean := builder.SelectOne()  //查询SQL: SELECT * FROM hstest where 1=1 and id=?  [ARGS][12]
   logger.Debug(bean)
}
```