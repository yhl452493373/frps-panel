package controller

import (
	"regexp"
)

const (
	Success          = 0
	ParamError       = 1
	UserExist        = 2
	SaveError        = 3
	UserFormatError  = 4
	TokenFormatError = 5
	FrpServerError   = 6

	SessionName      = "GOSESSION"
	AuthName         = "_PANEL_AUTH"
	LoginUrl         = "/login"
	LoginSuccessUrl  = "/"
	LogoutUrl        = "/logout"
	LogoutSuccessUrl = "/login"
)

var (
	UserFormatReg    = regexp.MustCompile("^\\w+$")
	TokenFormatReg   = regexp.MustCompile("^[\\w!@#$%^&*()]+$")
	TrimAllSpaceReg  = regexp.MustCompile("[\\n\\t\\r\\s]")
	TrimBreakLineReg = regexp.MustCompile("[\\n\\t\\r]")
)

type Response struct {
	Msg string `json:"msg"`
}

type HTTPError struct {
	Code int
	Err  error
}

type CommonInfo struct {
	PluginAddr    string
	PluginPort    int
	User          string
	Pwd           string
	KeepTime      int
	DashboardTLS  bool
	DashboardAddr string
	DashboardPort int
	DashboardUser string
	DashboardPwd  string
}

type TokenInfo struct {
	User       string `json:"user" form:"user"`
	Token      string `json:"token" form:"token"`
	Comment    string `json:"comment" form:"comment"`
	Ports      string `json:"ports" from:"ports"`
	Domains    string `json:"domains" from:"domains"`
	Subdomains string `json:"subdomains" from:"subdomains"`
	Status     bool   `json:"status" form:"status"`
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
