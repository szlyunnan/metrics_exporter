package utils

import (
	"github.com/prometheus/client_golang/prometheus"
	"metrics_exporter/model"
	"strconv"
)

type MetricsCollector struct {
	data           []model.QueryData
	signMetricDesc *prometheus.Desc
}

//{aax_code="430",aax_name="aax-sign",aax_pkg="com.aax.exchange.1.3.0.96e56842.ipa",status="Response verify OK"}
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
