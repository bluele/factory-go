package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

// Book holds the schema definition for the Book entity.
type Book struct {
	ent.Schema
}

// Fields of the Book.
func (Book) Fields() []ent.Field {
	return nil
}

// Edges of the Book.
func (Book) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", Customer.Type).Ref("books").Unique(),
	}
}
