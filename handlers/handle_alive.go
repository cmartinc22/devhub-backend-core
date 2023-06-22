package handlers

import (
	"net/http"

	"github.com/pedidosya/peya-go/server"
)

type ServiceInfo struct {
	Name string      `json:"name" yaml:"name"`
	Info *EngineInfo `json:"info" yaml:"info"`
}

type EngineInfo struct {
	Environment string `json:"env,omitempty" yaml:"env,omitempty"`
	Version     string `json:"version,omitempty" yaml:"version,omitempty"`
}

func HandleAlive(w http.ResponseWriter, r *http.Request, serviceName string, info EngineInfo) {
	// Do not want to debug Alive
	server.OK(w, r, &ServiceInfo{
		Name: "devhub",
		Info: &info,
	})
}
