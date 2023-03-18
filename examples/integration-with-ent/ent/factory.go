package ent

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/Pallinder/go-randomdata"
	"github.com/hyuti/factory-go/factory"
	_ "github.com/mattn/go-sqlite3"
)

type (
	ClientCtxType string
	Opt           struct {
		Key   string
		Value any
	}
	IFactory[ModelType any] interface {
		Create(context.Context) (ModelType, error)
		CreateWithClient(context.Context, *Client) (ModelType, error)
	}

	Factory[ModelType any] struct {
		worker *factory.Factory
	}
)

func Zero[T any]() T {
	var t T
	return t
}

const ClientCtxKey ClientCtxType = "client"

func (s *Factory[ModelType]) Create(ctx context.Context) (ModelType, error) {
	eAny, err := s.worker.CreateWithContext(ctx)
	if err != nil {
		return Zero[ModelType](), err
	}
	e, ok := eAny.(ModelType)
	if !ok {
		return Zero[ModelType](), fmt.Errorf("unexpected type %t", eAny)
	}
	return e, nil
}

func (s *Factory[ModelType]) CreateWithClient(ctx context.Context, client *Client) (ModelType, error) {
	EmbedClient(&ctx, client)
	return s.Create(ctx)
}

func getClient(ctx context.Context) (*Client, error) {
	client, ok := ctx.Value(ClientCtxKey).(*Client)
	if !ok || client == nil {
		return nil, fmt.Errorf("cannot find client in context")
	}
	return client, nil
}
func convertInputToOutput[ModelInputType, ModelType any](
	ctx context.Context,
	args factory.Args,
	factory *factory.Factory,
	saver func(context.Context, *Client, ModelInputType) (ModelType, error),
	opts ...Opt,
) error {
	optMap := make(map[string]any)
	for _, opt := range opts {
		optMap[opt.Key] = opt.Value
	}
	iAny, err := factory.CreateWithContextAndOption(ctx, optMap)
	if err != nil {
		return err
	}
	i, ok := iAny.(ModelInputType)
	if !ok {
		return fmt.Errorf("unexpected type %t", iAny)
	}
	client, err := getClient(ctx)
	if err != nil {
		return err
	}
	e, err := saver(ctx, client, i)
	if err != nil {
		return err
	}
	inst := args.Instance()
	dst := reflect.ValueOf(inst)
	src := reflect.ValueOf(e).Elem()
	dst.Elem().Set(src)
	return nil
}
func factoryTemplate[ModelType, ModelInputType any](
	model ModelType,
	f *factory.Factory,
	saver func(context.Context, *Client, ModelInputType) (ModelType, error),
	opts ...Opt,
) *factory.Factory {
	return factory.NewFactory(
		model,
	).OnCreate(func(a factory.Args) error {
		ctx := a.Context()
		err := convertInputToOutput(
			ctx,
			a,
			f,
			saver,
			opts...,
		)
		if err != nil {
			return err
		}
		return nil
	})
}
func EmbedClient(ctx *context.Context, v *Client) {
	c := *ctx
	client := c.Value(ClientCtxKey)
	if client == nil {
		*ctx = context.WithValue(*ctx, ClientCtxKey, v)
	}
}

// All you have to do is define a ModelFactory, you don't need to pay attention to all stuff above. But it's nice for you if you give it a shot.
// Also this example assumes you have ent installed. if not refer to https://github.com/ent/ent to get everything inplace
// For example, here is my best practice so far.
type (
	BookInput struct {
		OwnerID int
	}
	CustomerInput struct {
		Name string
	}
)

var customerFactory = factory.NewFactory(
	&CustomerInput{},
).Attr("Name", func(a factory.Args) (interface{}, error) {
	return randomdata.FullName(randomdata.RandomGender), nil
})
var bookFactory = factory.NewFactory(
	&BookInput{},
).SubFactory("OwnerID", CustomerFactory(), func(i interface{}) (interface{}, error) {
	e, ok := i.(*Customer)
	if !ok {
		return nil, fmt.Errorf("unexpected type %t", i)
	}
	return e.ID, nil
})

func CustomerFactory(opts ...Opt) *factory.Factory {
	return factoryTemplate(
		new(Customer),
		customerFactory,
		func(ctx context.Context, client *Client, i *CustomerInput) (*Customer, error) {
			return client.Customer.Create().SetName(i.Name).Save(ctx)
		},
		opts...,
	)
}
func BookFactory(opts ...Opt) *factory.Factory {
	return factoryTemplate(
		new(Book),
		bookFactory,
		func(ctx context.Context, client *Client, i *BookInput) (*Book, error) {
			return client.Book.Create().SetOwnerID(i.OwnerID).Save(ctx)
		},
		opts...,
	)
}
func MustCustomerFactory(opts ...Opt) IFactory[*Customer] {
	return &Factory[*Customer]{
		worker: CustomerFactory(opts...),
	}
}
func MustBookFactory(opts ...Opt) IFactory[*Book] {
	return &Factory[*Book]{
		worker: BookFactory(opts...),
	}
}
func OpenClient() (*Client, error) {
	db, err := sql.Open("sqlite3", "file:ent?mode=memory&_fk=1")
	if err != nil {
		return nil, err
	}
	drv := entsql.OpenDB(dialect.SQLite, db)
	return NewClient(Driver(drv)), nil
}

func New() (*Client, error) {
	return OpenClient()
}
