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
			reportDurationSecondsAvailable.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(report.DurationSecondsAvailable)
			reportDurationSecondsReportingError.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(report.DurationSecondsReportingError)
			reportDurationSecondsNotResponding.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(report.DurationSecondsNotResponding)
			reportOutagesCountReportingError.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(report.OutagesCountReportingError)
			reportOutagesCountNotResponding.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(report.OutagesCountNotResponding)
			reportAverageResponseTimeMillisecondsAvailable.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(report.AverageResponseTimeMillisecondsAvailable)
			reportAverageResponseTimeMillisecondsReportingError.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(report.AverageResponseTimeMillisecondsReportingError)
			reportMinimumResponseTimeMillisecondsAvailable.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(report.MinimumResponseTimeMillisecondsAvailable)
			reportMinimumResponseTimeMillisecondsReportingError.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(report.MinimumResponseTimeMillisecondsReportingError)
			reportMaximumResponseTimeMillisecondsAvailable.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(report.MaximumResponseTimeMillisecondsAvailable)
			reportMaximumResponseTimeMillisecondsReportingError.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(report.MaximumResponseTimeMillisecondsReportingError)
			reportDurationSecondsPerformanceGood.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(report.DurationSecondsPerformanceGood)
			reportDurationSecondsPerformanceDegraded.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(report.DurationSecondsPerformanceDegraded)
			reportCountPerformanceDegradations.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(report.CountPerformanceDegradations)
			reportApdexResultsCountSatisfied.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(report.ApdexResultsCountSatisfied)
			reportApdexResultsCountTolerating.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(report.ApdexResultsCountTolerating)
			reportApdexResultsCountFrustrated.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(report.ApdexResultsCountFrustrated)
			reportAvailabilityPercentage.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(availabilityPercentage)
			reportDurationSecondsTotalDowntime.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(report.DurationSecondsTotalDowntime)
			reportOutagesCountTotal.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(report.OutagesCountTotal)
			reportApdexScore.WithLabelValues(strconv.Itoa(check.CheckID), check.CheckName).Set(apdexScore)
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
			freshpingURL = string(content)
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
