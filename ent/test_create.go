// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"coresamples/ent/orderinfo"
	"coresamples/ent/sampletype"
	"coresamples/ent/test"
	"coresamples/ent/testdetail"
	"coresamples/ent/tubetype"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// TestCreate is the builder for creating a Test entity.
type TestCreate struct {
	config
	mutation *TestMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetTestName sets the "test_name" field.
func (tc *TestCreate) SetTestName(s string) *TestCreate {
	tc.mutation.SetTestName(s)
	return tc
}

// SetTestCode sets the "test_code" field.
func (tc *TestCreate) SetTestCode(s string) *TestCreate {
	tc.mutation.SetTestCode(s)
	return tc
}

// SetDisplayName sets the "display_name" field.
func (tc *TestCreate) SetDisplayName(s string) *TestCreate {
	tc.mutation.SetDisplayName(s)
	return tc
}

// SetTestDescription sets the "test_description" field.
func (tc *TestCreate) SetTestDescription(s string) *TestCreate {
	tc.mutation.SetTestDescription(s)
	return tc
}

// SetAssayName sets the "assay_name" field.
func (tc *TestCreate) SetAssayName(s string) *TestCreate {
	tc.mutation.SetAssayName(s)
	return tc
}

// SetIsActive sets the "isActive" field.
func (tc *TestCreate) SetIsActive(b bool) *TestCreate {
	tc.mutation.SetIsActive(b)
	return tc
}

// SetCreatedTime sets the "created_time" field.
func (tc *TestCreate) SetCreatedTime(t time.Time) *TestCreate {
	tc.mutation.SetCreatedTime(t)
	return tc
}

// SetNillableCreatedTime sets the "created_time" field if the given value is not nil.
func (tc *TestCreate) SetNillableCreatedTime(t *time.Time) *TestCreate {
	if t != nil {
		tc.SetCreatedTime(*t)
	}
	return tc
}

// SetUpdatedTime sets the "updated_time" field.
func (tc *TestCreate) SetUpdatedTime(t time.Time) *TestCreate {
	tc.mutation.SetUpdatedTime(t)
	return tc
}

// SetNillableUpdatedTime sets the "updated_time" field if the given value is not nil.
func (tc *TestCreate) SetNillableUpdatedTime(t *time.Time) *TestCreate {
	if t != nil {
		tc.SetUpdatedTime(*t)
	}
	return tc
}

// SetID sets the "id" field.
func (tc *TestCreate) SetID(i int) *TestCreate {
	tc.mutation.SetID(i)
	return tc
}

// AddTestDetailIDs adds the "test_details" edge to the TestDetail entity by IDs.
func (tc *TestCreate) AddTestDetailIDs(ids ...int) *TestCreate {
	tc.mutation.AddTestDetailIDs(ids...)
	return tc
}

// AddTestDetails adds the "test_details" edges to the TestDetail entity.
func (tc *TestCreate) AddTestDetails(t ...*TestDetail) *TestCreate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tc.AddTestDetailIDs(ids...)
}

// AddOrderInfoIDs adds the "order_info" edge to the OrderInfo entity by IDs.
func (tc *TestCreate) AddOrderInfoIDs(ids ...int) *TestCreate {
	tc.mutation.AddOrderInfoIDs(ids...)
	return tc
}

// AddOrderInfo adds the "order_info" edges to the OrderInfo entity.
func (tc *TestCreate) AddOrderInfo(o ...*OrderInfo) *TestCreate {
	ids := make([]int, len(o))
	for i := range o {
		ids[i] = o[i].ID
	}
	return tc.AddOrderInfoIDs(ids...)
}

// AddSampleTypeIDs adds the "sample_types" edge to the SampleType entity by IDs.
func (tc *TestCreate) AddSampleTypeIDs(ids ...int) *TestCreate {
	tc.mutation.AddSampleTypeIDs(ids...)
	return tc
}

// AddSampleTypes adds the "sample_types" edges to the SampleType entity.
func (tc *TestCreate) AddSampleTypes(s ...*SampleType) *TestCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return tc.AddSampleTypeIDs(ids...)
}

