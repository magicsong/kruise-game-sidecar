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
	"encoding/json"
	"fmt"
	"reflect"
)

func ConvertJsonObjectToStruct(source interface{}, target interface{}) error {
	if source == nil {
		return fmt.Errorf("source is nil")
	}
	if target == nil {
		return fmt.Errorf("target is nil")
	}
	// source must be a map[string]
	// obj must be a pointer to a struct
	_, ok := source.(map[string]interface{})
	if !ok {
		return fmt.Errorf("source must be a map[string]interface{}")
	}
	if !isPointerToStruct(target) {
		return fmt.Errorf("target must be a pointer to a struct")
	}
	// then convert
	bytes, err := json.Marshal(source)
	if err != nil {
		return fmt.Errorf("failed to marshal source: %w", err)
	}
	err = json.Unmarshal(bytes, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}
	return nil
}

func isPointerToStruct(obj interface{}) bool {
	v := reflect.ValueOf(obj)
	return v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct
}
