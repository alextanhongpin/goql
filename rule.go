package goql

/*

By default, all datatypes are comparable: =, <>, <, <=, >, >=
Since WHERE col IN (...) is pretty common in SQL, OpIn is added for most types.
If the datatype can be NULL, then OpIs will be appended.
*/
const (
	RuleComparable = OpEq | OpNeq | OpLt | OpLte | OpGt | OpGte
	RuleNull       = OpIs | OpNot
	RuleList       = OpIn | OpNot
	RuleWhere      = RuleComparable | RuleList
	RuleNegate     = OpLike | OpIlike | OpIn | OpIs
	RuleText       = RuleWhere | OpLike | OpIlike | OpIn | OpFts | OpPlFts | OpPhFts | OpWFts
	RuleInt        = RuleWhere
	RuleFloat      = RuleWhere
	RuleBool       = OpEq | OpNeq | OpIs | OpNot
	RuleJSON       = RuleWhere
	RuleRange      = OpCs | OpCd | OpOv | OpAdj
)

var opsByPgType = map[string]Op{
	pgTypeTimestamp:       RuleWhere,
	pgTypeTimestampTz:     RuleWhere,
	pgTypeDate:            RuleWhere,
	pgTypeTime:            RuleWhere,
	pgTypeTimeTz:          RuleWhere,
	pgTypeInterval:        RuleWhere | RuleRange,
	pgTypeInet:            RuleWhere,
	pgTypeCidr:            RuleWhere,
	pgTypeMacaddr:         RuleWhere,
	pgTypeBoolean:         RuleBool,
	pgTypeReal:            RuleInt,
	pgTypeDoublePrecision: RuleFloat,
	pgTypeSmallint:        RuleInt,
	pgTypeInteger:         RuleInt,
	pgTypeBigint:          RuleInt,
	pgTypeSmallserial:     RuleInt,
	pgTypeSerial:          RuleInt,
	pgTypeBigserial:       RuleInt,
	pgTypeVarchar:         RuleText,
	pgTypeChar:            RuleText,
	pgTypeText:            RuleText,
	pgTypeJSON:            RuleJSON,
	pgTypeJSONB:           RuleJSON,
	pgTypeBytea:           RuleWhere,
}
