// Code generated by "stringer -type RuleType -trimprefix RuleType"; DO NOT EDIT.

package goql

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[RuleTypeString-0]
	_ = x[RuleTypeInt-1]
	_ = x[RuleTypeFloat-2]
	_ = x[RuleTypeBool-3]
	_ = x[RuleTypeDate-4]
}

const _RuleType_name = "StringIntFloatBoolDate"

var _RuleType_index = [...]uint8{0, 6, 9, 14, 18, 22}

func (i RuleType) String() string {
	if i < 0 || i >= RuleType(len(_RuleType_index)-1) {
		return "RuleType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _RuleType_name[_RuleType_index[i]:_RuleType_index[i+1]]
}
