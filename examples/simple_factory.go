package main

import (
	"fmt"
	"github.com/bluele/factory-go"
)

type User struct {
	ID       int
	Name     string
	Location string
}

// 'Location: "Tokyo"' is default value.
var UserFactory = factory.NewFactory(
	&User{Location: "Tokyo"},
).SeqInt("ID", func(n int) (interface{}, error) {
	return n, nil
}).Attr("Name", func(args factory.Args) (interface{}, error) {
	user := args.Instance().(*User)
	return fmt.Sprintf("user-%d", user.ID), nil
})

func main() {
	for i := 0; i < 3; i++ {
		user := UserFactory.MustCreate().(*User)
		fmt.Println("ID:", user.ID, " Name:", user.Name, " Location:", user.Location)
	}
}
