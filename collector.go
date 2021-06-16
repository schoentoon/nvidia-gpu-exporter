package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var metrics = make(map[string]*prometheus.Desc)
var labels = []string{"uuid", "name"}
var driverDesc = prometheus.NewDesc(
	"gpu_driver",
	"The version of the installed NVIDIA display driver. This is an alphanumeric string.",
	[]string{"driver"}, nil,
)

// in here we're initializing the map with all the metrics that we're interested in
// the key of the map will be the value that will be requested from `nvidia-smi`
// for a full list of the possible values, have a look at the output of `nvidia-smi --help-query-gpu`
// to add a new metric simply consult that output and look up the metric that you're intrested in
// afterwards simply copy one of the prometheus.NewDesc functions and adjust the variables accordingly
// do make sure to adjust the metric name (first variable for prometheus.NewDesc) to something unique
// otherwise you will likely get a runtime failure. it is also important to leave the labels variable intact as is
// as this is responsible for putting gpu name and uuid in the eventual prometheus metrics.
func init() {
	metrics["temperature.gpu"] = prometheus.NewDesc(
		"gpu_temperature",
		"Core GPU temperature. in degrees C.",
		labels, nil,
	)
	metrics["temperature.memory"] = prometheus.NewDesc(
		"gpu_memory_temperature",
		" HBM memory temperature. in degrees C.",
		labels, nil,
	)
	metrics["utilization.gpu"] = prometheus.NewDesc(
		"gpu_utilization",
		"Percent of time over the past sample period during which one or more kernels was executing on the GPU.",
		labels, nil,
	)
	metrics["memory.total"] = prometheus.NewDesc(
		"gpu_memory_total",
		"Total installed GPU memory. In MiB.",
		labels, nil,
	)
	metrics["memory.used"] = prometheus.NewDesc(
		"gpu_memory_used",
		"Total memory allocated by active contexts. In MiB.",
		labels, nil,
	)
	metrics["memory.free"] = prometheus.NewDesc(
		"gpu_memory_free",
		"Total free memory. In MiB.",
		labels, nil,
	)
	metrics["fan.speed"] = prometheus.NewDesc(
		"gpu_fan_speed",
		"The fan speed value is the percent of the product's maximum noise tolerance fan speed that the device's fan is currently intended to run at.",
		labels, nil,
	)
	metrics["power.draw"] = prometheus.NewDesc(
		"gpu_power_draw",
		"The last measured power draw for the entire board, in watts. Only available if power management is supported. This reading is accurate to within +/- 5 watts.",
		labels, nil,
	)
	metrics["clocks.current.graphics"] = prometheus.NewDesc(
		"gpu_graphics_clock_speed",
		"Current frequency of graphics (shader) clock. In megahertz",
		labels, nil,
	)
	metrics["clocks.current.sm"] = prometheus.NewDesc(
		"gpu_sm_clock_speed",
		"Current frequency of SM (Streaming Multiprocessor) clock. In megahertz",
		labels, nil,
	)
	metrics["clocks.current.memory"] = prometheus.NewDesc(
		"gpu_memory_clock_speed",
		"Current frequency of memory clock. In megahertz",
		labels, nil,
	)
	metrics["clocks.current.video"] = prometheus.NewDesc(
		"gpu_video_clock_speed",
		"Current frequency of video encoder/decoder clock. In megahertz",
		labels, nil,
	)
}

type GpuCollector struct {
}

func NewGpuCollector() *GpuCollector {
	return &GpuCollector{}
}

func (cc GpuCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range metrics {
		ch <- desc
	}
	ch <- driverDesc
}

func (cc GpuCollector) Collect(ch chan<- prometheus.Metric) {
	cc.driverVersion(ch)

	query := []string{"name", "uuid"}
	for key := range metrics {
		query = append(query, key)
	}

	out, err := exec.Command(
		"nvidia-smi",
		fmt.Sprintf("--query-gpu=%s", strings.Join(query, ",")),
		"--format=csv,noheader,nounits").Output()

	if err != nil {
		logrus.Errorf("%s\n", err)
		return
	}

	csvReader := csv.NewReader(bytes.NewReader(out))
	csvReader.TrimLeadingSpace = true
	records, err := csvReader.ReadAll()

	if err != nil {
		logrus.Errorf("%s\n", err)
		return
	}

	for _, row := range records {
		cc.handleRow(row, query, ch)
	}
}

func (cc GpuCollector) driverVersion(ch chan<- prometheus.Metric) {
	out, err := exec.Command(
		"nvidia-smi",
		"--query-gpu=driver_version",
		"--format=csv,noheader,nounits").Output()

	if err != nil {
		logrus.Errorf("%s\n", err)
		return
	}

	csvReader := csv.NewReader(bytes.NewReader(out))
	csvReader.TrimLeadingSpace = true
	records, err := csvReader.ReadAll()

	if err != nil {
		logrus.Errorf("%s\n", err)
		return
	}

	for _, row := range records {
		ch <- prometheus.MustNewConstMetric(
			driverDesc,
			prometheus.GaugeValue,
			1,
			row[0],
		)
		// we return here as we might get multiple results in the case of multiple gpu's
		// but the result is going to be the same.. and prometheus would give a runtime error
		// if you return a second value with the same label combination.
		return
	}
}

func (cc GpuCollector) handleRow(row, query []string, ch chan<- prometheus.Metric) {
	name := row[0]
	uuid := row[1]

	// we do a [2:] here to cut off the first 2 items that we already extracted (name and uuid)
	for index, rawValue := range row[2:] {
		// check for `N/A` which they return if something is not supported etc..
		// we do a contains instead of a direct check as it somehow can't be consist
		// and sometimes returns with [] around it and sometimes doesn't...
		if strings.Contains(rawValue, "N/A") {
			continue
		}

		// label will be the according nvidia metric name, therefore the key in the global metrics map.
		label := query[index+2] // offset by 2 as we cut off the name and uuid part
		desc := metrics[label]
		value, err := strconv.ParseFloat(rawValue, 64)
		if err != nil {
			logrus.Warnf("error with %s with value %v: %v\n", label, rawValue, err)
		} else {
			ch <- prometheus.MustNewConstMetric(
				desc,
				prometheus.GaugeValue,
				value,
				uuid, name,
			)
		}
	}
}
