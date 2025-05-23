// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"coresamples/ent/betaprogramparticipation"
	"coresamples/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// BetaProgramParticipationDelete is the builder for deleting a BetaProgramParticipation entity.
type BetaProgramParticipationDelete struct {
	config
	hooks    []Hook
	mutation *BetaProgramParticipationMutation
}

// Where appends a list predicates to the BetaProgramParticipationDelete builder.
func (bppd *BetaProgramParticipationDelete) Where(ps ...predicate.BetaProgramParticipation) *BetaProgramParticipationDelete {
	bppd.mutation.Where(ps...)
	return bppd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (bppd *BetaProgramParticipationDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, bppd.sqlExec, bppd.mutation, bppd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (bppd *BetaProgramParticipationDelete) ExecX(ctx context.Context) int {
	n, err := bppd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (bppd *BetaProgramParticipationDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(betaprogramparticipation.Table, sqlgraph.NewFieldSpec(betaprogramparticipation.FieldID, field.TypeInt))
	if ps := bppd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, bppd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	bppd.mutation.done = true
	return affected, err
}

// BetaProgramParticipationDeleteOne is the builder for deleting a single BetaProgramParticipation entity.
type BetaProgramParticipationDeleteOne struct {
	bppd *BetaProgramParticipationDelete
}

// Where appends a list predicates to the BetaProgramParticipationDelete builder.
func (bppdo *BetaProgramParticipationDeleteOne) Where(ps ...predicate.BetaProgramParticipation) *BetaProgramParticipationDeleteOne {
	bppdo.bppd.mutation.Where(ps...)
	return bppdo
}

// Exec executes the deletion query.
func (bppdo *BetaProgramParticipationDeleteOne) Exec(ctx context.Context) error {
	n, err := bppdo.bppd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{betaprogramparticipation.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (bppdo *BetaProgramParticipationDeleteOne) ExecX(ctx context.Context) {
	if err := bppdo.Exec(ctx); err != nil {
		panic(err)
	}
}
