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

package manager

import (
	"fmt"

	"github.com/magicsong/kidecar/api"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type sidecarManager struct {
	ctrl.Manager
	api.DBManager
	kubernetes.Interface
}

func NewManager() (api.SidecarManager, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to get kubeconfig: %w", err)
	}

	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		return nil, fmt.Errorf("unable to create manager: %w", err)
	}
	kube, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create kubernetes client: %w", err)
	}
	return sidecarManager{Manager: mgr, Interface: kube}, nil
}
