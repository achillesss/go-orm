## Start A Connection
### Base Configuration
#### Options on connection
Options below will be used while connecting to database.
- WithDriverName. `mysql` by default
- WithUser.
- WithPassword.
- WithAddress.
- WithDB.
- WithCharSet. `utf8mb4` by default
- WithParseTime. `true` by default
- WithLoc. `UTC` by default
- WithTimeout `1m` by default
- WithReadTimeout `1m` by default
- WithWriteTimeout `1m` by default

#### Options to control 
- WithGetTableNameMethod. `Name` by default. for details, see comment on it.
```golang
orm.WithGetTableNameMethod("TableName")
```
- WithHandleError. Register a function to handle the error that returned by any query.
```golang
orm.WithHandleError(func(err error){
    fmt.Printf("Execution meets an error: %v\n", err)
})
```
- WithHandleCommitError. Register a function to handle the error that returned by sql transaction commit.
```golang
orm.WithHandleCommitError(func(err error) {
    fmt.Printf("Transaction commit meets an error: %v\n", err)
})
```
- WithDBStatsMonitor. Register a monitor which monitoring db status.
```golang
orm.WithDBStatsMonitor(func(f func() orm.DBStatus) {
    var ticker = time.NewTicker(time.Minute)
    for range ticker.C {
        status = f()
        fmt.Printf("DB Status: %+v\n", status)
    }
})
```
- WithStartQueryMonitor. Register a monitor which monitoring all the querys before executed.
```golang
orm.WithStartQueryMonitor(func(queries <-chan *orm.StartQuery) {
    for query range queries {
        fmt.Printf("query start: %+v\n", query)
    }
})
```
- WithEndQueryMonitor. Register a monitor which monitoring all the querys after executed.
```golang
orm.WithEndQueryMonitor(func(queries <-chan *orm.EndQuery) {
    for query range queries {
        fmt.Printf("query end: %+v\n", query)
    }
})
```
- WithBeginTxMonitor. Register a monitor which monitoring all the sql transactions on beginning.
```golang
orm.WithBeginTxMonitor(func(txs <-chan *orm.BeginTx) {
    for tx range txs {
        fmt.Printf("tx start: %+v\n", tx)
    }
})
```
- WithEndTxMonitor. Register a monitor which monitoring all the sql transactions on commit or rollback.
```golang
orm.WithEndTxMonitor(func(txs <-chan *orm.EndTx) {
    for tx range txs {
        fmt.Printf("tx end: %+v\n", tx)
    }
})
```

### Connect to database
```golang
// create config with options
var conf = orm.NewConnConfig(
    orm.WithUser("root"), // set user
    orm.WithPassword(""), // set password
    orm.WithDB("db_name"), // set database name
    orm.WithAddress("127.0.0.1:3306"), // set address
)

// connect to database
var db, err = conf.Open()
if err != nil {
    ...
}
```

## Table
### A table is a struct
```golang
type User struct{
    ID int64 `gorm:"column:id;auto_increment"` // only struct field with tag "gorm" is a table column
    Name string `gorm:"column:name"`
    Age int `gorm:"column:age"`
    IsSleeping bool `gorm:"column:is_sleeping;default:false"`
    Deskmate *int64 `gorm:"desk_mate"`
}
```
### Table name

#### default table name
```golang
// table `User` name is `user` by default
```

#### custom table name by default method
```golang
// you can custom it's name by creating default method on it:
func (t User) Name() { return "student" }
```
#### custom table name by specified method on it:
```golang
// you can custom table name by:
func (t User) MyTableName() { return "student" }

// and you must connect to database by using `WithGetTableNameMethod` option:
var conf = orm.NewConnConfig(
    ...,
    orm.WithGetTableNameMethod("MyTableName"),
    ...,
)
```

## INSERT
```golang
var Jhon User
Jhon.Name = "Jhon"
Jhon.Age = 17

// Insert will write the `ID` value back to table
// only if the table has a `id` column.
// ORM treats `id` column as primary key.
db.Insert(&John).Do()

// Batch Insert
var tony User
var sara User
db.Insert([]*User{&tony, &sara}).Do()
```

