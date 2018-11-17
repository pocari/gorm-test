package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// User モデル
type User struct {
	gorm.Model
	Name         string `gorm:"size:255"`
	Password     string `gorm:"size:255"`
	Email        string `gorm:"size:255"`
	RegisterDate time.Time
}

func insertSampleRecords(db *gorm.DB) {
	now := time.Now()
	fmt.Printf("Now: %v\n", now)
	users := []User{
		{
			Name:         "日本語名01",
			Password:     "password01",
			Email:        "email01@example.com",
			RegisterDate: now,
		},
		{
			Name:         "日本語名02",
			Password:     "password02",
			Email:        "email02@example.com",
			RegisterDate: now,
		},
		{
			Name:         "日本語名03",
			Password:     "password03-updated",
			Email:        "email03@example.com",
			RegisterDate: now,
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
		} else {
			// あったら更新
			db.Model(&tmp).Update(&u)
		}
	}
}

// ConnectionInfo struct
type ConnectionInfo struct {
	user         string
	pass         string
	host         string
	port         int
	databaseName string
}

func connect(ci *ConnectionInfo) (*gorm.DB, error) {
	connectinString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		ci.user,
		ci.pass,
		ci.host,
		ci.port,
		ci.databaseName,
	)
	return gorm.Open("mysql", connectinString)
}

func main() {
	ci := &ConnectionInfo{
		user:         "user01",
		pass:         "user01",
		host:         "127.0.0.1",
		port:         13306,
		databaseName: "test_database",
	}
	db, err := connect(ci)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// ログ出力を有効にする
	db.LogMode(true)

	if os.Getenv("RUN_MIGRATION") != "" {
		db.AutoMigrate(&User{})
	}
	insertSampleRecords(db)

	var u User
	if !db.Where("name = ?", "日本語名02").First(&u).RecordNotFound() {
		fmt.Printf("user: %+v\n", u)
		fmt.Printf("utc: %v\n", u.RegisterDate.UTC())
		fmt.Printf("local: %v\n", u.RegisterDate.Local())
	} else {
		fmt.Println("user not found")
	}
}
