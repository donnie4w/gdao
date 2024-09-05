## GDAO —— Go Persistence Framework  [[中文](https://github.com/donnie4w/gdao/blob/master/README_ZH.md)]

### Introduction

GDAO is a comprehensive persistence layer solution. Its primary purpose is to reduce coding effort, enhance productivity, improve performance, support the integration of multiple data sources, facilitate read-write separation, and establish a standard for persistence layer programming. By leveraging GDAO effectively, it's possible to reduce the amount of coding required in the persistence layer by 30%, even up to 50%, while simultaneously establishing a unified programming standard for the persistence layer, reducing errors, and making maintenance and expansion easier.
For Go, GDAO serves a similar role as [Hibernate](https://github.com/hibernate/hibernate-orm) + [MyBatis](https://github.com/mybatis/mybatis-3) do for Java. The GDAO framework combines the abstraction of Hibernate with the flexibility of MyBatis, addressing long-standing pain points in ORM frameworks. For more information about the pain points of Hibernate and MyBatis, refer to the [JDAO usage documentation](https://tlnet.top/jdaodoc).

GDAO fully implements the core features of MyBatis in Go, achieving the separation of SQL from the application code, and implementing powerful dynamic SQL functionality. It is currently the only ORM in Go that completely separates SQL from the application.

### [Usage Documentation](https://tlnet.top/gdaodoc)

### [Example Program](https://github.com/donnie4w/gdaodemo)

---

### Key Features
1. Full Table Object Mapping: Supports mapping table operations to standard struct object operations.
2. Full SQL Object Mapping: Supports complete SQL operation mapping to object operations.
3. Separation of SQL and Code: Supports the separation of SQL from the program, akin to [MyBatis](https://github.com/mybatis/mybatis-3).

### Main Functionality

1. Code Generation: Run the GDAO code generation tool to create standardized entity classes for database tables. Similar to thrift/protobuf.
2. Efficient Serialization: Standardized entity classes for tables implement efficient serialization and deserialization. Higher performance, smaller data size.
3. Read-Write Splitting: GDAO supports binding multiple slave data source connections and implementing read-write separation for bound tables or SQL queries.
4. Data Caching: GDAO supports data caching and provides detailed control over cache validity and data eviction characteristics.
5. Broad Compatibility: GDAO theoretically supports all databases that implement the Go database driver interface.
6. Advanced Features: Supports transactions, stored procedures, batch processing, and other database operations.
7. Dynamic SQL: GDAO implements rich dynamic SQL construction features. Supports dynamic SQL tag mapping construction and native SQL dynamic construction.
8. MyBatis Characteristics: GDAO's mapping module is the implementation of MyBatis in Go. It is also the only known ORM in Go that supports the separation of SQL from the application.

### GDAO Persistence Layer Solution

**Note:** Full table mapping and full SQL mapping can potentially lead to some issues such as over-encapsulation of objects and large numbers of repetitive, similar SQL statements. GDAO uses the abstraction of table mapping combined with the semi-automated management of SQL mapping to solve problems in the persistence layer.

1. **Full Table Object Mapping for CRUD Operations on Single Tables**: Mapping all operations to table object operations is a classic approach used by `Hibernate`, which can lead to over-encapsulation problems.
   GDAO supports auto-generating mapping entity classes for tables specifically to handle CRUD operations on single tables. Its underlying dynamic SQL construction effectively addresses the issue of large numbers of similar SQL statements resulting from similar table operation behaviors.
2. **Full SQL Object Mapping for Complex SQL Operations**: In practice, complex SQL, especially multi-table joins, often requires optimization, which requires understanding of the database schema, index properties, etc.
   Concatenating complex SQL strings using Go objects can increase the difficulty of comprehension. Eventually, developers may not know what the final executed SQL statement looks like after object concatenation, which increases risk and maintenance difficulty.
   GDAO implements the core feature of SQL mapping from `MyBatis`, allowing SQL to be configured in XML and mapped to Go objects.
3. **`SqlBuilder` for Native SQL Dynamic Construction**: GDAO implements dynamic construction of native SQL. For extremely complex SQL, including functions or stored procedures, GDAO supports dynamic construction using `SqlBuilder`.
   It is the procedural implementation of dynamic tags, based on Go programs to construct any form of SQL.

---

# Quick Start

### 1. Installation

```bash
go get github.com/donnie4w/gdao
```

### 2. Configuring Data Source

```go
gdao.Init(mysqldb, gdao.MYSQL) // Initialize gdao data source; mysqldb is an sql.DB object, gdao.MYSQL is the database type
```

### 3. Standard Entity Class Operations

```go
// Read
hs := dao.NewHstest()
hs.Where(hs.Id.EQ(10))
h, err := hs.Select(hs.Id, hs.Value, hs.Rowname)
//[DEBUG][SELETE ONE][ select id,value,rowname from hstest where id=?][10] 
```

```go
// Update
hs := dao.NewHstest()
hs.SetRowname("hello10")
hs.Where(hs.Id.EQ(10))
hs.Update()
//[DEBUG][UPDATE][update hstest set rowname=? where id=?][hello10 10]
```

```go
// Delete
hs := dao.NewHstest()
hs.Where(hs.Id.EQ(10))
hs.Delete()
//[DEBUG][UPDATE][delete from hstest where id=?][10]
```

```go
// Insert
hs := dao.NewHstest()
hs.SetValue("hello123")
hs.SetLevel(12345)
hs.SetBody([]byte("hello"))
hs.SetRowname("hello1234")
hs.SetUpdatetime(time.Now())
hs.Insert()
//[DEBUG][INSERT][insert  into hstest(floa,age,value,level,body,rowname,updatetime )values(?,?,?,?,?)][hello123 12345 hello hello1234 2024-07-17 19:36:44]
```

### 4. Native SQL Operations

```go
// Query, return a single row; gdao is the entry point for native SQL operations
gdao.ExecuteQueryBean("select id,value,rowname from hstest where id=?", 10)

// Insert
gdao.ExecuteUpdate("insert into hstest2(rowname,value) values(?,?)", "helloWorld", "123456789");

// Update
gdao.ExecuteUpdate("update hstest set value=? where id=1", "hello");

// Delete
gdao.ExecuteUpdate("delete from hstest where id = ?", 1);
```

### 5. Configuring Cache

```go
// Bind Hstest and enable caching; default cache expiration is 300 seconds; gdaoCache is the entry point for cache operations
gdaoCache.BindClass[dao.Hstest]()
```

### 6. Read-Write Splitting

```go
mysqldb := getDataSource("mysql.json") // Get slave data source: mysqldb
gdaoSlave.BindClass[dao.Hstest](mysqldb, gdao.MYSQL) // Here the master database is sqlite, the slave database is mysql; Hstest reads from the mysql data source; gdaoSlave is the entry point for read-write splitting operations
```

### 7. SQL Mapping

```xml
<!-- MyBatis style XML configuration file -->
<mapper namespace="user">
   <select id="selectHstest1" parameterType="int64" resultType="hstest1">
      SELECT * FROM hstest1  order by id desc limit #{limit}
   </select>
</mapper>
```

```go
// Read and parse XML configuration
hs1, _ := gdaoMapper.Select[dao.Hstest1]("user.selectHstest1", 1)
fmt.Println(hs1)
```

### 8. SqlBuilder

```go
func Test_sqlBuilder_if(t *testing.T) {
   context := map[string]any{"id": 12}
   
   builder := &SqlBuilder{}
   builder.Append("SELECT * FROM hstest where 1=1").
           AppendIf("id>0", context, "and id=?", context["id"])  // Dynamic SQL, when id>0, dynamically add and id=?
   
   bean := builder.SelectOne()  // Query SQL: SELECT * FROM hstest where 1=1 and id=?  [ARGS][12]
   logger.Debug(bean)
}
```