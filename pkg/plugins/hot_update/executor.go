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

package hot_update

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/magicsong/kidecar/pkg/store"
	"golang.org/x/mod/semver"

	"github.com/magicsong/kidecar/pkg/template"
)

func (h *hotUpdate) HotUpdateHandle(w http.ResponseWriter, r *http.Request) {

	err := h.DownloadHotUpdateFile(w, r)
	if err != nil {
		h.log.Error(err, "Failed to download file")
		return
	}

	h.log.Info("File downloaded and saved successfully")

	// According to the input Config, determine which way to trigger the update
	switch h.config.LoadPatchType {
	case LoadPatchTypeSignal:
		err := h.LoadHotUpdateFileBySignal()
		if err != nil {
			h.log.Error(err, "Failed to load hot update file by signal")
			h.result.Result = fmt.Sprintf("%s: Failed to load hot update file by signal: %s", h.result.Version, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "File downloaded and update successfully")
		h.result.Result = fmt.Sprintf("%s: Update success", h.result.Version)
	case LoadPatchTypeRequest:
		err := h.LoadHotUpdateFileByRequest()
		if err != nil {
			h.log.Error(err, "Failed to load hot update file by request")
			return
		}
	}

	err = h.StoreData()
	if err != nil {
		h.log.Error(err, "Failed to store data")
		return
	}

	err = h.StoreDataToConfigmap()
	if err != nil {
		h.log.Error(err, "failed to store data to configmap")
		return
	}
}

func (h *hotUpdate) DownloadHotUpdateFile(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "only POST requests are allowed")
		return fmt.Errorf("only POST requests are allowed")
	}

	// get version from request
	version := r.FormValue("version")
	if version == "" || !isValidVersion(version) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "please provide a valid version")
		return fmt.Errorf("please provide a valid version")
	}
	h.result.Version = version

	// get file url from request
	url := r.FormValue(URLKey)
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "please provide a valid URL")
		return fmt.Errorf("please provide a valid URL")
	}
	h.result.Url = url

	err := h.DownloadFileByUrl()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to download file: %v", err)
		return err
	}

	return nil
}

func (h *hotUpdate) DownloadFileByUrl() error {
	// Check if the hot update file  exists, if not, create it
	if _, err := os.Stat(FileDir); os.IsNotExist(err) {
		err := os.MkdirAll(FileDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	// get file name from url
	fileURL, err := http.NewRequest(http.MethodGet, h.result.Url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request for URL: %w", err)
	}
	filePath := filepath.Base(fileURL.URL.Path)

	// full save path
	savePath := filepath.Join(FileDir, filePath)

	// check if file exists, if exists, delete it
	if _, err := os.Stat(savePath); err == nil {
		err := os.Remove(savePath)
		if err != nil {
			return fmt.Errorf("failed to delete file: %w", err)
		}
	}

	// download file
	resp, err := http.Get(h.result.Url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	// save file
	out, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

func (h *hotUpdate) LoadHotUpdateFileBySignal() error {
	// get process list
	psout, err := exec.Command("ps", "aux").Output()
	if err != nil {
		return fmt.Errorf("failed to get process list, err: %v", err)
	}

	processName := h.config.Signal.ProcessName
	signal := h.config.Signal.SignalName
	processes := strings.Split(string(psout), "\n")
	for _, process := range processes {
		if strings.Contains(process, processName) {
			processInfos := strings.Fields(process)
			var pid string
			for _, processInfo := range processInfos {
				// Judge whether processInfo is a number
				if _, err := strconv.Atoi(processInfo); err != nil {
					continue
				}
				pid = processInfo
				cmd := exec.Command("kill", "-s", signal, pid)
				err := cmd.Run()
				if err != nil {
					return fmt.Errorf("failed to send signal to PID, signal: %v , pid: %v, processInfo: %v, err: %v", signal, pid, processInfos, err)
				}
				h.log.Info("Signal sent successfully", "signal", signal, "pid", pid, " processInfo: ", processInfos)
				return nil
			}
		}
	}
	return fmt.Errorf("process not found, processName: %v", processName)
}

func (h *hotUpdate) LoadHotUpdateFileByRequest() error {
	return nil
}

func (h *hotUpdate) StoreData() error {

	h.log.Info("store update result, ", "result: ", h.result.Result)
	err := template.ParseConfig(h.config.StorageConfig)
	if err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}

	return h.config.StorageConfig.StoreData(h.StorageFactory, h.result.Result)
}

func (h *hotUpdate) StoreDataToConfigmap() error {

	persistentResult := &store.PersistentConfig{
		Type: pluginName,
		Result: map[string]string{
			h.result.Version: h.result.Url,
		},
	}

	err := persistentResult.SetPersistenceInfo()
	if err != nil {
		return fmt.Errorf("failed to set hot update result to configmap")
	}

	return nil
}

func (h *hotUpdate) SetHotUpdateConfigWhenStart() error {

	persistentResult := &store.PersistentConfig{
		Type: pluginName,
	}
	err := persistentResult.GetPersistenceInfo()
	if err != nil {
		return fmt.Errorf("failed to GetPersistenceInfo of %v ", pluginName)
	}

	if len(persistentResult.Result) == 0 {
		h.result.Version = OriginVersion
		h.result.Url = OriginUrl
		err = h.StoreDataToConfigmap()
		if err != nil {
			return fmt.Errorf("failed to store data to configmap: %v", err)
		}

		h.log.Info("sidecar result has not the pod info")
		return nil
	}
	if len(persistentResult.Result) == 1 && persistentResult.Result[OriginVersion] == OriginUrl {
		h.result.Version = OriginVersion
		h.result.Url = OriginUrl
		err = h.StoreDataToConfigmap()
		if err != nil {
			return fmt.Errorf("failed to store data to configmap: %v", err)
		}

		h.log.Info("origin pod restart")
		return nil
	}

	version := ""
	url := ""
	for v, u := range persistentResult.Result {
		if semver.Compare(version, v) < 0 {
			version = v
			url = u
		}
	}

	h.result.Version = version
	h.result.Url = url

	// down load
	err = h.DownloadFileByUrl()
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}

	switch h.config.LoadPatchType {
	case LoadPatchTypeSignal:
		err := h.LoadHotUpdateFileBySignal()
		if err != nil {
			h.log.Error(err, "Failed to load hot update file by signal")
			h.result.Result = fmt.Sprintf("%s: Failed to load hot update file by signal: %s", h.result.Version, err)
			return fmt.Errorf("failed to load hot update file by signal")
		}
		h.result.Result = fmt.Sprintf("%s: Update success", h.result.Version)
	case LoadPatchTypeRequest:
		err := h.LoadHotUpdateFileByRequest()
		if err != nil {
			h.log.Error(err, "Failed to load hot update file by request")
			return err
		}
	}

	err = h.StoreData()
	if err != nil {
		return fmt.Errorf("failed to store data: %v", err)
	}

	err = h.StoreDataToConfigmap()
	if err != nil {
		return fmt.Errorf("failed to store data to configmap: %v", err)
	}
	return nil
}
