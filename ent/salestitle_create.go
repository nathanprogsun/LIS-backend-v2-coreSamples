// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"coresamples/ent/salesteam"
	"coresamples/ent/salestitle"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// SalesTitleCreate is the builder for creating a SalesTitle entity.
type SalesTitleCreate struct {
	config
	mutation *SalesTitleMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetTitleName sets the "title_name" field.
func (stc *SalesTitleCreate) SetTitleName(s string) *SalesTitleCreate {
	stc.mutation.SetTitleName(s)
	return stc
}

// SetOrder sets the "order" field.
func (stc *SalesTitleCreate) SetOrder(i int) *SalesTitleCreate {
	stc.mutation.SetOrder(i)
	return stc
}

// SetCreatedTime sets the "created_time" field.
func (stc *SalesTitleCreate) SetCreatedTime(t time.Time) *SalesTitleCreate {
	stc.mutation.SetCreatedTime(t)
	return stc
}

// SetNillableCreatedTime sets the "created_time" field if the given value is not nil.
func (stc *SalesTitleCreate) SetNillableCreatedTime(t *time.Time) *SalesTitleCreate {
	if t != nil {
		stc.SetCreatedTime(*t)
	}
	return stc
}

// SetUpdatedTime sets the "updated_time" field.
func (stc *SalesTitleCreate) SetUpdatedTime(t time.Time) *SalesTitleCreate {
	stc.mutation.SetUpdatedTime(t)
	return stc
}

// SetNillableUpdatedTime sets the "updated_time" field if the given value is not nil.
func (stc *SalesTitleCreate) SetNillableUpdatedTime(t *time.Time) *SalesTitleCreate {
	if t != nil {
		stc.SetUpdatedTime(*t)
	}
	return stc
}

// AddSaleIDs adds the "sales" edge to the SalesTeam entity by IDs.
func (stc *SalesTitleCreate) AddSaleIDs(ids ...int) *SalesTitleCreate {
	stc.mutation.AddSaleIDs(ids...)
	return stc
}

// AddSales adds the "sales" edges to the SalesTeam entity.
func (stc *SalesTitleCreate) AddSales(s ...*SalesTeam) *SalesTitleCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stc.AddSaleIDs(ids...)
}

// Mutation returns the SalesTitleMutation object of the builder.
func (stc *SalesTitleCreate) Mutation() *SalesTitleMutation {
	return stc.mutation
}

