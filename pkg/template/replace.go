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
	"os"
	"regexp"

	corev1 "k8s.io/api/core/v1"
)

// Expression format:
// ${SELF:VAR_NAME}: Indicates the environment variable of the sidecar itself.
// ${POD:VAR_NAME}: Indicates the environment variable of the Pod.

const (
	pattern = `\$\{(SELF|POD):([^}]+)\}`
)

func ReplaceValue(value string, container *corev1.Container) (string, error) {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(value)
	if matches == nil {
		return value, nil
	}

	envType := matches[1]
	envName := matches[2]

	var envValue string
	var found bool
	if envType == "SELF" {
		envValue, found = os.LookupEnv(envName)
	} else if envType == "POD" {
		// Search from the environment variables of the container
		for _, envVar := range container.Env {
			if envVar.Name == envName {
				envValue = envVar.Value
				found = true
				break
			}
		}
	} else {
		return "", fmt.Errorf("unknown environment variable type: %s", envType)
	}

	if !found {
		return "", fmt.Errorf("environment variable %s not found", envName)
	}

	return envValue, nil
}
