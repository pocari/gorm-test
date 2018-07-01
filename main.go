package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Userモデル
type User struct {
	gorm.Model
	Name     string `gorm:"size:255"`
	Password string `gorm:"size:255"`
	Email    string `gorm:"size:255"`
}

func insertSampleRecords(db *gorm.DB) {
	users := []User{
		{
			Name:     "日本語名01",
			Password: "password01",
			Email:    "email01@example.com",
		},
		{
			Name:     "日本語名02",
			Password: "password02",
			Email:    "email02@example.com",
		},
		{
			Name:     "日本語名03",
			Password: "password03",
			Email:    "email03@example.com",
		},
	}

	for _, u := range users {
		var tmp User
		nf := db.Where(
			"name = ?", u.Name,
		).First(&tmp).RecordNotFound()
		if nf {
			// なかったら作る
			db.NewRecord(&u)
			db.Create(&u)
		}
	}
}

func connect(
	user string,
	pass string,
	host string,
	port int,
	database string,
) (*gorm.DB, error) {
	connectinString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		pass,
		host,
		port,
		database,
	)
	return gorm.Open("mysql", connectinString)
}

func main() {
	db, err := connect("user01", "user01", "127.0.0.1", 13306, "test_database")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// ログ出力を有効にする
	db.LogMode(true)

	insertSampleRecords(db)

	var u User
	db.Where("name = ?", "日本語名02").First(&u).RecordNotFound()
	fmt.Printf("user: %+v", u)
}
