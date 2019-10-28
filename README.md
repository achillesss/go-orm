### TODO:
- Transactions
- Complete select conditions

### examples
```
var conf = gorm.NewConnConfig(
    gorm.WithUser("root"),
    gorm.WithPassword(""),
    gorm.WithDB("db_name"),
    gorm.WithAddress("127.0.0.1:3306"),
)

var db, err = conf.Open()

type User struct{
    ID int64 `gorm:"column:id"` // only struct field with tag "gorm" is a gorm table column
    Name string `gorm:"column:name"`
    Age int `gorm:"column:age"`
    IsSleeping bool `gorm:"column:is_sleeping;default:false"`
}

// insert
var Jhon User
Jhon.ID = 1
Jhon.Name = "Jhon"
Jhon.Age = 17

err = db.Insert(John).Do()
if err!=nil {
    ...
}

// select
// select one
var jhon User
err = db.Select().Where(User{Name: "Jhon"}).Do(&jhon)
if err!=nil {
    ...
}

// select many
var sevenTeenUsers []*User
err = db.Select().Where(User{Age: 17}).Do(&sevenTeenUsers)
if err!=nil {
    ...
}

// select columns
err = db.Select("id", "name").Where(User{Age: 17}).Do(&sevenTeenUsers)
if err!=nil {
    ...
}


// update

// on jhon's birthday
err = db.Update(User{Age: 18}).Where(User{Name: "Jhon"}).Do()
if err!=nil {
    ...
}

// jhon goes to sleep
err = db.Update(map[string]interface{"is_sleeping": true}).Where(User{Name: "Jhon"}).Do()
if err!=nil {
    ...
}

// jhon just wakes up
err = db.Update(map[string]interface{"is_sleeping": false}).Where(User{Name: "Jhon"}).Do()
if err!=nil {
    ...
}

// jhon wakes update and find out that he's age increased
err = db.Update(&User{Age: 19}, "is_sleeping").Where(User{Name: "Jhon"}).Do()
if err!=nil {
    ...
}

```
