package goql

import (
	"errors"
	"fmt"
	"net/url"
)

const (
	FilterTag = "q"
	SortTag   = "sort"
)

var (
	ErrUnknownOperator = errors.New("goql: unknown op")
	ErrInvalidOp       = errors.New("goql: invalid op")
	ErrUnknownField    = errors.New("goql: unknown field")
	ErrUnknownParser   = errors.New("goql: unknown parser")
	ErrTooManyValues   = errors.New("goql: too many values")
)

type Filter struct {
	Sort []Order
	And  []FieldSet
	Or   []FieldSet
}

type FieldSet struct {
	Tag    *Tag
	Name   string
	Value  any
	Values []string
	Op     string

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
	sort, err := ParseOrder(values["sort_by"])
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

	queries, err := ParseQuery(values, OpAnd.String(), OpOr.String(), "sort_by")
	if err != nil {
		return nil, err
	}

	for _, query := range queries {
		field, op, values := query.Field, query.Op, query.Values

		tag, ok := d.tagByField[field]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownField, field)
		}

		hasOp := tag.Ops.Has(op)
		if !hasOp {
			return nil, fmt.Errorf("%w: %s", ErrUnknownOperator, query)
		}

		if OpsNull.Has(op) && !sqlIs(values[0]) {
			// OpIs/OpIsNot must have value: true, false, unknown or null.
			return nil, fmt.Errorf("%w: %s", ErrInvalidOp, query)
		}

		parser, ok := d.parsers[tag.Type.Name]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownParser, tag.Type.Name)
		}

		fs := FieldSet{
			Tag:    tag,
			Name:   field,
			Op:     op.String(),
			Values: values,
		}

		switch {
		case OpsIn.Has(op), tag.Type.Array:
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
	for _, v := range values[conj.String()] {
		value, ok := Unquote(v, '(', ')')
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrInvalidOp, conj)
		}

		values := SplitOutsideBrackets(value)

		uvals := make(url.Values)
		avals := make(url.Values)
		ovals := make(url.Values)
		for _, value := range values {
			field, opv := Split2(value, ".")
			switch field {
			case OpOr.String():
				ovals.Add(OpOr.String(), opv)
			case OpAnd.String():
				avals.Add(OpAnd.String(), opv)
			default:
				rawOp, value := Split2(opv, ":")
				op, ok := ParseOp(rawOp)
				if !ok {
					return nil, fmt.Errorf("%w: %s", ErrUnknownOperator, rawOp)
				}

				uvals.Add(fmt.Sprintf("%s.%s", field, op), value)
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
			Name:   conj.String(),
			Op:     conj.String(),
			Value:  value,
			Values: values,
			And:    ands,
			Or:     ors,
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
