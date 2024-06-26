package alertsilence

import (
	"fmt"
	"time"

	validationV1Alpha "github.com/nobl9/nobl9-go/internal/manifest/v1alpha"
	"github.com/nobl9/nobl9-go/internal/validation"
	"github.com/nobl9/nobl9-go/manifest"
	"github.com/nobl9/nobl9-go/manifest/v1alpha"
)

func validate(s AlertSilence) *v1alpha.ObjectError {
	return v1alpha.ValidateObject(validator, s, manifest.KindAlertSilence)
}

var validator = validation.New[AlertSilence](
	validationV1Alpha.FieldRuleAPIVersion(func(a AlertSilence) manifest.Version { return a.APIVersion }),
	validationV1Alpha.FieldRuleKind(func(a AlertSilence) manifest.Kind { return a.Kind }, manifest.KindAlertSilence),
	validation.For(validation.GetSelf[AlertSilence]()).
		Cascade(validation.CascadeModeStop).
		Include(metadataValidation).
		Rules(alertPolicyProjectConsistencyRule),
	validation.For(func(s AlertSilence) Spec { return s.Spec }).
		WithName("spec").
		Include(specValidation),
)

var metadataValidation = validation.New[AlertSilence](
	validationV1Alpha.FieldRuleMetadataName(func(s AlertSilence) string { return s.Metadata.Name }),
	validationV1Alpha.FieldRuleMetadataProject(func(s AlertSilence) string { return s.Metadata.Project }),
)

var specValidation = validation.New[Spec](
	validation.For(func(s Spec) string { return s.SLO }).
		WithName("slo").
		Required().
		Rules(validation.StringIsDNSSubdomain()),
	validation.For(func(s Spec) string { return s.Description }).
		WithName("description").
		Rules(validation.StringDescription()),
	validation.For(func(s Spec) AlertPolicySource { return s.AlertPolicy }).
		WithName("alertPolicy").
		Include(alertPolicySourceValidation),
	validation.For(func(s Spec) Period { return s.Period }).
		WithName("period").
		Cascade(validation.CascadeModeStop).
		Rules(
			validation.MutuallyExclusive(true, map[string]func(p Period) any{
				"duration": func(p Period) any { return p.Duration },
				"endTime":  func(p Period) any { return p.EndTime },
			}),
		).
		Include(
			validation.New[Period](
				validation.Transform(func(p Period) string { return p.Duration }, time.ParseDuration).
					WithName("duration").
					Rules(validation.GreaterThan[time.Duration](0)),
			).When(func(p Period) bool { return p.Duration != "" }),
		).
		Include(
			validation.New[Period](
				validation.For(validation.GetSelf[Period]()).
					Rules(endTimeNotBeforeStartTimeRule),
			).When(func(p Period) bool { return p.EndTime != nil }),
		),
)

var alertPolicySourceValidation = validation.New[AlertPolicySource](
	validationV1Alpha.FieldRuleMetadataName(func(s AlertPolicySource) string { return s.Name }).
		WithName("name"),
	validation.For(func(s AlertPolicySource) string { return s.Project }).
		WithName("project").
		OmitEmpty().
		Rules(validation.StringIsDNSSubdomain()),
)

const errorCodeEndTimeNotBeforeOrNotEqualStartTime = "end_time_not_before_or_not_equal_start_time"
const errorCodeInconsistentProject = "alert_policy_project_inconsistent"

var endTimeNotBeforeStartTimeRule = validation.NewSingleRule(func(p Period) error {
	if p.EndTime != nil && p.StartTime != nil && !p.EndTime.After(*p.StartTime) {
		return validation.NewRuleError(
			fmt.Sprintf(`endTime '%s' must be after startTime '%s'`, p.EndTime, p.StartTime),
			errorCodeEndTimeNotBeforeOrNotEqualStartTime,
		)
	}

	return nil
})

// alertPolicyProjectConsistencyRule validates if user provide the same project (or empty) for the alert policy
// as declared in metadata for AlertSilence. Should be removed when cross-project Alert Policy is allowed PI-622.
var alertPolicyProjectConsistencyRule = validation.NewSingleRule(func(s AlertSilence) error {
	if s.Spec.AlertPolicy.Project != "" && s.Spec.AlertPolicy.Project != s.Metadata.Project {
		return validation.NewPropertyError(
			"spec.alertPolicy.project",
			s.Spec.AlertPolicy.Project,
			validation.NewRuleError(
				fmt.Sprintf(
					`alertPolicy project '%s' must be the same as in alertSilence metadata project: '%s'`,
					s.Spec.AlertPolicy.Project, s.Metadata.Project,
				),
				errorCodeInconsistentProject,
			),
		)
	}

	return nil
})
