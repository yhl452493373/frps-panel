package main

type Config struct {
	Common CommonInfo
	Tokens []TokenInfo
}

type CommonInfo struct {
	PluginAddr    string `toml:"plugin_addr" commented:"true"`
	PluginPort    int    `toml:"plugin_port"`
	AdminUser     string `toml:"admin_user"`
	AdminPwd      string `toml:"admin_pwd"`
	AdminKeepTime int    `toml:"admin_keep_time"`
	TlsMode       bool   `toml:"tls_mode"`
	TlsCertFile   string `toml:"tls_cert_file"`
	TlsKeyFile    string `toml:"tls_key_file"`
	DashboardAddr string `toml:"dashboard_addr"`
	DashboardPort int    `toml:"dashboard_port"`
	DashboardUser string `toml:"dashboard_user"`
	DashboardPwd  string `toml:"dashboard_pwd"`
}

type TokenInfo struct {
	User       string   `toml:"user" json:"user" form:"user"`
	Token      string   `toml:"token" json:"token" form:"token"`
	Comment    string   `toml:"comment" json:"comment" form:"comment"`
	Ports      []any    `toml:"ports" json:"ports" from:"ports"`
	Domains    []string `toml:"domains" json:"domains" from:"domains"`
	Subdomains []string `toml:"subdomains" json:"subdomains" from:"subdomains"`
	Status     bool     `toml:"status" json:"status" form:"status"`
}
