// Package nobl9 provide an abstraction for communication with API
package nobl9

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// APIVersion is a value of valid apiVersions
const (
	APIVersion = "n9/v1alpha"
)

// HiddenValue can be used as a value of a secret field and is ignored during saving
const HiddenValue = "[hidden]"

// Possible values of field kind for valid Objects.
const (
	KindSLO         = "SLO"
	KindService     = "Service"
	KindAgent       = "Agent"
	KindProject     = "Project"
	KindAlertPolicy = "AlertPolicy"
	KindAlert       = "Alert"
	KindAlertMethod = "AlertMethod"
	KindDirect      = "Direct"
	KindDataExport  = "DataExport"
	KindRoleBinding = "RoleBinding"
)

// APIObjects - all Objects available for this version of API
// Sorted in order of applying
type APIObjects struct {
	SLOs          []SLO
	Services      []Service
	Agents        []Agent
	AlertPolicies []AlertPolicy
	Alerts        []Alert
	AlertMethods  []AlertMethod
	Directs       []Direct
	DataExports   []DataExport
	Projects      []Project
	RoleBindings  []RoleBinding
}

type Payload struct {
	objects []AnyJSONObj
}

func Newpayload(org string) {
}

func (p *Payload) AddObject(in interface{}) {
	p.objects = append(p.objects, toAnyJSONObj(in))
}

func (p *Payload) GetObjects() []AnyJSONObj {
	return p.objects
}

func toAnyJSONObj(in interface{}) AnyJSONObj {
	tmp, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	var out AnyJSONObj
	if err := json.Unmarshal(tmp, &out); err != nil {
		panic(err)
	}
	return out
}

type CompositeBurnRateCondition struct {
	Value    float64 `json:"value"`
	Operator string  `json:"op,omitempty" example:"gte"`
}

type Composite struct {
	BudgetTarget      float64                     `json:"target"`
	BurnRateCondition *CompositeBurnRateCondition `json:"burnRateCondition"`
}

// CountMetricsSpec represents set of two time series of good and total counts
type CountMetricsSpec struct {
	Incremental *bool       `json:"incremental"`
	GoodMetric  *MetricSpec `json:"good"`
	TotalMetric *MetricSpec `json:"total"`
}

// RawMetricSpec represents integration with a metric source for a particular threshold
type RawMetricSpec struct {
	MetricQuery *MetricSpec `json:"query" validate:"required"`
}

// MetricSpec defines single time series kobtained from data source
type MetricSpec struct {
	Prometheus          *PrometheusMetric          `json:"prometheus,omitempty"`
	Datadog             *DatadogMetric             `json:"datadog,omitempty"`
	NewRelic            *NewRelicMetric            `json:"newRelic,omitempty"`
	AppDynamics         *AppDynamicsMetric         `json:"appDynamics,omitempty"`
	Splunk              *SplunkMetric              `json:"splunk,omitempty"`
	Lightstep           *LightstepMetric           `json:"lightstep,omitempty"`
	SplunkObservability *SplunkObservabilityMetric `json:"splunkObservability,omitempty"`
	Dynatrace           *DynatraceMetric           `json:"dynatrace,omitempty"`
	Elasticsearch       *ElasticsearchMetric       `json:"elasticsearch,omitempty"`
	ThousandEyes        *ThousandEyesMetric        `json:"thousandEyes,omitempty"`
	Graphite            *GraphiteMetric            `json:"graphite,omitempty"`
	BigQuery            *BigQueryMetric            `json:"bigQuery,omitempty"`
	OpenTSDB            *OpenTSDBMetric            `json:"opentsdb,omitempty"`
	GrafanaLoki         *GrafanaLokiMetric         `json:"grafanaLoki,omitempty"`
}

// PrometheusMetric represents metric from Prometheus
type PrometheusMetric struct {
	PromQL *string `json:"promql"`
}

// DatadogMetric represents metric from Datadog
type DatadogMetric struct {
	Query *string `json:"query"`
}

// NewRelicMetric represents metric from NewRelic
type NewRelicMetric struct {
	NRQL *string `json:"nrql"`
}

// ThousandEyesMetric represents metric from ThousandEyes
type ThousandEyesMetric struct {
	TestID *int64 `json:"testID"`
}

// AppDynamicsMetric represents metric from AppDynamics
type AppDynamicsMetric struct {
	ApplicationName *string `json:"applicationName"`
	MetricPath      *string `json:"metricPath"`
}

// SplunkMetric represents metric from Splunk
type SplunkMetric struct {
	Query     *string `json:"query"`
	FieldName *string `json:"fieldName"`
}

// LightstepMetric represents metric from Lightstep
type LightstepMetric struct {
	StreamID   *string  `json:"streamId"`
	TypeOfData *string  `json:"typeOfData"`
	Percentile *float64 `json:"percentile,omitempty"`
}

// SplunkObservabilityMetric represents metric from SplunkObservability
type SplunkObservabilityMetric struct {
	Program *string `json:"program"`
}

// DynatraceMetric represents metric from Dynatrace.
type DynatraceMetric struct {
	MetricSelector *string `json:"metricSelector"`
}

// ElasticsearchMetric represents metric from Elasticsearch.
type ElasticsearchMetric struct {
	Index *string `json:"index"`
	Query *string `json:"query"`
}

// GraphiteMetric represents metric from Graphite.
type GraphiteMetric struct {
	MetricPath *string `json:"metricPath"`
}

// BigQueryMetric represents metric from BigQuery
type BigQueryMetric struct {
	Query     string `json:"query"`
	ProjectID string `json:"projectId"`
	Location  string `json:"location"`
}

// OpenTSDBMetric represents metric from OpenTSDB.
type OpenTSDBMetric struct {
	Query *string `json:"query"`
}

