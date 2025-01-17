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

package template

import (
	"fmt"
	"reflect"

	"github.com/magicsong/kidecar/pkg/info"
	corev1 "k8s.io/api/core/v1"
)

func expressionReplaceValue(value string, pod *corev1.Pod) (string, error) {
	// Add your own template parsing logic here
	return ReplaceValue(value, &pod.Spec.Containers[0])
}

// ParseConfig Parse the fields in the configuration structure recursively
func ParseConfig(config interface{}) error {
	pod, err := info.GetCurrentPod()
	if err != nil {
		return fmt.Errorf("failed to get current pod: %w", err)
	}
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() {
		return nil
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanInterface() {
			continue
		}

		// Check if there is a `parse:"true"` tag
		if tagValue, ok := fieldType.Tag.Lookup("parse"); ok && tagValue == "true" {
			// Parse the field value
			if field.Kind() == reflect.String {
				parsedValue, err := expressionReplaceValue(field.String(), pod)
				if err != nil {
					return fmt.Errorf("failed to parse field %s: %w", fieldType.Name, err)
				}
				field.SetString(parsedValue)
			}
		}

		// If it is a structure or a pointer to a structure, handle it recursively
		if field.Kind() == reflect.Struct {
			if err := ParseConfig(field.Addr().Interface()); err != nil {
				return err
			}
		} else if field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.Struct {
			if err := ParseConfig(field.Interface()); err != nil {
				return err
			}
		}
	}

	return nil
}
