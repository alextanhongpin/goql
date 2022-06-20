package goql

import "strings"

type Op int

func (o Op) Has(tgt Op) bool {
	return o&tgt == tgt
}

//go:generate stringer -type Op -trimprefix Op
const (
	OpEq Op = 1 << iota
	OpNeq
	OpLt
	OpLte
	OpGt
	OpGte
	OpLike
	OpIlike
	OpIn
	OpIs  // IS, checking for exact equality (null,true,false,unknown)
	OpFts // Full-Text search using to_tsquery
	OpPlFts
	OpPhFts
	OpWFts
	OpCs  // @>, contains, e.g. ?tags=cs.{example,new}
	OpCd  // <@, contained in e.g. ?values=cd.{1,2,3}
	OpOv  // &&, Overlap
	OpSl  // <<, strictly left of
	OpSr  // >>, strictly right of
	OpNxr // &<
	OpNxl // &>
	OpAdj // -|-
	OpNot
	OpOr
	OpAnd
)

func ParseOp(op string) (Op, bool) {
	for i := OpAnd; i > 0; i = i >> 1 {
		if strings.EqualFold(i.String(), op) {
			return i, true
		}
	}

	return 0, false
}
