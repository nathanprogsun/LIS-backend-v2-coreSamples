// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"coresamples/ent/predicate"
	"coresamples/ent/zipcode"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// ZipcodeUpdate is the builder for updating Zipcode entities.
type ZipcodeUpdate struct {
	config
	hooks    []Hook
	mutation *ZipcodeMutation
}

// Where appends a list predicates to the ZipcodeUpdate builder.
func (zu *ZipcodeUpdate) Where(ps ...predicate.Zipcode) *ZipcodeUpdate {
	zu.mutation.Where(ps...)
	return zu
}

// SetZipCodeType sets the "ZipCodeType" field.
func (zu *ZipcodeUpdate) SetZipCodeType(s string) *ZipcodeUpdate {
	zu.mutation.SetZipCodeType(s)
	return zu
}

// SetNillableZipCodeType sets the "ZipCodeType" field if the given value is not nil.
func (zu *ZipcodeUpdate) SetNillableZipCodeType(s *string) *ZipcodeUpdate {
	if s != nil {
		zu.SetZipCodeType(*s)
	}
	return zu
}

// SetCity sets the "City" field.
func (zu *ZipcodeUpdate) SetCity(s string) *ZipcodeUpdate {
	zu.mutation.SetCity(s)
	return zu
}

// SetNillableCity sets the "City" field if the given value is not nil.
func (zu *ZipcodeUpdate) SetNillableCity(s *string) *ZipcodeUpdate {
	if s != nil {
		zu.SetCity(*s)
	}
	return zu
}

// SetState sets the "State" field.
func (zu *ZipcodeUpdate) SetState(s string) *ZipcodeUpdate {
	zu.mutation.SetState(s)
	return zu
}

// SetNillableState sets the "State" field if the given value is not nil.
func (zu *ZipcodeUpdate) SetNillableState(s *string) *ZipcodeUpdate {
	if s != nil {
		zu.SetState(*s)
	}
	return zu
}

// SetLocationType sets the "LocationType" field.
func (zu *ZipcodeUpdate) SetLocationType(s string) *ZipcodeUpdate {
	zu.mutation.SetLocationType(s)
	return zu
}

// SetNillableLocationType sets the "LocationType" field if the given value is not nil.
func (zu *ZipcodeUpdate) SetNillableLocationType(s *string) *ZipcodeUpdate {
	if s != nil {
		zu.SetLocationType(*s)
	}
	return zu
}

// SetLat sets the "Lat" field.
func (zu *ZipcodeUpdate) SetLat(f float64) *ZipcodeUpdate {
	zu.mutation.ResetLat()
	zu.mutation.SetLat(f)
	return zu
}

// SetNillableLat sets the "Lat" field if the given value is not nil.
func (zu *ZipcodeUpdate) SetNillableLat(f *float64) *ZipcodeUpdate {
	if f != nil {
		zu.SetLat(*f)
	}
	return zu
}

// AddLat adds f to the "Lat" field.
func (zu *ZipcodeUpdate) AddLat(f float64) *ZipcodeUpdate {
	zu.mutation.AddLat(f)
	return zu
}

// ClearLat clears the value of the "Lat" field.
func (zu *ZipcodeUpdate) ClearLat() *ZipcodeUpdate {
	zu.mutation.ClearLat()
	return zu
}

// SetLong sets the "Long" field.
func (zu *ZipcodeUpdate) SetLong(f float64) *ZipcodeUpdate {
	zu.mutation.ResetLong()
	zu.mutation.SetLong(f)
	return zu
}

// SetNillableLong sets the "Long" field if the given value is not nil.
func (zu *ZipcodeUpdate) SetNillableLong(f *float64) *ZipcodeUpdate {
	if f != nil {
		zu.SetLong(*f)
	}
	return zu
}

// AddLong adds f to the "Long" field.
func (zu *ZipcodeUpdate) AddLong(f float64) *ZipcodeUpdate {
	zu.mutation.AddLong(f)
	return zu
}

// ClearLong clears the value of the "Long" field.
func (zu *ZipcodeUpdate) ClearLong() *ZipcodeUpdate {
	zu.mutation.ClearLong()
	return zu
}

