package controller

import (
	"regexp"
)

const (
	Success int = iota
	ParamError
	UserExist
	UserNotExist
	SaveError
	UserFormatError
	TokenFormatError
	CommentFormatError
	PortsFormatError
	DomainsFormatError
	SubdomainsFormatError
	FrpServerError
)

const (
	TOKEN_ADD int = iota
	TOKEN_UPDATE
	TOKEN_REMOVE
	TOKEN_ENABLE
	TOKEN_DISABLE
)

const (
	SessionName      = "GOSESSION"
	AuthName         = "_PANEL_AUTH"
	LoginUrl         = "/login"
	LoginSuccessUrl  = "/"
	LogoutUrl        = "/logout"
	LogoutSuccessUrl = "/login"
)

var (
	userFormat        = regexp.MustCompile("^\\w+$")
	tokenFormat       = regexp.MustCompile("^[\\w!@#$%^&*()]+$")
	commentFormat     = regexp.MustCompile("[\\n\\t\\r]")
	portsFormatSingle = regexp.MustCompile("^\\s*\\d{1,5}\\s*$")
	portsFormatRange  = regexp.MustCompile("^\\s*\\d{1,5}\\s*-\\s*\\d{1,5}\\s*$")
	domainFormat      = regexp.MustCompile("^([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*\\.)+[a-zA-Z]{2,}$")
	subdomainFormat   = regexp.MustCompile("^[a-zA-z0-9-]{1,20}$")
	trimAllSpace      = regexp.MustCompile("[\\n\\t\\r\\s]")
)

type Response struct {
	Msg string `json:"msg"`
}

type HTTPError struct {
	Code int
	Err  error
}

type Common struct {
	Common CommonInfo
}

type CommonInfo struct {
	PluginAddr    string `toml:"plugin_addr"`
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
	DashboardTls  bool
}

type Tokens struct {
	Tokens map[string]TokenInfo `toml:"tokens"`
}

type TokenInfo struct {
	User       string   `toml:"user" json:"user" form:"user"`
	Token      string   `toml:"token" json:"token" form:"token"`
	Comment    string   `toml:"comment" json:"comment" form:"comment"`
	Ports      []string `toml:"ports" json:"ports" from:"ports"`
	Domains    []string `toml:"domains" json:"domains" from:"domains"`
	Subdomains []string `toml:"subdomains" json:"subdomains" from:"subdomains"`
	Status     bool     `toml:"status" json:"status" form:"status"`
}

type TokenResponse struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Count int         `json:"count"`
	Data  []TokenInfo `json:"data"`
}

type OperationResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ProxyResponse struct {
	OperationResponse
	Data string `json:"data"`
}

type TokenSearch struct {
	TokenInfo
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

type TokenUpdate struct {
	Before TokenInfo `json:"before"`
	After  TokenInfo `json:"after"`
}

type TokenRemove struct {
	Users []TokenInfo `json:"users"`
}

type TokenDisable struct {
	TokenRemove
}

type TokenEnable struct {
	TokenDisable
}

func (e *HTTPError) Error() string {
	return e.Err.Error()
}