// AddTubeTypeIDs adds the "tube_types" edge to the TubeType entity by IDs.
func (tc *TestCreate) AddTubeTypeIDs(ids ...int) *TestCreate {
	tc.mutation.AddTubeTypeIDs(ids...)
	return tc
}

// AddTubeTypes adds the "tube_types" edges to the TubeType entity.
func (tc *TestCreate) AddTubeTypes(t ...*TubeType) *TestCreate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tc.AddTubeTypeIDs(ids...)
}

// Mutation returns the TestMutation object of the builder.
func (tc *TestCreate) Mutation() *TestMutation {
	return tc.mutation
}

// Save creates the Test in the database.
func (tc *TestCreate) Save(ctx context.Context) (*Test, error) {
	tc.defaults()
	return withHooks(ctx, tc.sqlSave, tc.mutation, tc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (tc *TestCreate) SaveX(ctx context.Context) *Test {
	v, err := tc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tc *TestCreate) Exec(ctx context.Context) error {
	_, err := tc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tc *TestCreate) ExecX(ctx context.Context) {
	if err := tc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tc *TestCreate) defaults() {
	if _, ok := tc.mutation.CreatedTime(); !ok {
		v := test.DefaultCreatedTime()
		tc.mutation.SetCreatedTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tc *TestCreate) check() error {
	if _, ok := tc.mutation.TestName(); !ok {
		return &ValidationError{Name: "test_name", err: errors.New(`ent: missing required field "Test.test_name"`)}
	}
	if _, ok := tc.mutation.TestCode(); !ok {
		return &ValidationError{Name: "test_code", err: errors.New(`ent: missing required field "Test.test_code"`)}
	}
	if _, ok := tc.mutation.DisplayName(); !ok {
		return &ValidationError{Name: "display_name", err: errors.New(`ent: missing required field "Test.display_name"`)}
	}
	if _, ok := tc.mutation.TestDescription(); !ok {
		return &ValidationError{Name: "test_description", err: errors.New(`ent: missing required field "Test.test_description"`)}
	}
	if _, ok := tc.mutation.AssayName(); !ok {
		return &ValidationError{Name: "assay_name", err: errors.New(`ent: missing required field "Test.assay_name"`)}
	}
	if _, ok := tc.mutation.IsActive(); !ok {
		return &ValidationError{Name: "isActive", err: errors.New(`ent: missing required field "Test.isActive"`)}
	}
	if _, ok := tc.mutation.CreatedTime(); !ok {
		return &ValidationError{Name: "created_time", err: errors.New(`ent: missing required field "Test.created_time"`)}
	}
	return nil
}

func (tc *TestCreate) sqlSave(ctx context.Context) (*Test, error) {
	if err := tc.check(); err != nil {
		return nil, err
	}
	_node, _spec := tc.createSpec()
	if err := sqlgraph.CreateNode(ctx, tc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = int(id)
	}
	tc.mutation.id = &_node.ID
	tc.mutation.done = true
	return _node, nil
}

func (tc *TestCreate) createSpec() (*Test, *sqlgraph.CreateSpec) {
	var (
		_node = &Test{config: tc.config}
		_spec = sqlgraph.NewCreateSpec(test.Table, sqlgraph.NewFieldSpec(test.FieldID, field.TypeInt))
	)
	_spec.OnConflict = tc.conflict
	if id, ok := tc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := tc.mutation.TestName(); ok {
		_spec.SetField(test.FieldTestName, field.TypeString, value)
		_node.TestName = value
	}
	if value, ok := tc.mutation.TestCode(); ok {
		_spec.SetField(test.FieldTestCode, field.TypeString, value)
		_node.TestCode = value
	}
	if value, ok := tc.mutation.DisplayName(); ok {
		_spec.SetField(test.FieldDisplayName, field.TypeString, value)
		_node.DisplayName = value
	}
	if value, ok := tc.mutation.TestDescription(); ok {
		_spec.SetField(test.FieldTestDescription, field.TypeString, value)
		_node.TestDescription = &value
	}
	if value, ok := tc.mutation.AssayName(); ok {
		_spec.SetField(test.FieldAssayName, field.TypeString, value)
		_node.AssayName = &value
	}
	if value, ok := tc.mutation.IsActive(); ok {
		_spec.SetField(test.FieldIsActive, field.TypeBool, value)
		_node.IsActive = value
	}
	if value, ok := tc.mutation.CreatedTime(); ok {
		_spec.SetField(test.FieldCreatedTime, field.TypeTime, value)
		_node.CreatedTime = value
	}
	if value, ok := tc.mutation.UpdatedTime(); ok {
		_spec.SetField(test.FieldUpdatedTime, field.TypeTime, value)
		_node.UpdatedTime = value
	}
	if nodes := tc.mutation.TestDetailsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   test.TestDetailsTable,
			Columns: []string{test.TestDetailsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(testdetail.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tc.mutation.OrderInfoIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   test.OrderInfoTable,
			Columns: test.OrderInfoPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(orderinfo.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tc.mutation.SampleTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   test.SampleTypesTable,
			Columns: test.SampleTypesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(sampletype.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tc.mutation.TubeTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   test.TubeTypesTable,
			Columns: test.TubeTypesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tubetype.FieldID, field.TypeInt),
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
//	client.Test.Create().
//		SetTestName(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.TestUpsert) {
//			SetTestName(v+v).
//		}).
//		Exec(ctx)
func (tc *TestCreate) OnConflict(opts ...sql.ConflictOption) *TestUpsertOne {
	tc.conflict = opts
	return &TestUpsertOne{
		create: tc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Test.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (tc *TestCreate) OnConflictColumns(columns ...string) *TestUpsertOne {
	tc.conflict = append(tc.conflict, sql.ConflictColumns(columns...))
	return &TestUpsertOne{
		create: tc,
	}
}

type (
	// TestUpsertOne is the builder for "upsert"-ing
	//  one Test node.
	TestUpsertOne struct {
		create *TestCreate
	}

	// TestUpsert is the "OnConflict" setter.
	TestUpsert struct {
		*sql.UpdateSet
	}
)

// SetTestName sets the "test_name" field.
func (u *TestUpsert) SetTestName(v string) *TestUpsert {
	u.Set(test.FieldTestName, v)
	return u
}

// UpdateTestName sets the "test_name" field to the value that was provided on create.
func (u *TestUpsert) UpdateTestName() *TestUpsert {
	u.SetExcluded(test.FieldTestName)
	return u
}

// SetTestCode sets the "test_code" field.
func (u *TestUpsert) SetTestCode(v string) *TestUpsert {
	u.Set(test.FieldTestCode, v)
	return u
}

// UpdateTestCode sets the "test_code" field to the value that was provided on create.
func (u *TestUpsert) UpdateTestCode() *TestUpsert {
	u.SetExcluded(test.FieldTestCode)
	return u
}

// SetDisplayName sets the "display_name" field.
func (u *TestUpsert) SetDisplayName(v string) *TestUpsert {
	u.Set(test.FieldDisplayName, v)
	return u
}

// UpdateDisplayName sets the "display_name" field to the value that was provided on create.
func (u *TestUpsert) UpdateDisplayName() *TestUpsert {
	u.SetExcluded(test.FieldDisplayName)
	return u
}

// SetTestDescription sets the "test_description" field.
func (u *TestUpsert) SetTestDescription(v string) *TestUpsert {
	u.Set(test.FieldTestDescription, v)
	return u
}

// UpdateTestDescription sets the "test_description" field to the value that was provided on create.
func (u *TestUpsert) UpdateTestDescription() *TestUpsert {
	u.SetExcluded(test.FieldTestDescription)
	return u
}

// SetAssayName sets the "assay_name" field.
func (u *TestUpsert) SetAssayName(v string) *TestUpsert {
	u.Set(test.FieldAssayName, v)
	return u
}

// UpdateAssayName sets the "assay_name" field to the value that was provided on create.
func (u *TestUpsert) UpdateAssayName() *TestUpsert {
	u.SetExcluded(test.FieldAssayName)
	return u
}

// SetIsActive sets the "isActive" field.
func (u *TestUpsert) SetIsActive(v bool) *TestUpsert {
	u.Set(test.FieldIsActive, v)
	return u
}

// UpdateIsActive sets the "isActive" field to the value that was provided on create.
func (u *TestUpsert) UpdateIsActive() *TestUpsert {
	u.SetExcluded(test.FieldIsActive)
	return u
}

// SetCreatedTime sets the "created_time" field.
func (u *TestUpsert) SetCreatedTime(v time.Time) *TestUpsert {
	u.Set(test.FieldCreatedTime, v)
	return u
}

// UpdateCreatedTime sets the "created_time" field to the value that was provided on create.
func (u *TestUpsert) UpdateCreatedTime() *TestUpsert {
	u.SetExcluded(test.FieldCreatedTime)
	return u
}

// SetUpdatedTime sets the "updated_time" field.
func (u *TestUpsert) SetUpdatedTime(v time.Time) *TestUpsert {
	u.Set(test.FieldUpdatedTime, v)
	return u
}

// UpdateUpdatedTime sets the "updated_time" field to the value that was provided on create.
func (u *TestUpsert) UpdateUpdatedTime() *TestUpsert {
	u.SetExcluded(test.FieldUpdatedTime)
	return u
}

// ClearUpdatedTime clears the value of the "updated_time" field.
func (u *TestUpsert) ClearUpdatedTime() *TestUpsert {
	u.SetNull(test.FieldUpdatedTime)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Test.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(test.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *TestUpsertOne) UpdateNewValues() *TestUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(test.FieldID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Test.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *TestUpsertOne) Ignore() *TestUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *TestUpsertOne) DoNothing() *TestUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the TestCreate.OnConflict
// documentation for more info.
func (u *TestUpsertOne) Update(set func(*TestUpsert)) *TestUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&TestUpsert{UpdateSet: update})
	}))
	return u
}

// SetTestName sets the "test_name" field.
func (u *TestUpsertOne) SetTestName(v string) *TestUpsertOne {
	return u.Update(func(s *TestUpsert) {
		s.SetTestName(v)
	})
}

// UpdateTestName sets the "test_name" field to the value that was provided on create.
func (u *TestUpsertOne) UpdateTestName() *TestUpsertOne {
	return u.Update(func(s *TestUpsert) {
		s.UpdateTestName()
	})
}

// SetTestCode sets the "test_code" field.
func (u *TestUpsertOne) SetTestCode(v string) *TestUpsertOne {
	return u.Update(func(s *TestUpsert) {
		s.SetTestCode(v)
	})
}

// UpdateTestCode sets the "test_code" field to the value that was provided on create.
func (u *TestUpsertOne) UpdateTestCode() *TestUpsertOne {
	return u.Update(func(s *TestUpsert) {
		s.UpdateTestCode()
	})
}

// SetDisplayName sets the "display_name" field.
func (u *TestUpsertOne) SetDisplayName(v string) *TestUpsertOne {
	return u.Update(func(s *TestUpsert) {
		s.SetDisplayName(v)
	})
}

// UpdateDisplayName sets the "display_name" field to the value that was provided on create.
func (u *TestUpsertOne) UpdateDisplayName() *TestUpsertOne {
	return u.Update(func(s *TestUpsert) {
		s.UpdateDisplayName()
	})
}

// SetTestDescription sets the "test_description" field.
func (u *TestUpsertOne) SetTestDescription(v string) *TestUpsertOne {
	return u.Update(func(s *TestUpsert) {
		s.SetTestDescription(v)
	})
}

// UpdateTestDescription sets the "test_description" field to the value that was provided on create.
func (u *TestUpsertOne) UpdateTestDescription() *TestUpsertOne {
	return u.Update(func(s *TestUpsert) {
		s.UpdateTestDescription()
	})
}

// SetAssayName sets the "assay_name" field.
func (u *TestUpsertOne) SetAssayName(v string) *TestUpsertOne {
	return u.Update(func(s *TestUpsert) {
		s.SetAssayName(v)
	})
}

// UpdateAssayName sets the "assay_name" field to the value that was provided on create.
func (u *TestUpsertOne) UpdateAssayName() *TestUpsertOne {
	return u.Update(func(s *TestUpsert) {
		s.UpdateAssayName()
	})
}

// SetIsActive sets the "isActive" field.
func (u *TestUpsertOne) SetIsActive(v bool) *TestUpsertOne {
	return u.Update(func(s *TestUpsert) {
		s.SetIsActive(v)
	})
}

// UpdateIsActive sets the "isActive" field to the value that was provided on create.
func (u *TestUpsertOne) UpdateIsActive() *TestUpsertOne {
	return u.Update(func(s *TestUpsert) {
		s.UpdateIsActive()
	})
}

// SetCreatedTime sets the "created_time" field.
func (u *TestUpsertOne) SetCreatedTime(v time.Time) *TestUpsertOne {
	return u.Update(func(s *TestUpsert) {
		s.SetCreatedTime(v)
	})
}

// UpdateCreatedTime sets the "created_time" field to the value that was provided on create.
func (u *TestUpsertOne) UpdateCreatedTime() *TestUpsertOne {
	return u.Update(func(s *TestUpsert) {
		s.UpdateCreatedTime()
	})
}

// SetUpdatedTime sets the "updated_time" field.
func (u *TestUpsertOne) SetUpdatedTime(v time.Time) *TestUpsertOne {
	return u.Update(func(s *TestUpsert) {
		s.SetUpdatedTime(v)
	})
}

// UpdateUpdatedTime sets the "updated_time" field to the value that was provided on create.
func (u *TestUpsertOne) UpdateUpdatedTime() *TestUpsertOne {
	return u.Update(func(s *TestUpsert) {
		s.UpdateUpdatedTime()
	})
}

// ClearUpdatedTime clears the value of the "updated_time" field.
func (u *TestUpsertOne) ClearUpdatedTime() *TestUpsertOne {
	return u.Update(func(s *TestUpsert) {
		s.ClearUpdatedTime()
	})
}

// Exec executes the query.
func (u *TestUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for TestCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *TestUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *TestUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *TestUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// TestCreateBulk is the builder for creating many Test entities in bulk.
type TestCreateBulk struct {
	config
	err      error
	builders []*TestCreate
	conflict []sql.ConflictOption
}

// Save creates the Test entities in the database.
func (tcb *TestCreateBulk) Save(ctx context.Context) ([]*Test, error) {
	if tcb.err != nil {
		return nil, tcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(tcb.builders))
	nodes := make([]*Test, len(tcb.builders))
	mutators := make([]Mutator, len(tcb.builders))
	for i := range tcb.builders {
		func(i int, root context.Context) {
			builder := tcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TestMutation)
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
					_, err = mutators[i+1].Mutate(root, tcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = tcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, tcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, tcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (tcb *TestCreateBulk) SaveX(ctx context.Context) []*Test {
	v, err := tcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tcb *TestCreateBulk) Exec(ctx context.Context) error {
	_, err := tcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcb *TestCreateBulk) ExecX(ctx context.Context) {
	if err := tcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Test.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.TestUpsert) {
//			SetTestName(v+v).
//		}).
//		Exec(ctx)
func (tcb *TestCreateBulk) OnConflict(opts ...sql.ConflictOption) *TestUpsertBulk {
	tcb.conflict = opts
	return &TestUpsertBulk{
		create: tcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Test.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (tcb *TestCreateBulk) OnConflictColumns(columns ...string) *TestUpsertBulk {
	tcb.conflict = append(tcb.conflict, sql.ConflictColumns(columns...))
	return &TestUpsertBulk{
		create: tcb,
	}
}

// TestUpsertBulk is the builder for "upsert"-ing
// a bulk of Test nodes.
type TestUpsertBulk struct {
	create *TestCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Test.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(test.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *TestUpsertBulk) UpdateNewValues() *TestUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(test.FieldID)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Test.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *TestUpsertBulk) Ignore() *TestUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *TestUpsertBulk) DoNothing() *TestUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the TestCreateBulk.OnConflict
// documentation for more info.
func (u *TestUpsertBulk) Update(set func(*TestUpsert)) *TestUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&TestUpsert{UpdateSet: update})
	}))
	return u
}

// SetTestName sets the "test_name" field.
func (u *TestUpsertBulk) SetTestName(v string) *TestUpsertBulk {
	return u.Update(func(s *TestUpsert) {
		s.SetTestName(v)
	})
}

// UpdateTestName sets the "test_name" field to the value that was provided on create.
func (u *TestUpsertBulk) UpdateTestName() *TestUpsertBulk {
	return u.Update(func(s *TestUpsert) {
		s.UpdateTestName()
	})
}

// SetTestCode sets the "test_code" field.
func (u *TestUpsertBulk) SetTestCode(v string) *TestUpsertBulk {
	return u.Update(func(s *TestUpsert) {
		s.SetTestCode(v)
	})
}

// UpdateTestCode sets the "test_code" field to the value that was provided on create.
func (u *TestUpsertBulk) UpdateTestCode() *TestUpsertBulk {
	return u.Update(func(s *TestUpsert) {
		s.UpdateTestCode()
	})
}

// SetDisplayName sets the "display_name" field.
func (u *TestUpsertBulk) SetDisplayName(v string) *TestUpsertBulk {
	return u.Update(func(s *TestUpsert) {
		s.SetDisplayName(v)
	})
}

// UpdateDisplayName sets the "display_name" field to the value that was provided on create.
func (u *TestUpsertBulk) UpdateDisplayName() *TestUpsertBulk {
	return u.Update(func(s *TestUpsert) {
		s.UpdateDisplayName()
	})
}

// SetTestDescription sets the "test_description" field.
func (u *TestUpsertBulk) SetTestDescription(v string) *TestUpsertBulk {
	return u.Update(func(s *TestUpsert) {
		s.SetTestDescription(v)
	})
}

// UpdateTestDescription sets the "test_description" field to the value that was provided on create.
func (u *TestUpsertBulk) UpdateTestDescription() *TestUpsertBulk {
	return u.Update(func(s *TestUpsert) {
		s.UpdateTestDescription()
	})
}

// SetAssayName sets the "assay_name" field.
func (u *TestUpsertBulk) SetAssayName(v string) *TestUpsertBulk {
	return u.Update(func(s *TestUpsert) {
		s.SetAssayName(v)
	})
}

// UpdateAssayName sets the "assay_name" field to the value that was provided on create.
func (u *TestUpsertBulk) UpdateAssayName() *TestUpsertBulk {
	return u.Update(func(s *TestUpsert) {
		s.UpdateAssayName()
	})
}

// SetIsActive sets the "isActive" field.
func (u *TestUpsertBulk) SetIsActive(v bool) *TestUpsertBulk {
	return u.Update(func(s *TestUpsert) {
		s.SetIsActive(v)
	})
}

// UpdateIsActive sets the "isActive" field to the value that was provided on create.
func (u *TestUpsertBulk) UpdateIsActive() *TestUpsertBulk {
	return u.Update(func(s *TestUpsert) {
		s.UpdateIsActive()
	})
}

// SetCreatedTime sets the "created_time" field.
func (u *TestUpsertBulk) SetCreatedTime(v time.Time) *TestUpsertBulk {
	return u.Update(func(s *TestUpsert) {
		s.SetCreatedTime(v)
	})
}

// UpdateCreatedTime sets the "created_time" field to the value that was provided on create.
func (u *TestUpsertBulk) UpdateCreatedTime() *TestUpsertBulk {
	return u.Update(func(s *TestUpsert) {
		s.UpdateCreatedTime()
	})
}

// SetUpdatedTime sets the "updated_time" field.
func (u *TestUpsertBulk) SetUpdatedTime(v time.Time) *TestUpsertBulk {
	return u.Update(func(s *TestUpsert) {
		s.SetUpdatedTime(v)
	})
}

// UpdateUpdatedTime sets the "updated_time" field to the value that was provided on create.
func (u *TestUpsertBulk) UpdateUpdatedTime() *TestUpsertBulk {
	return u.Update(func(s *TestUpsert) {
		s.UpdateUpdatedTime()
	})
}

// ClearUpdatedTime clears the value of the "updated_time" field.
func (u *TestUpsertBulk) ClearUpdatedTime() *TestUpsertBulk {
	return u.Update(func(s *TestUpsert) {
		s.ClearUpdatedTime()
	})
}

// Exec executes the query.
func (u *TestUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the TestCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for TestCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *TestUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
