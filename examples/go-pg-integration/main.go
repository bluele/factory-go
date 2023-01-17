package main

import (
	"context"
	"fmt"

	"github.com/bluele/factory-go/factory"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type Group struct {
	ID   int
	Name string
}

type User struct {
	ID    int
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
	tx := args.Context().Value("tx").(*pg.Tx)
	return tx.Insert(args.Instance())
}).SubFactory("Group", GroupFactory)

var GroupFactory = factory.NewFactory(
	&Group{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return n, nil
}).Attr("Name", func(args factory.Args) (interface{}, error) {
	group := args.Instance().(*Group)
	return fmt.Sprintf("group-%d", group.ID), nil
}).OnCreate(func(args factory.Args) error {
	tx := args.Context().Value("tx").(*pg.Tx)
	return tx.Insert(args.Instance())
})

func createTestSchema(db *pg.DB) error {
	tables := []interface{}{
		&Group{},
		&User{},
	}
	for _, table := range tables {
		err := db.DropTable(table, &orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			return err
		}

		err = db.CreateTable(table, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func openDB() *pg.DB {
	db := pg.Connect(&pg.Options{
		User: "postgres",
	})

	err := createTestSchema(db)
	if err != nil {
		panic(err)
	}

	return db
}

func main() {
	db := openDB()
	for i := 0; i < 3; i++ {
		tx, err := db.Begin()
		if err != nil {
			panic(err)
		}

		ctx := context.WithValue(context.Background(), "tx", tx)
		v, err := UserFactory.CreateWithContext(ctx)
		if err != nil {
			panic(err)
		}
		user := v.(*User)
		fmt.Println(user, *user.Group)
		tx.Commit()
	}
}
