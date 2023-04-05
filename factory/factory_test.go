package factory

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {
	type User struct {
		ID       int
		Name     string
		Location string
	}

	userFactory := NewFactory(&User{Location: "Tokyo"}).
		SeqInt("ID", func(n int) (any, error) {
			return n, nil
		}).
		Attr("Name", func(args Args) (any, error) {
			return "hyuti", nil
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
	if user.Name != "hyuti" {
		t.Error(`user.Name should be "hyuti".`)
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
		SeqInt("ID", func(n int) (any, error) {
			return n, nil
		}).
		Attr("Ext", func(args Args) (any, error) {
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
		SeqInt("ID", func(n int) (any, error) {
			return n, nil
		})

	userFactory := NewFactory(&User{}).
		SeqInt("ID", func(n int) (any, error) {
			return n, nil
		}).
		Attr("Name", func(args Args) (any, error) {
			return "hyuti", nil
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
		SeqInt("ID", func(n int) (any, error) {
			return n, nil
		})

	userFactory := NewFactory(&User{}).
		SeqInt("ID", func(n int) (any, error) {
			return n, nil
		}).
		Attr("Name", func(args Args) (any, error) {
			return "hyuti", nil
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
		SeqInt("ID", func(n int) (any, error) {
			return n, nil
		}).
		Attr("Name", func(args Args) (any, error) {
			return "hyuti", nil
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
		SeqInt("ID", func(n int) (any, error) {
			return n, nil
		}).
		Attr("Name", func(args Args) (any, error) {
			return "hyuti", nil
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
	if err := userFactory.ConstructWithOption(user, map[string]any{"Name": "jun"}); err != nil {
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
		Attr("Name", func(args Args) (any, error) {
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

func TestFactoryWithOptions(t *testing.T) {
	type (
		Group struct {
			Name string
		}
		User struct {
			ID     int
			Name   string
			Group1 Group
			Group2 *Group
		}
	)

	var userFactory = NewFactory(&User{})
	user := userFactory.MustCreateWithOption(map[string]any{
		"ID":          1,
		"Name":        "hyuti",
		"Group1.Name": "programmer",
		"Group2.Name": "web",
	}).(*User)

	if user.ID != 1 {
		t.Errorf("user.ID should be 1, not %v", user.ID)
	}

	if user.Name != "hyuti" {
		t.Errorf("user.Name should be hyuti, not %v", user.Name)
	}

	if user.Group1.Name != "programmer" {
		t.Errorf("user.Group1.Name should be programmer, not %v", user.Group1.Name)
	}

	if user.Group2.Name != "web" {
		t.Errorf("user.Group2.Name should be web, not %v", user.Group2.Name)
	}
}

func TestFactoryMuctCreateWithContextAndOptions(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	type ctxField int
	const nameField ctxField = 1

	var userFactory = NewFactory(&User{})

	t.Run("with valid options", func(t *testing.T) {
		user := userFactory.MustCreateWithContextAndOption(context.Background(), map[string]any{
			"ID":   1,
			"Name": "hyuti",
		}).(*User)

		if user.ID != 1 {
			t.Errorf("user.ID should be 1, not %v", user.ID)
		}

		if user.Name != "hyuti" {
			t.Errorf("user.Name should be hyuti, not %v", user.Name)
		}
	})

	t.Run("with broken options", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Errorf("func should panic")
			}
		}()

		userFactory.MustCreateWithContextAndOption(context.Background(), map[string]any{
			"ID":   1,
			"Name": 3,
		})
	})

	t.Run("with filled context", func(t *testing.T) {
		userFactory := NewFactory(&User{}).Attr("Name", func(args Args) (any, error) {
			return args.Context().Value(nameField), nil
		})

		ctx := context.WithValue(context.Background(), nameField, "hyuti from ctx")
		user := userFactory.MustCreateWithContextAndOption(ctx, map[string]any{
			"ID": 1,
		}).(*User)

		if user.Name != "hyuti from ctx" {
			t.Errorf("user.Name should be hyuti from ctx, not %v", user.Name)
		}
	})

	t.Run("with nil context", func(t *testing.T) {
		userFactory := NewFactory(&User{}).Attr("Name", func(args Args) (any, error) {
			return args.Context().Value(nameField), nil
		})

		defer func() {
			if recover() == nil {
				t.Errorf("func should panic")
			}
		}()

		userFactory.MustCreateWithContextAndOption(nil, map[string]any{
			"ID": 1,
		})
	})
}

func TestFactorySeqConcurrency(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	var userFactory = NewFactory(&User{}).
		SeqInt("ID", func(n int) (any, error) {
			return n, nil
		}).
		SeqString("Name", func(s string) (any, error) {
			return "user-" + s, nil
		})

	var wg sync.WaitGroup
	users := make([]*User, 1000)

	// Concurrently construct many different Users
	for i := range users {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			user, err := userFactory.Create()
			if err != nil {
				t.Errorf("constructing a User shouldn't have failed: %v", err)
			} else {
				users[i] = user.(*User)
			}
		}()
	}
	wg.Wait()

	// Check that each ID and Name value is unique
	ids := make(map[int]bool)
	names := make(map[string]bool)

	for _, user := range users {
		if ids[user.ID] {
			t.Errorf("found a repeated integer sequence value %d (user.ID)", user.ID)
		} else {
			ids[user.ID] = true
		}

		if names[user.Name] {
			t.Errorf("found a repeated string sequence value %s (user.Name)", user.Name)
		} else {
			names[user.Name] = true
		}
	}
}

func TestFactorySeqIntStartsAt1(t *testing.T) {
	type User struct {
		ID int
	}

	var userFactory = NewFactory(&User{}).
		SeqInt("ID", func(n int) (any, error) {
			return n, nil
		})

	user, err := userFactory.Create()
	if err != nil {
		t.Errorf("failed to create a User: %v", err)
	}

	if id := user.(*User).ID; id != 1 {
		t.Errorf("the starting number for SeqInt was %d, not 1", id)
	}
}

func TestFactorySeqInt64StartsAt1(t *testing.T) {
	type User struct {
		ID int64
	}

	var userFactory = NewFactory(&User{}).
		SeqInt64("ID", func(n int64) (any, error) {
			return n, nil
		})

	user, err := userFactory.Create()
	if err != nil {
		t.Errorf("failed to create a User: %v", err)
	}

	if id := user.(*User).ID; id != 1 {
		t.Errorf("the starting number for SeqInt was %d, not 1", id)
	}
}

func TestFactorySeqStringStartsAt1(t *testing.T) {
	type User struct {
		Name string
	}

	var userFactory = NewFactory(&User{}).
		SeqString("Name", func(s string) (any, error) {
			return s, nil
		})

	user, err := userFactory.Create()
	if err != nil {
		t.Errorf("failed to create a User: %v", err)
	}

	if name := user.(*User).Name; name != "1" {
		t.Errorf("the starting number for SeqString was %s, not 1", name)
	}
}

// this func tests create, not Create function which is exported one.
func TestCreate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setUp  func(*testing.T, context.Context)
		expect func(*testing.T, context.Context)
	}{
		{
			name:  "client not defines any attr gen",
			setUp: func(t *testing.T, ctx context.Context) {},
			expect: func(t *testing.T, ctx context.Context) {
				type Foo struct {
					Bar int
				}
				fooFactory := NewFactory(&Foo{})
				require.Len(t, fooFactory.orderingAttrGens, fooFactory.numField)
				fooAny, err := fooFactory.create(ctx, nil, nil)
				foo, _ := fooAny.(*Foo)
				for _, attr := range fooFactory.orderingAttrGens {
					require.True(t, attr.isFilled)
				}
				require.Nil(t, err)
				require.Empty(t, foo.Bar)
			},
		},
		{
			name:  "client defines ordering attr gen",
			setUp: func(t *testing.T, ctx context.Context) {},
			expect: func(t *testing.T, ctx context.Context) {
				type Foo struct {
					FooBar string
					Bar    int
				}
				fooFactory := NewFactory(&Foo{}).
					Attr("Bar", func(a Args) (any, error) {
						return 0, nil
					}).
					Attr("FooBar", func(a Args) (any, error) {
						return "foobar", nil
					})
				_, err := fooFactory.create(ctx, nil, nil)
				require.Nil(t, err)
				require.Equal(t, fooFactory.orderingAttrGens[0].key, "Bar")
				require.Equal(t, fooFactory.orderingAttrGens[1].key, "FooBar")
			},
		},
		{
			name:  "client defines partial ordering attr gen",
			setUp: func(t *testing.T, ctx context.Context) {},
			expect: func(t *testing.T, ctx context.Context) {
				type Foo struct {
					FooBar string
					Bar    int
					BarFoo int
				}
				fooFactory := NewFactory(&Foo{}).
					Attr("Bar", func(a Args) (any, error) {
						return 0, nil
					}).
					Attr("FooBar", func(a Args) (any, error) {
						return "foobar", nil
					})
				_, err := fooFactory.create(ctx, nil, nil)
				require.Nil(t, err)
				require.Equal(t, fooFactory.orderingAttrGens[0].key, "Bar")
				require.Equal(t, fooFactory.orderingAttrGens[1].key, "FooBar")
				require.Equal(t, fooFactory.orderingAttrGens[2].key, "BarFoo")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tt.setUp(t, ctx)
			tt.expect(t, ctx)
		})
	}
}

func TestFormatter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setUp  func(*testing.T, context.Context)
		expect func(*testing.T, context.Context)
	}{
		{
			name:  "subfactory with formatters",
			setUp: func(t *testing.T, ctx context.Context) {},
			expect: func(t *testing.T, ctx context.Context) {
				type Bar struct {
					ID int
				}
				type Foo struct {
					BarID int
				}
				barFactory := NewFactory(&Bar{ID: 1})
				fooFactory := NewFactory(&Foo{}).
					SubFactory("BarID", barFactory, func(i any) (any, error) {
						inst, ok := i.(*Bar)
						if !ok {
							return nil, fmt.Errorf("unexpected type %t", i)
						}
						return inst.ID, nil
					})
				fooAny, err := fooFactory.Create()
				foo, ok := fooAny.(*Foo)
				require.True(t, ok)
				require.Nil(t, err)
				require.Equal(t, foo.BarID, 1)
			},
		},
		{
			name:  "non-subfactory with formatters",
			setUp: func(t *testing.T, ctx context.Context) {},
			expect: func(t *testing.T, ctx context.Context) {
				type Bar struct {
					Foo string
				}
				barFactory := NewFactory(&Bar{}).
					Attr("Foo", func(a Args) (any, error) {
						return "foo", nil
					}, func(i any) (any, error) {
						inst, ok := i.(string)
						if !ok {
							return nil, fmt.Errorf("unexpected type %t", i)
						}
						return fmt.Sprintf("%sbar", inst), nil
					})
				barAny, err := barFactory.Create()
				bar, ok := barAny.(*Bar)
				require.True(t, ok)
				require.Nil(t, err)
				require.Equal(t, bar.Foo, "foobar")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tt.setUp(t, ctx)
			tt.expect(t, ctx)
		})
	}
}

func TestFileWithDir(t *testing.T) {
	t.Parallel()
	type kitTest struct {
		gen  func(*os.File) (any, error)
		ctrl *gomock.Controller
		args *MockArgs
	}
	tests := [...]struct {
		name     string
		setUp    func(*testing.T) *kitTest
		expect   func(*testing.T, *kitTest)
		tearDown func(*testing.T, *kitTest)
	}{
		{
			name: "create unsuccess",
			setUp: func(t *testing.T) *kitTest {
				gen := func(f *os.File) (any, error) {
					return f, errors.New("foo")
				}
				ctrl := gomock.NewController(t)
				args := NewMockArgs(ctrl)
				return &kitTest{
					gen:  gen,
					ctrl: ctrl,
					args: args,
				}
			},
			expect: func(t *testing.T, kt *kitTest) {
				type foobar struct {
					foo string
				}
				f := NewFactory(&foobar{})
				f = f.FileWithDir("foo", "txt", "temp", kt.gen)
				_, err := f.orderingAttrGens[0].genFunc(kt.args)
				require.NotNil(t, err)
			},
			tearDown: func(t *testing.T, kt *kitTest) {
				kt.ctrl.Finish()
			},
		},
		{
			name: "create success",
			setUp: func(t *testing.T) *kitTest {
				gen := func(f *os.File) (any, error) {
					return f, nil
				}
				ctrl := gomock.NewController(t)
				args := NewMockArgs(ctrl)
				return &kitTest{
					gen:  gen,
					ctrl: ctrl,
					args: args,
				}
			},
			expect: func(t *testing.T, kt *kitTest) {
				type foobar struct {
					foo string
				}
				f := NewFactory(&foobar{})
				f = f.FileWithDir("foo", "txt", "temp", kt.gen)
				fn, err := f.orderingAttrGens[0].genFunc(kt.args)
				require.Nil(t, err)
				require.Contains(t, fn, ".txt")
			},
			tearDown: func(t *testing.T, kt *kitTest) {
				kt.ctrl.Finish()
			},
		},
		{
			name: "clean",
			setUp: func(t *testing.T) *kitTest {
				gen := func(f *os.File) (any, error) {
					return f, nil
				}
				ctrl := gomock.NewController(t)
				args := NewMockArgs(ctrl)
				return &kitTest{
					gen:  gen,
					ctrl: ctrl,
					args: args,
				}
			},
			expect: func(t *testing.T, kt *kitTest) {
				type foobar struct {
					foo string
				}
				f := NewFactory(&foobar{})
				f = f.FileWithDir("foo", "txt", "temp", kt.gen)
				fnAny, _ := f.orderingAttrGens[0].genFunc(kt.args)
				fn, _ := fnAny.(string)
				f.Clean()
				err := os.Remove(fn)
				require.NotNil(t, err)
			},
			tearDown: func(t *testing.T, kt *kitTest) {
				kt.ctrl.Finish()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kt := tt.setUp(t)
			tt.expect(t, kt)
			tt.tearDown(t, kt)
		})
	}
}
