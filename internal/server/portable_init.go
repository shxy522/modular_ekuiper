// Copyright 2022 EMQ Technologies Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build portable || !core
// +build portable !core

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"

	"github.com/lf-edge/ekuiper/internal/binder"
	"github.com/lf-edge/ekuiper/internal/conf"
	"github.com/lf-edge/ekuiper/internal/plugin"
	"github.com/lf-edge/ekuiper/internal/plugin/portable"
	"github.com/lf-edge/ekuiper/internal/plugin/portable/runtime"
	"github.com/lf-edge/ekuiper/pkg/errorx"
)

type portableStatusManager struct {
	sync.RWMutex
	status map[string]string
}

var psmManager *portableStatusManager

func init() {
	psmManager = &portableStatusManager{
		status: make(map[string]string),
	}
}

func (psm *portableStatusManager) StartInstall(pluginName string) {
	psm.Lock()
	defer psm.Unlock()
	psm.status[pluginName] = "installing"
}

func (psm *portableStatusManager) Installed(pluginName string) {
	psm.Lock()
	defer psm.Unlock()
	psm.status[pluginName] = "installed"
}

func (psm *portableStatusManager) GetPluginInstallStatus(name string) (string, bool) {
	psm.RLock()
	defer psm.RUnlock()
	s, ok := psm.status[name]
	if ok {
		return s, true
	}
	return "", false
}

var portableManager *portable.Manager

func init() {
	components["portable"] = portableComp{}
}

type portableComp struct{}

func (p portableComp) register() {
	var err error
	portableManager, err = portable.InitManager()
	if err != nil {
		panic(err)
	}
	entries = append(entries, binder.FactoryEntry{Name: "portable plugin", Factory: portableManager, Weight: 8})
}

func (p portableComp) rest(r *mux.Router) {
	r.HandleFunc("/plugins/portables", portablesHandler).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/plugins/portables/{name}", portableHandler).Methods(http.MethodGet, http.MethodDelete, http.MethodPut)
	r.HandleFunc("/plugins/portables/{name}/status", portableStatusHandler).Methods(http.MethodGet)
}

func portablesHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	switch r.Method {
	case http.MethodGet:
		content := portableManager.List()
		k := r.URL.Query().Get("kind")
		nc := make([]*portable.PluginInfo, 0)
		for _, c := range content {
			switch strings.ToLower(k) {
			case "source":
				if len(c.Sources) > 0 {
					nc = append(nc, c)
				}
			case "sink":
				if len(c.Sinks) > 0 {
					nc = append(nc, c)
				}
			case "function":
				if len(c.Functions) > 0 {
					nc = append(nc, c)
				}
			default:
				nc = append(nc, c)
				jsonResponse(content, w, logger)
			}
		}
		jsonResponse(nc, w, logger)
	case http.MethodPost:
		sd := plugin.NewPluginByType(plugin.PORTABLE)
		err := json.NewDecoder(r.Body).Decode(sd)
		// Problems decoding
		if err != nil {
			handleError(w, err, "Invalid body: Error decoding the portable plugin json", logger)
			return
		}
		conf.Log.Infof("recv install portable plugin %v request", sd.GetName())
		psmManager.StartInstall(sd.GetName())
		err = portableManager.Register(sd)
		if err != nil {
			conf.Log.Errorf("install portable plugin %v request err:%v", sd.GetName(), err)
			errMsg := fmt.Errorf("portable plugin create failed, err:%v", err.Error())
			handleError(w, errMsg, "", logger)
			return
		}
		psmManager.Installed(sd.GetName())
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("portable plugin %s is created", sd.GetName())))
	}
}

func portableStatusHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	name := vars["name"]
	switch r.Method {
	case http.MethodGet:
		status, err := runtime.GetPluginInsManager().GetPluginInsStatus(name)
		if err != nil {
			handleError(w, err, fmt.Sprintf("query portable plugin %s status error: %v", name, err), logger)
			return
		}
		_, foundInManager := portableManager.GetPluginInfo(name)
		installStatus, foundInstallStatus := psmManager.GetPluginInstallStatus(name)
		hasPluginRunningStatus := status != nil
		if hasPluginRunningStatus {
			jsonResponse(status, w, logger)
			return
		}
		switch {
		case !foundInstallStatus && !foundInManager:
			w.WriteHeader(http.StatusNotFound)
		case foundInstallStatus && !foundInManager:
			w.Write([]byte(installStatus))
			w.WriteHeader(http.StatusOK)
		case foundInManager:
			w.Write([]byte("installed"))
			w.WriteHeader(http.StatusOK)
		}
	}
}

func portableHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	name := vars["name"]
	switch r.Method {
	case http.MethodDelete:
		conf.Log.Infof("recv delete portable plugin %v request", name)
		err := portableManager.Delete(name)
		if err != nil {
			conf.Log.Errorf("delete portable plugin %v request err:%v", name, err)
			handleError(w, err, fmt.Sprintf("delete portable plugin %s error", name), logger)
			return
		}
		w.WriteHeader(http.StatusOK)
		result := fmt.Sprintf("portable plugin %s is deleted", name)
		w.Write([]byte(result))
	case http.MethodGet:
		j, ok := portableManager.GetPluginInfo(name)
		if !ok {
			handleError(w, errorx.NewWithCode(errorx.NOT_FOUND, "not found"), fmt.Sprintf("describe portable plugin %s error", name), logger)
			return
		}
		jsonResponse(j, w, logger)
	case http.MethodPut:
		sd := plugin.NewPluginByType(plugin.PORTABLE)
		err := json.NewDecoder(r.Body).Decode(sd)
		// Problems decoding
		if err != nil {
			handleError(w, err, "Invalid body: Error decoding the portable plugin json", logger)
			return
		}
		conf.Log.Infof("recv update portable plugin %v request", name)
		err = portableManager.Delete(name)
		if err != nil {
			conf.Log.Errorf("delete portable plugin %s error: %v", name, err)
		}
		err = portableManager.Register(sd)
		if err != nil {
			conf.Log.Errorf("update portable plugin %v request err:%v", name, err)
			handleError(w, err, "portable plugin update command error", logger)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("portable plugin %s is updated", sd.GetName())))
	}
}

func portablePluginsReset() {
	portableManager.UninstallAllPlugins()
}

func portablePluginExport() map[string]string {
	return portableManager.GetAllPlugins()
}

func portablePluginStatusExport() map[string]string {
	return portableManager.GetAllPlugins()
}

func portablePluginImport(plugins map[string]string) map[string]string {
	return portableManager.PluginImport(plugins)
}

func portablePluginPartialImport(plugins map[string]string) map[string]string {
	return portableManager.PluginPartialImport(plugins)
}
