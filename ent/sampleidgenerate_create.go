// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"coresamples/ent/sampleidgenerate"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// SampleIDGenerateCreate is the builder for creating a SampleIDGenerate entity.
type SampleIDGenerateCreate struct {
	config
	mutation *SampleIDGenerateMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetBarcode sets the "barcode" field.
func (sigc *SampleIDGenerateCreate) SetBarcode(s string) *SampleIDGenerateCreate {
	sigc.mutation.SetBarcode(s)
	return sigc
}

// SetNillableBarcode sets the "barcode" field if the given value is not nil.
func (sigc *SampleIDGenerateCreate) SetNillableBarcode(s *string) *SampleIDGenerateCreate {
	if s != nil {
		sigc.SetBarcode(*s)
	}
	return sigc
}

// SetID sets the "id" field.
func (sigc *SampleIDGenerateCreate) SetID(i int) *SampleIDGenerateCreate {
	sigc.mutation.SetID(i)
	return sigc
}

// Mutation returns the SampleIDGenerateMutation object of the builder.
func (sigc *SampleIDGenerateCreate) Mutation() *SampleIDGenerateMutation {
	return sigc.mutation
}

// Save creates the SampleIDGenerate in the database.
func (sigc *SampleIDGenerateCreate) Save(ctx context.Context) (*SampleIDGenerate, error) {
	return withHooks(ctx, sigc.sqlSave, sigc.mutation, sigc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (sigc *SampleIDGenerateCreate) SaveX(ctx context.Context) *SampleIDGenerate {
	v, err := sigc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sigc *SampleIDGenerateCreate) Exec(ctx context.Context) error {
	_, err := sigc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sigc *SampleIDGenerateCreate) ExecX(ctx context.Context) {
	if err := sigc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sigc *SampleIDGenerateCreate) check() error {
	if v, ok := sigc.mutation.ID(); ok {
		if err := sampleidgenerate.IDValidator(v); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`ent: validator failed for field "SampleIDGenerate.id": %w`, err)}
		}
	}
	return nil
}

func (sigc *SampleIDGenerateCreate) sqlSave(ctx context.Context) (*SampleIDGenerate, error) {
	if err := sigc.check(); err != nil {
		return nil, err
	}
	_node, _spec := sigc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sigc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = int(id)
	}
	sigc.mutation.id = &_node.ID
	sigc.mutation.done = true
	return _node, nil
}

func (sigc *SampleIDGenerateCreate) createSpec() (*SampleIDGenerate, *sqlgraph.CreateSpec) {
	var (
		_node = &SampleIDGenerate{config: sigc.config}
		_spec = sqlgraph.NewCreateSpec(sampleidgenerate.Table, sqlgraph.NewFieldSpec(sampleidgenerate.FieldID, field.TypeInt))
	)
	_spec.OnConflict = sigc.conflict
	if id, ok := sigc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := sigc.mutation.Barcode(); ok {
		_spec.SetField(sampleidgenerate.FieldBarcode, field.TypeString, value)
		_node.Barcode = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.SampleIDGenerate.Create().
//		SetBarcode(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.SampleIDGenerateUpsert) {
//			SetBarcode(v+v).
//		}).
//		Exec(ctx)
func (sigc *SampleIDGenerateCreate) OnConflict(opts ...sql.ConflictOption) *SampleIDGenerateUpsertOne {
	sigc.conflict = opts
	return &SampleIDGenerateUpsertOne{
		create: sigc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.SampleIDGenerate.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (sigc *SampleIDGenerateCreate) OnConflictColumns(columns ...string) *SampleIDGenerateUpsertOne {
	sigc.conflict = append(sigc.conflict, sql.ConflictColumns(columns...))
	return &SampleIDGenerateUpsertOne{
		create: sigc,
	}
}

type (
	// SampleIDGenerateUpsertOne is the builder for "upsert"-ing
	//  one SampleIDGenerate node.
	SampleIDGenerateUpsertOne struct {
		create *SampleIDGenerateCreate
	}

	// SampleIDGenerateUpsert is the "OnConflict" setter.
	SampleIDGenerateUpsert struct {
		*sql.UpdateSet
	}
)

// SetBarcode sets the "barcode" field.
func (u *SampleIDGenerateUpsert) SetBarcode(v string) *SampleIDGenerateUpsert {
	u.Set(sampleidgenerate.FieldBarcode, v)
	return u
}

// UpdateBarcode sets the "barcode" field to the value that was provided on create.
func (u *SampleIDGenerateUpsert) UpdateBarcode() *SampleIDGenerateUpsert {
	u.SetExcluded(sampleidgenerate.FieldBarcode)
	return u
}

// ClearBarcode clears the value of the "barcode" field.
func (u *SampleIDGenerateUpsert) ClearBarcode() *SampleIDGenerateUpsert {
	u.SetNull(sampleidgenerate.FieldBarcode)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.SampleIDGenerate.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(sampleidgenerate.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *SampleIDGenerateUpsertOne) UpdateNewValues() *SampleIDGenerateUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(sampleidgenerate.FieldID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.SampleIDGenerate.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *SampleIDGenerateUpsertOne) Ignore() *SampleIDGenerateUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *SampleIDGenerateUpsertOne) DoNothing() *SampleIDGenerateUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the SampleIDGenerateCreate.OnConflict
// documentation for more info.
func (u *SampleIDGenerateUpsertOne) Update(set func(*SampleIDGenerateUpsert)) *SampleIDGenerateUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&SampleIDGenerateUpsert{UpdateSet: update})
	}))
	return u
}

