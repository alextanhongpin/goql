package goql

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	StructTag    = "q"
	ValSeparator = ":"
)

var (
	ErrMultipleOperator = errors.New("goql: multiple op")
	ErrUnknownOperator  = errors.New("goql: unknown op")
	ErrInvalidIs        = errors.New("goql: 'is' must be followed by {true, false, null, unknown}")
	ErrInvalidArray     = errors.New("goql: invalid array")
	ErrUnknownField     = errors.New("goql: unknown field")
	ErrUnknownParser    = errors.New("goql: unknown parser")
)

type FieldSet struct {
	Tag      *Tag
	Name     string
	Value    any
	RawValue string
	Op       string
}

type Decoder[T any] struct {
	tagByField map[string]*Tag
	parsers    map[string]ParserFn
	tag        string
}

func NewDecoder[T any]() (*Decoder[T], error) {
	var t T

	parserByType := NewParsers()
	tagByField, err := ParseStruct(t, StructTag)
	if err != nil {
		return nil, err
	}

	return &Decoder[T]{
		tagByField: tagByField,
		tag:        StructTag,
		parsers:    parserByType,
	}, nil
}

func (d *Decoder[T]) SetStructTag(tag string) error {
	if tag == "" {
		panic("tag cannot be empty")
	}

	var t T
	tagByField, err := ParseStruct(t, tag)
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

func (d *Decoder[T]) Decode(values url.Values) ([]FieldSet, error) {
	return Decode(d.tagByField, d.parsers, values)
}

func Decode(tagByField map[string]*Tag, parsers map[string]ParserFn, values url.Values) ([]FieldSet, error) {
	cache := make(map[string]bool)

	var sets []FieldSet

	queries := ParseQuery(values)
	for _, query := range queries {
		field, op, value := query.Field, query.Op, query.Value

		tag, ok := tagByField[field]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownField, field)
		}

		hasOp := tag.Ops.Has(op)
		if !hasOp {
			return nil, fmt.Errorf("%w: %s", ErrUnknownOperator, query)
		}

		if OpsNull.Has(op) && !sqlIs(value) {
			// OpIs/OpIsNot must have value: true, false, unknown or null.
			return nil, fmt.Errorf("%w: %s", ErrInvalidIs, query)
		}

		cacheKey := fmt.Sprintf("%s:%s", field, op)
		if cache[cacheKey] {
			return nil, fmt.Errorf("%w: %q.%q", ErrMultipleOperator, field, strings.ToLower(op.String()))
		}
		cache[cacheKey] = true

		parser, ok := parsers[tag.Type.Name]
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

			vals, err := splitString(value)
			if err != nil {
				return nil, err
			}

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

		sets = append(sets, fs)
	}

	return sets, nil
}
