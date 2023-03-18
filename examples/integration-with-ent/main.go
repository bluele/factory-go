package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hyuti/factory-go/examples/integration-with-ent/ent"
)

func main() {
	c, _ := ent.New()
	defer c.Close()
	ctx := context.Background()

	if err := c.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	customer, _ := ent.MustCustomerFactory(ent.Opt{Key: "Name", Value: "foobar"}).CreateWithClient(ctx, c)
	fmt.Println("Customer.ID: ", customer.ID, " Customer.Name: ", customer.Name)
	// Output:
	// Customer.ID: 1 Customer.Name: foobar

	bookWithOwner, _ := ent.MustBookFactory(ent.Opt{Key: "OwnerID", Value: customer.ID}).CreateWithClient(ctx, c)
	fmt.Println("Book.ID: ", bookWithOwner.ID)
	fmt.Println("\tOwner.ID: ", bookWithOwner.QueryOwner().FirstIDX(ctx), " Owner.Name: ", bookWithOwner.QueryOwner().FirstX(ctx).Name)
	// Output:
	// Book.ID: 1
	// 		Owner.ID: 1 Owner.Name: foobar

	book, _ := ent.MustBookFactory().CreateWithClient(ctx, c)
	fmt.Println("Book.ID: ", book.ID)
	fmt.Println("\tOwner.ID: ", book.QueryOwner().FirstIDX(ctx), " Owner.Name: ", book.QueryOwner().FirstX(ctx).Name)
	// Output:
	// Book.ID: 2
	// 		Owner.ID: 2 Owner.Name: random-name
}
