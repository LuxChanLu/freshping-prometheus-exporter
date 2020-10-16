package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
)

func (exporter *FreshpingExporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// List fresphing checks
	res, err := http.Get(exporter.freshpingURL.String())
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	statusPage := FreshpingStatusPage{}
	if err := json.Unmarshal([]byte(doc.Find("script[id='__NEXT_DATA__']").Contents().Text()), &statusPage); err != nil {
		log.Println(err)
		return
	}
	// Get metrics details
	g, _ := errgroup.WithContext(context.Background())
	for _, check := range statusPage.Props.PageProps.Checks.AssociatedChecks {
		check := check
		g.Go(func() error {
			// Basic report
			resReport, err := http.Get("https://api.freshping.io/v1/public-check-stats-reports/" + strconv.Itoa(check.CheckID))
			if err != nil {
				return err
			}
			defer resReport.Body.Close()
			if resReport.StatusCode != 200 {
				return err
			}
			resReportBytes, err := ioutil.ReadAll(resReport.Body)
			if err != nil {
				return err
			}
			report := FreshpingCheckReport{}
			if err := json.Unmarshal(resReportBytes, &report); err != nil {
				return err
			}

			// Transform to prometheus metrics
			availabilityPercentage, err := strconv.ParseFloat(report.AvailabilityPercentage, 64)
			if err != nil {
				return err
			}
			apdexScore, err := strconv.ParseFloat(report.ApdexScore, 64)
			if err != nil {
				return err
			}
			labels := []string{
				strconv.Itoa(statusPage.Props.PageProps.Checks.OrganizationID),
				statusPage.Props.PageProps.Checks.OrganizationName,
				strconv.Itoa(statusPage.Props.PageProps.Checks.StatusPageID),
				statusPage.Props.PageProps.Checks.StatusPageName,
				strconv.Itoa(check.CheckID),
				check.CheckName,
			}

			reportDurationSecondsAvailable.WithLabelValues(labels...).Set(report.DurationSecondsAvailable)
			reportDurationSecondsReportingError.WithLabelValues(labels...).Set(report.DurationSecondsReportingError)
			reportDurationSecondsNotResponding.WithLabelValues(labels...).Set(report.DurationSecondsNotResponding)
			reportOutagesCountReportingError.WithLabelValues(labels...).Set(report.OutagesCountReportingError)
			reportOutagesCountNotResponding.WithLabelValues(labels...).Set(report.OutagesCountNotResponding)
			reportAverageResponseTimeMillisecondsAvailable.WithLabelValues(labels...).Set(report.AverageResponseTimeMillisecondsAvailable)
			reportAverageResponseTimeMillisecondsReportingError.WithLabelValues(labels...).Set(report.AverageResponseTimeMillisecondsReportingError)
			reportMinimumResponseTimeMillisecondsAvailable.WithLabelValues(labels...).Set(report.MinimumResponseTimeMillisecondsAvailable)
			reportMinimumResponseTimeMillisecondsReportingError.WithLabelValues(labels...).Set(report.MinimumResponseTimeMillisecondsReportingError)
			reportMaximumResponseTimeMillisecondsAvailable.WithLabelValues(labels...).Set(report.MaximumResponseTimeMillisecondsAvailable)
			reportMaximumResponseTimeMillisecondsReportingError.WithLabelValues(labels...).Set(report.MaximumResponseTimeMillisecondsReportingError)
			reportDurationSecondsPerformanceGood.WithLabelValues(labels...).Set(report.DurationSecondsPerformanceGood)
			reportDurationSecondsPerformanceDegraded.WithLabelValues(labels...).Set(report.DurationSecondsPerformanceDegraded)
			reportCountPerformanceDegradations.WithLabelValues(labels...).Set(report.CountPerformanceDegradations)
			reportApdexResultsCountSatisfied.WithLabelValues(labels...).Set(report.ApdexResultsCountSatisfied)
			reportApdexResultsCountTolerating.WithLabelValues(labels...).Set(report.ApdexResultsCountTolerating)
			reportApdexResultsCountFrustrated.WithLabelValues(labels...).Set(report.ApdexResultsCountFrustrated)
			reportAvailabilityPercentage.WithLabelValues(labels...).Set(availabilityPercentage)
			reportDurationSecondsTotalDowntime.WithLabelValues(labels...).Set(report.DurationSecondsTotalDowntime)
			reportOutagesCountTotal.WithLabelValues(labels...).Set(report.OutagesCountTotal)
			reportApdexScore.WithLabelValues(labels...).Set(apdexScore)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		log.Println(err)
		return
	}
	// Export metrics
	exporter.handler.ServeHTTP(w, r)
}

func main() {
	freshpingURL := os.Getenv("FRESHPING_URL")
	if freshpingURL == "" {
		urlFile := os.Getenv("FRESHPING_URL_FILE")
		if urlFile != "" {
			content, err := ioutil.ReadFile(urlFile)
			if err != nil {
				log.Fatalln(err)
				return
			}
			freshpingURL = strings.TrimSpace(string(content))
		}
	}
	if freshpingURL == "" {
		log.Fatalln("Missing FRESHPING_URL or FRESHPING_URL_FILE")
		return
	}
	urlParsed, err := url.Parse(freshpingURL)
	if err != nil {
		log.Fatal(err)
		return
	}
	registry := prometheus.NewRegistry()
	registry.MustRegister(reportDurationSecondsAvailable)
	registry.MustRegister(reportDurationSecondsReportingError)
	registry.MustRegister(reportDurationSecondsNotResponding)
	registry.MustRegister(reportOutagesCountReportingError)
	registry.MustRegister(reportOutagesCountNotResponding)
	registry.MustRegister(reportAverageResponseTimeMillisecondsAvailable)
	registry.MustRegister(reportAverageResponseTimeMillisecondsReportingError)
	registry.MustRegister(reportMinimumResponseTimeMillisecondsAvailable)
	registry.MustRegister(reportMinimumResponseTimeMillisecondsReportingError)
	registry.MustRegister(reportMaximumResponseTimeMillisecondsAvailable)
	registry.MustRegister(reportMaximumResponseTimeMillisecondsReportingError)
	registry.MustRegister(reportDurationSecondsPerformanceGood)
	registry.MustRegister(reportDurationSecondsPerformanceDegraded)
	registry.MustRegister(reportCountPerformanceDegradations)
	registry.MustRegister(reportApdexResultsCountSatisfied)
	registry.MustRegister(reportApdexResultsCountTolerating)
	registry.MustRegister(reportApdexResultsCountFrustrated)
	registry.MustRegister(reportAvailabilityPercentage)
	registry.MustRegister(reportDurationSecondsTotalDowntime)
	registry.MustRegister(reportOutagesCountTotal)
	registry.MustRegister(reportApdexScore)
	log.Println("Server started and listen on 9705")
	http.ListenAndServe(":9705", &FreshpingExporter{
		freshpingURL: urlParsed,
		handler:      promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
	})
}
