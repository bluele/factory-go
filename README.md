# factory-go

[![wercker status](https://app.wercker.com/status/1331f3d73cd25c2b45f76475b40c6a9c/m/master "wercker status")](https://app.wercker.com/project/bykey/1331f3d73cd25c2b45f76475b40c6a9c)

factory-go is a is a fixtures replacement inspired by factory_boy and factory_girl.

It can be generated easily complex objects by using this, and maitain easily those objects generaters.

## Install

```
$ go get -u github.com/bluele/factory-go/factory
```

## Usage

All of the following code on [examples](https://github.com/bluele/factory-go/tree/master/examples).

### Define a simple factory

Declare an factory has a set of simple attribute, and generate a fixture object.

```go
package main

import (
  "fmt"
  "github.com/bluele/factory-go/factory"
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
```

Output:

```
ID: 1  Name: user-1  Location: Tokyo
ID: 2  Name: user-2  Location: Tokyo
ID: 3  Name: user-3  Location: Tokyo
```

### Use factory with random yet realistic values.

Tests look better with random yet realistic values. For example, you can use [go-randomdata](https://github.com/Pallinder/go-randomdata) library to get them:

```go
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

// 'Location: "Tokyo"' is default value.
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
```

Output:

```
ID: 1  Name: Benjamin Thomas  Location: Burrton
ID: 2  Name: Madison Davis  Location: Brandwell
ID: 3  Name: Aubrey Robinson  Location: Campden
```

### Define a factory includes sub-factory

```go
package main

import (
  "fmt"
  "github.com/bluele/factory-go/factory"
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
).SeqInt("ID", func(n int) (interface{}, error) {
  return 2 - n%2, nil
}).Attr("Name", func(args factory.Args) (interface{}, error) {
  group := args.Instance().(*Group)
  return fmt.Sprintf("group-%d", group.ID), nil
})

// 'Location: "Tokyo"' is default value.
var UserFactory = factory.NewFactory(
  &User{Location: "Tokyo"},
).SeqInt("ID", func(n int) (interface{}, error) {
  return n, nil
}).Attr("Name", func(args factory.Args) (interface{}, error) {
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
```

Output:

```
ID: 1  Name: user-1  Location: Tokyo  Group.ID: 1  Group.Name group-1
ID: 2  Name: user-2  Location: Tokyo  Group.ID: 2  Group.Name group-2
ID: 3  Name: user-3  Location: Tokyo  Group.ID: 1  Group.Name group-1
```

### Define a factory includes a slice for sub-factory.

```go
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
```

Output:

```
ID: 1  Name: user-1
        Post.ID: 1  Post.Content: post-1
        Post.ID: 2  Post.Content: post-2
        Post.ID: 3  Post.Content: post-3
ID: 2  Name: user-2
        Post.ID: 4  Post.Content: post-4
        Post.ID: 5  Post.Content: post-5
        Post.ID: 6  Post.Content: post-6
ID: 3  Name: user-3
        Post.ID: 7  Post.Content: post-7
        Post.ID: 8  Post.Content: post-8
        Post.ID: 9  Post.Content: post-9
```

# Author

**Jun Kimura**

* <http://github.com/bluele>
* <junkxdev@gmail.com>
