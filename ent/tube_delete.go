// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"coresamples/ent/predicate"
	"coresamples/ent/tube"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// TubeDelete is the builder for deleting a Tube entity.
type TubeDelete struct {
	config
	hooks    []Hook
	mutation *TubeMutation
}

// Where appends a list predicates to the TubeDelete builder.
func (td *TubeDelete) Where(ps ...predicate.Tube) *TubeDelete {
	td.mutation.Where(ps...)
	return td
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (td *TubeDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, td.sqlExec, td.mutation, td.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (td *TubeDelete) ExecX(ctx context.Context) int {
	n, err := td.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (td *TubeDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(tube.Table, sqlgraph.NewFieldSpec(tube.FieldID, field.TypeInt))
	if ps := td.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, td.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	td.mutation.done = true
	return affected, err
}

// TubeDeleteOne is the builder for deleting a single Tube entity.
type TubeDeleteOne struct {
	td *TubeDelete
}

// Where appends a list predicates to the TubeDelete builder.
func (tdo *TubeDeleteOne) Where(ps ...predicate.Tube) *TubeDeleteOne {
	tdo.td.mutation.Where(ps...)
	return tdo
}

// Exec executes the deletion query.
func (tdo *TubeDeleteOne) Exec(ctx context.Context) error {
	n, err := tdo.td.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{tube.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (tdo *TubeDeleteOne) ExecX(ctx context.Context) {
	if err := tdo.Exec(ctx); err != nil {
		panic(err)
	}
}