// GrafanaLokiMetric represents metric from GrafanaLokiMetric.
type GrafanaLokiMetric struct {
	Logql *string `json:"logql"`
}

// ThresholdBase base structure representing a threshold
type ThresholdBase struct {
	DisplayName string  `json:"displayName"`
	Value       float64 `json:"value"`
}

// Threshold represents single threshold for SLO, for internal usage
type Threshold struct {
	ThresholdBase
	// <!-- Go struct field and type names renaming budgetTarget to target has been postponed after GA as requested
	// in PC-1240. -->
	BudgetTarget *float64 `json:"target"`
	// <!-- Go struct field and type names renaming thresholds to objectives has been postponed after GA as requested
	// in PC-1240. -->
	TimeSliceTarget *float64          `json:"timeSliceTarget,omitempty" example:"0.9"`
	CountMetrics    *CountMetricsSpec `json:"countMetrics,omitempty"`
	RawMetric       *RawMetricSpec    `json:"rawMetric,omitempty"`
	Operator        *string           `json:"op,omitempty" example:"lte"`
}

// Indicator represents integration with metric source can be. e.g. Prometheus, Datadog, for internal usage
type Indicator struct {
	MetricSource *MetricSourceSpec `json:"metricSource"`
}

type MetricSourceSpec struct {
	Project string `json:"project,omitempty"`
	Name    string `json:"name"`
	Kind    string `json:"kind"`
}

// SLOSpec represents content of Spec typical for SLO Object
type SLOSpec struct {
	Description     string       `json:"description"` //nolint:lll
	Indicator       Indicator    `json:"indicator"`
	BudgetingMethod string       `json:"budgetingMethod"`
	Thresholds      []Threshold  `json:"objectives"`
	Service         string       `json:"service"`
	TimeWindows     []TimeWindow `json:"timeWindows"`
	AlertPolicies   []string     `json:"alertPolicies"`
	Attachments     []Attachment `json:"attachments,omitempty"`
	CreatedAt       string       `json:"createdAt,omitempty"`
	Composite       *Composite   `json:"composite,omitempty"`
}

// SLO struct which mapped one to one with kind: slo yaml definition, external usage
type SLO struct {
	ObjectHeader
	Spec SLOSpec `json:"spec"`
}

// Time Series

type SLOTimeSeries struct {
	MetadataHolder
	TimeWindows                 []TimeWindowTimeSeries `json:"timewindows,omitempty"`
	RawSLIPercentilesTimeSeries Percentile             `json:"percentiles,omitempty"`
}

type ThresholdTimeSeries struct {
	ThresholdBase
	InstantaneousBurnRateTimeSeries
	CumulativeBurnedTimeSeries
	Status   ThresholdTimeSeriesStatus `json:"status"`
	Operator *string                   `json:"op,omitempty"`
	CountsSLITimeSeries
	BurnDownTimeSeries
}

type ThresholdTimeSeriesStatus struct {
	BurnedBudget            *float64 `json:"burnedBudget,omitempty" example:"0.25"`
	RemainingBudget         *float64 `json:"errorBudgetRemainingPercentage,omitempty" example:"0.25"`
	RemainingBudgetDuration *float64 `json:"errorBudgetRemaining,omitempty" example:"300"`
	InstantaneousBurnRate   *float64 `json:"instantaneousBurnRate,omitempty" example:"1.25"`
	Condition               *string  `json:"condition,omitempty" example:"ok"`
}

type TimeWindowTimeSeries struct {
	TimeWindow `json:"timewindow,omitempty"`
	// <!-- Go struct field and type names renaming thresholds to objectives has been postponed after GA as requested
	// in PC-1240. -->
	Thresholds []ThresholdTimeSeries `json:"objectives,omitempty"`
}

const (
	P1  string = "p1"
	P5  string = "p5"
	P10 string = "p10"
	P50 string = "p50"
	P90 string = "p90"
	P95 string = "p95"
	P99 string = "p99"
)

func GetAvailablePercentiles() []string {
	return []string{P1, P5, P10, P50, P90, P95, P99}
}

type Percentile struct {
	P1  TimeSeriesData `json:"p1,omitempty"`
	P5  TimeSeriesData `json:"p5,omitempty"`
	P10 TimeSeriesData `json:"p10,omitempty"`
	P50 TimeSeriesData `json:"p50,omitempty"`
	P90 TimeSeriesData `json:"p90,omitempty"`
	P95 TimeSeriesData `json:"p95,omitempty"`
	P99 TimeSeriesData `json:"p99,omitempty"`
}

type CountsSLITimeSeries struct {
	GoodCount  TimeSeriesData `json:"goodCount,omitempty"`
	TotalCount TimeSeriesData `json:"totalCount,omitempty"`
}

type InstantaneousBurnRateTimeSeries struct {
	InstantaneousBurnRate TimeSeriesData `json:"instantaneousBurnRate,omitempty"`
}

type CumulativeBurnedTimeSeries struct {
	CumulativeBurned TimeSeriesData `json:"cumulativeBurned,omitempty"`
}

// SLO History Report

type SLOHistoryReport struct {
	MetadataHolder
	TimeWindows []TimeWindowHistoryReport `json:"timewindows,omitempty"`
}

type TimeWindowHistoryReport struct {
	TimeWindow `json:"timewindow,omitempty"`
	Thresholds []ThresholdHistoryReport `json:"objectives,omitempty"`
}

type ThresholdHistoryReport struct {
	ThresholdBase
	BurnDownTimeSeries
}

// Common

type TimeSeriesData [][]interface{}

type BurnDownTimeSeries struct {
	BurnDown []TimeSeriesData `json:"burnDown,omitempty"`
}

