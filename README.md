# factory-go

![Test](https://github.com/bluele/factory-go/workflows/Test/badge.svg)
[![GoDoc](https://godoc.org/github.com/bluele/factory-go?status.svg)](https://pkg.go.dev/github.com/bluele/factory-go?tab=doc)

factory-go is a is a fixtures replacement inspired by [factory_boy](https://github.com/rbarrois/factory_boy) and [factory_girl](https://github.com/thoughtbot/factory_girl).

It can be generated easily complex objects by using this, and maintain easily those objects generaters.

## Install

```
$ go get -u github.com/bluele/factory-go/factory
```

## Usage

All of the following code on [examples](https://github.com/bluele/factory-go/tree/master/examples).

* [Define a simple factory](https://github.com/bluele/factory-go#define-a-simple-factory)
* [Use factory with random yet realistic values](https://github.com/bluele/factory-go#use-factory-with-random-yet-realistic-values)
* [Define a factory includes sub-factory](https://github.com/bluele/factory-go#define-a-factory-includes-sub-factory)
* [Define a factory includes a slice for sub-factory](https://github.com/bluele/factory-go#define-a-factory-includes-a-slice-for-sub-factory)
* [Define a factory includes sub-factory that contains self-reference](https://github.com/bluele/factory-go#define-a-factory-includes-sub-factory-that-contains-self-reference)
* [Define a sub-factory refers to parent factory](https://github.com/bluele/factory-go#define-a-sub-factory-refers-to-parent-factory)

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

### Define a factory includes sub-factory that contains self-reference.

```go
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
```

Output:
```
ID: 1  Name: Mia Williams  CloseFriend.ID: 2  CloseFriend.Name: Joseph Wilson
&{3 Liam Wilson <nil>} <nil>
```

### Define a sub-factory refers to parent factory

```go
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
```

Output:
```
Group.ID: 1
        User.ID: 1  User.Name: user-1  User.Group.ID: 1
        User.ID: 2  User.Name: user-2  User.Group.ID: 1
        User.ID: 3  User.Name: user-3  User.Group.ID: 1
```

## Persistent models

Currently this project has no support for directly integration with ORM like [gorm](https://github.com/jinzhu/gorm), so you need to do manually.

Here is an example: https://github.com/bluele/factory-go/blob/master/examples/gorm_integration.go

# Author

**Jun Kimura**

* <http://github.com/bluele>
* <junkxdev@gmail.com>
