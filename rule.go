package goql

type RuleType int

//go:generate stringer -type RuleType -trimprefix RuleType
const (
	RuleTypeString RuleType = iota
	RuleTypeInt
	RuleTypeFloat
	RuleTypeBool
	RuleTypeDate
)

type Rule struct {
	typ RuleType
	ops Ops
}

func (r Rule) Copy() Rule {
	return Rule{
		typ: r.typ,
		ops: r.ops.Copy(),
	}
}

func newPrimitiveRule(typ RuleType) Rule {
	return Rule{
		typ: typ,
		ops: Ops{
			OpEq:      true,
			OpNeq:     true,
			OpLt:      true,
			OpLte:     true,
			OpGt:      true,
			OpGte:     true,
			OpBetween: true,
			OpIn:      true,
			OpIs:      true,
			OpIsNot:   true,
		},
	}
}

var (
	StringRule = newPrimitiveRule(RuleTypeString)
	IntRule    = newPrimitiveRule(RuleTypeInt)
	FloatRule  = newPrimitiveRule(RuleTypeFloat)
	BoolRule   = newPrimitiveRule(RuleTypeBool)
)

// Some op cannot be used together, e.g. gt and gte
