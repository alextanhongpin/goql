package goql

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	FilterTag = "q"
	SortTag   = "sort"
	And       = "and"
	Or        = "or"
)

var (
	ErrMultipleOperator = errors.New("goql: multiple op")
	ErrUnknownOperator  = errors.New("goql: unknown op")
	ErrInvalidOp        = errors.New("goql: invalid op")
	ErrInvalidArray     = errors.New("goql: invalid array")
	ErrUnknownField     = errors.New("goql: unknown field")
	ErrUnknownParser    = errors.New("goql: unknown parser")
)

type Filter struct {
	Sort []Order
	And  []FieldSet
	Or   []FieldSet
}

type FieldSet struct {
	Tag      *Tag
	Name     string
	Value    any
	RawValue string
	Op       string

	Or  []FieldSet
	And []FieldSet
}

type Decoder[T any] struct {
	tagByField map[string]*Tag
	parsers    map[string]ParserFn
	tag        string
}

func NewDecoder[T any]() (*Decoder[T], error) {
	var t T

	parserByType := NewParsers()
	tagByField, err := ParseStruct(t, FilterTag, SortTag)
	if err != nil {
		return nil, err
	}

	return &Decoder[T]{
		tagByField: tagByField,
		tag:        FilterTag,
		parsers:    parserByType,
	}, nil
}

func (d *Decoder[T]) SetStructTag(tag string) error {
	if tag == "" {
		panic("tag cannot be empty")
	}

	var t T
	tagByField, err := ParseStruct(t, tag, SortTag)
	if err != nil {
		return err
	}

	d.tag = tag
	d.tagByField = tagByField

	return nil
}

func (d *Decoder[T]) SetParsers(parserByType map[string]ParserFn) {
	d.parsers = parserByType
}

func (d *Decoder[T]) SetOps(field string, ops Op) error {
	if _, ok := d.tagByField[field]; !ok {
		return fmt.Errorf("%w: %s", ErrUnknownField, field)
	}

	d.tagByField[field].Ops = ops

	return nil
}

func (d *Decoder[T]) Decode(values url.Values) (*Filter, error) {
	base, err := d.decodeFields(values)
	if err != nil {
		return nil, err
	}

	ands, err := d.decodeConjunction(OpAnd, values)
	if err != nil {
		return nil, err
	}

	ands = append(base, ands...)

	ors, err := d.decodeConjunction(OpOr, values)
	if err != nil {
		return nil, err
	}

	sorts, err := d.parseSort(values)
	if err != nil {
		return nil, err
	}

	return &Filter{
		Sort: sorts,
		And:  ands,
		Or:   ors,
	}, nil
}

func (d *Decoder[T]) parseSort(values url.Values) ([]Order, error) {
	sort, err := ParseOrder(values.Get("sort_by"))
	if err != nil {
		return nil, err
	}

	validSortByField := make(map[string]bool)
	for field, tag := range d.tagByField {
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

func (d *Decoder[T]) decodeFields(values url.Values) ([]FieldSet, error) {
	ands := make([]FieldSet, 0, len(values))
	cache := make(map[string]bool)
	queries := ParseQuery(values)

	for _, query := range queries {
		field, op, value := query.Field, query.Op, query.Value

		tag, ok := d.tagByField[field]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownField, field)
		}

		hasOp := tag.Ops.Has(op)
		if !hasOp {
			return nil, fmt.Errorf("%w: %s", ErrUnknownOperator, query)
		}

		if OpsNull.Has(op) && !sqlIs(value) {
			// OpIs/OpIsNot must have value: true, false, unknown or null.
			return nil, fmt.Errorf("%w: %s", ErrInvalidOp, query)
		}

		cacheKey := fmt.Sprintf("%s:%s", field, op)
		if cache[cacheKey] {
			return nil, fmt.Errorf("%w: %q.%q", ErrMultipleOperator, field, strings.ToLower(op.String()))
		}
		cache[cacheKey] = true

		parser, ok := d.parsers[tag.Type.Name]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownParser, tag.Type.Name)
		}

		fs := FieldSet{
			Tag:      tag,
			Name:     field,
			Op:       op.String(),
			RawValue: value,
		}

		switch {
		case OpsIn.Has(op), tag.Type.Array:
			value, ok = Unquote(value, '{', '}')
			if !ok {
				return nil, fmt.Errorf("%w: missing parantheses: %s", ErrInvalidArray, value)
			}

			// Strings may contain commas, which interferes with the splitting.
			vals := SplitCsv(value)

			res, err := Map(vals, parser)
			if err != nil {
				return nil, err
			}

			fs.Value = res
		default:
			res, err := parser(value)
			if err != nil {
				return nil, err
			}
			fs.Value = res
		}

		ands = append(ands, fs)
	}

	return ands, nil
}

func (d *Decoder[T]) decodeConjunction(conj Op, values url.Values) ([]FieldSet, error) {
	switch conj {
	case OpAnd, OpOr:
	default:
		panic("goql: invalid conj")
	}

	conjs := make([]FieldSet, 0, len(values))
	for _, v := range values[strings.ToLower(conj.String())] {
		value, ok := Unquote(v, '(', ')')
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrInvalidOp, conj)
		}

		values := SplitOutsideBrackets(value)
		fmt.Println("OUTSIDE", values, len(values))

		uvals := make(url.Values)
		avals := make(url.Values)
		ovals := make(url.Values)
		for _, value := range values {
			field, opv := Split2(value, ".")
			switch field {
			case Or:
				ovals.Add(Or, opv)
			case And:
				avals.Add(And, opv)
			default:
				uvals[field] = append(uvals[field], opv)
			}
		}

		ors, err := d.decodeConjunction(OpOr, ovals)
		if err != nil {
			return nil, err
		}

		ands, err := d.decodeConjunction(OpAnd, avals)
		if err != nil {
			return nil, err
		}

		innerConj, err := d.decodeFields(uvals)
		if err != nil {
			return nil, err
		}

		fs := FieldSet{
			Name:     conj.String(),
			Op:       conj.String(),
			Value:    value,
			RawValue: v,
			And:      ands,
			Or:       ors,
		}

		switch conj {
		case OpAnd:
			fs.And = append(fs.And, innerConj...)
		case OpOr:
			fs.Or = append(fs.Or, innerConj...)
		}
		conjs = append(conjs, fs)
	}

	return conjs, nil
}
