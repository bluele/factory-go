package main

import (
	"fmt"
	"github.com/bluele/factory-go/factory"
)

type User struct {
	ID    int
	Name  string
	Group *Group
}

type Group struct {
	ID    int
	Name  string
	Users []*User
}

var UserFactory = factory.NewFactory(
	&User{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return n, nil
}).Attr("Name", func(args factory.Args) (interface{}, error) {
	user := args.Instance().(*User)
	return fmt.Sprintf("user-%d", user.ID), nil
}).Attr("Group", func(args factory.Args) (interface{}, error) {
	if parent := args.Parent(); parent != nil {
		// if args have parent, use it.
		return parent.Instance(), nil
	}
	return nil, nil
})

var GroupFactory = factory.NewFactory(
	&Group{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return 2 - n%2, nil
}).Attr("Name", func(args factory.Args) (interface{}, error) {
	group := args.Instance().(*Group)
	return fmt.Sprintf("group-%d", group.ID), nil
}).SubSliceFactory("Users", UserFactory, func() int { return 3 })

func main() {
	group := GroupFactory.MustCreate().(*Group)
	fmt.Println("Group.ID:", group.ID)
	for _, user := range group.Users {
		fmt.Println("\tUser.ID:", user.ID, " User.Name:", user.Name, " User.Group.ID:", user.Group.ID)
	}
}
