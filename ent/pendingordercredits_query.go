// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"coresamples/ent/pendingordercredits"
	"coresamples/ent/predicate"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// PendingOrderCreditsQuery is the builder for querying PendingOrderCredits entities.
type PendingOrderCreditsQuery struct {
	config
	ctx        *QueryContext
	order      []pendingordercredits.OrderOption
	inters     []Interceptor
	predicates []predicate.PendingOrderCredits
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the PendingOrderCreditsQuery builder.
func (pocq *PendingOrderCreditsQuery) Where(ps ...predicate.PendingOrderCredits) *PendingOrderCreditsQuery {
	pocq.predicates = append(pocq.predicates, ps...)
	return pocq
}

// Limit the number of records to be returned by this query.
func (pocq *PendingOrderCreditsQuery) Limit(limit int) *PendingOrderCreditsQuery {
	pocq.ctx.Limit = &limit
	return pocq
}

// Offset to start from.
func (pocq *PendingOrderCreditsQuery) Offset(offset int) *PendingOrderCreditsQuery {
	pocq.ctx.Offset = &offset
	return pocq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (pocq *PendingOrderCreditsQuery) Unique(unique bool) *PendingOrderCreditsQuery {
	pocq.ctx.Unique = &unique
	return pocq
}

// Order specifies how the records should be ordered.
func (pocq *PendingOrderCreditsQuery) Order(o ...pendingordercredits.OrderOption) *PendingOrderCreditsQuery {
	pocq.order = append(pocq.order, o...)
	return pocq
}

// First returns the first PendingOrderCredits entity from the query.
// Returns a *NotFoundError when no PendingOrderCredits was found.
func (pocq *PendingOrderCreditsQuery) First(ctx context.Context) (*PendingOrderCredits, error) {
	nodes, err := pocq.Limit(1).All(setContextOp(ctx, pocq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{pendingordercredits.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (pocq *PendingOrderCreditsQuery) FirstX(ctx context.Context) *PendingOrderCredits {
	node, err := pocq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first PendingOrderCredits ID from the query.
// Returns a *NotFoundError when no PendingOrderCredits ID was found.
func (pocq *PendingOrderCreditsQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = pocq.Limit(1).IDs(setContextOp(ctx, pocq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{pendingordercredits.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (pocq *PendingOrderCreditsQuery) FirstIDX(ctx context.Context) int {
	id, err := pocq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single PendingOrderCredits entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one PendingOrderCredits entity is found.
// Returns a *NotFoundError when no PendingOrderCredits entities are found.
func (pocq *PendingOrderCreditsQuery) Only(ctx context.Context) (*PendingOrderCredits, error) {
	nodes, err := pocq.Limit(2).All(setContextOp(ctx, pocq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{pendingordercredits.Label}
	default:
		return nil, &NotSingularError{pendingordercredits.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (pocq *PendingOrderCreditsQuery) OnlyX(ctx context.Context) *PendingOrderCredits {
	node, err := pocq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only PendingOrderCredits ID in the query.
// Returns a *NotSingularError when more than one PendingOrderCredits ID is found.
// Returns a *NotFoundError when no entities are found.
func (pocq *PendingOrderCreditsQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = pocq.Limit(2).IDs(setContextOp(ctx, pocq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{pendingordercredits.Label}
	default:
		err = &NotSingularError{pendingordercredits.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (pocq *PendingOrderCreditsQuery) OnlyIDX(ctx context.Context) int {
	id, err := pocq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of PendingOrderCreditsSlice.
func (pocq *PendingOrderCreditsQuery) All(ctx context.Context) ([]*PendingOrderCredits, error) {
	ctx = setContextOp(ctx, pocq.ctx, "All")
	if err := pocq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*PendingOrderCredits, *PendingOrderCreditsQuery]()
	return withInterceptors[[]*PendingOrderCredits](ctx, pocq, qr, pocq.inters)
}

// AllX is like All, but panics if an error occurs.
func (pocq *PendingOrderCreditsQuery) AllX(ctx context.Context) []*PendingOrderCredits {
	nodes, err := pocq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of PendingOrderCredits IDs.
func (pocq *PendingOrderCreditsQuery) IDs(ctx context.Context) (ids []int, err error) {
	if pocq.ctx.Unique == nil && pocq.path != nil {
		pocq.Unique(true)
	}
	ctx = setContextOp(ctx, pocq.ctx, "IDs")
	if err = pocq.Select(pendingordercredits.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (pocq *PendingOrderCreditsQuery) IDsX(ctx context.Context) []int {
	ids, err := pocq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (pocq *PendingOrderCreditsQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, pocq.ctx, "Count")
	if err := pocq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, pocq, querierCount[*PendingOrderCreditsQuery](), pocq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (pocq *PendingOrderCreditsQuery) CountX(ctx context.Context) int {
	count, err := pocq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (pocq *PendingOrderCreditsQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, pocq.ctx, "Exist")
	switch _, err := pocq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (pocq *PendingOrderCreditsQuery) ExistX(ctx context.Context) bool {
	exist, err := pocq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the PendingOrderCreditsQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (pocq *PendingOrderCreditsQuery) Clone() *PendingOrderCreditsQuery {
	if pocq == nil {
		return nil
	}
	return &PendingOrderCreditsQuery{
		config:     pocq.config,
		ctx:        pocq.ctx.Clone(),
		order:      append([]pendingordercredits.OrderOption{}, pocq.order...),
		inters:     append([]Interceptor{}, pocq.inters...),
		predicates: append([]predicate.PendingOrderCredits{}, pocq.predicates...),
		// clone intermediate query.
		sql:  pocq.sql.Clone(),
		path: pocq.path,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		OrderID int64 `json:"order_id,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.PendingOrderCredits.Query().
//		GroupBy(pendingordercredits.FieldOrderID).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (pocq *PendingOrderCreditsQuery) GroupBy(field string, fields ...string) *PendingOrderCreditsGroupBy {
	pocq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &PendingOrderCreditsGroupBy{build: pocq}
	grbuild.flds = &pocq.ctx.Fields
	grbuild.label = pendingordercredits.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		OrderID int64 `json:"order_id,omitempty"`
//	}
//
//	client.PendingOrderCredits.Query().
//		Select(pendingordercredits.FieldOrderID).
//		Scan(ctx, &v)
func (pocq *PendingOrderCreditsQuery) Select(fields ...string) *PendingOrderCreditsSelect {
	pocq.ctx.Fields = append(pocq.ctx.Fields, fields...)
	sbuild := &PendingOrderCreditsSelect{PendingOrderCreditsQuery: pocq}
	sbuild.label = pendingordercredits.Label
	sbuild.flds, sbuild.scan = &pocq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a PendingOrderCreditsSelect configured with the given aggregations.
func (pocq *PendingOrderCreditsQuery) Aggregate(fns ...AggregateFunc) *PendingOrderCreditsSelect {
	return pocq.Select().Aggregate(fns...)
}

func (pocq *PendingOrderCreditsQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range pocq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, pocq); err != nil {
				return err
			}
		}
	}
	for _, f := range pocq.ctx.Fields {
		if !pendingordercredits.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if pocq.path != nil {
		prev, err := pocq.path(ctx)
		if err != nil {
			return err
		}
		pocq.sql = prev
	}
	return nil
}

func (pocq *PendingOrderCreditsQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*PendingOrderCredits, error) {
	var (
		nodes = []*PendingOrderCredits{}
		_spec = pocq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*PendingOrderCredits).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &PendingOrderCredits{config: pocq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, pocq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (pocq *PendingOrderCreditsQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := pocq.querySpec()
	_spec.Node.Columns = pocq.ctx.Fields
	if len(pocq.ctx.Fields) > 0 {
		_spec.Unique = pocq.ctx.Unique != nil && *pocq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, pocq.driver, _spec)
}

func (pocq *PendingOrderCreditsQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(pendingordercredits.Table, pendingordercredits.Columns, sqlgraph.NewFieldSpec(pendingordercredits.FieldID, field.TypeInt))
	_spec.From = pocq.sql
	if unique := pocq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if pocq.path != nil {
		_spec.Unique = true
	}
	if fields := pocq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, pendingordercredits.FieldID)
		for i := range fields {
			if fields[i] != pendingordercredits.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := pocq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := pocq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := pocq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := pocq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (pocq *PendingOrderCreditsQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(pocq.driver.Dialect())
	t1 := builder.Table(pendingordercredits.Table)
	columns := pocq.ctx.Fields
	if len(columns) == 0 {
		columns = pendingordercredits.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if pocq.sql != nil {
		selector = pocq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if pocq.ctx.Unique != nil && *pocq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range pocq.predicates {
		p(selector)
	}
	for _, p := range pocq.order {
		p(selector)
	}
	if offset := pocq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := pocq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// PendingOrderCreditsGroupBy is the group-by builder for PendingOrderCredits entities.
type PendingOrderCreditsGroupBy struct {
	selector
	build *PendingOrderCreditsQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (pocgb *PendingOrderCreditsGroupBy) Aggregate(fns ...AggregateFunc) *PendingOrderCreditsGroupBy {
	pocgb.fns = append(pocgb.fns, fns...)
	return pocgb
}

// Scan applies the selector query and scans the result into the given value.
func (pocgb *PendingOrderCreditsGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, pocgb.build.ctx, "GroupBy")
	if err := pocgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*PendingOrderCreditsQuery, *PendingOrderCreditsGroupBy](ctx, pocgb.build, pocgb, pocgb.build.inters, v)
}

func (pocgb *PendingOrderCreditsGroupBy) sqlScan(ctx context.Context, root *PendingOrderCreditsQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(pocgb.fns))
	for _, fn := range pocgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*pocgb.flds)+len(pocgb.fns))
		for _, f := range *pocgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*pocgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := pocgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// PendingOrderCreditsSelect is the builder for selecting fields of PendingOrderCredits entities.
type PendingOrderCreditsSelect struct {
	*PendingOrderCreditsQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (pocs *PendingOrderCreditsSelect) Aggregate(fns ...AggregateFunc) *PendingOrderCreditsSelect {
	pocs.fns = append(pocs.fns, fns...)
	return pocs
}

// Scan applies the selector query and scans the result into the given value.
func (pocs *PendingOrderCreditsSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, pocs.ctx, "Select")
	if err := pocs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*PendingOrderCreditsQuery, *PendingOrderCreditsSelect](ctx, pocs.PendingOrderCreditsQuery, pocs, pocs.inters, v)
}

func (pocs *PendingOrderCreditsSelect) sqlScan(ctx context.Context, root *PendingOrderCreditsQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(pocs.fns))
	for _, fn := range pocs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*pocs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := pocs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
