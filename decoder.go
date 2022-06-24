package goql

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

const (
	StructTag    = "q"
	ValSeparator = ":"
)

var (
	ErrMultipleOperator = errors.New("goql: multiple op")
	ErrUnknownOperator  = errors.New("goql: unknown op")
	ErrBadOperator      = errors.New("goql: bad op")
	ErrInvalidIs        = errors.New("goql: 'IS' must be followed by {true, false, null, unknown}")
	ErrInvalidArray     = errors.New("goql: invalid array")
	ErrUnknownField     = errors.New("goql: unknown field")
	ErrUnknownParser    = errors.New("goql: unknown parser")
)

type FieldSet struct {
	Tag      Tag
	Name     string
	Value    any
	RawValue string
	Op       string
}

type Decoder[T any] struct {
	opsByField map[string]Op
	tagByField map[string]Tag
	parsers    map[string]ParserFn
	tag        string
}

func NewDecoder[T any]() *Decoder[T] {
	var t T

	parserByType := NewParsers()
	opsByField := make(map[string]Op)
	tagByField := ParseStruct(t, StructTag)

	for field, tag := range tagByField {
		t := tag.Type

		// All types are comparable.
		ops := OpsComparable

		// Null type have special operators.
		if t.Null {
			ops |= OpsNull
		}

		// Array type have special operators.
		if t.Array {
			ops |= OpsRange
		}

		// String types have special operators.
		if t.Name == reflect.String.String() {
			ops |= OpsText
		}

		opsByField[field] = ops
	}

	return &Decoder[T]{
		tagByField: tagByField,
		opsByField: opsByField,
		tag:        StructTag,
		parsers:    parserByType,
	}
}

func (d *Decoder[T]) SetStructTag(tag string) {
	if tag == "" {
		panic("tag cannot be empty")
	}

	var t T
	d.tag = tag
	d.tagByField = ParseStruct(t, tag)
}

func (d *Decoder[T]) SetFieldOps(opsByField map[string]Op) {
	d.opsByField = opsByField
}

func (d *Decoder[T]) SetParsers(parserByType map[string]ParserFn) {
	d.parsers = parserByType
}

func (d *Decoder[T]) Decode(values url.Values) ([]FieldSet, error) {
	return Decode(d.opsByField, d.tagByField, d.parsers, values)
}

func Decode(opsByField map[string]Op, tagByField map[string]Tag, parsers map[string]ParserFn, values url.Values) ([]FieldSet, error) {
	cache := make(map[string]bool)

	var sets []FieldSet

	for field, rule := range opsByField {
		tag, ok := tagByField[field]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownField, field)
		}

		for _, v := range values[field] {
			opsByField, val := split2(v, ValSeparator)

			op, ok := ParseOp(opsByField)
			if !ok {
				return nil, fmt.Errorf("%w: %s", ErrUnknownOperator, v)
			}

			if !rule.Has(op) {
				return nil, fmt.Errorf("%w: %s", ErrBadOperator, v)
			}

			if OpsIs.Has(op) && !sqlIs(val) {
				// OpIs/OpIsNot must have value: true, false, unknown or null.
				return nil, fmt.Errorf("%w: %s", ErrInvalidIs, v)
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
				RawValue: val,
			}

			switch {
			case OpsIn.Has(op), tag.Type.Array:
				val, ok = Unquote(val, '{', '}')
				if !ok {
					return nil, fmt.Errorf("%w: missing parantheses: %s", ErrInvalidArray, val)
				}

				vals, err := splitString(val)
				if err != nil {
					return nil, err
				}

				res, err := Map(vals, parser)
				if err != nil {
					return nil, err
				}

				fs.Value = res
			default:
				res, err := parser(val)
				if err != nil {
					return nil, err
				}
				fs.Value = res
			}

			sets = append(sets, fs)
		}
	}

	return sets, nil
}
