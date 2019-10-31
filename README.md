## Start A Connection
### Base Configuration
```golang
var conf = gorm.NewConnConfig(
    gorm.WithUser("root"), // set user
    gorm.WithPassword(""), // set password
    gorm.WithDB("db_name"), // set database name
    gorm.WithAddress("127.0.0.1:3306"), // set address
)
```

### Other Options
- WithDriverName. `mysql` by default
- WithCharSet. `utf8mb4` by default
- WithParseTime. `true` by default
- WithLoc. `UTC` by default
- WithTimeout `1m` by default
- WithReadTimeout `1m` by default
- WithWriteTimeout `1m` by default
- WithDebug. `false` by default
- WithGetTableNameMethod. `Name` by default

### Connect To Database
```golang
var db, err = conf.Open()
if err != nil {
    ...
}
```

## Table Model
```golang
type User struct{
    ID int64 `gorm:"column:id"` // only struct field with tag "gorm" is a gorm table column
    Name string `gorm:"column:name"`
    Age int `gorm:"column:age"`
    IsSleeping bool `gorm:"column:is_sleeping;default:false"`
}

// table name is users by default

// you can custom table name by:
func (t User) Name() { return "student" }

// or if you:
var conf = gorm.NewConnConfig(
    ...,
    gorm.WithGetTableNameMethod("MyTableName"),
    ...,
)

// you can custom table name by:
func (t User) MyTableName() { return "student" }
```

## INSERT
```golang
var Jhon User
Jhon.ID = 1
Jhon.Name = "Jhon"
Jhon.Age = 17

err = db.Insert(John).Do()
if err!=nil {
    ...
}
```

## SELECT
```golang

// select * from user where name = 'Jhon';
var jhon User
err = db.Select().Where(User{Name: "Jhon"}).Do(&jhon)
if err!=nil {
    ...
}

// select * from user where `age` = 17;
var teenagers []*User
err = db.Select().Where(User{Age: 17}).Do(&teenagers)
if err!=nil {
    ...
}

// select `id`, `name` from user where `age` = 17;
err = db.Select("id", "name").Where(User{Age: 17}).Do(&teenagers)
if err!=nil {
    ...
}

// select * from user where `age` > 16 and id < 100;
err = db.Table(User{}).Select().Where(map[string]interface{}{"age": 16}, ">").And(map[string]interface{}{"id": 100}, "<").Do(&teenagers)
if err!=nil {
    ...
}

```

## UPDATE
```golang
// update user set age = 18 where name = 'Jhon';
err = db.Update(User{Age: 18}).Where(User{Name: "Jhon"}).Do()
if err!=nil {
    ...
}

// update user set is_sleeping = true where name = 'Jhon';
err = db.Update(map[string]interface{"is_sleeping": true}).Where(User{Name: "Jhon"}).Do()
if err!=nil {
    ...
}

// update user set is_sleeping = true;
err = db.Table(User{}).Update(map[string]interface{"is_sleeping": true}).Do()
if err!=nil {
    ...
}

// update user set is_sleeping = false where name = 'Jhon';
err = db.Update(map[string]interface{"is_sleeping": false}).Where(User{Name: "Jhon"}).Do()
if err!=nil {
    ...
}

// update user set aget = 19, is_sleeping = DEFAULT(is_sleeping) where name = 'Jhon';
err = db.Update(&User{Age: 19}, "is_sleeping").Where(User{Name: "Jhon"}).Do()
if err!=nil {
    ...
}
```

## Transaction

### Begin A Transaction
```golang
var tx = db.Begin()

// then use tx like using db
tx.Insert().Do()
tx.Update().Do()
tx.Select().Do()
```

### End A Transaction
```golang
// commit && rollback
tx.Commit()
tx.Rollback()

// transaction end
// if tx.err == nil, tx.Commmit(); else tx.Rollback()
tx.End(tx.err == nil)
```

## TODO:
- Complete select conditions
- create tables, drop tables, auto migrate tables, create db

