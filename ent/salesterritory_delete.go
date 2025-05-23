// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"coresamples/ent/predicate"
	"coresamples/ent/salesterritory"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// SalesTerritoryDelete is the builder for deleting a SalesTerritory entity.
type SalesTerritoryDelete struct {
	config
	hooks    []Hook
	mutation *SalesTerritoryMutation
}

// Where appends a list predicates to the SalesTerritoryDelete builder.
func (std *SalesTerritoryDelete) Where(ps ...predicate.SalesTerritory) *SalesTerritoryDelete {
	std.mutation.Where(ps...)
	return std
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (std *SalesTerritoryDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, std.sqlExec, std.mutation, std.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (std *SalesTerritoryDelete) ExecX(ctx context.Context) int {
	n, err := std.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (std *SalesTerritoryDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(salesterritory.Table, sqlgraph.NewFieldSpec(salesterritory.FieldID, field.TypeInt))
	if ps := std.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, std.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	std.mutation.done = true
	return affected, err
}

// SalesTerritoryDeleteOne is the builder for deleting a single SalesTerritory entity.
type SalesTerritoryDeleteOne struct {
	std *SalesTerritoryDelete
}

// Where appends a list predicates to the SalesTerritoryDelete builder.
func (stdo *SalesTerritoryDeleteOne) Where(ps ...predicate.SalesTerritory) *SalesTerritoryDeleteOne {
	stdo.std.mutation.Where(ps...)
	return stdo
}

// Exec executes the deletion query.
func (stdo *SalesTerritoryDeleteOne) Exec(ctx context.Context) error {
	n, err := stdo.std.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{salesterritory.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (stdo *SalesTerritoryDeleteOne) ExecX(ctx context.Context) {
	if err := stdo.Exec(ctx); err != nil {
		panic(err)
	}
}
