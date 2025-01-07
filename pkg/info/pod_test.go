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

package info

import (
	"os"
	"testing"

	"k8s.io/apimachinery/pkg/types"
)

func TestGetCurrentPodNamespaceAndName(t *testing.T) {
	tests := []struct {
		name           string
		namespaceEnv   string
		nameEnv        string
		wantErr        bool
		expectedResult *types.NamespacedName
	}{
		{
			name:         "BothSet",
			namespaceEnv: "namespace1",
			nameEnv:      "pod1",
			wantErr:      false,
			expectedResult: &types.NamespacedName{
				Namespace: "namespace1",
				Name:      "pod1",
			},
		},
		{
			name:         "NamespaceNotSet",
			namespaceEnv: "",
			nameEnv:      "pod1",
			wantErr:      true,
		},
		{
			name:         "NameNotSet",
			namespaceEnv: "namespace1",
			nameEnv:      "",
			wantErr:      true,
		},
		{
			name:         "BothNotSet",
			namespaceEnv: "",
			nameEnv:      "",
			wantErr:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("POD_NAMESPACE", tc.namespaceEnv)
			os.Setenv("POD_NAME", tc.nameEnv)

			result, err := GetCurrentPodNamespaceAndName()
			if (err != nil) != tc.wantErr {
				t.Errorf("Expected error: %v, but got: %v", tc.wantErr, err)
			}
			if err == nil && !compareNamespacedNames(result, tc.expectedResult) {
				t.Errorf("Expected result: %v, but got: %v", tc.expectedResult, result)
			}
		})
	}
}

func compareNamespacedNames(a, b *types.NamespacedName) bool {
	return a.Namespace == b.Namespace && a.Name == b.Name
}

func TestGetCurrentPodInfo(t *testing.T) {
	tests := []struct {
		name           string
		namespaceEnv   string
		nameEnv        string
		wantErr        bool
		expectedResult string
	}{
		{
			name:           "BothSet",
			namespaceEnv:   "namespace1",
			nameEnv:        "pod1",
			wantErr:        false,
			expectedResult: "namespace1-pod1",
		},
		{
			name:         "NamespaceNotSet",
			namespaceEnv: "",
			nameEnv:      "pod1",
			wantErr:      true,
		},
		{
			name:         "NameNotSet",
			namespaceEnv: "namespace1",
			nameEnv:      "",
			wantErr:      true,
		},
		{
			name:         "BothNotSet",
			namespaceEnv: "",
			nameEnv:      "",
			wantErr:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("POD_NAMESPACE", tc.namespaceEnv)
			os.Setenv("POD_NAME", tc.nameEnv)

			result, err := GetCurrentPodInfo()
			if (err != nil) != tc.wantErr {
				t.Errorf("Expected error: %v, but got: %v", tc.wantErr, err)
			}
			if err == nil && result != tc.expectedResult {
				t.Errorf("Expected result: %s, but got: %s", tc.expectedResult, result)
			}
		})
	}
}
