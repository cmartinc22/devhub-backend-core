package models

type DocsConfiguration struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}
type DiagnosticsConfig struct {
	LogLevel    string `json:"log_level" yaml:"log_level"`
	EnablePprof bool   `json:"enable_pprof" yaml:"enable_pprof"`
}

type DbConfiguration struct {
	Server        string `json:"server" yaml:"server"`
	Port          int    `json:"port" yaml:"port"`
	Database      string `json:"database" yaml:"database"`
	User          string `json:"user" yaml:"user"`
	Password      string `json:"pwd" yaml:"pwd"`
	ConnectionStr string `json:"connStr" yaml:"connStr"`
}

type STSConfiguration struct {
	Enabled    bool   `json:"enabled" yaml:"enabled"`
	URL        string `json:"url" yaml:"url"`
	PrivateKey string `json:"private_key" yaml:"private_key"`
	Timeout    int    `json:"timeout" yaml:"timeout"`
	ClientID   string `json:"client_id" yaml:"client_id"`
	KeyID      string `json:"key_id" yaml:"key_id"`
}

type CFConfiguration struct {
	Enabled             bool   `json:"enabled" yaml:"enabled"`
	AUD                 string `json:"AUD" yaml:"AUD"`
	ServicePublicDomain string `json:"servicePublicDomain" yaml:"servicePublicDomain"`
}
