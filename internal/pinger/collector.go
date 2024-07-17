package pinger

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type pingCollector struct {
	responseTimeMetric *prometheus.Desc
	responseCodeMetric *prometheus.Desc
	statusMetric       *prometheus.Desc
	authToken          string
	targets            []string
	locations          []string
}

// NewPingCollector takes in ping targets, locations to ping from and an PEEAAOO
// auth token that will be used to make a ping call.
func NewPingCollector(targets, locations, authToken string) *pingCollector {
	targetSlice := strings.Split(targets, ",")
	locationSlice := strings.Split(locations, ",")

	if targetSlice == nil {
		return nil
	}

	if locationSlice == nil {
		return nil
	}

	return &pingCollector{
		targets:   targetSlice,
		authToken: authToken,
		locations: locationSlice,
		statusMetric: prometheus.NewDesc("peeaao_response_status",
			"Status of Target (up/down/mixed - 0 down, 1 up, 2 mixed)",
			[]string{"target"}, nil,
		),
		responseCodeMetric: prometheus.NewDesc("peeaao_response_code",
			"Response Code of Ping",
			[]string{"target", "location"}, nil,
		),
		responseTimeMetric: prometheus.NewDesc("peeaao_response_time",
			"Response Time of Ping",
			[]string{"target", "location"}, nil,
		),
	}
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *pingCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.responseTimeMetric
	ch <- collector.responseCodeMetric
	ch <- collector.statusMetric
}

// PingerPayload represents the payload returned from PEEAAO's API call
type PingerPayload struct {
	Type      string              `json:"type"`
	Target    string              `json:"target"`
	Locations map[string][]Result `json:"locations"`
	Errors    []string            `json:"errors"`
	Status    string              `json:"status"`
}

type Result struct {
	Target     string `json:"target"`
	RunnerId   string `json:"runner_id"`
	Location   string `json:"location"`
	ResultInMs int64  `json:"result_in_ms"`
	Code       int64  `json:"code"`
}

//Collect implements required collect function for all prometheus collectors
func (collector *pingCollector) Collect(ch chan<- prometheus.Metric) {

	const (
		StatusDown float64 = iota
		StatusUp
		StatusMixed
	)

	var (
		wg      sync.WaitGroup
		payload PingerPayload

		authToken string = collector.authToken
	)

	for _, target := range collector.targets {

		wg.Add(1)
		go func() {
			defer wg.Done()
			locations := collector.locations
			body, statusCode, err := MakePing(target, locations, authToken)
			if err != nil {
				log.Printf("error pinging (%s) with locations (%s): %#v", target, locations, err)
				return
			}

			if statusCode != 200 {
				log.Printf("error pinging (%s) with locations (%s), expecting to get 200 status code but got %d", target, locations, statusCode)
				return

			}

			err = json.NewDecoder(body).Decode(&payload)
			if err != nil {
				fmt.Printf("error decoding: %#v", err)
				return
			}

			var metricValue float64
			var responseCodeValue float64
			var statusValue float64

			for _, location := range locations {

				// continue to the next loop if there's no result in there
				if len(payload.Locations[location]) == 0 {
					continue
				}

				metricValue = float64(payload.Locations[location][0].ResultInMs)
				responseCodeValue = float64(payload.Locations[location][0].Code)

				responseTimeMetric := prometheus.MustNewConstMetric(collector.responseTimeMetric, prometheus.GaugeValue, metricValue, target, location)
				responseTimeMetric = prometheus.NewMetricWithTimestamp(time.Now(), responseTimeMetric)
				responseCodeMetric := prometheus.MustNewConstMetric(collector.responseCodeMetric, prometheus.GaugeValue, responseCodeValue, target, location)
				responseCodeMetric = prometheus.NewMetricWithTimestamp(time.Now(), responseCodeMetric)
				ch <- responseTimeMetric
				ch <- responseCodeMetric
			}
			switch payload.Status {

			case "up":
				statusValue = StatusUp
			case "down":
				statusValue = StatusDown
			case "mixed":
				statusValue = StatusMixed
			}
			statusMetric := prometheus.MustNewConstMetric(collector.statusMetric, prometheus.GaugeValue, statusValue, target)
			statusMetric = prometheus.NewMetricWithTimestamp(time.Now(), statusMetric)
			ch <- statusMetric
		}()

	}
	wg.Wait()
}
