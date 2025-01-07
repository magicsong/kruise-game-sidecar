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

package plugins

import (
	"github.com/magicsong/kidecar/api"
	"github.com/magicsong/kidecar/pkg/plugins/hot_update"
	httpprobe "github.com/magicsong/kidecar/pkg/plugins/http_probe"
)

var PluginRegistry = make(map[string]api.Plugin)

func RegisterPlugin(plugin api.Plugin) {
	if plugin.Name() == "" {
		panic("plugin name is empty")
	}
	PluginRegistry[plugin.Name()] = plugin
}

func init() {
	RegisterPlugin(httpprobe.NewPlugin())
	RegisterPlugin(hot_update.NewPlugin())
}
