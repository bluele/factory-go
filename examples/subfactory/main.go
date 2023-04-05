package main

import (
	"fmt"

	"github.com/hyuti/factory-go/factory"
)

type Group struct {
	ID   int
	Name string
}

type User struct {
	ID       int
	Name     string
	Location string
	Group    *Group
}

var GroupFactory = factory.NewFactory(
	&Group{},
).SeqInt("ID", func(n int) (any, error) {
	return 2 - n%2, nil
}).Attr("Name", func(args factory.Args) (any, error) {
	group := args.Instance().(*Group)
	return fmt.Sprintf("group-%d", group.ID), nil
})

// 'Location: "Tokyo"' is default value.
var UserFactory = factory.NewFactory(
	&User{Location: "Tokyo"},
).SeqInt("ID", func(n int) (any, error) {
	return n, nil
}).Attr("Name", func(args factory.Args) (any, error) {
	user := args.Instance().(*User)
	return fmt.Sprintf("user-%d", user.ID), nil
}).SubFactory("Group", GroupFactory)

func main() {
	for i := 0; i < 3; i++ {
		user := UserFactory.MustCreate().(*User)
		fmt.Println(
			"ID:", user.ID, " Name:", user.Name, " Location:", user.Location,
			" Group.ID:", user.Group.ID, " Group.Name", user.Group.Name)
	}
}
