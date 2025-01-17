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
	"github.com/agiledragon/gomonkey/v2"
	"github.com/go-logr/logr"
	"github.com/magicsong/kidecar/pkg/info"
	"github.com/magicsong/kidecar/pkg/store"
	"testing"
)

func Test_hotUpdate_setHotUpdateConfigWhenStart(t *testing.T) {
	patchGetCurrentPodInfo := gomonkey.ApplyFunc(info.GetCurrentPodInfo, func() (string, error) {
		return "pod1-default", nil
	})
	defer patchGetCurrentPodInfo.Reset()

	p := &store.PersistentConfig{
		Result: map[string]string{
			"v1.0":   "url1",
			"v3":     "url3",
			"v0.0.5": "url5",
			"v6.0.1": "url6",
			"v2.0":   "url2",
		},
	}

	patchGetPersistenceInfo := gomonkey.ApplyMethod(p, "GetPersistenceInfo", func() error {
		p.Result = map[string]string{
			"v1.0": "url1",
			"v2.0": "url2",
		}
		return nil
	})
	defer patchGetPersistenceInfo.Reset()

	h := &hotUpdate{
		result: &HotUpdateResult{},
	}
	patchStoreDataToConfigmap := gomonkey.ApplyMethod(h, "StoreDataToConfigmap", func() error {
		return nil
	})
	defer patchStoreDataToConfigmap.Reset()
	patchdownloadFileByUrl := gomonkey.ApplyMethod(h, "DownloadFileByUrl", func() error {
		return nil
	})
	defer patchdownloadFileByUrl.Reset()
	loadHotUpdateFileBySignal := gomonkey.ApplyMethod(h, "LoadHotUpdateFileBySignal", func() error {
		return nil
	})
	defer loadHotUpdateFileBySignal.Reset()
	loadHotUpdateFileByRequest := gomonkey.ApplyMethod(h, "LoadHotUpdateFileByRequest", func() error {
		return nil
	})
	defer loadHotUpdateFileByRequest.Reset()

	storeData := gomonkey.ApplyMethod(h, "StoreData", func() error {
		return nil
	})
	defer storeData.Reset()

	type fields struct {
		config         HotUpdateConfig
		StorageFactory store.StorageFactory
		status         *HotUpdateStatus
		result         *HotUpdateResult
		log            logr.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestCase1",
			fields: fields{
				config: HotUpdateConfig{},
				status: &HotUpdateStatus{},
				result: &HotUpdateResult{},
				log:    logr.Logger{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := h.SetHotUpdateConfigWhenStart(); (err != nil) != tt.wantErr {
				t.Errorf("SetHotUpdateConfigWhenStart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
