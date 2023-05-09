# Rushia DB

Rushia DB 是依靠傳入現有的資料庫連線來將 Rushia 作為資料庫互動函式使用。

### 初始化語法

透過 `NewDB(...)` 來初始化 Rushia DB，你需要傳入一個有實作 `rushia.DB` 的結構體，舉例來說 Rushia DB 有 `rushiadb` 套件內建支援 [go-gorm/gorm](https://github.com/go-gorm/gorm)。

```go
// 初始化 Gorm 的連線。
db, err := gorm.Open(mysql.Open("root:password@tcp(localhost:3306)/db"), &gorm.Config{})

// 初始化一個基於現有 Gorm 連線的 Rushia DB。
q := rushia.NewDB(rushiadb.NewGorm(db))

// 開始使用 Rushia DB。
var user User
err = q.NewQuery("Users").Where("Username = ?", "YamiOdymel").Select().Query(&user)
// 等效於：SELECT * FROM Users WHERE Username = ?
```

### 執行與取得

Rushia DB 主要的用法是在 SQL 語法建置完的最後呼叫 `Exec()` 與 `Query(dest)`，取決於你正在使用哪種 SQL 語法。通常來說 `Update`、`Delete` 就是 `Exec`，因為他們不會回傳任何結果；而 `Select` 能使用 `Query` 來將結果映射到某個變數。

```go
type User struct {
	Username string
	Password string
}
u := User{
	Username: "YamiOdymel",
	Password: "test",
}
err = q.NewQuery("Users").Insert(u).Exec()
// 等效於：INSERT INTO Users (Username, Password) VALUES (?, ?)
```

### 事務交易

每個實作 `rushia.DB` 的資料庫連線都需要支援事務交易，當使用 `Transaction` 時，裡面回傳的 `error` 不是 `nil` 的時候，就會觸發相關的回溯機制，將剛才的所有變動全部復原。

```go
q.Transaction(func (tx *rushia.Query) error {
    return tx.NewQuery("Users").Insert(u).Exec()
})
```

你也可以手動透過 `Begin` 開始記錄並繼續你的資料庫寫入行為，如果途中發生錯誤，你能透過 `Rollback` 回到紀錄之前的狀態，即為回溯（或滾回、退回），如果這筆交易已經沒有問題了，透過 `Commit` 將這次的變更永久地儲存到資料庫中。

```go
// 當交易開始時請使用回傳的 `tx` 而不是原先的 `q`，這樣才能確保交易繼續。
tx, err := q.Begin()
if err != nil {
	panic(err)
}

// 如果插入資料時發生錯誤，則呼叫 `Rollback()` 回到交易剛開始的時候。
if _, err = tx.NewQuery("Wallets").Insert(data).Exec(); err != nil {
	tx.Rollback()
	panic(err)
}
if _, err = tx.NewQuery("Users").Insert(data).Exec(); err != nil {
	tx.Rollback()
	panic(err)
}

// 透過 `Commit()` 確保上列變更都已經永久地儲存到資料庫。
if err := tx.Commit(); err != nil {
	panic(err)
}
```
