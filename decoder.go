package goql

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
)

const (
	// Tag name for struct parsing, customizable.
	TagFilter = "q"
	TagSort   = "sort"

	// Reserved query string fields, customizable, since 'sort_by' or 'limit'
	// could be a valid field name.
	QuerySort   = "sort_by"
	QueryLimit  = "limit"
	QueryOffset = "offset"
	QueryAnd    = "and"
	QueryOr     = "or"

	// Pagination limit.
	LimitMin = 1
	LimitMax = 20
)

var (
	ErrUnknownOperator    = errors.New("goql: unknown op")
	ErrInvalidOp          = errors.New("goql: invalid op")
	ErrUnknownField       = errors.New("goql: unknown field")
	ErrUnknownParser      = errors.New("goql: unknown parser")
	ErrInvalidConjunction = errors.New("goql: invalid conjunction")
	ErrBadValue           = errors.New("goql: bad value")
	ErrTooManyValues      = errors.New("goql: too many values")
)

type Filter struct {
	Sort   []Order
	And    []FieldSet
	Or     []FieldSet
	Limit  *int
	Offset *int
}

type FieldSet struct {
	Tag    *Tag
	Name   string
	Value  any
	Values []string
	Op     Op

	Or  []FieldSet
	And []FieldSet
}

func (f FieldSet) String() string {
	return fmt.Sprintf("%v %v %#v", f.Name, f.Op, f.Value)
}

type Decoder[T any] struct {
	tags        map[string]*Tag
	parsers     map[string]ParserFn
	sortTag     string
	filterTag   string
	limitMin    int
	limitMax    int
	querySort   string
	queryLimit  string
	queryOffset string
}

func NewDecoder[T any]() *Decoder[T] {
	var t T

	parsers := NewParsers()
	tags, err := ParseStruct(t, TagFilter, TagSort)
	if err != nil {
		panic(err)
	}

	return &Decoder[T]{
		tags:        tags,
		parsers:     parsers,
		sortTag:     TagSort,
		filterTag:   TagFilter,
		limitMin:    LimitMin,
		limitMax:    LimitMax,
		querySort:   QuerySort,
		queryLimit:  QueryLimit,
		queryOffset: QueryOffset,
	}
}

// Validate checks if the parser exists for all the inferred types. This is
// called internally before decode is called.
func (d *Decoder[T]) Validate() error {
	for _, tag := range d.tags {
		if _, ok := d.parsers[tag.Type.Name]; !ok {
			return fmt.Errorf("%w: missing parser for type %q", ErrUnknownParser, tag.Type.Name)
		}
	}

	return nil
}

func (d *Decoder[T]) SetFilterTag(filterTag string) *Decoder[T] {
	if filterTag == "" {
		panic("goql: filter tag cannot be empty")
	}

	var t T
	tags, err := ParseStruct(t, filterTag, d.sortTag)

	if err != nil {
		panic(err)
	}

	d.filterTag = filterTag
	d.tags = tags

	return d
}

func (d *Decoder[T]) SetSortTag(sortTag string) *Decoder[T] {
	if sortTag == "" {
		panic("goql: sort tag cannot be empty")
	}

	var t T
	tags, err := ParseStruct(t, d.filterTag, sortTag)
	if err != nil {
		panic(err)
	}

	d.sortTag = sortTag
	d.tags = tags

	return d
}

func (d *Decoder[T]) SetLimitRange(min, max int) *Decoder[T] {
	if min == 0 || max == 0 {
		panic("goql: limit and offset cannot be 0")
	}

	d.limitMin = min
	d.limitMax = max

	return d
}

func (d *Decoder[T]) SetParsers(parsers map[string]ParserFn) *Decoder[T] {
	if len(parsers) == 0 {
		panic("goql: no parsers specified")
	}

	d.parsers = parsers

	return d
}

func (d *Decoder[T]) SetParser(name string, parserFn ParserFn) *Decoder[T] {
	d.parsers[name] = parserFn

	return d
}

func (d *Decoder[T]) SetOps(field string, ops Op) *Decoder[T] {
	if field == "" {
		panic("goql: set ops field cannot be empty")
	}

	if _, ok := d.tags[field]; !ok {
		panic(fmt.Errorf("%w: %q", ErrUnknownField, field))
	}

	if !ops.Valid() {
		panic(fmt.Errorf("%w: %q", ErrInvalidOp, field))
	}

	d.tags[field].Ops = ops

	return d
}

func (d *Decoder[T]) SetQuerySortName(name string) *Decoder[T] {
	if name == "" {
		panic("goql: query sort name cannot be empty")
	}

	d.querySort = name

	return d
}

func (d *Decoder[T]) SetQueryLimitName(name string) *Decoder[T] {
	if name == "" {
		panic("goql: query limit name cannot be empty")
	}

	d.queryLimit = name

	return d
}

func (d *Decoder[T]) SetQueryOffsetName(name string) *Decoder[T] {
	if name == "" {
		panic("goql: query offset name cannot be empty")
	}

	d.queryOffset = name

	return d
}

