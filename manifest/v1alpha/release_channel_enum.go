// Code generated by go-enum DO NOT EDIT.
// Version: 0.5.6
// Revision: 97611fddaa414f53713597918c5e954646cb8623
// Build Date: 2023-03-26T21:38:06Z
// Built By: goreleaser

package v1alpha

import (
	"fmt"
	"strings"
)

const (
	// ReleaseChannelStable is a ReleaseChannel of type Stable.
	ReleaseChannelStable ReleaseChannel = iota + 1
	// ReleaseChannelBeta is a ReleaseChannel of type Beta.
	ReleaseChannelBeta
	// ReleaseChannelAlpha is a ReleaseChannel of type Alpha.
	ReleaseChannelAlpha
)

var ErrInvalidReleaseChannel = fmt.Errorf("not a valid ReleaseChannel, try [%s]", strings.Join(_ReleaseChannelNames, ", "))

const _ReleaseChannelName = "stablebetaalpha"

var _ReleaseChannelNames = []string{
	_ReleaseChannelName[0:6],
	_ReleaseChannelName[6:10],
	_ReleaseChannelName[10:15],
}

// ReleaseChannelNames returns a list of possible string values of ReleaseChannel.
func ReleaseChannelNames() []string {
	tmp := make([]string, len(_ReleaseChannelNames))
	copy(tmp, _ReleaseChannelNames)
	return tmp
}

// ReleaseChannelValues returns a list of the values for ReleaseChannel
func ReleaseChannelValues() []ReleaseChannel {
	return []ReleaseChannel{
		ReleaseChannelStable,
		ReleaseChannelBeta,
		ReleaseChannelAlpha,
	}
}

var _ReleaseChannelMap = map[ReleaseChannel]string{
	ReleaseChannelStable: _ReleaseChannelName[0:6],
	ReleaseChannelBeta:   _ReleaseChannelName[6:10],
	ReleaseChannelAlpha:  _ReleaseChannelName[10:15],
}

// String implements the Stringer interface.
func (x ReleaseChannel) String() string {
	if str, ok := _ReleaseChannelMap[x]; ok {
		return str
	}
	return fmt.Sprintf("ReleaseChannel(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ReleaseChannel) IsValid() bool {
	_, ok := _ReleaseChannelMap[x]
	return ok
}

var _ReleaseChannelValue = map[string]ReleaseChannel{
	_ReleaseChannelName[0:6]:                    ReleaseChannelStable,
	strings.ToLower(_ReleaseChannelName[0:6]):   ReleaseChannelStable,
	_ReleaseChannelName[6:10]:                   ReleaseChannelBeta,
	strings.ToLower(_ReleaseChannelName[6:10]):  ReleaseChannelBeta,
	_ReleaseChannelName[10:15]:                  ReleaseChannelAlpha,
	strings.ToLower(_ReleaseChannelName[10:15]): ReleaseChannelAlpha,
}

// ParseReleaseChannel attempts to convert a string to a ReleaseChannel.
func ParseReleaseChannel(name string) (ReleaseChannel, error) {
	if x, ok := _ReleaseChannelValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _ReleaseChannelValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return ReleaseChannel(0), fmt.Errorf("%s is %w", name, ErrInvalidReleaseChannel)
}
