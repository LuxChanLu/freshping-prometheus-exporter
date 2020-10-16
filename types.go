package main

import (
	"net/http"
	"net/url"
	"time"
)

// FreshpingExporter global exporter structure
type FreshpingExporter struct {
	freshpingURL *url.URL
	handler      http.Handler
}

// FreshpingStatusPage list of check
type FreshpingStatusPage struct {
	Props struct {
		PageProps struct {
			Checks struct {
				StatusPageID     int    `json:"status_page_id"`
				StatusPageName   string `json:"status_page_name"`
				OrganizationID   int    `json:"organization_id"`
				OrganizationName string `json:"organization_name"`
				LogoURL          string `json:"logo_url"`
				AssociatedChecks []struct {
					ApplicationID                 int       `json:"application_id"`
					CheckID                       int       `json:"check_id"`
					OldestAvailableCheckStartTime time.Time `json:"oldest_available_check_start_time"`
					CheckName                     string    `json:"check_name"`
					RecentAvailableCheckStartTime time.Time `json:"recent_available_check_start_time"`
					LastCheckStateChangedTime     time.Time `json:"last_check_state_changed_time"`
					ApplicationName               string    `json:"application_name"`
					CheckState                    string    `json:"check_state"`
				} `json:"associated_checks"`
				GaTrackingID string `json:"ga_tracking_id"`
				ForceSsl     bool   `json:"force_ssl"`
				BrandingData string `json:"branding_data"`
			} `json:"checks"`
			UserReferer string `json:"userReferer"`
			Sid         string `json:"sid"`
			Apihost     string `json:"apihost"`
			FullURL     string `json:"fullUrl"`
		} `json:"pageProps"`
		NSSP bool `json:"__N_SSP"`
	} `json:"props"`
	Page  string `json:"page"`
	Query struct {
		Sid     string `json:"sid"`
		Apihost string `json:"apihost"`
		FullURL string `json:"fullUrl"`
	} `json:"query"`
	BuildID      string `json:"buildId"`
	IsFallback   bool   `json:"isFallback"`
	Gssp         bool   `json:"gssp"`
	CustomServer bool   `json:"customServer"`
}

// FreshpingCheckReport one check report
type FreshpingCheckReport struct {
	CheckID                                       int     `json:"check_id"`
	DurationSecondsAvailable                      float64 `json:"duration_seconds_available"`
	DurationSecondsReportingError                 float64 `json:"duration_seconds_reporting_error"`
	DurationSecondsNotResponding                  float64 `json:"duration_seconds_not_responding"`
	OutagesCountReportingError                    float64 `json:"outages_count_reporting_error"`
	OutagesCountNotResponding                     float64 `json:"outages_count_not_responding"`
	AverageResponseTimeMillisecondsAvailable      float64 `json:"average_response_time_milliseconds_available"`
	AverageResponseTimeMillisecondsReportingError float64 `json:"average_response_time_milliseconds_reporting_error"`
	MinimumResponseTimeMillisecondsAvailable      float64 `json:"minimum_response_time_milliseconds_available"`
	MinimumResponseTimeMillisecondsReportingError float64 `json:"minimum_response_time_milliseconds_reporting_error"`
	MaximumResponseTimeMillisecondsAvailable      float64 `json:"maximum_response_time_milliseconds_available"`
	MaximumResponseTimeMillisecondsReportingError float64 `json:"maximum_response_time_milliseconds_reporting_error"`
	DurationSecondsPerformanceGood                float64 `json:"duration_seconds_performance_good"`
	DurationSecondsPerformanceDegraded            float64 `json:"duration_seconds_performance_degraded"`
	CountPerformanceDegradations                  float64 `json:"count_performance_degradations"`
	ApdexResultsCountSatisfied                    float64 `json:"apdex_results_count_satisfied"`
	ApdexResultsCountTolerating                   float64 `json:"apdex_results_count_tolerating"`
	ApdexResultsCountFrustrated                   float64 `json:"apdex_results_count_frustrated"`
	AvailabilityPercentage                        string  `json:"availability_percentage"`
	DurationSecondsTotalDowntime                  float64 `json:"duration_seconds_total_downtime"`
	OutagesCountTotal                             float64 `json:"outages_count_total"`
	ApdexScore                                    string  `json:"apdex_score"`
}

// FreshpingTimingsReport one check report
type FreshpingTimingsReport struct {
	ReportStart   time.Time `json:"report_start"`
	ResponseTimes []struct {
		DurationSecondsReportingError                 float64   `json:"duration_seconds_reporting_error"`
		AverageResponseTimeMillisecondsAvailable      float64   `json:"average_response_time_milliseconds_available"`
		DurationSecondsNotResponding                  float64   `json:"duration_seconds_not_responding"`
		End                                           time.Time `json:"end"`
		Start                                         time.Time `json:"start"`
		DurationSecondsPaused                         float64   `json:"duration_seconds_paused"`
		DurationSecondsAvailable                      float64   `json:"duration_seconds_available"`
		AverageResponseTimeMillisecondsReportingError float64   `json:"average_response_time_milliseconds_reporting_error"`
	} `json:"response_times"`
	ReportEnd   time.Time `json:"report_end"`
	AggregateBy string    `json:"aggregate_by"`
}
