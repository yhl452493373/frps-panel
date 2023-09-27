package controller

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	plugin "github.com/fatedier/frp/pkg/plugin/server"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func (c *HandleController) MakeHandlerFunc() gin.HandlerFunc {
	return func(context *gin.Context) {
		var response plugin.Response
		var err error

		request := plugin.Request{}
		if err := context.BindJSON(&request); err != nil {
			_ = context.AbortWithError(http.StatusBadRequest, err)
			return
		}

		jsonStr, err := json.Marshal(request.Content)
		if err != nil {
			_ = context.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if request.Op == "Login" {
			content := plugin.LoginContent{}
			err = json.Unmarshal(jsonStr, &content)
			response = c.HandleLogin(&content)
		} else if request.Op == "NewProxy" {
			content := plugin.NewProxyContent{}
			err = json.Unmarshal(jsonStr, &content)
			response = c.HandleNewProxy(&content)
		} else if request.Op == "Ping" {
			content := plugin.PingContent{}
			err = json.Unmarshal(jsonStr, &content)
			response = c.HandlePing(&content)
		} else if request.Op == "NewWorkConn" {
			content := plugin.NewWorkConnContent{}
			err = json.Unmarshal(jsonStr, &content)
			response = c.HandleNewWorkConn(&content)
		} else if request.Op == "NewUserConn" {
			content := plugin.NewUserConnContent{}
			err = json.Unmarshal(jsonStr, &content)
			response = c.HandleNewUserConn(&content)
		}

		if err != nil {
			log.Printf("handle %s error: %v", context.Request.URL.Path, err)
			var e *HTTPError
			switch {
			case errors.As(err, &e):
				context.JSON(e.Code, &Response{Msg: e.Err.Error()})
			default:
				context.JSON(http.StatusInternalServerError, &Response{Msg: err.Error()})
			}
			return
		} else {
			resStr, _ := json.Marshal(response)
			log.Printf("handle:%v , result: %v", request.Op, string(resStr))
		}

		context.JSON(http.StatusOK, response)
	}
}

func (c *HandleController) MakeLoginFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		if context.Request.Method == "GET" {
			if c.LoginAuth("", "", context) {
				if context.Request.RequestURI == LoginUrl {
					context.Redirect(http.StatusTemporaryRedirect, LoginSuccessUrl)
				}
				return
			}
			context.HTML(http.StatusOK, "login.html", gin.H{
				"version":             c.Version,
				"FrpsPanel":           ginI18n.MustGetMessage(context, "Frps Panel"),
				"Username":            ginI18n.MustGetMessage(context, "Username"),
				"Password":            ginI18n.MustGetMessage(context, "Password"),
				"Login":               ginI18n.MustGetMessage(context, "Login"),
				"PleaseInputUsername": ginI18n.MustGetMessage(context, "Please input username"),
				"PleaseInputPassword": ginI18n.MustGetMessage(context, "Please input password"),
			})
		} else if context.Request.Method == "POST" {
			username := context.PostForm("username")
			password := context.PostForm("password")
			if c.LoginAuth(username, password, context) {
				context.JSON(http.StatusOK, gin.H{
					"success": true,
					"message": ginI18n.MustGetMessage(context, "Login success"),
				})
			} else {
				context.JSON(http.StatusOK, gin.H{
					"success": false,
					"message": ginI18n.MustGetMessage(context, "Username or password incorrect"),
				})
			}
		}
	}
}

func (c *HandleController) MakeLogoutFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		ClearAuth(context)
		context.Redirect(http.StatusTemporaryRedirect, LogoutSuccessUrl)
	}
}

