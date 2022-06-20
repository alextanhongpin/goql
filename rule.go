package goql

const (
	RuleComparable = OpEq | OpNeq | OpLt | OpLte | OpGt | OpGte
	RuleNull       = OpIs | OpNot
	RuleList       = OpIn | OpNot
	RuleWhere      = RuleComparable | RuleList | RuleNull
	RuleNegate     = OpLike | OpIlike | OpIn | OpIs
	RuleText       = RuleWhere | OpLike | OpIlike | OpIn | OpFts | OpPlFts | OpPhFts | OpWFts
	RuleInt        = RuleWhere
	RuleFloat      = RuleWhere
	RuleBool       = OpEq | OpNeq | OpIs | OpNot
)
