package main

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/bluele/factory-go/factory"
)

type User struct {
	ID       int
	Name     string
	Location string
}

var UserFactory = factory.NewFactory(
	&User{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return n, nil
}).Attr("Name", func(args factory.Args) (interface{}, error) {
	return randomdata.FullName(randomdata.RandomGender), nil
}).Attr("Location", func(args factory.Args) (interface{}, error) {
	return randomdata.City(), nil
})

func main() {
	for i := 0; i < 3; i++ {
		user := UserFactory.MustCreate().(*User)
		fmt.Println("ID:", user.ID, " Name:", user.Name, " Location:", user.Location)
	}
}
