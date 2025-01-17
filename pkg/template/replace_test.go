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

package template

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"os"
	"testing"
)

func TestReplaceValue(t *testing.T) {
	tests := []struct {
		name          string
		value         string
		container     *corev1.Container
		expectedValue string
		expectedError error
	}{
		{
			name:          "NoMatch",
			value:         "hello world",
			container:     &corev1.Container{},
			expectedValue: "hello world",
			expectedError: nil,
		},
		{
			name:  "SelfEnvFound",
			value: "${SELF:ENV_VAR}",
			container: &corev1.Container{
				Env: []corev1.EnvVar{
					{Name: "OTHER_VAR", Value: "other_value"},
				},
			},
			expectedValue: "value_of_env_var",
			expectedError: nil,
		},
		{
			name:  "PodEnvFound",
			value: "${POD:ENV_VAR}",
			container: &corev1.Container{
				Env: []corev1.EnvVar{
					{Name: "ENV_VAR", Value: "pod_value"},
				},
			},
			expectedValue: "pod_value",
			expectedError: nil,
		},
		{
			name:  "SelfEnvNotFound",
			value: "${SELF:NOT_FOUND_VAR}",
			container: &corev1.Container{
				Env: []corev1.EnvVar{
					{Name: "OTHER_VAR", Value: "other_value"},
				},
			},
			expectedValue: "",
			expectedError: fmt.Errorf("environment variable NOT_FOUND_VAR not found"),
		},
		{
			name:  "PodEnvNotFound",
			value: "${POD:NOT_FOUND_VAR}",
			container: &corev1.Container{
				Env: []corev1.EnvVar{
					{Name: "OTHER_VAR", Value: "other_value"},
				},
			},
			expectedValue: "",
			expectedError: fmt.Errorf("environment variable NOT_FOUND_VAR not found"),
		},
		{
			name:  "UnknownEnvType",
			value: "${UNKNOWN:ENV_VAR}",
			container: &corev1.Container{
				Env: []corev1.EnvVar{
					{Name: "OTHER_VAR", Value: "other_value"},
				},
			},
			expectedValue: "${UNKNOWN:ENV_VAR}",
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("ENV_VAR", "value_of_env_var")
			value, err := ReplaceValue(tc.value, tc.container)
			if err != nil && tc.expectedError == nil {
				t.Errorf("Expected no error, but got: %v", err)
			}
			if err == nil && tc.expectedError != nil {
				t.Errorf("Expected error: %v, but got none", tc.expectedError)
			}
			if value != tc.expectedValue {
				t.Errorf("Expected value: %s, but got: %s", tc.expectedValue, value)
			}
		})
	}
}