// Save creates the SalesTitle in the database.
func (stc *SalesTitleCreate) Save(ctx context.Context) (*SalesTitle, error) {
	stc.defaults()
	return withHooks(ctx, stc.sqlSave, stc.mutation, stc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (stc *SalesTitleCreate) SaveX(ctx context.Context) *SalesTitle {
	v, err := stc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (stc *SalesTitleCreate) Exec(ctx context.Context) error {
	_, err := stc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (stc *SalesTitleCreate) ExecX(ctx context.Context) {
	if err := stc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (stc *SalesTitleCreate) defaults() {
	if _, ok := stc.mutation.CreatedTime(); !ok {
		v := salestitle.DefaultCreatedTime()
		stc.mutation.SetCreatedTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (stc *SalesTitleCreate) check() error {
	if _, ok := stc.mutation.TitleName(); !ok {
		return &ValidationError{Name: "title_name", err: errors.New(`ent: missing required field "SalesTitle.title_name"`)}
	}
	if _, ok := stc.mutation.Order(); !ok {
		return &ValidationError{Name: "order", err: errors.New(`ent: missing required field "SalesTitle.order"`)}
	}
	if _, ok := stc.mutation.CreatedTime(); !ok {
		return &ValidationError{Name: "created_time", err: errors.New(`ent: missing required field "SalesTitle.created_time"`)}
	}
	return nil
}

func (stc *SalesTitleCreate) sqlSave(ctx context.Context) (*SalesTitle, error) {
	if err := stc.check(); err != nil {
		return nil, err
	}
	_node, _spec := stc.createSpec()
	if err := sqlgraph.CreateNode(ctx, stc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	stc.mutation.id = &_node.ID
	stc.mutation.done = true
	return _node, nil
}

func (stc *SalesTitleCreate) createSpec() (*SalesTitle, *sqlgraph.CreateSpec) {
	var (
		_node = &SalesTitle{config: stc.config}
		_spec = sqlgraph.NewCreateSpec(salestitle.Table, sqlgraph.NewFieldSpec(salestitle.FieldID, field.TypeInt))
	)
	_spec.OnConflict = stc.conflict
	if value, ok := stc.mutation.TitleName(); ok {
		_spec.SetField(salestitle.FieldTitleName, field.TypeString, value)
		_node.TitleName = value
	}
	if value, ok := stc.mutation.Order(); ok {
		_spec.SetField(salestitle.FieldOrder, field.TypeInt, value)
		_node.Order = value
	}
	if value, ok := stc.mutation.CreatedTime(); ok {
		_spec.SetField(salestitle.FieldCreatedTime, field.TypeTime, value)
		_node.CreatedTime = value
	}
	if value, ok := stc.mutation.UpdatedTime(); ok {
		_spec.SetField(salestitle.FieldUpdatedTime, field.TypeTime, value)
		_node.UpdatedTime = value
	}
	if nodes := stc.mutation.SalesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   salestitle.SalesTable,
			Columns: []string{salestitle.SalesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(salesteam.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.SalesTitle.Create().
//		SetTitleName(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.SalesTitleUpsert) {
//			SetTitleName(v+v).
//		}).
//		Exec(ctx)
func (stc *SalesTitleCreate) OnConflict(opts ...sql.ConflictOption) *SalesTitleUpsertOne {
	stc.conflict = opts
	return &SalesTitleUpsertOne{
		create: stc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.SalesTitle.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (stc *SalesTitleCreate) OnConflictColumns(columns ...string) *SalesTitleUpsertOne {
	stc.conflict = append(stc.conflict, sql.ConflictColumns(columns...))
	return &SalesTitleUpsertOne{
		create: stc,
	}
}

type (
	// SalesTitleUpsertOne is the builder for "upsert"-ing
	//  one SalesTitle node.
	SalesTitleUpsertOne struct {
		create *SalesTitleCreate
	}

	// SalesTitleUpsert is the "OnConflict" setter.
	SalesTitleUpsert struct {
		*sql.UpdateSet
	}
)

// SetTitleName sets the "title_name" field.
func (u *SalesTitleUpsert) SetTitleName(v string) *SalesTitleUpsert {
	u.Set(salestitle.FieldTitleName, v)
	return u
}

// UpdateTitleName sets the "title_name" field to the value that was provided on create.
func (u *SalesTitleUpsert) UpdateTitleName() *SalesTitleUpsert {
	u.SetExcluded(salestitle.FieldTitleName)
	return u
}

// SetOrder sets the "order" field.
func (u *SalesTitleUpsert) SetOrder(v int) *SalesTitleUpsert {
	u.Set(salestitle.FieldOrder, v)
	return u
}

// UpdateOrder sets the "order" field to the value that was provided on create.
func (u *SalesTitleUpsert) UpdateOrder() *SalesTitleUpsert {
	u.SetExcluded(salestitle.FieldOrder)
	return u
}

// AddOrder adds v to the "order" field.
func (u *SalesTitleUpsert) AddOrder(v int) *SalesTitleUpsert {
	u.Add(salestitle.FieldOrder, v)
	return u
}

// SetCreatedTime sets the "created_time" field.
func (u *SalesTitleUpsert) SetCreatedTime(v time.Time) *SalesTitleUpsert {
	u.Set(salestitle.FieldCreatedTime, v)
	return u
}

// UpdateCreatedTime sets the "created_time" field to the value that was provided on create.
func (u *SalesTitleUpsert) UpdateCreatedTime() *SalesTitleUpsert {
	u.SetExcluded(salestitle.FieldCreatedTime)
	return u
}

// SetUpdatedTime sets the "updated_time" field.
func (u *SalesTitleUpsert) SetUpdatedTime(v time.Time) *SalesTitleUpsert {
	u.Set(salestitle.FieldUpdatedTime, v)
	return u
}

// UpdateUpdatedTime sets the "updated_time" field to the value that was provided on create.
func (u *SalesTitleUpsert) UpdateUpdatedTime() *SalesTitleUpsert {
	u.SetExcluded(salestitle.FieldUpdatedTime)
	return u
}

// ClearUpdatedTime clears the value of the "updated_time" field.
func (u *SalesTitleUpsert) ClearUpdatedTime() *SalesTitleUpsert {
	u.SetNull(salestitle.FieldUpdatedTime)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.SalesTitle.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *SalesTitleUpsertOne) UpdateNewValues() *SalesTitleUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.SalesTitle.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *SalesTitleUpsertOne) Ignore() *SalesTitleUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *SalesTitleUpsertOne) DoNothing() *SalesTitleUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the SalesTitleCreate.OnConflict
// documentation for more info.
func (u *SalesTitleUpsertOne) Update(set func(*SalesTitleUpsert)) *SalesTitleUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&SalesTitleUpsert{UpdateSet: update})
	}))
	return u
}

// SetTitleName sets the "title_name" field.
func (u *SalesTitleUpsertOne) SetTitleName(v string) *SalesTitleUpsertOne {
	return u.Update(func(s *SalesTitleUpsert) {
		s.SetTitleName(v)
	})
}

// UpdateTitleName sets the "title_name" field to the value that was provided on create.
func (u *SalesTitleUpsertOne) UpdateTitleName() *SalesTitleUpsertOne {
	return u.Update(func(s *SalesTitleUpsert) {
		s.UpdateTitleName()
	})
}

// SetOrder sets the "order" field.
func (u *SalesTitleUpsertOne) SetOrder(v int) *SalesTitleUpsertOne {
	return u.Update(func(s *SalesTitleUpsert) {
		s.SetOrder(v)
	})
}

// AddOrder adds v to the "order" field.
func (u *SalesTitleUpsertOne) AddOrder(v int) *SalesTitleUpsertOne {
	return u.Update(func(s *SalesTitleUpsert) {
		s.AddOrder(v)
	})
}

// UpdateOrder sets the "order" field to the value that was provided on create.
func (u *SalesTitleUpsertOne) UpdateOrder() *SalesTitleUpsertOne {
	return u.Update(func(s *SalesTitleUpsert) {
		s.UpdateOrder()
	})
}

// SetCreatedTime sets the "created_time" field.
func (u *SalesTitleUpsertOne) SetCreatedTime(v time.Time) *SalesTitleUpsertOne {
	return u.Update(func(s *SalesTitleUpsert) {
		s.SetCreatedTime(v)
	})
}

// UpdateCreatedTime sets the "created_time" field to the value that was provided on create.
func (u *SalesTitleUpsertOne) UpdateCreatedTime() *SalesTitleUpsertOne {
	return u.Update(func(s *SalesTitleUpsert) {
		s.UpdateCreatedTime()
	})
}

// SetUpdatedTime sets the "updated_time" field.
func (u *SalesTitleUpsertOne) SetUpdatedTime(v time.Time) *SalesTitleUpsertOne {
	return u.Update(func(s *SalesTitleUpsert) {
		s.SetUpdatedTime(v)
	})
}

// UpdateUpdatedTime sets the "updated_time" field to the value that was provided on create.
func (u *SalesTitleUpsertOne) UpdateUpdatedTime() *SalesTitleUpsertOne {
	return u.Update(func(s *SalesTitleUpsert) {
		s.UpdateUpdatedTime()
	})
}

// ClearUpdatedTime clears the value of the "updated_time" field.
func (u *SalesTitleUpsertOne) ClearUpdatedTime() *SalesTitleUpsertOne {
	return u.Update(func(s *SalesTitleUpsert) {
		s.ClearUpdatedTime()
	})
}

// Exec executes the query.
func (u *SalesTitleUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for SalesTitleCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *SalesTitleUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *SalesTitleUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *SalesTitleUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// SalesTitleCreateBulk is the builder for creating many SalesTitle entities in bulk.
type SalesTitleCreateBulk struct {
	config
	err      error
	builders []*SalesTitleCreate
	conflict []sql.ConflictOption
}

// Save creates the SalesTitle entities in the database.
func (stcb *SalesTitleCreateBulk) Save(ctx context.Context) ([]*SalesTitle, error) {
	if stcb.err != nil {
		return nil, stcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(stcb.builders))
	nodes := make([]*SalesTitle, len(stcb.builders))
	mutators := make([]Mutator, len(stcb.builders))
	for i := range stcb.builders {
		func(i int, root context.Context) {
			builder := stcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SalesTitleMutation)
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
					_, err = mutators[i+1].Mutate(root, stcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = stcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, stcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
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
		if _, err := mutators[0].Mutate(ctx, stcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (stcb *SalesTitleCreateBulk) SaveX(ctx context.Context) []*SalesTitle {
	v, err := stcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (stcb *SalesTitleCreateBulk) Exec(ctx context.Context) error {
	_, err := stcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (stcb *SalesTitleCreateBulk) ExecX(ctx context.Context) {
	if err := stcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.SalesTitle.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.SalesTitleUpsert) {
//			SetTitleName(v+v).
//		}).
//		Exec(ctx)
func (stcb *SalesTitleCreateBulk) OnConflict(opts ...sql.ConflictOption) *SalesTitleUpsertBulk {
	stcb.conflict = opts
	return &SalesTitleUpsertBulk{
		create: stcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.SalesTitle.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (stcb *SalesTitleCreateBulk) OnConflictColumns(columns ...string) *SalesTitleUpsertBulk {
	stcb.conflict = append(stcb.conflict, sql.ConflictColumns(columns...))
	return &SalesTitleUpsertBulk{
		create: stcb,
	}
}

// SalesTitleUpsertBulk is the builder for "upsert"-ing
// a bulk of SalesTitle nodes.
type SalesTitleUpsertBulk struct {
	create *SalesTitleCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.SalesTitle.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *SalesTitleUpsertBulk) UpdateNewValues() *SalesTitleUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.SalesTitle.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *SalesTitleUpsertBulk) Ignore() *SalesTitleUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *SalesTitleUpsertBulk) DoNothing() *SalesTitleUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the SalesTitleCreateBulk.OnConflict
// documentation for more info.
func (u *SalesTitleUpsertBulk) Update(set func(*SalesTitleUpsert)) *SalesTitleUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&SalesTitleUpsert{UpdateSet: update})
	}))
	return u
}

// SetTitleName sets the "title_name" field.
func (u *SalesTitleUpsertBulk) SetTitleName(v string) *SalesTitleUpsertBulk {
	return u.Update(func(s *SalesTitleUpsert) {
		s.SetTitleName(v)
	})
}

// UpdateTitleName sets the "title_name" field to the value that was provided on create.
func (u *SalesTitleUpsertBulk) UpdateTitleName() *SalesTitleUpsertBulk {
	return u.Update(func(s *SalesTitleUpsert) {
		s.UpdateTitleName()
	})
}

// SetOrder sets the "order" field.
func (u *SalesTitleUpsertBulk) SetOrder(v int) *SalesTitleUpsertBulk {
	return u.Update(func(s *SalesTitleUpsert) {
		s.SetOrder(v)
	})
}

// AddOrder adds v to the "order" field.
func (u *SalesTitleUpsertBulk) AddOrder(v int) *SalesTitleUpsertBulk {
	return u.Update(func(s *SalesTitleUpsert) {
		s.AddOrder(v)
	})
}

// UpdateOrder sets the "order" field to the value that was provided on create.
func (u *SalesTitleUpsertBulk) UpdateOrder() *SalesTitleUpsertBulk {
	return u.Update(func(s *SalesTitleUpsert) {
		s.UpdateOrder()
	})
}

// SetCreatedTime sets the "created_time" field.
func (u *SalesTitleUpsertBulk) SetCreatedTime(v time.Time) *SalesTitleUpsertBulk {
	return u.Update(func(s *SalesTitleUpsert) {
		s.SetCreatedTime(v)
	})
}

// UpdateCreatedTime sets the "created_time" field to the value that was provided on create.
func (u *SalesTitleUpsertBulk) UpdateCreatedTime() *SalesTitleUpsertBulk {
	return u.Update(func(s *SalesTitleUpsert) {
		s.UpdateCreatedTime()
	})
}

// SetUpdatedTime sets the "updated_time" field.
func (u *SalesTitleUpsertBulk) SetUpdatedTime(v time.Time) *SalesTitleUpsertBulk {
	return u.Update(func(s *SalesTitleUpsert) {
		s.SetUpdatedTime(v)
	})
}

// UpdateUpdatedTime sets the "updated_time" field to the value that was provided on create.
func (u *SalesTitleUpsertBulk) UpdateUpdatedTime() *SalesTitleUpsertBulk {
	return u.Update(func(s *SalesTitleUpsert) {
		s.UpdateUpdatedTime()
	})
}

// ClearUpdatedTime clears the value of the "updated_time" field.
func (u *SalesTitleUpsertBulk) ClearUpdatedTime() *SalesTitleUpsertBulk {
	return u.Update(func(s *SalesTitleUpsert) {
		s.ClearUpdatedTime()
	})
}

// Exec executes the query.
func (u *SalesTitleUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the SalesTitleCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for SalesTitleCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *SalesTitleUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
