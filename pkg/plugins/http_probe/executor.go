/*
Copyright 2024

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

package httpprobe

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"time"

	"github.com/magicsong/kidecar/pkg/store"
	"github.com/magicsong/kidecar/pkg/template"
)

// Executor holds the HTTP client and provides methods for probing
type Executor struct {
	client *http.Client
	store.StorageFactory
}

// NewExecutor creates a new Prober with the provided timeout
func NewExecutor(timeout int, factory store.StorageFactory) *Executor {
	return &Executor{
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		StorageFactory: factory,
	}
}

// Probe performs the HTTP request based on the provided configuration
func (p *Executor) Probe(config EndpointConfig) error {
	req, err := http.NewRequest(config.Method, config.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	// Perform the request
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	// Check expected status code
	if resp.StatusCode != config.ExpectedStatusCode {
		return fmt.Errorf("unexpected status code: got %v, expected %v, body: %s", resp.StatusCode, config.ExpectedStatusCode, string(body))
	}

	// Extract data
	data, err := p.extractData(body, config.JSONPathConfig)
	if err != nil {
		return fmt.Errorf("failed to extract data: %v", err)
	}
	// Store data
	if err := p.storeData(data.(string), &config.StorageConfig); err != nil {
		return fmt.Errorf("failed to store data: %v", err)
	}
	return nil
}

func (p *Executor) extractData(data []byte, extractorConfig *store.JSONPathConfig) (interface{}, error) {
	if extractorConfig != nil {
		return getDataFromJsonText(string(data), extractorConfig.JSONPath)
	}
	return string(data), nil
}

func (p *Executor) storeData(data string, storeConfig *store.StorageConfig) error {
	err := template.ParseConfig(storeConfig)
	if err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}
	return storeConfig.StoreData(p.StorageFactory, data)
}

func getDataFromJsonText(json, path string) (interface{}, error) {
	if !gjson.Valid(json) {
		return nil, fmt.Errorf("invalid json")
	}
	value := gjson.Get(json, path)
	if !value.Exists() {
		return nil, fmt.Errorf("path not found")
	}
	return value.Value(), nil
}
