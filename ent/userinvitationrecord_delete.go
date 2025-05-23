// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"coresamples/ent/predicate"
	"coresamples/ent/userinvitationrecord"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// UserInvitationRecordDelete is the builder for deleting a UserInvitationRecord entity.
type UserInvitationRecordDelete struct {
	config
	hooks    []Hook
	mutation *UserInvitationRecordMutation
}

// Where appends a list predicates to the UserInvitationRecordDelete builder.
func (uird *UserInvitationRecordDelete) Where(ps ...predicate.UserInvitationRecord) *UserInvitationRecordDelete {
	uird.mutation.Where(ps...)
	return uird
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (uird *UserInvitationRecordDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, uird.sqlExec, uird.mutation, uird.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (uird *UserInvitationRecordDelete) ExecX(ctx context.Context) int {
	n, err := uird.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (uird *UserInvitationRecordDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(userinvitationrecord.Table, sqlgraph.NewFieldSpec(userinvitationrecord.FieldID, field.TypeInt))
	if ps := uird.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, uird.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	uird.mutation.done = true
	return affected, err
}

// UserInvitationRecordDeleteOne is the builder for deleting a single UserInvitationRecord entity.
type UserInvitationRecordDeleteOne struct {
	uird *UserInvitationRecordDelete
}

// Where appends a list predicates to the UserInvitationRecordDelete builder.
func (uirdo *UserInvitationRecordDeleteOne) Where(ps ...predicate.UserInvitationRecord) *UserInvitationRecordDeleteOne {
	uirdo.uird.mutation.Where(ps...)
	return uirdo
}

// Exec executes the deletion query.
func (uirdo *UserInvitationRecordDeleteOne) Exec(ctx context.Context) error {
	n, err := uirdo.uird.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{userinvitationrecord.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (uirdo *UserInvitationRecordDeleteOne) ExecX(ctx context.Context) {
	if err := uirdo.Exec(ctx); err != nil {
		panic(err)
	}
}
