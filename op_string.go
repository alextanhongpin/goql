// Code generated by "stringer -type Op -trimprefix Op"; DO NOT EDIT.

package goql

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OpEq-0]
	_ = x[OpNeq-1]
	_ = x[OpLt-2]
	_ = x[OpLte-3]
	_ = x[OpGt-4]
	_ = x[OpGte-5]
	_ = x[OpIn-6]
	_ = x[OpNin-7]
	_ = x[OpBetween-8]
	_ = x[OpIs-9]
	_ = x[OpNot-10]
	_ = x[OpOverlap-11]
	_ = x[OpAt-12]
	_ = x[OpContainedBy-13]
	_ = x[OpContainedAt-14]
	_ = x[OpNow-15]
	_ = x[OpToday-16]
	_ = x[OpYesterday-17]
	_ = x[OpAny-18]
	_ = x[OpAll-19]
	_ = x[OpContains-20]
	_ = x[OpSubset-21]
	_ = x[OpSuperSet-22]
}

const _Op_name = "EqNeqLtLteGtGteInNinBetweenIsNotOverlapAtContainedByContainedAtNowTodayYesterdayAnyAllContainsSubsetSuperSet"

var _Op_index = [...]uint8{0, 2, 5, 7, 10, 12, 15, 17, 20, 27, 29, 32, 39, 41, 52, 63, 66, 71, 80, 83, 86, 94, 100, 108}

func (i Op) String() string {
	if i < 0 || i >= Op(len(_Op_index)-1) {
		return "Op(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Op_name[_Op_index[i]:_Op_index[i+1]]
}
