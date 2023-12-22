// Code generated by go-enum DO NOT EDIT.
// Version: 0.5.10
// Revision: f2bc4fe0810071d365417dd2e3ce1b4d64b82c47
// Build Date: 2023-11-13T16:50:24Z
// Built By: goreleaser

package sdk

import (
	"fmt"
	"strings"
)

const (
	// ObjectSourceTypeFile is a ObjectSourceType of type File.
	ObjectSourceTypeFile ObjectSourceType = iota + 1
	// ObjectSourceTypeDirectory is a ObjectSourceType of type Directory.
	ObjectSourceTypeDirectory
	// ObjectSourceTypeGlobPattern is a ObjectSourceType of type GlobPattern.
	ObjectSourceTypeGlobPattern
	// ObjectSourceTypeURL is a ObjectSourceType of type URL.
	ObjectSourceTypeURL
	// ObjectSourceTypeReader is a ObjectSourceType of type Reader.
	ObjectSourceTypeReader
)

var ErrInvalidObjectSourceType = fmt.Errorf("not a valid ObjectSourceType, try [%s]", strings.Join(_ObjectSourceTypeNames, ", "))

const _ObjectSourceTypeName = "FileDirectoryGlobPatternURLReader"

var _ObjectSourceTypeNames = []string{
	_ObjectSourceTypeName[0:4],
	_ObjectSourceTypeName[4:13],
	_ObjectSourceTypeName[13:24],
	_ObjectSourceTypeName[24:27],
	_ObjectSourceTypeName[27:33],
}

// ObjectSourceTypeNames returns a list of possible string values of ObjectSourceType.
func ObjectSourceTypeNames() []string {
	tmp := make([]string, len(_ObjectSourceTypeNames))
	copy(tmp, _ObjectSourceTypeNames)
	return tmp
}

var _ObjectSourceTypeMap = map[ObjectSourceType]string{
	ObjectSourceTypeFile:        _ObjectSourceTypeName[0:4],
	ObjectSourceTypeDirectory:   _ObjectSourceTypeName[4:13],
	ObjectSourceTypeGlobPattern: _ObjectSourceTypeName[13:24],
	ObjectSourceTypeURL:         _ObjectSourceTypeName[24:27],
	ObjectSourceTypeReader:      _ObjectSourceTypeName[27:33],
}

// String implements the Stringer interface.
func (x ObjectSourceType) String() string {
	if str, ok := _ObjectSourceTypeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("ObjectSourceType(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ObjectSourceType) IsValid() bool {
	_, ok := _ObjectSourceTypeMap[x]
	return ok
}

var _ObjectSourceTypeValue = map[string]ObjectSourceType{
	_ObjectSourceTypeName[0:4]:   ObjectSourceTypeFile,
	_ObjectSourceTypeName[4:13]:  ObjectSourceTypeDirectory,
	_ObjectSourceTypeName[13:24]: ObjectSourceTypeGlobPattern,
	_ObjectSourceTypeName[24:27]: ObjectSourceTypeURL,
	_ObjectSourceTypeName[27:33]: ObjectSourceTypeReader,
}

// ParseObjectSourceType attempts to convert a string to a ObjectSourceType.
func ParseObjectSourceType(name string) (ObjectSourceType, error) {
	if x, ok := _ObjectSourceTypeValue[name]; ok {
		return x, nil
	}
	return ObjectSourceType(0), fmt.Errorf("%s is %w", name, ErrInvalidObjectSourceType)
}
