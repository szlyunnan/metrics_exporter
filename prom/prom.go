package prom

import (
	"github.com/prometheus/client_golang/prometheus"
	"metrics_exporter/model"
	"strconv"
)

type MetricsCollector struct {
	data           []model.QueryData
	signMetricDesc *prometheus.Desc
}

func NewMetricsCollector(data []model.QueryData) *MetricsCollector {
	return &MetricsCollector{
		data: data,
		signMetricDesc: prometheus.NewDesc("sign_Metrics",
			"Shows whether sign Metrics",
			// aax_code si_build
			// aax_name si_name
			// aax_pkg si_pkg
			// status si_message
			[]string{"aax_code", "aax_name", "aax_pkg", "status"},
			nil,
		),
	}
}

func (mc *MetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- mc.signMetricDesc
}

func (mc *MetricsCollector) Collect(ch chan<- prometheus.Metric) {
	for _, v := range mc.data {
		ch <- prometheus.MustNewConstMetric(
			mc.signMetricDesc,
			prometheus.GaugeValue,
			float64(v.SiCode),
			// labels
			strconv.FormatInt(v.SiBuild, 10),
			"aax-sign",
			v.SiPkg,
			v.SiMsg,
		)
	}
}
