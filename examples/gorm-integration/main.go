package main

import (
	"context"
	"fmt"

	"github.com/bluele/factory-go/factory"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Group struct {
	ID   int `gorm:"primary_key"`
	Name string
}

type User struct {
	ID    int `gorm:"primary_key"`
	Name  string
	Group *Group
}

var UserFactory = factory.NewFactory(
	&User{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return n, nil
}).Attr("Name", func(args factory.Args) (interface{}, error) {
	user := args.Instance().(*User)
	return fmt.Sprintf("user-%d", user.ID), nil
}).OnCreate(func(args factory.Args) error {
	db := args.Context().Value("db").(*gorm.DB)
	return db.Create(args.Instance()).Error
}).SubFactory("Group", GroupFactory)

var GroupFactory = factory.NewFactory(
	&Group{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return n, nil
}).Attr("Name", func(args factory.Args) (interface{}, error) {
	group := args.Instance().(*Group)
	return fmt.Sprintf("group-%d", group.ID), nil
}).OnCreate(func(args factory.Args) error {
	db := args.Context().Value("db").(*gorm.DB)
	return db.Create(args.Instance()).Error
})

func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	db.LogMode(true)
	db.AutoMigrate(&Group{}, &User{})

	for i := 0; i < 3; i++ {
		tx := db.Begin()
		ctx := context.WithValue(context.Background(), "db", tx)
		v, err := UserFactory.CreateWithContext(ctx)
		if err != nil {
			panic(err)
		}
		user := v.(*User)
		fmt.Println(user, *user.Group)
		tx.Commit()
	}
}
