package annotation

import (
	"fmt"

	"github.com/nobl9/nobl9-go/manifest/v1alpha"
	"github.com/nobl9/nobl9-go/validation"
)

var annotationValidation = validation.New[Annotation](
	validation.For(func(p Annotation) Metadata { return p.Metadata }).
		Include(metadataValidation),
	validation.For(func(p Annotation) Spec { return p.Spec }).
		WithName("spec").
		Include(specValidation),
)

var metadataValidation = validation.New[Metadata](
	v1alpha.FieldRuleMetadataName(func(m Metadata) string { return m.Name }),
	validation.For(func(m Metadata) string { return m.Project }).
		WithName("metadata.project").
		Omitempty().
		Rules(validation.StringIsDNSSubdomain()),
)

var specValidation = validation.New[Spec](
	validation.For(func(s Spec) string { return s.Slo }).
		WithName("slo").
		Required().
		Rules(validation.StringIsDNSSubdomain()),
	validation.For(func(s Spec) string { return s.ObjectiveName }).
		WithName("objectiveName").
		Omitempty().
		Rules(validation.StringIsDNSSubdomain()),
	validation.For(func(s Spec) string { return s.Description }).
		WithName("description").
		Required().
		Rules(validation.StringLength(0, 1000)),
	validation.For(validation.GetSelf[Spec]()).
		Rules(endTimeAfterStartTime),
)

func validate(p Annotation) *v1alpha.ObjectError {
	return v1alpha.ValidateObject(annotationValidation, p)
}

const errorCodeEndTimeAfterStartTime validation.ErrorCode = "end_time_after_start_time"

var endTimeAfterStartTime = validation.NewSingleRule(func(s Spec) error {
	if !s.EndTime.After(s.StartTime) {
		return &validation.RuleError{
			Message: fmt.Sprintf(`endTime '%s' must be after startTime '%s'`, s.EndTime, s.StartTime),
			Code:    errorCodeEndTimeAfterStartTime,
		}
	}

	return nil
})