// UnsupportedKindErr returns appropriate error for missing value in field kind
// for not empty field kind returns always that is not supported for this apiVersion
// so have to be validated before
func UnsupportedKindErr(o ObjectGeneric) error {
	if strings.TrimSpace(o.Kind) == "" {
		return EnhanceError(o, errors.New("missing or empty field kind for an Object"))
	}
	return EnhanceError(o, fmt.Errorf("invalid Object kind: %s for apiVersion: %s", o.Kind, o.APIVersion))
}

// ObjectInternal represents part of object which is only for internal usage,
// not exposed to the client, for internal usage
type ObjectInternal struct {
	Organization string `json:",omitempty" example:"nobl9-dev"`
	ManifestSrc  string `json:",omitempty" example:"x.yml"`
	OktaClientID string `json:"-"` // used only by kind Agent
}

// Metadata represents part of object which is is common for all available Objects, for internal usage
type Metadata struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	Project     string `json:"project,omitempty"`
	Labels      Labels `json:"labels,omitempty"`
}

type Labels map[string][]string

// MetadataHolder is an intermediate structure that can provides metadata related
// field to other structures
type MetadataHolder struct {
	Metadata Metadata `json:"metadata"`
}

// ObjectHeader represents Header which is common for all available Objects
type ObjectHeader struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	MetadataHolder
	ObjectInternal
}

// ObjectGeneric represents struct to which every Objects is parsable
// Specific types of Object have different structures as Spec
type ObjectGeneric struct {
	ObjectHeader
	Spec json.RawMessage `json:"spec"`
}

// EnhanceError annotates error with path of manifest source, if it exists
// if not returns the same error as passed as argument
func EnhanceError(o ObjectGeneric, err error) error {
	if err != nil && o.ManifestSrc != "" {
		err = fmt.Errorf("%s:\n%w", o.ManifestSrc, err)
	}
	return err
}

// genericToSLO converts ObjectGeneric to Object SLO
func genericToSLO(o ObjectGeneric, onlyHeader bool) (SLO, error) {
	res := SLO{
		ObjectHeader: o.ObjectHeader,
	}
	if onlyHeader {
		return res, nil
	}
	var resSpec SLOSpec
	if err := json.Unmarshal(o.Spec, &resSpec); err != nil {
		return res, EnhanceError(o, err)
	}
	res.Spec = resSpec
	return res, nil
}

// Calendar struct represents calendar time window
type Calendar struct {
	StartTime string `json:"startTime"`
	TimeZone  string `json:"timeZone"`
}

// Period represents period of time
type Period struct {
	Begin string `json:"begin"`
	End   string `json:"end"`
}

// TimeWindow represents content of time window
type TimeWindow struct {
	Unit      string    `json:"unit"`
	Count     int       `json:"count"`
	IsRolling bool      `json:"isRolling" example:"true"`
	Calendar  *Calendar `json:"calendar,omitempty"`

	// Period is only returned in `/get/slo` requests it is ignored for `/apply`
	Period *Period `json:"period"`
}

// Attachment represents user defined URL attached to SLO
type Attachment struct {
	URL         string  `json:"url"`
	DisplayName *string `json:"displayName,omitempty"`
}

// DataSourceStatus represents content of Status optional for DataSource Object
type DataSourceStatus struct {
	DataSourceType string `json:"dataSourceType" example:"Prometheus"`
}

// PrometheusConfig represents content of Prometheus Configuration typical for DataSource Object
type PrometheusConfig struct {
	URL              *string                     `json:"url,omitempty" example:"http://prometheus-service.monitoring:8080"`
	ServiceDiscovery *PrometheusServiceDiscovery `json:"serviceDiscovery,omitempty"`
}

// PrometheusServiceDiscovery provides settings for mechanism of auto Service discovery
type PrometheusServiceDiscovery struct {
	// empty is treated as once, later support 1m, 2d, etc. (for now not validated, skipped)
	Interval string                    `json:"interval,omitempty"`
	Rules    []PrometheusDiscoveryRule `json:"rules,omitempty"`
}

// PrometheusDiscoveryRule provides struct for storing rule for single Service discovery rule from Prometheus
type PrometheusDiscoveryRule struct {
	Discovery          string        `json:"discovery"`
	ServiceNamePattern string        `json:"serviceNamePattern"`
	Filter             []FilterEntry `json:"filter,omitempty"`
}

// FilterEntry represents single metric label to be matched against value
type FilterEntry struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// DatadogConfig represents content of Datadog Configuration typical for DataSource Object
type DatadogConfig struct {
	Site string `json:"site,omitempty"`
}

// DatadogAgentConfig represents content of Datadog Configuration typical for Agent Object
type DatadogAgentConfig struct {
	Site string `json:"site,omitempty"`
}

// DatadogDirectConfig represents content of Datadog Configuration typical for Direct Object
type DatadogDirectConfig struct {
	Site           string `json:"site,omitempty"`
	APIKey         string `json:"apiKey" example:"secret"`
	ApplicationKey string `json:"applicationKey" example:"secret"`
}

// NewRelicConfig represents content of NewRelic Configuration typical for DataSource Object
type NewRelicConfig struct {
	AccountID json.Number `json:"accountId,omitempty" example:"123654"`
}

// NewRelicAgentConfig represents content of NewRelic Configuration typical for Agent Object
type NewRelicAgentConfig struct {
	AccountID json.Number `json:"accountId,omitempty" example:"123654"`
}

// NewRelicDirectConfig represents content of NewRelic Configuration typical for Direct Object
type NewRelicDirectConfig struct {
	AccountID        json.Number `json:"accountId"`
	InsightsQueryKey string      `json:"insightsQueryKey" example:"secret"`
}

