package slo

import (
	"regexp"

	"github.com/nobl9/nobl9-go/internal/validation"
)

// NewRelicMetric represents metric from NewRelic
type NewRelicMetric struct {
	NRQL *string `json:"nrql"`
}

var newRelicValidation = validation.New[NewRelicMetric](
	validation.ForPointer(func(n NewRelicMetric) *string { return n.NRQL }).
		WithName("nrql").
		Required().
		Cascade(validation.CascadeModeStop).
		Rules(validation.StringNotEmpty()).
		Rules(validation.StringDenyRegexp(regexp.MustCompile(`(?i)[\n\s](since|until)([\n\s]|$)`)).
			WithDetails("query must not contain 'since' or 'until' keywords (case insensitive)")),
)
