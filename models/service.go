package models

type AuthResult struct {
	IsValid  bool
	Enabled  bool
	Identity *string
}

type ServiceInfo struct {
	Name string      `json:"name" yaml:"name"`
	Info *EngineInfo `json:"info" yaml:"info"`
}

type EngineInfo struct {
	Environment string `json:"env,omitempty" yaml:"env,omitempty"`
	Version     string `json:"version,omitempty" yaml:"version,omitempty"`
}
