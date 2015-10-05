package main

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/bluele/factory-go/factory"
)

type User struct {
	ID          int
	Name        string
	CloseFriend *User
}

var UserFactory = factory.NewFactory(
	&User{},
)

func init() {
	UserFactory.SeqInt("ID", func(n int) (interface{}, error) {
		return n, nil
	}).Attr("Name", func(args factory.Args) (interface{}, error) {
		return randomdata.FullName(randomdata.RandomGender), nil
	}).SubRecursiveFactory("CloseFriend", UserFactory, func() int { return 2 }) // recursive depth is always 2
}

func main() {
	user := UserFactory.MustCreate().(*User)
	fmt.Println("ID:", user.ID, " Name:", user.Name,
		" CloseFriend.ID:", user.CloseFriend.ID, " CloseFriend.Name:", user.CloseFriend.Name)
	// `user.CloseFriend.CloseFriend.CloseFriend ` depth is 3, so this value is always nil.
	fmt.Printf("%v %v\n", user.CloseFriend.CloseFriend, user.CloseFriend.CloseFriend.CloseFriend)
}