// AppDynamicsConfig represents content of AppDynamics Configuration typical for DataSource Object
type AppDynamicsConfig struct {
	URL string `json:"url,omitempty" example:"https://nobl9.saas.appdynamics.com"`
}

// AppDynamicsAgentConfig represents content of AppDynamics Configuration typical for Agent Object
type AppDynamicsAgentConfig struct {
	URL *string `json:"url,omitempty" example:"https://nobl9.saas.appdynamics.com"`
}

// AppDynamicsDirectConfig represents content of AppDynamics Configuration typical for Direct Object
type AppDynamicsDirectConfig struct {
	URL          string `json:"url,omitempty"`
	ClientID     string `json:"clientID,omitempty" example:"apiClientID@accountID"`
	ClientSecret string `json:"clientSecret,omitempty" example:"secret"`
}

// SplunkConfig represents content of Splunk Configuration typical for DataSource Object
type SplunkConfig struct {
	URL string `json:"url,omitempty" example:"https://localhost:8089/servicesNS/admin/"`
}

// SplunkAgentConfig represents content of Splunk Configuration typical for Agent Object
type SplunkAgentConfig struct {
	URL string `json:"url,omitempty" example:"https://localhost:8089/servicesNS/admin/"`
}

// LightstepConfig represents content of Lightstep Configuration typical for DataSource Object
type LightstepConfig struct {
	Organization string `json:"organization,omitempty" example:"LightStep-Play"`
	Project      string `json:"project,omitempty" example:"play"`
}

// LightstepAgentConfig represents content of Lightstep Configuration typical for Agent Object
type LightstepAgentConfig struct {
	Organization string `json:"organization,omitempty" example:"LightStep-Play"`
	Project      string `json:"project,omitempty" example:"play"`
}

// SplunkObservabilityAgentConfig represents content of SplunkObservability Configuration typical for Agent Object.
type SplunkObservabilityAgentConfig struct {
	Realm string `json:"realm,omitempty" example:"us1"`
}

// SplunkObservabilityDirectConfig represents content of SplunkObservability Configuration typical for Direct Object.
type SplunkObservabilityDirectConfig struct {
	Realm       string `json:"realm,omitempty"`
	AccessToken string `json:"accessToken,omitempty"`
}

// ThousandEyesDirectConfig represents content of ThousandEyes Configuration typical for Direct Object
type ThousandEyesDirectConfig struct {
	OauthBearerToken string `json:"oauthBearerToken,omitempty"`
}

// ThousandEyesAgentConfig represents content of ThousandEyes Configuration typical for Agent Object
type ThousandEyesAgentConfig struct {
	// ThousandEyes agent doesn't require any additional parameters.
}

// DynatraceAgentConfig represents content of Dynatrace Configuration typical for Agent Object.
type DynatraceAgentConfig struct {
	URL string `json:"url,omitempty"` //nolint: lll
}

// DynatraceConfig represents content of Dynatrace Configuration typical for DataSource Object.
type DynatraceConfig struct {
	URL string `json:"url,omitempty"` //nolint: lll
}

// ElasticsearchAgentConfig represents content of Elasticsearch Configuration typical for Agent Object.
type ElasticsearchAgentConfig struct {
	URL string `json:"url,omitempty"`
}

// ElasticsearchConfig represents content of Elasticsearch Configuration typical for DataSource Object.
type ElasticsearchConfig struct {
	URL string `json:"url,omitempty"`
}

// GraphiteAgentConfig represents content of Graphite Configuration typical for Agent Object.
type GraphiteAgentConfig struct {
	URL string `json:"url,omitempty"`
}

// BigQueryAgentConfig represents content of BigQuery configuration.
// Since the agent does not require additional configuration this is just a marker struct.
type BigQueryAgentConfig struct{}

type BigQueryDirectConfig struct {
	ServiceAccountKey string `json:"serviceAccountKey,omitempty"`
}

// OpenTSDBAgentConfig represents content of OpenTSDB Configuration typical for Agent Object.
type OpenTSDBAgentConfig struct {
	URL string `json:"url,omitempty" example:"example of OpenTSDB cluster URL"`
}

// GrafanaLokiAgentConfig represents content of GrafanaLoki Configuration typical for Agent Object.
type GrafanaLokiAgentConfig struct {
	URL string `json:"url,omitempty" example:"example of GrafanaLoki cluster URL"`
}

// genericToAgent converts ObjectGeneric to ObjectAgent
func genericToAgent(o ObjectGeneric, onlyHeader bool) (Agent, error) {
	res := Agent{
		ObjectHeader: o.ObjectHeader,
	}
	if onlyHeader {
		return res, nil
	}
	var resSpec AgentSpec
	if err := json.Unmarshal(o.Spec, &resSpec); err != nil {
		err = EnhanceError(o, err)
		return res, err
	}
	res.Spec = resSpec
	return res, nil
}

// Agent struct which mapped one to one with kind: Agent yaml definition
type Agent struct {
	ObjectHeader
	Spec   AgentSpec   `json:"spec"`
	Status AgentStatus `json:"status"`
}

// AgentWithSLOs struct which mapped one to one with kind: agent and slo yaml definition
type AgentWithSLOs struct {
	Agent Agent `json:"agent"`
	SLOs  []SLO `json:"slos"`
}

// AgentStatus represents content of Status optional for Agent Object
type AgentStatus struct {
	AgentType      string `json:"agentType" example:"Prometheus"`
	AgentVersion   string `json:"agentVersion,omitempty" example:"0.0.9"`
	LastConnection string `json:"lastConnection,omitempty" example:"2020-08-31T14:26:13Z"`
}

