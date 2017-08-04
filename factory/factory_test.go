package factory

import (
	"testing"
)

func TestFactory(t *testing.T) {
	type User struct {
		ID         int
		Name       string
		Location   string
		unexported string
	}

	userFactory := NewFactory(&User{Location: "Tokyo"}).
		SeqInt("ID", func(n int) (interface{}, error) {
			return n, nil
		}).
		Attr("Name", func(args Args) (interface{}, error) {
			return "bluele", nil
		})

	iuser, err := userFactory.Create()
	if err != nil {
		t.Error(err)
	}
	user, ok := iuser.(*User)
	if !ok {
		t.Error("It should be *User type.")
		return
	}
	if user.ID != 1 {
		t.Error("user.ID should be 1.")
		return
	}
	if user.Name != "bluele" {
		t.Error(`user.Name should be "bluele".`)
		return
	}
	if user.Location != "Tokyo" {
		t.Error(`user.Location should be "Tokyo".`)
		return
	}
}

func TestMapAttrFactory(t *testing.T) {
	type User struct {
		ID  int
		Ext map[string]string
	}
	var userFactory = NewFactory(&User{}).
		SeqInt("ID", func(n int) (interface{}, error) {
			return n, nil
		}).
		Attr("Ext", func(args Args) (interface{}, error) {
			return map[string]string{"test": "ok"}, nil
		})
	user := &User{}
	if err := userFactory.Construct(user); err != nil {
		t.Error(err)
		return
	}
	if user.ID == 0 {
		t.Error("user.ID should not be 0.")
	}
	if v, ok := user.Ext["test"]; !ok {
		t.Error("user.Ext[\"test\"] should not be empty.")
	} else if v != "ok" {
		t.Error("user.Ext[\"test\"] should be ok.")
	}
}

func TestSubFactory(t *testing.T) {
	type Group struct {
		ID int
	}
	type User struct {
		ID    int
		Name  string
		Group *Group
	}

	groupFactory := NewFactory(&Group{}).
		SeqInt("ID", func(n int) (interface{}, error) {
			return n, nil
		})

	userFactory := NewFactory(&User{}).
		SeqInt("ID", func(n int) (interface{}, error) {
			return n, nil
		}).
		Attr("Name", func(args Args) (interface{}, error) {
			return "bluele", nil
		}).
		SubFactory("Group", groupFactory)

	iuser, err := userFactory.Create()
	if err != nil {
		t.Error(err)
	}
	user, ok := iuser.(*User)
	if !ok {
		t.Error("It should be *User type.")
		return
	}

	if user.Group == nil {
		t.Error("user.Group should be *Group type.")
		return
	}

	if user.Group.ID != 1 {
		t.Error("user.Group.ID should be 1.")
		return
	}
}

func TestSubSliceFactory(t *testing.T) {
	type Group struct {
		ID int
	}
	type User struct {
		ID     int
		Name   string
		Groups []*Group
	}

	groupFactory := NewFactory(&Group{}).
		SeqInt("ID", func(n int) (interface{}, error) {
			return n, nil
		})

	userFactory := NewFactory(&User{}).
		SeqInt("ID", func(n int) (interface{}, error) {
			return n, nil
		}).
		Attr("Name", func(args Args) (interface{}, error) {
			return "bluele", nil
		}).
		SubSliceFactory("Groups", groupFactory, func() int { return 3 })

	iuser, err := userFactory.Create()
	if err != nil {
		t.Error(err)
	}
	user, ok := iuser.(*User)
	if !ok {
		t.Error("It should be *User type.")
		return
	}

	if user.Groups == nil {
		t.Error("user.Groups should be []*Group type.")
		return
	}

	if len(user.Groups) != 3 {
		t.Error("len(user.Groups) should be 3.")
		return
	}

	for i := 0; i < 3; i++ {
		if user.Groups[i].ID != i+1 {
			t.Errorf("user.Groups[%v].ID should be %v", i, i+1)
			return
		}
	}
}

func TestSubRecursiveFactory(t *testing.T) {
	type User struct {
		ID     int
		Name   string
		Friend *User
	}

	var userFactory = NewFactory(&User{})
	userFactory.
		SeqInt("ID", func(n int) (interface{}, error) {
			return n, nil
		}).
		Attr("Name", func(args Args) (interface{}, error) {
			return "bluele", nil
		}).
		SubRecursiveFactory("Friend", userFactory, func() int { return 2 })

	iuser, err := userFactory.Create()
	if err != nil {
		t.Error(err)
		return
	}
	user, ok := iuser.(*User)
	if !ok {
		t.Error("It should be *User type.")
		return
	}

	if user.Friend.Friend == nil {
		t.Error("user.Friend.Friend should not be nil.")
		return
	}

	if user.Friend.Friend.Friend != nil {
		t.Error("user.Friend.Friend.Friend should be nil.")
		return
	}
}

func TestFactoryConstruction(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	var userFactory = NewFactory(&User{}).
		SeqInt("ID", func(n int) (interface{}, error) {
			return n, nil
		}).
		Attr("Name", func(args Args) (interface{}, error) {
			return "bluele", nil
		})

	var user *User

	user = &User{}
	if err := userFactory.Construct(user); err != nil {
		t.Error(err)
		return
	}
	if user.ID == 0 {
		t.Error("user.ID should not be 0.")
	}
	if user.Name == "" {
		t.Error("user.ID should not be empty.")
	}

	user = &User{}
	if err := userFactory.ConstructWithOption(user, map[string]interface{}{"Name": "jun"}); err != nil {
		t.Error(err)
		return
	}
	if user.ID == 0 {
		t.Error("user.ID should not be 0.")
	}
	if user.Name == "" {
		t.Error("user.ID should not be empty.")
	}
}

func TestFactoryWhenCallArgsParent(t *testing.T) {
	type User struct {
		Name      string
		GroupUUID string
	}

	var userFactory = NewFactory(&User{})
	userFactory.
		Attr("Name", func(args Args) (interface{}, error) {
			if parent := args.Parent(); parent != nil {
				if pUser, ok := parent.Instance().(*User); ok {
					return pUser.GroupUUID, nil
				}
			}
			return "", nil
		})

	if err := userFactory.Construct(&User{}); err != nil {
		t.Error(err)
		return
	}
}
