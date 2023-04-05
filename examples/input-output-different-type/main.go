package main

import (
	"fmt"
	"reflect"

	"github.com/Pallinder/go-randomdata"
	"github.com/hyuti/factory-go/factory"
)

// This example illustrates on how to use factory in a case where you have different type of input and output
type (
	CustomerInput struct {
		Name string
	}
	Customer struct {
		ID   string
		Name string
	}
)

func CustomerRepository(i *CustomerInput) *Customer {
	return &Customer{
		ID:   "foo",
		Name: i.Name,
	}
}

var customerFactory = factory.NewFactory(
	&CustomerInput{},
).Attr("Name", func(a factory.Args) (any, error) {
	return randomdata.FullName(randomdata.RandomGender), nil
})

func CustomerFactory(opts map[string]any) *factory.Factory {
	return factory.NewFactory(
		&Customer{},
	).OnCreate(func(a factory.Args) error {
		ctx := a.Context()
		iAny, err := customerFactory.CreateWithContextAndOption(ctx, opts)
		if err != nil {
			return err
		}
		i, ok := iAny.(*CustomerInput)
		if !ok {
			return fmt.Errorf("unexpected type %t", iAny)
		}
		e := CustomerRepository(i)
		inst := a.Instance()
		dst := reflect.ValueOf(inst)
		src := reflect.ValueOf(e).Elem()
		dst.Elem().Set(src)
		return nil
	})
}

func main() {
	customerAny, err := CustomerFactory(nil).Create()
	if err != nil {
		fmt.Println(err)
	} else {
		customer, ok := customerAny.(*Customer)
		if !ok {
			fmt.Printf("unexpected type %t\n", customerAny)
		} else {
			fmt.Printf("Name: %s, ID: %s\n", customer.Name, customer.ID)
		}
	}
}
