// Code generated by "stringer -type=Class"; DO NOT EDIT.

package hl

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Invalid-0]
	_ = x[Bool-1]
	_ = x[Const-2]
	_ = x[Cond-3]
	_ = x[Repeat-4]
	_ = x[Var-5]
	_ = x[Def-6]
	_ = x[Fn-7]
	_ = x[Op-8]
	_ = x[Special-9]
	_ = x[Import-10]
}

const _Class_name = "InvalidBoolConstCondRepeatVarDefFnOpSpecialImport"

var _Class_index = [...]uint8{0, 7, 11, 16, 20, 26, 29, 32, 34, 36, 43, 49}

func (i Class) String() string {
	if i < 0 || i >= Class(len(_Class_index)-1) {
		return "Class(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Class_name[_Class_index[i]:_Class_index[i+1]]
}
