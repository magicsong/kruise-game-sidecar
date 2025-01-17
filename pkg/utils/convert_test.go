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

package utils

import (
	"testing"
)

type TestStruct struct {
	Field1 string
	Field2 int
}

func TestConvertJsonObjectToStruct(t *testing.T) {
	tests := []struct {
		name    string
		source  interface{}
		target  interface{}
		wantErr bool
	}{
		{
			name:    "Success",
			source:  map[string]interface{}{"Field1": "value1", "Field2": 2},
			target:  &TestStruct{},
			wantErr: false,
		},
		{
			name:    "SourceNil",
			source:  nil,
			target:  &TestStruct{},
			wantErr: true,
		},
		{
			name:    "TargetNil",
			source:  map[string]interface{}{"Field1": "value1", "Field2": 2},
			target:  nil,
			wantErr: true,
		},
		{
			name:    "SourceInvalid",
			source:  123,
			target:  &TestStruct{},
			wantErr: true,
		},
		{
			name:    "TargetInvalid",
			source:  map[string]interface{}{"Field1": "value1", "Field2": 2},
			target:  TestStruct{},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ConvertJsonObjectToStruct(tc.source, tc.target)
			if (err != nil) != tc.wantErr {
				t.Errorf("Expected error: %v, but got: %v", tc.wantErr, err)
			}
		})
	}
}
