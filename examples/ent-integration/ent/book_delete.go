// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/hyuti/factory-go/examples/ent-integration/ent/book"
	"github.com/hyuti/factory-go/examples/ent-integration/ent/predicate"
)

// BookDelete is the builder for deleting a Book entity.
type BookDelete struct {
	config
	hooks    []Hook
	mutation *BookMutation
}

// Where appends a list predicates to the BookDelete builder.
func (bd *BookDelete) Where(ps ...predicate.Book) *BookDelete {
	bd.mutation.Where(ps...)
	return bd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (bd *BookDelete) Exec(ctx context.Context) (int, error) {
	return withHooks[int, BookMutation](ctx, bd.sqlExec, bd.mutation, bd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (bd *BookDelete) ExecX(ctx context.Context) int {
	n, err := bd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (bd *BookDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(book.Table, sqlgraph.NewFieldSpec(book.FieldID, field.TypeInt))
	if ps := bd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, bd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	bd.mutation.done = true
	return affected, err
}

// BookDeleteOne is the builder for deleting a single Book entity.
type BookDeleteOne struct {
	bd *BookDelete
}

// Where appends a list predicates to the BookDelete builder.
func (bdo *BookDeleteOne) Where(ps ...predicate.Book) *BookDeleteOne {
	bdo.bd.mutation.Where(ps...)
	return bdo
}

// Exec executes the deletion query.
func (bdo *BookDeleteOne) Exec(ctx context.Context) error {
	n, err := bdo.bd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{book.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (bdo *BookDeleteOne) ExecX(ctx context.Context) {
	if err := bdo.Exec(ctx); err != nil {
		panic(err)
	}
}