## SELECT
```golang
// select * from user;
var all []*User
db.Table(User{}).Select().Do(&all)

// select * from user where name = 'Jhon';
var jhon User
db.Select().Where(User{Name: "Jhon"}).Do(&jhon)

// select * from user where `age` = 17;
var teenagers []*User
db.Select().Where(User{Age: 17}).Do(&teenagers)

// select * from user where `age` = 17 and is_sleep = false;
// 1.
var teenagers []*User
db.Select().Where(User{Age: 17}).And(map[string]interface{}{"is_sleep": false}).Do(&teenagers)

// 2. is_sleep is false by default, so we can
db.Select().Where(User{Age: 17}, "is_sleep").Do(&teenagers)

// select `id`, `name` from user where `age` = 17;
db.Select("id", "name").Where(User{Age: 17}).Do(&teenagers)

// select * from user where `age` > 16 and id < 100;
db.Table(User{}).Select().Where(map[string]interface{}{"age": 16}, ">").And(map[string]interface{}{"id": 100}, "<").Do(&teenagers)

// select * from user where age in(10, 20);
db.Table(User{}).Select().Where(map[string][]interface{}{"age": []interface{}{10, 20}}).Do(&teenagers)

// select * from user where age between 10 and 20;
db.Table(User{}).Select().Where(map[string][]interface{}{"age": []interface{}{10, 20}}, "BETWEEN").Do(&teenagers)

// select count(*) from user where desk_mate is NULL;
var count int
db.Talble(User{}).Select("count(*)").Where(map[string]interface{}{"desk_mate": nil}).Do(&count)
```

## UPDATE
```golang
// update user set age = 18 where name = 'Jhon';
db.Update(User{Age: 18}).Where(User{Name: "Jhon"}).Do()

// update user set is_sleeping = true where name = 'Jhon';
db.Update(map[string]interface{"is_sleeping": true}).Where(User{Name: "Jhon"}).Do()

// update user set is_sleeping = true;
db.Table(User{}).Update(map[string]interface{"is_sleeping": true}).Do()

// update user set is_sleeping = false where name = 'Jhon';
db.Update(map[string]interface{"is_sleeping": false}).Where(User{Name: "Jhon"}).Do()

// update user set aget = 19, is_sleeping = DEFAULT(is_sleeping) where name = 'Jhon';
db.Update(&User{Age: 19}, "is_sleeping").Where(User{Name: "Jhon"}).Do()
```

## DELETE
```golang
// delete from user where name = 'Jhon';
db.Delete().Where(User{Name: "Jhon"}).Do()
```

## RAW
```golang
// alter table use drop column desk_mate;
db.Raw("alter table %s drop column %s", User{}.Name(), "desk_mate").Do()
```

## Transaction

### Begin 
```golang
var tx = db.Begin()

// then use tx like using db
tx.Insert().Do()
tx.Update().Do()
tx.Select().Do()
```

### End 
```golang
// commit && rollback
tx.Commit()
tx.Rollback()

// transaction end
// if tx.err == nil, tx.Commmit(); else tx.Rollback()
// this is recommended
tx.End(tx.err == nil)

func doTransaction(){
    var tx = db.Begin()
    var err error
    defer func(){
        tx.End(err == nil)
    }()

    // just do it
    ...
}
```

## Callback 
you can put a callback function like `func()` into Do() like:
```golang
db.Table(User{}).Select().Do(func(){
    fmt.Printf("Found all users.\n")
}, &all)

db.Update(User{Age: 18}).Where(User{Name: "Jhon"}).Do(func(){
    fmt.Printf("Jhon's age updated.\n")
})

db.Insert([]*User{&tony, &sara}).Do(func(){
    fmt.Printf("Welcome tow freshmen.\n")
})
```
callback function will be called AFTER THE QUERY IS EXECUTED.

## TODO:
- Complete select conditions
- create tables, drop tables, auto migrate tables, create db