// AgentSpec represents content of Spec typical for Agent Object
type AgentSpec struct {
	Description         string                          `json:"description,omitempty" example:"Prometheus description"` //nolint:lll
	SourceOf            []string                        `json:"sourceOf" example:"Metrics,Services"`
	Prometheus          *PrometheusConfig               `json:"prometheus,omitempty"`
	Datadog             *DatadogAgentConfig             `json:"datadog,omitempty"`
	NewRelic            *NewRelicAgentConfig            `json:"newRelic,omitempty"`
	AppDynamics         *AppDynamicsAgentConfig         `json:"appDynamics,omitempty"`
	Splunk              *SplunkAgentConfig              `json:"splunk,omitempty"`
	Lightstep           *LightstepAgentConfig           `json:"lightstep,omitempty"`
	SplunkObservability *SplunkObservabilityAgentConfig `json:"splunkObservability,omitempty"`
	Dynatrace           *DynatraceAgentConfig           `json:"dynatrace,omitempty"`
	Elasticsearch       *ElasticsearchAgentConfig       `json:"elasticsearch,omitempty"`
	ThousandEyes        *ThousandEyesAgentConfig        `json:"thousandEyes,omitempty"`
	Graphite            *GraphiteAgentConfig            `json:"graphite,omitempty"`
	BigQuery            *BigQueryAgentConfig            `json:"bigQuery,omitempty"`
	OpenTSDB            *OpenTSDBAgentConfig            `json:"opentsdb,omitempty"`
	GrafanaLoki         *GrafanaLokiAgentConfig         `json:"grafanaLoki,omitempty"`
}

// genericToDirect converts ObjectGeneric to ObjectDirect
func genericToDirect(o ObjectGeneric, onlyHeader bool) (Direct, error) {
	res := Direct{
		ObjectHeader: o.ObjectHeader,
	}
	if onlyHeader {
		return res, nil
	}
	var resSpec DirectSpec
	if err := json.Unmarshal(o.Spec, &resSpec); err != nil {
		err = EnhanceError(o, err)
		return res, err
	}
	res.Spec = resSpec
	return res, nil
}

// Direct struct which mapped one to one with kind: Direct yaml definition
type Direct struct {
	ObjectHeader
	Spec   DirectSpec   `json:"spec"`
	Status DirectStatus `json:"status"`
}

// DirectSpec represents content of Spec typical for Direct Object
type DirectSpec struct {
	Description         string                           `json:"description,omitempty" example:"Datadog description"` //nolint:lll
	SourceOf            []string                         `json:"sourceOf" example:"Metrics,Services"`
	Datadog             *DatadogDirectConfig             `json:"datadog,omitempty"`
	NewRelic            *NewRelicDirectConfig            `json:"newRelic,omitempty"`
	AppDynamics         *AppDynamicsDirectConfig         `json:"appDynamics,omitempty"`
	SplunkObservability *SplunkObservabilityDirectConfig `json:"splunkObservability,omitempty"`
	ThousandEyes        *ThousandEyesDirectConfig        `json:"thousandEyes,omitempty"`
	BigQuery            *BigQueryDirectConfig            `json:"bigQuery,omitempty"`
}

// DirectStatus represents content of Status optional for Direct Object
type DirectStatus struct {
	DirectType string `json:"directType" example:"Datadog"`
}

// Service struct which mapped one to one with kind: service yaml definition
type Service struct {
	ObjectHeader
	Spec   ServiceSpec   `json:"spec"`
	Status ServiceStatus `json:"status"`
}

// ServiceWithSLOs struct which mapped one to one with kind: service and slo yaml definition
type ServiceWithSLOs struct {
	Service Service `json:"service"`
	SLOs    []SLO   `json:"slos"`
}

// ServiceStatus represents content of Status optional for Service Object
type ServiceStatus struct {
	SloCount int `json:"sloCount"`
}

// ServiceSpec represents content of Spec typical for Service Object
type ServiceSpec struct {
	Description string `json:"description"`
}

// Project struct which mapped one to one with kind: project yaml definition.
type Project struct {
	ObjectInternal
	APIVersion string          `json:"apiVersion"`
	Kind       string          `json:"kind"`
	Metadata   ProjectMetadata `json:"metadata"`
	Spec       ProjectSpec     `json:"spec"`
}

type ProjectMetadata struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
}

// ProjectSpec represents content of Spec typical for Project Object.
type ProjectSpec struct {
	Description string `json:"description"`
}

// genericToService converts ObjectGeneric to Object Service
func genericToService(o ObjectGeneric, onlyHeader bool) (Service, error) {
	res := Service{
		ObjectHeader: o.ObjectHeader,
	}
	if onlyHeader {
		return res, nil
	}

	var resSpec ServiceSpec
	if err := json.Unmarshal(o.Spec, &resSpec); err != nil {
		err = EnhanceError(o, err)
		return res, err
	}
	res.Spec = resSpec
	return res, nil
}

// AlertPolicy represents a set of conditions that can trigger an alert.
type AlertPolicy struct {
	ObjectHeader
	Spec AlertPolicySpec `json:"spec"`
}

// AlertPolicySpec represents content of AlertPolicy's Spec.
type AlertPolicySpec struct {
	Description  string              `json:"description"`
	Severity     string              `json:"severity"`
	Conditions   []AlertCondition    `json:"conditions"`
	AlertMethods []PublicAlertMethod `json:"alertMethods"`
}

func (spec AlertPolicySpec) GetAlertMethods() []PublicAlertMethod {
	return spec.AlertMethods
}

// AlertCondition represents a condition to meet to trigger an alert.
type AlertCondition struct {
	Measurement      string      `json:"measurement"`
	Value            interface{} `json:"value"`
	LastsForDuration string      `json:"lastsFor,omitempty"` //nolint:lll
	CoolDownDuration string      `json:"coolDown,omitempty"` //nolint:lll
	Operation        string      `json:"op"`
}

