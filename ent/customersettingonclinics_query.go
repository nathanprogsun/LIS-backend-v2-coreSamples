// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"coresamples/ent/clinic"
	"coresamples/ent/customer"
	"coresamples/ent/customersettingonclinics"
	"coresamples/ent/predicate"
	"coresamples/ent/setting"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// CustomerSettingOnClinicsQuery is the builder for querying CustomerSettingOnClinics entities.
type CustomerSettingOnClinicsQuery struct {
	config
	ctx          *QueryContext
	order        []customersettingonclinics.OrderOption
	inters       []Interceptor
	predicates   []predicate.CustomerSettingOnClinics
	withCustomer *CustomerQuery
	withClinic   *ClinicQuery
	withSetting  *SettingQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the CustomerSettingOnClinicsQuery builder.
func (csocq *CustomerSettingOnClinicsQuery) Where(ps ...predicate.CustomerSettingOnClinics) *CustomerSettingOnClinicsQuery {
	csocq.predicates = append(csocq.predicates, ps...)
	return csocq
}

// Limit the number of records to be returned by this query.
func (csocq *CustomerSettingOnClinicsQuery) Limit(limit int) *CustomerSettingOnClinicsQuery {
	csocq.ctx.Limit = &limit
	return csocq
}

// Offset to start from.
func (csocq *CustomerSettingOnClinicsQuery) Offset(offset int) *CustomerSettingOnClinicsQuery {
	csocq.ctx.Offset = &offset
	return csocq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (csocq *CustomerSettingOnClinicsQuery) Unique(unique bool) *CustomerSettingOnClinicsQuery {
	csocq.ctx.Unique = &unique
	return csocq
}

// Order specifies how the records should be ordered.
func (csocq *CustomerSettingOnClinicsQuery) Order(o ...customersettingonclinics.OrderOption) *CustomerSettingOnClinicsQuery {
	csocq.order = append(csocq.order, o...)
	return csocq
}

// QueryCustomer chains the current query on the "customer" edge.
func (csocq *CustomerSettingOnClinicsQuery) QueryCustomer() *CustomerQuery {
	query := (&CustomerClient{config: csocq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := csocq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := csocq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(customersettingonclinics.Table, customersettingonclinics.FieldID, selector),
			sqlgraph.To(customer.Table, customer.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, customersettingonclinics.CustomerTable, customersettingonclinics.CustomerColumn),
		)
		fromU = sqlgraph.SetNeighbors(csocq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryClinic chains the current query on the "clinic" edge.
func (csocq *CustomerSettingOnClinicsQuery) QueryClinic() *ClinicQuery {
	query := (&ClinicClient{config: csocq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := csocq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := csocq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(customersettingonclinics.Table, customersettingonclinics.FieldID, selector),
			sqlgraph.To(clinic.Table, clinic.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, customersettingonclinics.ClinicTable, customersettingonclinics.ClinicColumn),
		)
		fromU = sqlgraph.SetNeighbors(csocq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QuerySetting chains the current query on the "setting" edge.
func (csocq *CustomerSettingOnClinicsQuery) QuerySetting() *SettingQuery {
	query := (&SettingClient{config: csocq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := csocq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := csocq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(customersettingonclinics.Table, customersettingonclinics.FieldID, selector),
			sqlgraph.To(setting.Table, setting.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, customersettingonclinics.SettingTable, customersettingonclinics.SettingColumn),
		)
		fromU = sqlgraph.SetNeighbors(csocq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first CustomerSettingOnClinics entity from the query.
// Returns a *NotFoundError when no CustomerSettingOnClinics was found.
func (csocq *CustomerSettingOnClinicsQuery) First(ctx context.Context) (*CustomerSettingOnClinics, error) {
	nodes, err := csocq.Limit(1).All(setContextOp(ctx, csocq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{customersettingonclinics.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (csocq *CustomerSettingOnClinicsQuery) FirstX(ctx context.Context) *CustomerSettingOnClinics {
	node, err := csocq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first CustomerSettingOnClinics ID from the query.
// Returns a *NotFoundError when no CustomerSettingOnClinics ID was found.
func (csocq *CustomerSettingOnClinicsQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = csocq.Limit(1).IDs(setContextOp(ctx, csocq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{customersettingonclinics.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (csocq *CustomerSettingOnClinicsQuery) FirstIDX(ctx context.Context) int {
	id, err := csocq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single CustomerSettingOnClinics entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one CustomerSettingOnClinics entity is found.
// Returns a *NotFoundError when no CustomerSettingOnClinics entities are found.
func (csocq *CustomerSettingOnClinicsQuery) Only(ctx context.Context) (*CustomerSettingOnClinics, error) {
	nodes, err := csocq.Limit(2).All(setContextOp(ctx, csocq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{customersettingonclinics.Label}
	default:
		return nil, &NotSingularError{customersettingonclinics.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (csocq *CustomerSettingOnClinicsQuery) OnlyX(ctx context.Context) *CustomerSettingOnClinics {
	node, err := csocq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only CustomerSettingOnClinics ID in the query.
// Returns a *NotSingularError when more than one CustomerSettingOnClinics ID is found.
// Returns a *NotFoundError when no entities are found.
func (csocq *CustomerSettingOnClinicsQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = csocq.Limit(2).IDs(setContextOp(ctx, csocq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{customersettingonclinics.Label}
	default:
		err = &NotSingularError{customersettingonclinics.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (csocq *CustomerSettingOnClinicsQuery) OnlyIDX(ctx context.Context) int {
	id, err := csocq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of CustomerSettingOnClinicsSlice.
func (csocq *CustomerSettingOnClinicsQuery) All(ctx context.Context) ([]*CustomerSettingOnClinics, error) {
	ctx = setContextOp(ctx, csocq.ctx, "All")
	if err := csocq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*CustomerSettingOnClinics, *CustomerSettingOnClinicsQuery]()
	return withInterceptors[[]*CustomerSettingOnClinics](ctx, csocq, qr, csocq.inters)
}

// AllX is like All, but panics if an error occurs.
func (csocq *CustomerSettingOnClinicsQuery) AllX(ctx context.Context) []*CustomerSettingOnClinics {
	nodes, err := csocq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of CustomerSettingOnClinics IDs.
func (csocq *CustomerSettingOnClinicsQuery) IDs(ctx context.Context) (ids []int, err error) {
	if csocq.ctx.Unique == nil && csocq.path != nil {
		csocq.Unique(true)
	}
	ctx = setContextOp(ctx, csocq.ctx, "IDs")
	if err = csocq.Select(customersettingonclinics.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (csocq *CustomerSettingOnClinicsQuery) IDsX(ctx context.Context) []int {
	ids, err := csocq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (csocq *CustomerSettingOnClinicsQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, csocq.ctx, "Count")
	if err := csocq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, csocq, querierCount[*CustomerSettingOnClinicsQuery](), csocq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (csocq *CustomerSettingOnClinicsQuery) CountX(ctx context.Context) int {
	count, err := csocq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (csocq *CustomerSettingOnClinicsQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, csocq.ctx, "Exist")
	switch _, err := csocq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (csocq *CustomerSettingOnClinicsQuery) ExistX(ctx context.Context) bool {
	exist, err := csocq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the CustomerSettingOnClinicsQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (csocq *CustomerSettingOnClinicsQuery) Clone() *CustomerSettingOnClinicsQuery {
	if csocq == nil {
		return nil
	}
	return &CustomerSettingOnClinicsQuery{
		config:       csocq.config,
		ctx:          csocq.ctx.Clone(),
		order:        append([]customersettingonclinics.OrderOption{}, csocq.order...),
		inters:       append([]Interceptor{}, csocq.inters...),
		predicates:   append([]predicate.CustomerSettingOnClinics{}, csocq.predicates...),
		withCustomer: csocq.withCustomer.Clone(),
		withClinic:   csocq.withClinic.Clone(),
		withSetting:  csocq.withSetting.Clone(),
		// clone intermediate query.
		sql:  csocq.sql.Clone(),
		path: csocq.path,
	}
}

// WithCustomer tells the query-builder to eager-load the nodes that are connected to
// the "customer" edge. The optional arguments are used to configure the query builder of the edge.
func (csocq *CustomerSettingOnClinicsQuery) WithCustomer(opts ...func(*CustomerQuery)) *CustomerSettingOnClinicsQuery {
	query := (&CustomerClient{config: csocq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	csocq.withCustomer = query
	return csocq
}

// WithClinic tells the query-builder to eager-load the nodes that are connected to
// the "clinic" edge. The optional arguments are used to configure the query builder of the edge.
func (csocq *CustomerSettingOnClinicsQuery) WithClinic(opts ...func(*ClinicQuery)) *CustomerSettingOnClinicsQuery {
	query := (&ClinicClient{config: csocq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	csocq.withClinic = query
	return csocq
}

// WithSetting tells the query-builder to eager-load the nodes that are connected to
// the "setting" edge. The optional arguments are used to configure the query builder of the edge.
func (csocq *CustomerSettingOnClinicsQuery) WithSetting(opts ...func(*SettingQuery)) *CustomerSettingOnClinicsQuery {
	query := (&SettingClient{config: csocq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	csocq.withSetting = query
	return csocq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		CustomerID int `json:"customer_id,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.CustomerSettingOnClinics.Query().
//		GroupBy(customersettingonclinics.FieldCustomerID).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (csocq *CustomerSettingOnClinicsQuery) GroupBy(field string, fields ...string) *CustomerSettingOnClinicsGroupBy {
	csocq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &CustomerSettingOnClinicsGroupBy{build: csocq}
	grbuild.flds = &csocq.ctx.Fields
	grbuild.label = customersettingonclinics.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		CustomerID int `json:"customer_id,omitempty"`
//	}
//
//	client.CustomerSettingOnClinics.Query().
//		Select(customersettingonclinics.FieldCustomerID).
//		Scan(ctx, &v)
func (csocq *CustomerSettingOnClinicsQuery) Select(fields ...string) *CustomerSettingOnClinicsSelect {
	csocq.ctx.Fields = append(csocq.ctx.Fields, fields...)
	sbuild := &CustomerSettingOnClinicsSelect{CustomerSettingOnClinicsQuery: csocq}
	sbuild.label = customersettingonclinics.Label
	sbuild.flds, sbuild.scan = &csocq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a CustomerSettingOnClinicsSelect configured with the given aggregations.
func (csocq *CustomerSettingOnClinicsQuery) Aggregate(fns ...AggregateFunc) *CustomerSettingOnClinicsSelect {
	return csocq.Select().Aggregate(fns...)
}

func (csocq *CustomerSettingOnClinicsQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range csocq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, csocq); err != nil {
				return err
			}
		}
	}
	for _, f := range csocq.ctx.Fields {
		if !customersettingonclinics.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if csocq.path != nil {
		prev, err := csocq.path(ctx)
		if err != nil {
			return err
		}
		csocq.sql = prev
	}
	return nil
}

func (csocq *CustomerSettingOnClinicsQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*CustomerSettingOnClinics, error) {
	var (
		nodes       = []*CustomerSettingOnClinics{}
		_spec       = csocq.querySpec()
		loadedTypes = [3]bool{
			csocq.withCustomer != nil,
			csocq.withClinic != nil,
			csocq.withSetting != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*CustomerSettingOnClinics).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &CustomerSettingOnClinics{config: csocq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, csocq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := csocq.withCustomer; query != nil {
		if err := csocq.loadCustomer(ctx, query, nodes, nil,
			func(n *CustomerSettingOnClinics, e *Customer) { n.Edges.Customer = e }); err != nil {
			return nil, err
		}
	}
	if query := csocq.withClinic; query != nil {
		if err := csocq.loadClinic(ctx, query, nodes, nil,
			func(n *CustomerSettingOnClinics, e *Clinic) { n.Edges.Clinic = e }); err != nil {
			return nil, err
		}
	}
	if query := csocq.withSetting; query != nil {
		if err := csocq.loadSetting(ctx, query, nodes, nil,
			func(n *CustomerSettingOnClinics, e *Setting) { n.Edges.Setting = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (csocq *CustomerSettingOnClinicsQuery) loadCustomer(ctx context.Context, query *CustomerQuery, nodes []*CustomerSettingOnClinics, init func(*CustomerSettingOnClinics), assign func(*CustomerSettingOnClinics, *Customer)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*CustomerSettingOnClinics)
	for i := range nodes {
		fk := nodes[i].CustomerID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(customer.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "customer_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (csocq *CustomerSettingOnClinicsQuery) loadClinic(ctx context.Context, query *ClinicQuery, nodes []*CustomerSettingOnClinics, init func(*CustomerSettingOnClinics), assign func(*CustomerSettingOnClinics, *Clinic)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*CustomerSettingOnClinics)
	for i := range nodes {
		fk := nodes[i].ClinicID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(clinic.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "clinic_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (csocq *CustomerSettingOnClinicsQuery) loadSetting(ctx context.Context, query *SettingQuery, nodes []*CustomerSettingOnClinics, init func(*CustomerSettingOnClinics), assign func(*CustomerSettingOnClinics, *Setting)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*CustomerSettingOnClinics)
	for i := range nodes {
		fk := nodes[i].SettingID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(setting.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "setting_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (csocq *CustomerSettingOnClinicsQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := csocq.querySpec()
	_spec.Node.Columns = csocq.ctx.Fields
	if len(csocq.ctx.Fields) > 0 {
		_spec.Unique = csocq.ctx.Unique != nil && *csocq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, csocq.driver, _spec)
}

func (csocq *CustomerSettingOnClinicsQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(customersettingonclinics.Table, customersettingonclinics.Columns, sqlgraph.NewFieldSpec(customersettingonclinics.FieldID, field.TypeInt))
	_spec.From = csocq.sql
	if unique := csocq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if csocq.path != nil {
		_spec.Unique = true
	}
	if fields := csocq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, customersettingonclinics.FieldID)
		for i := range fields {
			if fields[i] != customersettingonclinics.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
		if csocq.withCustomer != nil {
			_spec.Node.AddColumnOnce(customersettingonclinics.FieldCustomerID)
		}
		if csocq.withClinic != nil {
			_spec.Node.AddColumnOnce(customersettingonclinics.FieldClinicID)
		}
		if csocq.withSetting != nil {
			_spec.Node.AddColumnOnce(customersettingonclinics.FieldSettingID)
		}
	}
	if ps := csocq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := csocq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := csocq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := csocq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (csocq *CustomerSettingOnClinicsQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(csocq.driver.Dialect())
	t1 := builder.Table(customersettingonclinics.Table)
	columns := csocq.ctx.Fields
	if len(columns) == 0 {
		columns = customersettingonclinics.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if csocq.sql != nil {
		selector = csocq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if csocq.ctx.Unique != nil && *csocq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range csocq.predicates {
		p(selector)
	}
	for _, p := range csocq.order {
		p(selector)
	}
	if offset := csocq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := csocq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// CustomerSettingOnClinicsGroupBy is the group-by builder for CustomerSettingOnClinics entities.
type CustomerSettingOnClinicsGroupBy struct {
	selector
	build *CustomerSettingOnClinicsQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (csocgb *CustomerSettingOnClinicsGroupBy) Aggregate(fns ...AggregateFunc) *CustomerSettingOnClinicsGroupBy {
	csocgb.fns = append(csocgb.fns, fns...)
	return csocgb
}

// Scan applies the selector query and scans the result into the given value.
func (csocgb *CustomerSettingOnClinicsGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, csocgb.build.ctx, "GroupBy")
	if err := csocgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*CustomerSettingOnClinicsQuery, *CustomerSettingOnClinicsGroupBy](ctx, csocgb.build, csocgb, csocgb.build.inters, v)
}

func (csocgb *CustomerSettingOnClinicsGroupBy) sqlScan(ctx context.Context, root *CustomerSettingOnClinicsQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(csocgb.fns))
	for _, fn := range csocgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*csocgb.flds)+len(csocgb.fns))
		for _, f := range *csocgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*csocgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := csocgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// CustomerSettingOnClinicsSelect is the builder for selecting fields of CustomerSettingOnClinics entities.
type CustomerSettingOnClinicsSelect struct {
	*CustomerSettingOnClinicsQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (csocs *CustomerSettingOnClinicsSelect) Aggregate(fns ...AggregateFunc) *CustomerSettingOnClinicsSelect {
	csocs.fns = append(csocs.fns, fns...)
	return csocs
}

// Scan applies the selector query and scans the result into the given value.
func (csocs *CustomerSettingOnClinicsSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, csocs.ctx, "Select")
	if err := csocs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*CustomerSettingOnClinicsQuery, *CustomerSettingOnClinicsSelect](ctx, csocs.CustomerSettingOnClinicsQuery, csocs, csocs.inters, v)
}

func (csocs *CustomerSettingOnClinicsSelect) sqlScan(ctx context.Context, root *CustomerSettingOnClinicsQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(csocs.fns))
	for _, fn := range csocs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*csocs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := csocs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
