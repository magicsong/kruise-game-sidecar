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

package store

import (
	"fmt"

	"github.com/magicsong/kidecar/api"
)

type StorageFactory interface {
	GetStorage(storageType StorageType) (Storage, error)
}

type defaultStorageFactory struct {
	storageMap map[StorageType]Storage
	manager    api.SidecarManager
}

func NewStorageFactory(mgr api.SidecarManager) StorageFactory {
	f := &defaultStorageFactory{
		storageMap: make(map[StorageType]Storage),
	}
	f.storageMap[StorageTypeInKube] = &inKube{}
	f.storageMap[StorageTypeHTTPMetric] = &promMetric{}
	f.manager = mgr
	return f
}

func (f *defaultStorageFactory) GetStorage(storageType StorageType) (Storage, error) {
	s := f.storageMap[storageType]
	if s == nil {
		return nil, fmt.Errorf("storage type %s not found", storageType)
	}
	if !s.IsInitialized() {
		if err := s.SetupWithManager(f.manager); err != nil {
			return nil, fmt.Errorf("failed to setup storage: %w", err)
		}
	}
	return s, nil
}