// AlertPolicyWithSLOs struct which mapped one to one with kind: alert policy and slo yaml definition
type AlertPolicyWithSLOs struct {
	AlertPolicy AlertPolicy `json:"alertPolicy"`
	SLOs        []SLO       `json:"slos"`
}

// AlertMethodAssignment represents an AlertMethod assigned to AlertPolicy.
type AlertMethodAssignment struct {
	Project string `json:"project,omitempty"`
	Name    string `json:"name"`
}

type AlertMethodWithAlertPolicy struct {
	AlertMethod   PublicAlertMethod `json:"alertMethod"`
	AlertPolicies []AlertPolicy     `json:"alertPolicies"`
}

// genericToAlertPolicy converts ObjectGeneric to ObjectAlertPolicy
func genericToAlertPolicy(o ObjectGeneric, onlyHeader bool) (AlertPolicy, error) {
	res := AlertPolicy{
		ObjectHeader: o.ObjectHeader,
	}
	if onlyHeader {
		return res, nil
	}
	var resSpec AlertPolicySpec
	if err := json.Unmarshal(o.Spec, &resSpec); err != nil {
		err = EnhanceError(o, err)
		return res, err
	}
	res.Spec = resSpec
	return res, nil
}

// Alert represents triggered alert
type Alert struct {
	ObjectHeader
	Spec AlertSpec `json:"spec"`
}

// AlertSpec represents content of Alert's Spec
type AlertSpec struct {
	AlertPolicy    Metadata `json:"alertPolicy"`
	SLO            Metadata `json:"slo"`
	Service        Metadata `json:"service"`
	ThresholdValue float64  `json:"thresholdValue,omitempty"`
	ClockTime      string   `json:"clockTime,omitempty"`
	Severity       string   `json:"severity"`
}

// AlertMethod represents the configuration required to send a notification to an external service
// when an alert is triggered.
type AlertMethod struct {
	ObjectHeader
	Spec AlertMethodSpec `json:"spec"`
}

// PublicAlertMethod represents the configuration required to send a notification to an external service
// when an alert is triggered.
type PublicAlertMethod struct {
	ObjectHeader
	Spec   PublicAlertMethodSpec    `json:"spec"`
	Status *PublicAlertMethodStatus `json:"status,omitempty"`
}

// PublicAlertMethodStatus represents content of Status optional for PublicAlertMethod Object
type PublicAlertMethodStatus struct {
	LastTestDate       string `json:"lastTestDate,omitempty"`
	NextTestPossibleAt string `json:"nextTestPossibleAt,omitempty"`
}

// AlertMethodSpec represents content of AlertMethod's Spec.
type AlertMethodSpec struct {
	Description string                 `json:"description"`
	Webhook     *WebhookAlertMethod    `json:"webhook,omitempty"`
	PagerDuty   *PagerDutyAlertMethod  `json:"pagerduty,omitempty"`
	Slack       *SlackAlertMethod      `json:"slack,omitempty"`
	Discord     *DiscordAlertMethod    `json:"discord,omitempty"`
	Opsgenie    *OpsgenieAlertMethod   `json:"opsgenie,omitempty"`
	ServiceNow  *ServiceNowAlertMethod `json:"servicenow,omitempty"`
	Jira        *JiraAlertMethod       `json:"jira,omitempty"`
	Teams       *TeamsAlertMethod      `json:"msteams,omitempty"`
	Email       *EmailAlertMethod      `json:"email,omitempty"`
}

// PublicAlertMethodSpec represents content of AlertMethod's Spec without secrets.
type PublicAlertMethodSpec struct {
	Description string                       `json:"description"`
	Webhook     *PublicWebhookAlertMethod    `json:"webhook,omitempty"`
	PagerDuty   *PublicPagerDutyAlertMethod  `json:"pagerduty,omitempty"`
	Slack       *PublicSlackAlertMethod      `json:"slack,omitempty"`
	Discord     *PublicDiscordAlertMethod    `json:"discord,omitempty"`
	Opsgenie    *PublicOpsgenieAlertMethod   `json:"opsgenie,omitempty"`
	ServiceNow  *PublicServiceNowAlertMethod `json:"servicenow,omitempty"`
	Jira        *PublicJiraAlertMethod       `json:"jira,omitempty"`
	Teams       *PublicTeamsAlertMethod      `json:"msteams,omitempty"`
	Email       *EmailAlertMethod            `json:"email,omitempty"`
}

// WebhookAlertMethod represents a set of properties required to send a webhook request.
type WebhookAlertMethod struct {
	URL            string          `json:"url"` // Field required when AlertMethod is created.
	Template       *string         `json:"template,omitempty"`
	TemplateFields []string        `json:"templateFields,omitempty"` //nolint:lll
	Headers        []WebhookHeader `json:"headers,omitempty"`
}
type WebhookHeader struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	IsSecret bool   `json:"isSecret"`
}

// PublicWebhookAlertMethod represents a set of properties required to send a webhook request without secrets.
type PublicWebhookAlertMethod struct {
	HiddenURL      string          `json:"url"`
	Template       *string         `json:"template,omitempty"`
	TemplateFields []string        `json:"templateFields,omitempty"` //nolint:lll
	Headers        []WebhookHeader `json:"headers,omitempty"`
}

// SendResolution If user set SendResolution, then “Send a notification after the cooldown period is over"
type SendResolution struct {
	Message *string `json:"message"`
}

// PagerDutyAlertMethod represents a set of properties required to open an Incident in PagerDuty.
type PagerDutyAlertMethod struct {
	IntegrationKey string          `json:"integrationKey"`
	SendResolution *SendResolution `json:"sendResolution,omitempty"`
}