// Mutation returns the ZipcodeMutation object of the builder.
func (zu *ZipcodeUpdate) Mutation() *ZipcodeMutation {
	return zu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (zu *ZipcodeUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, zu.sqlSave, zu.mutation, zu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (zu *ZipcodeUpdate) SaveX(ctx context.Context) int {
	affected, err := zu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (zu *ZipcodeUpdate) Exec(ctx context.Context) error {
	_, err := zu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (zu *ZipcodeUpdate) ExecX(ctx context.Context) {
	if err := zu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (zu *ZipcodeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(zipcode.Table, zipcode.Columns, sqlgraph.NewFieldSpec(zipcode.FieldID, field.TypeInt))
	if ps := zu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := zu.mutation.ZipCodeType(); ok {
		_spec.SetField(zipcode.FieldZipCodeType, field.TypeString, value)
	}
	if value, ok := zu.mutation.City(); ok {
		_spec.SetField(zipcode.FieldCity, field.TypeString, value)
	}
	if value, ok := zu.mutation.State(); ok {
		_spec.SetField(zipcode.FieldState, field.TypeString, value)
	}
	if value, ok := zu.mutation.LocationType(); ok {
		_spec.SetField(zipcode.FieldLocationType, field.TypeString, value)
	}
	if value, ok := zu.mutation.Lat(); ok {
		_spec.SetField(zipcode.FieldLat, field.TypeFloat64, value)
	}
	if value, ok := zu.mutation.AddedLat(); ok {
		_spec.AddField(zipcode.FieldLat, field.TypeFloat64, value)
	}
	if zu.mutation.LatCleared() {
		_spec.ClearField(zipcode.FieldLat, field.TypeFloat64)
	}
	if value, ok := zu.mutation.Long(); ok {
		_spec.SetField(zipcode.FieldLong, field.TypeFloat64, value)
	}
	if value, ok := zu.mutation.AddedLong(); ok {
		_spec.AddField(zipcode.FieldLong, field.TypeFloat64, value)
	}
	if zu.mutation.LongCleared() {
		_spec.ClearField(zipcode.FieldLong, field.TypeFloat64)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, zu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{zipcode.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	zu.mutation.done = true
	return n, nil
}

// ZipcodeUpdateOne is the builder for updating a single Zipcode entity.
type ZipcodeUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ZipcodeMutation
}

// SetZipCodeType sets the "ZipCodeType" field.
func (zuo *ZipcodeUpdateOne) SetZipCodeType(s string) *ZipcodeUpdateOne {
	zuo.mutation.SetZipCodeType(s)
	return zuo
}

// SetNillableZipCodeType sets the "ZipCodeType" field if the given value is not nil.
func (zuo *ZipcodeUpdateOne) SetNillableZipCodeType(s *string) *ZipcodeUpdateOne {
	if s != nil {
		zuo.SetZipCodeType(*s)
	}
	return zuo
}

// SetCity sets the "City" field.
func (zuo *ZipcodeUpdateOne) SetCity(s string) *ZipcodeUpdateOne {
	zuo.mutation.SetCity(s)
	return zuo
}

// SetNillableCity sets the "City" field if the given value is not nil.
func (zuo *ZipcodeUpdateOne) SetNillableCity(s *string) *ZipcodeUpdateOne {
	if s != nil {
		zuo.SetCity(*s)
	}
	return zuo
}

// SetState sets the "State" field.
func (zuo *ZipcodeUpdateOne) SetState(s string) *ZipcodeUpdateOne {
	zuo.mutation.SetState(s)
	return zuo
}

// SetNillableState sets the "State" field if the given value is not nil.
func (zuo *ZipcodeUpdateOne) SetNillableState(s *string) *ZipcodeUpdateOne {
	if s != nil {
		zuo.SetState(*s)
	}
	return zuo
}

// SetLocationType sets the "LocationType" field.
func (zuo *ZipcodeUpdateOne) SetLocationType(s string) *ZipcodeUpdateOne {
	zuo.mutation.SetLocationType(s)
	return zuo
}

// SetNillableLocationType sets the "LocationType" field if the given value is not nil.
func (zuo *ZipcodeUpdateOne) SetNillableLocationType(s *string) *ZipcodeUpdateOne {
	if s != nil {
		zuo.SetLocationType(*s)
	}
	return zuo
}

// SetLat sets the "Lat" field.
func (zuo *ZipcodeUpdateOne) SetLat(f float64) *ZipcodeUpdateOne {
	zuo.mutation.ResetLat()
	zuo.mutation.SetLat(f)
	return zuo
}

// SetNillableLat sets the "Lat" field if the given value is not nil.
func (zuo *ZipcodeUpdateOne) SetNillableLat(f *float64) *ZipcodeUpdateOne {
	if f != nil {
		zuo.SetLat(*f)
	}
	return zuo
}

// AddLat adds f to the "Lat" field.
func (zuo *ZipcodeUpdateOne) AddLat(f float64) *ZipcodeUpdateOne {
	zuo.mutation.AddLat(f)
	return zuo
}

// ClearLat clears the value of the "Lat" field.
func (zuo *ZipcodeUpdateOne) ClearLat() *ZipcodeUpdateOne {
	zuo.mutation.ClearLat()
	return zuo
}

// SetLong sets the "Long" field.
func (zuo *ZipcodeUpdateOne) SetLong(f float64) *ZipcodeUpdateOne {
	zuo.mutation.ResetLong()
	zuo.mutation.SetLong(f)
	return zuo
}

// SetNillableLong sets the "Long" field if the given value is not nil.
func (zuo *ZipcodeUpdateOne) SetNillableLong(f *float64) *ZipcodeUpdateOne {
	if f != nil {
		zuo.SetLong(*f)
	}
	return zuo
}

// AddLong adds f to the "Long" field.
func (zuo *ZipcodeUpdateOne) AddLong(f float64) *ZipcodeUpdateOne {
	zuo.mutation.AddLong(f)
	return zuo
}

// ClearLong clears the value of the "Long" field.
func (zuo *ZipcodeUpdateOne) ClearLong() *ZipcodeUpdateOne {
	zuo.mutation.ClearLong()
	return zuo
}

// Mutation returns the ZipcodeMutation object of the builder.
func (zuo *ZipcodeUpdateOne) Mutation() *ZipcodeMutation {
	return zuo.mutation
}

// Where appends a list predicates to the ZipcodeUpdate builder.
func (zuo *ZipcodeUpdateOne) Where(ps ...predicate.Zipcode) *ZipcodeUpdateOne {
	zuo.mutation.Where(ps...)
	return zuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (zuo *ZipcodeUpdateOne) Select(field string, fields ...string) *ZipcodeUpdateOne {
	zuo.fields = append([]string{field}, fields...)
	return zuo
}

// Save executes the query and returns the updated Zipcode entity.
func (zuo *ZipcodeUpdateOne) Save(ctx context.Context) (*Zipcode, error) {
	return withHooks(ctx, zuo.sqlSave, zuo.mutation, zuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (zuo *ZipcodeUpdateOne) SaveX(ctx context.Context) *Zipcode {
	node, err := zuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (zuo *ZipcodeUpdateOne) Exec(ctx context.Context) error {
	_, err := zuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (zuo *ZipcodeUpdateOne) ExecX(ctx context.Context) {
	if err := zuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (zuo *ZipcodeUpdateOne) sqlSave(ctx context.Context) (_node *Zipcode, err error) {
	_spec := sqlgraph.NewUpdateSpec(zipcode.Table, zipcode.Columns, sqlgraph.NewFieldSpec(zipcode.FieldID, field.TypeInt))
	id, ok := zuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Zipcode.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := zuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, zipcode.FieldID)
		for _, f := range fields {
			if !zipcode.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != zipcode.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := zuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := zuo.mutation.ZipCodeType(); ok {
		_spec.SetField(zipcode.FieldZipCodeType, field.TypeString, value)
	}
	if value, ok := zuo.mutation.City(); ok {
		_spec.SetField(zipcode.FieldCity, field.TypeString, value)
	}
	if value, ok := zuo.mutation.State(); ok {
		_spec.SetField(zipcode.FieldState, field.TypeString, value)
	}
	if value, ok := zuo.mutation.LocationType(); ok {
		_spec.SetField(zipcode.FieldLocationType, field.TypeString, value)
	}
	if value, ok := zuo.mutation.Lat(); ok {
		_spec.SetField(zipcode.FieldLat, field.TypeFloat64, value)
	}
	if value, ok := zuo.mutation.AddedLat(); ok {
		_spec.AddField(zipcode.FieldLat, field.TypeFloat64, value)
	}
	if zuo.mutation.LatCleared() {
		_spec.ClearField(zipcode.FieldLat, field.TypeFloat64)
	}
	if value, ok := zuo.mutation.Long(); ok {
		_spec.SetField(zipcode.FieldLong, field.TypeFloat64, value)
	}
	if value, ok := zuo.mutation.AddedLong(); ok {
		_spec.AddField(zipcode.FieldLong, field.TypeFloat64, value)
	}
	if zuo.mutation.LongCleared() {
		_spec.ClearField(zipcode.FieldLong, field.TypeFloat64)
	}
	_node = &Zipcode{config: zuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, zuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{zipcode.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	zuo.mutation.done = true
	return _node, nil
}
