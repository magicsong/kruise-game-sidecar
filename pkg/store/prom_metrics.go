/*
Copyright 2024  .

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package store

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/magicsong/kidecar/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type promMetric struct {
	registry  *prometheus.Registry
	metrics   map[string]prometheus.Gauge
	metricsMu sync.Mutex
}

// IsInitialized implements Storage.
func (p *promMetric) IsInitialized() bool {
	return p.registry != nil
}

// SetupWithManager implements Storage.
func (p *promMetric) SetupWithManager(mgr api.SidecarManager) error {
	reg := prometheus.NewRegistry()
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	// 启动HTTP服务器
	fmt.Println("Starting server at :8080")
	p.registry = reg
	go http.ListenAndServe(":8080", nil)
	return nil
}

// Store implements Storage.
func (p *promMetric) Store(data string, config interface{}) error {
	myconfig, ok := config.(*HTTPMetricConfig)
	if !ok {
		return fmt.Errorf("bad config of httpMetricConfig")
	}
	f, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return fmt.Errorf("bad data of httpMetricConfig, err: %w", err)
	}
	p.getOrCreateGauge(myconfig).Set(f)
	return nil
}

// getOrCreateGauge ...
func (p *promMetric) getOrCreateGauge(config *HTTPMetricConfig) prometheus.Gauge {
	p.metricsMu.Lock()
	defer p.metricsMu.Unlock()

	if gauge, exists := p.metrics[config.MetricName]; exists {
		return gauge
	}

	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: config.MetricName,
		Help: "Automatically generated metric from collected data",
	})
	p.registry.MustRegister(gauge)
	p.metrics[config.MetricName] = gauge
	return gauge
}