// PublicPagerDutyAlertMethod represents a set of properties required to open an Incident in PagerDuty without secrets.
type PublicPagerDutyAlertMethod struct {
	HiddenIntegrationKey string          `json:"integrationKey"`
	SendResolution       *SendResolution `json:"sendResolution,omitempty"`
}

// SlackAlertMethod represents a set of properties required to send message to Slack.
type SlackAlertMethod struct {
	URL string `json:"url"` // Required when AlertMethod is created.
}

// PublicSlackAlertMethod represents a set of properties required to send message to Slack without secrets.
type PublicSlackAlertMethod struct {
	HiddenURL string `json:"url"`
}

// OpsgenieAlertMethod represents a set of properties required to send message to Opsgenie.
type OpsgenieAlertMethod struct {
	Auth string `json:"auth"` // Field required when AlertMethod is created.
	URL  string `json:"url"`
}

// PublicOpsgenieAlertMethod represents a set of properties required to send message to Opsgenie without secrets.
type PublicOpsgenieAlertMethod struct {
	HiddenAuth string `json:"auth"`
	URL        string `json:"url"`
}

// ServiceNowAlertMethod represents a set of properties required to send message to ServiceNow.
type ServiceNowAlertMethod struct {
	Username   string `json:"username"`
	Password   string `json:"password"` // Field required when AlertMethod is created.
	InstanceID string `json:"instanceid"`
}

// PublicServiceNowAlertMethod represents a set of properties required to send message to ServiceNow without secrets.
type PublicServiceNowAlertMethod struct {
	Username       string `json:"username"`
	InstanceID     string `json:"instanceid"`
	HiddenPassword string `json:"password"`
}

// DiscordAlertMethod represents a set of properties required to send message to Discord.
type DiscordAlertMethod struct {
	URL string `json:"url"` // Field required when AlertMethod is created.
}

// PublicDiscordAlertMethod represents a set of properties required to send message to Discord without secrets.
type PublicDiscordAlertMethod struct {
	HiddenURL string `json:"url"`
}

// JiraAlertMethod represents a set of properties required create tickets in Jira.
type JiraAlertMethod struct {
	URL        string `json:"url"`
	Username   string `json:"username"`
	APIToken   string `json:"apiToken"` // Field required when AlertMethod is created.
	ProjectKey string `json:"projectKey"`
}

// PublicJiraAlertMethod represents a set of properties required create tickets in Jira without secrets.
type PublicJiraAlertMethod struct {
	URL            string `json:"url"`
	Username       string `json:"username"`
	ProjectKey     string `json:"projectKey"`
	HiddenAPIToken string `json:"apiToken"`
}

// TeamsAlertMethod represents a set of properties required create Microsoft Teams notifications.
type TeamsAlertMethod struct {
	URL string `json:"url"`
}

// PublicTeamsAlertMethod represents a set of properties required create Microsoft Teams notifications.
type PublicTeamsAlertMethod struct {
	HiddenURL string `json:"url"`
}

