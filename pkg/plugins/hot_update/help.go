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

package hot_update

import (
	"fmt"
	"regexp"
)

const (
	// pluginName is the name of the plugin.
	pluginName = "hot_update"

	// URLKey is the key of the URL in the request body.
	URLKey = "url"

	// LoadPatchTypeSignal is the type of load patch that sends a signal to the main container.
	LoadPatchTypeSignal = "signal"
	// LoadPatchTypeRequest is the type of load patch that sends a request to the main container.
	LoadPatchTypeRequest = "request"

	// OriginVersion is the version of the original file.
	OriginVersion = "OriginVersion"
	// OriginUrl is the url of the original file.
	OriginUrl = "OriginUrl"

	// FileDir ...
	FileDir = "/app/downloads"
)

func validateConfig(hotUpdateConfig *HotUpdateConfig) error {
	if hotUpdateConfig.LoadPatchType != LoadPatchTypeRequest && hotUpdateConfig.LoadPatchType != LoadPatchTypeSignal {
		return fmt.Errorf("loadPatchType is empty")
	}
	if hotUpdateConfig.LoadPatchType == LoadPatchTypeSignal {
		if hotUpdateConfig.Signal.SignalName == "" {
			return fmt.Errorf("SignalName is empty")
		}
	}
	if hotUpdateConfig.LoadPatchType == LoadPatchTypeRequest {
		if hotUpdateConfig.Request.Address == "" {
			return fmt.Errorf("address is empty")
		}
		if hotUpdateConfig.Request.Port == 0 {
			return fmt.Errorf("port is empty")
		}
	}

	return nil
}

func isValidVersion(version string) bool {
	re := regexp.MustCompile(`^v\d+(?:\.\d+)*$`)
	return re.MatchString(version)
}
