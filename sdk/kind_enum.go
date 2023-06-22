// Code generated by go-enum DO NOT EDIT.
// Version: 0.5.6
// Revision: 97611fddaa414f53713597918c5e954646cb8623
// Build Date: 2023-03-26T21:38:06Z
// Built By: goreleaser

package sdk

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const (
	KindSLO Kind = iota
	KindService
	KindAgent
	KindAlertPolicy
	KindAlertSilence
	KindAlert
	KindProject
	KindAlertMethod
	KindMetricSource
	KindDirect
	KindDataExport
	KindUsageSummary
	KindRoleBinding
	KindSLOErrorBudgetStatus
	KindAnnotation
	KindGroup
)

var ErrInvalidKind = errors.New("not a valid Kind")

const _KindName = "SLOServiceAgentAlertPolicyAlertSilenceAlertProjectAlertMethodMetricSourceDirectDataExportUsageSummaryRoleBindingSLOErrorBudgetStatusAnnotationGroup"

var _KindMap = map[Kind]string{
	KindSLO:                  _KindName[0:3],
	KindService:              _KindName[3:10],
	KindAgent:                _KindName[10:15],
	KindAlertPolicy:          _KindName[15:26],
	KindAlertSilence:         _KindName[26:38],
	KindAlert:                _KindName[38:43],
	KindProject:              _KindName[43:50],
	KindAlertMethod:          _KindName[50:61],
	KindMetricSource:         _KindName[61:73],
	KindDirect:               _KindName[73:79],
	KindDataExport:           _KindName[79:89],
	KindUsageSummary:         _KindName[89:101],
	KindRoleBinding:          _KindName[101:112],
	KindSLOErrorBudgetStatus: _KindName[112:132],
	KindAnnotation:           _KindName[132:142],
	KindGroup:                _KindName[142:147],
}

// String implements the Stringer interface.
func (x Kind) String() string {
	if str, ok := _KindMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Kind(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Kind) IsValid() bool {
	_, ok := _KindMap[x]
	return ok
}

var _KindValue = map[string]Kind{
	_KindName[0:3]:                      KindSLO,
	strings.ToLower(_KindName[0:3]):     KindSLO,
	_KindName[3:10]:                     KindService,
	strings.ToLower(_KindName[3:10]):    KindService,
	_KindName[10:15]:                    KindAgent,
	strings.ToLower(_KindName[10:15]):   KindAgent,
	_KindName[15:26]:                    KindAlertPolicy,
	strings.ToLower(_KindName[15:26]):   KindAlertPolicy,
	_KindName[26:38]:                    KindAlertSilence,
	strings.ToLower(_KindName[26:38]):   KindAlertSilence,
	_KindName[38:43]:                    KindAlert,
	strings.ToLower(_KindName[38:43]):   KindAlert,
	_KindName[43:50]:                    KindProject,
	strings.ToLower(_KindName[43:50]):   KindProject,
	_KindName[50:61]:                    KindAlertMethod,
	strings.ToLower(_KindName[50:61]):   KindAlertMethod,
	_KindName[61:73]:                    KindMetricSource,
	strings.ToLower(_KindName[61:73]):   KindMetricSource,
	_KindName[73:79]:                    KindDirect,
	strings.ToLower(_KindName[73:79]):   KindDirect,
	_KindName[79:89]:                    KindDataExport,
	strings.ToLower(_KindName[79:89]):   KindDataExport,
	_KindName[89:101]:                   KindUsageSummary,
	strings.ToLower(_KindName[89:101]):  KindUsageSummary,
	_KindName[101:112]:                  KindRoleBinding,
	strings.ToLower(_KindName[101:112]): KindRoleBinding,
	_KindName[112:132]:                  KindSLOErrorBudgetStatus,
	strings.ToLower(_KindName[112:132]): KindSLOErrorBudgetStatus,
	_KindName[132:142]:                  KindAnnotation,
	strings.ToLower(_KindName[132:142]): KindAnnotation,
	_KindName[142:147]:                  KindGroup,
	strings.ToLower(_KindName[142:147]): KindGroup,
}

// ParseKind attempts to convert a string to a Kind.
func ParseKind(name string) (Kind, error) {
	if x, ok := _KindValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _KindValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return Kind(0), fmt.Errorf("%s is %w", name, ErrInvalidKind)
}