// SetBarcode sets the "barcode" field.
func (u *SampleIDGenerateUpsertOne) SetBarcode(v string) *SampleIDGenerateUpsertOne {
	return u.Update(func(s *SampleIDGenerateUpsert) {
		s.SetBarcode(v)
	})
}

// UpdateBarcode sets the "barcode" field to the value that was provided on create.
func (u *SampleIDGenerateUpsertOne) UpdateBarcode() *SampleIDGenerateUpsertOne {
	return u.Update(func(s *SampleIDGenerateUpsert) {
		s.UpdateBarcode()
	})
}

// ClearBarcode clears the value of the "barcode" field.
func (u *SampleIDGenerateUpsertOne) ClearBarcode() *SampleIDGenerateUpsertOne {
	return u.Update(func(s *SampleIDGenerateUpsert) {
		s.ClearBarcode()
	})
}

// Exec executes the query.
func (u *SampleIDGenerateUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for SampleIDGenerateCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *SampleIDGenerateUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *SampleIDGenerateUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *SampleIDGenerateUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// SampleIDGenerateCreateBulk is the builder for creating many SampleIDGenerate entities in bulk.
type SampleIDGenerateCreateBulk struct {
	config
	err      error
	builders []*SampleIDGenerateCreate
	conflict []sql.ConflictOption
}

// Save creates the SampleIDGenerate entities in the database.
func (sigcb *SampleIDGenerateCreateBulk) Save(ctx context.Context) ([]*SampleIDGenerate, error) {
	if sigcb.err != nil {
		return nil, sigcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(sigcb.builders))
	nodes := make([]*SampleIDGenerate, len(sigcb.builders))
	mutators := make([]Mutator, len(sigcb.builders))
	for i := range sigcb.builders {
		func(i int, root context.Context) {
			builder := sigcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SampleIDGenerateMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, sigcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = sigcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, sigcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil && nodes[i].ID == 0 {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, sigcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (sigcb *SampleIDGenerateCreateBulk) SaveX(ctx context.Context) []*SampleIDGenerate {
	v, err := sigcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sigcb *SampleIDGenerateCreateBulk) Exec(ctx context.Context) error {
	_, err := sigcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sigcb *SampleIDGenerateCreateBulk) ExecX(ctx context.Context) {
	if err := sigcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.SampleIDGenerate.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.SampleIDGenerateUpsert) {
//			SetBarcode(v+v).
//		}).
//		Exec(ctx)
func (sigcb *SampleIDGenerateCreateBulk) OnConflict(opts ...sql.ConflictOption) *SampleIDGenerateUpsertBulk {
	sigcb.conflict = opts
	return &SampleIDGenerateUpsertBulk{
		create: sigcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.SampleIDGenerate.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (sigcb *SampleIDGenerateCreateBulk) OnConflictColumns(columns ...string) *SampleIDGenerateUpsertBulk {
	sigcb.conflict = append(sigcb.conflict, sql.ConflictColumns(columns...))
	return &SampleIDGenerateUpsertBulk{
		create: sigcb,
	}
}

// SampleIDGenerateUpsertBulk is the builder for "upsert"-ing
// a bulk of SampleIDGenerate nodes.
type SampleIDGenerateUpsertBulk struct {
	create *SampleIDGenerateCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.SampleIDGenerate.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(sampleidgenerate.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *SampleIDGenerateUpsertBulk) UpdateNewValues() *SampleIDGenerateUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(sampleidgenerate.FieldID)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.SampleIDGenerate.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *SampleIDGenerateUpsertBulk) Ignore() *SampleIDGenerateUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *SampleIDGenerateUpsertBulk) DoNothing() *SampleIDGenerateUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the SampleIDGenerateCreateBulk.OnConflict
// documentation for more info.
func (u *SampleIDGenerateUpsertBulk) Update(set func(*SampleIDGenerateUpsert)) *SampleIDGenerateUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&SampleIDGenerateUpsert{UpdateSet: update})
	}))
	return u
}

// SetBarcode sets the "barcode" field.
func (u *SampleIDGenerateUpsertBulk) SetBarcode(v string) *SampleIDGenerateUpsertBulk {
	return u.Update(func(s *SampleIDGenerateUpsert) {
		s.SetBarcode(v)
	})
}

// UpdateBarcode sets the "barcode" field to the value that was provided on create.
func (u *SampleIDGenerateUpsertBulk) UpdateBarcode() *SampleIDGenerateUpsertBulk {
	return u.Update(func(s *SampleIDGenerateUpsert) {
		s.UpdateBarcode()
	})
}

// ClearBarcode clears the value of the "barcode" field.
func (u *SampleIDGenerateUpsertBulk) ClearBarcode() *SampleIDGenerateUpsertBulk {
	return u.Update(func(s *SampleIDGenerateUpsert) {
		s.ClearBarcode()
	})
}

// Exec executes the query.
func (u *SampleIDGenerateUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the SampleIDGenerateCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for SampleIDGenerateCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *SampleIDGenerateUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