func (c *HandleController) MakeIndexFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{
			"version":                      c.Version,
			"showExit":                     trimString(c.CommonInfo.AdminUser) != "" && trimString(c.CommonInfo.AdminPwd) != "",
			"FrpsPanel":                    ginI18n.MustGetMessage(context, "Frps Panel"),
			"User":                         ginI18n.MustGetMessage(context, "User"),
			"Token":                        ginI18n.MustGetMessage(context, "Token"),
			"Notes":                        ginI18n.MustGetMessage(context, "Notes"),
			"Search":                       ginI18n.MustGetMessage(context, "Search"),
			"Reset":                        ginI18n.MustGetMessage(context, "Reset"),
			"NewUser":                      ginI18n.MustGetMessage(context, "New user"),
			"RemoveUser":                   ginI18n.MustGetMessage(context, "Remove user"),
			"DisableUser":                  ginI18n.MustGetMessage(context, "Disable user"),
			"EnableUser":                   ginI18n.MustGetMessage(context, "Enable user"),
			"Remove":                       ginI18n.MustGetMessage(context, "Remove"),
			"Enable":                       ginI18n.MustGetMessage(context, "Enable"),
			"Disable":                      ginI18n.MustGetMessage(context, "Disable"),
			"PleaseInputUserAccount":       ginI18n.MustGetMessage(context, "Please input user account"),
			"PleaseInputUserToken":         ginI18n.MustGetMessage(context, "Please input user token"),
			"PleaseInputUserNotes":         ginI18n.MustGetMessage(context, "Please input user notes"),
			"AllowedPorts":                 ginI18n.MustGetMessage(context, "Allowed ports"),
			"PleaseInputAllowedPorts":      ginI18n.MustGetMessage(context, "Please input allowed ports"),
			"AllowedDomains":               ginI18n.MustGetMessage(context, "Allowed domains"),
			"PleaseInputAllowedDomains":    ginI18n.MustGetMessage(context, "Please input allowed domains"),
			"AllowedSubdomains":            ginI18n.MustGetMessage(context, "Allowed subdomains"),
			"PleaseInputAllowedSubdomains": ginI18n.MustGetMessage(context, "Please input allowed subdomains"),
			"NotLimit":                     ginI18n.MustGetMessage(context, "Not limit"),
			"None":                         ginI18n.MustGetMessage(context, "None"),
			"ServerInfo":                   ginI18n.MustGetMessage(context, "Server Info"),
			"Users":                        ginI18n.MustGetMessage(context, "Users"),
			"Proxies":                      ginI18n.MustGetMessage(context, "Proxies"),
			"TrafficStatistics":            ginI18n.MustGetMessage(context, "Traffic Statistics"),
			"Name":                         ginI18n.MustGetMessage(context, "Name"),
			"Type":                         ginI18n.MustGetMessage(context, "Type"),
			"Domains":                      ginI18n.MustGetMessage(context, "Domains"),
			"SubDomain":                    ginI18n.MustGetMessage(context, "SubDomain"),
			"Locations":                    ginI18n.MustGetMessage(context, "Locations"),
			"HostRewrite":                  ginI18n.MustGetMessage(context, "HostRewrite"),
			"Encryption":                   ginI18n.MustGetMessage(context, "Encryption"),
			"Compression":                  ginI18n.MustGetMessage(context, "Compression"),
			"Addr":                         ginI18n.MustGetMessage(context, "Addr"),
			"LastStart":                    ginI18n.MustGetMessage(context, "Last Start"),
			"LastClose":                    ginI18n.MustGetMessage(context, "Last Close"),
			"Version":                      ginI18n.MustGetMessage(context, "Version"),
			"BindPort":                     ginI18n.MustGetMessage(context, "Bind Port"),
			"KCPBindPort":                  ginI18n.MustGetMessage(context, "KCP Bind Port"),
			"QUICBindPort":                 ginI18n.MustGetMessage(context, "QUIC Bind Port"),
			"HTTPPort":                     ginI18n.MustGetMessage(context, "HTTP Port"),
			"HTTPSPort":                    ginI18n.MustGetMessage(context, "HTTPS Port"),
			"TCPMUXPort":                   ginI18n.MustGetMessage(context, "TCPMUX Port"),
			"SubdomainHost":                ginI18n.MustGetMessage(context, "Subdomain Host"),
			"MaxPoolCount":                 ginI18n.MustGetMessage(context, "Max Pool Count"),
			"MaxPortsPerClient":            ginI18n.MustGetMessage(context, "Max Ports Per Client"),
			"HeartBeatTimeout":             ginI18n.MustGetMessage(context, "Heart Beat Timeout"),
			"AllowPorts":                   ginI18n.MustGetMessage(context, "Allow Ports"),
			"TLSOnly":                      ginI18n.MustGetMessage(context, "TLS Only"),
			"CurrentConnections":           ginI18n.MustGetMessage(context, "Current Connections"),
			"ClientCounts":                 ginI18n.MustGetMessage(context, "Client Counts"),
			"ProxyCounts":                  ginI18n.MustGetMessage(context, "Proxy Counts"),
		})
	}
}

