// Code generated by "stringer -trimprefix SourceType -type SourceType pkg/definitions/source.go"; DO NOT EDIT.

package definitions

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SourceTypeFile-0]
	_ = x[SourceTypeDirectory-1]
	_ = x[SourceTypeGlobPattern-2]
	_ = x[SourceTypeURL-3]
	_ = x[SourceTypeInput-4]
}

const _SourceType_name = "FileDirectoryGlobPatternURLInput"

var _SourceType_index = [...]uint8{0, 4, 13, 24, 27, 32}

func (i SourceType) String() string {
	if i < 0 || i >= SourceType(len(_SourceType_index)-1) {
		return "SourceType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SourceType_name[_SourceType_index[i]:_SourceType_index[i+1]]
}