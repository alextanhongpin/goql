package goql

var opsByText map[string]Op

func init() {
	opsByText = make(map[string]Op)
	for k, v := range opsText {
		opsByText[v] = k
	}
}

// By default, all datatypes are comparable: =, <>, <, <=, >, >=
// Since WHERE col IN (...) is pretty common in SQL, OpIn is added for most types.
// If the datatype can be NULL, then OpIs will be appended.
const (
	OpsComparable     = OpEq | OpNeq | OpLt | OpLte | OpGt | OpGte | OpIn | OpNotIn
	OpsNull           = OpIs | OpIsNot
	OpsIn             = OpIn | OpNotIn
	OpsLike           = OpLike | OpNotLike | OpIlike | OpNotIlike
	OpsFullTextSearch = OpFts | OpPlFts | OpPhFts | OpWFts
	OpsRange          = OpCs | OpCd | OpOv | OpSl | OpSr | OpNxr | OpNxl | OpAdj
)

type Op int

func (op Op) String() string {
	return opsText[op]
}

// Valid returns true if the ops is set.
func (o Op) Valid() bool {
	return o != 0
}

// Has checks if the ops is part of the rule.
func (o Op) Has(tgt Op) bool {
	return o&tgt == tgt
}

// Is checks if the op is equal the other value.
func (o Op) Is(tgt Op) bool {
	return o == tgt
}

// Useful reference:
// range operators:  https://www.postgresql.org/docs/14/functions-range.html
const (
	OpEq       Op = 1 << iota // =, equals, e.g. name=eq:john becomes name = 'john'
	OpNeq                     // <> or !=, not equals, e.g. name=neq:john becomes name <> 'john'
	OpLt                      // <, less than
	OpLte                     // <=, less than equals
	OpGt                      // >, greater than
	OpGte                     // >=, greater than equals
	OpLike                    // like, e.g. like:john* becomes name like 'john%' (use * instead of % in query string)
	OpIlike                   // ilike, same as like, but case insensitive, e.g. name=ilike:john%
	OpNotLike                 // not like
	OpNotIlike                // not ilike
	OpIn                      // in, e.g. name=in:{1,2,3| becomes name in (1,2,3). The curly brackets is needed to indicate that the whole value is an array, since for array columns there's no way to differentiate for a single value.
	OpNotIn                   // not in
	OpIs                      // is, checking for exact equality (null,true,false,unknown), e.g. age=is:null
	OpIsNot                   // is not {null, true, false, unknown}
	OpFts                     // Full-Text search using to_tsquery
	OpPlFts                   // Full-Text search using plain to tsquery
	OpPhFts                   // Full-Text search using phrase to tsquery
	OpWFts                    // Full-Text search using word.
	OpCs                      // @>, contains, e.g. ?tags=cs:{example,new}
	OpCd                      // <@, contained in e.g. ?values=cd:{1,2,3}
	OpOv                      // &&, Overlap
	OpSl                      // <<, strictly left of
	OpSr                      // >>, strictly right of
	OpNxr                     // &<
	OpNxl                     // &>
	OpAdj                     // -|-
	OpNot                     // not
	OpOr                      // or
	OpAnd                     // and
)

func ParseOp(unk string) (Op, bool) {
	op, ok := opsByText[unk]
	return op, ok
}

var opsText = map[Op]string{
	OpEq:       "eq",
	OpNeq:      "neq",
	OpLt:       "lt",
	OpLte:      "lte",
	OpGt:       "gt",
	OpGte:      "gte",
	OpLike:     "like",
	OpIlike:    "ilike",
	OpNotLike:  "notlike",
	OpNotIlike: "notilike",
	OpIn:       "in",
	OpNotIn:    "notin",
	OpIs:       "is",
	OpIsNot:    "isnot",
	OpFts:      "fts",
	OpPlFts:    "plfts",
	OpPhFts:    "phfts",
	OpWFts:     "wfts",
	OpCs:       "cs",
	OpCd:       "cd",
	OpOv:       "ov",
	OpSl:       "sl",
	OpSr:       "sr",
	OpNxr:      "nxr",
	OpNxl:      "nxl",
	OpAdj:      "adj",
	OpNot:      "not",
	OpOr:       "or",
	OpAnd:      "and",
}
