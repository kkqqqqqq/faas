// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package metrics

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricOptions to be used by web handlers
type MetricOptions struct {
	GatewayFunctionInvocation        *prometheus.CounterVec
	GatewayFunctionsHistogram        *prometheus.HistogramVec
	GatewayFunctionInvocationStarted *prometheus.CounterVec

	ServiceReplicasGauge *prometheus.GaugeVec

	ServiceMetrics *ServiceMetricOptions

	//cold start
	GatewayFunctionColdStartMetrics *ServiceMetricOptions
}

// ServiceMetricOptions provides RED metrics
type ServiceMetricOptions struct {
	Histogram *prometheus.HistogramVec
	Counter   *prometheus.CounterVec
}

// Synchronize to make sure MustRegister only called once
var once = sync.Once{}

// RegisterExporter registers with Prometheus for tracking
func RegisterExporter(exporter *Exporter) {
	once.Do(func() {
		prometheus.MustRegister(exporter)
	})
}

// PrometheusHandler Bootstraps prometheus for metrics collection
func PrometheusHandler() http.Handler {
	return promhttp.Handler()
}

// BuildMetricsOptions builds metrics for tracking functions in the API gateway
func BuildMetricsOptions() MetricOptions {
	gatewayFunctionsHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "gateway_functions_seconds",
		Help: "Function time taken",
	}, []string{"function_name", "code"})

	gatewayFunctionInvocation := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gateway",
			Subsystem: "function",
			Name:      "invocation_total",
			Help:      "Function metrics",
		},
		[]string{"function_name", "code"},
	)

	serviceReplicas := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "gateway",
			Name:      "service_count",
			Help:      "Current count of replicas for function",
		},
		[]string{"function_name"},
	)

	// For automatic monitoring and alerting (RED method)
	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "Seconds spent serving HTTP requests.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "path", "status"})

	// Can be used Kubernetes HPA v2
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "The total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	gatewayFunctionInvocationStarted := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gateway",
			Subsystem: "function",
			Name:      "invocation_started",
			Help:      "The total number of function HTTP requests started.",
		},
		[]string{"function_name"},
	)

	//cold start 
	gatewayFunctionColdStartMetricOptions := &ServiceMetricOptions{
		Counter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "gateway",
				Subsystem: "function",
				Name:      "cold_start_total",
				Help:      "The total number of function cold start.",
			},
			[]string{"function_name"},
		),
		Histogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "gateway",
			Subsystem: "function",
			Name:      "cold_start_duration_seconds",
			Help:      "Seconds spent function cold start.",
			Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		}, []string{"function_name"}),
	}



	serviceMetricOptions := &ServiceMetricOptions{
		Counter:   counter,
		Histogram: histogram,
	}

	metricsOptions := MetricOptions{
		GatewayFunctionsHistogram:        gatewayFunctionsHistogram,
		GatewayFunctionInvocation:        gatewayFunctionInvocation,
		ServiceReplicasGauge:             serviceReplicas,
		ServiceMetrics:                   serviceMetricOptions,
		GatewayFunctionInvocationStarted: gatewayFunctionInvocationStarted,

		//cold start 
		GatewayFunctionColdStartMetrics:  gatewayFunctionColdStartMetricOptions,
	}

	return metricsOptions
}
