package api

import (
	"context"

	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Plugin ...
type Plugin interface {
	Name() string
	Init(config interface{}, mgr SidecarManager) error
	Start(ctx context.Context, errCh chan<- error)
	Stop(ctx context.Context) error
	Version() string
	Status() (*PluginStatus, error)
	// must return a pointer of your config
	GetConfigType() interface{}
}

// PluginConfig ...
type PluginConfig struct {
	Name      string      `json:"name"`
	Config    interface{} `json:"config"`
	BootOrder int         `json:"bootOrder"`
}

// SidecarConfig ...
type SidecarConfig struct {
	Plugins           []PluginConfig    `json:"plugins"`           // plugins and  configurations
	RestartPolicy     string            `json:"restartPolicy"`     // Restart policy
	Resources         map[string]string `json:"resources"`         // The resources required by Sidecar
	SidecarStartOrder string            `json:"sidecarStartOrder"` // The startup sequence of Sidecar, is it after or before the main container
}

// PluginStatus =
type PluginStatus struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Running     bool     `json:"running"`
	LastChecked string   `json:"lastChecked"` //  YYYY-MM-DD HH:MM:SS
	Health      string   `json:"health"`      //  "Healthy", "Unhealthy"
	Infos       []string `json:"infos"`
}

// Sidecar ...
type Sidecar interface {
	InitPlugins() error
	RemovePlugin(pluginName string) error
	GetVersion() string
	PluginStatus(pluginName string) (*PluginStatus, error)
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	SetupWithManager(mgr SidecarManager) error
	LoadConfig(path string) error
}

type SidecarManager interface {
	ctrl.Manager
	DBManager
	kubernetes.Interface
}

type DBManager interface {
}
