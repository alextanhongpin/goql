package goql

var opsByText map[string]Op

func init() {
	opsByText = make(map[string]Op)
	for k, v := range opsText {
		opsByText[v] = k
	}
}

const (
	// OpsComparable represents comparable operators.
	// Most (or all) datatypes in SQL are comparable.
	OpsComparable = OpEq | OpNeq | OpLt | OpLte | OpGt | OpGte | OpIn | OpNotIn

	// OpsNull represents nullable operators.
	OpsNull = OpIs | OpIsNot

	// OpsIn represents the variation of `in`.
	OpsIn = OpIn | OpNotIn

	// OpsLike represents the variation of `like`/`not like`.
	OpsLike = OpLike | OpNotLike | OpIlike | OpNotIlike

	// OpsFullTextSearch represents full-text-search operators.
	OpsFullTextSearch = OpFts | OpPlFts | OpPhFts | OpWFts

	// OpsRange operators supports range operations.
	OpsRange = OpCs | OpCd | OpOv | OpSl | OpSr | OpNxr | OpNxl | OpAdj

	// OpsMany operators supports multiple values.
	OpsMany = OpsIn | OpsLike
)

// Op represents a SQL operator.
type Op int

func (op Op) String() string {
	return opsText[op]
}

// Valid returns true if the ops is set.
func (op Op) Valid() bool {
	return op != 0
}

// Has checks if the ops is part of the rule.
func (op Op) Has(tgt Op) bool {
	return op&tgt == tgt
}

// Is checks if the op is equal the other value.
func (op Op) Is(tgt Op) bool {
	return op == tgt
}

// Useful reference:
// functions: https://www.postgresql.org/docs/14/functions.html
// range operators:  https://www.postgresql.org/docs/14/functions-range.html
// array operators: https://www.postgresql.org/docs/current/functions-array.html
const (
	OpEq       Op = 1 << iota // =, equals, e.g. name.eq=john appleseed
	OpNeq                     // <> or !=, not equals, e.g. name.neq=john appleseed
	OpLt                      // <, less than
	OpLte                     // <=, less than equals
	OpGt                      // >, greater than
	OpGte                     // >=, greater than equals
	OpLike                    // like, multi-values, e.g. name.like=john*
	OpIlike                   // ilike, multi-values, same as like, but case insensitive, e.g. name.ilike=john%
	OpNotLike                 // not like, multi-values, e.g. name.notlike=john*
	OpNotIlike                // not ilike, multi-values, e.g. name.notilike=john*
	OpIn                      // in, multi-values, name.in=alice&name.in=bob
	OpNotIn                   // not in, multi-values, name.notin=alice&name.notin=bob
	OpIs                      // is, checking for exact equality (null,true,false,unknown), e.g. age.is=null
	OpIsNot                   // is not, e.g. age.isnot=null
	OpFts                     // Full-Text search using to_tsquery
	OpPlFts                   // Full-Text search using plain to tsquery
	OpPhFts                   // Full-Text search using phrase to tsquery
	OpWFts                    // Full-Text search using word.
	OpCs                      // @>, contains, e.g. ?tags.cs=apple&tags.cs=orange
	OpCd                      // <@, contained in e.g. ?values.cd=1&values.cd=2
	OpOv                      // &&, overlap
	OpSl                      // <<, strictly left of
	OpSr                      // >>, strictly right of
	OpNxr                     // &<
	OpNxl                     // &>
	OpAdj                     // -|-
	OpNot                     // not
	OpOr                      // or, e.g. or=(age.gt:10,age.lt:100)
	OpAnd                     // and, e.g. and=(or.(married_at.isnot:null, married_at.gt:now))
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
