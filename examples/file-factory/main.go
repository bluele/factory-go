package main

import (
	"fmt"
	"os"

	"github.com/hyuti/factory-go/factory"
)

type User struct {
	Img string
}

var UserFactory = factory.NewFactory(
	&User{},
).Png("Img", func(f *os.File) (any, error) {
	return f, nil
})

func main() {
	for i := 0; i < 3; i++ {
		user := UserFactory.MustCreate().(*User)
		fmt.Println("img:", user.Img)
	}
	defer UserFactory.Clean()
}
