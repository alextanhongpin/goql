package goql

const (
	RuleString = OpEq | OpNeq | OpLt | OpLte | OpGt | OpGte | OpLike | OpIlike | OpIn | OpIs | OpFts | OpPlFts | OpPhFts | OpWFts | OpNot
	RuleInt    = OpEq | OpNeq | OpLt | OpLte | OpGt | OpGte | OpIn | OpNot
	RuleFloat  = RuleInt
	RuleBool   = OpEq | OpNeq | OpIs | OpNot
	RuleNot    = OpNot | OpLike | OpIlike | OpIs | OpIn
)
