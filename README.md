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
// select Jhon
var jhon User
err = db.Select().Where(User{Name: "Jhon"}).Do(&jhon)
if err!=nil {
    ...
}

// select 17 year old teenagers
var teenagers []*User
err = db.Select().Where(User{Age: 17}).Do(&teenagers)
if err!=nil {
    ...
}

// select 17 year old teenagers's `id` and `name`
err = db.Select("id", "name").Where(User{Age: 17}).Do(&sevenTeenUsers)
if err!=nil {
    ...
}
```

## UPDATE
```golang
// on jhon's birthday
err = db.Update(User{Age: 18}).Where(User{Name: "Jhon"}).Do()
if err!=nil {
    ...
}

// he goes to sleep
err = db.Update(map[string]interface{"is_sleeping": true}).Where(User{Name: "Jhon"}).Do()
if err!=nil {
    ...
}

// then everyone goes to sleep
err = db.Table(User{}).Update(map[string]interface{"is_sleeping": true}).Do()
if err!=nil {
    ...
}

// after a while, jhon wakes up
err = db.Update(map[string]interface{"is_sleeping": false}).Where(User{Name: "Jhon"}).Do()
if err!=nil {
    ...
}

// or jhon wakes update and find out that he's age increased
err = db.Update(&User{Age: 19}, "is_sleeping").Where(User{Name: "Jhon"}).Do()
if err!=nil {
    ...
}
```

## TODO:
- Transactions
- Complete select conditions
- create tables, drop tables, auto migrate tables, create db