func (c *HandleController) MakeLangFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"User":                  ginI18n.MustGetMessage(context, "User"),
			"Token":                 ginI18n.MustGetMessage(context, "Token"),
			"Notes":                 ginI18n.MustGetMessage(context, "Notes"),
			"Status":                ginI18n.MustGetMessage(context, "Status"),
			"Operation":             ginI18n.MustGetMessage(context, "Operation"),
			"Enable":                ginI18n.MustGetMessage(context, "Enable"),
			"Disable":               ginI18n.MustGetMessage(context, "Disable"),
			"NewUser":               ginI18n.MustGetMessage(context, "New user"),
			"Confirm":               ginI18n.MustGetMessage(context, "Confirm"),
			"Cancel":                ginI18n.MustGetMessage(context, "Cancel"),
			"RemoveUser":            ginI18n.MustGetMessage(context, "Remove user"),
			"DisableUser":           ginI18n.MustGetMessage(context, "Disable user"),
			"ConfirmRemoveUser":     ginI18n.MustGetMessage(context, "Confirm to remove user"),
			"ConfirmDisableUser":    ginI18n.MustGetMessage(context, "Confirm to disable user"),
			"TakeTimeMakeEffective": ginI18n.MustGetMessage(context, "will take sometime to make effective"),
			"ConfirmEnableUser":     ginI18n.MustGetMessage(context, "Confirm to enable user"),
			"OperateSuccess":        ginI18n.MustGetMessage(context, "Operate success"),
			"OperateError":          ginI18n.MustGetMessage(context, "Operate error"),
			"OperateFailed":         ginI18n.MustGetMessage(context, "Operate failed"),
			"UserExist":             ginI18n.MustGetMessage(context, "User exist"),
			"UserNotExist":          ginI18n.MustGetMessage(context, "User not exist"),
			"UserFormatError":       ginI18n.MustGetMessage(context, "User format error"),
			"TokenFormatError":      ginI18n.MustGetMessage(context, "Token format error"),
			"ShouldCheckUser":       ginI18n.MustGetMessage(context, "Please check at least one user"),
			"OperationConfirm":      ginI18n.MustGetMessage(context, "Operation confirm"),
			"EmptyData":             ginI18n.MustGetMessage(context, "Empty data"),
			"NotLimit":              ginI18n.MustGetMessage(context, "Not limit"),
			"AllowedPorts":          ginI18n.MustGetMessage(context, "Allowed ports"),
			"AllowedDomains":        ginI18n.MustGetMessage(context, "Allowed domains"),
			"AllowedSubdomains":     ginI18n.MustGetMessage(context, "Allowed subdomains"),
			"PortsInvalid":          ginI18n.MustGetMessage(context, "Ports is invalid"),
			"DomainsInvalid":        ginI18n.MustGetMessage(context, "Domains is invalid"),
			"SubdomainsInvalid":     ginI18n.MustGetMessage(context, "Subdomains is invalid"),
			"CommentInvalid":        ginI18n.MustGetMessage(context, "Comment is invalid"),
			"ParamError":            ginI18n.MustGetMessage(context, "Param error"),
			"OtherError":            ginI18n.MustGetMessage(context, "Other error"),
			"Name":                  ginI18n.MustGetMessage(context, "Name"),
			"Port":                  ginI18n.MustGetMessage(context, "Port"),
			"Connections":           ginI18n.MustGetMessage(context, "Connections"),
			"TrafficIn":             ginI18n.MustGetMessage(context, "Traffic In"),
			"TrafficOut":            ginI18n.MustGetMessage(context, "Traffic Out"),
			"ClientVersion":         ginI18n.MustGetMessage(context, "Client Version"),
			"TrafficStatistics":     ginI18n.MustGetMessage(context, "Traffic Statistics"),
			"online":                ginI18n.MustGetMessage(context, "online"),
			"offline":               ginI18n.MustGetMessage(context, "offline"),
			"true":                  ginI18n.MustGetMessage(context, "true"),
			"false":                 ginI18n.MustGetMessage(context, "false"),
			"NetworkTraffic":        ginI18n.MustGetMessage(context, "Network Traffic"),
			"today":                 ginI18n.MustGetMessage(context, "today"),
			"now":                   ginI18n.MustGetMessage(context, "now"),
			"Proxies":               ginI18n.MustGetMessage(context, "Proxies"),
			"NotSet":                ginI18n.MustGetMessage(context, "Not Set"),
			"Proxy":                 ginI18n.MustGetMessage(context, "Proxy"),
			"TokenInvalid":          ginI18n.MustGetMessage(context, "Token invalid"),
		})
	}
}

