package main

import (
	"fmt"
	"github.com/bluele/factory-go/factory"
)

type Post struct {
	ID      int
	Content string
}

type User struct {
	ID    int
	Name  string
	Posts []*Post
}

var PostFactory = factory.NewFactory(
	&Post{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return n, nil
}).Attr("Content", func(args factory.Args) (interface{}, error) {
	post := args.Instance().(*Post)
	return fmt.Sprintf("post-%d", post.ID), nil
})

var UserFactory = factory.NewFactory(
	&User{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return n, nil
}).Attr("Name", func(args factory.Args) (interface{}, error) {
	user := args.Instance().(*User)
	return fmt.Sprintf("user-%d", user.ID), nil
}).SubSliceFactory("Posts", PostFactory, func() int { return 3 })

func main() {
	for i := 0; i < 3; i++ {
		user := UserFactory.MustCreate().(*User)
		fmt.Println("ID:", user.ID, " Name:", user.Name)
		for _, post := range user.Posts {
			fmt.Printf("\tPost.ID: %v  Post.Content: %v\n", post.ID, post.Content)
		}
	}
}
