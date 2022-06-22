package goql

/*

By default, all datatypes are comparable: =, <>, <, <=, >, >=
Since WHERE col IN (...) is pretty common in SQL, OpIn is added for most types.
If the datatype can be NULL, then OpIs will be appended.
*/
const (
	OpsComparable = OpEq | OpNeq | OpLt | OpLte | OpGt | OpGte | OpIn
	OpsNull       = OpIs | OpNot
	OpsNot        = OpLike | OpIlike | OpIn | OpIs
	OpsText       = OpsComparable | OpLike | OpIlike | OpIn | OpFts | OpPlFts | OpPhFts | OpWFts
	OpsRange      = OpCs | OpCd | OpOv | OpSl | OpSr | OpNxr | OpNxl | OpAdj
)