func (c *HandleController) MakeQueryTokensFunc() func(context *gin.Context) {
	return func(context *gin.Context) {

		search := TokenSearch{}
		search.Limit = 0

		err := context.BindQuery(&search)
		if err != nil {
			return
		}

		var tokenList []TokenInfo
		for _, tokenInfo := range c.Tokens {
			tokenList = append(tokenList, tokenInfo)
		}
		sort.Slice(tokenList, func(i, j int) bool {
			return strings.Compare(tokenList[i].User, tokenList[j].User) < 0
		})

		var filtered []TokenInfo
		for _, tokenInfo := range tokenList {
			if filter(tokenInfo, search.TokenInfo) {
				filtered = append(filtered, tokenInfo)
			}
		}
		if filtered == nil {
			filtered = []TokenInfo{}
		}

		count := len(filtered)
		if search.Limit > 0 {
			start := max((search.Page-1)*search.Limit, 0)
			end := min(search.Page*search.Limit, len(filtered))
			filtered = filtered[start:end]
		}

		context.JSON(http.StatusOK, &TokenResponse{
			Code:  0,
			Msg:   "query Tokens success",
			Count: count,
			Data:  filtered,
		})
	}
}

func (c *HandleController) MakeAddTokenFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		info := TokenInfo{
			Enable: true,
		}
		response := OperationResponse{
			Success: true,
			Code:    Success,
			Message: "user add success",
		}
		err := context.BindJSON(&info)
		if err != nil {
			response.Success = false
			response.Code = ParamError
			response.Message = fmt.Sprintf("user add failed, param error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		result := c.verifyToken(info, TOKEN_ADD)

		if !result.Success {
			context.JSON(http.StatusOK, &result)
			return
		}

		info.Comment = cleanString(info.Comment)
		info.Ports = cleanPorts(info.Ports)
		info.Domains = cleanStrings(info.Domains)
		info.Subdomains = cleanStrings(info.Subdomains)

		c.Tokens[info.User] = info

		err = c.saveToken()
		if err != nil {
			response.Success = false
			response.Code = SaveError
			response.Message = fmt.Sprintf("user add failed, error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		context.JSON(0, &response)
	}
}

func (c *HandleController) MakeUpdateTokensFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		response := OperationResponse{
			Success: true,
			Code:    Success,
			Message: "user update success",
		}
		update := TokenUpdate{}
		err := context.BindJSON(&update)
		if err != nil {
			response.Success = false
			response.Code = ParamError
			response.Message = fmt.Sprintf("update failed, param error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		before := update.Before
		after := update.After

		if before.User != after.User {
			response.Success = false
			response.Code = ParamError
			response.Message = fmt.Sprintf("update failed, user should be same : before -> %v, after -> %v", before.User, after.User)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		result := c.verifyToken(after, TOKEN_UPDATE)

		if !result.Success {
			context.JSON(http.StatusOK, &result)
			return
		}

		after.Comment = cleanString(after.Comment)
		after.Ports = cleanPorts(after.Ports)
		after.Domains = cleanStrings(after.Domains)
		after.Subdomains = cleanStrings(after.Subdomains)

		c.Tokens[after.User] = after

		err = c.saveToken()
		if err != nil {
			response.Success = false
			response.Code = SaveError
			response.Message = fmt.Sprintf("user update failed, error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		context.JSON(http.StatusOK, &response)
	}
}

func (c *HandleController) MakeRemoveTokensFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		response := OperationResponse{
			Success: true,
			Code:    Success,
			Message: "user remove success",
		}
		remove := TokenRemove{}
		err := context.BindJSON(&remove)
		if err != nil {
			response.Success = false
			response.Code = ParamError
			response.Message = fmt.Sprintf("user remove failed, param error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		for _, user := range remove.Users {
			result := c.verifyToken(user, TOKEN_REMOVE)

			if !result.Success {
				context.JSON(http.StatusOK, &result)
				return
			}
		}

		for _, user := range remove.Users {
			delete(c.Tokens, user.User)
		}

		err = c.saveToken()
		if err != nil {
			response.Success = false
			response.Code = SaveError
			response.Message = fmt.Sprintf("user update failed, error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		context.JSON(http.StatusOK, &response)
	}
}

func (c *HandleController) MakeDisableTokensFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		response := OperationResponse{
			Success: true,
			Code:    Success,
			Message: "remove success",
		}
		disable := TokenDisable{}
		err := context.BindJSON(&disable)
		if err != nil {
			response.Success = false
			response.Code = ParamError
			response.Message = fmt.Sprintf("disable failed, param error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		for _, user := range disable.Users {
			result := c.verifyToken(user, TOKEN_DISABLE)

			if !result.Success {
				context.JSON(http.StatusOK, &result)
				return
			}
		}

		for _, user := range disable.Users {
			token := c.Tokens[user.User]
			token.Enable = false
			c.Tokens[user.User] = token
		}

		err = c.saveToken()

		if err != nil {
			response.Success = false
			response.Code = SaveError
			response.Message = fmt.Sprintf("disable failed, error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		context.JSON(http.StatusOK, &response)
	}
}

func (c *HandleController) MakeEnableTokensFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		response := OperationResponse{
			Success: true,
			Code:    Success,
			Message: "remove success",
		}
		enable := TokenEnable{}
		err := context.BindJSON(&enable)
		if err != nil {
			response.Success = false
			response.Code = ParamError
			response.Message = fmt.Sprintf("enable failed, param error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		for _, user := range enable.Users {
			result := c.verifyToken(user, TOKEN_ENABLE)

			if !result.Success {
				context.JSON(http.StatusOK, &result)
				return
			}
		}

		for _, user := range enable.Users {
			token := c.Tokens[user.User]
			token.Enable = true
			c.Tokens[user.User] = token
		}

		err = c.saveToken()

		if err != nil {
			log.Printf("enable failed, error : %v", err)
			response.Success = false
			response.Code = SaveError
			response.Message = "enable failed"
			context.JSON(http.StatusOK, &response)
			return
		}

		context.JSON(http.StatusOK, &response)
	}
}

func (c *HandleController) MakeProxyFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		var client *http.Client
		var protocol string

		if c.CommonInfo.DashboardTls {
			client = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			}
			protocol = "https://"
		} else {
			client = http.DefaultClient
			protocol = "http://"
		}

		res := ProxyResponse{}
		host := c.CommonInfo.DashboardAddr
		port := c.CommonInfo.DashboardPort

		host, _ = strings.CutPrefix(host, protocol)

		requestUrl := protocol + host + ":" + strconv.Itoa(port) + context.Param("serverApi")
		request, _ := http.NewRequest("GET", requestUrl, nil)
		username := c.CommonInfo.DashboardUser
		password := c.CommonInfo.DashboardPwd
		if trimString(username) != "" && trimString(password) != "" {
			request.SetBasicAuth(username, password)
			log.Printf("Proxy to %s", requestUrl)
		}

		response, err := client.Do(request)

		if err != nil {
			res.Code = FrpServerError
			res.Success = false
			res.Message = err.Error()
			log.Print(err)
			context.JSON(http.StatusOK, &res)
			return
		}

		res.Code = response.StatusCode
		body, err := io.ReadAll(response.Body)

		if err != nil {
			res.Success = false
			res.Message = err.Error()
		} else {
			if res.Code == http.StatusOK {
				res.Success = true
				res.Data = string(body)
				res.Message = fmt.Sprintf("Proxy to %s success", requestUrl)
			} else {
				res.Success = false
				if res.Code == http.StatusNotFound {
					res.Message = fmt.Sprintf("Proxy to %s error: url not found", requestUrl)
				} else {
					res.Message = fmt.Sprintf("Proxy to %s error: %s", requestUrl, string(body))
				}
			}
		}
		log.Printf(res.Message)
		context.JSON(http.StatusOK, &res)
	}
}
