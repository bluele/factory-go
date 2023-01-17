package main

import (
	"fmt"

	"github.com/bluele/factory-go/factory"
)

type (
	Customer struct {
		ID   int
		Name string
	}
	Order struct {
		ID           int
		CustomerName string
		CustomerID   int
	}
)

func GetCustomerByID(id int) *Customer {
	return &Customer{
		ID:   id,
		Name: "foo",
	}
}

var CustomerFactory = factory.NewFactory(
	&Customer{ID: 1},
).Attr("Name", func(a factory.Args) (interface{}, error) {
	return "foo", nil
})

var OrderFactory = factory.NewFactory(
	&Order{},

// define CustomerID before CustomerName so that CustomerID field is generated before CustomerName
// user a formatter to retrieve ID from Customer after Customer is created by CustomerFactory
).SubFactory("CustomerID", CustomerFactory, func(i interface{}) (interface{}, error) {
	e, ok := i.(*Customer)
	if !ok {
		return nil, fmt.Errorf("unexpected type %t", i)
	}
	return e.ID, nil
}).Attr("CustomerName", func(a factory.Args) (interface{}, error) {
	inst, ok := a.Instance().(*Order)
	if !ok {
		return nil, fmt.Errorf("unexpected type %t", a.Instance())
	}
	e := GetCustomerByID(inst.CustomerID)
	return e.Name, nil
})

func main() {
	order := OrderFactory.MustCreate().(*Order)
	fmt.Println("ID:", order.ID, " CustomerName:", order.CustomerName, " CustomerID:", order.CustomerID)
}