func (d *Decoder[T]) Decode(u url.Values) (*Filter, error) {
	if err := d.Validate(); err != nil {
		return nil, err
	}

	limit, offset, err := d.parseLimit(u)
	if err != nil {
		return nil, err
	}

	ands, ors, err := d.parseFilter(u)
	if err != nil {
		return nil, err
	}

	sorts, err := d.parseSort(u)
	if err != nil {
		return nil, err
	}

	return &Filter{
		Sort:   sorts,
		And:    ands,
		Or:     ors,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (d *Decoder[T]) parseLimit(u url.Values) (limit, offset *int, err error) {
	if v, ok := u[d.queryLimit]; ok && len(v) > 0 {
		var n int
		n, err = strconv.Atoi(v[0])
		if err != nil {
			return
		}

		if n < d.limitMin {
			n = d.limitMin
		}

		if n > d.limitMax {
			n = d.limitMax
		}

		limit = &n
	}

	if v, ok := u[d.queryOffset]; ok && len(v) > 0 {
		var n int
		n, err = strconv.Atoi(v[0])
		if err != nil {
			return
		}

		if n < 0 {
			n = 0
		}

		offset = &n
	}

	return
}

func (d *Decoder[T]) reservedKeys() []string {
	return []string{QueryAnd, QueryOr, d.querySort, d.queryLimit, d.queryOffset}
}

func (d *Decoder[T]) parseFilter(values url.Values) (ands, ors []FieldSet, err error) {
	baseValues := FilterValues(values, d.reservedKeys()...)

	// Base values are the same as AND values.
	// `name.eq=john` and `and=name.eq:john` is equivalent.
	// We merge them in order to remove duplicate values.
	andValues := make([]string, 0, len(baseValues))

	for key, values := range baseValues {
		for _, val := range values {
			andValues = append(andValues, fmt.Sprintf("%s:%s", key, val))
		}
	}

	andValues = append(andValues, values[QueryAnd]...)

	ands, err = d.decodeConjunction(OpAnd, andValues)
	if err != nil {
		return
	}

	ors, err = d.decodeConjunction(OpOr, values[QueryOr])
	if err != nil {
		return
	}

	return
}

func (d *Decoder[T]) parseSort(values url.Values) ([]Order, error) {
	sort, err := ParseOrder(values[d.querySort])
	if err != nil {
		return nil, err
	}

	validSortByField := make(map[string]bool)
	for field, tag := range d.tags {
		validSortByField[field] = tag.Sort
	}

	sorts := make([]Order, 0, len(sort))
	for _, s := range sort {
		if validSortByField[s.Field] {
			sorts = append(sorts, s)
		}
	}

	return sorts, nil
}

func (d *Decoder[T]) decodeField(query Query) (*FieldSet, error) {
	field, op, values := query.Field, query.Op, query.Values

	tag, ok := d.tags[field]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownField, field)
	}

	if ok := tag.Ops.Has(op); !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownOperator, query)
	}

	parser, ok := d.parsers[tag.Type.Name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownParser, tag.Type.Name)
	}

	fs := FieldSet{
		Tag:    tag,
		Name:   field,
		Op:     op,
		Values: values,
	}

	switch {
	case OpsMany.Has(op), tag.Type.Array:
		res, err := Map(values, parser)
		if err != nil {
			return nil, err
		}

		fs.Value = res
	default:
		if len(values) > 1 {
			return nil, fmt.Errorf("%w: %s", ErrTooManyValues, query)
		}

		value, _ := Unquote(values[0], '"', '"')
		res, err := parser(value)
		if err != nil {
			return nil, err
		}

		fs.Value = res
	}

	return &fs, nil
}

func (d *Decoder[T]) decodeFields(values url.Values) ([]FieldSet, error) {
	res := make([]FieldSet, 0, len(values))

	values = FilterValues(values, d.reservedKeys()...)
	queries, err := ParseQuery(values)
	if err != nil {
		return nil, err
	}

	for _, query := range queries {
		fs, err := d.decodeField(query)
		if err != nil {
			return nil, err
		}

		res = append(res, *fs)
	}

	return res, nil
}

func (d *Decoder[T]) decodeConjunction(conj Op, values []string) ([]FieldSet, error) {
	switch conj {
	case OpAnd, OpOr:
	default:
		panic("goql: invalid conj")
	}

	values = Unique(values)
	sort.Strings(values)

	conjs := make([]FieldSet, 0, len(values))

	uvals := make(url.Values)

	for _, value := range values {
		field, opv := Split2(value, ".")

		switch field {
		case OpOr.String():
			vl, ok := Unquote(opv, '(', ')')
			if !ok {
				return nil, fmt.Errorf("%w: %s", ErrInvalidConjunction, opv)
			}

			vals := SplitOutsideBrackets(vl)
			ors, err := d.decodeConjunction(OpOr, vals)
			if err != nil {
				return nil, err
			}

			fs := FieldSet{
				Name:   conj.String(),
				Op:     OpOr,
				Values: vals,
				Or:     ors,
			}

			conjs = append(conjs, fs)

		case OpAnd.String():
			vl, ok := Unquote(opv, '(', ')')
			if !ok {
				return nil, fmt.Errorf("%w: %s", ErrInvalidConjunction, opv)
			}

			vals := SplitOutsideBrackets(vl)

			ands, err := d.decodeConjunction(OpAnd, vals)
			if err != nil {
				return nil, err
			}

			fs := FieldSet{
				Name:   conj.String(),
				Op:     OpAnd,
				Values: vals,
				And:    ands,
			}

			conjs = append(conjs, fs)

		default:
			// The `AND` may contain `in` operators, e.g.
			// and=name.in:alice&and=name.in:bob
			// We need to combine them to `name.in=[]string{alice, bob}` before
			// parsing.
			k, v := Split2(value, ":")
			uvals.Add(k, v)
		}
	}

	innerConjs, err := d.decodeFields(uvals)
	if err != nil {
		return nil, err
	}

	conjs = append(innerConjs, conjs...)

	return conjs, nil
}
