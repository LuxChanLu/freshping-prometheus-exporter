package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const metricsNamespace = "freshping"
const reportSubsystem = "report"
const timingsSubsystem = "report"

var labels = []string{"check_id", "check_name"}

var (
	// Basic report metrics
	reportDurationSecondsAvailable                      = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "duration_seconds_available"}, labels)
	reportDurationSecondsReportingError                 = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "duration_seconds_reporting_error"}, labels)
	reportDurationSecondsNotResponding                  = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "duration_seconds_not_responding"}, labels)
	reportOutagesCountReportingError                    = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "outages_count_reporting_error"}, labels)
	reportOutagesCountNotResponding                     = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "outages_count_not_responding"}, labels)
	reportAverageResponseTimeMillisecondsAvailable      = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "average_response_time_milliseconds_available"}, labels)
	reportAverageResponseTimeMillisecondsReportingError = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "average_response_time_milliseconds_reporting_error"}, labels)
	reportMinimumResponseTimeMillisecondsAvailable      = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "minimum_response_time_milliseconds_available"}, labels)
	reportMinimumResponseTimeMillisecondsReportingError = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "minimum_response_time_milliseconds_reporting_error"}, labels)
	reportMaximumResponseTimeMillisecondsAvailable      = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "maximum_response_time_milliseconds_available"}, labels)
	reportMaximumResponseTimeMillisecondsReportingError = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "maximum_response_time_milliseconds_reporting_error"}, labels)
	reportDurationSecondsPerformanceGood                = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "duration_seconds_performance_good"}, labels)
	reportDurationSecondsPerformanceDegraded            = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "duration_seconds_performance_degraded"}, labels)
	reportCountPerformanceDegradations                  = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "count_performance_degradations"}, labels)
	reportApdexResultsCountSatisfied                    = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "apdex_results_count_satisfied"}, labels)
	reportApdexResultsCountTolerating                   = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "apdex_results_count_tolerating"}, labels)
	reportApdexResultsCountFrustrated                   = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "apdex_results_count_frustrated"}, labels)
	reportAvailabilityPercentage                        = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "availability_percentage"}, labels)
	reportDurationSecondsTotalDowntime                  = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "duration_seconds_total_downtime"}, labels)
	reportOutagesCountTotal                             = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "outages_count_total"}, labels)
	reportApdexScore                                    = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: metricsNamespace, Subsystem: reportSubsystem, Name: "apdex_score"}, labels)
)
