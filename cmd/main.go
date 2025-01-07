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
package main

import (
	"context"
	"os"

	"github.com/magicsong/kidecar/pkg/assembler"
	"github.com/magicsong/kidecar/pkg/info"
	"github.com/magicsong/kidecar/pkg/manager"
	flag "github.com/spf13/pflag"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "/opt/kidecar/config.yaml", "config file path")
}

func main() {
	logf.SetLogger(zap.New())
	log := logf.Log.WithName("manager-examples")
	flag.Parse()
	sidecar := assembler.NewSidecar()
	if err := sidecar.LoadConfig(configPath); err != nil {
		log.Error(err, "failed to load config")
		os.Exit(1)
	}
	mgr, err := manager.NewManager()
	if err != nil {
		log.Error(err, "failed to create manager")
		panic(err)
	}
	info.SetGlobalKubeInterface(mgr)
	sidecar.SetupWithManager(mgr)
	// add plugins
	if err := sidecar.InitPlugins(); err != nil {
		panic(err)
	}
	ctx := context.TODO()
	if err := sidecar.Start(ctx); err != nil {
		panic(err)
	}
}