// EmailAlertMethod represents a set of properties required to send an email.
type EmailAlertMethod struct {
	To      []string `json:"to,omitempty"`
	Cc      []string `json:"cc,omitempty"`
	Bcc     []string `json:"bcc,omitempty"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

// genericToAlertMethod converts ObjectGeneric to ObjectAlertMethod
func genericToAlertMethod(o ObjectGeneric, onlyHeader bool) (AlertMethod, error) {
	res := AlertMethod{
		ObjectHeader: o.ObjectHeader,
	}
	if onlyHeader {
		return res, nil
	}
	var resSpec AlertMethodSpec
	if err := json.Unmarshal(o.Spec, &resSpec); err != nil {
		err = EnhanceError(o, err)
		return res, err
	}
	res.Spec = resSpec
	return res, nil
}

// DataExport struct which mapped one to one with kind: DataExport yaml definition
type DataExport struct {
	ObjectHeader
	Spec   DataExportSpec   `json:"spec"`
	Status DataExportStatus `json:"status"`
}

// DataExportSpec represents content of DataExport's Spec
type DataExportSpec struct {
	ExportType string      `json:"exportType"`
	Spec       interface{} `json:"spec"`
}

// S3DataExportSpec represents content of Amazon S3 export type spec.
type S3DataExportSpec struct {
	BucketName string `json:"bucketName"`
	RoleARN    string `json:"roleArn"` //nolint:lll
}

// GCSDataExportSpec represents content of GCP Cloud Storage export type spec.
type GCSDataExportSpec struct {
	BucketName string `json:"bucketName"`
}

// DataExportStatus represents content of Status optional for DataExport Object
type DataExportStatus struct {
	ExportJob     DataExportStatusJob `json:"exportJob"`
	AWSExternalID *string             `json:"awsExternalID,omitempty"`
}

// DataExportStatusJob represents content of ExportJob status
type DataExportStatusJob struct {
	Timestamp string `json:"timestamp,omitempty" example:"2021-02-09T10:43:07Z"`
	State     string `json:"state" example:"finished"`
}

// dataExportGeneric represents struct to which every DataExport is parsable.
// Specific types of DataExport have different structures as Spec.
type dataExportGeneric struct {
	ExportType string          `json:"exportType"`
	Spec       json.RawMessage `json:"spec"`
}

// genericToDataExport converts ObjectGeneric to ObjectDataExport
func genericToDataExport(o ObjectGeneric, onlyHeader bool) (DataExport, error) {
	res := DataExport{
		ObjectHeader: o.ObjectHeader,
	}
	if onlyHeader {
		return res, nil
	}
	deg := dataExportGeneric{}
	if err := json.Unmarshal(o.Spec, &deg); err != nil {
		err = EnhanceError(o, err)
		return res, err
	}

	resSpec := DataExportSpec{ExportType: deg.ExportType}
	switch resSpec.ExportType {
	case "S3", "Snowflake":
		resSpec.Spec = &S3DataExportSpec{}
	case "GCS":
		resSpec.Spec = &GCSDataExportSpec{}
	}
	if deg.Spec != nil {
		if err := json.Unmarshal(deg.Spec, &resSpec.Spec); err != nil {
			err = EnhanceError(o, err)
			return res, err
		}
	}
	res.Spec = resSpec
	return res, nil
}

// genericToProject converts ObjectGeneric to Project
func genericToProject(o ObjectGeneric, onlyHeader bool) (Project, error) {
	res := Project{
		APIVersion: o.ObjectHeader.APIVersion,
		Kind:       o.ObjectHeader.Kind,
		Metadata: ProjectMetadata{
			Name:        o.Metadata.Name,
			DisplayName: o.Metadata.DisplayName,
		},
		ObjectInternal: ObjectInternal{
			Organization: o.ObjectHeader.Organization,
			ManifestSrc:  o.ObjectHeader.ManifestSrc,
		},
	}
	if onlyHeader {
		return res, nil
	}

	var resSpec ProjectSpec
	if err := json.Unmarshal(o.Spec, &resSpec); err != nil {
		err = EnhanceError(o, err)
		return res, err
	}
	res.Spec = resSpec

	return res, nil
}

// RoleBinding represents relation of User and Role
type RoleBinding struct {
	ObjectInternal
	APIVersion string              `json:"apiVersion"`
	Kind       string              `json:"kind"`
	Metadata   RoleBindingMetadata `json:"metadata"`
	Spec       RoleBindingSpec     `json:"spec"`
}

type RoleBindingSpec struct {
	User       string `json:"user"`
	RoleRef    string `json:"roleRef"`
	ProjectRef string `json:"projectRef,omitempty"`
}

type RoleBindingMetadata struct {
	Name string `json:"name"`
}

// genericToRoleBinding converts ObjectGeneric to ObjectRoleBinding
// onlyHeader parameter is not supported for RoleBinding since ProjectRef is defined on Spec section.
func genericToRoleBinding(o ObjectGeneric) (RoleBinding, error) {
	res := RoleBinding{
		APIVersion: o.ObjectHeader.APIVersion,
		Kind:       o.ObjectHeader.Kind,
		Metadata: RoleBindingMetadata{
			Name: o.Metadata.Name,
		},
		ObjectInternal: ObjectInternal{
			Organization: o.ObjectHeader.Organization,
			ManifestSrc:  o.ObjectHeader.ManifestSrc,
		},
	}
	var resSpec RoleBindingSpec
	if err := json.Unmarshal(o.Spec, &resSpec); err != nil {
		err = EnhanceError(o, err)
		return res, err
	}
	res.Spec = resSpec

	return res, nil
}

// Parse takes care of all Object supported by n9/v1alpha apiVersion
func Parse(o ObjectGeneric, parsedObjects *APIObjects, onlyHeaders bool) error {
	var allErrors []string
	switch o.Kind {
	case KindSLO:
		slo, err := genericToSLO(o, onlyHeaders)
		if err != nil {
			allErrors = append(allErrors, err.Error())
		}
		parsedObjects.SLOs = append(parsedObjects.SLOs, slo)
	case KindService:
		service, err := genericToService(o, onlyHeaders)
		if err != nil {
			allErrors = append(allErrors, err.Error())
		}
		parsedObjects.Services = append(parsedObjects.Services, service)
	case KindAgent:
		agent, err := genericToAgent(o, onlyHeaders)
		if err != nil {
			allErrors = append(allErrors, err.Error())
		}
		parsedObjects.Agents = append(parsedObjects.Agents, agent)
	case KindAlertPolicy:
		alertPolicy, err := genericToAlertPolicy(o, onlyHeaders)
		if err != nil {
			allErrors = append(allErrors, err.Error())
		}
		parsedObjects.AlertPolicies = append(parsedObjects.AlertPolicies, alertPolicy)
	case KindAlertMethod:
		alertMethod, err := genericToAlertMethod(o, onlyHeaders)
		if err != nil {
			allErrors = append(allErrors, err.Error())
		}
		parsedObjects.AlertMethods = append(parsedObjects.AlertMethods, alertMethod)
	case KindDirect:
		direct, err := genericToDirect(o, onlyHeaders)
		if err != nil {
			allErrors = append(allErrors, err.Error())
		}
		parsedObjects.Directs = append(parsedObjects.Directs, direct)
	case KindDataExport:
		dataExport, err := genericToDataExport(o, onlyHeaders)
		if err != nil {
			allErrors = append(allErrors, err.Error())
		}
		parsedObjects.DataExports = append(parsedObjects.DataExports, dataExport)
	case KindProject:
		project, err := genericToProject(o, onlyHeaders)
		if err != nil {
			allErrors = append(allErrors, err.Error())
		}
		parsedObjects.Projects = append(parsedObjects.Projects, project)
	case KindRoleBinding:
		roleBinding, err := genericToRoleBinding(o)
		if err != nil {
			allErrors = append(allErrors, err.Error())
		}
		parsedObjects.RoleBindings = append(parsedObjects.RoleBindings, roleBinding)
	// catching invalid kinds of objects for this apiVersion
	default:
		err := UnsupportedKindErr(o)
		allErrors = append(allErrors, err.Error())
	}
	if len(allErrors) > 0 {
		return errors.New(strings.Join(allErrors, "\n"))
	}
	return nil
}
